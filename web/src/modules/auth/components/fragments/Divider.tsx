import React from "react"

interface DividerProps {
  children?: React.ReactNode
}

export function Divider({ children }: DividerProps) {
  return (
    <div className="relative my-5 flex items-center justify-center">
      <div className="absolute inset-0 flex items-center">
        <div className="w-full border-t border-zinc-200 dark:border-zinc-800"></div>
      </div>
      {children && (
        <span className="relative px-3 bg-white dark:bg-zinc-900 text-[10px] font-bold text-zinc-400 dark:text-zinc-500 uppercase tracking-widest">
          {children}
        </span>
      )}
    </div>
  )
}
