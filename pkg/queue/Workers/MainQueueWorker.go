package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/r1i2t3/agni/pkg/db"
	"github.com/r1i2t3/agni/pkg/notification"
	"github.com/r1i2t3/agni/pkg/notification/channels/email"
	"github.com/r1i2t3/agni/pkg/queue"
)

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int, queueName string) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		workers:    make([]*NotificationWorker, 0, numWorkers),
		ctx:        ctx,
		cancel:     cancel,
		numWorkers: numWorkers,
		queueName:  queueName,
	}
}

func (wp *WorkerPool) Start() {
	log.Printf("🚀 Starting %d notification workers for queue: %s", wp.numWorkers, wp.queueName)

	for i := 0; i < wp.numWorkers; i++ {
		worker := &NotificationWorker{
			WorkerID:   i + 1,
			QueueName:  wp.queueName,
			MaxRetries: 3,
			RetryDelay: time.Second * 5,
		}

		worker.ctx, worker.cancel = context.WithCancel(wp.ctx)
		worker.wg = &wp.wg

		wp.workers = append(wp.workers, worker)

		// Start worker in goroutine
		wp.wg.Add(1)
		go worker.Start()
	}

	log.Printf("✅ All %d workers started successfully", wp.numWorkers)
}

// Stop gracefully stops all workers
func (wp *WorkerPool) Stop() {
	log.Println("🛑 Stopping worker pool...")

	// Cancel all workers
	wp.cancel()

	// Wait for all workers to finish
	wp.wg.Wait()

	log.Println("✅ All workers stopped")
}

// Start begins the worker's processing loop
func (w *NotificationWorker) Start() {
	defer w.wg.Done()

	log.Printf("🔄 Worker %d started", w.WorkerID)

	for {
		select {
		case <-w.ctx.Done():
			log.Printf("🛑 Worker %d stopping...", w.WorkerID)
			return
		default:
			if err := w.processNext(); err != nil {
				if err.Error() != "queue empty: no notifications available" {
					log.Printf("❌ Worker %d error: %v", w.WorkerID, err)
				}
				// Brief pause if no work available
				time.Sleep(time.Second * 2)
			}
		}
	}
}

// processNext dequeues and processes the next notification
func (w *NotificationWorker) processNext() error {
	// Dequeue notification with 5 second timeout
	fmt.Println("Attempting to dequeue notification...")
	queuedNotif, err := queue.DequeueNotification(w.QueueName, time.Second*5)
	if err != nil {
		fmt.Println("Error dequeuing notification:", err)
		return err
	}

	log.Printf("📝 Worker %d processing notification %s for %s, subject: %s",
		w.WorkerID, queuedNotif.ID, queuedNotif.Recipient, queuedNotif.Subject)

	// Process the notification
	return w.processNotification(queuedNotif)
}
func (w *NotificationWorker) processNotification(notif *queue.QueuedNotification) error {
	log.Printf("🔔 Worker %d processing notification %s", w.WorkerID, notif.ID)
	var err error
	var sentNotification *notification.Notification
	switch notif.Channel {
	case "email":
		// Process email notification
		log.Printf("📧 Worker %d sending email to %s with subject %s", w.WorkerID, notif.Recipient, notif.Subject)

		sentNotification, err = email.ProcessEmailNotifications(notif)
	case "sms":
		// Process SMS notification
		log.Printf("📱 Worker %d sending SMS to %s", w.WorkerID, notif.Recipient)
		// TODO: Add SMS sending logic here
	case "push":
		// Process push notification
		log.Printf("📲 Worker %d sending push notification to %s", w.WorkerID, notif.Recipient)
		// TODO: Add push notification logic here
	default:
		log.Printf("⚠️ Worker %d unknown notification channel: %s", w.WorkerID, notif.Channel)
		return fmt.Errorf("unknown notification channel: %s", notif.Channel)
	}

	if err == nil {
		database := db.GetMySQLDB()
		err := database.Create(&sentNotification).Error
		if err != nil {
			log.Printf("Error saving to DB: %v", err)
		}
	}

	return nil
}
