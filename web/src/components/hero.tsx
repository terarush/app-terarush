import { useState, useEffect, useRef } from "react";
import {
	ArrowRight,
	CheckCircle2,
	BarChart3,
	Code2,
	Bell,
	TrendingUp,
	Users,
	Zap,
	Shield,
	MapPin,
	Clock,
	Globe,
	Copy,
	Check,
	X,
} from "lucide-react";
import gsap from "gsap";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { heroContent, siteConfig } from "@/content/config";

export default function Hero() {
	const [showDialog, setShowDialog] = useState(false);
	const [copied, setCopied] = useState(false);
	const heroRef = useRef<HTMLDivElement>(null);
	const leftColumnRef = useRef<HTMLDivElement>(null);
	const rightColumnRef = useRef<HTMLDivElement>(null);
	const locationRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		// GSAP Timeline for animations
		const tl = gsap.timeline({ defaults: { ease: "power3.out" } });

		// Animate left column content
		if (leftColumnRef.current) {
			const elements =
				leftColumnRef.current.querySelectorAll(".animate-element");
			tl.fromTo(
				elements,
				{ y: 40, opacity: 0 },
				{ y: 0, opacity: 1, duration: 0.8, stagger: 0.15 },
			);
		}

		// Animate right column illustration
		if (rightColumnRef.current) {
			const cards =
				rightColumnRef.current.querySelectorAll(".floating-card");
			tl.fromTo(
				cards,
				{ scale: 0, opacity: 0, rotation: -10 },
				{
					scale: 1,
					opacity: 1,
					rotation: 0,
					duration: 0.8,
					stagger: 0.15,
					ease: "back.out(1.7)",
				},
				"-=0.6",
			);
		}

		// Animate location
		if (locationRef.current) {
			const locationItems =
				locationRef.current.querySelectorAll(".location-item");
			tl.fromTo(
				locationItems,
				{ x: -20, opacity: 0 },
				{
					x: 0,
					opacity: 1,
					duration: 0.6,
					stagger: 0.1,
				},
				"-=0.4",
			);
		}

		// Floating animation for cards
		gsap.to(".floating-card-1", {
			y: -15,
			x: 10,
			duration: 3,
			repeat: -1,
			yoyo: true,
			ease: "sine.inOut",
		});

		gsap.to(".floating-card-2", {
			y: 15,
			x: -10,
			duration: 4,
			repeat: -1,
			yoyo: true,
			ease: "sine.inOut",
		});

		gsap.to(".floating-card-3", {
			y: -10,
			x: 15,
			duration: 3.5,
			repeat: -1,
			yoyo: true,
			ease: "sine.inOut",
		});

		// Pulse animation for icons
		gsap.to(".pulse-icon", {
			scale: 1.1,
			duration: 2,
			repeat: -1,
			yoyo: true,
			ease: "sine.inOut",
			stagger: 0.3,
		});
	}, []);

	const handleStartBuilding = () => {
		// Open contact dialog
		setShowDialog(true);
	};

	const handleCopyEmail = () => {
		navigator.clipboard.writeText(siteConfig.email);
		setCopied(true);
		setTimeout(() => setCopied(false), 2000);
	};

	const handleSendEmail = () => {
		const subject = encodeURIComponent("New Project Inquiry - Terarush");
		const body = encodeURIComponent(
			`Hello Terarush Team,\n\nI'm interested in discussing a project with you.\n\nBest regards`,
		);
		window.open(
			`mailto:${siteConfig.email}?subject=${subject}&body=${body}`,
		);
	};

	return (
		<div
			ref={heroRef}
			className="relative overflow-hidden bg-background pt-24 pb-16 lg:pt-32 lg:pb-24 min-h-screen flex items-center"
			id="home"
		>
			<div className="absolute inset-0 overflow-hidden opacity-40">
				<div className="absolute top-0 right-0 w-[500px] h-[500px] bg-primary/5 rounded-full blur-3xl" />
				<div className="absolute bottom-0 left-0 w-[500px] h-[500px] bg-primary/5 rounded-full blur-3xl" />
				<div className="absolute inset-0 bg-[linear-gradient(to_right,#80808008_1px,transparent_1px),linear-gradient(to_bottom,#80808008_1px,transparent_1px)] bg-[size:40px_40px]" />
			</div>

			<div className="container mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
				<div className="grid lg:grid-cols-2 gap-12 lg:gap-16 items-center max-w-7xl mx-auto">
					<div ref={leftColumnRef} className="space-y-8">
						<div className="animate-element">
							<Badge
								variant="secondary"
								className="px-4 py-2 text-sm font-medium rounded-full bg-primary/10 text-primary border border-primary/20 inline-flex items-center gap-2"
							>
								<Zap className="w-3.5 h-3.5" />
								{heroContent.badge.text}
							</Badge>
						</div>

						<h1 className="animate-element text-4xl sm:text-5xl lg:text-6xl font-bold tracking-tight text-foreground leading-tight">
							{heroContent.title}
						</h1>

						<p className="animate-element text-lg sm:text-xl text-muted-foreground leading-relaxed max-w-xl">
							{heroContent.description}
						</p>

						{/* CTA Buttons */}
						<div className="animate-element flex flex-col sm:flex-row gap-4">
							<Button
								onClick={handleStartBuilding}
								size="lg"
								className="h-14 px-8 rounded-xl font-bold text-lg bg-primary hover:bg-primary/90 text-primary-foreground shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-300 flex items-center justify-center gap-2 group"
							>
								Start Building
								<ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
							</Button>
							<Button
								asChild
								size="lg"
								variant="outline"
								className="h-14 px-8 rounded-xl font-bold text-lg border-2 border-border hover:border-primary hover:bg-accent hover:scale-105 transition-all duration-300 flex items-center justify-center gap-2"
							>
								<a href="#contact">Contact Us</a>
							</Button>
						</div>

						<div className="animate-element flex flex-wrap gap-4">
							{heroContent.features.map((feature, index) => (
								<div
									key={index}
									className="flex items-center gap-2 text-sm text-muted-foreground"
								>
									<CheckCircle2 className="w-4 h-4 text-primary" />
									<span>{feature}</span>
								</div>
							))}
						</div>

						{/* Location Info */}
						<div className="animate-element pt-8 border-t border-border space-y-4">
							<p className="text-sm font-semibold text-foreground">
								{heroContent.location.title}
							</p>
							<div
								ref={locationRef}
								className="grid grid-cols-1 sm:grid-cols-3 gap-4"
							>
								<div className="location-item flex items-start gap-3 opacity-0">
									<div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
										<MapPin className="w-5 h-5 text-primary" />
									</div>
									<div>
										<p className="text-sm font-medium text-foreground">
											Location
										</p>
										<p className="text-xs text-muted-foreground">
											{heroContent.location.address}
										</p>
									</div>
								</div>
								<div className="location-item flex items-start gap-3 opacity-0">
									<div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
										<Clock className="w-5 h-5 text-primary" />
									</div>
									<div>
										<p className="text-sm font-medium text-foreground">
											Timezone
										</p>
										<p className="text-xs text-muted-foreground">
											{heroContent.location.timezone}
										</p>
									</div>
								</div>
								<div className="location-item flex items-start gap-3 opacity-0">
									<div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
										<Globe className="w-5 h-5 text-primary" />
									</div>
									<div>
										<p className="text-sm font-medium text-foreground">
											Work Mode
										</p>
										<p className="text-xs text-muted-foreground">
											{heroContent.location.availability}
										</p>
									</div>
								</div>
							</div>
						</div>
					</div>

					<div
						ref={rightColumnRef}
						className="relative h-[500px] lg:h-[600px] hidden lg:block"
					>
						<div className="relative w-full h-full perspective-1000">
							<div className="floating-card floating-card-1 absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-80 h-96 bg-card rounded-3xl shadow-2xl border border-border overflow-hidden">
								<div className="p-6 h-full flex flex-col">
									<div className="flex items-center justify-between mb-6">
										<div className="flex items-center gap-3">
											<div className="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center">
												<BarChart3 className="w-5 h-5 text-primary pulse-icon" />
											</div>
											<div>
												<h3 className="font-semibold text-sm text-foreground">
													Analytics
												</h3>
												<p className="text-xs text-muted-foreground">
													Real-time data
												</p>
											</div>
										</div>
										<div className="flex gap-1">
											<div className="w-2 h-2 rounded-full bg-green-500" />
											<div className="w-2 h-2 rounded-full bg-yellow-500" />
											<div className="w-2 h-2 rounded-full bg-red-500" />
										</div>
									</div>

									<div className="flex-1 bg-muted/30 rounded-2xl p-4 mb-4">
										<div className="flex items-end justify-between h-full gap-2">
											{[40, 70, 50, 80, 60, 90, 75].map(
												(height, i) => (
													<div
														key={i}
														className="flex-1 bg-primary/20 rounded-t-lg relative overflow-hidden"
														style={{
															height: `${height}%`,
														}}
													>
														<div className="absolute bottom-0 left-0 right-0 h-1/2 bg-primary/40" />
													</div>
												),
											)}
										</div>
									</div>

									{/* Stats row */}
									<div className="grid grid-cols-2 gap-3">
										<div className="bg-muted/30 rounded-xl p-3">
											<div className="flex items-center gap-2 mb-1">
												<TrendingUp className="w-3.5 h-3.5 text-green-500" />
												<span className="text-xs text-muted-foreground">
													Growth
												</span>
											</div>
											<p className="text-lg font-bold text-foreground">
												+23%
											</p>
										</div>
										<div className="bg-muted/30 rounded-xl p-3">
											<div className="flex items-center gap-2 mb-1">
												<Users className="w-3.5 h-3.5 text-primary" />
												<span className="text-xs text-muted-foreground">
													Users
												</span>
											</div>
											<p className="text-lg font-bold text-foreground">
												2.5K
											</p>
										</div>
									</div>
								</div>
							</div>

							{/* Floating Card 1 - Code Editor */}
							<div className="floating-card floating-card-2 absolute top-12 left-8 w-64 h-48 bg-card rounded-2xl shadow-xl border border-border overflow-hidden rotate-6">
								<div className="p-4 h-full">
									<div className="flex items-center gap-2 mb-4">
										<Code2 className="w-4 h-4 text-primary pulse-icon" />
										<span className="text-xs font-medium text-foreground">
											terarush.tsx
										</span>
									</div>
									<div className="space-y-2">
										<div className="h-2 bg-primary/20 rounded w-3/4" />
										<div className="h-2 bg-muted/50 rounded w-full" />
										<div className="h-2 bg-muted/50 rounded w-5/6" />
										<div className="h-2 bg-primary/20 rounded w-2/3" />
										<div className="h-2 bg-muted/50 rounded w-4/5" />
										<div className="h-2 bg-muted/50 rounded w-full" />
									</div>
								</div>
							</div>

							{/* Floating Card 2 - Notifications */}
							<div className="floating-card floating-card-3 absolute bottom-16 right-8 w-72 bg-card rounded-2xl shadow-xl border border-border overflow-hidden -rotate-3">
								<div className="p-4">
									<div className="flex items-center gap-2 mb-4">
										<Bell className="w-4 h-4 text-primary pulse-icon" />
										<span className="text-xs font-medium text-foreground">
											Notifications
										</span>
									</div>
									<div className="space-y-3">
										{[
											{
												icon: CheckCircle2,
												text: "Project deployed",
												color: "text-green-500",
											},
											{
												icon: Shield,
												text: "Security updated",
												color: "text-blue-500",
											},
											{
												icon: TrendingUp,
												text: "Performance +15%",
												color: "text-primary",
											},
										].map((notif, i) => (
											<div
												key={i}
												className="flex items-center gap-3 p-2 bg-muted/20 rounded-xl"
											>
												<notif.icon
													className={`w-4 h-4 ${notif.color}`}
												/>
												<span className="text-xs text-muted-foreground">
													{notif.text}
												</span>
											</div>
										))}
									</div>
								</div>
							</div>

							{/* Decorative circles */}
							<div className="absolute -top-4 -right-4 w-24 h-24 bg-primary/5 rounded-full blur-2xl" />
							<div className="absolute -bottom-4 -left-4 w-32 h-32 bg-primary/5 rounded-full blur-2xl" />
						</div>
					</div>
				</div>
			</div>

			{/* Contact Dialog */}
			<Dialog open={showDialog} onOpenChange={setShowDialog}>
				<DialogContent className="sm:max-w-md">
					<DialogHeader>
						<DialogTitle className="flex items-center justify-between">
							<span>Get In Touch</span>
							<Button
								variant="ghost"
								size="icon"
								onClick={() => setShowDialog(false)}
								className="h-6 w-6"
							>
								<X className="h-4 w-4" />
							</Button>
						</DialogTitle>
						<DialogDescription>
							Ready to start your project? Reach out to us and
							let's build something amazing together.
						</DialogDescription>
					</DialogHeader>
					<div className="space-y-6">
						{/* Email Contact */}
						<div className="space-y-3">
							<div className="flex items-center gap-3 p-4 bg-muted/50 rounded-lg">
								<div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
									<ArrowRight className="w-5 h-5 text-primary" />
								</div>
								<div className="flex-1 min-w-0">
									<p className="text-sm font-medium text-muted-foreground">
										Email us at
									</p>
									<p className="text-base font-semibold text-foreground truncate">
										{siteConfig.email}
									</p>
								</div>
							</div>
						</div>

						{/* Action Buttons */}
						<div className="flex flex-col gap-3">
							<Button
								onClick={handleSendEmail}
								size="lg"
								className="w-full gap-2"
							>
								<ArrowRight className="w-4 h-4" />
								Send Email
							</Button>
							<Button
								onClick={handleCopyEmail}
								size="lg"
								variant="outline"
								className="w-full gap-2"
							>
								{copied ? (
									<>
										<Check className="w-4 h-4" />
										Copied!
									</>
								) : (
									<>
										<Copy className="w-4 h-4" />
										Copy Email Address
									</>
								)}
							</Button>
						</div>

						{/* Additional Info */}
						<div className="pt-4 border-t border-border space-y-2">
							<p className="text-sm font-medium text-foreground">
								Quick Links
							</p>
							<div className="flex flex-col gap-2">
								<a
									href={siteConfig.links.github}
									target="_blank"
									rel="noopener noreferrer"
									className="text-sm text-muted-foreground hover:text-primary transition-colors flex items-center gap-2"
								>
									<span>→</span> GitHub
								</a>
								<a
									href={siteConfig.links.linkedin}
									target="_blank"
									rel="noopener noreferrer"
									className="text-sm text-muted-foreground hover:text-primary transition-colors flex items-center gap-2"
								>
									<span>→</span> LinkedIn
								</a>
								<a
									href={siteConfig.links.twitter}
									target="_blank"
									rel="noopener noreferrer"
									className="text-sm text-muted-foreground hover:text-primary transition-colors flex items-center gap-2"
								>
									<span>→</span> Twitter
								</a>
							</div>
						</div>
					</div>
				</DialogContent>
			</Dialog>
		</div>
	);
}
