import { ArrowRight } from "lucide-react";
import { motion, useReducedMotion } from "motion/react";
import { Button } from "@/components/ui/button";
import { heroContent } from "@/content/config";

export default function Hero() {
	const reduce = useReducedMotion();

	return (
		<div
			className="relative bg-background pt-32 pb-16 lg:pt-40 lg:pb-24 min-h-[100dvh] flex items-center justify-center"
			id="home"
		>
			<div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-7xl">
				<div className="grid lg:grid-cols-12 gap-6 items-stretch">
					{/* Left Bento Cell: Text Content */}
					<motion.div
						initial={reduce ? false : { opacity: 0, y: 20 }}
						animate={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.5, ease: "easeOut" }}
						className="lg:col-span-7 flex flex-col justify-center p-8 md:p-12 rounded-3xl bg-secondary/30 border border-border/50"
					>
						<div className="space-y-8 max-w-2xl">
							<div className="inline-flex items-center gap-2 px-3 py-1 text-xs font-semibold uppercase tracking-wider rounded-full bg-primary/10 text-primary border border-primary/20">
								{heroContent.badge.text}
							</div>

							<h1 className="text-4xl md:text-5xl lg:text-6xl font-bold tracking-tighter text-foreground leading-[1.1]">
								{heroContent.title}
							</h1>

							<p className="text-lg md:text-xl text-muted-foreground leading-relaxed">
								{heroContent.description}
							</p>

							<div className="flex flex-col sm:flex-row gap-4 pt-4">
								<Button
									size="lg"
									className="h-12 px-8 rounded-full font-medium text-base bg-foreground text-background hover:bg-foreground/90 transition-all flex items-center gap-2 group"
								>
									Start Building
									<ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
								</Button>
								<Button
									asChild
									size="lg"
									variant="outline"
									className="h-12 px-8 rounded-full font-medium text-base border-border hover:bg-accent transition-all"
								>
									<a href="#contact">Contact Us</a>
								</Button>
							</div>
						</div>
					</motion.div>

					{/* Right Bento Cell: Image */}
					<motion.div
						initial={reduce ? false : { opacity: 0, y: 20 }}
						animate={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.5, delay: 0.1, ease: "easeOut" }}
						className="lg:col-span-5 relative min-h-[400px] lg:min-h-full rounded-3xl overflow-hidden border border-border/50"
					>
						<img
							src="/company/hello_world.jpeg"
							alt="Team workspace"
							className="absolute inset-0 w-full h-full object-cover"
						/>
						{/* Subtle gradient overlay to match bento aesthetic */}
						<div className="absolute inset-0 bg-gradient-to-t from-background/20 to-transparent pointer-events-none" />
					</motion.div>
				</div>
			</div>
		</div>
	);
}