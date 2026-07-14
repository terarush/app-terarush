import React, { useState, useEffect } from "react"
import { Mail, Lock, User, Check, ArrowRight } from "lucide-react"
import { Link, useNavigate } from "@tanstack/react-router"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { AuthField } from "./fragments/AuthField"
import { SocialButton } from "./elements/SocialButton"
import { Divider } from "./fragments/Divider"
import { authContent } from "../content/auth"

export function RegisterForm() {
  const navigate = useNavigate()
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [agreeTerms, setAgreeTerms] = useState(false)
  const [passwordStrength, setPasswordStrength] = useState(0)
  const [isLoading, setIsLoading] = useState(false)
  const [errors, setErrors] = useState<{ name?: string; email?: string; password?: string; confirmPassword?: string; agreeTerms?: string }>({})
  const [errorMessage, setErrorMessage] = useState("")

  // Calculate password strength
  useEffect(() => {
    let strength = 0
    if (password.length >= 8) strength++
    if (/[a-z]/.test(password) && /[A-Z]/.test(password)) strength++
    if (/\d/.test(password)) strength++
    if (/[^a-zA-Z0-9]/.test(password)) strength++
    setPasswordStrength(strength)
  }, [password])

  const validate = () => {
    const newErrors: typeof errors = {}
    
    if (!name.trim()) {
      newErrors.name = "Full name is required"
    } else if (name.trim().length < 2) {
      newErrors.name = "Name must be at least 2 characters"
    }

    if (!email) {
      newErrors.email = "Email is required"
    } else if (!/\S+@\S+\.\S+/.test(email)) {
      newErrors.email = "Please enter a valid email address"
    }

    if (!password) {
      newErrors.password = "Password is required"
    } else if (password.length < 8) {
      newErrors.password = "Password must be at least 8 characters"
    }

    if (!confirmPassword) {
      newErrors.confirmPassword = "Please confirm your password"
    } else if (password !== confirmPassword) {
      newErrors.confirmPassword = "Passwords do not match"
    }

    if (!agreeTerms) {
      newErrors.agreeTerms = "You must agree to the Terms and Privacy Policy"
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
      toast.success("Account created successfully!", {
        description: "Please sign in to proceed.",
      })
      navigate({ to: "/login" })
    } catch (error: any) {
      console.error("Registration error:", error)
      setErrorMessage("Registration failed. Please try again.")
    } finally {
      setIsLoading(false)
    }
  }

  const handleGitHubLogin = () => {
    toast.info("Redirecting to GitHub...", {
      description: "Connecting to OAuth authentication provider.",
    })
  }

  const getPasswordStrengthColor = () => {
    if (passwordStrength === 0) return "bg-muted"
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

  return (
    <Card className="border border-border bg-card shadow-xl rounded-2xl">
      <CardContent className="p-8">
        <form onSubmit={onSubmit} className="space-y-5">
          {/* Error Message */}
          {errorMessage && (
            <div className="p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive text-sm">
              {errorMessage}
            </div>
          )}

          {/* Name Field */}
          <AuthField
            id="name"
            label={authContent.register.nameLabel}
            type="text"
            placeholder={authContent.register.namePlaceholder}
            value={name}
            onChange={(e) => setName(e.target.value)}
            error={errors.name}
            icon={<User className="h-5 w-5" />}
            disabled={isLoading}
          />

          {/* Email Field */}
          <AuthField
            id="email"
            label={authContent.register.emailLabel}
            type="email"
            placeholder={authContent.register.emailPlaceholder}
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            error={errors.email}
            icon={<Mail className="h-5 w-5" />}
            disabled={isLoading}
          />

          {/* Password Field */}
          <div>
            <AuthField
              id="password"
              label={authContent.register.passwordLabel}
              type="password"
              placeholder={authContent.register.passwordPlaceholder}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              error={errors.password}
              icon={<Lock className="h-5 w-5" />}
              disabled={isLoading}
            />
            {/* Password Strength Indicator */}
            {password && (
              <div className="mt-2">
                <div className="flex gap-1 mb-1">
                  {[1, 2, 3, 4].map((level) => (
                    <div
                      key={level}
                      className={`h-1 flex-1 rounded-full transition-all duration-300 ${
                        level <= passwordStrength
                          ? getPasswordStrengthColor()
                          : "bg-muted"
                      }`}
                    />
                  ))}
                </div>
                <p className="text-xs text-muted-foreground">
                  Password strength: {getPasswordStrengthText()}
                </p>
              </div>
            )}
          </div>

          {/* Confirm Password Field */}
          <div>
            <AuthField
              id="confirmPassword"
              label={authContent.register.confirmPasswordLabel}
              type="password"
              placeholder={authContent.register.confirmPasswordPlaceholder}
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              error={errors.confirmPassword}
              icon={<Lock className="h-5 w-5" />}
              disabled={isLoading}
            />
            {confirmPassword && password === confirmPassword && (
              <p className="text-xs text-green-600 dark:text-green-400 mt-1 flex items-center gap-1">
                <Check className="h-3 w-3" />
                Passwords match
              </p>
            )}
          </div>

          {/* Terms Checkbox */}
          <div className="flex items-start space-x-2">
            <input
              type="checkbox"
              id="terms"
              checked={agreeTerms}
              onChange={(e) => setAgreeTerms(e.target.checked)}
              className="w-4 h-4 mt-0.5 rounded border-border text-primary focus:ring-primary focus:ring-offset-0 cursor-pointer"
              disabled={isLoading}
            />
            <label htmlFor="terms" className="text-sm text-muted-foreground cursor-pointer select-none">
              {authContent.register.termsText}{" "}
              <Link to="/login" className="text-primary hover:text-primary/80 transition-colors font-medium">
                {authContent.register.termsLink}
              </Link>{" "}
              {authContent.register.andText}{" "}
              <Link to="/login" className="text-primary hover:text-primary/80 transition-colors font-medium">
                {authContent.register.privacyLink}
              </Link>
            </label>
          </div>
          {errors.agreeTerms && (
            <p className="text-sm text-destructive">{errors.agreeTerms}</p>
          )}

          {/* Submit Button */}
          <Button
            type="submit"
            disabled={isLoading}
            className="w-full bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl h-12 font-semibold shadow-lg hover:shadow-xl hover:scale-[1.02] transition-all duration-300 cursor-pointer"
          >
            {isLoading ? (
              authContent.register.submittingButton
            ) : (
              <span className="flex items-center justify-center">
                {authContent.register.submitButton}
                <ArrowRight className="ml-2 h-5 w-5" />
              </span>
            )}
          </Button>
        </form>

        {/* Divider */}
        <Divider>{authContent.register.dividerText}</Divider>

        {/* Social Login */}
        <div className="space-y-3">
          <SocialButton
            provider="github"
            onClick={handleGitHubLogin}
          >
            {authContent.register.githubButton}
          </SocialButton>
        </div>

        {/* Sign In Link */}
        <div className="text-center mt-6">
          <p className="text-muted-foreground">
            {authContent.register.hasAccountText}{" "}
            <Link
              to="/login"
              className="text-primary hover:text-primary/80 font-semibold transition-colors"
            >
              {authContent.register.signInLink}
            </Link>
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
