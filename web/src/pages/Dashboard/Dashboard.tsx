import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useUser, useIsAdmin } from "@/hooks";
import { Server, Package, Receipt, Settings, Users } from "lucide-react";

export default function Dashboard() {
	const user = useUser();
	const isAdmin = useIsAdmin();
	const navigate = useNavigate();

	const quickActions = [
		{
			title: "Profile",
			description: "Update your information",
			icon: Users,
			path: "/dashboard/profile",
			show: true,
		},
		{
			title: "My Transactions",
			description: "View your purchase history",
			icon: Receipt,
			path: "/dashboard/transactions",
			show: true,
		},
		{
			title: "Nodes",
			description: "Manage your containers",
			icon: Server,
			path: "/dashboard/nodes",
			show: isAdmin,
		},
		{
			title: "Products",
			description: "Manage container products",
			icon: Package,
			path: "/dashboard/products",
			show: isAdmin,
		},
		{
			title: "Settings",
			description: "Manage preferences",
			icon: Settings,
			path: "/dashboard/settings",
			show: true,
		},
	].filter((action) => action.show);

	return (
		<>
			<div className="mb-8">
				<h1 className="text-3xl font-bold">Dashboard</h1>
				<p className="text-muted-foreground mt-1">
					Welcome back, {user?.name}! Here's your overview.
				</p>
			</div>

			<div className="space-y-8">
				{/* User Info */}
				<div className="border rounded-lg p-6">
					<h2 className="text-xl font-semibold mb-4">Account Information</h2>
					<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
						<div>
							<p className="text-sm text-muted-foreground">Name</p>
							<p className="text-lg font-medium">{user?.name}</p>
						</div>
						<div>
							<p className="text-sm text-muted-foreground">Email</p>
							<p className="text-lg font-medium">{user?.email}</p>
						</div>
						<div>
							<p className="text-sm text-muted-foreground">Role</p>
							<p className="text-lg font-medium capitalize">{user?.role}</p>
						</div>
						<div>
							<p className="text-sm text-muted-foreground">Member Since</p>
							<p className="text-lg font-medium">
								{user?.created_at &&
									new Date(user.created_at).toLocaleDateString("id-ID", {
										year: "numeric",
										month: "short",
										day: "numeric",
									})}
							</p>
						</div>
					</div>
				</div>

				{/* Quick Actions */}
				<div className="border rounded-lg p-6">
					<h2 className="text-xl font-semibold mb-4">Quick Actions</h2>
					<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
						{quickActions.map((action) => (
							<Button
								key={action.path}
								variant="outline"
								className="h-auto p-6 flex flex-col items-start gap-2 text-left"
								onClick={() => navigate(action.path)}
							>
								<action.icon className="h-6 w-6 text-primary" />
								<div>
									<p className="font-semibold">{action.title}</p>
									<p className="text-xs text-muted-foreground">
										{action.description}
									</p>
								</div>
							</Button>
						))}
					</div>
				</div>

				{/* Admin Badge */}
				{isAdmin && (
					<div className="border border-primary rounded-lg p-6 bg-primary/5">
						<div className="flex items-center gap-3">
							<div className="w-12 h-12 rounded-full bg-primary/20 flex items-center justify-center">
								<Settings className="h-6 w-6 text-primary" />
							</div>
							<div>
								<h3 className="font-semibold">Admin Access</h3>
								<p className="text-sm text-muted-foreground">
									You have administrator privileges
								</p>
							</div>
						</div>
					</div>
				)}
			</div>
		</>
	);
}
