import { ArrowRight } from "lucide-react";
import { motion, useReducedMotion } from "motion/react";

import { Button } from "@/components/ui/button";
import { heroContent } from "@/content/config";

export default function Hero() {
	const reduce = useReducedMotion();

	return (
		<div
			className="relative bg-background pt-24 pb-8 lg:pt-32 lg:pb-16 min-h-[100dwh] flex items-center"
			id="home"
		>
			<div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-7xl">
				<div className="grid lg:grid-cols-12 gap-4 md:gap-6 items-stretch">
					{/* Left Bento Cell: Text Content */}
					<motion.div
						initial={reduce ? false : { opacity: 0, y: 16 }}
						animate={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.6, ease: [0.16, 1, 0.3, 1] }}
						className="lg:col-span-6 flex flex-col justify-center p-8 md:p-10 lg:p-12 rounded-3xl bg-secondary/40 border border-border/60 shadow-sm"
					>
						<div className="space-y-6 max-w-xl">
							<div className="inline-flex items-center gap-2 px-3 py-1 text-xs font-semibold uppercase tracking-widest rounded-full bg-primary/10 text-primary border border-primary/20">
								{heroContent.badge.text}
							</div>

							<h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold tracking-tighter text-foreground leading-[1.08]">
								{heroContent.title}
							</h1>

							<p className="text-base sm:text-lg text-muted-foreground leading-relaxed max-w-lg">
								{heroContent.description}
							</p>

							<div className="flex flex-col sm:flex-row gap-3 pt-2">
								<Button
									size="lg"
									className="h-11 px-7 rounded-full font-medium text-sm bg-foreground text-background hover:bg-foreground/90 transition-all active:scale-[0.97] inline-flex items-center justify-center"
								>
									Start Building
									<ArrowRight className="w-4 h-4 ml-1.5" />
								</Button>
								<Button
									size="lg"
									variant="outline"
									className="h-11 px-7 rounded-full font-medium text-sm border-border hover:bg-accent hover:border-foreground/20 transition-all active:scale-[0.97]"
									onClick={() => {
										const el = document.getElementById("contact");
										if (el) el.scrollIntoView({ behavior: "smooth" });
									}}
								>
									Contact Us
								</Button>
							</div>
						</div>
					</motion.div>

					{/* Right Bento Cell: Image */}
					<motion.div
						initial={reduce ? false : { opacity: 0, y: 16 }}
						animate={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.6, delay: 0.1, ease: [0.16, 1, 0.3, 1] }}
						className="lg:col-span-6 relative min-h-[400px] lg:min-h-full rounded-3xl overflow-hidden border border-border/60 shadow-sm"
					>
						<img
							src="/company/hello_world.jpeg"
							alt="Team workspace"
							className="absolute inset-0 w-full h-full object-cover"
						/>
						<div className="absolute inset-0 shadow-[inset_0_2px_8px_rgba(0,0,0,0.06)] pointer-events-none rounded-3xl" />
						<div className="absolute inset-x-0 bottom-0 h-1/3 bg-gradient-to-t from-background/10 to-transparent pointer-events-none" />
					</motion.div>
				</div>
			</div>
		</div>
	);
}
