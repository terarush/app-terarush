import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./assets/globals.css";
import Index from "./pages/Index";
import NotFound from "./pages/NotFound";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { ThemeProvider } from "./components/theme-provider";
import { Toaster } from "sonner";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

function AppContent() {
  // useSmoothScroll moved to specific pages where needed
  if (import.meta.env.VITE_NODE_ENV === "development") {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center">
          <h1 className="text-4xl font-bold mb-4">
            Development Mode
          </h1>
          <p className="text-lg text-gray-600">
            The application is running in development mode.
          </p>
        </div>
      </div>
    );
  }

  if (import.meta.env.VITE_NODE_ENV === "maintenance") {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center">
          <h1 className="text-4xl font-bold mb-4">
            Maintenance Mode
          </h1>
          <p className="text-lg text-gray-600">
            The application is currently under maintenance. Please
            check back later.
          </p>
        </div>
      </div>
    );
  }

  return (
    <BrowserRouter>
      <Toaster position="top-right" richColors />
      <Routes>
        <Route path="" element={<Index />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </BrowserRouter>
  );
}

createRoot(document.getElementById("root")!).render(
  <QueryClientProvider client={queryClient}>
    <ThemeProvider defaultTheme="system">
      <StrictMode>
        <AppContent />
      </StrictMode>
    </ThemeProvider>
  </QueryClientProvider>,
);
