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
  appName: string
  setAppName: (name: string) => void
  onCreate: () => void
}
export function CreateAppDialog({
  open,
  onOpenChange,
  appName,
  setAppName,
  onCreate,
}: CreateAppDialogProps) {
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
            Create
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
