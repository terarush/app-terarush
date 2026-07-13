import { ArrowRight } from "lucide-react";
import { motion, useReducedMotion } from "motion/react";

import { Button } from "@/components/ui/button";
import { heroContent } from "@/content/config";

export default function Hero() {
	const reduce = useReducedMotion();

	return (
		<div
			className="relative bg-background pt-32 pb-16 lg:pt-40 lg:pb-24 min-h-[100dvh] flex items-center justify-center overflow-hidden"
			id="home"
		>
			{/* Subtle animated mesh blob - not AI purple, uses primary */}
			{!reduce && (
				<div className="absolute inset-0 overflow-hidden pointer-events-none">
					<motion.div
						className="absolute -top-1/4 -right-1/4 w-[600px] h-[600px] rounded-full"
						style={{
							background:
								"radial-gradient(circle, hsl(var(--primary) / 0.08), transparent 60%)",
						}}
						animate={{
							x: [0, -30, 20, -10, 0],
							y: [0, 20, -30, 10, 0],
						}}
						transition={{
							duration: 20,
							ease: "linear",
							repeat: Infinity,
						}}
					/>
					<motion.div
						className="absolute -bottom-1/4 -left-1/4 w-[500px] h-[500px] rounded-full"
						style={{
							background:
								"radial-gradient(circle, hsl(var(--primary) / 0.05), transparent 60%)",
						}}
						animate={{
							x: [0, 25, -15, 10, 0],
							y: [0, -15, 25, -10, 0],
						}}
						transition={{
							duration: 25,
							ease: "linear",
							repeat: Infinity,
						}}
					/>
				</div>
			)}

			{/* Fine grain overlay */}
			<div
				className="absolute inset-0 pointer-events-none opacity-[0.015]"
				style={{
					backgroundImage:
						"url(\"data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E\")",
				}}
			/>

			<div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-7xl relative z-10">
				<div className="grid lg:grid-cols-12 gap-6 items-stretch">
					{/* Left Bento Cell: Text Content */}
					<motion.div
						initial={reduce ? false : { opacity: 0, y: 20 }}
						animate={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.5, ease: "easeOut" }}
						whileHover={reduce ? {} : { y: -2 }}
						className="lg:col-span-7 flex flex-col justify-center p-8 md:p-12 rounded-3xl bg-secondary/30 border border-border/50 transition-shadow duration-500"
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
									className="h-12 px-8 rounded-full font-medium text-base bg-foreground text-background hover:bg-foreground/90 transition-all active:scale-[0.97] inline-flex items-center justify-center"
								>
									Start Building
								</Button>
								<Button
									size="lg"
									variant="outline"
									className="h-12 px-8 rounded-full font-medium text-base border-border hover:bg-accent transition-all active:scale-[0.97]"
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

					{/* Right Bento Cell: Image - no hover animation */}
					<div className="lg:col-span-5 relative min-h-[400px] lg:min-h-full rounded-3xl overflow-hidden border border-border/50">
						<img
							src="/company/hello_world.jpeg"
							alt="Team workspace"
							className="absolute inset-0 w-full h-full object-cover"
						/>
						<div className="absolute inset-0 bg-gradient-to-t from-background/20 to-transparent pointer-events-none" />
					</div>
				</div>
			</div>
		</div>
	);
}
