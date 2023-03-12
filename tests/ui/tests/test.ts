import { expect, test } from '@playwright/test';

test('example test', async ({ page }) => {
	await page.goto('/app/example');
	// expect(await page.textContent('//*[@data-testid="fetcher"]')).toBe("Fetch");
	await page.click('data-testid=fetcher');
	let username = await page.textContent('data-testid=username');
	expect(username).toBe('Username: username');
});
