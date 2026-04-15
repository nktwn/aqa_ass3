'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import api from '@/lib/axios';
import { getStatusBadge } from '@/components/StatusBadge';
import SignatureModal from '@/components/SignatureModal';

interface OrderProduct {
    id: number;
    name: string;
    image: string;
    price: number;
    quantity: number;
}

interface Order {
    id: number;
    status: string;
    order_date: string;
    delivery_address: string;
    delivery_comment: string;
    payment_status: string;
    supplier: {
        id: number;
        name: string;
    };
    product_list: OrderProduct[];
}

interface Contract {
    id: number;
    content: string;
    status: number;
    supplier_signature?: string;
    customer_signature?: string;
}

export default function OrderDetailPage() {
    const { id } = useParams();
    const [order, setOrder] = useState<Order | null>(null);
    const [loading, setLoading] = useState(true);
    const [role, setRole] = useState<number | null>(null);
    const [customerName, setCustomerName] = useState<string>('');
    const [updating, setUpdating] = useState(false);
    const [contracts, setContracts] = useState<Contract[]>([]);
    const [signingContract, setSigningContract] = useState<{ id: number, role: 'supplier' | 'user' } | null>(null);

    useEffect(() => {
        const fetchOrderAndRole = async () => {
            try {
                const [orderRes, roleRes, contractRes] = await Promise.all([
                    api.get(`/order/${id}`),
                    api.get(`/user/role`),
                    api.get(`/contract`),
                ]);
                setOrder(orderRes.data);
                setRole(roleRes.data.role);
                setContracts(contractRes.data);

                if (roleRes.data.role === 0) {
                    const profileRes = await api.get("/user/profile");
                    setCustomerName(profileRes.data.user.name);
                }
            } catch (err) {
                console.error('Ошибка загрузки данных:', err);
            } finally {
                setLoading(false);
            }
        };

        fetchOrderAndRole();
    }, [id]);

    const refreshOrder = async () => {
        if (!order) return;
        const res = await api.get(`/order/${order.id}`);
        setOrder(res.data);
    };

    const refreshContracts = async () => {
        const res = await api.get(`/contract`);
        setContracts(res.data);
    };

    const cancelOrder = async () => {
        if (!order) return;
        setUpdating(true);
        try {
            await api.post("/order/cancel", { order_id: order.id });
            await refreshOrder();
        } catch (err) {
            console.error("Ошибка при отмене заказа:", err);
        } finally {
            setUpdating(false);
        }
    };

    const updateStatus = async (newStatus: 'In Progress' | 'Completed' | 'Cancelled') => {
        if (!order) return;

        const statusMap: Record<string, number> = {
            'In Progress': 2,
            'Completed': 3,
            'Cancelled': 4
        };

        setUpdating(true);
        try {
            await api.post("/order/status", {
                order_id: order.id,
                new_status_id: statusMap[newStatus],
            });
            await refreshOrder();
        } catch (err) {
            console.error("Ошибка при смене статуса:", err);
        } finally {
            setUpdating(false);
        }
    };

    const confirmSignature = async (code: string) => {
        if (!signingContract) return;
        setUpdating(true);
        try {
            void code;
            await api.post("/contract/sign", {
                contract_id: signingContract.id,
                signature: signingContract.role,
            });
            await refreshContracts();
        } catch (err) {
            console.error("Ошибка при подписании контракта:", err);
            alert("Не удалось подписать документ");
        } finally {
            setUpdating(false);
            setSigningContract(null);
        }
    };

    if (loading) return <p className="text-center mt-20">Загрузка заказа...</p>;
    if (!order) return <p className="text-center mt-20 text-red-500">Заказ не найден.</p>;

    const relatedContracts = contracts.filter(c => c.content.includes(`#${order.id}`));
    const totalAmount = order.product_list.reduce((sum, item) => sum + item.price * item.quantity, 0);

    return (
        <div className="max-w-3xl mx-auto space-y-8">
            <div className="flex justify-between items-start flex-wrap gap-4">
                <div>
                    <h1 className="text-3xl font-bold mb-2">📦 Заказ #{order.id}</h1>
                    <div className="text-sm text-gray-600 space-y-1">
                        <p className="flex items-center gap-2">Статус: {getStatusBadge(order.status)}</p>
                        <p>Поставщик: <span className="font-medium">{order.supplier.name}</span></p>
                        <p>Дата заказа: {new Date(order.order_date).toLocaleString()}</p>
                        <p>Оплата: <span className="font-medium">{order.payment_status || 'pending'}</span></p>
                        <p>Адрес доставки: <span className="font-medium">{order.delivery_address}</span></p>
                        {order.delivery_comment && <p>Комментарий: <span className="font-medium">{order.delivery_comment}</span></p>}
                    </div>
                </div>

                {role === 0 && order.status === "Pending" && (
                    <button
                        onClick={cancelOrder}
                        disabled={updating}
                        className="btn-danger text-sm"
                    >
                        {updating ? "Отмена..." : "❌ Отменить заказ"}
                    </button>
                )}

                {role === 1 && order.status === 'Pending' && (
                    <button
                        onClick={() => updateStatus('In Progress')}
                        disabled={updating}
                        className="btn-outline-primary text-sm"
                    >
                        🔄 Принять в работу
                    </button>
                )}

                {role === 1 && order.status === 'In Progress' && (
                    <div className="flex gap-2">
                        <button
                            onClick={() => updateStatus('Completed')}
                            disabled={updating}
                            className="btn-primary text-sm"
                        >
                            ✅ Завершить
                        </button>
                        <button
                            onClick={() => updateStatus('Cancelled')}
                            disabled={updating}
                            className="btn-danger text-sm"
                        >
                            ❌ Отменить
                        </button>
                    </div>
                )}
            </div>

            <ul className="divide-y border rounded-xl bg-white shadow">
                {order.product_list.map((item) => (
                    <li key={item.id} className="flex gap-4 p-4">
                        <img
                            src={item.image}
                            alt={item.name}
                            className="w-20 h-20 object-cover border rounded"
                        />
                        <div className="flex-1">
                            <p className="font-medium">{item.name}</p>
                            <p className="text-sm text-gray-500">
                                {item.price.toLocaleString()} ₸ × {item.quantity} шт
                            </p>
                        </div>
                        <div className="text-right font-semibold">
                            {(item.price * item.quantity).toLocaleString()} ₸
                        </div>
                    </li>
                ))}
            </ul>

            {relatedContracts.length > 0 && (
                <div className="space-y-6">
                    <h2 className="text-xl font-semibold text-[var(--foreground)]">📃 Документы</h2>

                    {relatedContracts.map((contract) => {
                        const supplierSigned = !!contract.supplier_signature;
                        const customerSigned = !!contract.customer_signature;

                        const canSupplierSign =
                            role === 1 && !supplierSigned && order.status === 'In Progress';
                        const canCustomerSign =
                            role === 0 && supplierSigned && !customerSigned;

                        return (
                            <div
                                key={contract.id}
                                className="border border-gray-300 rounded-xl bg-white shadow p-8 font-serif space-y-6 leading-relaxed"
                            >
                                <h3 className="text-2xl font-bold text-center underline mb-4">
                                    Акт приёма-передачи товара
                                </h3>

                                <p>
                                    Настоящий акт составлен между поставщиком <strong>{order.supplier.name}</strong> и покупателем <strong>{customerName || 'покупатель'}</strong> по заказу №{order.id} от {new Date(order.order_date).toLocaleDateString()}.
                                </p>

                                <p>
                                    Я, <strong>{order.supplier.name}</strong>, подтверждаю, что передал(а) покупателю следующий товар:
                                </p>

                                <ul className="list-disc pl-6">
                                    {order.product_list.map((item) => (
                                        <li key={item.id}>
                                            {item.name} — {item.quantity} шт по {item.price.toLocaleString()} ₸
                                        </li>
                                    ))}
                                </ul>

                                <p>
                                    Общая сумма поставки составляет: <strong>{totalAmount.toLocaleString()} ₸</strong>.
                                </p>

                                <p>
                                    Я, <strong>{customerName || 'покупатель'}</strong>, подтверждаю, что получил(а) товар в полном объёме и надлежащего качества.
                                </p>

                                <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 pt-6">
                                    <div className="border-t pt-2">
                                        <p className="text-sm font-medium">Подпись поставщика:</p>
                                        <p className="mt-1 italic">{supplierSigned ? '✅ Подписано' : '—'}</p>
                                    </div>
                                    <div className="border-t pt-2">
                                        <p className="text-sm font-medium">Подпись клиента:</p>
                                        <p className="mt-1 italic">{customerSigned ? '✅ Подписано' : '—'}</p>
                                    </div>
                                </div>

                                {(canSupplierSign || canCustomerSign) && (
                                    <div className="pt-4">
                                        <button
                                            onClick={() =>
                                                setSigningContract({
                                                    id: contract.id,
                                                    role: role === 1 ? 'supplier' : 'user',
                                                })
                                            }
                                            disabled={updating}
                                            className="btn-primary w-full sm:w-auto"
                                        >
                                            ✍️ Подписать акт
                                        </button>
                                    </div>
                                )}
                            </div>
                        );
                    })}
                </div>
            )}

            <SignatureModal
                isOpen={!!signingContract}
                onClose={() => setSigningContract(null)}
                onConfirm={confirmSignature}
            />
        </div>
    );
}
