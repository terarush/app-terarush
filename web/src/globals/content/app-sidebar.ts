import {
  LayoutDashboard,
  BookOpen,
  MessageSquare,
  Heart,
  Users,
  type LucideIcon,
} from "lucide-react"
import type { ComponentType } from "react"

export interface SidebarItem {
  title: string
  href: string
  icon: LucideIcon | ComponentType<{ className?: string }>
  badge?: string
}

export interface SidebarGroupType {
  groupName?: string
  admin?: boolean
  items: SidebarItem[]
}

export const sidebarContentList: SidebarGroupType[] = [
  {
    groupName: "Workspace",
    items: [
      { title: "Dashboard", href: "/app", icon: LayoutDashboard },
      { title: "Blogs", href: "/app/blogs", icon: BookOpen },
    ],
  },
  {
    groupName: "Interactions",
    items: [
      { title: "Comments", href: "/app/comments", icon: MessageSquare, badge: "New" },
      { title: "Favorites", href: "/app/favorites", icon: Heart },
    ],
  },
  {
    groupName: "Administration",
    admin: true,
    items: [
      { title: "Users", href: "/app/users", icon: Users },
    ],
  },
]
