// Simple, type-safe API client wrapper around the standard Fetch API

const BASE_URL = import.meta.env.VITE_API_URL || 'https://jsonplaceholder.typicode.com';

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const response = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || `HTTP error! status: ${response.status}`);
  }

  // Handle empty responses
  if (response.status === 204) {
    return {} as T;
  }

  return response.json();
}

export const apiClient = {
  get: <T>(path: string, options?: RequestInit) => 
    request<T>(path, { ...options, method: 'GET' }),
  
  post: <T>(path: string, data: unknown, options?: RequestInit) => 
    request<T>(path, { 
      ...options, 
      method: 'POST', 
      body: JSON.stringify(data) 
    }),
  
  put: <T>(path: string, data: unknown, options?: RequestInit) => 
    request<T>(path, { 
      ...options, 
      method: 'PUT', 
      body: JSON.stringify(data) 
    }),
  
  delete: <T>(path: string, options?: RequestInit) => 
    request<T>(path, { ...options, method: 'DELETE' }),
};
