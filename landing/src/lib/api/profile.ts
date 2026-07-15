import { apiClient } from "./client";

export interface User {
  id: number;
  name: string;
  email: string;
  avatar: string;
  bio?: string;
  banner?: string;
  role: string;
  is_banned?: boolean;
  provider?: string;
  created_at: string;
  updated_at: string;
}

export interface UpdateProfileData {
  name: string;
  avatar?: string;
  bio?: string;
  banner?: string;
}

export const profileApi = {
  // Get current user profile
  getCurrentUser: async (): Promise<User> => {
    const response = await apiClient.get("/profile");
    return response.data;
  },

  // Update profile
  updateProfile: async (data: UpdateProfileData): Promise<User> => {
    const response = await apiClient.put("/profile", data);
    return response.data;
  },

  // Upload avatar
  uploadAvatar: async (file: File): Promise<{ avatar_url: string; message: string }> => {
    const formData = new FormData();
    formData.append("avatar", file);

    const response = await apiClient.post("/profile/avatar", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
    return response.data;
  },

  // Upload banner
  uploadBanner: async (file: File): Promise<{ banner_url: string; message: string }> => {
    const formData = new FormData();
    formData.append("banner", file);

    const response = await apiClient.post("/profile/banner", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
    return response.data;
  },
};
