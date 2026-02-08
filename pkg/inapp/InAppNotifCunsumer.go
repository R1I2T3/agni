package inapp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const StreamName = "stream:inapp"
const GroupName = "inapp-group"
const DLQ = "stream:inapp:dlq"
const ProcessedSet = "inapp:processed_ids"
const BroadcastChannelPrefix = "inapp:broadcast:" // New constant

func StartConsumer(ctx context.Context, rdb *redis.Client, group, consumer string) {
	for {
		entries, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{StreamName, ">"},
			Count:    10,
			Block:    5 * time.Second,
		}).Result()
		if err != nil {
			if err == redis.Nil {
				continue
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
					_ = rdb.XAck(ctx, StreamName, group, msg.ID)
					_, _ = rdb.XDel(ctx, StreamName, msg.ID).Result()
					continue
				}

				id, _ := payload["id"].(string)
				recipient, _ := payload["recipient"].(string)

				// Idempotency check
				already, _ := rdb.SIsMember(ctx, ProcessedSet, id).Result()
				if already {
					_ = rdb.XAck(ctx, StreamName, group, msg.ID)
					_, _ = rdb.XDel(ctx, StreamName, msg.ID).Result()
					continue
				}

				// CHANGED: Publish to Redis Pub/Sub instead of calling hub directly
				// This broadcasts to ALL containers so the one with the WS connection can deliver
				broadcastChannel := fmt.Sprintf("%s%s", BroadcastChannelPrefix, recipient)
				if err := rdb.Publish(ctx, broadcastChannel, raw).Err(); err != nil {
					log.Printf("failed to publish broadcast for %s: %v", recipient, err)
					continue
				}

				log.Printf("✓ Published notification %s for %s to pub/sub", id, recipient)

				// Mark processed and ack
				_, _ = rdb.SAdd(ctx, ProcessedSet, id).Result()
				_ = rdb.XAck(ctx, StreamName, group, msg.ID)
				_, _ = rdb.XDel(ctx, StreamName, msg.ID).Result()
			}
		}
	}
}
