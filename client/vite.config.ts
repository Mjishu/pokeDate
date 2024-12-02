import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
  server: {
    proxy: {
      "/kittens": {
        target: "https://api.pexels.com/v1/search?query=nature&per_page=1",
        changeOrigin: true
      }
    }
  }
})
