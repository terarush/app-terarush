import { useEffect, useRef } from "react";
import { Code, Eye, Award, Github } from "lucide-react";
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { whyTerarushContent } from "@/content/config";

gsap.registerPlugin(ScrollTrigger);

const iconMap = {
	code: Code,
	eye: Eye,
	award: Award,
	github: Github,
};

export default function WhyUs() {
	const sectionRef = useRef<HTMLElement>(null);
	const headerRef = useRef<HTMLDivElement>(null);
	const reasonsRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		const ctx = gsap.context(() => {
			// Header
			if (headerRef.current) {
				gsap.fromTo(
					headerRef.current.children,
					{ y: 30, opacity: 0 },
					{
						y: 0,
						opacity: 1,
						duration: 0.8,
						stagger: 0.15,
						ease: "power3.out",
						scrollTrigger: {
							trigger: headerRef.current,
							start: "top 80%",
							toggleActions: "play none none none",
						},
					}
				);
			}

			// Reasons cards
			if (reasonsRef.current) {
				const cards = reasonsRef.current.querySelectorAll(".reason-card");
				gsap.fromTo(
					cards,
					{ y: 50, opacity: 0, scale: 0.95 },
					{
						y: 0,
						opacity: 1,
						scale: 1,
						duration: 0.8,
						stagger: 0.15,
						ease: "back.out(1.2)",
						scrollTrigger: {
							trigger: reasonsRef.current,
							start: "top 75%",
							toggleActions: "play none none none",
						},
					}
				);
			}
		}, sectionRef);

		return () => ctx.revert();
	}, []);

	return (
		<section
			ref={sectionRef}
			id="why"
			className="relative overflow-hidden bg-background py-24"
		>
			{/* Background */}
			<div className="absolute inset-0 overflow-hidden opacity-30">
				<div className="absolute top-0 left-0 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
				<div className="absolute bottom-0 right-0 w-96 h-96 bg-accent/10 rounded-full blur-3xl" />
			</div>

			<div className="container mx-auto px-4 sm:px-6 relative z-10 max-w-7xl">
				{/* Header */}
				<div ref={headerRef} className="text-center mb-20">
					<Badge
						variant="secondary"
						className="px-4 py-1.5 text-sm font-medium rounded-full bg-primary/10 text-primary border border-primary/20 mb-6 opacity-0"
					>
						Why Choose Us
					</Badge>
					<h2 className="text-4xl sm:text-5xl md:text-6xl font-bold text-foreground mb-6 opacity-0">
						{whyTerarushContent.title}
					</h2>
					<p className="text-lg sm:text-xl text-muted-foreground max-w-3xl mx-auto opacity-0">
						{whyTerarushContent.subtitle}
					</p>
				</div>

				{/* Reasons Grid */}
				<div ref={reasonsRef} className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
					{whyTerarushContent.reasons.map((reason, index) => {
						const Icon = iconMap[reason.icon as keyof typeof iconMap];
						return (
							<Card
								key={index}
								className="reason-card opacity-0 group border border-border bg-card hover:shadow-xl hover:border-primary/50 transition-all duration-300"
							>
								<CardContent className="p-8">
									<div className="mb-6">
										<div className="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-gradient-to-br from-primary/20 to-primary/10 group-hover:scale-110 transition-transform duration-300">
											<Icon className="h-7 w-7 text-primary" />
										</div>
									</div>
									<h3 className="text-xl font-bold text-foreground mb-3 group-hover:text-primary transition-colors duration-300">
										{reason.title}
									</h3>
									<p className="text-muted-foreground leading-relaxed">
										{reason.description}
									</p>
								</CardContent>
							</Card>
						);
					})}
				</div>
			</div>
		</section>
	);
}
