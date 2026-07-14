import React, { useState } from "react"
import { Mail, Lock, ArrowRight } from "lucide-react"
import { Link, useNavigate } from "@tanstack/react-router"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { AuthField } from "./fragments/AuthField"
import { SocialButton } from "./elements/SocialButton"
import { Divider } from "./fragments/Divider"
import { authContent } from "../content/auth"

export function LoginForm() {
  const navigate = useNavigate()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [rememberMe, setRememberMe] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [errors, setErrors] = useState<{ email?: string; password?: string }>({})
  const [errorMessage, setErrorMessage] = useState("")

  const validate = () => {
    const newErrors: typeof errors = {}
    if (!email) {
      newErrors.email = "Email is required"
    } else if (!/\S+@\S+\.\S+/.test(email)) {
      newErrors.email = "Please enter a valid email address"
    }

    if (!password) {
      newErrors.password = "Password is required"
    } else if (password.length < 6) {
      newErrors.password = "Password must be at least 6 characters"
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrorMessage("")
    
    if (!validate()) return

    setIsLoading(true)
    try {
      await new Promise((resolve) => setTimeout(resolve, 1500))
      toast.success("Welcome back!", {
        description: "You have successfully signed in.",
      })
      navigate({ to: "/" })
    } catch (error: any) {
      console.error("Login error:", error)
      setErrorMessage("Login failed. Please check your credentials.")
    } finally {
      setIsLoading(false)
    }
  }

  const handleGitHubLogin = () => {
    toast.info("Redirecting to GitHub...", {
      description: "Connecting to OAuth authentication provider.",
    })
  }

  return (
    <Card className="border border-border bg-card shadow-xl rounded-2xl">
      <CardContent className="p-8">
        <form onSubmit={onSubmit} className="space-y-6">
          {/* Error Message */}
          {errorMessage && (
            <div className="p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive text-sm">
              {errorMessage}
            </div>
          )}

          {/* Email Field */}
          <AuthField
            id="email"
            label={authContent.login.emailLabel}
            type="email"
            placeholder={authContent.login.emailPlaceholder}
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            error={errors.email}
            icon={<Mail className="h-5 w-5" />}
            disabled={isLoading}
          />

          {/* Password Field */}
          <AuthField
            id="password"
            label={authContent.login.passwordLabel}
            type="password"
            placeholder={authContent.login.passwordPlaceholder}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            error={errors.password}
            icon={<Lock className="h-5 w-5" />}
            disabled={isLoading}
          />

          {/* Remember Me & Forgot Password */}
          <div className="flex items-center justify-between">
            <label className="flex items-center space-x-2 cursor-pointer select-none">
              <input
                type="checkbox"
                checked={rememberMe}
                onChange={(e) => setRememberMe(e.target.checked)}
                className="w-4 h-4 rounded border-border text-primary focus:ring-primary focus:ring-offset-0 cursor-pointer"
                disabled={isLoading}
              />
              <span className="text-sm text-muted-foreground">
                {authContent.login.rememberMe}
              </span>
            </label>
            <Link
              to="/login"
              className="text-sm text-primary hover:text-primary/80 transition-colors font-medium"
            >
              {authContent.login.forgotPassword}
            </Link>
          </div>

          {/* Submit Button */}
          <Button
            type="submit"
            disabled={isLoading}
            className="w-full bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl h-12 font-semibold shadow-lg hover:shadow-xl hover:scale-[1.02] transition-all duration-300 cursor-pointer"
          >
            {isLoading ? (
              authContent.login.submittingButton
            ) : (
              <span className="flex items-center justify-center">
                {authContent.login.submitButton}
                <ArrowRight className="ml-2 h-5 w-5" />
              </span>
            )}
          </Button>
        </form>

        {/* Divider */}
        <Divider>{authContent.login.dividerText}</Divider>

        {/* Social Login */}
        <div className="space-y-3">
          <SocialButton
            provider="github"
            onClick={handleGitHubLogin}
          >
            {authContent.login.githubButton}
          </SocialButton>
        </div>

        {/* Sign Up Link */}
        <div className="text-center mt-6">
          <p className="text-muted-foreground">
            {authContent.login.noAccountText}{" "}
            <Link
              to="/register"
              className="text-primary hover:text-primary/80 font-semibold transition-colors"
            >
              {authContent.login.signUpLink}
            </Link>
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
