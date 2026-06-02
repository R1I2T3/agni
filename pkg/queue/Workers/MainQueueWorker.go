package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/r1i2t3/agni/pkg/db"
	"github.com/r1i2t3/agni/pkg/notification"
	"github.com/r1i2t3/agni/pkg/notification/channels/email"
	inapp "github.com/r1i2t3/agni/pkg/notification/channels/in-app"
	"github.com/r1i2t3/agni/pkg/notification/channels/sms"
	"github.com/r1i2t3/agni/pkg/notification/channels/webpush"
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
	queuedNotif, err := queue.DequeueNotification(w.QueueName, time.Second*5)
	if err != nil {
		return err
	}

	log.Printf("📝 Worker %d processing notification %s for %s", w.WorkerID, queuedNotif.ID, queuedNotif.Recipient)

	err = w.processNotification(queuedNotif)
	if err != nil {
		log.Printf("❌ Worker %d error sending notification %s: %v", w.WorkerID, queuedNotif.ID, err)

		// Retry Logic
		if queuedNotif.Attempts < w.MaxRetries {
			queuedNotif.Attempts++
			log.Printf("🔄 Rescheduling notification %s for retry (Attempt %d/%d) in %v",
				queuedNotif.ID, queuedNotif.Attempts, w.MaxRetries, w.RetryDelay)

			_, retryErr := queue.DelayReEnqueueNotification(queuedNotif, w.RetryDelay)
			if retryErr != nil {
				log.Printf("❌ Failed to reschedule notification: %v", retryErr)
			}
		} else {
			log.Printf("💀 Notification %s reached max retries (%d). Marking as failed.",
				queuedNotif.ID, w.MaxRetries)
			// Optional: save to database as permanently failed
		}
	}
	return nil
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
		sentNotification, err = sms.ProcessSMSNotifications(notif)
	case "webpush":
		// Process web push notification
		log.Printf("📲 Worker %d sending web push notification to %s", w.WorkerID, notif.Recipient)
		sentNotification, err = webpush.ProcessWebPushNotifications(notif)

	case "InApp":
		// Process in App notification
		log.Printf("📲 Worker %d sending InApp notification to %s", w.WorkerID, notif.Recipient)
		sentNotification, err = inapp.ProcessInAppNotifications(notif)
	default:
		log.Printf("⚠️ Worker %d unknown notification channel: %s", w.WorkerID, notif.Channel)
		return fmt.Errorf("unknown notification channel: %s", notif.Channel)
	}

	if err == nil {
		database := db.GetMySQLDB()
		dbNotif := db.Notification{
			ID:                 uuid.MustParse(sentNotification.ID),
			ApplicationID:      uuid.MustParse(sentNotification.ApplicationID),
			QueueID:            sentNotification.QueueID,
			Channel:            string(sentNotification.Channel),
			Provider:           sentNotification.Provider,
			Recipient:          sentNotification.Recipient,
			Subject:            sentNotification.Subject,
			Message:            sentNotification.Message,
			Status:             sentNotification.Status,
			Attempts:           sentNotification.Attempts,
			MessageContentType: sentNotification.MessageContentType,
			TemplateID:         sentNotification.TemplateID,
		}
		dbErr := database.Create(&dbNotif).Error
		if dbErr != nil {
			log.Printf("Error saving notification record to DB: %v", dbErr)
		}
	}

	return nil
}
