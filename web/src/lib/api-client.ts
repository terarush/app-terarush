import axios from "axios"
import Cookies from "js-cookie"

const BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api"

export const apiClient = axios.create({
  baseURL: `${BASE_URL}/v1`,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 10000,
})

apiClient.interceptors.request.use(
  (config) => {
    const token = Cookies.get("accessToken")
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

apiClient.interceptors.response.use(
  (response) => {
    if (
      response.data &&
      typeof response.data === "object" &&
      "data" in response.data
    ) {
      response.data = response.data.data
    }
    return response
  },
  async (error) => {
    const originalRequest = error.config

    const skipRefreshUrls = [
      "/auth/login",
      "/auth/register",
      "/auth/refresh",
      "/auth/github/callback",
    ]
    const isSkipUrl = skipRefreshUrls.some((url) =>
      originalRequest?.url?.includes(url)
    )

    if (
      error.response?.status === 401 &&
      !originalRequest._retry &&
      !isSkipUrl
    ) {
      originalRequest._retry = true

      const refreshToken = Cookies.get("refreshToken")
      if (refreshToken) {
        try {
          const response = await axios.post(
            `${BASE_URL}/v1/auth/refresh`,
            {
              refresh_token: refreshToken,
            }
          )

          const data = response.data?.data || response.data
          const access_token = data.access_token || data.AccessToken

          if (access_token) {
            Cookies.set("accessToken", access_token, {
              expires: 7,
            })

            originalRequest.headers.Authorization = `Bearer ${access_token}`
            return apiClient(originalRequest)
          }
        } catch (refreshError) {
          Cookies.remove("accessToken")
          Cookies.remove("refreshToken")

          if (
            !window.location.pathname.includes("/login") &&
            !window.location.pathname.includes("/register")
          ) {
            window.location.href = "/login"
          }
          return Promise.reject(refreshError)
        }
      } else {
        Cookies.remove("accessToken")
        Cookies.remove("refreshToken")

        if (
          !window.location.pathname.includes("/login") &&
          !window.location.pathname.includes("/register") &&
          window.location.pathname !== "/"
        ) {
          window.location.href = "/login"
        }
      }
    }

    return Promise.reject(error)
  }
)
