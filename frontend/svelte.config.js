import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/kit/vite';
// import path from "path";

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://kit.svelte.dev/docs/integrations#preprocessors
	// for more information about preprocessors
	preprocess: vitePreprocess(),


	kit: {
		adapter: adapter({
			assets: "assets",
			fallback: 'index.html',
		}),
		prerender: {
			entries: []
		},
		paths: {
			assets: "http://localhost:3000/assets",
			base: "/spa"
		},
	},

};

export default config;
