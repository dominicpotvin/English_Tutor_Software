import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// During `npm run dev`, /api is proxied to the Go backend so the frontend
// runs same-origin. In production nginx performs the same proxying.
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5174,
    proxy: {
      "/api": "http://localhost:8096",
    },
  },
});
