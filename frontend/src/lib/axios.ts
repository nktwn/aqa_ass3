// src/lib/axios.ts

import axios from "axios";

const defaultApiBaseUrl = "http://localhost:8080/api";
const configuredApiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || defaultApiBaseUrl;

const api = axios.create({
    baseURL: configuredApiBaseUrl.replace(/\/+$/, ""),
    withCredentials: false,
});

api.interceptors.request.use((config) => {
    if (typeof window === "undefined") {
        return config;
    }

    const expiresAt = localStorage.getItem("expires_at");
    if (expiresAt && Date.now() >= parseInt(expiresAt, 10)) {
        localStorage.removeItem("access_token");
        localStorage.removeItem("refresh_token");
        localStorage.removeItem("expires_at");
    }

    const token = localStorage.getItem("access_token");
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (typeof window !== "undefined" && error?.response?.status === 401) {
            localStorage.removeItem("access_token");
            localStorage.removeItem("refresh_token");
            localStorage.removeItem("expires_at");
            window.dispatchEvent(new Event("auth:unauthorized"));

            if (!window.location.pathname.startsWith("/login") && !window.location.pathname.startsWith("/register")) {
                window.location.href = "/login";
            }
        }

        return Promise.reject(error);
    }
);

export default api;
