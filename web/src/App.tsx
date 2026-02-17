import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./assets/globals.css";
import Index from "./pages/Index";
import NotFound from "./pages/NotFound";
import Login from "./pages/Auth/Login";
import Register from "./pages/Auth/Register";
import Dashboard from "./pages/Dashboard/Dashboard";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { ThemeProvider } from "./components/theme-provider";
import { AuthProvider } from "./contexts/AuthContext";
import { ProtectedRoute } from "./components/ProtectedRoute";

createRoot(document.getElementById("root")!).render(
	<ThemeProvider defaultTheme="system">
		<StrictMode>
			<BrowserRouter>
				<AuthProvider>
					<Routes>
						<Route path="" element={<Index />} />
						<Route path="login" element={<Login />} />
						<Route path="register" element={<Register />} />
						<Route 
							path="dashboard" 
							element={
								<ProtectedRoute>
									<Dashboard />
								</ProtectedRoute>
							} 
						/>
						<Route path="*" element={<NotFound />} />
					</Routes>
				</AuthProvider>
			</BrowserRouter>
		</StrictMode>
	</ThemeProvider>,
);
