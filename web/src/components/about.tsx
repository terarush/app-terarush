import { useEffect, useRef } from "react";
import {
	Users,
	Target,
	Zap,
	Heart,
	Award,
	Quote,
	Linkedin,
	Github,
} from "lucide-react";
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";

import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent } from "@/components/ui/card";
import { aboutContent } from "@/content/config";

gsap.registerPlugin(ScrollTrigger);

const iconMap = {
	users: Users,
	target: Target,
	zap: Zap,
	heart: Heart,
	award: Award,
};

export default function About() {
	const sectionRef = useRef<HTMLElement>(null);
	const headerRef = useRef<HTMLDivElement>(null);
	const storyRef = useRef<HTMLDivElement>(null);
	const valuesRef = useRef<HTMLDivElement>(null);
	const teamRef = useRef<HTMLDivElement>(null);
	const ctaRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		const ctx = gsap.context(() => {
			// Header animation
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
					},
				);
			}

			// Story section - ScrollTrigger
			if (storyRef.current) {
				const storyElements =
					storyRef.current.querySelectorAll(".story-element");
				gsap.fromTo(
					storyElements,
					{ y: 50, opacity: 0 },
					{
						y: 0,
						opacity: 1,
						duration: 0.8,
						stagger: 0.2,
						ease: "power3.out",
						scrollTrigger: {
							trigger: storyRef.current,
							start: "top 80%",
							toggleActions: "play none none none",
						},
					},
				);
			}

			// Values - ScrollTrigger
			if (valuesRef.current) {
				const valueCards =
					valuesRef.current.querySelectorAll(".value-card");
				gsap.fromTo(
					valueCards,
					{ y: 40, opacity: 0, scale: 0.95 },
					{
						y: 0,
						opacity: 1,
						scale: 1,
						duration: 0.7,
						stagger: 0.12,
						ease: "back.out(1.2)",
						scrollTrigger: {
							trigger: valuesRef.current,
							start: "top 80%",
							toggleActions: "play none none none",
						},
					},
				);
			}

			// Team - ScrollTrigger
			if (teamRef.current) {
				const teamCards =
					teamRef.current.querySelectorAll(".team-card");
				gsap.fromTo(
					teamCards,
					{ y: 50, opacity: 0, scale: 0.9 },
					{
						y: 0,
						opacity: 1,
						scale: 1,
						duration: 0.8,
						stagger: 0.15,
						ease: "back.out(1.4)",
						scrollTrigger: {
							trigger: teamRef.current,
							start: "top 80%",
							toggleActions: "play none none none",
						},
					},
				);
			}

			// CTA - ScrollTrigger
			if (ctaRef.current) {
				gsap.fromTo(
					ctaRef.current,
					{ y: 40, opacity: 0, scale: 0.95 },
					{
						y: 0,
						opacity: 1,
						scale: 1,
						duration: 0.8,
						ease: "power3.out",
						scrollTrigger: {
							trigger: ctaRef.current,
							start: "top 85%",
							toggleActions: "play none none none",
						},
					},
				);
			}
		}, sectionRef);

		return () => ctx.revert();
	}, []);

	return (
		<section
			ref={sectionRef}
			className="relative overflow-hidden bg-background py-24"
			id="about"
		>
			{/* Background gradients */}
			<div className="absolute inset-0 overflow-hidden">
				<div className="absolute top-1/4 -right-40 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
				<div className="absolute bottom-1/4 -left-40 w-96 h-96 bg-accent/10 rounded-full blur-3xl" />
			</div>

			<div className="container mx-auto px-4 sm:px-6 relative z-10">
				{/* Header */}
				<div ref={headerRef} className="text-center mb-20">
					<Badge
						variant="secondary"
						className="px-4 py-1.5 text-sm font-medium rounded-full bg-accent text-accent-foreground mb-6 opacity-0"
					>
						{aboutContent.badge.text}
					</Badge>
					<h2 className="text-4xl sm:text-5xl md:text-6xl font-bold tracking-tight mb-6 opacity-0">
						<span className="block text-foreground">
							{aboutContent.headline.line1}
						</span>
						<span className="text-primary">
							{aboutContent.headline.line2}
						</span>
					</h2>
					<p className="text-lg sm:text-xl text-muted-foreground max-w-3xl mx-auto opacity-0">
						{aboutContent.description}
					</p>
				</div>

				{/* Our Story */}
				<div
					ref={storyRef}
					className="grid lg:grid-cols-2 gap-16 items-center mb-24"
				>
					<div className="story-element opacity-0">
						<h3 className="text-3xl sm:text-4xl font-bold text-foreground mb-6">
							{aboutContent.story.title}
						</h3>
						<div className="space-y-6">
							{aboutContent.story.paragraphs.map(
								(paragraph, index) => (
									<p
										key={index}
										className="text-lg text-muted-foreground"
									>
										{paragraph}
									</p>
								),
							)}
							<div className="flex items-start space-x-4 p-6 bg-card rounded-2xl border border-border">
								<Quote className="h-6 w-6 text-primary mt-1 flex-shrink-0" />
								<div>
									<p className="text-foreground font-medium italic">
										"{aboutContent.story.quote.text}"
									</p>
									<p className="text-sm text-muted-foreground mt-2">
										— {aboutContent.story.quote.author}
									</p>
								</div>
							</div>
						</div>
					</div>
					<div className="story-element opacity-0 relative">
						<div className="aspect-[4/3] rounded-2xl overflow-hidden shadow-2xl">
							<img
								src={aboutContent.story.image.src}
								alt={aboutContent.story.image.alt}
								className="w-full h-full object-cover"
							/>
						</div>
						<div className="absolute -bottom-6 -right-6 bg-card rounded-2xl p-6 shadow-xl border border-border">
							<div className="flex items-center space-x-3">
								<Award className="h-8 w-8 text-primary" />
								<div>
									<p className="font-semibold text-foreground">
										{aboutContent.story.award.title}
									</p>
									<p className="text-sm text-muted-foreground">
										{aboutContent.story.award.subtitle}
									</p>
								</div>
							</div>
						</div>
					</div>
				</div>
				{/* CTA */}
				<div
					ref={ctaRef}
					className="opacity-0 text-center bg-primary rounded-3xl p-12 text-primary-foreground shadow-2xl"
				>
					<h3 className="text-3xl sm:text-4xl font-bold mb-4">
						{aboutContent.cta.title}
					</h3>
					<p className="text-xl mb-5 max-w-2xl mx-auto">
						{aboutContent.cta.description}
					</p>
				</div>
			</div>
		</section>
	);
}
