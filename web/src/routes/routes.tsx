import { createRootRoute, createRoute, Outlet, Link } from '@tanstack/react-router'
import { TooltipProvider } from '@/components/ui/tooltip'
import { Button } from '@/components/ui/button'
import LoginPage from '@/modules/auth/login'
import RegisterPage from '@/modules/auth/register'

import '../styles.css'

export const rootRoute = createRootRoute({
  component: RootComponent,
})

function RootComponent() {
  return (
    <TooltipProvider>
      <div className="min-h-screen bg-background text-foreground font-sans antialiased">
        <Outlet />
      </div>
    </TooltipProvider>
  )
}

function Home() {
  return (
    <div className="p-8">
      <div className="max-w-md mx-auto space-y-6 text-center mt-20">
        <h1 className="text-4xl font-extrabold tracking-tight">TeraRush React Starter</h1>
        <p className="text-muted-foreground text-sm">
          A clean, programmatic routing setup using TanStack Router, Tailwind CSS v4, and shadcn/ui.
        </p>
        <div className="flex justify-center gap-4">
          <Link to="/login">
            <Button className="cursor-pointer">Sign In</Button>
          </Link>
          <Link to="/register">
            <Button variant="outline" className="cursor-pointer">Sign Up</Button>
          </Link>
        </div>
      </div>
    </div>
  )
}

export const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: Home,
})

export const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  component: LoginPage,
})

export const registerRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/register',
  component: RegisterPage,
})

export const routeTree = rootRoute.addChildren([indexRoute, loginRoute, registerRoute])

