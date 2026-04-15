'use client';

import { useEffect, useState } from 'react';
import {
    Bar, Line, Doughnut,
} from 'react-chartjs-2';
import {
    Chart as ChartJS, CategoryScale, LinearScale, BarElement, PointElement, LineElement, ArcElement, Tooltip, Legend,
} from 'chart.js';
import api from '@/lib/axios';

ChartJS.register(CategoryScale, LinearScale, BarElement, PointElement, LineElement, ArcElement, Tooltip, Legend);

export default function AnalyticsPage() {
    const [role, setRole] = useState<number | null>(null);

    useEffect(() => {
        const fetchRole = async () => {
            try {
                const res = await api.get("/user/role");
                setRole(res.data.role);
            } catch (e) {
                console.error("Ошибка определения роли:", e);
            }
        };

        fetchRole();
    }, []);

    if (role === null) return <p className="text-center mt-20">Загрузка роли...</p>;

    return (
        <div className="max-w-6xl mx-auto space-y-12 pt-20">
            <h1 className="text-3xl font-bold text-[var(--foreground)]">📊 Моя аналитика</h1>

            {role === 0 ? <ClientAnalytics /> : <SupplierAnalytics />}
        </div>
    );
}

function ClientAnalytics() {
    const data = {
        ordersCount: 28,
        completedCount: 22,
        totalSpent: 185000,
        orderStats: [2, 4, 3, 5, 6, 4, 4], // по дням недели
        statusBreakdown: { pending: 3, inProgress: 2, completed: 22, cancelled: 1 },
    };

    return (
        <>
            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
                <SummaryCard title="Мои заказы" value={data.ordersCount.toString()} />
                <SummaryCard title="Завершены" value={data.completedCount.toString()} />
                <SummaryCard title="Потрачено (₸)" value={data.totalSpent.toLocaleString()} />
            </div>

            <section className="bg-white rounded-xl shadow p-6">
                <h2 className="text-xl font-semibold mb-4">📈 Заказы по дням</h2>
                <Line
                    data={{
                        labels: ["Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"],
                        datasets: [{
                            label: 'Кол-во заказов',
                            data: data.orderStats,
                            borderColor: '#3b82f6',
                            backgroundColor: '#3b82f620',
                            fill: true,
                        }],
                    }}
                />
            </section>

            <section className="bg-white rounded-xl shadow p-6">
                <h2 className="text-xl font-semibold mb-4">📦 Статусы моих заказов</h2>

                <div className="flex flex-col md:flex-row items-center justify-between gap-6">
                    <Doughnut
                        data={{
                            labels: ['Ожидает', 'В процессе', 'Завершён', 'Отменён'],
                            datasets: [{
                                data: [
                                    data.statusBreakdown.pending,
                                    data.statusBreakdown.inProgress,
                                    data.statusBreakdown.completed,
                                    data.statusBreakdown.cancelled,
                                ],
                                backgroundColor: ['#facc15', '#3b82f6', '#10b981', '#ef4444'],
                            }],
                        }}
                        options={{ responsive: false, plugins: { legend: { display: true, position: 'right' } } }}
                        width={200}
                        height={200}
                    />

                    <div className="flex-1 space-y-4">
                        <div>
                            <p className="text-sm font-medium text-gray-700">Прогресс завершения заказов</p>
                            <div className="w-full bg-gray-200 h-4 rounded-full overflow-hidden mt-1">
                                <div
                                    className="bg-green-500 h-4"
                                    style={{
                                        width: `${Math.round((data.completedCount / data.ordersCount) * 100)}%`,
                                    }}
                                />
                            </div>
                            <p className="text-sm text-gray-600 mt-1">
                                {data.completedCount} из {data.ordersCount} завершены
                            </p>
                        </div>

                        {/* Можно добавить ещё карточку */}
                        <div className="text-sm text-gray-600">
                            Среднее количество заказов в неделю: <strong>~{Math.round(data.ordersCount / 4)}</strong>
                        </div>
                    </div>
                </div>
            </section>

        </>
    );
}

function SupplierAnalytics() {
    const data = {
        ordersReceived: 134,
        totalRevenue: 2980000,
        monthlyEarnings: [180000, 220000, 190000, 250000, 280000, 350000],
        topProducts: [
            { name: 'Товар A', sales: 84 },
            { name: 'Товар B', sales: 61 },
            { name: 'Товар C', sales: 45 },
        ],
    };

    return (
        <>
            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
                <SummaryCard title="Заказов получено" value={data.ordersReceived.toString()} />
                <SummaryCard title="Доход (₸)" value={data.totalRevenue.toLocaleString()} />
            </div>

            <section className="bg-white rounded-xl shadow p-6">
                <h2 className="text-xl font-semibold mb-4">💸 Доход по месяцам</h2>
                <Bar
                    data={{
                        labels: ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн'],
                        datasets: [{
                            label: '₸ Доход',
                            data: data.monthlyEarnings,
                            backgroundColor: '#14b8a6',
                        }],
                    }}
                />
            </section>

            <section className="bg-white rounded-xl shadow p-6">
                <h2 className="text-xl font-semibold mb-4">🔥 Топ-продажи</h2>
                <Bar
                    data={{
                        labels: data.topProducts.map(p => p.name),
                        datasets: [{
                            label: 'Продано (шт)',
                            data: data.topProducts.map(p => p.sales),
                            backgroundColor: '#f59e0b',
                        }],
                    }}
                    options={{ indexAxis: 'y' }}
                />
            </section>
        </>
    );
}

function SummaryCard({ title, value }: { title: string, value: string }) {
    return (
        <div className="bg-white rounded-xl shadow p-6 flex flex-col items-center justify-center text-center">
            <p className="text-sm text-gray-500">{title}</p>
            <h3 className="text-2xl font-bold text-[var(--foreground)]">{value}</h3>
        </div>
    );
}
