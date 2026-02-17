import Cookies from "js-cookie";
import { apiClient } from "./client";

// Types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
  confirm_password: string;
  role?: "admin" | "user";
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
  confirm_password: string;
}

export interface UpdateProfileRequest {
  name: string;
  email: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface UserResponse {
  id: number;
  name: string;
  email: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface ApiError {
  error: string;
}

// Auth API functions
export const authApi = {
  /**
   * Register a new user
   */
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/register", data);
    const { access_token, refresh_token } = response.data;
    
    // Store tokens in cookies
    Cookies.set("accessToken", access_token, { expires: 7 });
    Cookies.set("refreshToken", refresh_token, { expires: 30 });
    
    return response.data;
  },

  /**
   * Login user
   */
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/login", data);
    const { access_token, refresh_token } = response.data;
    
    // Store tokens in cookies
    Cookies.set("accessToken", access_token, { expires: 7 });
    Cookies.set("refreshToken", refresh_token, { expires: 30 });
    
    return response.data;
  },

  /**
   * Get current user profile
   */
  getProfile: async (): Promise<UserResponse> => {
    const response = await apiClient.get<UserResponse>("/auth/profile");
    return response.data;
  },

  /**
   * Update current user profile
   */
  updateProfile: async (data: UpdateProfileRequest): Promise<UserResponse> => {
    const response = await apiClient.put<UserResponse>("/auth/profile", data);
    return response.data;
  },

  /**
   * Change password
   */
  changePassword: async (data: ChangePasswordRequest): Promise<void> => {
    await apiClient.post("/auth/change-password", data);
  },

  /**
   * Refresh access token
   */
  refreshToken: async (refreshToken: string): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/refresh", {
      refresh_token: refreshToken,
    });
    const { access_token } = response.data;
    Cookies.set("accessToken", access_token, { expires: 7 });
    return response.data;
  },

  /**
   * Logout user
   */
  logout: async (): Promise<void> => {
    try {
      await apiClient.post("/auth/logout");
    } finally {
      // Clear tokens from cookies
      Cookies.remove("accessToken");
      Cookies.remove("refreshToken");
    }
  },
};
