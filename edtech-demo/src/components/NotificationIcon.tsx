import { Bell, Mail, Smartphone, Zap } from 'lucide-react'

export interface NotificationIconProps {
  channel: 'email' | 'sms' | 'InApp'
  className?: string
}

export function NotificationIcon({ channel, className }: NotificationIconProps) {
  const icons = {
    email: <Mail className={className} />,
    sms: <Smartphone className={className} />,
    InApp: <Bell className={className} />,
  }
  return icons[channel] || <Zap className={className} />
}
