import type { Notification } from "./type"

// Mock data generator for demonstration
export const generateMockNotifications = (count = 100): Array<Notification> => {
  const statuses = ["delivered", "failed", "pending", "processing", "bounced"]
  const channels = ["email", "sms", "push", "inapp"]
  const providers = ["gmail", "twilio", "resend", "webpush", "in-app"]
  const types = ["welcome", "notification", "alert", "reminder", "marketing"]

  return Array.from({ length: count }, (_, i) => {
    const createdAt = new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000)
    const processedAt = Math.random() > 0.3 ? new Date(createdAt.getTime() + Math.random() * 60 * 60 * 1000) : null

    return {
      id: `notif-${i + 1}`,
      application_id: `app-${Math.floor(Math.random() * 5) + 1}`,
      queue_id: `queue-${Math.floor(Math.random() * 10) + 1}`,
      type: types[Math.floor(Math.random() * types.length)],
      channel: channels[Math.floor(Math.random() * channels.length)],
      provider: providers[Math.floor(Math.random() * providers.length)],
      template_id: Math.random() > 0.5 ? `template-${Math.floor(Math.random() * 20) + 1}` : undefined,
      message_content_type: Math.random() > 0.5 ? "html" : "text",
      recipient: `user${i + 1}@example.com`,
      subject: `Notification Subject ${i + 1}`,
      message: `This is notification message ${i + 1}`,
      status: statuses[Math.floor(Math.random() * statuses.length)],
      attempts: Math.floor(Math.random() * 5) + 1,
      created_at: createdAt.toISOString(),
      updated_at: new Date(createdAt.getTime() + Math.random() * 24 * 60 * 60 * 1000).toISOString(),
      persisted_at: Math.random() > 0.2 ? createdAt.toISOString() : undefined,
      processed_at: processedAt?.toISOString(),
    }
  })
}
