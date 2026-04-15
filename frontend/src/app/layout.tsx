import type { Metadata } from "next";
import "./globals.css";
import { AuthProvider } from "@/context/AuthContext";
import Navbar from "@/components/Navbar";
import { CartProvider } from "@/context/CartContext";
import Footer from "@/components/Footer";

export const metadata: Metadata = {
    title: "E-commerce Diplom",
    description: "Frontend for diplom project",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
    return (
        <html lang="ru" className="overflow-x-hidden">
        <body className="antialiased bg-[var(--background)] text-[var(--foreground)] overflow-x-hidden">
        <AuthProvider>
            <CartProvider>
                <Navbar />
                <main className="pt-24 px-4 sm:px-6 lg:px-12 py-6 min-h-screen">
                    {children}
                </main>

                <Footer />
            </CartProvider>
        </AuthProvider>
        </body>
        </html>
    );
}

