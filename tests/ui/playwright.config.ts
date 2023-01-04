import type {PlaywrightTestConfig} from '@playwright/test';

const config: PlaywrightTestConfig = {
    testDir: 'tests',
    retries: 2,
    reporter: [['html', {open: 'never'}]],
    use: {
        baseURL: 'http://localhost:3000',
        trace: 'on-first-retry',
        video: 'on-first-retry',
    },
};

export default config;
