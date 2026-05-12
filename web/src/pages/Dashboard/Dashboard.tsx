import { useUser, useIsAdmin } from "@/hooks";

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
					<div>
						hello world
					</div>
				)}
			</div>
		</>
	);
}
