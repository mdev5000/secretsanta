import { expect, test } from '@playwright/test';

test('example test', async ({ page }) => {
	await page.goto('/app/example');
	await page.click('data-testid=fetcher >> button');
	let username = await page.textContent('data-testid=username');
	expect(username).toBe('Username: username');
});
