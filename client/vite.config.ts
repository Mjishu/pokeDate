import { defineConfig } from 'vitest/config';
import { sveltekit } from '@sveltejs/kit/vite';

const isProd = process.env.NODE_ENV === "production";
const apiPath = isProd ? "https://pokefind-server.fly.dev" : "http://localhost:8080"

export default defineConfig({
	plugins: [sveltekit()],

	test: {
		include: ['src/**/*.{test,spec}.{js,ts}']
	}, server: {
		proxy: {
			'/api': {
				target: apiPath,
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/api/, ''),
			}
		}
	}
});