import { Outlet } from "react-router-dom";
import { AppSidebar } from "@/components/app-sidebar";
import { ProtectedRoute } from "@/components/ProtectedRoute";

export function DashboardLayout() {
	return (
		<ProtectedRoute>
			<div className="flex h-screen overflow-hidden bg-background">
				<AppSidebar />

				<main className="flex-1 overflow-y-auto">
					<Outlet />
				</main>
			</div>
		</ProtectedRoute>
	);
}
