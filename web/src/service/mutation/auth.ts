import { useMutation } from "@tanstack/react-query"
import { useAuth } from "@/contexts/auth-context"
import type { LoginRequest, RegisterRequest } from "../api/auth"

export function useLoginMutation() {
  const { login } = useAuth()
  
  return useMutation({
    mutationFn: (credentials: LoginRequest) => login(credentials),
  })
}

export function useRegisterMutation() {
  const { register } = useAuth()
  
  return useMutation({
    mutationFn: (data: RegisterRequest) => register(data),
  })
}

export function useGitHubLoginMutation() {
  const { loginWithGitHub } = useAuth()
  
  return useMutation({
    mutationFn: (code: string) => loginWithGitHub(code),
  })
}
