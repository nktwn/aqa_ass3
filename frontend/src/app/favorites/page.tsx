'use client';

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import ProductCard from "@/components/ProductCard";
import api from "@/lib/axios";
import { Product } from "@/types";
import { useAuth } from "@/hooks/useAuth";

export default function FavoritesPage() {
    const { user, isReady } = useAuth();
    const router = useRouter();
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (!isReady) {
            return;
        }

        if (!user) {
            router.replace("/login");
            return;
        }

        const fetchFavorites = async () => {
            try {
                const res = await api.get("/product/favorites");
                setProducts(res.data.product_list ?? []);
            } catch (error) {
                console.error("Не удалось загрузить избранное:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchFavorites();
    }, [isReady, router, user]);

    if (loading) {
        return <div className="mx-auto max-w-6xl py-12 text-center text-slate-500">Загрузка избранного...</div>;
    }

    return (
        <div className="mx-auto max-w-6xl space-y-6 py-8">
            <div>
                <h1 className="text-3xl font-bold text-[var(--foreground)]">Избранное</h1>
                <p className="mt-2 text-sm text-slate-500">Товары, которые вы сохранили для быстрого возврата.</p>
            </div>

            {products.length === 0 ? (
                <div className="rounded-2xl border border-dashed border-[var(--card-border)] bg-white p-10 text-center text-slate-500">
                    В избранном пока ничего нет.
                </div>
            ) : (
                <div className="grid grid-cols-2 gap-6 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-5">
                    {products.map((product) => (
                        <ProductCard key={product.id} product={product} />
                    ))}
                </div>
            )}
        </div>
    );
}
