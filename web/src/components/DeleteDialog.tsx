import { useTransition } from 'react'
import { toast } from 'sonner'
import { useQueryClient } from '@tanstack/react-query'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'

type DeleteDialogProps = {
  open: boolean
  onCancel: () => void
  onOpenChange: (arg: boolean) => void
  appId: string
}
export function DeleteDialog({
  open,
  onCancel,
  onOpenChange,
  appId,
}: DeleteDialogProps) {
  const [pending, startTransition] = useTransition()
  const queryClient = useQueryClient()
  const onDelete = () => {
    startTransition(async () => {
      // Simulate a delete action
      const res = await fetch('/api/admin/delete-application', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ application_name: appId }),
      })
      if (res.status !== 200) {
        toast.error('Failed to delete application')
      } else {
        toast.success('Application deleted successfully')
        queryClient.invalidateQueries({
          queryKey: ['applications'],
        })
      }
      onOpenChange(false)
    })
  }
  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent className="bg-gradient-to-br from-red-950 to-orange-950 border-red-800/50 text-orange-100">
        <AlertDialogHeader>
          <AlertDialogTitle className="text-red-300">
            🔥 Are you sure?
          </AlertDialogTitle>
          <AlertDialogDescription className="text-orange-300/70">
            This action cannot be undone. This will permanently delete the
            application and remove all associated data.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel
            onClick={onCancel}
            className="border-red-600/50 text-red-300 hover:bg-red-600/20"
          >
            Cancel
          </AlertDialogCancel>
          <AlertDialogAction
            onClick={onDelete}
            className="bg-gradient-to-r from-red-700 to-red-600 hover:from-red-600 hover:to-red-500"
          >
            {pending ? 'Deleting...' : 'Delete'}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
