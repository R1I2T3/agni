
import { Copy, Eye, EyeOff } from 'lucide-react'
import { Input } from './ui/input'
import { Button } from './ui/button'
import { Label } from './ui/label'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

interface Application {
  id: string
  name: string
  token: string
  secret: string
  createdAt: string
}

type EditDialogProps = {
  open: boolean
  onClose: () => void
  app: Application | null
  showSecret: boolean
  setShowSecret: (show: boolean) => void
  copyToClipboard: (text: string, type: string) => void
  onSave: () => void
}
export function EditDialog({ open, onClose, app, showSecret, setShowSecret, copyToClipboard, onSave }: EditDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="max-w-2xl bg-gradient-to-br from-red-950 to-orange-950 border-red-800/50 text-orange-100">
        <DialogHeader>
          <DialogTitle className="text-orange-200">🔥 Edit Application</DialogTitle>
          <DialogDescription className="text-orange-300/70">
            New credentials have been generated for this application.
          </DialogDescription>
        </DialogHeader>
        {app && (
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label className="text-orange-200">Application Name</Label>
              <Input value={app.name} readOnly className="bg-red-900/30 border-red-700/50 text-orange-100" />
            </div>
            <div className="space-y-2">
              <Label className="text-orange-200">New Application Token</Label>
              <div className="flex gap-2">
                <Input
                  value={app.token}
                  readOnly
                  className="bg-red-900/30 border-red-700/50 text-orange-100"
                />
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => copyToClipboard(app.token, "Token")}
                  className="border-orange-600/50 text-orange-300 hover:bg-orange-600/20"
                >
                  <Copy className="h-4 w-4" />
                </Button>
              </div>
            </div>
            <div className="space-y-2">
              <Label className="text-orange-200">New Application Secret</Label>
              <div className="flex gap-2">
                <Input
                  type={showSecret ? "text" : "password"}
                  value={app.secret}
                  readOnly
                  className="bg-red-900/30 border-red-700/50 text-orange-100"
                />
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => setShowSecret(!showSecret)}
                  className="border-orange-600/50 text-orange-300 hover:bg-orange-600/20"
                >
                  {showSecret ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </Button>
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => copyToClipboard(app.secret, "Secret")}
                  className="border-orange-600/50 text-orange-300 hover:bg-orange-600/20"
                >
                  <Copy className="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>
        )}
        <DialogFooter>
          <Button
            variant="outline"
            onClick={onClose}
            className="border-red-600/50 text-red-300 hover:bg-red-600/20"
          >
            Cancel
          </Button>
          <Button
            onClick={onSave}
            className="bg-gradient-to-r from-red-600 to-orange-600 hover:from-red-500 hover:to-orange-500"
          >
            Save Changes
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
