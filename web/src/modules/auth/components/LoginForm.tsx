import React, { useState } from "react"
import { Link, useNavigate } from "@tanstack/react-router"
import { toast } from "sonner"

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
      navigate({ to: "/app" })
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
    <div className="w-full">
      <div className="text-left space-y-1 mb-6">
        <h1 className="text-2xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50">
          {authContent.login.title}
        </h1>
        <p className="text-xs text-zinc-500 dark:text-zinc-400">
          {authContent.login.subtitle}
        </p>
      </div>

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

        <button
          type="submit"
          disabled={isPending}
          className="w-full bg-zinc-950 hover:bg-zinc-900 dark:bg-zinc-50 dark:hover:bg-zinc-200 text-zinc-50 dark:text-zinc-950 rounded-md h-10 text-sm font-semibold transition-all duration-200 cursor-pointer flex items-center justify-center border-0 shadow-sm focus:outline-hidden focus:ring-3 focus:ring-zinc-950/10 dark:focus:ring-zinc-100/15 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isPending ? (
            authContent.login.submittingButton
          ) : (
            authContent.login.submitButton
          )}
        </button>
      </form>

      <div className="mt-5 space-y-5">
        <Divider>{authContent.login.dividerText}</Divider>

        <SocialButton
          provider="github"
          onClick={handleGitHubLogin}
        >
          {authContent.login.githubButton}
        </SocialButton>

        <div className="text-center pt-1">
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
    </div>
  )
}
