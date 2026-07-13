import { useState, useEffect } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import {
	Bell,
	Search,
	ChevronLeft,
	ChevronRight,
	LogOut,
	Crown,
	Moon,
	Sun,
} from "lucide-react";
import { cn, getAvatarUrl } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from "@/components/ui/tooltip";
import { Badge } from "@/components/ui/badge";
import { useUser, useIsAdmin, useLogout } from "@/hooks";
import { useTheme } from "@/components/theme-provider";
import { siteConfig } from "@/content/config";
import { useIsMobile } from "@/hooks/use-mobile";
import { sidebarNavItems } from "@/content/sidebar";

export function AppSidebar() {
	const [collapsed, setCollapsed] = useState(false);
	const location = useLocation();
	const navigate = useNavigate();
	const user = useUser();
	const isAdmin = useIsAdmin();
	const logout = useLogout();
	const { theme, setTheme } = useTheme();
	const isMobile = useIsMobile();

	// Auto-collapse sidebar on mobile
	useEffect(() => {
		if (isMobile && !collapsed) {
			setCollapsed(true);
		}
	}, [isMobile, collapsed]);

	const toggleTheme = () => {
		setTheme(theme === "dark" ? "light" : "dark");
	};

	const handleLogout = async () => {
		await logout();
	};

	const getInitials = (name: string) => {
		return name
			.split(" ")
			.map((n) => n[0])
			.join("")
			.toUpperCase()
			.slice(0, 2);
	};

	const filteredNavItems = sidebarNavItems.filter(
		(item) => !item.adminOnly || isAdmin,
	);

	return (
		<TooltipProvider delayDuration={0}>
			<aside
				className={cn(
					"relative flex flex-col h-screen bg-card border-r border-border transition-all duration-300 ease-in-out",
					collapsed ? "w-16" : "w-64",
				)}
			>
				{/* Header */}
				<div className="flex items-center justify-between h-16 px-4 border-b border-border">
					{!collapsed && (
						<Link
							to="/"
							className="flex items-center space-x-2 group"
						>
							<img
								src={siteConfig.logo}
								alt={siteConfig.name}
								className="h-8 w-8 rounded-lg group-hover:scale-110 transition-transform duration-200"
							/>
							<span className="text-xl font-bold text-foreground">
								{siteConfig.name}
							</span>
						</Link>
					)}
					{collapsed && (
						<Link
							to="/"
							className="flex items-center justify-center w-full"
						>
							<img
								src="/assets/logo.png"
								alt={siteConfig.name}
								className="h-8 w-8 rounded-lg hover:scale-110 transition-transform duration-200"
							/>
						</Link>
					)}
				</div>

				{/* Search Bar */}
				{!collapsed && (
					<div className="p-4">
						<Button
							variant="outline"
							className="w-full justify-start text-muted-foreground hover:text-foreground"
							onClick={() => navigate("/dashboard/search")}
						>
							<Search className="h-4 w-4 mr-2" />
							Search...
						</Button>
					</div>
				)}

			{/* Navigation */}
			<nav className="flex-1 overflow-y-auto px-3 py-4 space-y-1">
				{filteredNavItems.map((item) => {
					const Icon = item.icon;
					// For Dashboard, use exact match. For others, check if path starts with the link
					const isActive = item.link === "/dashboard" 
						? location.pathname === item.link 
						: location.pathname.startsWith(item.link);

						const navButton = (
							<Button
								key={item.link}
								variant={isActive ? "secondary" : "ghost"}
								className={cn(
									"w-full justify-start group relative transition-all duration-200",
									collapsed ? "px-0 justify-center" : "px-3",
									isActive &&
										"bg-primary/10 text-primary hover:bg-primary/15",
								)}
								onClick={() => navigate(item.link)}
							>
								{isActive && (
									<div className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 bg-primary rounded-r-full" />
								)}
								<Icon
									className={cn(
										"h-5 w-5 transition-transform duration-200 group-hover:scale-110",
										collapsed ? "" : "mr-3",
										isActive && "text-primary",
									)}
								/>
								{!collapsed && (
									<>
										<span className="flex-1 text-left">
											{item.name}
										</span>
										{item.badge && (
											<Badge
												variant={
													item.badge === "New"
														? "default"
														: "secondary"
												}
												className="ml-auto text-xs"
											>
												{item.badge}
											</Badge>
										)}
									</>
								)}
							</Button>
						);

						if (collapsed) {
							return (
								<Tooltip key={item.link}>
									<TooltipTrigger asChild>
										{navButton}
									</TooltipTrigger>
									<TooltipContent
										side="right"
										className="flex items-center gap-2"
									>
										{item.name}
										{item.badge && (
											<Badge
												variant={
													item.badge === "New"
														? "default"
														: "secondary"
												}
												className="text-xs"
											>
												{item.badge}
											</Badge>
										)}
									</TooltipContent>
								</Tooltip>
							);
						}

						return navButton;
					})}
				</nav>

				<Separator />

				{/* User Profile & Controls */}
				<div className="p-3 space-y-2">
					{/* Notifications & Theme Toggle */}
					<div className={cn(
						"flex gap-2",
						collapsed && "flex-col"
					)}>
						{!collapsed ? (
							<>
								<Button
									variant="ghost"
									className="flex-1 justify-start"
									onClick={() => navigate("/dashboard/notifications")}
								>
									<Bell className="h-5 w-5 mr-3" />
									<span className="flex-1 text-left">
										Notifications
									</span>
									<Badge variant="destructive" className="text-xs">
										5
									</Badge>
								</Button>
								<Button
									variant="ghost"
									size="sm"
									onClick={toggleTheme}
								>
									{theme === "light" ? (
										<Moon className="h-4 w-4" />
									) : (
										<Sun className="h-4 w-4" />
									)}
								</Button>
							</>
						) : (
							<>
								<Tooltip>
									<TooltipTrigger asChild>
										<Button
											variant="ghost"
											className="w-full justify-center relative"
											onClick={() =>
												navigate("/dashboard/notifications")
											}
										>
											<Bell className="h-5 w-5" />
											<span className="absolute -top-1 -right-1 h-4 w-4 rounded-full bg-destructive text-destructive-foreground text-xs flex items-center justify-center">
												5
											</span>
										</Button>
									</TooltipTrigger>
									<TooltipContent side="right">
										Notifications (5)
									</TooltipContent>
								</Tooltip>
								<Tooltip>
									<TooltipTrigger asChild>
										<Button
											variant="ghost"
											size="sm"
											className="w-full"
											onClick={toggleTheme}
										>
											{theme === "light" ? (
												<Moon className="h-4 w-4" />
											) : (
												<Sun className="h-4 w-4" />
											)}
										</Button>
									</TooltipTrigger>
									<TooltipContent side="right">
										{theme === "light" ? "Dark Mode" : "Light Mode"}
									</TooltipContent>
								</Tooltip>
							</>
						)}
					</div>

					{/* User Profile */}
					<div
						className={cn(
							"flex items-center gap-3 p-2 rounded-lg hover:bg-muted/50 transition-colors cursor-pointer",
							collapsed && "justify-center",
						)}
						onClick={() => navigate("/dashboard/profile")}
					>
					<Avatar className="h-9 w-9 border-2 border-primary/20">
						{user?.avatar && (
							<AvatarImage src={getAvatarUrl(user.avatar)} alt={user.name} />
						)}
						<AvatarFallback className="bg-primary/10 text-primary font-semibold">
							{user?.name ? getInitials(user.name) : "U"}
						</AvatarFallback>
					</Avatar>
						{!collapsed && (
							<div className="flex-1 min-w-0">
								<div className="flex items-center gap-2">
									<p className="text-sm font-medium truncate">
										{user?.name || "User"}
									</p>
									{isAdmin && (
										<Crown className="h-3.5 w-3.5 text-yellow-500" />
									)}
								</div>
								<p className="text-xs text-muted-foreground truncate">
									{user?.email}
								</p>
							</div>
						)}
					</div>

					<Separator />

					{/* Logout & Collapse */}
					<div className={cn(
						"flex gap-2",
						collapsed && "flex-col"
					)}>
						{!collapsed ? (
							<>
								<Button
									variant="ghost"
									size="sm"
									className="flex-1 text-destructive hover:text-destructive hover:bg-destructive/10"
									onClick={handleLogout}
								>
									<LogOut className="h-4 w-4 mr-2" />
									Logout
								</Button>
								<Button
									variant="ghost"
									size="sm"
									onClick={() => setCollapsed(true)}
								>
									<ChevronLeft className="h-4 w-4" />
								</Button>
							</>
						) : (
							<>
								<Tooltip>
									<TooltipTrigger asChild>
										<Button
											variant="ghost"
											size="sm"
											className="w-full text-destructive hover:text-destructive hover:bg-destructive/10"
											onClick={handleLogout}
										>
											<LogOut className="h-4 w-4" />
										</Button>
									</TooltipTrigger>
									<TooltipContent side="right">
										Logout
									</TooltipContent>
								</Tooltip>
								<Tooltip>
									<TooltipTrigger asChild>
										<Button
											variant="ghost"
											size="sm"
											className="w-full"
											onClick={() => setCollapsed(false)}
										>
											<ChevronRight className="h-4 w-4" />
										</Button>
									</TooltipTrigger>
									<TooltipContent side="right">
										Expand
									</TooltipContent>
								</Tooltip>
							</>
						)}
					</div>
				</div>
			</aside>
		</TooltipProvider>
	);
}
