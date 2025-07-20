// next.config.ts
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://backend:8080/api/:path*", // имя сервиса из docker-compose
      },
    ];
  },
};

export default nextConfig;
