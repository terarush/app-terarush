import React from "react"
import { Link } from "@tanstack/react-router"
import { BrandLogo } from "../components/elements/BrandLogo"
import { authContent } from "../content/auth"
import { companyMeta } from "@/meta"

interface AuthLayoutProps {
  children: React.ReactNode
  title: string
  subtitle: string
}

export function AuthLayout({ children, title, subtitle }: AuthLayoutProps) {
  return (
    <div className="min-h-screen w-full grid grid-cols-1 md:grid-cols-12 bg-white dark:bg-zinc-950 transition-colors duration-200">

      <div className="hidden md:flex md:col-span-5 lg:col-span-4 bg-primary text-primary-foreground flex flex-col justify-between p-12 relative border-r border-zinc-200 dark:border-zinc-800/20">

        <div>
          <Link to="/" className="inline-flex items-center space-x-2.5 group">
            <div className="h-9 w-9 rounded-md bg-primary-foreground text-primary flex items-center justify-center font-bold tracking-tighter text-sm">
              TR
            </div>
            <span className="text-xl font-bold tracking-tight text-primary-foreground">
              {companyMeta.name}
            </span>
          </Link>
        </div>

        <div className="space-y-6 my-auto">
          <span className="text-[10px] font-bold tracking-widest uppercase bg-white/10 text-primary-foreground/90 px-2.5 py-1 rounded-md border border-white/5">
            Custom Web & App Development Agency
          </span>
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
            <p className="text-[10px] font-bold text-primary-foreground/75 tracking-wider uppercase">— Our Mission</p>
          </div>
        </div>

        <div className="text-[10px] text-primary-foreground/40 font-medium">
          &copy; {new Date().getFullYear()} {companyMeta.name}. All rights reserved.
        </div>
      </div>

      <div className="col-span-1 md:col-span-7 lg:col-span-8 bg-zinc-50 dark:bg-zinc-950 flex flex-col justify-center items-center py-12 px-6 sm:px-12 lg:px-24">

        <div className="md:hidden mb-8 text-center">
          <div className="flex justify-center mb-4">
            <BrandLogo />
          </div>
        </div>

        <div className="w-full max-w-md space-y-6">
          <div className="text-center md:text-left space-y-1.5">
            <h1 className="text-2xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50">
              {title}
            </h1>
            <p className="text-sm text-zinc-500 dark:text-zinc-400">
              {subtitle}
            </p>
          </div>

          <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800/80 p-8 shadow-xs rounded-xl">
            {children}
          </div>

          <div className="text-center md:text-left">
            <Link
              to="/"
              className="text-xs font-semibold text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200 transition-colors"
            >
              {authContent.backToHome}
            </Link>
          </div>
        </div>

      </div>

    </div>
  )
}
