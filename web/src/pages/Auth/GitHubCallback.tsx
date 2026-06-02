import { useEffect, useState, useRef } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { authApi } from "@/lib/api/auth";
import { useAuth } from "@/contexts/AuthContext";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";

export default function GitHubCallback() {
	const [searchParams] = useSearchParams();
	const navigate = useNavigate();
	const { refreshUser } = useAuth();
	const [error, setError] = useState<string | null>(null);
	const hasRun = useRef(false);

	useEffect(() => {
		// Prevent double execution in React Strict Mode
		if (hasRun.current) return;
		hasRun.current = true;

		const handleCallback = async () => {
			const code = searchParams.get("code");
			const state = searchParams.get("state");
			const error = searchParams.get("error");

			// Check for GitHub OAuth errors
			if (error) {
				const errorDescription =
					searchParams.get("error_description") ||
					"Authentication failed";
				setError(errorDescription);
				toast.error(errorDescription);
				setTimeout(() => navigate("/login"), 3000);
				return;
			}

			// Validate required parameters
			if (!code) {
				setError("No authorization code received");
				toast.error("No authorization code received");
				setTimeout(() => navigate("/login"), 3000);
				return;
			}

			// Verify state (CSRF protection)
			const savedState = sessionStorage.getItem("github_oauth_state");
			if (state !== savedState) {
				setError("Invalid state parameter - possible CSRF attack");
				toast.error("Authentication failed - security check");
				setTimeout(() => navigate("/login"), 3000);
				return;
			}

			// Clear state from storage
			sessionStorage.removeItem("github_oauth_state");

			try {
				// Exchange code for tokens
				const authResponse = await authApi.githubCallback(code);

				// Refresh user from context (tokens already saved by API client)
				await refreshUser();

				// Show success message
				toast.success(`Welcome, ${authResponse.user.name}!`);

				// Redirect to dashboard
				navigate("/dashboard");
			} catch (err: any) {
				const errorMessage =
					err.response?.data?.message ||
					err.response?.data?.error ||
					"Failed to authenticate with GitHub";
				setError(errorMessage);
				toast.error(errorMessage);
				console.error("GitHub auth error details:", err.response?.data);
				setTimeout(() => navigate("/login"), 3000);
			}
		};

		handleCallback();
	}, [searchParams, navigate, refreshUser]);

	return (
		<div className="min-h-screen flex items-center justify-center bg-background">
			<div className="text-center space-y-4 p-8">
				{!error ? (
					<>
						<Loader2 className="h-12 w-12 animate-spin text-primary mx-auto" />
						<h2 className="text-2xl font-semibold text-gray-900 dark:text-white">
							Authenticating with GitHub
						</h2>
						<p className="text-gray-600 dark:text-gray-400">
							Please wait while we complete your login...
						</p>
					</>
				) : (
					<>
						<div className="h-12 w-12 rounded-full bg-destructive/10 flex items-center justify-center mx-auto">
							<svg
								className="h-6 w-6 text-destructive"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									strokeWidth={2}
									d="M6 18L18 6M6 6l12 12"
								/>
							</svg>
						</div>
						<h2 className="text-2xl font-semibold text-gray-900 dark:text-white">
							Authentication Failed
						</h2>
						<p className="text-gray-600 dark:text-gray-400 max-w-md">
							{error}
						</p>
						<p className="text-sm text-gray-500 dark:text-gray-500">
							Redirecting to login page...
						</p>
					</>
				)}
			</div>
		</div>
	);
}
