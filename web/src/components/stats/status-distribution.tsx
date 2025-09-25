import { Cell, Pie, PieChart, ResponsiveContainer } from "recharts"
import type { Notification } from "@/types"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "@/components/ui/chart"

interface StatusDistributionProps {
  notifications: Array<Notification>
}

const COLORS = {
  delivered: "#eab308",   // elegant amber
  failed: "#dc2626",      // deep red
  pending: "#f97316",     // warm orange
  processing: "#facc15",  // golden yellow
  bounced: "#fb923c",     // soft orange
}

export function StatusDistribution({ notifications }: StatusDistributionProps) {
  const statusCounts = notifications.reduce(
    (acc, notification) => {
      acc[notification.status] = (acc[notification.status] || 0) + 1
      return acc
    },
    {} as Record<string, number>,
  )

  const data = Object.entries(statusCounts).map(([status, count]) => ({
    name: status,
    value: count,
    percentage: ((count / notifications.length) * 100).toFixed(1),
  }))

  return (
    <Card className="bg-gradient-to-br from-neutral-900 via-red-950 to-amber-900 text-neutral-100 border border-amber-700/40 shadow-md shadow-amber-900/30 rounded-2xl">
      <CardHeader>
        <CardTitle className="text-xl font-semibold text-amber-300 tracking-wide">
          Status Distribution
        </CardTitle>
        <CardDescription className="text-sm text-neutral-400">
          Breakdown of notification statuses
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer
          config={{
            delivered: { label: "Delivered", color: COLORS.delivered },
            failed: { label: "Failed", color: COLORS.failed },
            pending: { label: "Pending", color: COLORS.pending },
            processing: { label: "Processing", color: COLORS.processing },
            bounced: { label: "Bounced", color: COLORS.bounced },
          }}
          className="h-[300px]"
        >
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={data}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percentage }) => `${name}: ${percentage}%`}
                outerRadius={95}
                dataKey="value"
              >
                {data.map((entry, index) => (
                  <Cell
                    key={`cell-${index}`}
                    fill={COLORS[entry.name as keyof typeof COLORS] || "#fbbf24"}
                    stroke="rgba(255,255,255,0.1)"
                    strokeWidth={1}
                  />
                ))}
              </Pie>
              <ChartTooltip content={<ChartTooltipContent />} />
            </PieChart>
          </ResponsiveContainer>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}
