import { expect, test } from '@playwright/test';

test('example test', async ({ page }) => {
	await page.goto('/app/example');
	// Wait till we can see the button before trying to click.
	await page.textContent('data-testid=fetcher >> button');
	await page.click('data-testid=fetcher >> button');
	let username = await page.textContent('data-testid=username');
	expect(username).toBe('Username: username');
});
