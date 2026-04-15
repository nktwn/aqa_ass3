'use client';

import { useEffect, useState } from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import api from "@/lib/axios";

export default function Navbar() {
    const { user, isReady } = useAuth();
    const router = useRouter();
    const pathname = usePathname();
    const [query, setQuery] = useState("");
    const [suggestions, setSuggestions] = useState<string[]>([]);
    const [showSuggestions, setShowSuggestions] = useState(false);
    const [unreadCount, setUnreadCount] = useState(0);
    const [role, setRole] = useState<number | null>(null);

    useEffect(() => {
        const trimmed = query.trim();
        if (trimmed.length < 2) {
            setSuggestions([]);
            return;
        }

        const timer = setTimeout(async () => {
            try {
                const res = await api.get("/product/suggest", {
                    params: { q: trimmed, limit: 5 },
                });
                setSuggestions(res.data.suggestions ?? []);
            } catch (error) {
                console.error("Не удалось загрузить подсказки поиска:", error);
                setSuggestions([]);
            }
        }, 250);

        return () => clearTimeout(timer);
    }, [query]);

    useEffect(() => {
        if (!isReady || !user) {
            setUnreadCount(0);
            setRole(null);
            return;
        }

        const fetchUnreadCount = async () => {
            try {
                const [notificationRes, roleRes] = await Promise.all([
                    api.get("/notifications/unread-count"),
                    api.get("/user/role"),
                ]);
                setUnreadCount(notificationRes.data.count ?? 0);
                setRole(roleRes.data.role ?? null);
            } catch (error) {
                console.error("Не удалось загрузить количество уведомлений:", error);
            }
        };

        fetchUnreadCount();
    }, [isReady, user, pathname]);

    const submitSearch = (value: string) => {
        const trimmed = value.trim();
        if (!trimmed) {
            router.push("/catalog");
        } else {
            router.push(`/catalog?q=${encodeURIComponent(trimmed)}`);
        }
        setShowSuggestions(false);
    };

    return (
        <nav className="fixed top-0 left-0 w-full z-50 border-b border-white/40 bg-[var(--navbar-bg)] shadow-[0_12px_40px_rgba(15,23,42,0.08)] backdrop-blur-xl">
            <div className="mx-auto flex max-w-7xl flex-col gap-3 px-4 py-3">
                <div className="flex flex-wrap items-center gap-3 lg:flex-nowrap lg:justify-between">
                    <Link href="/" className="min-w-fit text-xl font-extrabold tracking-tight text-[var(--primary)]">
                        Zhan.Store
                    </Link>
                <div className="relative min-w-0 flex-1 lg:max-w-xl">
                    <form
                        onSubmit={(event) => {
                            event.preventDefault();
                            submitSearch(query);
                        }}
                        className="flex items-center gap-2"
                    >
                        <input
                            value={query}
                            onChange={(event) => setQuery(event.target.value)}
                            onFocus={() => setShowSuggestions(true)}
                            placeholder="Поиск товаров по названию"
                            className="w-full rounded-xl border border-white/60 bg-white/90 px-4 py-2.5 text-sm outline-none transition focus:border-[var(--primary)]"
                        />
                        <button
                            type="submit"
                            className="rounded-xl bg-[var(--primary)] px-4 py-2.5 text-sm font-semibold text-white transition hover:bg-[var(--primary-hover)]"
                        >
                            Найти
                        </button>
                    </form>

                    {showSuggestions && suggestions.length > 0 && (
                        <div className="absolute left-0 right-0 top-[calc(100%+8px)] rounded-2xl border border-white/70 bg-white shadow-lg">
                            {suggestions.map((suggestion) => (
                                <button
                                    key={suggestion}
                                    type="button"
                                    onMouseDown={() => submitSearch(suggestion)}
                                    className="block w-full px-4 py-3 text-left text-sm transition hover:bg-slate-50"
                                >
                                    {suggestion}
                                </button>
                            ))}
                        </div>
                    )}
                </div>
                </div>

                <div className="flex flex-wrap items-center gap-2 rounded-xl border border-white/50 bg-white/65 p-1.5 text-sm font-semibold text-[var(--foreground)] shadow-sm">
                    <Link href="/catalog" className={`rounded-xl px-4 py-2 transition ${pathname === "/catalog" ? "bg-[var(--primary-soft)] text-[var(--primary)]" : "hover:bg-white hover:text-[var(--primary)]"}`}>Каталог</Link>
                    <Link href="/favorites" className={`rounded-xl px-4 py-2 transition ${pathname === "/favorites" ? "bg-[var(--primary-soft)] text-[var(--primary)]" : "hover:bg-white hover:text-[var(--primary)]"}`}>Избранное</Link>
                    <Link href="/cart" className="rounded-xl px-4 py-2 transition hover:bg-white hover:text-[var(--primary)]">Корзина</Link>
                    {user ? (
                        <>
                            <Link href="/orders" className={`rounded-xl px-4 py-2 transition ${pathname === "/orders" ? "bg-[var(--primary-soft)] text-[var(--primary)]" : "hover:bg-white hover:text-[var(--primary)]"}`}>Заказы</Link>
                            <Link href="/notifications" className={`relative rounded-xl px-4 py-2 transition ${pathname === "/notifications" ? "bg-[var(--primary-soft)] text-[var(--primary)]" : "hover:bg-white hover:text-[var(--primary)]"}`}>
                                Уведомления
                                {unreadCount > 0 && (
                                    <span className="ml-2 rounded-full bg-red-500 px-2 py-0.5 text-xs text-white">
                                        {unreadCount}
                                    </span>
                                )}
                            </Link>
                            {role === 2 && <Link href="/admin/products" className="rounded-xl px-4 py-2 transition hover:bg-white hover:text-[var(--primary)]">Админ</Link>}
                            <Link href="/profile" className={`rounded-xl px-4 py-2 transition ${pathname === "/profile" ? "bg-[var(--primary-soft)] text-[var(--primary)]" : "hover:bg-white hover:text-[var(--primary)]"}`}>Профиль</Link>
                        </>
                    ) : (
                        <Link href="/login" className="rounded-xl px-4 py-2 transition hover:bg-white hover:text-[var(--primary)]">Войти</Link>
                    )}
                </div>

            </div>
        </nav>

    );
}
