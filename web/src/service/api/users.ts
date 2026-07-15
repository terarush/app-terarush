import { apiClient } from "@/lib/api-client"
import type { User } from "@/types"

export const usersApi = {
  list: async () => {
    const res = await apiClient.get<User[]>("/users")
    return res.data
  },

  get: async (id: number | string) => {
    const res = await apiClient.get<User>(`/users/${id}`)
    return res.data
  },

  create: async (user: Omit<User, "id">) => {
    const res = await apiClient.post<User>("/users", user)
    return res.data
  },
}
