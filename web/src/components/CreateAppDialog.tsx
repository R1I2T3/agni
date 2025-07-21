import { useState } from 'react'
import { toast } from 'sonner'
import { useMutation, useQueryClient } from '@tanstack/react-query'
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

type CreateAppDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
}

interface CreateApplicationRequest {
  application_name: string
}

interface CreateApplicationResponse {
  id: string
  name: string
  // Add other response properties based on your API
}

export function CreateAppDialog({ open, onOpenChange }: CreateAppDialogProps) {
  const [appName, setAppName] = useState('')
  const queryClient = useQueryClient()

  const createApplicationMutation = useMutation({
    mutationFn: async (
      data: CreateApplicationRequest,
    ): Promise<CreateApplicationResponse> => {
      const res = await fetch('/api/admin/create-application', {
        method: 'POST',
        headers: {
          'Content-type': 'application/json',
        },
        body: JSON.stringify(data),
      })
      if (!res.ok) {
        throw new Error(`Failed to create application: ${res.status}`)
      }
      return res.json()
    },
    onSuccess: async () => {
      toast.success('Application Created successfully')
      await queryClient.invalidateQueries({
        queryKey: ['applications'],
      })

      setAppName('')
      onOpenChange(false)
    },
    onError: (error) => {
      console.error('Failed to create application:', error)
      toast.error('Failed to create application')
    },
  })

  const onCreate = () => {
    if (!appName.trim()) {
      toast.error('Please enter an application name')
      return
    }

    createApplicationMutation.mutate({
      application_name: appName.trim(),
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
              disabled={createApplicationMutation.isPending}
              className="bg-red-900/30 border-red-700/50 text-orange-100 placeholder:text-orange-400/50"
            />
          </div>
        </div>
        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={createApplicationMutation.isPending}
            className="border-red-600/50 text-red-300 hover:bg-red-600/20"
          >
            Cancel
          </Button>
          <Button
            onClick={onCreate}
            disabled={createApplicationMutation.isPending}
            className="bg-gradient-to-r from-red-600 to-orange-600 hover:from-red-500 hover:to-orange-500"
          >
            {createApplicationMutation.isPending ? 'Creating...' : 'Create'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
