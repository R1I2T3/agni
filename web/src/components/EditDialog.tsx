import { useTransition } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { Button } from './ui/button'
import type { Application } from '@/lib/type'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

type EditDialogProps = {
  open: boolean
  onClose: () => void
  app: Application | null
  showSecret: boolean
  setShowSecret: (show: boolean) => void
  copyToClipboard: (text: string, type: string) => void
  onSave: () => void
}
export function EditDialog({ open, onClose, app }: EditDialogProps) {
  const [pending, startTransition] = useTransition()
  const queryClient = useQueryClient()

  const onUpdate = () => {
    startTransition(async () => {
      // Simulate a delete action
      const res = await fetch('/api/admin/regenerate-token', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ application_name: app?.name }),
      })
      if (res.status !== 200) {
        toast.error('Failed to regenerate token and secret')
      } else {
        toast.success('Application Updated successfully')
        queryClient.invalidateQueries({
          queryKey: ['applications'],
        })
      }
      onClose()
    })
  }
  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="max-w-2xl bg-gradient-to-br from-red-950 to-orange-950 border-red-800/50 text-orange-100">
        <DialogHeader>
          <DialogTitle className="text-orange-200">
            🔥 Edit Application
          </DialogTitle>
          <DialogDescription className="text-orange-300/70">
            New credentials will be generated for this application.
          </DialogDescription>
        </DialogHeader>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={onClose}
            className="border-red-600/50 text-red-300 hover:bg-red-600/20"
          >
            Cancel
          </Button>
          <Button
            onClick={onUpdate}
            className="bg-gradient-to-r from-red-600 to-orange-600 hover:from-red-500 hover:to-orange-500"
          >
            {pending ? 'Updating...' : 'Update applications'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
