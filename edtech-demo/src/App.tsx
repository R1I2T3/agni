import { useState, useCallback } from 'react'
import { Toaster } from 'sonner'
import { BookOpen } from 'lucide-react'

import { AuthPanel } from '@/pages/AuthPanel'
import { CampaignSender } from '@/pages/CampaignSender'
import { StudentInbox } from '@/pages/StudentInbox'
import type { ApiCredentials, Channel } from '@/lib/api'

export default function App() {
  const [credentials, setCredentials] = useState<ApiCredentials | null>(null)
  const [userId, setUserId] = useState('')
  const [refreshInbox, setRefreshInbox] = useState(0)

  const handleAuthenticated = useCallback((creds: ApiCredentials, id: string) => {
    setCredentials(creds)
    setUserId(id)
  }, [])

  const handleLogout = useCallback(() => {
    setCredentials(null)
    setUserId('')
    setRefreshInbox(0)
  }, [])

  const handleSent = useCallback((_channels: Channel[]) => {
    // Only refresh inbox if InApp was sent
    if (_channels.includes('InApp')) {
      setRefreshInbox((prev) => prev + 1)
    }
  }, [])

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      {/* Header */}
      <header className="border-b border-gray-200 bg-white shadow-sm">
        <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
          <div className="flex items-center gap-3">
            <div className="rounded-lg bg-orange-100 p-2">
              <BookOpen className="h-6 w-6 text-orange-600" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">EdTech Notification Demo</h1>
              <p className="text-sm text-gray-600">
                Showcase the Agni microservice in action
              </p>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="space-y-6">
          {/* Auth Panel */}
          <AuthPanel
            onAuthenticated={handleAuthenticated}
            onLogout={handleLogout}
            isAuthenticated={!!credentials}
          />

          {/* Campaign & Inbox Grid */}
          {credentials && userId && (
            <div className="grid gap-6 lg:grid-cols-3">
              <div className="lg:col-span-1">
                <CampaignSender
                  credentials={credentials}
                  onSent={handleSent}
                />
              </div>
              <div className="lg:col-span-2">
                <StudentInbox userId={userId} refreshTrigger={refreshInbox} />
              </div>
            </div>
          )}

          {!credentials && (
            <div className="rounded-lg border-2 border-dashed border-gray-300 px-8 py-12 text-center">
              <h2 className="mb-2 text-xl font-semibold text-gray-900">Get Started</h2>
              <p className="text-gray-600">
                Sign in above with your application credentials to begin sending notifications.
              </p>
            </div>
          )}
        </div>
      </main>

      <Toaster />
    </div>
  )
}
