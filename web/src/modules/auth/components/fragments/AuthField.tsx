import React, { useState } from "react"
import { Eye, EyeOff } from "lucide-react"
import { Input } from "@/components/ui/input"

interface AuthFieldProps extends React.ComponentProps<typeof Input> {
  label: string
  error?: string
  icon?: React.ReactNode
}

export function AuthField({ label, error, icon, type = "text", className, id, ...props }: AuthFieldProps) {
  const [showPassword, setShowPassword] = useState(false)
  const isPassword = type === "password"
  const inputType = isPassword ? (showPassword ? "text" : "password") : type

  return (
    <div className="space-y-1.5 text-left">
      <label htmlFor={id} className="block text-xs font-semibold text-zinc-700 dark:text-zinc-300">
        {label}
      </label>
      <div className="relative">
        {icon && (
          <div className="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-zinc-500 pointer-events-none">
            {icon}
          </div>
        )}
        <Input
          id={id}
          type={inputType}
          className={`${icon ? "pl-9" : "pl-3"} ${isPassword ? "pr-9" : "pr-3"} rounded-lg border border-zinc-200 dark:border-zinc-800 bg-transparent text-sm h-10 w-full focus-visible:border-zinc-900 focus-visible:ring-1 focus-visible:ring-zinc-900/10 dark:focus-visible:border-zinc-100 dark:focus-visible:ring-zinc-100/10`}
          {...props}
        />
        {isPassword && (
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-zinc-500 hover:text-zinc-600 dark:hover:text-zinc-300 focus:outline-hidden transition-colors"
          >
            {showPassword ? (
              <EyeOff className="h-4.5 w-4.5" />
            ) : (
              <Eye className="h-4.5 w-4.5" />
            )}
          </button>
        )}
      </div>
      {error && (
        <p className="text-xs text-red-500 font-medium">{error}</p>
      )}
    </div>
  )
}
