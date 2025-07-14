"use client"

import { useCallback, useState } from "react"
import {  Plus} from "lucide-react"
import { toast } from "sonner"
import Sidebar from "./Sidebar"
import { CreateAppDialog } from "./CreateAppDialog"
import { TokenDialog } from "./TokenDialog"
import { EditDialog } from "./EditDialog"
import { AppTableRow } from "./ApplicationTableRow"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { DeleteDialog } from "@/components/DeleteDialog"
import { Table, TableBody, TableHead, TableHeader, TableRow } from "@/components/ui/table"

interface Application {
  id: string
  name: string
  token: string
  secret: string
  createdAt: string
}

function generateRandomString(length: number) {
  const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
  let result = ""
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
}
export default function AdminPanel() {
  const [applications, setApplications] = useState<Array<Application>>([
    {
      id: "1",
      name: "Mobile App",
      token: "app_token_1234567890abcdef",
      secret: "secret_abcdef1234567890",
      createdAt: "2024-01-15",
    },
    {
      id: "2",
      name: "Web Dashboard",
      token: "app_token_fedcba0987654321",
      secret: "secret_0987654321fedcba",
      createdAt: "2024-01-10",
    },
  ])
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const [isTokenModalOpen, setIsTokenModalOpen] = useState(false)
  const [isEditModalOpen, setIsEditModalOpen] = useState(false)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)
  const [newAppName, setNewAppName] = useState("")
  const [currentApp, setCurrentApp] = useState<Application | null>(null)
  const [showSecret, setShowSecret] = useState(false)
  const [deleteAppId, setDeleteAppId] = useState<string | null>(null)

  const copyToClipboard = useCallback((text: string, type: string) => {
    navigator.clipboard.writeText(text)
    toast(`Copied to ${type} Clipboard`)
  }, [])

  const handleCreateApplication = useCallback(() => {
    if (!newAppName.trim()) {
      toast("Please enter a valid application name.")
      return
    }
    const newApp: Application = {
      id: Date.now().toString(),
      name: newAppName,
      token: `app_token_${generateRandomString(16)}`,
      secret: `secret_${generateRandomString(16)}`,
      createdAt: new Date().toISOString().split("T")[0],
    }
    setApplications((apps) => [...apps, newApp])
    setCurrentApp(newApp)
    setNewAppName("")
    setIsCreateModalOpen(false)
    setIsTokenModalOpen(true)
  }, [newAppName])

  const handleEdit = useCallback((app: Application) => {
    const updatedApp = {
      ...app,
      token: `app_token_${generateRandomString(16)}`,
      secret: `secret_${generateRandomString(16)}`,
    }
    setCurrentApp(updatedApp)
    setIsEditModalOpen(true)
  }, [])

  const handleSaveEdit = useCallback(() => {
    if (currentApp) {
      setApplications((apps) => apps.map((app) => (app.id === currentApp.id ? currentApp : app)))
      setIsEditModalOpen(false)
      setCurrentApp(null)
      setShowSecret(false)
    }
  }, [currentApp])

  const handleDelete = useCallback((id: string) => {
    setDeleteAppId(id)
    setIsDeleteModalOpen(true)
  }, [])

  const confirmDelete = useCallback(() => {
    if (deleteAppId) {
      setApplications((apps) => apps.filter((app) => app.id !== deleteAppId))
      setIsDeleteModalOpen(false)
      setDeleteAppId(null)
      toast("Success,Application Deleted Successfully")
    }
  }, [deleteAppId])

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
              <p className="text-gray-600 mt-1">Manage your application credentials with fire power</p>
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
              <CardTitle className="text-orange-200">Application List</CardTitle>
              <CardDescription className="text-orange-300/70">
                All registered applications with their tokens and management options
              </CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow className="border-red-800/30 hover:bg-red-900/20">
                    <TableHead className="text-orange-200">Name</TableHead>
                    <TableHead className="text-orange-200">Application Token</TableHead>
                    <TableHead className="text-orange-200">Created</TableHead>
                    <TableHead className="text-right text-orange-200">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {applications.map((app) => (
                    <AppTableRow
                      key={app.id}
                      app={app}
                      onEdit={handleEdit}
                      onDelete={handleDelete}
                    />
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
        appName={newAppName}
        setAppName={setNewAppName}
        onCreate={handleCreateApplication}
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
        onCancel={() => {
          setIsDeleteModalOpen(false)
          setDeleteAppId(null)
        }}
        onDelete={confirmDelete}
      />
    </div>
  )
}
