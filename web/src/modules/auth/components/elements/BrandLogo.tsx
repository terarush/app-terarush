import React from "react"
import { Link } from "@tanstack/react-router"

export function BrandLogo() {
  return (
    <Link to="/" className="inline-flex items-center space-x-3 group">
      <div className="h-11 w-11 rounded-lg bg-[--brand-green] text-white flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-300">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" className="h-6 w-6">
          <path d="M4 4h16" />
          <path d="M12 4v16" />
          <path d="M12 11h4a3 3 0 0 0 0-6" />
          <path d="M14 11.5l5 8.5" />
        </svg>
      </div>
      <span className="text-3xl font-bold text-foreground">
        TeraRush
      </span>
    </Link>
  )
}
