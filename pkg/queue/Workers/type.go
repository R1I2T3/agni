package workers

import (
	"context"
	"sync"
	"time"
)

type NotificationWorker struct {
	WorkerID   int
	QueueName  string
	MaxRetries int
	RetryDelay time.Duration
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
}

type WorkerPool struct {
	workers    []*NotificationWorker
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	numWorkers int
	queueName  string
}
