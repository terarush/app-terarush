import type React from "react";
import { useState } from "react";
import {
	Mail,
	MapPin,
	ArrowRight,
	Linkedin,
	Github,
	Twitter,
	ExternalLink,
} from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { siteConfig } from "@/content/config";

const footerLinks = {
	product: [
		{ name: "Services", href: "#services" },
		{ name: "Open Source", href: "#opensource" },
		{ name: "Documentation", href: "/docs", external: true },
	],
	company: [
		{ name: "About", href: "#about" },
		{ name: "Why Us", href: "#why" },
		{ name: "Contact", href: "#contact" },
	],
	resources: [
		{ name: "GitHub", href: siteConfig.links.github, external: true },
		{ name: "Community", href: "/community" },
		{ name: "Blog", href: "/blog" },
	],
	legal: [
		{ name: "Privacy Policy", href: "/privacy" },
		{ name: "Terms of Service", href: "/terms" },
		{ name: "License", href: "/license" },
	],
};

const socialLinks = [
	{
		name: "Twitter",
		icon: Twitter,
		href: siteConfig.links.twitter,
		color: "hover:text-foreground",
	},
	{
		name: "LinkedIn",
		icon: Linkedin,
		href: siteConfig.links.linkedin,
		color: "hover:text-blue-500",
	},
	{
		name: "GitHub",
		icon: Github,
		href: siteConfig.links.github,
		color: "hover:text-foreground",
	},
];

