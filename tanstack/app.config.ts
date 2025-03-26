import { defineConfig } from "@tanstack/react-start/config";
import tsConfigPaths from "vite-tsconfig-paths";
import "dotenv/config";

export default defineConfig({
  tsr: {
    appDirectory: "src",
  },
  server: {
    preset: "vercel",
  },
  vite: {
    define: {
      "process.env": JSON.stringify(process.env),
    },
    plugins: [
      tsConfigPaths({
        projects: ["./tsconfig.json"],
      }),
    ],
  },
});
