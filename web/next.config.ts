import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  experimental: {
    dynamicIO: true,
    useCache: true,
  },
};

export default nextConfig;
