import { useAuth } from "@/contexts/auth-context"
import { BookOpen, MessageSquare, Heart, Activity } from "lucide-react"

export default function AppPage() {
  const { user } = useAuth()

  return (
    <div className="max-w-5xl w-full mx-auto space-y-8 animate-in fade-in duration-300">
      
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-zinc-200 dark:border-zinc-800/80 pb-6">
        <div className="space-y-1">
          <h1 className="text-3xl font-extrabold tracking-tight">
            Dashboard
          </h1>
          <p className="text-sm text-zinc-500 dark:text-zinc-400">
            Welcome back, <span className="font-semibold text-zinc-800 dark:text-zinc-200">{user?.name || "Developer"}</span>. Monitor connections, navigate, and manage assets.
          </p>
        </div>
        <div className="flex items-center space-x-2 text-xs font-semibold text-zinc-500 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 px-3 py-1.5 rounded-md shadow-xs shrink-0">
          <span className="h-2 w-2 rounded-full bg-green-500 animate-pulse" />
          <span>Server: Connected</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800/80 p-6 rounded-xl shadow-xs space-y-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-bold text-zinc-400 dark:text-zinc-500 uppercase tracking-wider">Blogs Published</span>
            <BookOpen className="h-4.5 w-4.5 text-zinc-400" />
          </div>
          <div className="space-y-1">
            <h3 className="text-2xl font-bold">12</h3>
            <p className="text-[10px] text-zinc-500 dark:text-zinc-400 font-medium">Active entries in repository</p>
          </div>
        </div>

        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800/80 p-6 rounded-xl shadow-xs space-y-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-bold text-zinc-400 dark:text-zinc-500 uppercase tracking-wider">Comments Count</span>
            <MessageSquare className="h-4.5 w-4.5 text-zinc-400" />
          </div>
          <div className="space-y-1">
            <h3 className="text-2xl font-bold">48</h3>
            <p className="text-[10px] text-zinc-500 dark:text-zinc-400 font-medium">Interactions recorded</p>
          </div>
        </div>

        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800/80 p-6 rounded-xl shadow-xs space-y-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-bold text-zinc-400 dark:text-zinc-500 uppercase tracking-wider">Favorites</span>
            <Heart className="h-4.5 w-4.5 text-zinc-400" />
          </div>
          <div className="space-y-1">
            <h3 className="text-2xl font-bold">8</h3>
            <p className="text-[10px] text-zinc-500 dark:text-zinc-400 font-medium">Saved items</p>
          </div>
        </div>

      </div>

      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800/80 rounded-xl shadow-xs overflow-hidden">
        <div className="border-b border-zinc-200 dark:border-zinc-800/80 p-6 flex items-center justify-between bg-zinc-50/50 dark:bg-zinc-950/20">
          <div>
            <h3 className="text-sm font-bold tracking-tight">System Environment</h3>
            <p className="text-xs text-zinc-400">Current connection status and parameters</p>
          </div>
          <Activity className="h-4 w-4 text-zinc-400" />
        </div>
        
        <div className="divide-y divide-zinc-100 dark:divide-zinc-800/50 text-xs">
          <div className="px-6 py-4 flex items-center justify-between">
            <span className="font-semibold text-zinc-500 dark:text-zinc-400">Database Driver</span>
            <span className="font-mono text-zinc-900 dark:text-zinc-100 bg-zinc-100 dark:bg-zinc-800 px-2 py-0.5 rounded">PostgreSQL (GORM)</span>
          </div>
          <div className="px-6 py-4 flex items-center justify-between">
            <span className="font-semibold text-zinc-500 dark:text-zinc-400">API Connection</span>
            <span className="text-green-600 dark:text-green-400 font-semibold flex items-center gap-1.5">
              <span className="h-1.5 w-1.5 rounded-full bg-green-500 animate-pulse" />
              Active / Connected
            </span>
          </div>
          <div className="px-6 py-4 flex items-center justify-between">
            <span className="font-semibold text-zinc-500 dark:text-zinc-400">Environment Mode</span>
            <span className="font-mono text-zinc-900 dark:text-zinc-100 bg-zinc-100 dark:bg-zinc-800 px-2 py-0.5 rounded">Development</span>
          </div>
          <div className="px-6 py-4 flex items-center justify-between">
            <span className="font-semibold text-zinc-500 dark:text-zinc-400">Current User Session</span>
            <span className="font-semibold">{user?.email || "Unknown user"}</span>
          </div>
        </div>
      </div>

    </div>
  )
}
