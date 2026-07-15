import React, { useState, useEffect } from "react"
import { Check } from "lucide-react"
import { Link, useNavigate } from "@tanstack/react-router"
import { toast } from "sonner"

import { AuthField } from "./fragments/AuthField"
import { SocialButton } from "./elements/SocialButton"
import { Divider } from "./fragments/Divider"
import { authContent } from "../content/auth"
import { registerSchema } from "@/schemas/auth"
import { useRegisterMutation } from "@/service/mutation/auth"
import { authApi } from "@/service/api/auth"

export function RegisterForm() {
  const registerMutation = useRegisterMutation()
  const navigate = useNavigate()
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [agreeTerms, setAgreeTerms] = useState(false)
  const [passwordStrength, setPasswordStrength] = useState(0)
  const [errors, setErrors] = useState<{ name?: string; email?: string; password?: string; confirmPassword?: string; agreeTerms?: string }>({})
  const [errorMessage, setErrorMessage] = useState("")

  useEffect(() => {
    let strength = 0
    if (password.length >= 8) strength++
    if (/[a-z]/.test(password) && /[A-Z]/.test(password)) strength++
    if (/\d/.test(password)) strength++
    if (/[^a-zA-Z0-9]/.test(password)) strength++
    setPasswordStrength(strength)
  }, [password])

  const validate = () => {
    const result = registerSchema.safeParse({ name, email, password, confirmPassword })
    const fieldErrors: typeof errors = {}
    
    if (!result.success) {
      result.error.issues.forEach((issue) => {
        const path = issue.path[0] as keyof typeof errors
        fieldErrors[path] = issue.message
      })
    }
    
    if (!agreeTerms) {
      fieldErrors.agreeTerms = "You must agree to the Terms and Privacy Policy"
    }
    
    setErrors(fieldErrors)
    return result.success && agreeTerms
  }

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrorMessage("")
    
    if (!validate()) return

    try {
      await registerMutation.mutateAsync({
        name,
        email,
        password,
        confirm_password: confirmPassword,
      })
      toast.success("Account created successfully!", {
        description: "Please sign in to proceed.",
      })
      navigate({ to: "/login" })
    } catch (error: any) {
      console.error("Registration error:", error)
      const errorMsg = error?.response?.data?.error || "Registration failed. Please try again."
      setErrorMessage(errorMsg)
      toast.error("Registration failed", {
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

  const getPasswordStrengthColor = () => {
    if (passwordStrength === 0) return "bg-zinc-200 dark:bg-zinc-800"
    if (passwordStrength === 1) return "bg-red-500"
    if (passwordStrength === 2) return "bg-orange-500"
    if (passwordStrength === 3) return "bg-yellow-500"
    return "bg-green-500"
  }

  const getPasswordStrengthText = () => {
    if (passwordStrength === 0) return ""
    if (passwordStrength === 1) return "Weak"
    if (passwordStrength === 2) return "Fair"
    if (passwordStrength === 3) return "Good"
    return "Strong"
  }

  const isPending = registerMutation.isPending

  return (
    <div className="w-full">
      <div className="text-left space-y-1 mb-6">
        <h1 className="text-2xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50">
          {authContent.register.title}
        </h1>
        <p className="text-xs text-zinc-500 dark:text-zinc-400">
          {authContent.register.subtitle}
        </p>
      </div>

      <form onSubmit={onSubmit} className="space-y-4">
        {errorMessage && (
          <div className="p-3 rounded-md bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-900/30 text-red-600 dark:text-red-400 text-xs font-medium">
            {errorMessage}
          </div>
        )}

        <AuthField
          id="name"
          label={authContent.register.nameLabel}
          type="text"
          placeholder={authContent.register.namePlaceholder}
          value={name}
          onChange={(e) => setName(e.target.value)}
          error={errors.name}
          disabled={isPending}
        />

        <AuthField
          id="email"
          label={authContent.register.emailLabel}
          type="email"
          placeholder={authContent.register.emailPlaceholder}
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          error={errors.email}
          disabled={isPending}
        />

        <div>
          <AuthField
            id="password"
            label={authContent.register.passwordLabel}
            type="password"
            placeholder={authContent.register.passwordPlaceholder}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            error={errors.password}
            disabled={isPending}
          />
          {password && (
            <div className="mt-2">
              <div className="flex gap-1 mb-1">
                {[1, 2, 3, 4].map((level) => (
                  <div
                    key={level}
                    className={`h-1 flex-1 rounded-full transition-all duration-300 ${
                      level <= passwordStrength
                        ? getPasswordStrengthColor()
                        : "bg-zinc-200 dark:bg-zinc-800"
                    }`}
                  />
                ))}
              </div>
              <p className="text-[10px] text-zinc-500 dark:text-zinc-400 font-medium">
                Password strength: {getPasswordStrengthText()}
              </p>
            </div>
          )}
        </div>

        <div>
          <AuthField
            id="confirmPassword"
            label={authContent.register.confirmPasswordLabel}
            type="password"
            placeholder={authContent.register.confirmPasswordPlaceholder}
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            error={errors.confirmPassword}
            disabled={isPending}
          />
          {confirmPassword && password === confirmPassword && (
            <p className="text-[10px] text-green-600 dark:text-green-400 mt-1 flex items-center gap-1 font-medium">
              <Check className="h-3.5 w-3.5" />
              Passwords match
            </p>
          )}
        </div>

        <div className="flex items-start space-x-2 pt-1">
          <input
            type="checkbox"
            id="terms"
            checked={agreeTerms}
            onChange={(e) => setAgreeTerms(e.target.checked)}
            className="w-3.5 h-3.5 mt-0.5 rounded border-zinc-300 dark:border-zinc-700 text-zinc-900 focus:ring-zinc-950 cursor-pointer"
            disabled={isPending}
          />
          <label htmlFor="terms" className="text-xs text-zinc-500 dark:text-zinc-400 cursor-pointer select-none font-medium leading-normal">
            {authContent.register.termsText}{" "}
            <Link to="/login" className="text-zinc-900 dark:text-zinc-100 hover:underline font-semibold">
              {authContent.register.termsLink}
            </Link>{" "}
            {authContent.register.andText}{" "}
            <Link to="/login" className="text-zinc-900 dark:text-zinc-100 hover:underline font-semibold">
              {authContent.register.privacyLink}
            </Link>
          </label>
        </div>
        {errors.agreeTerms && (
          <p className="text-xs text-red-500 font-medium">{errors.agreeTerms}</p>
        )}

        <button
          type="submit"
          disabled={isPending}
          className="w-full bg-zinc-950 hover:bg-zinc-900 dark:bg-zinc-50 dark:hover:bg-zinc-200 text-zinc-50 dark:text-zinc-950 rounded-md h-10 text-sm font-semibold transition-all duration-200 cursor-pointer flex items-center justify-center border-0 shadow-sm focus:outline-hidden focus:ring-3 focus:ring-zinc-950/10 dark:focus:ring-zinc-100/15 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isPending ? (
            authContent.register.submittingButton
          ) : (
            authContent.register.submitButton
          )}
        </button>
      </form>

      <div className="mt-5 space-y-5">
        <Divider>{authContent.register.dividerText}</Divider>

        <SocialButton
          provider="github"
          onClick={handleGitHubLogin}
        >
          {authContent.register.githubButton}
        </SocialButton>

        <div className="text-center pt-1">
          <p className="text-zinc-500 dark:text-zinc-400 text-xs">
            {authContent.register.hasAccountText}{" "}
            <Link
              to="/login"
              className="font-semibold text-primary hover:underline"
            >
              {authContent.register.signInLink}
            </Link>
          </p>
        </div>
      </div>
    </div>
  )
}
