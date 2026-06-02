package inapp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const ProcessedSet = "inapp:processed_ids"
const BroadcastChannelPrefix = "inapp:broadcast:" // New constant

func StartConsumer(ctx context.Context, rdb *redis.Client, stream, group, consumer string) {
	for {
		if ctx.Err() != nil {
			log.Printf("inapp consumer: stopping due to context cancellation: %v", ctx.Err())
			return
		}
		entries, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{stream, ">"},
			Count:    10,
			Block:    5 * time.Second,
		}).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			if ctx.Err() != nil {
				log.Printf("inapp consumer: stopping due to context cancellation: %v", ctx.Err())
				return
			}
			log.Printf("inapp consumer read err: %v", err)
			time.Sleep(time.Second)
			continue
		}
		for _, st := range entries {
			for _, msg := range st.Messages {
				raw, _ := msg.Values["payload"].(string)
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(raw), &payload); err != nil {
					log.Printf("inapp: invalid payload: %v", err)
					_ = rdb.XAck(ctx, stream, group, msg.ID)
					_, _ = rdb.XDel(ctx, stream, msg.ID).Result()
					continue
				}

				id, _ := payload["id"].(string)
				recipient, _ := payload["recipient"].(string)
				applicationID, _ := payload["application_id"].(string)

				// Idempotency check
				already, _ := rdb.SIsMember(ctx, ProcessedSet, id).Result()
				if already {
					_ = rdb.XAck(ctx, stream, group, msg.ID)
					_, _ = rdb.XDel(ctx, stream, msg.ID).Result()
					continue
				}

				// Format: inapp:broadcast:app_id:user_id
				broadcastChannel := fmt.Sprintf("%s%s:%s", BroadcastChannelPrefix, applicationID, recipient)
				if err := rdb.Publish(ctx, broadcastChannel, raw).Err(); err != nil {
					log.Printf("failed to publish broadcast for %s: %v", recipient, err)
					continue
				}

				log.Printf("✓ Published notification %s for app %s user %s", id, applicationID, recipient)

				// Mark processed and ack
				_, _ = rdb.SAdd(ctx, ProcessedSet, id).Result()
				_ = rdb.XAck(ctx, stream, group, msg.ID)
				_, _ = rdb.XDel(ctx, stream, msg.ID).Result()
			}
		}
	}
}
