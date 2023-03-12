import type {PlaywrightTestConfig} from '@playwright/test';
import * as process from "process";

let reportPath=undefined;

if (process.env.REPORT_PATH !== "") {
    reportPath = process.env.REPORT_PATH;
}

const config: PlaywrightTestConfig = {
    testDir: 'tests',
    retries: 2,
    reporter: [['html', {open: 'never', outputFolder: reportPath}]],
    use: {
        baseURL: 'http://localhost:5173',
        trace: 'on-first-retry',
        video: 'on-first-retry',
    },
};

export default config;
