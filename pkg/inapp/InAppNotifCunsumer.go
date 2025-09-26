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
					continue
				}
				//		id, _ := payload["id"].(string)
				recipient, _ := payload["recipient"].(string)

				// delete the entry from the stream to free space immediately
				_, _ = rdb.XDel(ctx, StreamName, msg.ID).Result()

				// deliver to connected clients
				fmt.Println("dilvering to %s", recipient)
				DefaultHub.BroadcastToUser(recipient, payload)

			}
		}
	}
}
