// app.config.ts
import { defineConfig } from "@tanstack/react-start/config";
import tsConfigPaths from "vite-tsconfig-paths";
import "dotenv/config";
var app_config_default = defineConfig({
  tsr: {
    appDirectory: "src"
  },
  server: {
    preset: "vercel"
  },
  vite: {
    define: {
      "process.env": JSON.stringify(process.env)
    },
    plugins: [
      tsConfigPaths({
        projects: ["./tsconfig.json"]
      })
    ]
  }
});
export {
  app_config_default as default
};
