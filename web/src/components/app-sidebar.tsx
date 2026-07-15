import { useRouterState, useNavigate } from "@tanstack/react-router"
import { LogOut, Moon, Sun } from "lucide-react"

import { cn } from "@/lib/utils"
import { useAuth } from "@/contexts/auth-context"
import { useTheme } from "@/components/theme-provider"
import { companyMeta } from "@/meta"
import { sidebarContentList } from "@/globals/content/app-sidebar"

import {
  Sidebar,
  SidebarHeader,
  SidebarContent,
  SidebarFooter,
  SidebarSeparator,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarMenuBadge,
  useSidebar,
} from "@/components/ui/sidebar"

export function AppSidebar() {
  const routerState = useRouterState()
  const currentPath = routerState.location.pathname
  const navigate = useNavigate()
  const { user, logout } = useAuth()
  const { isMobile, setOpenMobile } = useSidebar()
  const { theme, setTheme } = useTheme()

  const isAdmin = user?.role === "admin"

  const toggleTheme = () => {
    setTheme(theme === "dark" ? "light" : "dark")
  }

  const handleLogout = async () => {
    await logout()
    navigate({ to: "/" })
  }

  const handleItemClick = (href: string) => {
    if (isMobile) setOpenMobile(false)
    navigate({ to: href })
  }

  const filteredSidebarContent = sidebarContentList.filter((group) => {
    if (group.admin && !isAdmin) return false
    return true
  })

  return (
    <Sidebar collapsible="icon" variant="sidebar" className="border-r bg-sidebar text-sidebar-foreground">
      <SidebarHeader className="h-14 flex flex-row items-center justify-center border-b border-sidebar-border group-data-[collapsible=icon]:px-0">
        <div className="flex items-center gap-2.5 overflow-hidden w-full px-5 group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0">
          <div className="flex size-6 shrink-0 items-center justify-center rounded-md bg-sidebar-primary text-sidebar-primary-foreground text-[10px] font-bold">
            {companyMeta.name.charAt(0)}
          </div>
          <span className="text-sm font-semibold tracking-tight text-sidebar-foreground truncate group-data-[collapsible=icon]:hidden">
            {companyMeta.name}
          </span>
        </div>
      </SidebarHeader>

      <SidebarContent className="py-3 overflow-y-auto [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none]">
        {filteredSidebarContent.map((group) => (
          <SidebarGroup key={group.groupName || "default"} className="px-2 py-0.5">
            {group.groupName && (
              <SidebarGroupLabel className="px-3 py-1 text-[11px] font-medium tracking-wide text-sidebar-foreground/50 uppercase">
                {group.groupName}
              </SidebarGroupLabel>
            )}
            <SidebarGroupContent>
              <SidebarMenu>
                {group.items.map((item) => {
                  const Icon = item.icon
                  const isActive = currentPath === item.href

                  return (
                    <SidebarMenuItem key={item.title}>
                      <SidebarMenuButton
                        isActive={isActive}
                        onClick={() => handleItemClick(item.href)}
                        className={cn(
                          "relative flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-all duration-150 cursor-pointer group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0",
                          isActive
                            ? "bg-sidebar-accent text-sidebar-accent-foreground"
                            : "text-sidebar-foreground/70 hover:bg-sidebar-accent/50 hover:text-sidebar-foreground"
                        )}
                      >
                        <Icon className={cn(
                          "size-4 shrink-0 transition-colors",
                          isActive ? "text-sidebar-primary" : "text-sidebar-foreground/50"
                        )} />
                        <span className="group-data-[collapsible=icon]:hidden">
                          {item.title}
                        </span>
                      </SidebarMenuButton>
                      {item.badge && (
                        <SidebarMenuBadge className="group-data-[collapsible=icon]:hidden ml-auto bg-sidebar-primary/10 text-sidebar-primary text-[10px] font-semibold px-1.5 py-0.5 rounded-md">
                          {item.badge}
                        </SidebarMenuBadge>
                      )}
                    </SidebarMenuItem>
                  )
                })}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>

      <SidebarFooter className="border-t border-sidebar-border px-3 py-3">
        {user && (
          <div className="mb-2 flex items-center gap-3 rounded-lg px-3 py-2 group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0">
            <div className="relative shrink-0">
              <img
                src={user.avatar || `https://ui-avatars.com/api/?name=${encodeURIComponent(user.name)}&background=6366f1&color=fff`}
                alt={user.name}
                className="size-7 rounded-full object-cover"
              />
            </div>
            <div className="flex-1 overflow-hidden group-data-[collapsible=icon]:hidden">
              <p className="truncate text-sm font-medium text-sidebar-foreground leading-tight">
                {user.name}
              </p>
              <p className="truncate text-[11px] text-sidebar-foreground/50 leading-tight">
                {user.email}
              </p>
            </div>
          </div>
        )}

        <SidebarSeparator className="bg-sidebar-border/50 mb-2 group-data-[collapsible=icon]:hidden" />

        <div className="flex flex-col gap-0.5">
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton
                onClick={toggleTheme}
                tooltip={theme === "light" ? "Dark Mode" : "Light Mode"}
                className="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-sidebar-foreground/70 hover:bg-sidebar-accent/50 hover:text-sidebar-foreground transition-all duration-150 cursor-pointer group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0"
              >
                {theme === "light" ? (
                  <Moon className="size-4 shrink-0 text-sidebar-foreground/50" />
                ) : (
                  <Sun className="size-4 shrink-0 text-sidebar-foreground/50" />
                )}
                <span className="group-data-[collapsible=icon]:hidden">
                  {theme === "light" ? "Dark Mode" : "Light Mode"}
                </span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>

          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton
                onClick={handleLogout}
                tooltip="Sign Out"
                className="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-sidebar-foreground/70 hover:bg-destructive/10 hover:text-destructive transition-all duration-150 cursor-pointer group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0"
              >
                <LogOut className="size-4 shrink-0" />
                <span className="group-data-[collapsible=icon]:hidden">Sign Out</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </div>
      </SidebarFooter>
    </Sidebar>
  )
}
