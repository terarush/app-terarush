import { useState, useEffect, useRef } from "react";
import { Link, useLocation } from "react-router-dom";
import { Menu, Search, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
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

export default function Navbar() {
	const [isScrolled, setIsScrolled] = useState(false);
	const [isSearchOpen, setIsSearchOpen] = useState(false);
	const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
	const location = useLocation();
	const headerRef = useRef<HTMLElement>(null);
	const mobileMenuRef = useRef<HTMLDivElement>(null);
	const mobileBackdropRef = useRef<HTMLDivElement>(null);

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
								{navConfig.mainNav.map((item) => (
									<a
										key={item.href}
										href={item.href}
										className="px-4 py-2 rounded-xl text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200 font-medium"
									>
										{item.label}
									</a>
								))}
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
							<Button
								className="hidden sm:flex bg-gray-900 dark:bg-gray-100 text-white dark:text-gray-900 rounded-xl px-6 py-2 font-semibold shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-105"
								asChild
							>
								<Link to="/auth/login">Get Started</Link>
							</Button>
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
								<div className="space-y-2">
									{navConfig.mainNav.map((item) => (
										<a
											key={item.href}
											href={item.href}
											className="menu-item flex items-center space-x-3 p-3 rounded-xl hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200"
											onClick={() =>
												setIsMobileMenuOpen(false)
											}
										>
											<span className="font-medium text-gray-900 dark:text-white">
												{item.label}
											</span>
										</a>
									))}
								</div>
								<div className="pt-2 menu-item">
									<Button
										className="w-full bg-gray-900 dark:bg-gray-100 text-white dark:text-gray-900 rounded-xl py-3 font-semibold"
										asChild
									>
										<Link to="/auth/login">
											Get Started
										</Link>
									</Button>
								</div>
							</div>
						</div>
					</div>
				</div>
			)}
		</>
	);
}
