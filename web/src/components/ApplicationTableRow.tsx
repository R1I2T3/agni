import { Edit, Trash2 } from 'lucide-react'
import { TableCell, TableRow } from './ui/table'
import { Button } from './ui/button'
import type { Application } from '@/lib/type'
import { formatToDDMMYY } from '@/lib/utils'

type AppTableRowProps = {
  app: Application
  onEdit: (app: Application) => void
  onDelete: (app: Application) => void
}
export function AppTableRow({ app, onEdit, onDelete }: AppTableRowProps) {
  return (
    <TableRow key={app.name} className="border-red-800/30 hover:bg-red-900/20">
      <TableCell className="font-medium text-orange-100">{app.name}</TableCell>
      <TableCell>
        <code className="bg-red-900/50 border border-red-700/50 px-3 py-1 rounded text-sm text-orange-200">
          {app.api_token}
        </code>
      </TableCell>
      <TableCell className="text-orange-100">
        {formatToDDMMYY(app.created_at)}
      </TableCell>
      <TableCell className="text-right">
        <div className="flex justify-end gap-2">
          <Button
            variant="outline"
            size="icon"
            onClick={() => onEdit(app)}
            className="border-orange-600/50 text-orange-300 hover:bg-orange-600/20 hover:text-orange-200"
          >
            <Edit className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            size="icon"
            onClick={() => onDelete(app)}
            className="border-red-600/50 text-red-300 hover:bg-red-600/20 hover:text-red-200"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </TableCell>
    </TableRow>
  )
}
