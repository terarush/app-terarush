import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./assets/globals.css";
import Index from "./pages/Index";
import NotFound from "./pages/NotFound";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { ThemeProvider } from "./components/theme-provider";

createRoot(document.getElementById("root")!).render(
	<ThemeProvider defaultTheme="system">
		<StrictMode>
			<BrowserRouter>
				<Routes>
					<Route path="/" element={<Index />} />
					<Route path="*" element={<NotFound />} />
				</Routes>
			</BrowserRouter>
		</StrictMode>
	</ThemeProvider>,
);
