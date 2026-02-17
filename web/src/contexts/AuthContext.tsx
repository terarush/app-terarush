import React, { createContext, useContext, useState, useEffect } from "react";
import type { ReactNode } from "react";
import Cookies from "js-cookie";
import { authApi } from "@/lib/api/auth";
import type { UserResponse, LoginRequest, RegisterRequest } from "@/lib/api/auth";

interface AuthContextType {
  user: UserResponse | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<UserResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const isAuthenticated = !!user && !!Cookies.get("accessToken");

  // Fetch user profile on mount if token exists
  useEffect(() => {
    const initAuth = async () => {
      const token = Cookies.get("accessToken");
      
      if (token) {
        try {
          const userProfile = await authApi.getProfile();
          setUser(userProfile);
        } catch (error) {
          console.error("Failed to fetch user profile:", error);
          // Clear invalid tokens
          Cookies.remove("accessToken");
          Cookies.remove("refreshToken");
        }
      }
      
      setIsLoading(false);
    };

    initAuth();
  }, []);

  const login = async (credentials: LoginRequest) => {
    try {
      await authApi.login(credentials);
      const userProfile = await authApi.getProfile();
      setUser(userProfile);
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    }
  };

  const register = async (data: RegisterRequest) => {
    try {
      await authApi.register(data);
      const userProfile = await authApi.getProfile();
      setUser(userProfile);
    } catch (error) {
      console.error("Registration failed:", error);
      throw error;
    }
  };

  const logout = async () => {
    try {
      await authApi.logout();
    } catch (error) {
      console.error("Logout failed:", error);
    } finally {
      setUser(null);
      Cookies.remove("accessToken");
      Cookies.remove("refreshToken");
    }
  };

  const refreshUser = async () => {
    try {
      const userProfile = await authApi.getProfile();
      setUser(userProfile);
    } catch (error) {
      console.error("Failed to refresh user:", error);
      throw error;
    }
  };

  const value: AuthContextType = {
    user,
    isLoading,
    isAuthenticated,
    login,
    register,
    logout,
    refreshUser,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// Custom hook to use auth context
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  
  return context;
};
