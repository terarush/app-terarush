import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./assets/globals.css";
import Index from "./pages/Index";
import NotFound from "./pages/NotFound";
import Login from "./pages/Auth/Login";
import Register from "./pages/Auth/Register";
import GitHubCallback from "./pages/Auth/GitHubCallback";
import Dashboard from "./pages/Dashboard/Dashboard";
import Profile from "./pages/Dashboard/Profile";
import { BlogManagement } from "./pages/Dashboard/BlogManagement";
import { DashboardLayout } from "./pages/Dashboard/Layout";
import { BlogList } from "./pages/Blog/BlogList";
import { BlogDetail } from "./pages/Blog/BlogDetail";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { ThemeProvider } from "./components/theme-provider";
import { AuthProvider } from "./contexts/AuthContext";
import { Toaster } from "sonner";
import Settings from "./pages/Dashboard/Settings";
import { useSmoothScroll } from "./hooks/useSmoothScroll";

function AppContent() {
	useSmoothScroll();
	if (import.meta.env.VITE_NODE_ENV === "development") {
		return (
			<div className="flex items-center justify-center h-screen">
				<div className="text-center">
					<h1 className="text-4xl font-bold mb-4">Development Mode</h1>
					<p className="text-lg text-gray-600">The application is running in development mode.</p>
				</div>
			</div>
		)
	}

	if (import.meta.env.VITE_NODE_ENV === "maintenance") {
		return (
			<div className="flex items-center justify-center h-screen">
				<div className="text-center">
					<h1 className="text-4xl font-bold mb-4">Maintenance Mode</h1>
					<p className="text-lg text-gray-600">The application is currently under maintenance. Please check back later.</p>
				</div>
			</div>
		)
	}

	return (
		<BrowserRouter>
			<AuthProvider>
				<Toaster position="top-right" richColors />
				<Routes>
					<Route path="" element={<Index />} />
					<Route path="login" element={<Login />} />
					<Route path="register" element={<Register />} />
					<Route path="auth/github/callback" element={<GitHubCallback />} />
					<Route path="blog" element={<BlogList />} />
					<Route path="blog/:slug" element={<BlogDetail />} />
					<Route path="dashboard" element={<DashboardLayout />}>
						<Route index element={<Dashboard />} />
						<Route path="profile" element={<Profile />} />
						<Route path="blogs" element={<BlogManagement />} />
						<Route path="settings" element={<Settings />} />
					</Route>
					<Route path="*" element={<NotFound />} />
				</Routes>
			</AuthProvider>
		</BrowserRouter>
	);
}

createRoot(document.getElementById("root")!).render(
	<ThemeProvider defaultTheme="system">
		<StrictMode>
			<AppContent />
		</StrictMode>
	</ThemeProvider>,
);
