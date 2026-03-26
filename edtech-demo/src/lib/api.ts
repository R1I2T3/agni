export type Channel = 'email' | 'sms' | 'InApp'

export interface ApiCredentials {
  token: string
  secret: string
}

export interface SendNotificationPayload {
  channel: Channel
  provider: string
  recipient: string
  subject?: string
  message: string
}

export interface SendNotificationResponse {
  message: string
  queue_id: string
  status: string
}

export interface InAppNotification {
  id: string
  application_id: string
  recipient: string
  channel: string
  provider: string
  message: string
  subject?: string
  status: string
  read: boolean
  created_at: string
}

type RawInAppNotification = Partial<InAppNotification> & {
  ID?: string
  ApplicationID?: string
  Recipient?: string
  Channel?: string
  Provider?: string
  Message?: string
  Subject?: string
  Status?: string
  Read?: boolean
  CreatedAt?: string
}

interface GetNotificationsResponse {
  notifications: RawInAppNotification[]
  total: number
  limit: number
  offset: number
  has_more: boolean
}

function normalizeNotification(raw: RawInAppNotification, index: number): InAppNotification {
  const createdAt = raw.created_at ?? raw.CreatedAt ?? new Date().toISOString()
  const parsed = new Date(createdAt)

  return {
    id: raw.id ?? raw.ID ?? `notification-${index}`,
    application_id: raw.application_id ?? raw.ApplicationID ?? '',
    recipient: raw.recipient ?? raw.Recipient ?? '',
    channel: raw.channel ?? raw.Channel ?? 'InApp',
    provider: raw.provider ?? raw.Provider ?? '',
    message: raw.message ?? raw.Message ?? '',
    subject: raw.subject ?? raw.Subject,
    status: raw.status ?? raw.Status ?? 'sent',
    read: raw.read ?? raw.Read ?? false,
    created_at: Number.isNaN(parsed.getTime()) ? new Date().toISOString() : parsed.toISOString(),
  }
}

// HMAC-SHA256 signature generator
function toHex(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer)
  let hex = ''
  for (const byte of bytes) {
    hex += byte.toString(16).padStart(2, '0')
  }
  return hex
}

async function createSignature(
  secret: string,
  userId: string,
  timestamp: number,
): Promise<string> {
  const message = `${userId}:${timestamp}`
  const encoder = new TextEncoder()
  const secretKey = await crypto.subtle.importKey(
    'raw',
    encoder.encode(secret),
    { name: 'HMAC', hash: 'SHA-256' },
    false,
    ['sign'],
  )
  const signature = await crypto.subtle.sign('HMAC', secretKey, encoder.encode(message))
  return toHex(signature)
}

export async function authenticateUser(
  token: string,
  secret: string,
  userId: string,
): Promise<void> {
  const timestamp = Math.floor(Date.now() / 1000)
  const signature = await createSignature(secret, userId, timestamp)

  const response = await fetch('/api/auth/login', {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      api_token: token,
      user_id: userId,
      timestamp,
      signature,
    }),
  })

  if (!response.ok) {
    const data = await response.json() as { error?: string }
    throw new Error(data.error || 'Authentication failed')
  }
}

export async function sendNotification(
  credentials: ApiCredentials,
  payload: SendNotificationPayload,
): Promise<SendNotificationResponse> {
  const response = await fetch('/api/notification/send', {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      application_token: credentials.token,
      application_secret: credentials.secret,
      channel: payload.channel,
      provider: payload.provider,
      recipient: payload.recipient,
      subject: payload.subject,
      message: payload.message,
      message_content_type: 'text/plain',
    }),
  })

  const data = await response.json() as SendNotificationResponse & { error?: string }
  if (!response.ok) {
    throw new Error(data.error || 'Failed to send notification')
  }
  return data
}

export async function getInAppNotifications(userId: string): Promise<InAppNotification[]> {
  const params = new URLSearchParams({
    user_id: userId,
    limit: '50',
    offset: '0',
  })

  const response = await fetch(`/api/inapp/notifications?${params}`, {
    credentials: 'include',
  })

  const data = await response.json() as GetNotificationsResponse & { error?: string }
  if (!response.ok) {
    throw new Error(data.error || 'Failed to fetch notifications')
  }

  return (data.notifications ?? []).map((notification, index) =>
    normalizeNotification(notification, index),
  )
}

export async function getUnreadCount(userId: string): Promise<number> {
  const params = new URLSearchParams({ user_id: userId })
  const response = await fetch(`/api/inapp/notifications/unread-count?${params}`, {
    credentials: 'include',
  })

  const raw = await response.text()
  let data: { unread_count?: number; error?: string } = {}
  if (raw) {
    try {
      data = JSON.parse(raw) as { unread_count?: number; error?: string }
    } catch {
      data = { error: raw }
    }
  }
  if (!response.ok) {
    throw new Error(data.error || `Failed to get unread count (HTTP ${response.status})`)
  }

  return data.unread_count ?? 0
}

export async function markAsRead(notificationId: string): Promise<void> {
  const response = await fetch(`/api/inapp/notifications/${notificationId}/read`, {
    method: 'PUT',
    credentials: 'include',
  })

  if (!response.ok) {
    const data = await response.json() as { error?: string }
    throw new Error(data.error || 'Failed to mark as read')
  }
}

export function getWebSocketUrl(): string {
  const envUrl = import.meta.env.VITE_INAPP_WS_URL as string | undefined
  if (envUrl && envUrl.trim()) {
    return envUrl.trim()
  }

  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  const hostname = window.location.hostname || 'localhost'
  return `${protocol}://${hostname}/ws`
}
