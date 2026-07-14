import { createRootRoute, createRoute, Outlet } from '@tanstack/react-router'
import { TooltipProvider } from '@/components/ui/tooltip'
import { Button } from '@/components/ui/button'

import '../styles.css'

export const rootRoute = createRootRoute({
  component: RootComponent,
})

function RootComponent() {
  return (
    <TooltipProvider>
      <div className="min-h-screen bg-background text-foreground font-sans antialiased p-8">
        <Outlet />
      </div>
    </TooltipProvider>
  )
}

function Home() {
  return (
    <div className="max-w-md mx-auto space-y-4 text-center mt-20">
      <h1 className="text-3xl font-extrabold tracking-tight">TeraRush React Starter</h1>
      <p className="text-muted-foreground text-sm">
        A clean, programmatic routing setup using TanStack Router, Tailwind CSS v4, and shadcn/ui.
      </p>
      <Button>Get Started</Button>
    </div>
  )
}

export const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: Home,
})

export const routeTree = rootRoute.addChildren([indexRoute])
