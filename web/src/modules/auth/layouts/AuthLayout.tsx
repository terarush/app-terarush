import React from "react"
import { Link } from "@tanstack/react-router"
import { BrandLogo } from "../components/elements/BrandLogo"
import { authContent } from "../content/auth"
import { companyMeta } from "@/meta"

interface AuthLayoutProps {
  children: React.ReactNode
}

export function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="min-h-screen w-full grid grid-cols-1 md:grid-cols-12 bg-white dark:bg-zinc-950 transition-colors duration-200">

      <div className="hidden md:flex md:col-span-5 lg:col-span-4 bg-primary text-primary-foreground flex flex-col justify-between p-12 relative border-r border-zinc-200 dark:border-zinc-800/20">

        <div>
          <Link to="/" className="inline-flex items-center space-x-2.5 group">
            <span className="text-xl font-bold tracking-tight text-primary-foreground">
              {companyMeta.name}
            </span>
          </Link>
        </div>

        <div className="space-y-6 my-auto">
          <h2 className="text-3xl lg:text-4xl font-semibold tracking-tight leading-tight font-serif text-primary-foreground">
            Build Your Digital Future with Expert Development.
          </h2>
          <p className="text-sm text-primary-foreground/80 leading-relaxed font-normal">
            Professional web development, mobile apps, backend solutions, and bug fixes. From concept to deployment, we deliver exceptional results.
          </p>

          <div className="pt-6 border-t border-white/10 space-y-3">
            <p className="text-xs italic leading-relaxed text-primary-foreground/90">
              &quot;To empower businesses and developers by delivering cutting-edge web and application development services, solving technical challenges, and creating robust digital solutions.&quot;
            </p>
          </div>
        </div>

        <div className="text-[10px] text-primary-foreground/40 font-medium">
          &copy; {new Date().getFullYear()} {companyMeta.name}. All rights reserved.
        </div>
      </div>

      <div className="col-span-1 md:col-span-7 lg:col-span-8 bg-zinc-50 dark:bg-zinc-950 flex flex-col justify-center items-center py-12 px-6 sm:px-12 lg:px-24 relative">

        <div className="absolute top-6 left-6 md:left-8">
          <Link
            to="/"
            className="inline-flex items-center gap-1.5 text-xs font-semibold text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200 transition-colors"
          >
            &larr; {authContent.backToHome}
          </Link>
        </div>

        <div className="md:hidden mb-8 text-center">
          <div className="flex justify-center mb-4">
            <BrandLogo />
          </div>
        </div>

        <div className="w-full max-w-md">
          <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800/80 p-8 shadow-xs rounded-xl">
            {children}
          </div>
        </div>

      </div>

    </div>
  )
}
