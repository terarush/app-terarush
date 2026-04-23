import {
	LayoutDashboard,
	Settings,
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
		name: "Settings",
		link: "/dashboard/settings",
		icon: Settings,
	},
];
