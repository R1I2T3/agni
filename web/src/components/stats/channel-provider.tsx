import { Bar, BarChart, CartesianGrid, ResponsiveContainer, XAxis, YAxis } from "recharts"
import type { Notification } from "../type"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "@/components/ui/chart"

interface ChannelPerformanceProps {
  notifications: Notification[]
}

export function ChannelPerformance({ notifications }: ChannelPerformanceProps) {
  const channelStats = notifications.reduce(
    (acc, notification) => {
      const channel = notification.channel
      if (!acc[channel]) {
        acc[channel] = { total: 0, delivered: 0, failed: 0 }
      }
      acc[channel].total++
      if (notification.status === "delivered") acc[channel].delivered++
      if (notification.status === "failed") acc[channel].failed++
      return acc
    },
    {} as Record<string, { total: number; delivered: number; failed: number }>,
  )

  const data = Object.entries(channelStats).map(([channel, stats]) => ({
    channel,
    total: stats.total,
    delivered: stats.delivered,
    failed: stats.failed,
    deliveryRate: ((stats.delivered / stats.total) * 100).toFixed(1),
  }))

  return (
    <Card className="bg-gradient-to-br from-neutral-900 via-red-950 to-amber-900 text-neutral-100 border border-amber-700/40 shadow-md shadow-amber-900/30 rounded-2xl">
      <CardHeader>
        <CardTitle className="text-xl font-semibold text-amber-300 tracking-wide">
          Channel Performance
        </CardTitle>
        <CardDescription className="text-sm text-neutral-400">
          Delivery performance by communication channel
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer
          config={{
            delivered: { label: "Delivered", color: "#eab308" }, // amber
            failed: { label: "Failed", color: "#dc2626" }, // deep red
            total: { label: "Total", color: "#f97316" }, // warm orange
          }}
          className="h-[300px]"
        >
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={data} margin={{ top: 20, right: 20, left: 10, bottom: 5 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.08)" />
              <XAxis dataKey="channel" stroke="#fcd34d" tick={{ fill: "#fcd34d" }} />
              <YAxis stroke="#fcd34d" tick={{ fill: "#fcd34d" }} />
              <ChartTooltip content={<ChartTooltipContent />} />
              <Bar dataKey="delivered" stackId="a" fill="#eab308" radius={[4, 4, 0, 0]} />
              <Bar dataKey="failed" stackId="a" fill="#dc2626" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}
