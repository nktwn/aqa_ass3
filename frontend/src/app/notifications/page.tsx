'use client';

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import api from "@/lib/axios";
import { Notification } from "@/types";
import { useAuth } from "@/hooks/useAuth";

export default function NotificationsPage() {
    const { user, isReady } = useAuth();
    const router = useRouter();
    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [loading, setLoading] = useState(true);

    const loadNotifications = async () => {
        try {
            const res = await api.get("/notifications");
            setNotifications(res.data.notifications ?? []);
        } catch (error) {
            console.error("Не удалось загрузить уведомления:", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (!isReady) {
            return;
        }

        if (!user) {
            router.replace("/login");
            return;
        }

        loadNotifications();
    }, [isReady, router, user]);

    const markRead = async (notificationId: number) => {
        try {
            await api.post(`/notifications/${notificationId}/read`);
            await loadNotifications();
        } catch (error) {
            console.error("Не удалось отметить уведомление как прочитанное:", error);
        }
    };

    const markAllRead = async () => {
        try {
            await api.post("/notifications/read-all");
            await loadNotifications();
        } catch (error) {
            console.error("Не удалось отметить все уведомления как прочитанные:", error);
        }
    };

    if (loading) {
        return <div className="mx-auto max-w-4xl py-12 text-center text-slate-500">Загрузка уведомлений...</div>;
    }

    return (
        <div className="mx-auto max-w-4xl space-y-6 py-8">
            <div className="flex items-center justify-between gap-4 flex-wrap">
                <div>
                    <h1 className="text-3xl font-bold text-[var(--foreground)]">Уведомления</h1>
                    <p className="mt-2 text-sm text-slate-500">Новые статусы заказов и важные события по вашим действиям.</p>
                </div>
                <button
                    onClick={markAllRead}
                    className="rounded-xl border border-[var(--card-border)] px-4 py-2 text-sm font-semibold transition hover:bg-slate-50"
                >
                    Отметить всё как прочитанное
                </button>
            </div>

            {notifications.length === 0 ? (
                <div className="rounded-2xl border border-dashed border-[var(--card-border)] bg-white p-10 text-center text-slate-500">
                    Пока уведомлений нет.
                </div>
            ) : (
                <div className="space-y-4">
                    {notifications.map((notification) => (
                        <div
                            key={notification.id}
                            className={`rounded-2xl border p-5 shadow-sm ${notification.is_read ? "bg-white" : "border-[var(--primary)] bg-[var(--primary)]/5"}`}
                        >
                            <div className="flex items-start justify-between gap-4">
                                <div className="space-y-2">
                                    <div className="flex items-center gap-3">
                                        <h2 className="text-lg font-semibold text-slate-900">{notification.title}</h2>
                                        {!notification.is_read && (
                                            <span className="rounded-full bg-[var(--primary)] px-2 py-0.5 text-xs font-semibold text-white">
                                                Новое
                                            </span>
                                        )}
                                    </div>
                                    <p className="text-sm text-slate-700">{notification.message}</p>
                                    <p className="text-xs text-slate-500">
                                        {new Date(notification.created_at).toLocaleString()}
                                    </p>
                                    {notification.order_id && (
                                        <Link href={`/orders/${notification.order_id}`} className="text-sm font-medium text-[var(--primary)] hover:underline">
                                            Перейти к заказу #{notification.order_id}
                                        </Link>
                                    )}
                                </div>
                                {!notification.is_read && (
                                    <button
                                        onClick={() => markRead(notification.id)}
                                        className="rounded-lg border border-[var(--card-border)] px-3 py-2 text-xs font-semibold transition hover:bg-white"
                                    >
                                        Прочитано
                                    </button>
                                )}
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
