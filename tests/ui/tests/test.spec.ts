import { expect, test } from '@playwright/test';

test('example test', async ({ page }) => {
	await page.goto('/app/example');
	// Wait till we can see the button before trying to click.
	await page.textContent('data-testid=fetcher >> button', {timeout: 5000});
	await page.click('data-testid=fetcher >> button');
	let status = await page.textContent('data-testid=status');
	expect(status).toBe('Status: some status');
});
