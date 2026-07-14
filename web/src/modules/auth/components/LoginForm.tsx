import React, { useState } from "react"
import { Mail, Lock } from "lucide-react"
import { Link, useNavigate } from "@tanstack/react-router"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import { AuthField } from "./fragments/AuthField"
import { SocialButton } from "./elements/SocialButton"
import { Divider } from "./fragments/Divider"
import { authContent } from "../content/auth"
import { loginSchema } from "@/schemas/auth"
import { useLoginMutation } from "@/service/mutation/auth"
import { authApi } from "@/service/api/auth"

export function LoginForm() {
  const loginMutation = useLoginMutation()
  const navigate = useNavigate()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [rememberMe, setRememberMe] = useState(false)
  const [errors, setErrors] = useState<{ email?: string; password?: string }>({})
  const [errorMessage, setErrorMessage] = useState("")

  const validate = () => {
    const result = loginSchema.safeParse({ email, password })
    if (!result.success) {
      const fieldErrors: typeof errors = {}
      result.error.issues.forEach((issue) => {
        const path = issue.path[0] as keyof typeof errors
        fieldErrors[path] = issue.message
      })
      setErrors(fieldErrors)
      return false
    }
    setErrors({})
    return true
  }

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrorMessage("")
    
    if (!validate()) return

    try {
      await loginMutation.mutateAsync({ email, password })
      toast.success("Welcome back!", {
        description: "You have successfully signed in.",
      })
      navigate({ to: "/" })
    } catch (error: any) {
      console.error("Login error:", error)
      const errorMsg = error?.response?.data?.error || "Login failed. Please check your credentials."
      setErrorMessage(errorMsg)
      toast.error("Authentication failed", {
        description: errorMsg,
      })
    }
  }

  const handleGitHubLogin = () => {
    toast.info("Redirecting to GitHub...", {
      description: "Connecting to OAuth authentication provider.",
    })
    window.location.href = authApi.getGitHubAuthUrl()
  }

  const isPending = loginMutation.isPending

  return (
    <div className="w-full space-y-5">
      <form onSubmit={onSubmit} className="space-y-4">
        {errorMessage && (
          <div className="p-3 rounded-md bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-900/30 text-red-600 dark:text-red-400 text-xs font-medium">
            {errorMessage}
          </div>
        )}

        <AuthField
          id="email"
          label={authContent.login.emailLabel}
          type="email"
          placeholder={authContent.login.emailPlaceholder}
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          error={errors.email}
          icon={<Mail className="h-4 w-4" />}
          disabled={isPending}
        />

        <AuthField
          id="password"
          label={authContent.login.passwordLabel}
          type="password"
          placeholder={authContent.login.passwordPlaceholder}
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          error={errors.password}
          icon={<Lock className="h-4 w-4" />}
          disabled={isPending}
        />

        <div className="flex items-center justify-between text-xs pt-1">
          <label className="flex items-center space-x-2 cursor-pointer select-none">
            <input
              type="checkbox"
              checked={rememberMe}
              onChange={(e) => setRememberMe(e.target.checked)}
              className="w-3.5 h-3.5 rounded border-zinc-300 dark:border-zinc-700 text-zinc-900 focus:ring-zinc-950 cursor-pointer"
              disabled={isPending}
            />
            <span className="text-zinc-600 dark:text-zinc-400 font-medium">
              {authContent.login.rememberMe}
            </span>
          </label>
          <Link
            to="/login"
            className="text-zinc-900 dark:text-zinc-100 hover:underline font-semibold"
          >
            {authContent.login.forgotPassword}
          </Link>
        </div>

        <Button
          type="submit"
          disabled={isPending}
          className="w-full bg-primary hover:bg-primary/90 text-primary-foreground rounded-lg h-10 text-sm font-semibold transition-colors cursor-pointer flex items-center justify-center border border-transparent"
        >
          {isPending ? (
            authContent.login.submittingButton
          ) : (
            authContent.login.submitButton
          )}
        </Button>
      </form>

      <Divider>{authContent.login.dividerText}</Divider>

      <SocialButton
        provider="github"
        onClick={handleGitHubLogin}
      >
        {authContent.login.githubButton}
      </SocialButton>

      <div className="text-center pt-2">
        <p className="text-zinc-500 dark:text-zinc-400 text-xs">
          {authContent.login.noAccountText}{" "}
          <Link
            to="/register"
            className="font-semibold text-primary hover:underline"
          >
            {authContent.login.signUpLink}
          </Link>
        </p>
      </div>
    </div>
  )
}
