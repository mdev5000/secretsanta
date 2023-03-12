import type {PlaywrightTestConfig} from '@playwright/test';
import * as process from "process";

let reportPath=undefined;
let baseURL='http://localhost:5173';

if (process.env.REPORT_PATH !== "") {
    reportPath = process.env.REPORT_PATH;
}

if (process.env.BASE_URL !== "") {
    baseURL = process.env.BASE_URL;
}

const config: PlaywrightTestConfig = {
    testDir: 'tests',
    retries: 2,
    reporter: [['html', {open: 'never', outputFolder: reportPath}]],
    use: {
        baseURL: baseURL,
        trace: 'on-first-retry',
        video: 'on-first-retry',
    },
};

export default config;
