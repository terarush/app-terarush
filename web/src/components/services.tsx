import { useEffect, useRef } from "react";
import { Globe, Server, Shield, Check } from "lucide-react";
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";

import { Card, CardContent } from "@/components/ui/card";
import { servicesContent } from "@/content/config";

gsap.registerPlugin(ScrollTrigger);

const iconMap = {
	globe: Globe,
	server: Server,
	shield: Shield,
};

export default function Services() {
	const sectionRef = useRef<HTMLElement>(null);
	const headerRef = useRef<HTMLDivElement>(null);
	const servicesRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		const ctx = gsap.context(() => {
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

			if (servicesRef.current) {
				const cards = servicesRef.current.querySelectorAll(".service-card");
				gsap.fromTo(
					cards,
					{ y: 60, opacity: 0, rotateX: 15 },
					{
						y: 0,
						opacity: 1,
						rotateX: 0,
						duration: 0.8,
						stagger: 0.2,
						ease: "back.out(1.2)",
						scrollTrigger: {
							trigger: servicesRef.current,
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
			id="services"
			className="relative overflow-hidden bg-background py-24"
		>
			<div className="absolute inset-0 overflow-hidden opacity-30 pointer-events-none">
				<div className="absolute top-0 left-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
				<div className="absolute bottom-0 right-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
			</div>

			<div className="container mx-auto px-4 sm:px-6 relative z-10 max-w-7xl">
				<div ref={headerRef} className="text-center mb-20">
					<h2 className="text-4xl sm:text-5xl md:text-6xl font-bold text-foreground mb-6 opacity-0">
						{servicesContent.title}
					</h2>
					<p className="text-lg sm:text-xl text-muted-foreground max-w-3xl mx-auto opacity-0">
						{servicesContent.subtitle}
					</p>
				</div>

				<div ref={servicesRef} className="grid md:grid-cols-3 gap-8">
					{servicesContent.services.map((service, index) => {
						const Icon = iconMap[service.icon as keyof typeof iconMap];
						return (
							<Card
								key={index}
								className="service-card opacity-0 group relative overflow-hidden border border-border bg-gradient-to-br from-card to-card/50 hover:shadow-2xl hover:scale-[1.02] transition-all duration-500 flex flex-col"
							>
								<div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none" />

								<CardContent className="p-8 relative z-10 flex flex-col h-full">
									<div className="mb-6">
										<div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-primary/20 to-primary/10 group-hover:scale-110 transition-transform duration-300">
											<Icon className="h-8 w-8 text-primary" />
										</div>
									</div>

									<h3 className="text-2xl font-bold text-foreground mb-4 group-hover:text-primary transition-colors duration-300">
										{service.title}
									</h3>

									<p className="text-muted-foreground mb-6 leading-relaxed flex-grow">
										{service.description}
									</p>

									<ul className="space-y-3">
										{service.features.map((feature, idx) => (
											<li
												key={idx}
												className="flex items-start gap-3 text-sm text-muted-foreground"
											>
												<Check className="h-5 w-5 text-primary flex-shrink-0 mt-0.5" />
												<span>{feature}</span>
											</li>
										))}
									</ul>
								</CardContent>
							</Card>
						);
					})}
				</div>
			</div>
		</section>
	);
}
