import { useState, useEffect, useRef } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import {
	Menu,
	Search,
	X,
	User,
	LogOut,
	LayoutDashboard,
	Crown,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import gsap from "gsap";
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { ThemeToggle } from "./elements/toggle-theme";
import { siteConfig, navConfig } from "@/content/config";
import { useAuth } from "@/contexts/AuthContext";
import { getAvatarUrl } from "@/lib/utils";

export default function Navbar() {
	const [isScrolled, setIsScrolled] = useState(false);
	const [isSearchOpen, setIsSearchOpen] = useState(false);
	const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
	const location = useLocation();
	const navigate = useNavigate();
	const { user, isAuthenticated, logout } = useAuth();
	const headerRef = useRef<HTMLElement>(null);
	const mobileMenuRef = useRef<HTMLDivElement>(null);
	const mobileBackdropRef = useRef<HTMLDivElement>(null);

	const getInitials = (name: string) => {
		return name
			.split(" ")
			.map((n) => n[0])
			.join("")
			.toUpperCase()
			.slice(0, 2);
	};

	const handleNavClick = (href: string) => {
		// If it's a hash link and we're on home page, scroll to section
		if (href.startsWith("#") && location.pathname === "/") {
			const target = document.querySelector(href);
			if (target) {
				target.scrollIntoView({ behavior: "smooth" });
			}
		}
	};

	const handleLogout = async () => {
		await logout();
		navigate("/");
	};

	useEffect(() => {
		// eslint-disable-next-line react-hooks/exhaustive-deps
		setIsMobileMenuOpen(false);
	}, [location]);

	useEffect(() => {
		const handleScroll = () => {
			setIsScrolled(window.scrollY > 20);
		};
		window.addEventListener("scroll", handleScroll);
		return () => window.removeEventListener("scroll", handleScroll);
	}, []);

	// GSAP animation for header on mount
	useEffect(() => {
		if (headerRef.current) {
			gsap.fromTo(
				headerRef.current,
				{ y: -100, opacity: 0 },
				{ y: 0, opacity: 1, duration: 0.8, ease: "power3.out" },
			);
		}
	}, []);

	// GSAP animation for mobile menu
	useEffect(() => {
		if (isMobileMenuOpen) {
			// Animate backdrop
			if (mobileBackdropRef.current) {
				gsap.fromTo(
					mobileBackdropRef.current,
					{ opacity: 0 },
					{ opacity: 1, duration: 0.3, ease: "power2.out" },
				);
			}

			// Animate menu
			if (mobileMenuRef.current) {
				gsap.fromTo(
					mobileMenuRef.current,
					{ y: 100, opacity: 0 },
					{ y: 0, opacity: 1, duration: 0.5, ease: "power3.out" },
				);

				// Animate menu items
				const menuItems =
					mobileMenuRef.current.querySelectorAll(".menu-item");
				gsap.fromTo(
					menuItems,
					{ y: 20, opacity: 0 },
					{
						y: 0,
						opacity: 1,
						duration: 0.4,
						stagger: 0.05,
						ease: "power2.out",
						delay: 0.2,
					},
				);
			}
		} else if (mobileMenuRef.current && mobileBackdropRef.current) {
			// Animate out
			gsap.to(mobileMenuRef.current, {
				y: 100,
				opacity: 0,
				duration: 0.3,
				ease: "power2.in",
			});
			gsap.to(mobileBackdropRef.current, {
				opacity: 0,
				duration: 0.2,
				ease: "power2.in",
			});
		}
	}, [isMobileMenuOpen]);

	return (
		<>
			<header
				ref={headerRef}
				className="fixed top-4 left-0 right-0 z-50 px-4"
				style={{ opacity: 0 }}
			>
				<nav
					className={[
						"relative backdrop-blur-xl border border-border shadow-md shadow-gray-900/10 transition-all duration-500 rounded-xl mx-auto",
						isScrolled ? "py-3 max-w-7xl" : "py-3 max-w-7xl",
						"bg-background/80",
					].join(" ")}
				>
					<div className="relative px-6 flex items-center justify-between">
						<div className="flex items-center space-x-8">
							<Link
								to="/"
								className="flex items-center space-x-3 group"
							>
								<div className="relative">
									<img
										src="/assets/logo.png"
										alt={siteConfig.name}
										className="h-8 w-8 object-contain transition-transform duration-300 group-hover:scale-110 rounded-md"
									/>
								</div>
								<span className="text-xl font-bold text-gray-900 dark:text-white">
									{siteConfig.name}
								</span>
							</Link>
							<div className="hidden lg:flex items-center space-x-1">
								{navConfig.mainNav.map((item) => {
									// Use Link component for internal routes
									if (item.href.startsWith("/")) {
										return (
											<Link
												key={item.href}
												to={item.href}
												className="px-4 py-2 rounded-xl text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200 font-medium"
											>
												{item.label}
											</Link>
										);
									}
									// Use anchor tag for hash links
									return (
										<a
											key={item.href}
											href={item.href}
											onClick={() => handleNavClick(item.href)}
											className="px-4 py-2 rounded-xl text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200 font-medium"
										>
											{item.label}
										</a>
									);
								})}
							</div>
						</div>
						<div className="flex items-center space-x-3">
							<Dialog
								open={isSearchOpen}
								onOpenChange={setIsSearchOpen}
							>
								<DialogTrigger asChild>
									<Button
										variant="ghost"
										size="sm"
										className="hidden lg:flex items-center space-x-2 px-3 py-2 rounded-xl bg-gray-100/50 dark:bg-gray-800/50 hover:bg-gray-200/50 dark:hover:bg-gray-700/50 transition-all duration-200 min-w-[200px] justify-start"
									>
										<Search className="h-4 w-4 text-gray-500" />
										<span className="text-gray-500 text-sm">
											Search...
										</span>
									</Button>
								</DialogTrigger>
								<DialogContent className="max-w-2xl">
									<DialogHeader>
										<DialogTitle>Search</DialogTitle>
									</DialogHeader>
									<div className="flex items-center space-x-2 border-b pb-4">
										<Search className="h-5 w-5 text-gray-400" />
										<Input
											placeholder="Search products, docs, and more..."
											className="border-0 focus-visible:ring-0 text-lg"
											autoFocus
										/>
									</div>
									<div className="py-4">
										<p className="text-sm text-gray-500">
											Start typing to search...
										</p>
									</div>
								</DialogContent>
							</Dialog>
							<Button
								variant="ghost"
								size="sm"
								className="lg:hidden p-2 rounded-xl hover:bg-gray-100/50 dark:hover:bg-gray-800/50"
							>
								<Search className="h-5 w-5" />
							</Button>
							<ThemeToggle variant="rounded" />

							{/* User Menu - Show when authenticated */}
							{isAuthenticated && user ? (
								<DropdownMenu>
									<DropdownMenuTrigger asChild>
										<Button
											variant="ghost"
											className="hidden sm:flex items-center space-x-2 px-3 py-2 rounded-xl hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200"
										>
										<Avatar className="h-8 w-8 border-2 border-primary/20">
											{user.avatar && (
												<AvatarImage
													src={getAvatarUrl(user.avatar)}
													alt={user.name}
												/>
											)}
											<AvatarFallback className="bg-primary/10 text-primary text-xs font-semibold">
												{getInitials(user.name)}
											</AvatarFallback>
										</Avatar>
											<div className="flex flex-col items-start">
												<div className="flex items-center gap-1.5">
													<span className="text-sm font-medium">
														{user.name}
													</span>
													{user.role === "admin" && (
														<Crown className="h-3 w-3 text-yellow-500" />
													)}
												</div>
											</div>
										</Button>
									</DropdownMenuTrigger>
									<DropdownMenuContent
										align="end"
										className="w-56"
									>
										<DropdownMenuLabel>
											<div className="flex flex-col space-y-1">
												<div className="flex items-center gap-1.5">
													<p className="text-sm font-medium leading-none">
														{user.name}
													</p>
													{user.role === "admin" && (
														<Crown className="h-3.5 w-3.5 text-yellow-500" />
													)}
												</div>
												<p className="text-xs leading-none text-muted-foreground">
													{user.email}
												</p>
											</div>
										</DropdownMenuLabel>
										<DropdownMenuSeparator />
										<DropdownMenuItem
											onClick={() =>
												navigate("/dashboard")
											}
										>
											<LayoutDashboard className="mr-2 h-4 w-4" />
											Dashboard
										</DropdownMenuItem>
										<DropdownMenuItem
											onClick={() =>
												navigate("/dashboard/profile")
											}
										>
											<User className="mr-2 h-4 w-4" />
											Profile
										</DropdownMenuItem>
										<DropdownMenuSeparator />
										<DropdownMenuItem
											onClick={handleLogout}
											className="text-destructive focus:text-destructive"
										>
											<LogOut className="mr-2 h-4 w-4" />
											Logout
										</DropdownMenuItem>
									</DropdownMenuContent>
								</DropdownMenu>
							) : (
								<Button
									className="hidden sm:flex bg-gray-900 dark:bg-gray-100 text-white dark:text-gray-900 rounded-xl px-6 py-2 font-semibold shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-105"
									asChild
								>
									<Link to="/login">Get Started</Link>
								</Button>
							)}

							<Button
								variant="ghost"
								className="lg:hidden p-2 rounded-md hover:bg-gray-100/50 dark:hover:bg-gray-800/50"
								onClick={() =>
									setIsMobileMenuOpen(!isMobileMenuOpen)
								}
							>
								{isMobileMenuOpen ? (
									<X className="h-6 w-6" />
								) : (
									<Menu className="h-6 w-6" />
								)}
							</Button>
						</div>
					</div>
				</nav>
			</header>
			{isMobileMenuOpen && (
				<div className="fixed inset-0 z-40 lg:hidden">
					<div
						ref={mobileBackdropRef}
						className="fixed inset-0 bg-black/20 backdrop-blur-sm opacity-0"
						onClick={() => setIsMobileMenuOpen(false)}
					/>
					<div className="fixed bottom-0 left-0 right-0 p-4 pointer-events-none">
						<div
							ref={mobileMenuRef}
							className="relative bg-background/95 backdrop-blur-xl border border-border shadow-2xl rounded-2xl p-6 pointer-events-auto translate-y-full opacity-0"
						>
							<div className="space-y-4">
								{/* User Info in Mobile Menu */}
								{isAuthenticated && user && (
									<div className="menu-item flex items-center space-x-3 p-3 rounded-xl bg-gray-100/50 dark:bg-gray-800/50">
									<Avatar className="h-10 w-10 border-2 border-primary/20">
										{user.avatar && (
											<AvatarImage
												src={getAvatarUrl(user.avatar)}
												alt={user.name}
											/>
										)}
										<AvatarFallback className="bg-primary/10 text-primary font-semibold">
											{getInitials(user.name)}
										</AvatarFallback>
									</Avatar>
										<div className="flex-1">
											<div className="flex items-center gap-1.5">
												<p className="text-sm font-medium">
													{user.name}
												</p>
												{user.role === "admin" && (
													<Crown className="h-3 w-3 text-yellow-500" />
												)}
											</div>
											<p className="text-xs text-muted-foreground">
												{user.email}
											</p>
										</div>
									</div>
								)}

								<div className="space-y-2">
									{navConfig.mainNav.map((item) => {
										// Use Link component for internal routes
										if (item.href.startsWith("/")) {
											return (
												<Link
													key={item.href}
													to={item.href}
													className="menu-item flex items-center space-x-3 p-3 rounded-xl hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200"
													onClick={() =>
														setIsMobileMenuOpen(false)
													}
												>
													<span className="font-medium text-gray-900 dark:text-white">
														{item.label}
													</span>
												</Link>
											);
										}
										// Use anchor tag for hash links
										return (
											<a
												key={item.href}
												href={item.href}
												className="menu-item flex items-center space-x-3 p-3 rounded-xl hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200"
												onClick={(e) => {
													e.preventDefault();
													setIsMobileMenuOpen(false);
													handleNavClick(item.href);
												}}
											>
												<span className="font-medium text-gray-900 dark:text-white">
													{item.label}
												</span>
											</a>
										);
									})}

									{/* Additional menu items when authenticated */}
									{isAuthenticated && (
										<>
											<Link
												to="/dashboard"
												className="menu-item flex items-center space-x-3 p-3 rounded-xl hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200"
												onClick={() =>
													setIsMobileMenuOpen(false)
												}
											>
												<LayoutDashboard className="h-5 w-5" />
												<span className="font-medium text-gray-900 dark:text-white">
													Dashboard
												</span>
											</Link>
											<Link
												to="/dashboard/profile"
												className="menu-item flex items-center space-x-3 p-3 rounded-xl hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200"
												onClick={() =>
													setIsMobileMenuOpen(false)
												}
											>
												<User className="h-5 w-5" />
												<span className="font-medium text-gray-900 dark:text-white">
													Profile
												</span>
											</Link>
										</>
									)}
								</div>

								<div className="pt-2 menu-item">
									{isAuthenticated ? (
										<Button
											className="w-full bg-destructive text-destructive-foreground rounded-xl py-3 font-semibold"
											onClick={() => {
												handleLogout();
												setIsMobileMenuOpen(false);
											}}
										>
											<LogOut className="mr-2 h-5 w-5" />
											Logout
										</Button>
									) : (
										<Button
											className="w-full bg-gray-900 dark:bg-gray-100 text-white dark:text-gray-900 rounded-xl py-3 font-semibold"
											asChild
										>
											<Link to="/login">Get Started</Link>
										</Button>
									)}
								</div>
							</div>
						</div>
					</div>
				</div>
			)}
		</>
	);
}
