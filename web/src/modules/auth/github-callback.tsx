import { useEffect } from "react"
import { useNavigate } from "@tanstack/react-router"
import { useGitHubLoginMutation } from "@/service/mutation/auth"
import { toast } from "sonner"
import { Loader2 } from "lucide-react"

export default function GitHubCallback() {
  const navigate = useNavigate()
  const githubMutation = useGitHubLoginMutation()

  useEffect(() => {
    const handleCallback = async () => {
      const params = new URLSearchParams(window.location.search)
      const code = params.get("code")
      
      if (!code) {
        toast.error("GitHub Login Failed", {
          description: "No authorization code returned from GitHub.",
        })
        navigate({ to: "/login" })
        return
      }

      try {
        await githubMutation.mutateAsync(code)
        toast.success("GitHub Login Successful", {
          description: "Welcome to TeraRush!",
        })
        navigate({ to: "/" })
      } catch (error: any) {
        console.error("GitHub auth callback error:", error)
        toast.error("GitHub Login Failed", {
          description: error?.response?.data?.error || "Failed to exchange GitHub authorization code.",
        })
        navigate({ to: "/login" })
      }
    }

    handleCallback()
  }, [navigate, githubMutation])

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 flex flex-col justify-center items-center p-6 text-center">
      <div className="space-y-4">
        <Loader2 className="h-8 w-8 animate-spin text-primary mx-auto" />
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">
          Authenticating with GitHub
        </h2>
        <p className="text-sm text-zinc-500 dark:text-zinc-400">
          Please wait while we complete the secure handshake...
        </p>
      </div>
    </div>
  )
}
