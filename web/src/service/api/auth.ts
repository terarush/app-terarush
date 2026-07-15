import Cookies from "js-cookie"
import { apiClient } from "../../lib/api-client"

// Types
export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  name: string
  email: string
  password: string
  confirm_password: string
  role?: "admin" | "user"
}

export interface ChangePasswordRequest {
  old_password: string
  new_password: string
  confirm_password: string
}

export interface UpdateProfileRequest {
  name: string
  email: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
  user: UserResponse
}

export interface UserResponse {
  id: number
  name: string
  email: string
  avatar: string
  bio?: string
  banner?: string
  role: string
  is_banned?: boolean
  provider?: string
  created_at: string
  updated_at: string
}

export interface ApiError {
  error: string
}

// Auth API functions
export const authApi = {
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/register", data)
    const authData = response.data

    if (authData.access_token) {
      Cookies.set("accessToken", authData.access_token, { expires: 7 })
    }
    if (authData.refresh_token) {
      Cookies.set("refreshToken", authData.refresh_token, { expires: 30 })
    }

    return authData
  },

  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/login", data)
    const authData = response.data

    if (authData.access_token) {
      Cookies.set("accessToken", authData.access_token, { expires: 7 })
    }
    if (authData.refresh_token) {
      Cookies.set("refreshToken", authData.refresh_token, { expires: 30 })
    }

    return authData
  },

  getProfile: async (): Promise<UserResponse> => {
    const response = await apiClient.get<UserResponse>("/auth/profile")
    return response.data
  },

  updateProfile: async (data: UpdateProfileRequest): Promise<UserResponse> => {
    const response = await apiClient.put<UserResponse>("/auth/profile", data)
    return response.data
  },

  changePassword: async (data: ChangePasswordRequest): Promise<void> => {
    await apiClient.post("/auth/change-password", data)
  },

  refreshToken: async (refreshToken: string): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/refresh", {
      refresh_token: refreshToken,
    })
    const authData = response.data

    if (authData.access_token) {
      Cookies.set("accessToken", authData.access_token, { expires: 7 })
    }

    return authData
  },

  logout: async (): Promise<void> => {
    try {
      await apiClient.post("/auth/logout")
    } finally {
      Cookies.remove("accessToken")
      Cookies.remove("refreshToken")
    }
  },

  getGitHubAuthUrl: (): string => {
    const clientId = import.meta.env.VITE_GITHUB_CLIENT_ID;

    if (!clientId) {
      throw new Error("Missing required ");
    }

    const redirectUri = `${window.location.origin}/auth/github/callback`
    const scope = "user:email"
    const state = Math.random().toString(36).substring(7)

    sessionStorage.setItem("github_oauth_state", state)

    return `https://github.com/login/oauth/authorize?client_id=${clientId}&redirect_uri=${encodeURIComponent(redirectUri)}&scope=${scope}&state=${state}`
  },

  githubCallback: async (code: string): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/github/callback", {
      code,
    })
    const authData = response.data

    if (authData.access_token) {
      Cookies.set("accessToken", authData.access_token, { expires: 7 })
    }
    if (authData.refresh_token) {
      Cookies.set("refreshToken", authData.refresh_token, { expires: 30 })
    }

    return authData
  },
}
