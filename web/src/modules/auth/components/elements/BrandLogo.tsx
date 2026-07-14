import React from "react"
import { Link } from "@tanstack/react-router"
import { companyMeta } from "@/meta"

export function BrandLogo() {
  return (
    <Link to="/" className="inline-flex items-center space-x-2.5 group">
      <div className="h-9 w-9 rounded-md bg-zinc-900 dark:bg-zinc-50 text-white dark:text-zinc-950 flex items-center justify-center font-bold tracking-tighter text-sm">
        TR
      </div>
      <span className="text-xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50">
        {companyMeta.name}
      </span>
    </Link>
  )
}
