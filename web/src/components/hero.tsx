import { useEffect, useRef } from "react";
import { ArrowRight } from "lucide-react";
import gsap from "gsap";
import { Button } from "@/components/ui/button";
import { heroContent } from "@/content/config";

export default function Hero() {
	const heroRef = useRef<HTMLDivElement>(null);
	const leftRef = useRef<HTMLDivElement>(null);
	const rightRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		const ctx = gsap.context(() => {
			const tl = gsap.timeline({ defaults: { ease: "power3.out" } });
			const els = leftRef.current?.querySelectorAll(".anim-el");
			if (els) {
				tl.fromTo(
					els,
					{ y: 30, opacity: 0 },
					{ y: 0, opacity: 1, duration: 0.7, stagger: 0.12 },
				);
			}
			if (rightRef.current) {
				tl.fromTo(
					rightRef.current,
					{ scale: 0.95, opacity: 0 },
					{ scale: 1, opacity: 1, duration: 0.8 },
					"-=0.4",
				);
			}
		}, heroRef);
		return () => ctx.revert();
	}, []);

	return (
		<div
			ref={heroRef}
			className="relative overflow-hidden bg-background pt-24 pb-16 lg:pt-24 lg:pb-24 min-h-[100dvh] flex items-center"
			id="home"
		>
			{/* Subtle background */}
			<div className="absolute inset-0 overflow-hidden pointer-events-none">
				<div className="absolute top-0 right-0 w-[600px] h-[600px] bg-primary/5 rounded-full blur-3xl" />
				<div className="absolute bottom-0 left-0 w-[400px] h-[400px] bg-primary/5 rounded-full blur-3xl" />
			</div>

			<div className="container mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
				<div className="grid lg:grid-cols-2 gap-12 lg:gap-16 items-center max-w-7xl mx-auto">
					{/* Left copy */}
					<div ref={leftRef} className="space-y-8 max-w-xl">
						<div className="anim-el inline-flex items-center gap-2 px-4 py-1.5 text-sm font-medium rounded-full bg-primary/10 text-primary border border-primary/20">
							{heroContent.badge.text}
						</div>

						<h1 className="anim-el text-4xl sm:text-5xl lg:text-6xl font-bold tracking-tight text-foreground leading-tight">
							{heroContent.title}
						</h1>

						<p className="anim-el text-lg sm:text-xl text-muted-foreground leading-relaxed">
							{heroContent.description}
						</p>

						<div className="anim-el flex flex-col sm:flex-row gap-4">
							<Button
								size="lg"
								className="h-13 px-8 rounded-xl font-bold text-base bg-primary hover:bg-primary/90 text-primary-foreground shadow-lg hover:shadow-xl transition-all duration-300 flex items-center gap-2 group"
							>
								Start Building
								<ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
							</Button>
							<Button
								asChild
								size="lg"
								variant="outline"
								className="h-13 px-8 rounded-xl font-bold text-base border-2 border-border hover:border-primary hover:bg-accent transition-all duration-300"
							>
								<a href="#contact">Contact Us</a>
							</Button>
						</div>
					</div>

					{/* Right visual */}
					<div
						ref={rightRef}
						className="relative h-[400px] lg:h-[560px] rounded-2xl overflow-hidden shadow-2xl"
					>
						<img
							src="/company/hello_world.jpeg"
							alt="Team workspace"
							className="w-full h-full object-cover"
						/>
						<div className="absolute inset-0 bg-gradient-to-t from-background/20 to-transparent" />
					</div>
				</div>
			</div>
		</div>
	);
}
