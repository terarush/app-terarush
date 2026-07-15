import axios from "axios";
import Cookies from "js-cookie";

const API_BASE_URL = import.meta.env.VITE_API_URL;

// Create axios instance with default config
export const apiClient = axios.create({
	baseURL: `${API_BASE_URL}/api/v1`,
	headers: {
		"Content-Type": "application/json",
	},
	timeout: 10000,
});

// Request interceptor - add auth token to requests
apiClient.interceptors.request.use(
	(config) => {
		const token = Cookies.get("accessToken");
		if (token) {
			config.headers.Authorization = `Bearer ${token}`;
		}
		return config;
	},
	(error) => {
		return Promise.reject(error);
	},
);

// Response interceptor - unwrap response and handle token refresh on 401
apiClient.interceptors.response.use(
	(response) => {
		// Unwrap the response data if it's wrapped
		if (
			response.data &&
			typeof response.data === "object" &&
			"data" in response.data
		) {
			response.data = response.data.data;
		}
		return response;
	},
	async (error) => {
		const originalRequest = error.config;

		// If 401 and we haven't retried yet, try to refresh token
		// BUT skip refresh for login/register/refresh endpoints
		const skipRefreshUrls = [
			"/auth/login",
			"/auth/register",
			"/auth/refresh",
			"/auth/github/callback",
		];
		const isSkipUrl = skipRefreshUrls.some((url) =>
			originalRequest?.url?.includes(url),
		);

		if (
			error.response?.status === 401 &&
			!originalRequest._retry &&
			!isSkipUrl
		) {
			originalRequest._retry = true;

			const refreshToken = Cookies.get("refreshToken");
			if (refreshToken) {
				try {
					const response = await axios.post(
						`${API_BASE_URL}/api/v1/auth/refresh`,
						{
							refresh_token: refreshToken,
						},
					);

					// Extract access_token from wrapped response
					const data = response.data?.data || response.data;
					const access_token = data.access_token || data.AccessToken;

					if (access_token) {
						Cookies.set("accessToken", access_token, {
							expires: 7,
						});

						// Retry original request with new token
						originalRequest.headers.Authorization = `Bearer ${access_token}`;
						return apiClient(originalRequest);
					}
				} catch (refreshError) {
					// Refresh failed, clear tokens and redirect to login
					Cookies.remove("accessToken");
					Cookies.remove("refreshToken");

					// Only redirect if not already on login/register page
					if (
						!window.location.pathname.includes("/login") &&
						!window.location.pathname.includes("/register")
					) {
						window.location.href = "/login";
					}
					return Promise.reject(refreshError);
				}
			} else {
				// No refresh token, clear cookies and redirect only if authenticated page
				Cookies.remove("accessToken");
				Cookies.remove("refreshToken");

				if (
					!window.location.pathname.includes("/login") &&
					!window.location.pathname.includes("/register") &&
					window.location.pathname !== "/"
				) {
					window.location.href = "/login";
				}
			}
		}

		return Promise.reject(error);
	},
);
