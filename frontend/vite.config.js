import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// During local development, proxy /api requests to `vercel dev` (run
// alongside `npm run dev` on port 3000) so the React app can talk to the
// real Go functions without CORS friction.
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': 'http://localhost:3000'
    }
  }
})
