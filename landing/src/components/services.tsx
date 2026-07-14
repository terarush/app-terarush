import { Globe, Server, Shield, ArrowUpRight } from "lucide-react";
import { motion, useReducedMotion } from "motion/react";
import { Button } from "@/components/ui/button";
import { servicesContent } from "@/content/config";

const iconMap = {
	globe: Globe,
	server: Server,
	shield: Shield,
};

export default function Services() {
	const reduce = useReducedMotion();

	return (
		<section
			id="services"
			className="relative overflow-hidden bg-background pt-32 pb-16"
		>
			<div className="container mx-auto px-4 sm:px-6 relative z-10 max-w-7xl">
				<div className="mb-20 max-w-2xl">
					<h2 className="text-4xl md:text-5xl font-bold tracking-tight text-foreground mb-6">
						{servicesContent.title}
					</h2>
					<p className="text-xl text-muted-foreground leading-relaxed">
						{servicesContent.subtitle}
					</p>
				</div>

				<div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
					{servicesContent.services.map((service, index) => {
						const Icon = iconMap[service.icon as keyof typeof iconMap];
						return (
							<motion.div
								key={index}
								initial={reduce ? false : { opacity: 0, y: 20 }}
								whileInView={{ opacity: 1, y: 0 }}
								viewport={{ once: true, margin: "-100px" }}
								transition={{ duration: 0.5, delay: index * 0.1 }}
								className="group relative p-8 rounded-3xl bg-secondary/50 border border-border/50 hover:bg-secondary transition-colors duration-300 flex flex-col h-full"
							>
								<div className="mb-8">
									<div className="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-background shadow-sm border border-border mb-6">
										<Icon className="h-6 w-6 text-foreground" />
									</div>
									<h3 className="text-2xl font-bold text-foreground mb-3 tracking-tight">
										{service.title}
									</h3>
									<p className="text-muted-foreground leading-relaxed">
										{service.description}
									</p>
								</div>

								<div className="mt-auto pt-8">
									<ul className="space-y-3">
										{service.features.map((feature, idx) => (
											<li key={idx} className="flex items-center gap-3 text-sm text-muted-foreground">
												<div className="w-1.5 h-1.5 rounded-full bg-primary/50" />
												{feature}
											</li>
										))}
									</ul>
								</div>
							</motion.div>
						);
					})}
				</div>

				<div className="mt-20 flex justify-center">
					<Button asChild variant="outline" className="rounded-full h-12 px-8 font-medium">
						<a href="#contact">
							Explore Custom Solutions
							<ArrowUpRight className="ml-2 h-4 w-4" />
						</a>
					</Button>
				</div>
			</div>
		</section>
	);
}