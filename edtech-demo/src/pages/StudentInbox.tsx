import { useEffect, useState, useCallback } from 'react'
import { formatDistanceToNow } from 'date-fns'
import { Inbox, RefreshCw } from 'lucide-react'
import { toast } from 'sonner'

import { Button } from '@/components/Button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/Card'
import { NotificationIcon } from '@/components/NotificationIcon'
import {
  getInAppNotifications,
  getUnreadCount,
  getWebSocketUrl,
  markAsRead,
  type InAppNotification,
} from '@/lib/api'

interface StudentInboxProps {
  userId: string
  refreshTrigger?: number
}

function formatRelativeTime(value: string): string {
  const parsed = new Date(value)
  if (Number.isNaN(parsed.getTime())) {
    return 'Unknown time'
  }
  return formatDistanceToNow(parsed, { addSuffix: true })
}

type IncomingNotification = Partial<InAppNotification> & {
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

function normalizeIncomingNotification(raw: IncomingNotification): InAppNotification {
  const createdAt = raw.created_at ?? raw.CreatedAt ?? new Date().toISOString()
  const parsed = new Date(createdAt)

  return {
    id: raw.id ?? raw.ID ?? `ws-${Date.now()}`,
    application_id: raw.application_id ?? raw.ApplicationID ?? '',
    recipient: raw.recipient ?? raw.Recipient ?? '',
    channel: raw.channel ?? raw.Channel ?? 'InApp',
    provider: raw.provider ?? raw.Provider ?? 'inapp',
    message: raw.message ?? raw.Message ?? '',
    subject: raw.subject ?? raw.Subject,
    status: raw.status ?? raw.Status ?? 'sent',
    read: raw.read ?? raw.Read ?? false,
    created_at: Number.isNaN(parsed.getTime()) ? new Date().toISOString() : parsed.toISOString(),
  }
}

export function StudentInbox({ userId, refreshTrigger = 0 }: StudentInboxProps) {
  const [notifications, setNotifications] = useState<InAppNotification[]>([])
  const [unreadCount, setUnreadCount] = useState(0)
  const [loading, setLoading] = useState(true)

  const loadNotifications = useCallback(async () => {
    try {
      setLoading(true)
      const [notifs, unread] = await Promise.all([
        getInAppNotifications(userId),
        getUnreadCount(userId),
      ])
      setNotifications(notifs)
      setUnreadCount(unread)
    } catch (error) {
      toast.error(
        `Failed to load notifications: ${error instanceof Error ? error.message : 'Unknown error'}`,
      )
    } finally {
      setLoading(false)
    }
  }, [userId])

  useEffect(() => {
    void loadNotifications()
  }, [userId, refreshTrigger, loadNotifications])

  useEffect(() => {
    if (!userId) {
      return
    }

    let socket: WebSocket | null = null
    let reconnectTimer: number | null = null
    let reconnectAttempts = 0
    let isActive = true
    let isConnecting = false

    const scheduleReconnect = () => {
      if (!isActive || reconnectAttempts >= 5 || reconnectTimer !== null) {
        return
      }

      const delay = Math.min(1000 * 2 ** reconnectAttempts, 10000)
      reconnectAttempts += 1
      reconnectTimer = window.setTimeout(() => {
        reconnectTimer = null
        connect()
      }, delay)
    }

    const connect = () => {
      if (!isActive || isConnecting) {
        return
      }

      isConnecting = true
      socket = new WebSocket(getWebSocketUrl())

      socket.onopen = () => {
        if (!isActive) {
          socket?.close()
          return
        }

        isConnecting = false
        reconnectAttempts = 0
      }

      socket.onmessage = (event) => {
        try {
          const payload = JSON.parse(event.data) as IncomingNotification
          const incoming = normalizeIncomingNotification(payload)

          if (incoming.recipient && incoming.recipient !== userId) {
            return
          }

          setNotifications((prev) => {
            const exists = prev.some((n) => n.id === incoming.id)
            if (exists) {
              return prev
            }
            return [incoming, ...prev]
          })

          if (!incoming.read) {
            setUnreadCount((prev) => prev + 1)
          }

          toast.success(incoming.subject || 'New in-app notification', {
            description: incoming.message || 'A notification has arrived.',
          })
        } catch {
          // Ignore malformed websocket payloads.
        }
      }

      socket.onerror = () => {
        isConnecting = false
      }

      socket.onclose = () => {
        isConnecting = false
        if (!isActive) {
          return
        }

        scheduleReconnect()
      }
    }

    connect()

    return () => {
      isActive = false
      if (reconnectTimer !== null) {
        window.clearTimeout(reconnectTimer)
        reconnectTimer = null
      }
      if (socket) {
        socket.onclose = null
        socket.onerror = null
      }
      socket?.close()
    }
  }, [userId])

  const handleMarkAsRead = useCallback(
    async (notificationId: string) => {
      try {
        await markAsRead(notificationId)
        setNotifications((prev) =>
          prev.map((n) => (n.id === notificationId ? { ...n, read: true } : n)),
        )
        setUnreadCount((prev) => Math.max(0, prev - 1))
      } catch (error) {
        toast.error(`Error: ${error instanceof Error ? error.message : 'Unknown error'}`)
      }
    },
    [],
  )

  return (
    <Card className="border-blue-200 bg-linear-to-br from-blue-50 to-cyan-50">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Inbox className="h-5 w-5" />
            <CardTitle>Student Inbox</CardTitle>
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={() => void loadNotifications()}
            disabled={loading}
          >
            <RefreshCw className="h-4 w-4" />
            Refresh
          </Button>
        </div>
        <CardDescription>
          Live in-app notifications delivered in real-time
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="mb-4 flex items-center justify-between rounded-lg border border-blue-200 bg-blue-100 px-4 py-2">
          <span className="text-sm font-medium text-blue-900">Unread</span>
          <span className="text-lg font-bold text-blue-700">{unreadCount}</span>
        </div>

        {loading && (
          <div className="flex items-center justify-center py-8">
            <p className="text-gray-500">Loading inbox...</p>
          </div>
        )}

        {!loading && notifications.length === 0 && (
          <div className="rounded-lg border-2 border-dashed border-gray-300 px-4 py-8 text-center">
            <p className="text-gray-600">No notifications yet</p>
            <p className="text-sm text-gray-500">Send a notification from the Campaign Sender</p>
          </div>
        )}

        {!loading && notifications.length > 0 && (
          <div className="max-h-96 space-y-3 overflow-y-auto">
            {notifications.map((notif) => (
              <article
                key={notif.id}
                className={`rounded-lg border p-4 transition ${
                  notif.read
                    ? 'border-gray-200 bg-gray-50'
                    : 'border-blue-300 bg-white shadow-sm'
                }`}
              >
                <div className="mb-2 flex items-start justify-between">
                  <div className="flex items-start gap-3">
                    <NotificationIcon
                      channel={notif.channel as 'email' | 'sms' | 'InApp'}
                      className="mt-1 h-4 w-4 text-gray-600"
                    />
                    <div>
                      <h4 className="font-semibold text-gray-900">
                        {notif.subject || 'Notification'}
                      </h4>
                      <p className="line-clamp-2 text-sm text-gray-700">{notif.message}</p>
                    </div>
                  </div>
                  <span className="text-xs text-gray-500">
                    {formatRelativeTime(notif.created_at)}
                  </span>
                </div>

                <div className="flex items-center justify-between">
                  <div className="flex gap-2">
                    <span className="inline-block rounded-full bg-orange-100 px-2 py-1 text-xs font-medium text-orange-700">
                      {notif.channel}
                    </span>
                    <span className="inline-block rounded-full bg-gray-200 px-2 py-1 text-xs font-medium text-gray-700">
                      {notif.status}
                    </span>
                  </div>
                  {!notif.read && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => void handleMarkAsRead(notif.id)}
                    >
                      Mark Read
                    </Button>
                  )}
                </div>
              </article>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
