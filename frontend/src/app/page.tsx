'use client';

import { useEffect, useState } from "react";
import ProductCard from "@/components/ProductCard";
import api from "@/lib/axios";
import Link from "next/link";
import { Product } from "@/types";

const BANNER_IMAGES = ["/banner/img1.png", "/banner/img2.png", "/banner/img3.png", "/banner/img4.png"];

export default function HomePage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [topProducts, setTopProducts] = useState<Product[]>([]);
  const [catalogPage, setCatalogPage] = useState(0);
  const [prevPage, setPrevPage] = useState<number | null>(null);
  const [bannerIndex, setBannerIndex] = useState(0);
  const [prevBannerIndex, setPrevBannerIndex] = useState<number | null>(null);

  useEffect(() => {
    const fetch = async () => {
      try {
        const [catalogRes, featuredRes] = await Promise.all([
          api.get("/product/list?limit=50&offset=0&sort_by=newest&sort_order=desc"),
          api.get("/product/featured?limit=5"),
        ]);
        const productList: Product[] = catalogRes.data.product_list;
        setProducts(productList);
        setTopProducts(featuredRes.data.product_list ?? []);
      } catch (e) {
        console.error("Ошибка загрузки продуктов:", e);
      }
    };
    fetch();
  }, []);

  useEffect(() => {
    const interval = setInterval(() => {
      setPrevBannerIndex(bannerIndex);
      setBannerIndex((prev) => (prev + 1) % BANNER_IMAGES.length);
    }, 8000);
    return () => clearInterval(interval);
  }, [bannerIndex]);

  useEffect(() => {
    const interval = setInterval(() => {
      setPrevPage(catalogPage);
      setCatalogPage((prev) => (prev + 1) % 5);
    }, 8000);
    return () => clearInterval(interval);
  }, [catalogPage]);

  return (
    <div className="mx-auto w-[92%] space-y-10 px-4 py-8">
      <section className="relative overflow-hidden rounded-[32px] border border-white/50 bg-[linear-gradient(135deg,rgba(15,118,110,0.92),rgba(217,119,6,0.88))] shadow-[0_30px_100px_rgba(15,23,42,0.12)] w-full min-h-[440px]">
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(255,255,255,0.18),transparent_28%),radial-gradient(circle_at_bottom_right,rgba(255,255,255,0.14),transparent_30%)]" />
        <div className="relative z-20 grid h-full gap-8 px-8 py-10 lg:grid-cols-[0.95fr_1.05fr] lg:px-12">
          <div className="flex flex-col justify-center">
            <p className="text-sm font-semibold uppercase tracking-[0.28em] text-white/80">Digital Marketplace</p>
            <h1 className="mt-4 max-w-xl text-4xl font-extrabold leading-tight text-white md:text-5xl">
              Все нужные товары в одном аккуратном каталоге
            </h1>
            <p className="mt-5 max-w-lg text-sm leading-7 text-white/84">
              Быстрый поиск, избранное, рейтинги, адресная доставка и удобное оформление заказа без перегруженного интерфейса.
            </p>
            <div className="mt-8 flex flex-wrap gap-3">
              <Link href="/catalog" className="rounded-2xl bg-white px-5 py-3 text-sm font-semibold text-[var(--primary)] transition hover:bg-white/90">
                Открыть каталог
              </Link>
              <Link href="/checkout" className="rounded-2xl border border-white/40 px-5 py-3 text-sm font-semibold text-white transition hover:bg-white/10">
                Перейти к checkout
              </Link>
            </div>
          </div>

          <div className="relative h-[340px] overflow-hidden rounded-[28px] border border-white/30 bg-white/12 backdrop-blur-sm lg:h-full">
            {BANNER_IMAGES.map((img, idx) => {
              const baseStyle = "absolute inset-0 h-full w-full transition-all duration-[1000ms] ease-in-out";
              const isCurrent = idx === bannerIndex;
              const isPrev = idx === prevBannerIndex;

              return (
                <div
                  key={idx}
                  className={`${baseStyle} ${isCurrent ? "translate-x-0 opacity-100 z-10" : isPrev ? "-translate-x-full opacity-0 z-0" : "translate-x-full opacity-0 z-0"}`}
                >
                  <img
                    src={img}
                    alt=""
                    className="absolute inset-0 h-full w-full object-cover blur-lg scale-110 opacity-35"
                  />
                  <div className="relative z-10 flex h-full items-center justify-center">
                    <img
                      src={img}
                      alt={`Banner ${idx}`}
                      className="max-h-[82%] max-w-[75%] object-contain"
                    />
                  </div>
                </div>
              );
            })}

            <div className="absolute bottom-4 left-1/2 flex -translate-x-1/2 gap-2">
              {BANNER_IMAGES.map((_, idx) => (
                <span
                  key={idx}
                  className={`h-2 w-2 rounded-full ${idx === bannerIndex ? "bg-white" : "bg-white/40"} transition`}
                />
              ))}
            </div>
          </div>
        </div>
      </section>


      <section className="rounded-[28px] border border-[var(--card-border)] shadow p-6 bg-gradient-to-br from-[var(--primary-soft)] to-white w-full">
        <h2 className="text-2xl font-bold mb-4 text-[var(--foreground)]">🔥 Топ-продаж</h2>
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-4">
          {topProducts.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      </section>

      <section className="rounded-[28px] border border-[var(--card-border)] shadow p-6 bg-[var(--card-bg)] w-full">
        <h2 className="text-2xl font-bold mb-4 text-[var(--foreground)]">🛍️ Каталог</h2>
        <div className="relative h-[750px] overflow-hidden">
          {[...Array(5)].map((_, idx) => {
            let className = "absolute inset-0 transition-all duration-[1000ms] ease-in-out";
            if (idx === catalogPage) {
              className += " translate-x-0 opacity-100 z-10";
            } else if (idx === prevPage) {
              className += " -translate-x-full opacity-0 z-0";
            } else {
              className += " translate-x-full opacity-0 z-0";
            }

            return (
              <div key={idx} className={className}>
                <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-4">
                  {products.slice(idx * 10, idx * 10 + 10).map((product) => (
                    <ProductCard key={product.id} product={product} />
                  ))}
                </div>
              </div>
            );
          })}
        </div>
        <div className="flex justify-center mt-4 gap-2">
          {[...Array(5)].map((_, idx) => (
            <span
              key={idx}
              className={`w-2 h-2 rounded-full ${
                idx === catalogPage ? "bg-[var(--primary)]" : "bg-gray-300"
              } transition`}
            />
          ))}
        </div>
        <div className="text-center mt-6">
          <Link
            href="/catalog"
            className="inline-block text-[var(--primary)] text-base font-semibold border border-[var(--primary)] px-5 py-2 rounded hover:bg-[var(--primary)] hover:text-white transition"
          >
            Смотреть весь каталог →
          </Link>
        </div>
      </section>

    </div>
  );
}
