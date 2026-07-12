import { useEffect, useRef } from "react";
import { Users, Target, Zap, Heart, Linkedin, Github } from "lucide-react";
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { aboutContent } from "@/content/config";

gsap.registerPlugin(ScrollTrigger);

const iconMap = {
	users: Users,
	target: Target,
	zap: Zap,
	heart: Heart,
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
			const fadeIn = (el: HTMLElement, delay = 0) =>
				gsap.fromTo(
					el,
					{ y: 30, opacity: 0 },
					{
						y: 0,
						opacity: 1,
						duration: 0.7,
						delay,
						ease: "power3.out",
						scrollTrigger: {
							trigger: el,
							start: "top 80%",
							toggleActions: "play none none none",
						},
					},
				);

			if (headerRef.current) fadeIn(headerRef.current);
			if (storyRef.current) fadeIn(storyRef.current);
			if (valuesRef.current) fadeIn(valuesRef.current);
			if (teamRef.current) fadeIn(teamRef.current);
			if (ctaRef.current) fadeIn(ctaRef.current);
		}, sectionRef);
		return () => ctx.revert();
	}, []);

	return (
		<section
			ref={sectionRef}
			className="relative overflow-hidden bg-background py-24"
			id="about"
		>
			<div className="absolute inset-0 overflow-hidden pointer-events-none">
				<div className="absolute top-1/4 -right-40 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
				<div className="absolute bottom-1/4 -left-40 w-96 h-96 bg-accent/10 rounded-full blur-3xl" />
			</div>

			<div className="container mx-auto px-4 sm:px-6 relative z-10">
				{/* Header */}
				<div ref={headerRef} className="text-center mb-20">
					<h2 className="text-4xl sm:text-5xl md:text-6xl font-bold tracking-tight mb-6">
						<span className="block text-foreground">
							{aboutContent.headline.line1}
						</span>
						<span className="text-primary">
							{aboutContent.headline.line2}
						</span>
					</h2>
					<p className="text-lg sm:text-xl text-muted-foreground max-w-3xl mx-auto">
						{aboutContent.description}
					</p>
				</div>

				{/* Our Story */}
				<div
					ref={storyRef}
					className="grid lg:grid-cols-2 gap-16 items-center mb-24"
				>
					<div>
						<h3 className="text-3xl sm:text-4xl font-bold text-foreground mb-6">
							{aboutContent.story.title}
						</h3>
						<div className="space-y-6">
							{aboutContent.story.paragraphs.map(
								(paragraph, index) => (
									<p
										key={index}
										className="text-lg text-muted-foreground leading-relaxed"
									>
										{paragraph}
									</p>
								),
							)}
							<div className="flex items-start space-x-4 p-6 bg-card rounded-2xl border border-border">
								<div>
									<p className="text-foreground font-medium italic leading-relaxed">
										"{aboutContent.story.quote.text}"
									</p>
									<p className="text-sm text-muted-foreground mt-2 font-semibold">
										{aboutContent.story.quote.author}
									</p>
								</div>
							</div>
						</div>
					</div>
					<div className="relative">
						<div className="aspect-[4/3] rounded-2xl overflow-hidden shadow-2xl">
							<img
								src={aboutContent.story.image.src}
								alt={aboutContent.story.image.alt}
								className="w-full h-full object-cover"
							/>
						</div>
					</div>
				</div>

				{/* Values Section */}
				<div ref={valuesRef} className="mb-24">
					<div className="text-center mb-16">
						<h3 className="text-3xl sm:text-4xl font-bold text-foreground mb-4">
							{aboutContent.values.title}
						</h3>
						<p className="text-lg text-muted-foreground max-w-2xl mx-auto">
							{aboutContent.values.subtitle}
						</p>
					</div>
					<div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
						{aboutContent.values.items.map((value, index) => {
							const Icon =
								iconMap[value.icon as keyof typeof iconMap];
							return (
								<Card
									key={index}
									className="group h-full border-0 bg-card/80 backdrop-blur-sm hover:shadow-xl hover:scale-[1.02] transition-all duration-300"
								>
									<CardContent className="p-6 text-center">
										<div
											className={`inline-flex items-center justify-center w-16 h-16 rounded-2xl ${value.color} mb-4 group-hover:scale-110 transition-transform duration-300`}
										>
											<Icon className="h-8 w-8 text-primary-foreground" />
										</div>
										<h4 className="text-xl font-semibold text-foreground mb-3">
											{value.title}
										</h4>
										<p className="text-muted-foreground">
											{value.description}
										</p>
									</CardContent>
								</Card>
							);
						})}
					</div>
				</div>

				{/* Team */}
				<div ref={teamRef} className="mb-24">
					<div className="text-center mb-16">
						<h3 className="text-3xl sm:text-4xl font-bold text-foreground mb-4">
							{aboutContent.team.title}
						</h3>
						<p className="text-lg text-muted-foreground max-w-2xl mx-auto">
							{aboutContent.team.subtitle}
						</p>
					</div>
					<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 max-w-5xl mx-auto">
						{aboutContent.team.members.map((member, index) => (
							<Card
								key={index}
								className="group relative overflow-hidden border border-border bg-gradient-to-br from-card to-card/50 hover:shadow-2xl hover:scale-[1.02] transition-all duration-500"
							>
								<CardContent className="p-8 relative z-10">
									<div className="flex flex-col items-center text-center">
										<div className="relative mb-6">
											<div className="absolute inset-0 bg-gradient-to-br from-primary to-primary/50 rounded-full blur-xl opacity-0 group-hover:opacity-70 transition-opacity duration-500" />
											<Avatar className="w-28 h-28 relative border-4 border-background shadow-xl ring-4 ring-primary/20 group-hover:ring-primary/50 transition-all duration-500">
												<AvatarImage
													src={
														member.avatar ||
														"/company/logo.png"
													}
													className="object-cover"
												/>
												<AvatarFallback className="bg-gradient-to-br from-primary to-primary/80 text-primary-foreground text-2xl font-bold">
													{member.initials}
												</AvatarFallback>
											</Avatar>
										</div>

										<div className="space-y-2 mb-6">
											<h4 className="text-xl font-bold text-foreground group-hover:text-primary transition-colors duration-300">
												{member.name}
											</h4>
											<p className="text-primary font-semibold text-sm uppercase tracking-wider">
												{member.role}
											</p>
											<p className="text-sm text-muted-foreground leading-relaxed pt-2">
												{member.bio}
											</p>
										</div>

										<div className="flex gap-3 pt-4 border-t border-border w-full justify-center">
											<Button
												variant="ghost"
												size="sm"
												className="h-10 w-10 rounded-full hover:bg-primary/10 hover:text-primary hover:scale-110 transition-all duration-300"
												asChild
											>
												<a
													href={
														member.social.linkedin
													}
													target="_blank"
													rel="noopener noreferrer"
													aria-label={`${member.name} LinkedIn`}
												>
													<Linkedin className="h-5 w-5" />
												</a>
											</Button>
											<Button
												variant="ghost"
												size="sm"
												className="h-10 w-10 rounded-full hover:bg-primary/10 hover:text-primary hover:scale-110 transition-all duration-300"
												asChild
											>
												<a
													href={member.social.github}
													target="_blank"
													rel="noopener noreferrer"
													aria-label={`${member.name} GitHub`}
												>
													<Github className="h-5 w-5" />
												</a>
											</Button>
										</div>
									</div>
								</CardContent>
							</Card>
						))}
					</div>
				</div>

				{/* CTA */}
				<div
					ref={ctaRef}
					className="text-center bg-primary rounded-3xl p-12 text-primary-foreground shadow-2xl"
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
