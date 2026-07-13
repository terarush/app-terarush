import { useEffect, useRef } from "react";
import {
	Mail,
	Phone,
	MapPin,
	Clock,
	MessageCircle,
	Calendar,
	Twitter,
	Linkedin,
	Github,
	Instagram,
	ExternalLink,
	Globe,
	Shield,
} from "lucide-react";
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import { useReducedMotion } from "motion/react";

import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { siteConfig } from "@/content/config";

gsap.registerPlugin(ScrollTrigger);

export default function Contact() {
	const sectionRef = useRef<HTMLElement>(null);
	const headerRef = useRef<HTMLDivElement>(null);
	const cardsRef = useRef<HTMLDivElement>(null);
	const reduce = useReducedMotion();

	useEffect(() => {
		if (reduce) return;
		const ctx = gsap.context(() => {
			if (headerRef.current) {
				gsap.fromTo(
					headerRef.current.children,
					{ y: 20, opacity: 0 },
					{
						y: 0,
						opacity: 1,
						duration: 0.6,
						stagger: 0.1,
						ease: "power3.out",
						scrollTrigger: {
							trigger: headerRef.current,
							start: "top 80%",
							toggleActions: "play none none none",
						},
					}
				);
			}

			if (cardsRef.current) {
				const cards = cardsRef.current.querySelectorAll(".contact-card");
				gsap.fromTo(
					cards,
					{ y: 20, opacity: 0 },
					{
						y: 0,
						opacity: 1,
						duration: 0.5,
						stagger: 0.1,
						ease: "power3.out",
						scrollTrigger: {
							trigger: cardsRef.current,
							start: "top 80%",
							toggleActions: "play none none none",
						},
					}
				);
			}
		}, sectionRef);

		return () => ctx.revert();
	}, [reduce]);

	return (
		<section
			ref={sectionRef}
			className="bg-background py-24"
			id="contact"
		>
			<div className="container mx-auto px-4 sm:px-6">
				<div ref={headerRef} className="text-center mb-16">
					<h2 className="text-4xl font-bold mb-4 text-foreground opacity-0">
						Let's <span className="text-primary">talk</span>
					</h2>
					<p className="text-lg text-muted-foreground max-w-2xl mx-auto opacity-0">
						Multiple ways to reach us. Pick what works best for you.
					</p>
				</div>

				<div className="max-w-6xl mx-auto">
					<div
						ref={cardsRef}
						className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 auto-rows-auto"
					>
						{/* Email Us - Large Card */}
						<div className="contact-card lg:col-span-2 lg:row-span-2 opacity-0">
							<Card className="h-full p-8 bg-card border border-border hover:shadow-lg transition-all duration-300">
								<CardContent className="p-0 h-full flex flex-col">
									<div className="flex items-center space-x-3 mb-6">
										<div className="w-12 h-12 bg-primary rounded-xl flex items-center justify-center">
											<Mail className="h-6 w-6 text-primary-foreground" />
										</div>
										<div>
											<h3 className="text-xl font-semibold text-foreground">
												Email Us
											</h3>
										</div>
									</div>
									<div className="flex-1">
										<p className="lg:text-3xl text-2xl font-bold mb-2 text-foreground">
											{siteConfig.primaryEmail}
										</p>
										<p className="text-muted-foreground mb-6">
											Our main inbox. We typically respond within 2 hours during
											business days.
										</p>
										<div className="space-y-3">
											<div className="flex items-center space-x-2 text-sm">
												<Clock className="h-4 w-4 text-primary" />
												<span>Average response: 2 hours</span>
											</div>
											<div className="flex items-center space-x-2 text-sm">
												<Shield className="h-4 w-4 text-primary" />
												<span>Secure and confidential</span>
											</div>
										</div>
									</div>
									<Button
										asChild
										className="w-full mt-6 bg-primary hover:bg-primary/80 text-primary-foreground"
									>
										<a href={`mailto:${siteConfig.primaryEmail}`}>
											Send Email
											<ExternalLink className="ml-2 h-4 w-4" />
										</a>
									</Button>
								</CardContent>
							</Card>
						</div>

						{/* Call Us */}
						<div className="contact-card opacity-0">
							<Card className="h-full p-6 hover:shadow-lg transition-all duration-300 border border-border bg-card">
								<CardContent className="p-0">
									<div className="w-10 h-10 bg-chart-1 rounded-lg flex items-center justify-center mb-4">
										<Phone className="h-5 w-5 text-white" />
									</div>
									<h3 className="font-semibold mb-2 text-foreground">
										Call Us
									</h3>
									<p className="text-lg font-mono text-foreground">
										{siteConfig.phone}
									</p>
									<p className="text-sm text-muted-foreground mt-2">
										Mon-Fri, 9AM-10PM PST
									</p>
									<Button
										asChild
										size="sm"
										className="w-full mt-4 bg-chart-1 hover:bg-chart-1/80 text-white"
									>
										<a href={`tel:${siteConfig.phone}`}>Call Now</a>
									</Button>
								</CardContent>
							</Card>
						</div>

						{/* WhatsApp */}
						<div className="contact-card opacity-0">
							<Card className="h-full p-6 hover:shadow-lg transition-all duration-300 border border-border bg-card">
								<CardContent className="flex flex-col h-full p-0">
									<div>
										<div className="w-10 h-10 bg-chart-2 rounded-lg flex items-center justify-center mb-4">
											<MessageCircle className="h-5 w-5 text-white" />
										</div>
										<h3 className="font-semibold mb-2 text-foreground">
											WhatsApp
										</h3>
										<p className="text-sm text-muted-foreground mb-4">
											Instant messaging
										</p>
									</div>
									<div className="flex flex-1 items-end">
										<Button
											asChild
											size="sm"
											className="w-full bg-chart-2 hover:bg-chart-2/80 text-white"
										>
											<a
												href={`https://wa.me/${siteConfig.whatsapp}`}
												target="_blank"
												rel="noopener noreferrer"
											>
												Message
												<ExternalLink className="ml-2 h-3 w-3" />
											</a>
										</Button>
									</div>
								</CardContent>
							</Card>
						</div>

						{/* Business Hours */}
						<div className="contact-card lg:col-span-2 opacity-0">
							<Card className="h-full p-6 bg-accent/50 border border-border hover:shadow-lg transition-all duration-300">
								<CardContent className="p-0">
									<div className="flex items-center justify-between mb-4">
										<div className="flex items-center space-x-2">
											<Clock className="h-5 w-5 text-primary" />
											<h3 className="font-semibold text-foreground">
												Business Hours
											</h3>
										</div>
									</div>
									<div className="grid grid-cols-2 gap-4 text-sm">
										<div>
											<p className="font-medium text-foreground">Weekdays</p>
											<p className="text-muted-foreground">
												{siteConfig.businessHours.weekdays}
											</p>
										</div>
										<div>
											<p className="font-medium text-foreground">Weekends</p>
											<p className="text-muted-foreground">
												{siteConfig.businessHours.weekends}
											</p>
										</div>
									</div>
								</CardContent>
							</Card>
						</div>

						{/* Schedule Call */}
						<div className="contact-card opacity-0">
							<Card className="h-full p-6 hover:shadow-lg transition-all duration-300 border border-border bg-card">
								<CardContent className="p-0">
									<div className="w-10 h-10 bg-primary rounded-lg flex items-center justify-center mb-4">
										<Calendar className="h-5 w-5 text-primary-foreground" />
									</div>
									<h3 className="font-semibold mb-2 text-foreground">
										Schedule Call
									</h3>
									<p className="text-sm text-muted-foreground mb-4">
										30-min consultation
									</p>
									<Button
										asChild
										size="sm"
										className="w-full bg-primary hover:bg-primary/80 text-primary-foreground"
									>
										<a
											href={`https://wa.me/${siteConfig.whatsapp}`}
											target="_blank"
											rel="noopener noreferrer"
										>
											Book Now
											<ExternalLink className="ml-2 h-3 w-3" />
										</a>
									</Button>
								</CardContent>
							</Card>
						</div>

						{/* Visit Us */}
						<div className="contact-card opacity-0">
							<Card className="h-full p-6 hover:shadow-lg transition-all duration-300 border border-border bg-card">
								<CardContent className="p-0">
									<div className="w-10 h-10 bg-chart-3 rounded-lg flex items-center justify-center mb-4">
										<MapPin className="h-5 w-5 text-white" />
									</div>
									<h3 className="font-semibold mb-2 text-foreground">
										Visit Us
									</h3>
									<p className="text-sm text-foreground">
										{siteConfig.location.country}
									</p>
									<p className="text-sm text-muted-foreground">
										{siteConfig.location.region}
									</p>
								</CardContent>
							</Card>
						</div>

						{/* Connect With Us */}
						<div className="contact-card lg:col-span-2 opacity-0">
							<Card className="h-full p-6 bg-muted border border-border hover:shadow-lg transition-all duration-300">
								<CardContent className="p-0">
									<div className="flex items-center space-x-2 mb-4">
										<Globe className="h-5 w-5 text-muted-foreground" />
										<h3 className="font-semibold text-foreground">
											Connect With Us
										</h3>
									</div>
									<p className="text-sm text-muted-foreground mb-4">
										Follow our journey and get the latest updates
									</p>
									<div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
										{[
												{ icon: Twitter, name: "Twitter", color: "text-primary", href: siteConfig.links.twitter },
												{ icon: Linkedin, name: "LinkedIn", color: "text-primary", href: siteConfig.links.linkedin },
												{ icon: Github, name: "GitHub", color: "text-foreground", href: siteConfig.links.github },
												{ icon: Instagram, name: "Instagram", color: "text-pink-500", href: siteConfig.links.instagram },
											]
											.map((social, index) => {
												const Icon = social.icon;
												return (
													<a
														key={index}
														href={social.href}
														target="_blank"
														rel="noopener noreferrer"
														className="flex flex-col items-center p-3 rounded-lg hover:bg-card/50 transition-colors duration-200 group"
													>
														<div className="w-8 h-8 mb-2 flex items-center justify-center">
															<Icon
																className={`h-5 w-5 ${social.color} group-hover:scale-110 transition-transform duration-200`}
															/>
														</div>
														<p className="text-xs font-medium text-foreground">
															{social.name}
														</p>
													</a>
												);
											})}
									</div>
								</CardContent>
							</Card>
						</div>
					</div>
				</div>
			</div>
		</section>
	);
}
