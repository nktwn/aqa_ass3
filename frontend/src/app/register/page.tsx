'use client';

import { useState } from "react";
import { useAuth } from "@/hooks/useAuth";
import { useRouter } from "next/navigation";
import Link from "next/link";

export default function RegisterPage() {
    const { register } = useAuth();
    const router = useRouter();

    const [form, setForm] = useState({
        name: '',
        phone_number: '',
        password: '',
        confirm_password: '',
    });
    const [error, setError] = useState("");

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (form.password !== form.confirm_password) {
            setError("Пароли не совпадают");
            return;
        }
        try {
            await register(form);
            router.push("/");
        } catch {
            setError("Ошибка регистрации");
        }
    };

    return (
        <div className="px-4 pt-28 pb-20 flex justify-center">
            <div className="w-full max-w-2xl border rounded-2xl shadow-xl bg-gradient-to-br from-[var(--primary)]/10 to-white p-10 space-y-8">
                <h1 className="text-3xl font-bold text-center text-[var(--foreground)]">📝 Регистрация</h1>
                <form onSubmit={handleSubmit} className="space-y-6">
                    <input
                        type="text"
                        placeholder="Имя"
                        value={form.name}
                        onChange={(e) => setForm({ ...form, name: e.target.value })}
                        className="w-full p-4 border rounded-lg text-lg"
                    />
                    <input
                        type="tel"
                        placeholder="Телефон"
                        value={form.phone_number}
                        onChange={(e) => setForm({ ...form, phone_number: e.target.value })}
                        className="w-full p-4 border rounded-lg text-lg"
                    />
                    <input
                        type="password"
                        placeholder="Пароль"
                        value={form.password}
                        onChange={(e) => setForm({ ...form, password: e.target.value })}
                        className="w-full p-4 border rounded-lg text-lg"
                    />
                    <input
                        type="password"
                        placeholder="Повторите пароль"
                        value={form.confirm_password}
                        onChange={(e) => setForm({ ...form, confirm_password: e.target.value })}
                        className="w-full p-4 border rounded-lg text-lg"
                    />
                    {error && <p className="text-red-600 text-center">{error}</p>}
                    <button className="w-full bg-blue-600 text-white py-3 text-lg rounded-lg hover:bg-blue-700 transition">
                        Зарегистрироваться
                    </button>
                </form>

                <p className="text-center text-gray-700 text-base">
                    Уже есть аккаунт?{" "}
                    <Link href="/login" className="text-[var(--primary)] hover:underline font-medium">
                        Войти →
                    </Link>
                </p>
            </div>
        </div>
    );
}
