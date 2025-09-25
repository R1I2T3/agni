import { Bar, BarChart, CartesianGrid, ResponsiveContainer, XAxis, YAxis } from "recharts"
import type { Notification } from "@/types"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "@/components/ui/chart"

interface RetryAnalysisProps {
  notifications: Array<Notification>
}

export function RetryAnalysis({ notifications }: RetryAnalysisProps) {
  const attemptStats = notifications.reduce(
    (acc, notification) => {
      const attempts = notification.attempts
      acc[attempts] = (acc[attempts] || 0) + 1
      return acc
    },
    {} as Record<number, number>,
  )

  const data = Object.entries(attemptStats)
    .map(([attempts, count]) => ({
      attempts: `${attempts} attempt${attempts === "1" ? "" : "s"}`,
      count,
      percentage: ((count / notifications.length) * 100).toFixed(1),
    }))
    .sort((a, b) => Number.parseInt(a.attempts) - Number.parseInt(b.attempts))

  return (
    <Card  className="bg-gradient-to-br from-neutral-900 via-red-950 to-amber-900 text-neutral-100 border border-amber-700/40 shadow-md shadow-amber-900/30 rounded-2xl">
      <CardHeader>
        <CardTitle className="text-xl font-semibold text-amber-400">
          Retry Analysis
        </CardTitle>
        <CardDescription className="text-sm text-neutral-400">
          Distribution of notification delivery attempts
        </CardDescription>
      </CardHeader>

      <CardContent>
        <ChartContainer
          config={{
            count: { label: "Count", color: "#f59e0b" }, // muted amber
          }}
          className="h-[300px]"
        >
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={data}
              margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
              barSize={36}
            >
              <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.08)" />
              <XAxis dataKey="attempts" stroke="#d1d5db" />
              <YAxis stroke="#d1d5db" />
              <ChartTooltip
                content={<ChartTooltipContent />}
                formatter={(value, name) => [`${value} notifications`, "Count"]}
                cursor={{ fill: "rgba(245, 158, 11, 0.1)" }}
              />
              <Bar
                dataKey="count"
                fill="url(#warmGradient)"
                radius={[6, 6, 0, 0]}
              />
              <defs>
                <linearGradient id="warmGradient" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stopColor="#f59e0b" /> {/* warm amber */}
                  <stop offset="100%" stopColor="#b45309" /> {/* muted burnt orange */}
                </linearGradient>
              </defs>
            </BarChart>
          </ResponsiveContainer>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}
