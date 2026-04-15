import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "e-catalog.gov.kz",
        pathname: "/media/products/**",
      },
    ],
  },
};

export default nextConfig;
