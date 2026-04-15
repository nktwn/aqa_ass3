// types/index.ts

export interface Supplier {
    id: number;
    name: string;
    order_amount: number;
    free_delivery_amount: number;
    delivery_fee: number;
}

export interface Product {
    id: number;
    name: string;
    image: string;
    rating_average: number;
    rating_count: number;
    lowest_product_supplier: {
        price: number;
        sell_amount: number;
        supplier: Supplier;
    };
}

export interface CartItem {
    productId: number;
    supplierId: number;
    quantity: number;
}

export interface Cart {
    supplier: Supplier;
    items: CartItem[];
}

export interface Order {
    id: number;
    status: string;
    createdAt: string;
    supplier: Supplier;
    items: CartItem[];
}

export interface User {
    id: number;
    email: string;
    name: string;
}

export interface Address {
    id?: number;
    description: string;
    street: string;
}

export interface Review {
    id: number;
    user_id: number;
    user_name: string;
    rating: number;
    comment: string;
    created_at: string;
}

export interface Notification {
    id: number;
    title: string;
    message: string;
    kind: string;
    is_read: boolean;
    order_id?: number;
    created_at: string;
}
