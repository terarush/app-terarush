import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { usersApi } from '@/service/api/users';

// Query Key Factory for clean cache management
export const userKeys = {
  all: ['users'] as const,
  lists: () => [...userKeys.all, 'list'] as const,
  details: () => [...userKeys.all, 'detail'] as const,
  detail: (id: number | string) => [...userKeys.details(), id] as const,
};

// Hook to get all users
export function useUsers() {
  return useQuery({
    queryKey: userKeys.lists(),
    queryFn: usersApi.list,
  });
}

// Hook to get a single user by ID
export function useUser(id: number | string) {
  return useQuery({
    queryKey: userKeys.detail(id),
    queryFn: () => usersApi.get(id),
    enabled: !!id, // Only run query if id is provided
  });
}

// Hook to create a new user
export function useCreateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: usersApi.create,
    onSuccess: (newUser) => {
      // Optimistically update or invalidate the query cache
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
      
      // Optionally pre-populate the detail cache for this new user
      queryClient.setQueryData(userKeys.detail(newUser.id), newUser);
    },
  });
}

