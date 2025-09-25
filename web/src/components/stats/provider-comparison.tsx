"use client"

import { AlertTriangle, TrendingDown, TrendingUp } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import type { Notification } from "@/types"

interface ProviderComparisonProps {
  notifications: Array<Notification>
}

export function ProviderComparison({ notifications }: ProviderComparisonProps) {
  const providerStats = notifications.reduce(
    (acc, notification) => {
      const provider = notification.provider
      if (!acc[provider]) {
        acc[provider] = { total: 0, delivered: 0, failed: 0, avgAttempts: 0 }
      }
      acc[provider].total++
      if (notification.status === "delivered") acc[provider].delivered++
      if (notification.status === "failed") acc[provider].failed++
      acc[provider].avgAttempts += notification.attempts
      return acc
    },
    {} as Record<string, { total: number; delivered: number; failed: number; avgAttempts: number }>,
  )

  const providerData = Object.entries(providerStats).map(([provider, stats]) => ({
    provider,
    total: stats.total,
    delivered: stats.delivered,
    failed: stats.failed,
    deliveryRate: ((stats.delivered / stats.total) * 100).toFixed(1),
    avgAttempts: (stats.avgAttempts / stats.total).toFixed(1),
  }))

  return (
    <Card  className="bg-gradient-to-br from-neutral-900 via-red-950 to-amber-900 text-neutral-100 border border-amber-700/40 shadow-md shadow-amber-900/30 rounded-2xl">
      <CardHeader>
        <CardTitle className="text-2xl font-extrabold text-yellow-300 drop-shadow-lg">
          🔥 Provider Comparison
        </CardTitle>
        <CardDescription className="text-orange-200/80">
          Performance metrics by notification provider
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-5">
          {providerData.map((provider) => {
            const rate = parseFloat(provider.deliveryRate)
            let badge, icon, barColor

            if (rate >= 80) {
              badge = "bg-green-500/20 text-green-300 border border-green-500/40"
              icon = <TrendingUp className="h-4 w-4 text-green-400" />
              barColor = "from-green-400 to-emerald-500"
            } else if (rate >= 50) {
              badge = "bg-yellow-500/20 text-yellow-300 border border-yellow-500/40"
              icon = <AlertTriangle className="h-4 w-4 text-yellow-400" />
              barColor = "from-yellow-400 to-orange-500"
            } else {
              badge = "bg-red-500/20 text-red-300 border border-red-500/40"
              icon = <TrendingDown className="h-4 w-4 text-red-400" />
              barColor = "from-red-400 to-pink-500"
            }

            return (
              <div
                key={provider.provider}
                className="p-4 rounded-xl border border-white/10 bg-white/5 backdrop-blur-sm hover:shadow-lg hover:shadow-orange-500/20 transition"
              >
                <div className="flex items-center justify-between mb-2">
                  <p className="font-semibold capitalize text-yellow-100 flex items-center gap-2">
                    {provider.provider} {icon}
                  </p>
                  <span
                    className={`px-2 py-0.5 rounded-full text-xs font-medium ${badge}`}
                  >
                    {rate >= 80 ? "Excellent" : rate >= 50 ? "Moderate" : "Poor"}
                  </span>
                </div>

                {/* Gradient Progress Bar */}
                <div className="w-full bg-white/10 h-2 rounded-full overflow-hidden">
                  <div
                    className={`h-2 rounded-full bg-gradient-to-r ${barColor} transition-all duration-700`}
                    style={{ width: `${provider.deliveryRate}%` }}
                  ></div>
                </div>

                {/* Stats row */}
                <div className="flex justify-between text-sm mt-2 text-orange-200/80">
                  <span>
                    {provider.delivered} delivered / {provider.failed} failed
                  </span>
                  <span>{provider.avgAttempts} avg attempts</span>
                </div>
              </div>
            )
          })}
        </div>
      </CardContent>
    </Card>
  )
}
