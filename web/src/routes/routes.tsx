import { createRootRoute, createRoute, Outlet, Link } from '@tanstack/react-router'
import { TooltipProvider } from '@/components/ui/tooltip'
import { Button } from '@/components/ui/button'
import { ThemeProvider } from '@/components/theme-provider'
import { AuthProvider } from '@/contexts/auth-context'
import LoginPage from '@/modules/auth/login'
import RegisterPage from '@/modules/auth/register'
import GitHubCallbackPage from '@/modules/auth/github-callback'

import '../styles.css'

import { Toaster } from '@/components/ui/sonner'

export const rootRoute = createRootRoute({
  component: RootComponent,
})

function RootComponent() {
  return (
    <AuthProvider>
      <ThemeProvider defaultTheme="system" storageKey="vite-ui-theme">
        <TooltipProvider>
          <div className="min-h-screen bg-background text-foreground font-sans antialiased">
            <Outlet />
            <Toaster />
          </div>
        </TooltipProvider>
      </ThemeProvider>
    </AuthProvider>
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

const routes = [
  { path: '/', component: Home },
  { path: '/login', component: LoginPage },
  { path: '/register', component: RegisterPage },
  { path: '/auth/github/callback', component: GitHubCallbackPage },
]

export const routeTree = rootRoute.addChildren(
  routes.map((r) =>
    createRoute({
      getParentRoute: () => rootRoute,
      path: r.path,
      component: r.component,
    })
  )
)
