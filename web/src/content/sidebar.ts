import {
	LayoutDashboard,
	Settings,
	BookOpen,
	ImageIcon,
	type LucideIcon,
} from "lucide-react";

export interface SidebarNavItem {
	name: string;
	icon: LucideIcon;
	badge?: string;
	link: string;
	adminOnly?: boolean;
}

export const sidebarNavItems: SidebarNavItem[] = [
	{
		name: "Dashboard",
		link: "/dashboard",
		icon: LayoutDashboard,
	},
	{
		name: "Blog Management",
		link: "/dashboard/blogs",
		icon: BookOpen,
		adminOnly: true,
	},
	{
		name: "Assets",
		link: "/dashboard/assets",
		icon: ImageIcon,
		adminOnly: true,
	},
	{
		name: "Settings",
		link: "/dashboard/settings",
		icon: Settings,
	},
];
