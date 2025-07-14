import { BarChart3,Flame,Settings,Shield,Users } from "lucide-react"
import { Button } from "@/components/ui/button"

export default function Sidebar() {
  return (
    <div className="w-64 bg-gradient-to-b from-red-950 to-orange-950 shadow-2xl border-r border-red-800/50 backdrop-blur-sm">
      <div className="p-6 border-b border-red-800/30">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-gradient-to-r from-red-500 to-orange-500 rounded-lg">
            <Flame className="h-6 w-6 text-white" />
          </div>
          <h2 className="text-xl font-bold bg-gradient-to-r from-orange-300 to-red-300 bg-clip-text text-transparent">
            Agni
          </h2>
        </div>
      </div>
      <nav className="mt-6">
        <div className="px-4 space-y-2">
          <Button className="w-full justify-start gap-3 bg-gradient-to-r from-red-600 to-orange-600 hover:from-red-500 hover:to-orange-500 text-white shadow-lg">
            <Shield className="h-4 w-4" />
            Applications
          </Button>
          <Button
            variant="ghost"
            className="w-full justify-start gap-3 text-orange-200 hover:text-white hover:bg-red-800/30"
          >
            <Users className="h-4 w-4" />
            Users
          </Button>
          <Button
            variant="ghost"
            className="w-full justify-start gap-3 text-orange-200 hover:text-white hover:bg-red-800/30"
          >
            <BarChart3 className="h-4 w-4" />
            Analytics
          </Button>
          <Button
            variant="ghost"
            className="w-full justify-start gap-3 text-orange-200 hover:text-white hover:bg-red-800/30"
          >
            <Settings className="h-4 w-4" />
            Settings
          </Button>
        </div>
      </nav>
    </div>
  )
}
