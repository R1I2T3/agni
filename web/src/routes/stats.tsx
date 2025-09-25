import { createFileRoute } from '@tanstack/react-router'
import  { useEffect,useState } from 'react'
import { OverviewCards } from '../components/stats/overview-card'
import { StatusDistribution } from '../components/stats/status-distribution'
import { ChannelPerformance } from '../components/stats/channel-provider'
import { RetryAnalysis } from '../components/stats/retry-analysis'
import { ProviderComparison } from '../components/stats/provider-comparison'
import { generateMockNotifications } from '../mock-data'
import type {Notification} from "../type"
import Sidebar from "@/components/Sidebar"

export const Route = createFileRoute('/stats')({
  component: RouteComponent,
})

function RouteComponent() {
  const [notifications, setNotifications] = useState<Array<Notification>>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchNotifications = async () => {
      try {
        const mockData = generateMockNotifications(500)
        setNotifications(mockData)
      } catch (error) {
        console.error("Failed to fetch notifications:", error)
      } finally {
        setLoading(false)
      }
    }

    fetchNotifications()
  }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading analytics...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="flex  bg-white">
      <Sidebar />
    <div className="min-h-screen bg-background overflow-auto flex-1">
      <div className="container mx-auto p-6 space-y-6">
        <div className="space-y-2">
          <h1 className="text-3xl font-bold tracking-tight">Notification Analytics</h1>
          <p className="text-muted-foreground">Comprehensive insights into your notification delivery performance</p>
        </div>

        <OverviewCards notifications={notifications} />

        <div className="grid gap-6 md:grid-cols-2">
          <StatusDistribution notifications={notifications} />
          <ChannelPerformance notifications={notifications} />
        </div>


        <div className="grid gap-6 md:grid-cols-2">
          <RetryAnalysis notifications={notifications} />
          <ProviderComparison notifications={notifications} />
        </div>
      </div>
    </div>
    </div>
  )
}
