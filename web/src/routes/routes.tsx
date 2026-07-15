import { createRootRoute, createRoute, Outlet, redirect } from '@tanstack/react-router'
import { TooltipProvider } from '@/components/ui/tooltip'
import { ThemeProvider } from '@/components/theme-provider'
import { AuthProvider } from '@/contexts/auth-context'
import LoginPage from '@/modules/auth/login'
import RegisterPage from '@/modules/auth/register'
import GitHubCallbackPage from '@/modules/auth/github-callback'
import NotFound from '@/modules/error/not-found'
import { GlobalsAppLayout } from '@/modules/app/layouts/globalsAppLayout'
import AppPage from '@/modules/app/index'

import '../styles.css'

import { Toaster } from '@/components/ui/sonner'

export const rootRoute = createRootRoute({
  component: RootComponent,
  notFoundComponent: NotFound,
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

const routeConfig = {
  root: [
    {
      path: '/',
      beforeLoad: () => {
        throw redirect({ to: '/app' })
      },
    },
    { path: '/login', component: LoginPage },
    { path: '/register', component: RegisterPage },
    { path: '/auth/github/callback', component: GitHubCallbackPage },
  ],
  app: {
    layout: GlobalsAppLayout,
    children: [
      { path: '/', component: AppPage },
    ],
  },
}

export const routeTree = rootRoute.addChildren([
  ...routeConfig.root.map((r) =>
    createRoute({
      getParentRoute: () => rootRoute,
      path: r.path,
      component: r.component,
      beforeLoad: r.beforeLoad,
    })
  ),
  (() => {
    const layout = createRoute({
      getParentRoute: () => rootRoute,
      path: '/app',
      component: routeConfig.app.layout,
    })
    return layout.addChildren(
      routeConfig.app.children.map((r) =>
        createRoute({
          getParentRoute: () => layout,
          path: r.path,
          component: r.component,
          beforeLoad: r.beforeLoad,
        })
      )
    )
  })(),
])
