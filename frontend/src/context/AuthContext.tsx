'use client';

import { createContext, useContext, useState, useEffect } from "react";
import api from "@/lib/axios";
import { User } from "@/types";

interface AuthContextType {
    user: User | null;
    isReady: boolean;
    login: (phone: string, password: string) => Promise<void>;
    logout: () => Promise<void>;
    register: (data: { name: string; phone_number: string; password: string; confirm_password: string }) => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [isReady, setIsReady] = useState(false);

    useEffect(() => {
        if (typeof window === "undefined") {
            return;
        }

        const syncAuthState = () => {
            const token = localStorage.getItem("access_token");
            const expiresAt = localStorage.getItem("expires_at");

            if (token && expiresAt) {
                const now = Date.now();
                if (now < parseInt(expiresAt, 10)) {
                    setUser({ id: 0, name: "Пользователь", email: "" });
                    setIsReady(true);
                    return;
                }
            }

            localStorage.removeItem("access_token");
            localStorage.removeItem("refresh_token");
            localStorage.removeItem("expires_at");
            setUser(null);
            setIsReady(true);
        };

        syncAuthState();
        window.addEventListener("auth:unauthorized", syncAuthState);

        return () => {
            window.removeEventListener("auth:unauthorized", syncAuthState);
        };
    }, []);

    const login = async (phone_number: string, password: string) => {
        const res = await api.post("/auth/login", { phone_number, password });

        const expiresAt = Date.now() + 6 * 60 * 60 * 1000; // 6 часов

        localStorage.setItem("access_token", res.data.access_token);
        localStorage.setItem("refresh_token", res.data.refresh_token);
        localStorage.setItem("expires_at", expiresAt.toString());

        setUser({ id: 0, name: "Пользователь", email: "" });
    };

    const register = async (data: { name: string; phone_number: string; password: string; confirm_password: string }) => {
        await api.post("/auth/register", data);
        const loginRes = await api.post("/auth/login", {
            phone_number: data.phone_number,
            password: data.password,
        });

        const expiresAt = Date.now() + 6 * 60 * 60 * 1000;
        localStorage.setItem("access_token", loginRes.data.access_token);
        localStorage.setItem("refresh_token", loginRes.data.refresh_token);
        localStorage.setItem("expires_at", expiresAt.toString());
        setUser({ id: 0, name: "Пользователь", email: "" });
    };

    const logout = async () => {
        setUser(null);
        localStorage.removeItem("access_token");
        localStorage.removeItem("refresh_token");
        localStorage.removeItem("expires_at");
    };

    return (
        <AuthContext.Provider value={{ user, isReady, login, logout, register }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) throw new Error("useAuth must be used within AuthProvider");
    return context;
};
