import { useUser, useIsAdmin } from "@/hooks";
import { BarChart3, Eye, FileText, Users, TrendingUp, Calendar, ArrowUpRight, ArrowDownRight } from "lucide-react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function Dashboard() {
	const user = useUser();
	const isAdmin = useIsAdmin();

	// Mock data - replace with actual API calls
	const stats = [
		{
			label: "Total Posts",
			value: "24",
			change: "+12%",
			isPositive: true,
			icon: FileText,
			color: "text-blue-500",
			bgColor: "bg-blue-500/10"
		},
		{
			label: "Total Views",
			value: "8,245",
			change: "+23%",
			isPositive: true,
			icon: Eye,
			color: "text-emerald-500",
			bgColor: "bg-emerald-500/10"
		},
		{
			label: "Active Users",
			value: "1,234",
			change: "+5%",
			isPositive: true,
			icon: Users,
			color: "text-purple-500",
			bgColor: "bg-purple-500/10"
		},
		{
			label: "Engagement Rate",
			value: "68%",
			change: "-2%",
			isPositive: false,
			icon: TrendingUp,
			color: "text-amber-500",
			bgColor: "bg-amber-500/10"
		}
	];

	const recentPosts = [
		{ id: 1, title: "Getting Started with React 18", views: 1203, date: "2 days ago", status: "Published" },
		{ id: 2, title: "Advanced TypeScript Patterns", views: 856, date: "5 days ago", status: "Published" },
		{ id: 3, title: "Web Performance Tips", views: 432, date: "1 week ago", status: "Draft" }
	];

	return (
		<div className="space-y-8">
			{/* Header */}
			<div className="space-y-2">
				<div className="flex items-center justify-between">
					<div>
						<h1 className="text-4xl font-bold bg-gradient-to-r from-primary to-blue-600 bg-clip-text text-transparent">
							Dashboard
						</h1>
						<p className="text-muted-foreground mt-2">
							Welcome back, <span className="font-semibold text-foreground">{user?.name}</span>! Here's your overview.
						</p>
					</div>
					<Calendar className="h-12 w-12 text-muted-foreground/20" />
				</div>
			</div>

			{isAdmin && (
				<div className="space-y-8">
					{/* Stats Grid */}
					<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
						{stats.map((stat, index) => {
							const Icon = stat.icon;
							return (
								<Card key={index} className="p-6 border-border/50 bg-card hover:border-primary/50 transition-all">
									<div className="flex items-start justify-between mb-4">
										<div className={`p-3 rounded-lg ${stat.bgColor}`}>
											<Icon className={`h-6 w-6 ${stat.color}`} />
										</div>
										<div className={`flex items-center gap-1 text-xs font-semibold ${stat.isPositive ? 'text-emerald-500' : 'text-red-500'}`}>
											{stat.isPositive ? <ArrowUpRight className="h-3 w-3" /> : <ArrowDownRight className="h-3 w-3" />}
											{stat.change}
										</div>
									</div>
									<p className="text-sm text-muted-foreground mb-1">{stat.label}</p>
									<p className="text-3xl font-bold">{stat.value}</p>
								</Card>
							);
						})}
					</div>

					{/* Content Section */}
					<div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
						{/* Recent Posts */}
						<div className="lg:col-span-2">
							<Card className="p-8 border-border/50">
								<div className="flex items-center justify-between mb-6">
									<div>
										<h2 className="text-2xl font-bold">Recent Posts</h2>
										<p className="text-sm text-muted-foreground mt-1">Your latest blog articles</p>
									</div>
									<Button className="rounded-lg">View All</Button>
								</div>

								<div className="space-y-4">
									{recentPosts.map((post) => (
										<div key={post.id} className="flex items-center justify-between p-4 rounded-lg border border-border/30 hover:border-primary/50 hover:bg-muted/30 transition-all group">
											<div className="flex-1 min-w-0">
												<h3 className="font-semibold text-sm group-hover:text-primary transition-colors truncate">
													{post.title}
												</h3>
												<p className="text-xs text-muted-foreground mt-1">{post.date}</p>
											</div>
											<div className="flex items-center gap-6 ml-4">
												<div className="text-right">
													<p className="text-sm font-semibold">{post.views}</p>
													<p className="text-xs text-muted-foreground">views</p>
												</div>
												<div className={`px-3 py-1 rounded-full text-xs font-medium ${
													post.status === 'Published'
														? 'bg-emerald-500/10 text-emerald-600'
														: 'bg-amber-500/10 text-amber-600'
												}`}>
													{post.status}
												</div>
											</div>
										</div>
									))}
								</div>
							</Card>
						</div>

						{/* Quick Actions */}
						<div className="space-y-6">
							{/* Activity Card */}
							<Card className="p-6 border-border/50 bg-gradient-to-br from-primary/10 via-accent/5 to-background">
								<h3 className="text-lg font-bold mb-4">Activity</h3>
								<div className="space-y-3">
									<div className="flex items-center justify-between text-sm">
										<span className="text-muted-foreground">This week</span>
										<span className="font-semibold">12 posts</span>
									</div>
									<div className="w-full bg-border/50 rounded-full h-2">
										<div className="bg-primary rounded-full h-2 w-1/3"></div>
									</div>
									<div className="flex items-center justify-between text-sm">
										<span className="text-muted-foreground">This month</span>
										<span className="font-semibold">45 posts</span>
									</div>
									<div className="w-full bg-border/50 rounded-full h-2">
										<div className="bg-primary rounded-full h-2 w-3/4"></div>
									</div>
								</div>
							</Card>

							{/* Team Card */}
							<Card className="p-6 border-border/50">
								<h3 className="text-lg font-bold mb-4">Team Members</h3>
								<div className="space-y-3">
									{[1, 2, 3].map((i) => (
										<div key={i} className="flex items-center gap-3">
											<div className="w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center text-xs font-bold text-primary">
												U{i}
											</div>
											<div className="flex-1 min-w-0">
												<p className="text-sm font-medium">User {i}</p>
												<p className="text-xs text-muted-foreground">Active now</p>
											</div>
											<div className="w-2 h-2 rounded-full bg-emerald-500"></div>
										</div>
									))}
								</div>
							</Card>

							{/* Help Card */}
							<Card className="p-6 border-border/50 bg-gradient-to-br from-blue-500/10 to-purple-500/10">
								<h3 className="text-lg font-bold mb-2">Need Help?</h3>
								<p className="text-sm text-muted-foreground mb-4">
									Check out our documentation to get started.
								</p>
								<Button variant="outline" className="w-full rounded-lg text-xs">
									View Docs
								</Button>
							</Card>
						</div>
					</div>

					{/* Analytics Section */}
					<Card className="p-8 border-border/50">
						<div className="flex items-center justify-between mb-6">
							<div>
								<h2 className="text-2xl font-bold flex items-center gap-2">
									<BarChart3 className="h-6 w-6 text-primary" />
									Analytics Overview
								</h2>
								<p className="text-sm text-muted-foreground mt-1">Performance metrics for the last 30 days</p>
							</div>
						</div>

						<div className="grid grid-cols-1 md:grid-cols-3 gap-6">
							<div className="p-6 rounded-lg border border-border/30 bg-muted/20">
								<p className="text-sm text-muted-foreground mb-2">Avg. Daily Views</p>
								<p className="text-3xl font-bold">275</p>
								<p className="text-xs text-emerald-500 mt-2">↑ 12% from last month</p>
							</div>
							<div className="p-6 rounded-lg border border-border/30 bg-muted/20">
								<p className="text-sm text-muted-foreground mb-2">Conversion Rate</p>
								<p className="text-3xl font-bold">3.2%</p>
								<p className="text-xs text-red-500 mt-2">↓ 0.5% from last month</p>
							</div>
							<div className="p-6 rounded-lg border border-border/30 bg-muted/20">
								<p className="text-sm text-muted-foreground mb-2">Avg. Time on Page</p>
								<p className="text-3xl font-bold">3m 24s</p>
								<p className="text-xs text-emerald-500 mt-2">↑ 8% from last month</p>
							</div>
						</div>
					</Card>
				</div>
			)}

			{/* Non-Admin Message */}
			{!isAdmin && (
				<Card className="p-12 text-center border-border/50">
					<div className="mb-4">
						<BarChart3 className="h-12 w-12 text-muted-foreground/30 mx-auto" />
					</div>
					<h2 className="text-2xl font-bold mb-2">Dashboard Unavailable</h2>
					<p className="text-muted-foreground mb-6">
						You don't have permission to access the admin dashboard. Please contact an administrator for access.
					</p>
				</Card>
			)}
		</div>
	);
}
