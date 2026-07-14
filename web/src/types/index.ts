// Domain models and common TS types for the application

export interface User {
  id: number;
  name: string;
  username: string;
  email: string;
  phone?: string;
  website?: string;
  company?: {
    name: string;
    catchPhrase: string;
    bs: string;
  };
}

export interface ApiResponse<T> {
  data: T;
  message?: string;
  status: number;
}
