'use client';

import { FormEvent, useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import api from "@/lib/axios";
import { Product } from "@/types";
import ProductCard from "@/components/ProductCard";

const LIMIT = 20;

export default function CatalogPage() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const [products, setProducts] = useState<Product[]>([]);
    const [page, setPage] = useState(0);
    const [hasNextPage, setHasNextPage] = useState(true);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(false);
    const [search, setSearch] = useState(searchParams.get("q") ?? "");
    const [minPrice, setMinPrice] = useState(searchParams.get("min_price") ?? "");
    const [maxPrice, setMaxPrice] = useState(searchParams.get("max_price") ?? "");
    const [sortBy, setSortBy] = useState(searchParams.get("sort_by") ?? "newest");
    const [sortOrder, setSortOrder] = useState(searchParams.get("sort_order") ?? "desc");
    const [appliedFilters, setAppliedFilters] = useState({
        search: searchParams.get("q") ?? "",
        minPrice: searchParams.get("min_price") ?? "",
        maxPrice: searchParams.get("max_price") ?? "",
        sortBy: searchParams.get("sort_by") ?? "newest",
        sortOrder: searchParams.get("sort_order") ?? "desc",
    });

    useEffect(() => {
        const fetchProducts = async () => {
            try {
                setLoading(true);
                const res = await api.get("/product/list", {
                    params: {
                        limit: LIMIT,
                        offset: page * LIMIT,
                        search: appliedFilters.search || undefined,
                        min_price: appliedFilters.minPrice || undefined,
                        max_price: appliedFilters.maxPrice || undefined,
                        sort_by: appliedFilters.sortBy,
                        sort_order: appliedFilters.sortOrder,
                    },
                });
                setProducts(res.data.product_list);
                setTotal(res.data.total ?? 0);
                setHasNextPage((page + 1) * LIMIT < (res.data.total ?? 0));
            } catch (e) {
                console.error("Не удалось загрузить продукты:", e);
            } finally {
                setLoading(false);
            }
        };

        fetchProducts();
    }, [page, appliedFilters]);

    const totalPages = Math.max(1, Math.ceil(total / LIMIT));

    const applyFilters = (event?: FormEvent) => {
        event?.preventDefault();

        const nextFilters = {
            search: search.trim(),
            minPrice: minPrice.trim(),
            maxPrice: maxPrice.trim(),
            sortBy,
            sortOrder,
        };

        setAppliedFilters(nextFilters);
        setPage(0);

        const params = new URLSearchParams();
        if (nextFilters.search) params.set("q", nextFilters.search);
        if (nextFilters.minPrice) params.set("min_price", nextFilters.minPrice);
        if (nextFilters.maxPrice) params.set("max_price", nextFilters.maxPrice);
        if (nextFilters.sortBy !== "newest") params.set("sort_by", nextFilters.sortBy);
        if (nextFilters.sortOrder !== "desc") params.set("sort_order", nextFilters.sortOrder);

        const query = params.toString();
        router.replace(query ? `/catalog?${query}` : "/catalog");
    };

    const resetFilters = () => {
        setSearch("");
        setMinPrice("");
        setMaxPrice("");
        setSortBy("newest");
        setSortOrder("desc");
        setAppliedFilters({
            search: "",
            minPrice: "",
            maxPrice: "",
            sortBy: "newest",
            sortOrder: "desc",
        });
        setPage(0);
        router.replace("/catalog");
    };

    // ближайшие страницы вокруг текущей
    const getNearbyPages = () => {
        const pages = [];
        for (let i = Math.max(0, page - 2); i <= Math.min(totalPages - 1, page + 2); i++) {
            pages.push(i);
        }
        return pages;
    };

    return (
        <div className="w-[92%] mx-auto px-4 py-8 space-y-10">
            <section className="rounded-[28px] border border-[var(--card-border)] shadow p-6 bg-[linear-gradient(180deg,rgba(255,255,255,0.9),rgba(255,247,237,0.86))] w-full">
                <div className="mb-6 flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
                    <div>
                        <p className="text-sm font-semibold uppercase tracking-[0.24em] text-[var(--primary)]">Curated Catalog</p>
                        <h1 className="text-3xl font-bold text-[var(--foreground)]">Каталог товаров</h1>
                    </div>
                    <div className="rounded-full bg-[var(--primary-soft)] px-4 py-2 text-sm font-semibold text-[var(--primary)]">
                        Найдено: {total}
                    </div>
                </div>
                <form onSubmit={applyFilters} className="mb-6 grid gap-4 rounded-[24px] border border-[var(--card-border)] bg-white/80 p-5 md:grid-cols-5">
                    <input
                        value={search}
                        onChange={(event) => setSearch(event.target.value)}
                        placeholder="Поиск по названию"
                        className="rounded-xl border border-[var(--card-border)] px-4 py-3 text-sm outline-none transition focus:border-[var(--primary)] md:col-span-2"
                    />
                    <input
                        value={minPrice}
                        onChange={(event) => setMinPrice(event.target.value)}
                        placeholder="Цена от"
                        inputMode="numeric"
                        className="rounded-xl border border-[var(--card-border)] px-4 py-3 text-sm outline-none transition focus:border-[var(--primary)]"
                    />
                    <input
                        value={maxPrice}
                        onChange={(event) => setMaxPrice(event.target.value)}
                        placeholder="Цена до"
                        inputMode="numeric"
                        className="rounded-xl border border-[var(--card-border)] px-4 py-3 text-sm outline-none transition focus:border-[var(--primary)]"
                    />
                    <div className="grid grid-cols-2 gap-3 md:col-span-5">
                        <select
                            value={sortBy}
                            onChange={(event) => setSortBy(event.target.value)}
                            className="rounded-xl border border-[var(--card-border)] px-4 py-3 text-sm outline-none transition focus:border-[var(--primary)]"
                        >
                            <option value="newest">Сначала новые</option>
                            <option value="price">По цене</option>
                            <option value="name">По названию</option>
                        </select>
                        <select
                            value={sortOrder}
                            onChange={(event) => setSortOrder(event.target.value)}
                            className="rounded-xl border border-[var(--card-border)] px-4 py-3 text-sm outline-none transition focus:border-[var(--primary)]"
                        >
                            <option value="desc">По убыванию</option>
                            <option value="asc">По возрастанию</option>
                        </select>
                    </div>
                    <div className="flex flex-wrap gap-3 md:col-span-5">
                        <button
                            type="submit"
                            className="rounded-xl bg-[var(--primary)] px-5 py-3 text-sm font-semibold text-white transition hover:bg-[var(--primary-hover)]"
                        >
                            Применить фильтры
                        </button>
                        <button
                            type="button"
                            onClick={resetFilters}
                            className="rounded-xl border border-[var(--card-border)] px-5 py-3 text-sm font-semibold transition hover:bg-slate-50"
                        >
                            Сбросить
                        </button>
                    </div>
                </form>

                {loading ? (
                    <div className="py-16 text-center text-slate-500">Загрузка товаров...</div>
                ) : products.length === 0 ? (
                    <div className="rounded-2xl border border-dashed border-[var(--card-border)] bg-white/70 py-16 text-center text-slate-500">
                        По текущим фильтрам товары не найдены.
                    </div>
                ) : (
                    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-5 gap-6">
                    {products.map((product) => (
                        <ProductCard key={product.id} product={product} />
                    ))}
                    </div>
                )}

                {/* Pagination */}
                <div className="flex justify-center mt-10">
                    <div className="inline-flex gap-2 flex-wrap justify-center items-center">
                        {/* ← */}
                        <button
                            onClick={() => setPage((p) => Math.max(0, p - 1))}
                            disabled={page === 0}
                            className="px-4 py-2 rounded-lg border text-sm font-medium transition disabled:opacity-50
                bg-white text-[#14b8a6] border-[#14b8a6] hover:bg-[#14b8a6] hover:text-white"
                        >
                            ←
                        </button>

                        {/* numbered pages */}
                        {getNearbyPages().map((p) => (
                            <button
                                key={p}
                                onClick={() => setPage(p)}
                                className={`px-4 py-2 rounded-lg text-sm font-medium transition border ${
                                    p === page
                                        ? "bg-[#14b8a6] text-white border-[#14b8a6]"
                                        : "bg-white text-[#14b8a6] border-[#14b8a6] hover:bg-[#14b8a6] hover:text-white"
                                }`}
                            >
                                {p + 1}
                            </button>
                        ))}

                        {/* → */}
                        <button
                            onClick={() => setPage((p) => p + 1)}
                            disabled={!hasNextPage}
                            className="px-4 py-2 rounded-lg border text-sm font-medium transition disabled:opacity-50
                bg-white text-[#14b8a6] border-[#14b8a6] hover:bg-[#14b8a6] hover:text-white"
                        >
                            →
                        </button>
                    </div>
                </div>
            </section>
        </div>
    );
}
