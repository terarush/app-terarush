import { useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import {
	User,
	Settings,
	TrendingUp,
	Package,
	FileText,
	Activity,
} from "lucide-react";
import gsap from "gsap";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useUser, useIsAdmin } from "@/hooks";

export default function Dashboard() {
	const containerRef = useRef<HTMLDivElement>(null);
	const user = useUser();
	const isAdmin = useIsAdmin();
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
					},
				);
			}
		}, containerRef);

		return () => ctx.revert();
	}, []);

	return (
		<div className="flex-1 overflow-y-auto">
			{/* Page Header */}
			<div className="border-b border-border bg-card/50 backdrop-blur-sm sticky top-0 z-10">
				<div className="px-8 py-6">
					<h1 className="text-3xl font-bold text-foreground">
						Dashboard
					</h1>
					<p className="text-muted-foreground mt-1">
						Welcome back, {user?.name}! Here's what's happening
						today.
					</p>
				</div>
			</div>

			{/* Main Content */}
			<div className="p-8">
				<div ref={containerRef} className="max-w-7xl mx-auto space-y-8">
					{/* Stats Grid */}
					<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 opacity-0">
						<Card className="border-border hover:shadow-lg transition-shadow">
							<CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
								<CardTitle className="text-sm font-medium">
									Total Projects
								</CardTitle>
								<Package className="h-4 w-4 text-muted-foreground" />
							</CardHeader>
							<CardContent>
								<div className="text-2xl font-bold">24</div>
								<p className="text-xs text-muted-foreground mt-1">
									<span className="text-green-600">+2</span>{" "}
									from last month
								</p>
							</CardContent>
						</Card>

						<Card className="border-border hover:shadow-lg transition-shadow">
							<CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
								<CardTitle className="text-sm font-medium">
									Active Tasks
								</CardTitle>
								<Activity className="h-4 w-4 text-muted-foreground" />
							</CardHeader>
							<CardContent>
								<div className="text-2xl font-bold">142</div>
								<p className="text-xs text-muted-foreground mt-1">
									<span className="text-green-600">+12</span>{" "}
									this week
								</p>
							</CardContent>
						</Card>

						<Card className="border-border hover:shadow-lg transition-shadow">
							<CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
								<CardTitle className="text-sm font-medium">
									Documents
								</CardTitle>
								<FileText className="h-4 w-4 text-muted-foreground" />
							</CardHeader>
							<CardContent>
								<div className="text-2xl font-bold">89</div>
								<p className="text-xs text-muted-foreground mt-1">
									<span className="text-blue-600">+5</span>{" "}
									added today
								</p>
							</CardContent>
						</Card>

						<Card className="border-border hover:shadow-lg transition-shadow">
							<CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
								<CardTitle className="text-sm font-medium">
									Growth Rate
								</CardTitle>
								<TrendingUp className="h-4 w-4 text-muted-foreground" />
							</CardHeader>
							<CardContent>
								<div className="text-2xl font-bold">+12.5%</div>
								<p className="text-xs text-muted-foreground mt-1">
									<span className="text-green-600">
										+2.1%
									</span>{" "}
									from last month
								</p>
							</CardContent>
						</Card>
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
									<p className="text-sm text-muted-foreground mb-1">
										Name
									</p>
									<p className="text-lg font-semibold text-foreground">
										{user?.name}
									</p>
								</div>
								<div>
									<p className="text-sm text-muted-foreground mb-1">
										Email
									</p>
									<p className="text-lg font-semibold text-foreground">
										{user?.email}
									</p>
								</div>
								<div>
									<p className="text-sm text-muted-foreground mb-1">
										Role
									</p>
									<p className="text-lg font-semibold text-foreground capitalize">
										{user?.role}
									</p>
								</div>
								<div>
									<p className="text-sm text-muted-foreground mb-1">
										Member Since
									</p>
									<p className="text-lg font-semibold text-foreground">
										{user?.created_at &&
											new Date(
												user.created_at,
											).toLocaleDateString()}
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
									onClick={() =>
										navigate("/dashboard/profile")
									}
									className="h-auto p-6 flex flex-col items-center gap-2 border-border hover:bg-muted"
								>
									<User className="h-8 w-8 text-primary" />
									<span className="font-semibold">
										Edit Profile
									</span>
									<span className="text-xs text-muted-foreground">
										Update your information
									</span>
								</Button>

								<Button
									variant="outline"
									onClick={() =>
										navigate("/dashboard/settings")
									}
									className="h-auto p-6 flex flex-col items-center gap-2 border-border hover:bg-muted"
								>
									<Settings className="h-8 w-8 text-primary" />
									<span className="font-semibold">
										Settings
									</span>
									<span className="text-xs text-muted-foreground">
										Manage preferences
									</span>
								</Button>

								<Button
									variant="outline"
									onClick={() =>
										navigate("/dashboard/projects")
									}
									className="h-auto p-6 flex flex-col items-center gap-2 border-border hover:bg-muted"
								>
									<Package className="h-8 w-8 text-primary" />
									<span className="font-semibold">
										Projects
									</span>
									<span className="text-xs text-muted-foreground">
										View all projects
									</span>
								</Button>
							</div>
						</CardContent>
					</Card>

					{/* Admin Badge (if admin) */}
					{isAdmin && (
						<Card className="opacity-0 border-primary bg-primary/5">
							<CardContent className="p-6">
								<div className="flex items-center gap-3">
									<div className="w-12 h-12 rounded-full bg-primary/20 flex items-center justify-center">
										<Settings className="h-6 w-6 text-primary" />
									</div>
									<div>
										<h3 className="font-semibold text-foreground">
											Admin Access
										</h3>
										<p className="text-sm text-muted-foreground">
											You have administrator privileges
										</p>
									</div>
								</div>
							</CardContent>
						</Card>
					)}
				</div>
			</div>
		</div>
	);
}
