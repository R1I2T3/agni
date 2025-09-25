export interface Notification {
  id: string
  application_id: string
  queue_id: string
  type: string
  channel: string
  provider: string
  template_id?: string
  message_content_type?: string
  recipient: string
  subject: string
  message: string
  status: string
  attempts: number
  created_at: string
  updated_at: string
  persisted_at?: string
  processed_at?: string
}