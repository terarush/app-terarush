import { apiClient } from '@/lib/api-client';
import type { User } from '@/types';

// API endpoints for User resources
export const usersApi = {
  // GET /users
  list: () => 
    apiClient.get<User[]>('/users'),

  // GET /users/:id
  get: (id: number | string) => 
    apiClient.get<User>(`/users/${id}`),

  // POST /users
  create: (user: Omit<User, 'id'>) => 
    apiClient.post<User>('/users', user),
};
