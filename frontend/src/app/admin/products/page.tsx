'use client';

import { FormEvent, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import api from "@/lib/axios";
import axios from "axios";
import { Product } from "@/types";

type ProductForm = {
    name: string;
    image: string;
    gtin: string;
};

const initialForm: ProductForm = {
    name: "",
    image: "",
    gtin: "",
};

export default function AdminProductsPage() {
    const router = useRouter();
    const [products, setProducts] = useState<Product[]>([]);
    const [form, setForm] = useState<ProductForm>(initialForm);
    const [editingId, setEditingId] = useState<number | null>(null);
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState("");

    const loadProducts = async () => {
        const [roleRes, productRes] = await Promise.all([
            api.get("/user/role"),
            api.get("/product/list?limit=50&offset=0&sort_by=newest&sort_order=desc"),
        ]);

        if (roleRes.data.role !== 2) {
            window.location.href = "/";
            return;
        }

        setProducts(productRes.data.product_list || []);
    };

    useEffect(() => {
        const init = async () => {
            try {
                await loadProducts();
            } catch (e) {
                if (axios.isAxiosError(e) && e.response?.status === 401) {
                    window.location.href = "/login";
                    return;
                }
                setError("Не удалось загрузить админ-панель.");
            } finally {
                setLoading(false);
            }
        };

        init();
    }, [router]);

    const submit = async (event: FormEvent) => {
        event.preventDefault();
        setSubmitting(true);
        setError("");

        try {
            const payload = {
                name: form.name.trim(),
                image: form.image.trim(),
                gtin: Number(form.gtin || 0),
            };

            if (editingId) {
                await api.put(`/product/admin/${editingId}`, payload);
            } else {
                await api.post("/product/admin", payload);
            }

            setForm(initialForm);
            setEditingId(null);
            await loadProducts();
        } catch (e) {
            console.error("Ошибка сохранения товара:", e);
            setError("Не удалось сохранить товар.");
        } finally {
            setSubmitting(false);
        }
    };

    const edit = (product: Product) => {
        setEditingId(product.id);
        setForm({
            name: product.name,
            image: product.image,
            gtin: "",
        });
    };

    const remove = async (productId: number) => {
        try {
            await api.delete(`/product/admin/${productId}`);
            await loadProducts();
        } catch (e) {
            console.error("Ошибка удаления товара:", e);
            setError("Не удалось удалить товар.");
        }
    };

    if (loading) {
        return <p className="mt-20 text-center">Загрузка админ-панели...</p>;
    }

    return (
        <div className="mx-auto max-w-6xl space-y-8">
            <section className="rounded-[28px] border border-[var(--card-border)] bg-[var(--card-bg)] p-8 shadow-[0_24px_80px_rgba(15,23,42,0.08)]">
                <p className="text-sm font-semibold uppercase tracking-[0.24em] text-[var(--primary)]">Admin API</p>
                <h1 className="mt-3 text-3xl font-bold text-[var(--foreground)]">Управление товарами</h1>
                <p className="mt-3 text-sm text-slate-600">
                    Эта страница использует новые admin endpoints: создание, обновление и удаление товаров.
                </p>
            </section>

            <section className="grid gap-8 lg:grid-cols-[0.9fr_1.1fr]">
                <form onSubmit={submit} className="rounded-[24px] border border-[var(--card-border)] bg-white/90 p-6 shadow-sm space-y-4">
                    <h2 className="text-xl font-semibold text-[var(--foreground)]">
                        {editingId ? `Редактирование товара #${editingId}` : "Новый товар"}
                    </h2>

                    <input
                        value={form.name}
                        onChange={(event) => setForm((prev) => ({ ...prev, name: event.target.value }))}
                        placeholder="Название товара"
                        className="w-full rounded-2xl border border-[var(--card-border)] px-4 py-3 outline-none transition focus:border-[var(--primary)]"
                        required
                    />
                    <input
                        value={form.image}
                        onChange={(event) => setForm((prev) => ({ ...prev, image: event.target.value }))}
                        placeholder="Ссылка на изображение"
                        className="w-full rounded-2xl border border-[var(--card-border)] px-4 py-3 outline-none transition focus:border-[var(--primary)]"
                    />
                    <input
                        value={form.gtin}
                        onChange={(event) => setForm((prev) => ({ ...prev, gtin: event.target.value }))}
                        placeholder="GTIN"
                        inputMode="numeric"
                        className="w-full rounded-2xl border border-[var(--card-border)] px-4 py-3 outline-none transition focus:border-[var(--primary)]"
                    />

                    {error && <p className="text-sm text-red-600">{error}</p>}

                    <div className="flex flex-wrap gap-3">
                        <button type="submit" disabled={submitting} className="btn-primary">
                            {submitting ? "Сохраняем..." : editingId ? "Обновить товар" : "Создать товар"}
                        </button>
                        {editingId && (
                            <button
                                type="button"
                                onClick={() => {
                                    setEditingId(null);
                                    setForm(initialForm);
                                }}
                                className="btn-outline-primary"
                            >
                                Отменить редактирование
                            </button>
                        )}
                    </div>
                </form>

                <div className="rounded-[24px] border border-[var(--card-border)] bg-white/90 p-6 shadow-sm">
                    <h2 className="text-xl font-semibold text-[var(--foreground)]">Последние товары</h2>
                    <div className="mt-5 space-y-3">
                        {products.map((product) => (
                            <div key={product.id} className="flex flex-col gap-3 rounded-2xl border border-[var(--card-border)] bg-[var(--secondary-soft)] p-4 md:flex-row md:items-center md:justify-between">
                                <div className="min-w-0">
                                    <p className="font-semibold text-[var(--foreground)]">#{product.id} {product.name}</p>
                                    <p className="mt-1 text-sm text-slate-500">{product.image || "Без изображения"}</p>
                                </div>
                                <div className="flex gap-3">
                                    <button type="button" onClick={() => edit(product)} className="btn-outline-primary">
                                        Изменить
                                    </button>
                                    <button type="button" onClick={() => remove(product.id)} className="btn-danger">
                                        Удалить
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </section>
        </div>
    );
}
