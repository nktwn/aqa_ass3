'use client';

import { createContext, useContext, useState } from "react";
import api from "@/lib/axios";
import { CartItem } from "@/types";
import { useAuth } from "@/hooks/useAuth";

interface CartContextType {
    cart: CartItem[];
    addToCart: (item: CartItem) => Promise<void>;
    clearCart: () => void;
    checkout: () => Promise<void>;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export const CartProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [cart, setCart] = useState<CartItem[]>([]);
    const { user, isReady } = useAuth();

    const addToCart = async (item: CartItem) => {
        if (!isReady || !user) {
            if (typeof window !== "undefined") {
                window.location.href = "/login";
            }
            return;
        }

        await api.post("/cart/add", {
            product_id: item.productId,
            supplier_id: item.supplierId,
            quantity: item.quantity,
        });

        setCart((prev) => {
            const exists = prev.find(
                (p) => p.productId === item.productId && p.supplierId === item.supplierId
            );
            if (exists) {
                return prev.map((p) =>
                    p.productId === item.productId && p.supplierId === item.supplierId
                        ? { ...p, quantity: p.quantity + item.quantity }
                        : p
                );
            }
            return [...prev, item];
        });
    };

    const clearCart = () => {
        setCart([]);
    };

    const checkout = async () => {
        if (!isReady || !user) {
            if (typeof window !== "undefined") {
                window.location.href = "/login";
            }
            return;
        }

        await api.post("/cart/checkout");
        clearCart();
    };

    return (
        <CartContext.Provider value={{ cart, addToCart, clearCart, checkout }}>
            {children}
        </CartContext.Provider>
    );
};

export const useCart = () => {
    const context = useContext(CartContext);
    if (!context) throw new Error("useCart must be used within CartProvider");
    return context;
};
