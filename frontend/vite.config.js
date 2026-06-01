import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import wails from "@wailsio/runtime/plugins/vite";
import tailwindcss from '@tailwindcss/vite';
import path from "path";  // 需要安装 @types/node

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),   // 将 @ 映射到 src 目录
    },
  },
  plugins: [
    vue(), 
    wails("./bindings"),
    tailwindcss()
  ],
});