import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import type { Notification } from "@/types"
import { TrendingUp, TrendingDown, AlertTriangle } from "lucide-react"

interface OverviewCardsProps {
  notifications: Array<Notification>
}

export function OverviewCards({ notifications }: OverviewCardsProps) {
  const totalNotifications = notifications.length
  const deliveredCount = notifications.filter((n) => n.status === "delivered").length
  const failedCount = notifications.filter((n) => n.status === "failed").length
  const pendingCount = notifications.filter((n) => n.status === "pending").length

  const deliveryRate =
    totalNotifications > 0 ? ((deliveredCount / totalNotifications) * 100).toFixed(1) : "0"
  const failureRate =
    totalNotifications > 0 ? ((failedCount / totalNotifications) * 100).toFixed(1) : "0"

  const avgAttempts =
    totalNotifications > 0
      ? (notifications.reduce((sum, n) => sum + n.attempts, 0) / totalNotifications).toFixed(1)
      : "0"

  const cards = [
    {
      title: "🔥 Total Notifications",
      value: totalNotifications.toLocaleString(),
      sub: `${pendingCount} pending`,
      gradient: "from-orange-700 via-red-700 to-yellow-600",
      text: "text-yellow-200",
      icon: <TrendingUp className="h-4 w-4 text-yellow-300" />,
    },
    {
      title: "⚡ Delivery Rate",
      value: `${deliveryRate}%`,
      sub: `${deliveredCount} delivered`,
      gradient: "from-yellow-600 via-orange-600 to-red-600",
      text: "text-yellow-100",
      icon: <TrendingUp className="h-4 w-4 text-yellow-300" />,
    },
    {
      title: "❌ Failure Rate",
      value: `${failureRate}%`,
      sub: `${failedCount} failed`,
      gradient: "from-red-800 via-orange-700 to-yellow-700",
      text: "text-red-200",
      icon: <TrendingDown className="h-4 w-4 text-red-400" />,
    },
    {
      title: "⚠️ Avg Attempts",
      value: avgAttempts,
      sub: "per notification",
      gradient: "from-orange-800 via-red-700 to-yellow-700",
      text: "text-orange-200",
      icon: <AlertTriangle className="h-4 w-4 text-yellow-300" />,
    },
  ]

  return (
    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
      {cards.map((card, i) => (
        <Card
          key={i}
          className={`bg-gradient-to-br ${card.gradient} text-white shadow-lg shadow-black/20 border border-white/10 rounded-2xl transition hover:scale-[1.02] hover:shadow-xl`}
        >
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className={`text-sm font-medium ${card.text} drop-shadow-sm`}>
              {card.title}
            </CardTitle>
            {card.icon}
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-extrabold tracking-tight">{card.value}</div>
            <p className={`text-xs opacity-80 ${card.text}`}>{card.sub}</p>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
