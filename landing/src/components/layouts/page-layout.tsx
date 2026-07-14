import { type ReactNode } from "react";

interface PageLayoutProps {
	children: ReactNode;
}

export function PageLayout({ children }: PageLayoutProps) {
	return <div className="min-h-screen bg-background flex flex-col">{children}</div>;
}

interface PageContentProps {
	children: ReactNode;
	className?: string;
}

export function PageContent({ children, className = "" }: PageContentProps) {
	return (
		<main className={`flex-1 pt-24 ${className}`}>
			{children}
		</main>
	);
}

interface HeroSectionProps {
	children: ReactNode;
}

export function HeroSection({ children }: HeroSectionProps) {
	return (
		<div className="relative overflow-hidden border-b border-border">
			<div className="container mx-auto px-4 py-16 md:py-24">
				{children}
			</div>
		</div>
	);
}

interface ContentSectionProps {
	children: ReactNode;
	className?: string;
}

export function ContentSection({ children, className = "" }: ContentSectionProps) {
	return (
		<div className={`container mx-auto px-4 py-16 ${className}`}>
			{children}
		</div>
	);
}
