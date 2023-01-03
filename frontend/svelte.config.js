import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/kit/vite';

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
			// Note: this path is replaced as part of the build process, see the assets.replacePaths makefile command
			// for more details.
			assets: "http://localhost:3000/replaceme/assets",
			base: "/app"
		},
	},

};

export default config;
