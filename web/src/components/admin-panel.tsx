'use client'

import { useCallback, useState } from 'react'
import { Plus } from 'lucide-react'
import { toast } from 'sonner'
import Sidebar from './Sidebar'
import { CreateAppDialog } from './CreateAppDialog'
import { TokenDialog } from './TokenDialog'
import { EditDialog } from './EditDialog'
import { AppTableRow } from './ApplicationTableRow'
import type { Application } from '@/lib/type'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { DeleteDialog } from '@/components/DeleteDialog'
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

export default function AdminPanel({
  props_application,
}: {
  props_application: Array<Application>
}) {
  const [applications, setApplications] =
    useState<Array<Application>>(props_application)
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const [isTokenModalOpen, setIsTokenModalOpen] = useState(false)
  const [isEditModalOpen, setIsEditModalOpen] = useState(false)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)
  const [currentApp, setCurrentApp] = useState<Application | null>(null)
  const [showSecret, setShowSecret] = useState(false)
  const copyToClipboard = useCallback((text: string, type: string) => {
    navigator.clipboard.writeText(text)
    toast(`Copied to ${type} Clipboard`)
  }, [])

  const handleSaveEdit = useCallback(() => {
    if (currentApp) {
      setApplications((apps) =>
        apps.map((app) => (app.name === currentApp.name ? currentApp : app)),
      )
      setIsEditModalOpen(false)
      setCurrentApp(null)
      setShowSecret(false)
    }
  }, [currentApp])

  return (
    <div className="flex h-screen bg-white">
      <Sidebar />
      <div className="flex-1 overflow-auto">
        <div className="p-8">
          <div className="flex justify-between items-center mb-8">
            <div>
              <h1 className="text-3xl font-bold bg-gradient-to-r from-orange-500 to-red-500 bg-clip-text text-transparent">
                Applications
              </h1>
              <p className="text-gray-600 mt-1">
                Manage your application credentials with fire power
              </p>
            </div>
            <Button
              onClick={() => setIsCreateModalOpen(true)}
              className="gap-2 bg-gradient-to-r from-red-600 to-orange-600 hover:from-red-500 hover:to-orange-500 text-white shadow-lg border border-red-500/50"
            >
              <Plus className="h-4 w-4" />
              Create Application
            </Button>
          </div>
          <Card className="bg-gradient-to-br from-red-950/50 to-orange-950/50 border-red-800/50 backdrop-blur-sm shadow-2xl">
            <CardHeader className="border-b border-red-800/30">
              <CardTitle className="text-orange-200">
                Application List
              </CardTitle>
              <CardDescription className="text-orange-300/70">
                All registered applications with their tokens and management
                options
              </CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow className="border-red-800/30 hover:bg-red-900/20">
                    <TableHead className="text-orange-200">Name</TableHead>
                    <TableHead className="text-orange-200">
                      Application Token
                    </TableHead>
                    <TableHead className="text-orange-200">Created</TableHead>
                    <TableHead className="text-right text-orange-200">
                      Actions
                    </TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {applications.map((app) => (
                    <AppTableRow key={app.name} app={app} />
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </div>
      </div>
      <CreateAppDialog
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
      />
      <TokenDialog
        open={isTokenModalOpen}
        onClose={() => {
          setIsTokenModalOpen(false)
          setCurrentApp(null)
          setShowSecret(false)
        }}
        app={currentApp}
        showSecret={showSecret}
        setShowSecret={setShowSecret}
        copyToClipboard={copyToClipboard}
      />
      <EditDialog
        open={isEditModalOpen}
        onClose={() => {
          setIsEditModalOpen(false)
          setCurrentApp(null)
          setShowSecret(false)
        }}
        app={currentApp}
        showSecret={showSecret}
        setShowSecret={setShowSecret}
        copyToClipboard={copyToClipboard}
        onSave={handleSaveEdit}
      />
      <DeleteDialog
        open={isDeleteModalOpen}
        onCancel={() => setIsDeleteModalOpen(false)}
        onOpenChange={setIsDeleteModalOpen}
      />
    </div>
  )
}
