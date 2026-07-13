import { useEffect, useRef } from "react";
import { Globe, Server, Shield, Check, ArrowUpRight } from "lucide-react";
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import { useReducedMotion } from "motion/react";

import { Button } from "@/components/ui/button";
import { servicesContent } from "@/content/config";

gsap.registerPlugin(ScrollTrigger);

const iconMap = {
	globe: Globe,
	server: Server,
	shield: Shield,
};

export default function Services() {
	const sectionRef = useRef<HTMLElement>(null);
	const reduce = useReducedMotion();

	useEffect(() => {
		if (reduce) return;
		const ctx = gsap.context(() => {
			const rows = sectionRef.current?.querySelectorAll(".service-row");
			if (!rows) return;
			gsap.fromTo(
				rows,
				{ y: 40, opacity: 0 },
				{
					y: 0,
					opacity: 1,
					duration: 0.7,
					stagger: 0.15,
					ease: "power3.out",
					scrollTrigger: {
						trigger: sectionRef.current,
						start: "top 75%",
						toggleActions: "play none none none",
					},
				},
			);
		}, sectionRef);
		return () => ctx.revert();
	}, [reduce]);

	return (
		<section
			ref={sectionRef}
			id="services"
			className="relative overflow-hidden bg-background py-24"
		>
			<div className="absolute inset-0 overflow-hidden opacity-30 pointer-events-none">
				<div className="absolute top-0 left-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
				<div className="absolute bottom-0 right-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
			</div>

			<div className="container mx-auto px-4 sm:px-6 relative z-10 max-w-5xl">
				<div className="text-center mb-20">
					<p className="text-sm font-semibold text-primary tracking-wide uppercase mb-4">
						What We Do
					</p>
					<h2 className="text-4xl sm:text-5xl font-bold text-foreground mb-6">
						{servicesContent.title}
					</h2>
					<p className="text-lg text-muted-foreground max-w-2xl mx-auto">
						{servicesContent.subtitle}
					</p>
				</div>

				<div className="space-y-12">
					{servicesContent.services.map((service, index) => {
						const Icon = iconMap[service.icon as keyof typeof iconMap];
						return (
							<div
								key={index}
								className="service-row grid md:grid-cols-12 gap-6 md:gap-10 items-start p-6 md:p-0 bg-card md:bg-transparent rounded-2xl border border-border md:border-0"
							>
								{/* Icon column */}
								<div className="md:col-span-2 flex md:justify-center pt-1">
									<div className="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-gradient-to-br from-primary/20 to-primary/10">
										<Icon className="h-7 w-7 text-primary" />
									</div>
								</div>

								{/* Content */}
								<div className="md:col-span-10 space-y-4">
									<h3 className="text-2xl font-bold text-foreground">
										{service.title}
									</h3>
									<p className="text-muted-foreground leading-relaxed max-w-prose">
										{service.description}
									</p>
									<div className="flex flex-wrap gap-2 pt-1">
										{service.features.map((feature, idx) => (
											<span
												key={idx}
												className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-primary/5 text-primary text-sm font-medium border border-primary/10"
											>
												<Check className="h-3.5 w-3.5" />
												{feature}
											</span>
										))}
									</div>
								</div>
							</div>
						);
					})}
				</div>

				<div className="text-center mt-16">
					<p className="text-muted-foreground text-sm mb-4">
						Not sure what you need?
					</p>
					<Button asChild variant="outline" className="rounded-xl">
						<a href="#contact">
							Let&apos;s Talk
							<ArrowUpRight className="ml-2 h-4 w-4" />
						</a>
					</Button>
				</div>
			</div>
		</section>
	);
}
