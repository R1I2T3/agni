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
  onDelete: () => void
}
export function DeleteDialog({ open, onCancel, onDelete }: DeleteDialogProps) {
  return (
    <AlertDialog open={open} onOpenChange={onCancel}>
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
            Delete
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
