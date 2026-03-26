import { useCallback, useState } from 'react'
import { Key, LogOut } from 'lucide-react'
import { toast } from 'sonner'

import { Button } from '@/components/Button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/Card'
import { Input } from '@/components/Input'
import { authenticateUser, type ApiCredentials } from '@/lib/api'

interface AuthPanelProps {
  onAuthenticated: (credentials: ApiCredentials, userId: string) => void
  onLogout: () => void
  isAuthenticated: boolean
}

export function AuthPanel({ onAuthenticated, onLogout, isAuthenticated }: AuthPanelProps) {
  const [token, setToken] = useState('')
  const [secret, setSecret] = useState('')
  const [userId, setUserId] = useState('student_001')
  const [loading, setLoading] = useState(false)

  const handleAuth = useCallback(async () => {
    if (!token || !secret || !userId) {
      toast.error('All fields are required')
      return
    }

    setLoading(true)
    try {
      await authenticateUser(token, secret, userId)
      onAuthenticated({ token, secret }, userId)
      toast.success('Authenticated successfully!')
    } catch (error) {
      toast.error(`Auth failed: ${error instanceof Error ? error.message : 'Unknown error'}`)
    } finally {
      setLoading(false)
    }
  }, [token, secret, userId, onAuthenticated])

  return (
    <Card className="border-teal-200 bg-gradient-to-br from-teal-50 to-green-50">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Key className="h-5 w-5" />
            <CardTitle>Authentication</CardTitle>
          </div>
          {isAuthenticated && (
            <Button
              variant="destructive"
              size="sm"
              onClick={onLogout}
            >
              <LogOut className="h-4 w-4" />
              Logout
            </Button>
          )}
        </div>
        <CardDescription>
          {isAuthenticated
            ? 'You are authenticated. Ready to send notifications!'
            : 'Sign in with application credentials from the admin panel'}
        </CardDescription>
      </CardHeader>
      <CardContent>
        {!isAuthenticated ? (
          <div className="space-y-4">
            <div>
              <label className="mb-2 block text-sm font-medium text-gray-900">
                Application Token
              </label>
              <Input
                type="password"
                placeholder="Paste your API token"
                value={token}
                onChange={(e) => setToken(e.target.value)}
              />
            </div>

            <div>
              <label className="mb-2 block text-sm font-medium text-gray-900">
                Application Secret
              </label>
              <Input
                type="password"
                placeholder="Paste your API secret"
                value={secret}
                onChange={(e) => setSecret(e.target.value)}
              />
            </div>

            <div>
              <label className="mb-2 block text-sm font-medium text-gray-900">
                Student User ID
              </label>
              <Input
                placeholder="e.g., student_001"
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
              />
            </div>

            <Button
              onClick={() => void handleAuth()}
              disabled={loading}
              className="w-full"
              size="lg"
            >
              {loading ? 'Authenticating...' : 'Sign In'}
            </Button>
          </div>
        ) : (
          <div className="space-y-3 rounded-lg border border-teal-200 bg-white p-4">
            <p className="text-sm text-gray-700">
              <strong>User ID:</strong> {userId}
            </p>
            <p className="text-xs text-gray-600">
              Ready to send and receive notifications through the Agni microservice.
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
