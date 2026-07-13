import { useState, useEffect, useRef } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import {
	Menu,
	X,
	User,
	LogOut,
	LayoutDashboard,
	Crown,
} from "lucide-react";
import { Button } from "@/components/ui/button";
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
import { useReducedMotion } from "motion/react";
import { ThemeToggle } from "./elements/toggle-theme";
import { siteConfig, navConfig } from "@/content/config";
import { useAuth } from "@/contexts/AuthContext";
import { getAvatarUrl } from "@/lib/utils";

export default function Navbar() {
	const reduce = useReducedMotion();
	const [isScrolled, setIsScrolled] = useState(false);
	const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
	const location = useLocation();
	const navigate = useNavigate();
	const { user, isAuthenticated, logout } = useAuth();
	const headerRef = useRef<HTMLElement>(null);
	const mobileMenuRef = useRef<HTMLDivElement>(null);
	const mobileBackdropRef = useRef<HTMLDivElement>(null);

	const getInitials = (name: string) =>
		name
			.split(" ")
			.map((n) => n[0])
			.join("")
			.toUpperCase()
			.slice(0, 2);

	const handleNavClick = (href: string) => {
		if (href.startsWith("#") && location.pathname === "/") {
			document.querySelector(href)?.scrollIntoView({ behavior: "smooth" });
		} else if (href.startsWith("#") && location.pathname !== "/") {
			navigate("/" + href);
		}
	};

	const handleLogout = async () => {
		await logout();
		navigate("/");
	};

	useEffect(() => {
		setIsMobileMenuOpen(false);
	}, [location]);

	// ponytail: IntersectionObserver over scroll listener — fires once on threshold cross, not every frame
	useEffect(() => {
		const sentinel = document.createElement("div");
		sentinel.style.position = "absolute";
		sentinel.style.top = "1px";
		sentinel.style.height = "1px";
		sentinel.style.width = "1px";
		sentinel.style.pointerEvents = "none";
		document.body.prepend(sentinel);

		const observer = new IntersectionObserver(
			([entry]) => setIsScrolled(!entry.isIntersecting),
			{ threshold: 0 },
		);
		observer.observe(sentinel);

		return () => {
			observer.disconnect();
			sentinel.remove();
		};
	}, []);

	// GSAP header mount animation
	useEffect(() => {
		if (reduce || !headerRef.current) return;
		gsap.fromTo(
			headerRef.current,
			{ y: -100, opacity: 0 },
			{ y: 0, opacity: 1, duration: 0.8, ease: "power3.out" },
		);
	}, []);

	// GSAP mobile menu animations
	useEffect(() => {
		if (reduce) return;
		if (isMobileMenuOpen) {
			if (mobileBackdropRef.current) {
				gsap.fromTo(
					mobileBackdropRef.current,
					{ opacity: 0 },
					{ opacity: 1, duration: 0.3, ease: "power2.out" },
				);
			}
			if (mobileMenuRef.current) {
				gsap.fromTo(
					mobileMenuRef.current,
					{ y: 100, opacity: 0 },
					{ y: 0, opacity: 1, duration: 0.5, ease: "power3.out" },
				);
				const items =
					mobileMenuRef.current.querySelectorAll(".menu-item");
				gsap.fromTo(
					items,
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
						"relative backdrop-blur-xl border border-border/60 shadow-lg shadow-black/5 rounded-2xl mx-auto max-w-7xl transition-all duration-500",
						isScrolled
							? "bg-background/90 shadow-md"
							: "bg-background/70",
					].join(" ")}
				>
					<div className="relative px-5 flex items-center justify-between h-14">
						<Link
							to="/"
							className="flex items-center gap-3 group shrink-0"
						>
							<img
								src={siteConfig.logo}
								alt={siteConfig.name}
								className="h-7 w-7 object-contain transition-transform duration-300 group-hover:scale-110 rounded-md"
							/>
							<span className="text-lg font-semibold tracking-tight">
								{siteConfig.name}
							</span>
						</Link>

						<div className="hidden lg:flex items-center gap-1">
							{navConfig.mainNav.map((item) => {
								if (item.href.startsWith("/")) {
									return (
										<Link
											key={item.href}
											to={item.href}
											className="px-3 py-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors duration-200 font-medium"
										>
											{item.label}
										</Link>
									);
								}
								return (
									<a
										key={item.href}
										href={item.href}
										onClick={() => handleNavClick(item.href)}
										className="px-3 py-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors duration-200 font-medium"
									>
										{item.label}
									</a>
								);
							})}
						</div>

						<div className="flex items-center gap-2">
							<ThemeToggle variant="rounded" />

							{isAuthenticated && user ? (
								<DropdownMenu>
									<DropdownMenuTrigger asChild>
										<Button
											variant="ghost"
											className="hidden sm:flex items-center gap-2 px-2 py-1.5 h-9 rounded-xl hover:bg-accent/50 transition-all duration-200"
										>
											<Avatar className="h-7 w-7 border-2 border-primary/20">
												{user.avatar && (
													<AvatarImage
														src={getAvatarUrl(user.avatar)}
														alt={user.name}
													/>
												)}
												<AvatarFallback className="bg-primary/10 text-primary text-[10px] font-semibold">
													{getInitials(user.name)}
												</AvatarFallback>
											</Avatar>
											<div className="flex items-center gap-1 leading-tight">
												<span className="text-sm font-medium">
													{user.name}
												</span>
												{user.role === "admin" && (
													<Crown className="h-3 w-3 text-yellow-500" />
												)}
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
									className="hidden sm:inline-flex h-8 px-4 text-sm font-semibold rounded-full bg-neutral-900 dark:bg-neutral-100 text-white dark:text-neutral-900 shadow-sm hover:shadow-md transition-all duration-300 hover:scale-[1.02]"
									asChild
								>
									<Link to="/login">Get Started</Link>
								</Button>
							)}

							<Button
								variant="ghost"
								size="icon"
								className="lg:hidden rounded-full"
								onClick={() =>
									setIsMobileMenuOpen(!isMobileMenuOpen)
								}
							>
								{isMobileMenuOpen ? (
									<X className="h-5 w-5" />
								) : (
									<Menu className="h-5 w-5" />
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
							className="relative bg-background/95 backdrop-blur-xl border border-border/60 shadow-2xl rounded-2xl p-5 pointer-events-auto translate-y-full opacity-0"
						>
							<div className="space-y-1">
								{isAuthenticated && user && (
									<div className="menu-item flex items-center gap-3 p-3 rounded-xl bg-accent/50 mb-3">
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

								{navConfig.mainNav.map((item) => {
									if (item.href.startsWith("/")) {
										return (
											<Link
												key={item.href}
												to={item.href}
												className="menu-item flex items-center p-3 rounded-xl hover:bg-accent/50 transition-all duration-200"
												onClick={() =>
													setIsMobileMenuOpen(false)
												}
											>
												<span className="font-medium">
													{item.label}
												</span>
											</Link>
										);
									}
									return (
										<a
											key={item.href}
											href={item.href}
											className="menu-item flex items-center p-3 rounded-xl hover:bg-accent/50 transition-all duration-200"
											onClick={(e) => {
												e.preventDefault();
												setIsMobileMenuOpen(false);
												handleNavClick(item.href);
											}}
										>
											<span className="font-medium">
												{item.label}
											</span>
										</a>
									);
								})}

								{isAuthenticated && (
									<>
										<div className="my-2 border-t border-border/60" />
										<Link
											to="/dashboard"
											className="menu-item flex items-center gap-3 p-3 rounded-xl hover:bg-accent/50 transition-all duration-200"
											onClick={() =>
												setIsMobileMenuOpen(false)
											}
										>
											<LayoutDashboard className="h-5 w-5" />
											<span className="font-medium">
												Dashboard
											</span>
										</Link>
										<Link
											to="/dashboard/profile"
											className="menu-item flex items-center gap-3 p-3 rounded-xl hover:bg-accent/50 transition-all duration-200"
											onClick={() =>
												setIsMobileMenuOpen(false)
											}
										>
											<User className="h-5 w-5" />
											<span className="font-medium">
												Profile
											</span>
										</Link>
									</>
								)}
							</div>

							<div className="mt-4 pt-4 border-t border-border/60 menu-item">
								{isAuthenticated ? (
									<Button
										className="w-full rounded-xl py-2.5"
										variant="destructive"
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
										className="w-full rounded-xl py-2.5 bg-neutral-900 dark:bg-neutral-100 text-white dark:text-neutral-900 font-semibold"
										asChild
									>
										<Link
											to="/login"
											onClick={() =>
												setIsMobileMenuOpen(false)
											}
										>
											Get Started
										</Link>
									</Button>
								)}
							</div>
						</div>
					</div>
				</div>
			)}
		</>
	);
}