export default function Footer() {
	const [isSubscribed, setIsSubscribed] = useState(false);
	const [subscribeError, setSubscribeError] = useState("");
	const [email, setEmail] = useState("");

	const handleNewsletterSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setSubscribeError("");

		// Simple email validation
		if (!email || !email.includes("@")) {
			setSubscribeError("Please enter a valid email");
			return;
		}

		try {
			// Simulate API call (replace with actual endpoint)
			await new Promise((resolve) => setTimeout(resolve, 1000));

			setEmail("");
			setIsSubscribed(true);
			setTimeout(() => setIsSubscribed(false), 3000);
		} catch (error) {
			setIsSubscribed(false);
			setSubscribeError("Failed to subscribe. Please try again.");
		}
	};

	return (
		<footer className="relative bg-background border-t border-border">
			{/* Background decorations */}
			<div className="absolute inset-0 overflow-hidden opacity-30">
				<div className="absolute top-0 left-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
				<div className="absolute bottom-0 right-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
			</div>

			<div className="container mx-auto px-4 sm:px-6 relative z-10 max-w-7xl">
				{/* Main Footer Content */}
				<div className="py-16">
					<div className="grid grid-cols-1 lg:grid-cols-12 gap-12">
						{/* Left Column - Brand & Newsletter */}
						<div className="lg:col-span-4">
							<div className="mb-8">
								<a
									href="#home"
									className="flex items-center space-x-3 mb-6 group"
								>
									<img
										src="/assets/logo.png"
										alt="Terarush"
										className="h-10 w-10 rounded-xl shadow-lg group-hover:scale-110 transition-transform duration-300"
									/>
									<span className="text-2xl font-bold text-foreground">
										{siteConfig.name}
									</span>
								</a>
								<p className="text-muted-foreground text-lg leading-relaxed mb-6">
									{siteConfig.description}
								</p>
							</div>

							{/* Newsletter */}
							<div className="mb-8">
								<h3 className="text-lg font-semibold text-foreground mb-4">
									Stay Updated
								</h3>
								<p className="text-muted-foreground mb-4">
									Get the latest updates on our open-source
									projects and tools.
								</p>
								<form
									onSubmit={handleNewsletterSubmit}
									className="space-y-3"
								>
									<div className="flex flex-col sm:flex-row gap-3">
										<Input
											type="email"
											placeholder="Enter your email"
											value={email}
											onChange={(e) =>
												setEmail(e.target.value)
											}
											className="flex-1 rounded-xl border-border focus:border-primary"
											required
										/>
										<Button
											type="submit"
											className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl px-6 font-medium shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-300"
										>
											Subscribe
											<ArrowRight className="ml-2 h-4 w-4" />
										</Button>
									</div>
									{subscribeError && (
										<p className="text-red-600 dark:text-red-400 text-sm font-medium">
											{subscribeError}
										</p>
									)}
									{isSubscribed && (
										<p className="text-green-600 dark:text-green-400 text-sm font-medium">
											✓ Successfully subscribed!
										</p>
									)}
								</form>
							</div>

							{/* Contact Info */}
							<div className="space-y-3">
								<div className="flex items-center space-x-3 text-muted-foreground hover:text-foreground transition-colors">
									<Mail className="h-5 w-5" />
									<a href={`mailto:${siteConfig.email}`}>
										{siteConfig.email}
									</a>
								</div>
								<div className="flex items-center space-x-3 text-muted-foreground">
									<MapPin className="h-5 w-5" />
									<span>Malang, Indonesia</span>
								</div>
							</div>
						</div>

						{/* Right Columns - Links */}
						<div className="lg:col-span-8">
							<div className="grid grid-cols-2 md:grid-cols-4 gap-8">
								{/* Product */}
								<div>
									<h3 className="text-lg font-semibold text-foreground mb-6">
										Product
									</h3>
									<ul className="space-y-3">
										{footerLinks.product.map((link) => (
											<li key={link.name}>
												{link.external ? (
													<a
														href={link.href}
														target="_blank"
														rel="noopener noreferrer"
														className="text-muted-foreground hover:text-foreground transition-colors duration-200 flex items-center group"
													>
														{link.name}
														<ExternalLink className="ml-1 h-3 w-3" />
													</a>
												) : (
													<a
														href={link.href}
														className="text-muted-foreground hover:text-foreground transition-colors duration-200 flex items-center group"
													>
														{link.name}
														<ArrowRight className="ml-1 h-3 w-3 opacity-0 group-hover:opacity-100 transition-opacity duration-200" />
													</a>
												)}
											</li>
										))}
									</ul>
								</div>

								{/* Company */}
								<div>
									<h3 className="text-lg font-semibold text-foreground mb-6">
										Company
									</h3>
									<ul className="space-y-3">
										{footerLinks.company.map((link) => (
											<li key={link.name}>
												<a
													href={link.href}
													className="text-muted-foreground hover:text-foreground transition-colors duration-200 flex items-center group"
												>
													{link.name}
													<ArrowRight className="ml-1 h-3 w-3 opacity-0 group-hover:opacity-100 transition-opacity duration-200" />
												</a>
											</li>
										))}
									</ul>
								</div>

								{/* Resources */}
								<div>
									<h3 className="text-lg font-semibold text-foreground mb-6">
										Resources
									</h3>
									<ul className="space-y-3">
										{footerLinks.resources.map((link) => (
											<li key={link.name}>
												{link.external ? (
													<a
														href={link.href}
														target="_blank"
														rel="noopener noreferrer"
														className="text-muted-foreground hover:text-foreground transition-colors duration-200 flex items-center group"
													>
														{link.name}
														<ExternalLink className="ml-1 h-3 w-3" />
													</a>
												) : (
													<a
														href={link.href}
														className="text-muted-foreground hover:text-foreground transition-colors duration-200 flex items-center group"
													>
														{link.name}
														<ArrowRight className="ml-1 h-3 w-3 opacity-0 group-hover:opacity-100 transition-opacity duration-200" />
													</a>
												)}
											</li>
										))}
									</ul>
								</div>

								{/* Legal */}
								<div>
									<h3 className="text-lg font-semibold text-foreground mb-6">
										Legal
									</h3>
									<ul className="space-y-3">
										{footerLinks.legal.map((link) => (
											<li key={link.name}>
												<a
													href={link.href}
													className="text-muted-foreground hover:text-foreground transition-colors duration-200 flex items-center group"
												>
													{link.name}
													<ArrowRight className="ml-1 h-3 w-3 opacity-0 group-hover:opacity-100 transition-opacity duration-200" />
												</a>
											</li>
										))}
									</ul>
								</div>
							</div>
						</div>
					</div>
				</div>

				{/* Bottom Bar */}
				<div className="border-t border-border py-8">
					<div className="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
						<div className="flex items-center space-x-2 text-muted-foreground">
							<span>
								© {new Date().getFullYear()} {siteConfig.name}
								. All rights reserved.
							</span>
						</div>

						<div className="flex items-center space-x-4">
							<span className="text-sm text-muted-foreground mr-2">
								Follow us:
							</span>
							{socialLinks.map((social) => {
								const Icon = social.icon;
								return (
									<a
										key={social.name}
										href={social.href}
										target="_blank"
										rel="noopener noreferrer"
										className={`p-2 rounded-full bg-card border border-border shadow-sm hover:shadow-md hover:border-primary/50 transition-all duration-300 text-muted-foreground ${social.color} hover:scale-110`}
										aria-label={social.name}
									>
										<Icon className="h-5 w-5" />
									</a>
								);
							})}
						</div>
					</div>
				</div>
			</div>
		</footer>
	);
}
