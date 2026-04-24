import axios from "axios";
import { API_BASE_URL } from "./config";

const apiClient = axios.create({
    baseURL: API_BASE_URL + "/api/v1",
    withCredentials: true,
    headers: {
        "Content-Type": "application/json",
    },
});

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {

    console.log("API Error:", error);
    const status = error.response?.status;

    const isAuthError = status === 401 || status === 403;

    const publicRoutes = [
      '/login',
      '/forgot-password',
      '/forgot-username',
      '/reset-password'
    ];

    if (
      isAuthError &&
      typeof window !== 'undefined' &&
      !publicRoutes.includes(window.location.pathname)
    ) {
      window.location.replace('/login?error=session_expired');
    }

    return Promise.reject(error);
  }
);

export default apiClient;