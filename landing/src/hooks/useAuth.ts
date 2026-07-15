import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth as useAuthContext } from "@/contexts/AuthContext";

// Re-export the main useAuth hook from AuthContext
export { useAuth } from "@/contexts/AuthContext";

/**
 * Hook to get the current authenticated user
 * @returns The current user object or null if not authenticated
 */
export const useUser = () => {
	const { user } = useAuthContext();
	return user;
};

/**
 * Hook to check if user is authenticated
 * @returns Boolean indicating authentication status
 */
export const useIsAuthenticated = () => {
	const { isAuthenticated } = useAuthContext();
	return isAuthenticated;
};

/**
 * Hook to get auth loading state
 * @returns Boolean indicating if auth is loading
 */
export const useAuthLoading = () => {
	const { isLoading } = useAuthContext();
	return isLoading;
};

/**
 * Hook that redirects to login if user is not authenticated
 * Use this in components that require authentication
 *
 * @param redirectTo - Optional path to redirect to after login (default: current path)
 * @returns Object with isAuthenticated and isLoading status
 */
export const useRequireAuth = (redirectTo?: string) => {
	const { isAuthenticated, isLoading } = useAuthContext();
	const navigate = useNavigate();

	useEffect(() => {
		if (!isLoading && !isAuthenticated) {
			const currentPath = redirectTo || window.location.pathname;
			const loginPath =
				currentPath !== "/login"
					? `/login?redirect=${encodeURIComponent(currentPath)}`
					: "/login";
			navigate(loginPath, { replace: true });
		}
	}, [isAuthenticated, isLoading, navigate, redirectTo]);

	return { isAuthenticated, isLoading };
};

/**
 * Hook that checks if user has a specific role
 * Redirects to home page if user doesn't have the required role
 *
 * @param requiredRole - The role required to access the component (e.g., "admin")
 * @param redirectTo - Optional path to redirect to if role check fails (default: "/")
 * @returns Object with hasRole and isLoading status
 */
export const useRequireRole = (
	requiredRole: string,
	redirectTo: string = "/",
) => {
	const { user, isLoading } = useAuthContext();
	const navigate = useNavigate();

	const hasRole = user?.role === requiredRole;

	useEffect(() => {
		if (!isLoading && (!user || !hasRole)) {
			navigate(redirectTo, { replace: true });
		}
	}, [user, hasRole, isLoading, navigate, redirectTo]);

	return { hasRole, isLoading };
};

/**
 * Hook that checks if user has any of the specified roles
 * Redirects to home page if user doesn't have any of the required roles
 *
 * @param requiredRoles - Array of roles that are allowed (e.g., ["admin", "moderator"])
 * @param redirectTo - Optional path to redirect to if role check fails (default: "/")
 * @returns Object with hasAnyRole and isLoading status
 */
export const useRequireAnyRole = (
	requiredRoles: string[],
	redirectTo: string = "/",
) => {
	const { user, isLoading } = useAuthContext();
	const navigate = useNavigate();

	const hasAnyRole = user && requiredRoles.includes(user.role);

	useEffect(() => {
		if (!isLoading && !hasAnyRole) {
			navigate(redirectTo, { replace: true });
		}
	}, [hasAnyRole, isLoading, navigate, redirectTo]);

	return { hasAnyRole, isLoading };
};

/**
 * Hook to get login function
 * @returns Login function
 */
export const useLogin = () => {
	const { login } = useAuthContext();
	return login;
};

/**
 * Hook to get register function
 * @returns Register function
 */
export const useRegister = () => {
	const { register } = useAuthContext();
	return register;
};

/**
 * Hook to get logout function with optional redirect
 * @returns Logout function with redirect capability
 */
export const useLogout = () => {
	const { logout } = useAuthContext();
	const navigate = useNavigate();

	const logoutWithRedirect = async (redirectTo: string = "/login") => {
		await logout();
		navigate(redirectTo, { replace: true });
	};

	return logoutWithRedirect;
};

/**
 * Hook to refresh user profile data
 * @returns Refresh user function
 */
export const useRefreshUser = () => {
	const { refreshUser } = useAuthContext();
	return refreshUser;
};

/**
 * Hook that redirects authenticated users away from auth pages
 * Use this on login/register pages to redirect already logged-in users
 *
 * @param redirectTo - Path to redirect authenticated users to (default: "/dashboard")
 */
export const useRedirectIfAuthenticated = (
	redirectTo: string = "/dashboard",
) => {
	const { isAuthenticated, isLoading } = useAuthContext();
	const navigate = useNavigate();

	useEffect(() => {
		if (!isLoading && isAuthenticated) {
			navigate(redirectTo, { replace: true });
		}
	}, [isAuthenticated, isLoading, navigate, redirectTo]);

	return { isAuthenticated, isLoading };
};

/**
 * Hook to check if current user is admin
 * @returns Boolean indicating if user is admin
 */
export const useIsAdmin = () => {
	const { user } = useAuthContext();
	return user?.role === "admin";
};

/**
 * Hook to check if current user matches a specific user ID
 * Useful for checking ownership of resources
 *
 * @param userId - User ID to compare with current user
 * @returns Boolean indicating if current user matches the provided ID
 */
export const useIsCurrentUser = (userId: number) => {
	const { user } = useAuthContext();
	return user?.id === userId;
};
