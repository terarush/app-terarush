import { useEffect, useRef } from "react";
import { Home, ArrowLeft, Search, Mail } from "lucide-react";
import gsap from "gsap";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { siteConfig } from "@/content/config";

export default function NotFound() {
	const containerRef = useRef<HTMLDivElement>(null);
	const numberRef = useRef<HTMLDivElement>(null);
	const contentRef = useRef<HTMLDivElement>(null);
	const cardsRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		const ctx = gsap.context(() => {
			// Animate 404 number
			if (numberRef.current) {
				gsap.fromTo(
					numberRef.current,
					{
						scale: 0.8,
						opacity: 0,
						rotateX: -90,
					},
					{
						scale: 1,
						opacity: 1,
						rotateX: 0,
						duration: 1,
						ease: "back.out(1.7)",
					}
				);
			}

			// Animate content
			if (contentRef.current) {
				gsap.fromTo(
					contentRef.current.children,
					{
						y: 30,
						opacity: 0,
					},
					{
						y: 0,
						opacity: 1,
						duration: 0.8,
						stagger: 0.15,
						delay: 0.3,
						ease: "power3.out",
					}
				);
			}

			// Animate cards
			if (cardsRef.current) {
				gsap.fromTo(
					cardsRef.current.children,
					{
						y: 40,
						opacity: 0,
						scale: 0.95,
					},
					{
						y: 0,
						opacity: 1,
						scale: 1,
						duration: 0.8,
						stagger: 0.1,
						delay: 0.6,
						ease: "back.out(1.2)",
					}
				);
			}

			// Floating animation for 404
			if (numberRef.current) {
				gsap.to(numberRef.current, {
					y: -20,
					duration: 2,
					repeat: -1,
					yoyo: true,
					ease: "sine.inOut",
					delay: 1,
				});
			}
		}, containerRef);

		return () => ctx.revert();
	}, []);

	const quickLinks = [
		{
			title: "Go Home",
			description: "Return to our homepage",
			icon: Home,
			href: "/",
			color: "bg-primary",
		},
		{
			title: "Go Back",
			description: "Return to previous page",
			icon: ArrowLeft,
			href: "#",
			color: "bg-chart-1",
			onClick: () => window.history.back(),
		},
		{
			title: "Search",
			description: "Find what you're looking for",
			icon: Search,
			href: "/#services",
			color: "bg-chart-2",
		},
		{
			title: "Contact Us",
			description: "We're here to help",
			icon: Mail,
			href: "/#contact",
			color: "bg-chart-3",
		},
	];

	return (
		<div
			ref={containerRef}
			className="min-h-screen bg-background flex items-center justify-center relative overflow-hidden"
		>
			{/* Background decorations */}
			<div className="absolute inset-0 overflow-hidden opacity-30">
				<div className="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
				<div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
			</div>

			<div className="container mx-auto px-4 sm:px-6 relative z-10 max-w-5xl py-12">
				{/* 404 Number */}
				<div
					ref={numberRef}
					className="text-center mb-8"
					style={{ perspective: "1000px" }}
				>
					<h1 className="text-[180px] sm:text-[240px] md:text-[300px] font-bold leading-none opacity-0">
						<span className="bg-gradient-to-br from-primary via-primary/80 to-primary/60 bg-clip-text text-transparent">
							404
						</span>
					</h1>
				</div>

				{/* Content */}
				<div ref={contentRef} className="text-center mb-12">
					<h2 className="text-3xl sm:text-4xl md:text-5xl font-bold text-foreground mb-4 opacity-0">
						Page Not Found
					</h2>
					<p className="text-lg sm:text-xl text-muted-foreground max-w-2xl mx-auto mb-6 opacity-0">
						Oops! The page you're looking for doesn't exist. It might have been
						moved or deleted.
					</p>
					<p className="text-md text-muted-foreground opacity-0">
						Don't worry, we'll help you get back on track.
					</p>
				</div>

				{/* Quick Links */}
				<div
					ref={cardsRef}
					className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-12"
				>
					{quickLinks.map((link) => {
						const Icon = link.icon;
						return (
							<Card
								key={link.title}
								className="opacity-0 border border-border bg-card hover:shadow-lg hover:scale-105 transition-all duration-300 group"
							>
								<CardContent className="p-6">
									{link.onClick ? (
										<button
											onClick={link.onClick}
											className="w-full text-left"
										>
											<div
												className={`w-12 h-12 ${link.color} rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300`}
											>
												<Icon className="h-6 w-6 text-white" />
											</div>
											<h3 className="text-lg font-semibold text-foreground mb-2">
												{link.title}
											</h3>
											<p className="text-sm text-muted-foreground">
												{link.description}
											</p>
										</button>
									) : (
										<a href={link.href} className="block">
											<div
												className={`w-12 h-12 ${link.color} rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300`}
											>
												<Icon className="h-6 w-6 text-white" />
											</div>
											<h3 className="text-lg font-semibold text-foreground mb-2">
												{link.title}
											</h3>
											<p className="text-sm text-muted-foreground">
												{link.description}
											</p>
										</a>
									)}
								</CardContent>
							</Card>
						);
					})}
				</div>

				{/* Help Section */}
				<Card className="bg-gradient-to-br from-primary/10 to-primary/5 border border-primary/20">
					<CardContent className="p-8 text-center">
						<h3 className="text-2xl font-bold text-foreground mb-3">
							Need Help?
						</h3>
						<p className="text-muted-foreground mb-6 max-w-2xl mx-auto">
							If you believe this is an error or need assistance, feel free to
							reach out to our team.
						</p>
						<div className="flex flex-col sm:flex-row gap-4 justify-center">
							<Button
								asChild
								className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl h-12 px-8 font-semibold shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-300"
							>
								<a href="/#contact">Contact Support</a>
							</Button>
							<Button
								asChild
								variant="outline"
								className="border-border hover:bg-primary/10 hover:text-primary rounded-xl h-12 px-8 font-semibold transition-all duration-300 hover:scale-105"
							>
								<a href={`mailto:${siteConfig.primaryEmail}`}>
									Email Us
								</a>
							</Button>
						</div>
					</CardContent>
				</Card>
			</div>
		</div>
	);
}
