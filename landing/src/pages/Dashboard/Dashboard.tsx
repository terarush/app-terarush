import { useUser, useIsAdmin } from "@/hooks";
import {
	BarChart3,
	Eye,
	FileText,
	Users,
	TrendingUp,
	Calendar,
	ArrowUpRight,
	ArrowDownRight,
} from "lucide-react";
import { Card } from "@/components/ui/card";

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
			bgColor: "bg-blue-500/10",
		},
		{
			label: "Total Views",
			value: "8,245",
			change: "+23%",
			isPositive: true,
			icon: Eye,
			color: "text-emerald-500",
			bgColor: "bg-emerald-500/10",
		},
		{
			label: "Active Users",
			value: "1,234",
			change: "+5%",
			isPositive: true,
			icon: Users,
			color: "text-purple-500",
			bgColor: "bg-purple-500/10",
		},
		{
			label: "Engagement Rate",
			value: "68%",
			change: "-2%",
			isPositive: false,
			icon: TrendingUp,
			color: "text-amber-500",
			bgColor: "bg-amber-500/10",
		},
	];

	return (
		<div className="space-y-8">
			{/* Header */}
			<div className="space-y-2">
				<div className="flex items-center justify-between">
					<div>
						<h1 className="text-4xl font-bold bg-linear-to-r from-primary to-blue-600 bg-clip-text text-transparent">
							Dashboard
						</h1>
						<p className="text-muted-foreground mt-2">
							Welcome back,{" "}
							<span className="font-semibold text-foreground">
								{user?.name}
							</span>
							! Here's your overview.
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
								<Card
									key={index}
									className="p-6 border-border/50 bg-card hover:border-primary/50 transition-all"
								>
									<div className="flex items-start justify-between mb-4">
										<div
											className={`p-3 rounded-lg ${stat.bgColor}`}
										>
											<Icon
												className={`h-6 w-6 ${stat.color}`}
											/>
										</div>
										<div
											className={`flex items-center gap-1 text-xs font-semibold ${stat.isPositive ? "text-emerald-500" : "text-red-500"}`}
										>
											{stat.isPositive ? (
												<ArrowUpRight className="h-3 w-3" />
											) : (
												<ArrowDownRight className="h-3 w-3" />
											)}
											{stat.change}
										</div>
									</div>
									<p className="text-sm text-muted-foreground mb-1">
										{stat.label}
									</p>
									<p className="text-3xl font-bold">
										{stat.value}
									</p>
								</Card>
							);
						})}
					</div>
				</div>
			)}

			{/* Non-Admin Message */}
			{!isAdmin && (
				<Card className="p-12 text-center border-border/50">
					<div className="mb-4">
						<BarChart3 className="h-12 w-12 text-muted-foreground/30 mx-auto" />
					</div>
					<h2 className="text-2xl font-bold mb-2">
						Dashboard Unavailable
					</h2>
					<p className="text-muted-foreground mb-6">
						You don't have permission to access the admin dashboard.
						Please contact an administrator for access.
					</p>
				</Card>
			)}
		</div>
	);
}
