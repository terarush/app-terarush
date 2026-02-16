import { useEffect, useRef } from "react";
import { Wrench, Layers, Package, Check } from "lucide-react";
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { openSourceContent } from "@/content/config";

gsap.registerPlugin(ScrollTrigger);

const iconMap = {
	wrench: Wrench,
	layers: Layers,
	package: Package,
};

export default function OpenSource() {
	const sectionRef = useRef<HTMLElement>(null);
	const headerRef = useRef<HTMLDivElement>(null);
	const philosophyRef = useRef<HTMLDivElement>(null);
	const projectsRef = useRef<HTMLDivElement>(null);
	const communityRef = useRef<HTMLDivElement>(null);

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

			// Philosophy
			if (philosophyRef.current) {
				gsap.fromTo(
					philosophyRef.current,
					{ y: 40, opacity: 0 },
					{
						y: 0,
						opacity: 1,
						duration: 0.8,
						ease: "power3.out",
						scrollTrigger: {
							trigger: philosophyRef.current,
							start: "top 80%",
							toggleActions: "play none none none",
						},
					}
				);
			}

			// Projects
			if (projectsRef.current) {
				const cards = projectsRef.current.querySelectorAll(".project-card");
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
							trigger: projectsRef.current,
							start: "top 75%",
							toggleActions: "play none none none",
						},
					}
				);
			}

			// Community
			if (communityRef.current) {
				gsap.fromTo(
					communityRef.current.children,
					{ y: 40, opacity: 0 },
					{
						y: 0,
						opacity: 1,
						duration: 0.8,
						stagger: 0.1,
						ease: "power3.out",
						scrollTrigger: {
							trigger: communityRef.current,
							start: "top 80%",
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
			id="opensource"
			className="relative overflow-hidden bg-background py-24"
		>
			{/* Background */}
			<div className="absolute inset-0 overflow-hidden opacity-30">
				<div className="absolute top-1/3 right-0 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
				<div className="absolute bottom-1/3 left-0 w-96 h-96 bg-accent/10 rounded-full blur-3xl" />
			</div>

			<div className="container mx-auto px-4 sm:px-6 relative z-10 max-w-7xl">
				{/* Header */}
				<div ref={headerRef} className="text-center mb-20">
					<Badge
						variant="secondary"
						className="px-4 py-1.5 text-sm font-medium rounded-full bg-primary/10 text-primary border border-primary/20 mb-6 opacity-0"
					>
						Open Source
					</Badge>
					<h2 className="text-4xl sm:text-5xl md:text-6xl font-bold text-foreground mb-6 opacity-0">
						{openSourceContent.title}
					</h2>
					<p className="text-lg sm:text-xl text-muted-foreground max-w-3xl mx-auto opacity-0">
						{openSourceContent.subtitle}
					</p>
				</div>

				{/* Philosophy */}
				<div
					ref={philosophyRef}
					className="mb-20 opacity-0 bg-card rounded-3xl p-8 sm:p-12 border border-border"
				>
					<h3 className="text-3xl font-bold text-foreground mb-4 text-center">
						{openSourceContent.philosophy.headline}
					</h3>
					<p className="text-lg text-muted-foreground text-center max-w-3xl mx-auto">
						{openSourceContent.philosophy.description}
					</p>
				</div>

				{/* Projects */}
				<div className="mb-20">
					<h3 className="text-3xl font-bold text-foreground mb-12 text-center">
						{openSourceContent.projects.title}
					</h3>
					<div ref={projectsRef} className="grid md:grid-cols-3 gap-8">
						{openSourceContent.projects.types.map((project, index) => {
							const Icon = iconMap[project.icon as keyof typeof iconMap];
							return (
								<Card
									key={index}
									className="project-card opacity-0 group border border-border bg-card hover:shadow-xl hover:scale-105 transition-all duration-300"
								>
									<CardContent className="p-8 text-center">
										<div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-primary/20 to-primary/10 mb-6 group-hover:scale-110 transition-transform duration-300">
											<Icon className="h-8 w-8 text-primary" />
										</div>
										<h4 className="text-xl font-bold text-foreground mb-3">
											{project.title}
										</h4>
										<p className="text-muted-foreground">
											{project.description}
										</p>
									</CardContent>
								</Card>
							);
						})}
					</div>
				</div>

				{/* Community */}
				<div
					ref={communityRef}
					className="bg-gradient-to-br from-primary/10 to-primary/5 rounded-3xl p-8 sm:p-12 border border-primary/20"
				>
					<h3 className="text-3xl font-bold text-foreground mb-4 text-center opacity-0">
						{openSourceContent.community.title}
					</h3>
					<p className="text-lg text-muted-foreground text-center mb-8 max-w-3xl mx-auto opacity-0">
						{openSourceContent.community.description}
					</p>
					<div className="grid md:grid-cols-2 gap-4 max-w-3xl mx-auto opacity-0">
						{openSourceContent.community.highlights.map((highlight, index) => (
							<div
								key={index}
								className="flex items-center gap-3 text-foreground"
							>
								<Check className="h-5 w-5 text-primary flex-shrink-0" />
								<span>{highlight}</span>
							</div>
						))}
					</div>
				</div>
			</div>
		</section>
	);
}
