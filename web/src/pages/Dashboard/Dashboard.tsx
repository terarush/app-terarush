import { useUser, useIsAdmin } from "@/hooks";
import { Settings } from "lucide-react";

export default function Dashboard() {
	const user = useUser();
	const isAdmin = useIsAdmin();
	return (
		<>
			<div className="mb-8">
				<h1 className="text-3xl font-bold">Dashboard</h1>
				<p className="text-muted-foreground mt-1">
					Welcome back, {user?.name}! Here's your overview.
				</p>
			</div>

			<div className="space-y-8">
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
