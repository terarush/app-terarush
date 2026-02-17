import { useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { LogOut, User, Settings, Home } from "lucide-react";
import gsap from "gsap";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useAuth } from "@/contexts/AuthContext";
import { siteConfig } from "@/content/config";

export default function Dashboard() {
	const containerRef = useRef<HTMLDivElement>(null);
	const { user, logout } = useAuth();
	const navigate = useNavigate();

	useEffect(() => {
		const ctx = gsap.context(() => {
			if (containerRef.current) {
				gsap.fromTo(
					containerRef.current.children,
					{ y: 30, opacity: 0 },
					{
						y: 0,
						opacity: 1,
						duration: 0.8,
						stagger: 0.15,
						ease: "power3.out",
					}
				);
			}
		}, containerRef);

		return () => ctx.revert();
	}, []);

	const handleLogout = async () => {
		try {
			await logout();
			navigate("/login");
		} catch (error) {
			console.error("Logout error:", error);
		}
	};

	return (
		<div className="min-h-screen bg-background">
			{/* Header */}
			<header className="border-b border-border bg-card/50 backdrop-blur-sm sticky top-0 z-50">
				<div className="container mx-auto px-4 py-4 flex items-center justify-between">
					<div className="flex items-center space-x-3">
						<img
							src="/assets/logo.png"
							alt={siteConfig.name}
							className="h-10 w-10 rounded-lg"
						/>
						<span className="text-2xl font-bold text-foreground">
							{siteConfig.name}
						</span>
					</div>
					
					<div className="flex items-center gap-4">
						<Button
							variant="ghost"
							onClick={() => navigate("/")}
							className="text-muted-foreground hover:text-foreground"
						>
							<Home className="h-5 w-5 mr-2" />
							Home
						</Button>
						<Button
							variant="destructive"
							onClick={handleLogout}
							className="rounded-lg"
						>
							<LogOut className="h-5 w-5 mr-2" />
							Logout
						</Button>
					</div>
				</div>
			</header>

			{/* Main Content */}
			<main className="container mx-auto px-4 py-12">
				<div ref={containerRef} className="max-w-4xl mx-auto space-y-8">
					{/* Welcome Section */}
					<div className="opacity-0">
						<h1 className="text-4xl font-bold text-foreground mb-2">
							Welcome back, {user?.name}!
						</h1>
						<p className="text-muted-foreground">
							You're logged in to your dashboard.
						</p>
					</div>

					{/* User Info Card */}
					<Card className="opacity-0 border-border">
						<CardHeader>
							<CardTitle className="flex items-center gap-2">
								<User className="h-5 w-5" />
								Your Profile
							</CardTitle>
						</CardHeader>
						<CardContent className="space-y-4">
							<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
								<div>
									<p className="text-sm text-muted-foreground mb-1">Name</p>
									<p className="text-lg font-semibold text-foreground">{user?.name}</p>
								</div>
								<div>
									<p className="text-sm text-muted-foreground mb-1">Email</p>
									<p className="text-lg font-semibold text-foreground">{user?.email}</p>
								</div>
								<div>
									<p className="text-sm text-muted-foreground mb-1">Role</p>
									<p className="text-lg font-semibold text-foreground capitalize">{user?.role}</p>
								</div>
								<div>
									<p className="text-sm text-muted-foreground mb-1">Member Since</p>
									<p className="text-lg font-semibold text-foreground">
										{user?.created_at && new Date(user.created_at).toLocaleDateString()}
									</p>
								</div>
							</div>
						</CardContent>
					</Card>

					{/* Quick Actions */}
					<Card className="opacity-0 border-border">
						<CardHeader>
							<CardTitle className="flex items-center gap-2">
								<Settings className="h-5 w-5" />
								Quick Actions
							</CardTitle>
						</CardHeader>
						<CardContent>
							<div className="grid grid-cols-1 md:grid-cols-3 gap-4">
								<Button
									variant="outline"
									className="h-auto p-6 flex flex-col items-center gap-2 border-border hover:bg-muted"
								>
									<User className="h-8 w-8 text-primary" />
									<span className="font-semibold">Edit Profile</span>
									<span className="text-xs text-muted-foreground">Update your information</span>
								</Button>
								
								<Button
									variant="outline"
									className="h-auto p-6 flex flex-col items-center gap-2 border-border hover:bg-muted"
								>
									<Settings className="h-8 w-8 text-primary" />
									<span className="font-semibold">Settings</span>
									<span className="text-xs text-muted-foreground">Manage preferences</span>
								</Button>
								
								<Button
									variant="outline"
									onClick={() => navigate("/")}
									className="h-auto p-6 flex flex-col items-center gap-2 border-border hover:bg-muted"
								>
									<Home className="h-8 w-8 text-primary" />
									<span className="font-semibold">Go Home</span>
									<span className="text-xs text-muted-foreground">Back to landing page</span>
								</Button>
							</div>
						</CardContent>
					</Card>

					{/* Admin Badge (if admin) */}
					{user?.role === "admin" && (
						<Card className="opacity-0 border-primary bg-primary/5">
							<CardContent className="p-6">
								<div className="flex items-center gap-3">
									<div className="w-12 h-12 rounded-full bg-primary/20 flex items-center justify-center">
										<Settings className="h-6 w-6 text-primary" />
									</div>
									<div>
										<h3 className="font-semibold text-foreground">Admin Access</h3>
										<p className="text-sm text-muted-foreground">
											You have administrator privileges
										</p>
									</div>
								</div>
							</CardContent>
						</Card>
					)}
				</div>
			</main>
		</div>
	);
}
