import React from "react"

interface DividerProps {
  children?: React.ReactNode
}

export function Divider({ children = "Or continue with" }: DividerProps) {
  return (
    <div className="relative my-6">
      <div className="absolute inset-0 flex items-center">
        <div className="w-full border-t border-border"></div>
      </div>
      <div className="relative flex justify-center text-sm">
        <span className="px-4 bg-card text-muted-foreground">
          {children}
        </span>
      </div>
    </div>
  )
}
