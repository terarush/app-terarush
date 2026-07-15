import { useEffect, useRef } from "react";
import { Code, Eye, Award, Github } from "lucide-react";
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import { useReducedMotion } from "motion/react";

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
	const reduce = useReducedMotion();

	useEffect(() => {
		if (reduce) return;
		const ctx = gsap.context(() => {
			const cards = sectionRef.current?.querySelectorAll(".pillar-card");
			if (!cards) return;
			gsap.fromTo(
				cards,
				{ y: 40, opacity: 0 },
				{
					y: 0,
					opacity: 1,
					duration: 0.8,
					stagger: 0.2,
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

	const pillars = [
		{
			title: "Engineering Excellence",
			items: [
				whyTerarushContent.reasons[0],
				whyTerarushContent.reasons[2],
			],
			icon: Code,
			gradient: "from-primary/5 via-primary/3 to-transparent",
		},
		{
			title: "Partnership First",
			items: [
				whyTerarushContent.reasons[1],
				whyTerarushContent.reasons[3],
			],
			icon: Eye,
			gradient: "from-accent via-accent/50 to-transparent",
		},
	];

	return (
		<section
			ref={sectionRef}
			id="why"
			className="relative overflow-hidden bg-background py-24"
		>
			<div className="absolute inset-0 overflow-hidden opacity-20 pointer-events-none">
				<div className="absolute top-0 left-0 w-[500px] h-[500px] bg-primary/10 rounded-full blur-3xl" />
				<div className="absolute bottom-0 right-0 w-[500px] h-[500px] bg-accent/10 rounded-full blur-3xl" />
			</div>

			<div className="container mx-auto px-4 sm:px-6 relative z-10 max-w-6xl">
				<div className="text-center mb-16">
					<h2 className="text-4xl sm:text-5xl font-bold text-foreground mb-6">
						{whyTerarushContent.title}
					</h2>
					<p className="text-lg text-muted-foreground max-w-2xl mx-auto">
						{whyTerarushContent.subtitle}
					</p>
				</div>

				<div className="grid md:grid-cols-2 gap-8">
					{pillars.map((pillar, pidx) => {
						const PIcon = pillar.icon;
						return (
							<div
								key={pidx}
								className="pillar-card relative overflow-hidden rounded-3xl border border-border bg-card p-8 md:p-10"
							>
								{/* Gradient overlay */}
								<div className={`absolute inset-0 bg-gradient-to-br ${pillar.gradient} pointer-events-none`} />

								<div className="relative z-10">
									<div className="flex items-center gap-4 mb-8">
										<div className="inline-flex items-center justify-center w-12 h-12 rounded-2xl bg-primary/10">
											<PIcon className="h-6 w-6 text-primary" />
										</div>
										<h3 className="text-2xl font-bold text-foreground">
											{pillar.title}
										</h3>
									</div>

									<div className="space-y-8">
										{pillar.items.map((item, iidx) => {
											const ItemIcon = iconMap[item.icon as keyof typeof iconMap];
											return (
												<div key={iidx} className="flex gap-5">
													<div className="flex-shrink-0 mt-1">
														<div className="flex items-center justify-center w-10 h-10 rounded-xl bg-primary/10">
															<ItemIcon className="h-5 w-5 text-primary" />
														</div>
													</div>
													<div>
														<h4 className="text-lg font-semibold text-foreground mb-2">
															{item.title}
														</h4>
														<p className="text-muted-foreground leading-relaxed text-sm">
															{item.description}
														</p>
													</div>
												</div>
											);
										})}
									</div>
								</div>
							</div>
						);
					})}
				</div>
			</div>
		</section>
	);
}
