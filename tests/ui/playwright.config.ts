import type { PlaywrightTestConfig } from '@playwright/test';

const config: PlaywrightTestConfig = {
	webServer: {
		command: 'npm run build && npm run preview',
		port: 4173
	},
	testDir: 'tests',
	retries: 2,
	reporter: [['html', {open: 'never'}]],
	use: {
		trace: 'on-first-retry',
		video: 'on-first-retry',
	},
};

export default config;
