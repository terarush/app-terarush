import {
	LayoutDashboard,
	Users,
	Settings,
	Server,
	Package,
	Receipt,
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
		name: "My Transactions",
		link: "/dashboard/transactions",
		icon: Receipt,
	},
	{
		name: "Nodes",
		link: "/dashboard/nodes",
		icon: Server,
		adminOnly: true,
	},
	{
		name: "Products",
		link: "/dashboard/admin/products",
		icon: Package,
		adminOnly: true,
	},
	{
		name: "Users",
		link: "/dashboard/users",
		icon: Users,
		adminOnly: true,
	},
	{
		name: "Settings",
		link: "/dashboard/settings",
		icon: Settings,
	},
];
