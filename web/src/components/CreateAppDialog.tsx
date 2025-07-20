import { useState, useTransition } from 'react'
import { toast } from 'sonner'
import { Label } from './ui/label'
import { Input } from './ui/input'
import { Button } from './ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

import type { Application } from '@/lib/type'


type CreateAppDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  setCurrentApp: (app: Application | null) => void
  setIsTokenModalOpen: (open: boolean) => void
}
export function CreateAppDialog({ open, onOpenChange ,setCurrentApp,setIsTokenModalOpen}: CreateAppDialogProps) {
  const [appName, setAppName] = useState('')
  const [pending, startTransition] = useTransition()
  const onCreate = () => {
    startTransition(async () => {
      const res = await fetch('/api/admin/create-application', {
        method: 'POST',
        headers: {
          'Content-type': 'application/json',
        },
        body: JSON.stringify({
          application_name: appName,
        }),
      })
      if (!res.ok || res.status !== 200) {
        toast.error('Failed to create application')
      } else {
        toast.success('Application Created successfully')
        const responseData = await res.json()
        console.log('Response Data:', responseData)
        const currentapp: Application = {
          name: appName,
          api_token: responseData['api-token'],
          api_secret: responseData['api-secret'],
          created_at: new Date().toISOString(),
        }
        console.log('Current App:', currentapp)
        setCurrentApp(currentapp)
        onOpenChange(false)
        setIsTokenModalOpen(true)
        
      
      }
    })
  }
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="bg-gradient-to-br from-red-950 to-orange-950 border-red-800/50 text-orange-100">
        <DialogHeader>
          <DialogTitle className="text-orange-200">
            Create New Application
          </DialogTitle>
          <DialogDescription className="text-orange-300/70">
            Enter a name for your new application to generate credentials.
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="appName" className="text-orange-200">
              Application Name
            </Label>
            <Input
              id="appName"
              placeholder="Enter application name"
              value={appName}
              onChange={(e) => setAppName(e.target.value)}
              className="bg-red-900/30 border-red-700/50 text-orange-100 placeholder:text-orange-400/50"
            />
          </div>
        </div>
        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            className="border-red-600/50 text-red-300 hover:bg-red-600/20"
          >
            Cancel
          </Button>
          <Button
            onClick={onCreate}
            className="bg-gradient-to-r from-red-600 to-orange-600 hover:from-red-500 hover:to-orange-500"
          >
            {pending ? 'Creating...' : 'Create'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
