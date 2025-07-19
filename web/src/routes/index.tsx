import { createFileRoute, redirect } from '@tanstack/react-router'
import AdminPanel from '@/components/admin-panel'

export const Route = createFileRoute('/')({
  component: App,
  loader: async () => {
    const res = await fetch('/api/admin/applications')
    if (!res.ok || res.status !== 200) {
      return redirect({
        to: '/login',
      })
    }
    const data = await res.json()
    return data
  },
})

function App() {
  const { applications } = Route.useLoaderData()
  console.log(applications)
  return <AdminPanel props_application={applications} />
}
