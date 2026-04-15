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
import Nodes from "./pages/Dashboard/Nodes";
import { DashboardLayout } from "./pages/Dashboard/Layout";
import { ProductList } from "./pages/Products/ProductList";
import { ProductDetail } from "./pages/Products/ProductDetail";
import { MyTransactions } from "./pages/Dashboard/MyTransactions";
import { Products } from "./pages/Dashboard/Products";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { ThemeProvider } from "./components/theme-provider";
import { AuthProvider } from "./contexts/AuthContext";
import { Toaster } from "sonner";
import Settings from "./pages/Dashboard/Settings";
import { useSmoothScroll } from "./hooks/useSmoothScroll";

function AppContent() {
	useSmoothScroll();

	return (
		<BrowserRouter>
			<AuthProvider>
				<Toaster position="top-right" richColors />
				<Routes>
					<Route path="" element={<Index />} />
					<Route path="login" element={<Login />} />
					<Route path="register" element={<Register />} />
					<Route path="auth/github/callback" element={<GitHubCallback />} />
					<Route path="products" element={<ProductList />} />
					<Route path="products/:id" element={<ProductDetail />} />
					<Route path="dashboard" element={<DashboardLayout />}>
						<Route index element={<Dashboard />} />
						<Route path="profile" element={<Profile />} />
						<Route path="nodes" element={<Nodes />} />
						<Route path="transactions" element={<MyTransactions />} />
						<Route path="products" element={<Products />} />
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
