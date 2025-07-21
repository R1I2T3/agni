import { createFileRoute } from '@tanstack/react-router'
import { queryOptions, useQuery } from '@tanstack/react-query'
import type { Application } from '@/lib/type'
import AdminPanel from '@/components/admin-panel'
import { queryClient } from '@/main'

export const applicationsQueryOptions = queryOptions({
  queryKey: ['applications'] as const,
  queryFn: async (): Promise<Array<Application>> => {
    const res = await fetch('/api/admin/applications')
    if (!res.ok) {
      throw new Error(`Failed to fetch applications: ${res.status}`)
    }
    const data = await res.json()
    const { applications } = data
    return applications
  },
  retry: (failureCount: number, error: Error) => {
    // Don't retry on 401/403 errors
    if (error.message.includes('401') || error.message.includes('403')) {
      return false
    }
    return failureCount < 3
  },
})

export const Route = createFileRoute('/')({
  component: App,
  loader: () => queryClient.ensureQueryData(applicationsQueryOptions),
})

function App() {
  const { data, isError, error, refetch, isLoading } = useQuery(
    applicationsQueryOptions,
  )
  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen">
        <p className="text-gray-500">Loading applications...</p>
      </div>
    )
  }
  const applications = data ?? []

  if (isError) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen gap-4">
        <p className="text-red-600">
          Failed to load applications:{' '}
          {error instanceof Error ? error.message : 'Unknown error'}
        </p>
        <button
          onClick={() => refetch()}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Retry
        </button>
      </div>
    )
  }

  return <AdminPanel props_application={applications} />
}
