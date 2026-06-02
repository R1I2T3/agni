import { defineConfig } from 'vite'
import viteReact from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

import { TanStackRouterVite } from '@tanstack/router-plugin/vite'
import { resolve } from 'node:path'

const isLocal = process.env.ENV_MODE === 'local'
const backendTarget = isLocal ? 'http://localhost:8080/' : 'http://agni-backend:8080/'

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
        target: backendTarget, // Adjust to your backend server
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
