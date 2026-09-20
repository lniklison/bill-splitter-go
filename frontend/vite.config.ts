import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import { fileURLToPath, URL } from "node:url"

export default defineConfig({
  plugins: [vue()],
  build: {
    outDir: fileURLToPath(new URL("../internal/web/dist", import.meta.url)),
    emptyOutDir: true,
  },
})
