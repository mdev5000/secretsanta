import { expect, test } from '@playwright/test';
import axios from 'axios'
import * as process from 'process'

let baseAPIUrl = process.env.BASE_URL;

function delay(time) {
    return new Promise(resolve => setTimeout(resolve, time));
}

async function retryBool(fn: () => Promise<bool>, config: {delay: number, count: number} = undefined): Promise<boolean> {
    return await retry(async () => {
        const isOk = await fn();
        if (!isOk) {
            throw "unsuccessful";
        }
    }, config);
}

async function retry(fn: () => Promise<undefined>, config: {delay: number, count: number} = undefined): Promise<boolean> {
    if (config === undefined) {
        config = {delay: 0, count: 0}
    }
    if (config.delay == 0) {
        config.delay = 200;
    }
    if (config.count == 0) {
        config.count = 10;
    }
    for (let i = 0; i < config.count; i++) {
        try {
            await fn();
            return true;
        } catch (e) {
            await delay(config.delay);
        }
    }
    return false;
}

test('can setup application', async ({ page }) => {

    // Reset application completely.
    await test.step('reset application completely', async () => {

        if (baseAPIUrl === undefined || baseAPIUrl == "") {
            baseAPIUrl = "http://localhost:3000"
        }
        const resp = await axios.post(baseAPIUrl + "/test/delete-all-restart");
        expect(resp.status, 'delete and restarts ok').toBe(200);

        // Wait for the server to be ready again.

        const ok = await retry(async () => {
            await axios.get(baseAPIUrl + "/ready")
        });
        expect(ok, 'is ready after restart').toBe(true)

    });

    // Go through the setup wizard.

    await page.goto('/app/setup', {timeout: 5000});
    let title = await page.textContent('h1', {timeout: 5000});
    await page.click('data-testid=next-btn', {timeout: 5000});

    await page.fill('data-testid=admin-username', 'test-admin');
    await page.fill('data-testid=admin-firstname', 'adminFirst');
    await page.fill('data-testid=admin-lastname', 'adminList');
    await page.fill('data-testid=admin-password', 'adminPassword');
    await page.click('data-testid=setup-submit');

    // Make sure we get setup successful message

    let status = await page.textContent('data-testid=status');
    expect(status, 'completes setup successfully').toBe('Setup completed, will redirect to app in a moment...');
    await delay(2000);

    // Should be redirected to login screen.

    // @todo figure out why redirect isn't working
    await page.goto('/app', {timeout: 5000});

    let loginTitle = await page.textContent('h1', {timeout: 5000});
    expect(loginTitle, 'goes to login page').toBe('Login');

    // Test admin user can login

    await page.fill('data-testid=username', 'test-admin');
    await page.fill('data-testid=password', 'adminPassword');
    await page.click('data-testid=login-submit');

    const loginSuccess = await retryBool(async () => {
        const status = await page.$('data-testid=login-submit');
        if (status && ((await status.textContent()) == 'Success')) {
            return true;
        }
        const h1 = await page.$('h1', {timeout: 5000});
        return h1 && ((await h1.textContent()) == 'Secret Santa'); 
    });

    expect(loginSuccess, 'logs in successfully').toBe(true);

    // Should login successful and see home page.

    await page.goto('/app', {timeout: 5000});

    let h1 = await page.textContent('h1', {timeout: 5000});
    expect(h1).toBe('Secret Santa', 'can see the homepage');
});
