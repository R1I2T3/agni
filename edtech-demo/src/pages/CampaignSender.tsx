import { useCallback, useState } from 'react'
import { Mail, Smartphone, Bell, Send } from 'lucide-react'
import { toast } from 'sonner'

import { Button } from '@/components/Button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/Card'
import { Input } from '@/components/Input'
import { Textarea } from '@/components/Textarea'
import { sendNotification, type Channel, type ApiCredentials } from '@/lib/api'

interface CampaignSenderProps {
  credentials: ApiCredentials | null
  onSent?: (channels: Channel[]) => void
}

const channelInfo = {
  email: { icon: Mail, label: 'Email', provider: 'Resend', color: 'bg-blue-100 text-blue-700' },
  sms: { icon: Smartphone, label: 'SMS', provider: 'twilio', color: 'bg-green-100 text-green-700' },
  InApp: { icon: Bell, label: 'In-App', provider: 'InApp', color: 'bg-purple-100 text-purple-700' },
}

export function CampaignSender({ credentials, onSent }: CampaignSenderProps) {
  const [channels, setChannels] = useState<Set<Channel>>(new Set(['InApp']))
  const [recipient, setRecipient] = useState('student_001')
  const [email, setEmail] = useState('student001@example.com')
  const [phone, setPhone] = useState('+15551234567')
  const [subject, setSubject] = useState('Important Class Update')
  const [message, setMessage] = useState(
    'Hi! Your assignment is due in 24 hours. Make sure to submit before the deadline.',
  )
  const [loading, setLoading] = useState(false)

  const toggleChannel = useCallback((channel: Channel) => {
    setChannels((prev) => {
      const next = new Set(prev)
      if (next.has(channel)) next.delete(channel)
      else next.add(channel)
      return next
    })
  }, [])

  const handleSend = useCallback(async () => {
    if (!credentials) {
      toast.error('Not authenticated. Please set credentials first.')
      return
    }

    if (channels.size === 0) {
      toast.error('Select at least one channel')
      return
    }

    if (!message.trim()) {
      toast.error('Message is required')
      return
    }

    if (channels.has('InApp') && !recipient.trim()) {
      toast.error('Recipient ID is required for In-App channel')
      return
    }

    if (channels.has('email') && !email.trim()) {
      toast.error('Email is required for Email channel')
      return
    }

    if (channels.has('sms') && !phone.trim()) {
      toast.error('Phone number is required for SMS channel')
      return
    }

    setLoading(true)
    let successCount = 0
    let failureCount = 0

    for (const channel of Array.from(channels)) {
      try {
        const channelRecipient =
          channel === 'email' ? email : channel === 'sms' ? phone : recipient

        await sendNotification(credentials, {
          channel,
          provider: channelInfo[channel].provider,
          recipient: channelRecipient,
          subject: channel === 'email' ? subject : undefined,
          message,
        })
        successCount += 1
        toast.success(`${channelInfo[channel].label} sent successfully`)
      } catch (error) {
        failureCount += 1
        toast.error(`Failed to send ${channelInfo[channel].label}: ${error instanceof Error ? error.message : 'Unknown error'}`)
      }
    }

    setLoading(false)

    if (successCount > 0) {
      toast.success(`${successCount} notification(s) queued successfully`)
      onSent?.(Array.from(channels))
    }
  }, [credentials, channels, recipient, email, phone, subject, message, onSent])

  return (
    <Card className="border-orange-200 bg-linear-to-br from-orange-50 to-amber-50">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Send className="h-5 w-5" />
          Send Campaign
        </CardTitle>
        <CardDescription>
          Choose channels and send notifications to students instantly
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div>
          <label className="mb-2 block text-sm font-medium text-gray-900">
            Select Channels
          </label>
          <div className="grid gap-3 sm:grid-cols-3">
            {Object.entries(channelInfo).map(([key, info]) => {
              const channel = key as Channel
              const IconComponent = info.icon
              const isSelected = channels.has(channel)
              return (
                <button
                  key={channel}
                  onClick={() => toggleChannel(channel)}
                  className={`relative flex items-center gap-2 rounded-lg border-2 p-3 transition ${
                    isSelected
                      ? 'border-orange-500 bg-orange-100'
                      : 'border-gray-200 hover:border-gray-300'
                  }`}
                >
                  <IconComponent className="h-4 w-4" />
                  <span className="text-sm font-medium">{info.label}</span>
                  {isSelected && (
                    <span className="ml-auto text-xs font-bold text-orange-600">✓</span>
                  )}
                </button>
              )
            })}
          </div>
        </div>

        <div className="grid gap-4 sm:grid-cols-2">
          <div>
            <label className="mb-2 block text-sm font-medium text-gray-900">
              Email
            </label>
            <Input
              type="email"
              placeholder="e.g., student@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>

          <div>
            <label className="mb-2 block text-sm font-medium text-gray-900">
              Phone Number
            </label>
            <Input
              placeholder="e.g., +15551234567"
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
            />
          </div>
        </div>

        <div>
          <label className="mb-2 block text-sm font-medium text-gray-900">
            Recipient ID (In-App)
          </label>
          <Input
            placeholder="e.g., student_001"
            value={recipient}
            onChange={(e) => setRecipient(e.target.value)}
          />
        </div>

        <div>
          <label className="mb-2 block text-sm font-medium text-gray-900">
            Subject (for email)
          </label>
          <Input
            placeholder="Email subject line"
            value={subject}
            onChange={(e) => setSubject(e.target.value)}
          />
        </div>

        <div>
          <label className="mb-2 block text-sm font-medium text-gray-900">
            Message
          </label>
          <Textarea
            placeholder="Your notification message..."
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            className="min-h-32"
          />
        </div>

        <Button
          onClick={() => void handleSend()}
          disabled={loading || !credentials}
          className="w-full"
          size="lg"
        >
          {loading ? 'Sending...' : 'Send Notifications'}
        </Button>
      </CardContent>
    </Card>
  )
}
