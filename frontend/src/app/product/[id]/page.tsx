'use client';

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import api from "@/lib/axios";
import { Product, Review, Supplier } from "@/types";
import { useCart } from "@/context/CartContext";
import { useAuth } from "@/hooks/useAuth";

interface SupplierOffer {
    price: number;
    sell_amount: number;
    supplier: Supplier;
}

interface ApiResponse {
    product: {
        product: Product;
        suppliers: SupplierOffer[];
        reviews: Review[];
    };
}

export default function ProductPage() {
    const { id } = useParams<{ id: string }>();
    const router = useRouter();
    const { user, isReady } = useAuth();
    const [product, setProduct] = useState<Product | null>(null);
    const [suppliers, setSuppliers] = useState<SupplierOffer[]>([]);
    const [reviews, setReviews] = useState<Review[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [isFavorite, setIsFavorite] = useState(false);
    const [favoriteLoading, setFavoriteLoading] = useState(false);
    const [reviewRating, setReviewRating] = useState(5);
    const [reviewComment, setReviewComment] = useState("");
    const [reviewSubmitting, setReviewSubmitting] = useState(false);
    const { addToCart } = useCart();

    const [selectedSupplier, setSelectedSupplier] = useState<number | null>(null);
    const [quantities, setQuantities] = useState<Record<number, number>>({});

    useEffect(() => {
        const fetchProduct = async () => {
            try {
                const res = await api.get<ApiResponse>(`/product/${id}`);
                setProduct(res.data.product.product);
                setSuppliers(res.data.product.suppliers);
                setReviews(res.data.product.reviews || []);
            } catch (err) {
                console.error("Ошибка загрузки продукта:", err);
                setError("Не удалось загрузить продукт");
            } finally {
                setLoading(false);
            }
        };

        if (id) fetchProduct();
    }, [id]);

    useEffect(() => {
        if (!isReady || !user || !id) {
            setIsFavorite(false);
            return;
        }

        const fetchFavoriteStatus = async () => {
            try {
                const res = await api.get(`/product/${id}/favorite`);
                setIsFavorite(Boolean(res.data.is_favorite));
            } catch (err) {
                console.error("Не удалось получить статус избранного:", err);
            }
        };

        fetchFavoriteStatus();
    }, [id, isReady, user]);

    const handleSelectQuantity = (supplierId: number) => {
        setSelectedSupplier((prev) => (prev === supplierId ? null : supplierId));
        setQuantities((prev) => ({
            ...prev,
            [supplierId]: prev[supplierId] ?? 1,
        }));
    };

    const handleChangeQuantity = (supplierId: number, delta: number) => {
        setQuantities((prev) => {
            const newQty = Math.max(1, (prev[supplierId] ?? 1) + delta);
            return { ...prev, [supplierId]: newQty };
        });
    };

    const handleAddToCart = async (supplierId: number) => {
        const quantity = quantities[supplierId] ?? 1;

        await addToCart({
            productId: product!.id,
            supplierId,
            quantity,
        });

        // сброс мини-интерфейса
        setSelectedSupplier(null);
        setQuantities((prev) => ({ ...prev, [supplierId]: 1 }));
    };

    const reloadProduct = async () => {
        const res = await api.get<ApiResponse>(`/product/${id}`);
        setProduct(res.data.product.product);
        setSuppliers(res.data.product.suppliers);
        setReviews(res.data.product.reviews || []);
    };

    const toggleFavorite = async () => {
        if (!user) {
            router.push("/login");
            return;
        }

        setFavoriteLoading(true);
        try {
            if (isFavorite) {
                await api.delete(`/product/${id}/favorite`);
                setIsFavorite(false);
            } else {
                await api.post(`/product/${id}/favorite`);
                setIsFavorite(true);
            }
        } catch (err) {
            console.error("Не удалось обновить избранное:", err);
        } finally {
            setFavoriteLoading(false);
        }
    };

    const submitReview = async () => {
        if (!user) {
            router.push("/login");
            return;
        }

        setReviewSubmitting(true);
        try {
            await api.post(`/product/${id}/reviews`, {
                rating: reviewRating,
                comment: reviewComment,
            });
            setReviewComment("");
            await reloadProduct();
        } catch (err) {
            console.error("Не удалось сохранить отзыв:", err);
        } finally {
            setReviewSubmitting(false);
        }
    };

    if (loading) return <p>Загрузка...</p>;
    if (error) return <p className="text-red-600">{error}</p>;
    if (!product) return <p>Продукт не найден</p>;

    return (
        <div className="max-w-4xl mx-auto space-y-8 pt-6">
            <Link href="/" className="text-blue-600 hover:underline">← Назад</Link>

            {/* Основной блок */}
            <div className="bg-white border rounded-xl shadow p-6 flex flex-col sm:flex-row gap-6">
                <img
                    src={product.image}
                    alt={product.name}
                    className="w-full sm:w-64 h-64 object-cover rounded-lg border"
                />
                <div className="flex-1 space-y-4">
                    <h1 className="text-2xl font-bold text-[var(--foreground)]">{product.name}</h1>

                    {product.lowest_product_supplier?.price > 0 && (
                        <p className="text-lg text-[var(--primary)] font-semibold">
                            от {product.lowest_product_supplier.price.toLocaleString()} ₸
                        </p>
                    )}

                    <div className="flex items-center gap-3">
                        <p className="text-sm text-amber-600 font-medium">
                            {"★".repeat(Math.max(1, Math.round(product.rating_average || 0)))} {product.rating_average.toFixed(1)}
                        </p>
                        <p className="text-sm text-slate-500">{product.rating_count} отзывов</p>
                    </div>

                    <p className="text-gray-600 text-sm">
                        Доступные предложения от поставщиков ниже
                    </p>

                    <button
                        onClick={toggleFavorite}
                        disabled={favoriteLoading}
                        className={`w-fit rounded-lg px-4 py-2 text-sm font-medium transition ${isFavorite ? "bg-rose-500 text-white" : "border border-slate-300 text-slate-700 hover:bg-slate-50"}`}
                    >
                        {favoriteLoading ? "Сохраняем..." : isFavorite ? "♥ В избранном" : "♡ Добавить в избранное"}
                    </button>
                </div>
            </div>

            {/* Поставщики */}
            <div className="space-y-6">
                <h2 className="text-xl font-semibold text-[var(--foreground)]">Поставщики</h2>

                {suppliers.length === 0 ? (
                    <p className="text-gray-500">Нет доступных поставщиков для этого товара.</p>
                ) : (
                    <ul className="space-y-4">
                        {suppliers.map((offer) => {
                            const isSelected = selectedSupplier === offer.supplier.id;
                            const quantity = quantities[offer.supplier.id] ?? 1;

                            return (
                                <li
                                    key={offer.supplier.id}
                                    className="border p-6 rounded-xl bg-white shadow flex flex-col sm:flex-row justify-between items-start sm:items-center gap-6"
                                >
                                    <div className="flex-1 space-y-2">
                                        <p className="text-lg font-semibold text-[var(--foreground)]">
                                            🏪 {offer.supplier.name}
                                        </p>
                                        <p className="text-gray-600 text-sm">
                                            Минимальный заказ: {offer.supplier.order_amount.toLocaleString()} ₸
                                        </p>
                                        <p className="text-gray-600 text-sm">
                                            Бесплатная доставка от: {offer.supplier.free_delivery_amount.toLocaleString()} ₸
                                        </p>
                                        <p className="text-gray-600 text-sm">
                                            Стоимость доставки: {offer.supplier.delivery_fee.toLocaleString()} ₸
                                        </p>
                                    </div>

                                    <div className="flex flex-col items-end gap-3 min-w-[160px] w-full sm:w-auto">
                                        <p className="text-xl font-bold text-[var(--primary)]">
                                            {offer.price.toLocaleString()} ₸
                                        </p>

                                        {!isSelected ? (
                                            <button
                                                onClick={() => handleSelectQuantity(offer.supplier.id)}
                                                className="bg-[var(--primary)] text-white text-sm px-4 py-2 rounded hover:bg-[var(--primary-hover)] transition"
                                            >
                                                Выбрать количество
                                            </button>
                                        ) : (
                                            <div className="w-full flex flex-col items-end gap-2">
                                                <div className="flex items-center gap-2">
                                                    <button
                                                        onClick={() => handleChangeQuantity(offer.supplier.id, -1)}
                                                        className="px-3 py-1 rounded border text-sm"
                                                    >
                                                        −
                                                    </button>
                                                    <span className="font-medium">{quantity}</span>
                                                    <button
                                                        onClick={() => handleChangeQuantity(offer.supplier.id, 1)}
                                                        className="px-3 py-1 rounded border text-sm"
                                                    >
                                                        +
                                                    </button>
                                                </div>
                                                <button
                                                    onClick={() => handleAddToCart(offer.supplier.id)}
                                                    className="w-full bg-[var(--primary)] text-white text-sm px-4 py-2 rounded hover:bg-[var(--primary-hover)] transition"
                                                >
                                                    Добавить
                                                </button>
                                            </div>
                                        )}
                                    </div>
                                </li>
                            );
                        })}
                    </ul>
                )}
            </div>

            <div className="space-y-6 rounded-xl border bg-white p-6 shadow">
                <div className="flex items-center justify-between gap-4 flex-wrap">
                    <h2 className="text-xl font-semibold text-[var(--foreground)]">Отзывы и рейтинг</h2>
                    <div className="text-sm text-slate-500">
                        Средняя оценка: <span className="font-semibold text-amber-600">{product.rating_average.toFixed(1)}</span> из 5
                    </div>
                </div>

                <div className="grid gap-3 md:grid-cols-[140px_1fr_auto]">
                    <select
                        value={reviewRating}
                        onChange={(event) => setReviewRating(Number(event.target.value))}
                        className="rounded-lg border px-3 py-2"
                    >
                        <option value={5}>5 звёзд</option>
                        <option value={4}>4 звезды</option>
                        <option value={3}>3 звезды</option>
                        <option value={2}>2 звезды</option>
                        <option value={1}>1 звезда</option>
                    </select>
                    <input
                        value={reviewComment}
                        onChange={(event) => setReviewComment(event.target.value)}
                        placeholder="Поделитесь впечатлением о товаре"
                        className="rounded-lg border px-3 py-2"
                    />
                    <button
                        onClick={submitReview}
                        disabled={reviewSubmitting}
                        className="rounded-lg bg-[var(--primary)] px-4 py-2 text-sm font-semibold text-white"
                    >
                        {reviewSubmitting ? "Сохраняем..." : "Оставить отзыв"}
                    </button>
                </div>

                {reviews.length === 0 ? (
                    <p className="text-sm text-slate-500">Пока нет отзывов. Будьте первым.</p>
                ) : (
                    <div className="space-y-4">
                        {reviews.map((review) => (
                            <div key={review.id} className="rounded-xl border border-slate-200 p-4">
                                <div className="flex items-center justify-between gap-3">
                                    <div>
                                        <p className="font-semibold text-slate-900">{review.user_name}</p>
                                        <p className="text-sm text-amber-600">{"★".repeat(review.rating)}</p>
                                    </div>
                                    <p className="text-xs text-slate-500">
                                        {new Date(review.created_at).toLocaleDateString()}
                                    </p>
                                </div>
                                <p className="mt-2 text-sm text-slate-700">{review.comment || "Без комментария"}</p>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
