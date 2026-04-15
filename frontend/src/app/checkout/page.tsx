'use client';

import { useEffect, useState } from "react";
import api from "@/lib/axios";
import axios from "axios";
import { Address } from "@/types";

interface CartProduct {
    id: number;
    name: string;
    image: string;
    price: number;
    quantity: number;
}

interface SupplierGroup {
    id: number;
    name: string;
    total_amount: number;
    delivery_fee: number;
    free_delivery_amount: number;
    order_amount: number;
    product_list: CartProduct[];
}

interface CartResponse {
    customer_id: number;
    total: number;
    suppliers: SupplierGroup[];
}

export default function CheckoutPage() {
    const [cartData, setCartData] = useState<CartResponse | null>(null);
    const [checkoutUrl, setCheckoutUrl] = useState<string | null>(null);
    const [paymentOrderId, setPaymentOrderId] = useState<string | null>(null);
    const [successMessage, setSuccessMessage] = useState(false);
    const [error, setError] = useState("");
    const [addresses, setAddresses] = useState<Address[]>([]);
    const [selectedAddressId, setSelectedAddressId] = useState<number | null>(null);
    const [addressForm, setAddressForm] = useState<Address>({ description: "", street: "" });
    const [isCreatingPayment, setIsCreatingPayment] = useState(false);
    const [isSavingAddress, setIsSavingAddress] = useState(false);

    useEffect(() => {
        if (sessionStorage.getItem("order_success") === "true") {
            setSuccessMessage(true);
            sessionStorage.removeItem("order_success");
        }

        const fetchData = async () => {
            try {
                const [cartRes, addressRes] = await Promise.all([
                    api.get("/cart/"),
                    api.get("/user/address"),
                ]);
                const nextCartData = cartRes.data as CartResponse;
                const supplierWithMinOrderIssue = nextCartData.suppliers.find(
                    (supplier) => supplier.total_amount < supplier.order_amount
                );

                setCartData(cartRes.data);
                setAddresses(addressRes.data.address_list || []);

                if ((addressRes.data.address_list || []).length > 0) {
                    setSelectedAddressId(addressRes.data.address_list[0].id);
                }

                if (supplierWithMinOrderIssue) {
                    setError(
                        `Минимальная сумма заказа у поставщика "${supplierWithMinOrderIssue.name}" составляет ${supplierWithMinOrderIssue.order_amount.toLocaleString()} ₸.`
                    );
                    return;
                }
            } catch (e: unknown) {
                if (axios.isAxiosError(e)) {
                    const responseData = e.response?.data;
                    const apiError =
                        responseData?.error ||
                        responseData?.err ||
                        (typeof responseData === "string" ? responseData : "");

                    if (apiError) {
                        setError(apiError);
                        return;
                    }

                    if (e.response?.status === 400) {
                        setError("Не удалось оформить заказ. Проверьте минимальную сумму заказа у поставщика.");
                        return;
                    }
                }

                console.error("Ошибка оформления:", e);
                setError("Не удалось создать ссылку на оплату. Проверьте корзину и попробуйте снова.");
            }
        };

        fetchData();
    }, []);

    const saveAddress = async () => {
        if (!addressForm.description.trim() || !addressForm.street.trim()) {
            setError("Заполните название адреса и полный адрес доставки.");
            return;
        }

        try {
            setError("");
            setIsSavingAddress(true);
            await api.post("/user/address", { address: addressForm });
            const addressRes = await api.get("/user/address");
            const nextAddresses = addressRes.data.address_list || [];
            setAddresses(nextAddresses);
            setAddressForm({ description: "", street: "" });
            if (nextAddresses.length > 0) {
                setSelectedAddressId(nextAddresses[nextAddresses.length - 1].id);
            }
        } catch (e) {
            console.error("Ошибка сохранения адреса:", e);
            setError("Не удалось сохранить адрес. Попробуйте снова.");
        } finally {
            setIsSavingAddress(false);
        }
    };

    const startPayment = async () => {
        if (!selectedAddressId) {
            setError("Выберите адрес доставки перед оплатой.");
            return;
        }

        try {
            setError("");
            setIsCreatingPayment(true);
            const checkoutRes = await api.post("/cart/checkout", { address_id: selectedAddressId });
            if (!checkoutRes.data?.checkout_url) {
                setError("Ссылка на оплату не была получена. Проверьте корзину и попробуйте снова.");
                return;
            }

            setCheckoutUrl(checkoutRes.data.checkout_url);
            setPaymentOrderId(checkoutRes.data.order_id ?? null);
        } catch (e: unknown) {
            if (axios.isAxiosError(e)) {
                const responseData = e.response?.data;
                const apiError =
                    responseData?.error ||
                    responseData?.err ||
                    (typeof responseData === "string" ? responseData : "");

                setError(apiError || "Не удалось создать платёжную ссылку.");
                return;
            }
            setError("Не удалось создать платёжную ссылку.");
        } finally {
            setIsCreatingPayment(false);
        }
    };

    const selectedAddress = addresses.find((address) => address.id === selectedAddressId);

    if (error) {
        return (
            <div className="text-center mt-20 text-red-600 font-medium">
                {error}
            </div>
        );
    }

    if (!cartData) {
        return <p className="text-center mt-20">Загрузка...</p>;
    }

    return (
        <div className="max-w-5xl mx-auto space-y-8">
            <div className="rounded-[28px] border border-[var(--card-border)] bg-[var(--card-bg)] p-8 shadow-[0_24px_80px_rgba(15,23,42,0.08)]">
                <p className="text-sm font-semibold uppercase tracking-[0.24em] text-[var(--primary)]">Checkout Flow</p>
                <h1 className="mt-3 text-3xl font-bold text-[var(--foreground)]">Подтверждение заказа и имитация оплаты</h1>
                <p className="mt-3 max-w-3xl text-sm text-slate-600">
                    Сначала выберите или сохраните адрес доставки, затем создайте платёжную ссылку и только после этого подтвердите оплату.
                </p>
            </div>

            {successMessage && (
                <div className="bg-green-100 border border-green-300 text-green-800 px-6 py-4 rounded-lg shadow text-center">
                    ✅ Заказ успешно создан и ожидает оплаты.
                </div>
            )}

            <div className="grid gap-8 xl:grid-cols-[1.25fr_0.75fr]">
                <div className="space-y-6">
                    <section className="rounded-[24px] border border-[var(--card-border)] bg-white/90 p-6 shadow-sm">
                        <div className="flex items-center justify-between gap-4">
                            <div>
                                <p className="text-xs font-semibold uppercase tracking-[0.24em] text-slate-500">Шаг 1</p>
                                <h2 className="text-xl font-semibold text-[var(--foreground)]">Адрес доставки</h2>
                            </div>
                            <span className="rounded-full bg-[var(--primary-soft)] px-3 py-1 text-xs font-semibold text-[var(--primary)]">
                                Сохраняется в профиле
                            </span>
                        </div>

                        <div className="mt-5 space-y-3">
                            {addresses.length === 0 ? (
                                <div className="rounded-2xl border border-dashed border-[var(--card-border)] bg-[var(--secondary-soft)] p-5 text-sm text-slate-600">
                                    У вас пока нет сохранённых адресов. Добавьте хотя бы один адрес, чтобы перейти к оплате.
                                </div>
                            ) : (
                                addresses.map((address) => (
                                    <label
                                        key={address.id}
                                        className={`flex cursor-pointer items-start gap-3 rounded-2xl border p-4 transition ${
                                            selectedAddressId === address.id
                                                ? "border-[var(--primary)] bg-[var(--primary-soft)]"
                                                : "border-[var(--card-border)] bg-white hover:border-[var(--primary)]/40"
                                        }`}
                                    >
                                        <input
                                            type="radio"
                                            name="address"
                                            checked={selectedAddressId === address.id}
                                            onChange={() => setSelectedAddressId(address.id ?? null)}
                                            className="mt-1"
                                        />
                                        <div>
                                            <p className="font-semibold text-[var(--foreground)]">{address.description}</p>
                                            <p className="mt-1 text-sm text-slate-600">{address.street}</p>
                                        </div>
                                    </label>
                                ))
                            )}
                        </div>

                        <div className="mt-6 grid gap-3 rounded-2xl border border-[var(--card-border)] bg-[var(--secondary-soft)] p-4">
                            <p className="text-sm font-semibold text-[var(--foreground)]">Новый адрес</p>
                            <input
                                value={addressForm.description}
                                onChange={(event) => setAddressForm((prev) => ({ ...prev, description: event.target.value }))}
                                placeholder="Например: Дом, офис, склад"
                                className="rounded-xl border border-[var(--card-border)] bg-white px-4 py-3 text-sm outline-none transition focus:border-[var(--primary)]"
                            />
                            <textarea
                                value={addressForm.street}
                                onChange={(event) => setAddressForm((prev) => ({ ...prev, street: event.target.value }))}
                                placeholder="Город, улица, дом, квартира, ориентир"
                                rows={3}
                                className="rounded-xl border border-[var(--card-border)] bg-white px-4 py-3 text-sm outline-none transition focus:border-[var(--primary)]"
                            />
                            <button
                                type="button"
                                onClick={saveAddress}
                                disabled={isSavingAddress}
                                className="btn-outline-primary w-full justify-center"
                            >
                                {isSavingAddress ? "Сохраняем адрес..." : "Сохранить адрес"}
                            </button>
                        </div>
                    </section>

                    <section className="rounded-[24px] border border-[var(--card-border)] bg-white/90 p-6 shadow-sm">
                        <div className="flex items-center justify-between gap-4">
                            <div>
                                <p className="text-xs font-semibold uppercase tracking-[0.24em] text-slate-500">Шаг 2</p>
                                <h2 className="text-xl font-semibold text-[var(--foreground)]">Состав заказа</h2>
                            </div>
                            <span className="text-sm text-slate-500">{cartData.suppliers.length} поставщиков</span>
                        </div>

                        <div className="mt-6 space-y-5">
                            {cartData.suppliers.map((supplier) => {
                                const needsMore = supplier.total_amount < supplier.free_delivery_amount;
                                const delivery = needsMore ? supplier.delivery_fee : 0;

                                return (
                                    <div
                                        key={supplier.id}
                                        className="rounded-2xl border border-[var(--card-border)] bg-[var(--secondary-soft)] p-5"
                                    >
                                        <h3 className="text-lg font-semibold text-[var(--foreground)]">{supplier.name}</h3>

                                        <ul className="mt-4 divide-y divide-[var(--card-border)]">
                                            {supplier.product_list.map((product) => (
                                                <li key={product.id} className="flex gap-4 py-4">
                                                    <img
                                                        src={product.image}
                                                        alt={product.name}
                                                        className="h-20 w-20 rounded-2xl border border-[var(--card-border)] object-cover"
                                                    />
                                                    <div className="flex-1">
                                                        <p className="font-medium">{product.name}</p>
                                                        <p className="text-sm text-slate-500">
                                                            {product.price.toLocaleString()} ₸ × {product.quantity} шт
                                                        </p>
                                                    </div>
                                                    <div className="text-right font-semibold">
                                                        {(product.price * product.quantity).toLocaleString()} ₸
                                                    </div>
                                                </li>
                                            ))}
                                        </ul>

                                        <div className="mt-4 border-t border-[var(--card-border)] pt-4 text-sm text-slate-600 space-y-1">
                                            {needsMore ? (
                                                <p>
                                                    Доставка: {delivery.toLocaleString()} ₸, до бесплатной ещё{" "}
                                                    {(supplier.free_delivery_amount - supplier.total_amount).toLocaleString()} ₸
                                                </p>
                                            ) : (
                                                <p className="font-medium text-emerald-700">Бесплатная доставка активна</p>
                                            )}
                                            <p className="font-semibold text-[var(--foreground)]">
                                                Сумма товаров: {supplier.total_amount.toLocaleString()} ₸
                                            </p>
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    </section>
                </div>

                <aside className="space-y-6">
                    <section className="rounded-[24px] border border-[var(--card-border)] bg-[var(--card-bg)] p-6 shadow-[0_18px_48px_rgba(2,6,23,0.08)]">
                        <p className="text-xs font-semibold uppercase tracking-[0.24em] text-slate-500">Шаг 3</p>
                        <h2 className="mt-2 text-xl font-semibold text-[var(--foreground)]">Подтверждение</h2>

                        <div className="mt-5 space-y-4 text-sm">
                            <div className="rounded-2xl bg-[var(--secondary-soft)] p-4">
                                <p className="text-slate-500">Адрес доставки</p>
                                {selectedAddress ? (
                                    <>
                                        <p className="mt-2 font-semibold text-[var(--foreground)]">{selectedAddress.description}</p>
                                        <p className="mt-1 text-slate-600">{selectedAddress.street}</p>
                                    </>
                                ) : (
                                    <p className="mt-2 text-slate-600">Адрес пока не выбран</p>
                                )}
                            </div>

                            <div className="rounded-2xl bg-[var(--secondary-soft)] p-4">
                                <p className="text-slate-500">Итого к оплате</p>
                                <p className="mt-2 text-3xl font-bold text-[var(--foreground)]">{cartData.total.toLocaleString()} ₸</p>
                            </div>

                            <button
                                type="button"
                                onClick={startPayment}
                                disabled={isCreatingPayment || !selectedAddressId}
                                className="btn-primary w-full justify-center"
                            >
                                {isCreatingPayment ? "Готовим оплату..." : "Создать платёжную ссылку"}
                            </button>

                            {checkoutUrl && (
                                <div className="rounded-3xl border border-[var(--primary)]/20 bg-[linear-gradient(135deg,rgba(15,118,110,0.08),rgba(217,119,6,0.10))] p-5">
                                    <p className="text-xs font-semibold uppercase tracking-[0.24em] text-[var(--primary)]">Имитация оплаты</p>
                                    <h3 className="mt-2 text-lg font-semibold text-[var(--foreground)]">Платёж готов к подтверждению</h3>
                                    <p className="mt-2 text-sm text-slate-600">
                                        Заказ будет создан только после перехода по тестовой ссылке оплаты.
                                    </p>
                                    {paymentOrderId && (
                                        <p className="mt-3 rounded-xl bg-white/80 px-3 py-2 text-xs text-slate-500">
                                            Payment order: {paymentOrderId}
                                        </p>
                                    )}
                                    <a
                                        href={checkoutUrl}
                                        target="_self"
                                        className="btn-primary mt-4 w-full justify-center"
                                    >
                                        Подтвердить тестовую оплату
                                    </a>
                                </div>
                            )}
                        </div>
                    </section>
                </aside>
            </div>
        </div>
    );
}
