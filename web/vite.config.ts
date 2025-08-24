import { defineConfig } from 'vite'
import viteReact from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

import { TanStackRouterVite } from '@tanstack/router-plugin/vite'
import { resolve } from 'node:path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    TanStackRouterVite({ autoCodeSplitting: true }),
    viteReact(),
    tailwindcss(),
  ],
  server: {
    host: '0.0.0.0', // Important for Docker
    port: 3000,
    strictPort: true,
    watch: {
      usePolling: true, // Important for Docker on Windows
    },
    proxy: {
      '/api': {
        target: 'http://agni-backend:8080/', // Adjust to your backend server
        changeOrigin: true,
        secure: false, // If your backend is not using HTTPS
      },
    },
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
    },
  },
})
