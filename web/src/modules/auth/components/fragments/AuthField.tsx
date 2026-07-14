import React, { useState } from "react"
import { Eye, EyeOff } from "lucide-react"

interface AuthFieldProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string
  error?: string
}

export function AuthField({ label, error, type = "text", className, id, ...props }: AuthFieldProps) {
  const [showPassword, setShowPassword] = useState(false)
  const isPassword = type === "password"
  const inputType = isPassword ? (showPassword ? "text" : "password") : type

  return (
    <div className="space-y-1.5 text-left">
      <label htmlFor={id} className="block text-xs font-medium text-zinc-700 dark:text-zinc-300">
        {label}
      </label>
      <div className="relative">
        <input
          id={id}
          type={inputType}
          className={`w-full bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800/80 rounded-md px-3 py-2.5 text-sm h-10 transition-all duration-200 shadow-xs placeholder:text-zinc-400 focus:bg-white dark:focus:bg-zinc-900 focus:border-zinc-950 dark:focus:border-zinc-300 focus:ring-3 focus:ring-zinc-950/5 dark:focus:ring-zinc-100/5 focus:outline-hidden ${isPassword ? "pr-10" : ""}`}
          {...props}
        />
        {isPassword && (
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-zinc-500 hover:text-zinc-600 dark:hover:text-zinc-300 focus:outline-hidden transition-colors cursor-pointer"
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
        <p className="text-xs text-red-500 font-medium mt-1">{error}</p>
      )}
    </div>
  )
}
