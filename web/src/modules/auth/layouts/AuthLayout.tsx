import React, { useState, useEffect } from "react"
import { Link } from "@tanstack/react-router"
import { BrandLogo } from "../components/elements/BrandLogo"
import { authContent } from "../content/auth"

interface AuthLayoutProps {
  children: React.ReactNode
  title: string
  subtitle: string
}

export function AuthLayout({ children, title, subtitle }: AuthLayoutProps) {
  const [isMounted, setIsMounted] = useState(false)

  useEffect(() => {
    setIsMounted(true)
  }, [])

  return (
    <div className="min-h-screen bg-background flex items-center justify-center relative overflow-hidden p-4">
      {/* Background decorations */}
      <div className="absolute inset-0 overflow-hidden opacity-30 pointer-events-none">
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-chart-1/10 rounded-full blur-3xl" />
      </div>

      <div className={`w-full max-w-md relative z-10 transition-all duration-700 transform ${isMounted ? 'translate-y-0 opacity-100' : 'translate-y-4 opacity-0'}`}>
        {/* Logo and Title */}
        <div className="text-center mb-8">
          <div className="mb-6 flex justify-center">
            <BrandLogo />
          </div>
          <h1 className="text-3xl font-bold text-foreground mb-2">
            {title}
          </h1>
          <p className="text-muted-foreground">
            {subtitle}
          </p>
        </div>

        {/* Active form card */}
        {children}

        {/* Back to Home */}
        <div className="text-center mt-6">
          <Link
            to="/"
            className="inline-flex items-center text-muted-foreground hover:text-foreground transition-colors font-medium"
          >
            {authContent.backToHome}
          </Link>
        </div>
      </div>
    </div>
  )
}
