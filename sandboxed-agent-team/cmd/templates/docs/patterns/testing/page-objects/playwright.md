# Playwright Page Objects

When writing Playwright E2E tests, encapsulate DOM selectors in a page object class so tests survive layout changes without modification.

```typescript
// e2e/page-objects/EmployeeFormPage.ts
import { Page, Locator } from '@playwright/test';

export class EmployeeFormPage {
    constructor(private readonly page: Page) {}

    readonly heading = (): Locator => this.page.getByRole('heading');
    readonly nameField = (): Locator => this.page.getByLabel('Name');
    readonly emailField = (): Locator => this.page.getByLabel('Email');
    readonly saveButton = (): Locator => this.page.getByRole('button', { name: 'Save' });

    async submit(name: string, email: string): Promise<void> {
        await this.nameField().fill(name);
        await this.emailField().fill(email);
        await this.saveButton().click();
    }
}
```

```typescript
// e2e/employee.spec.ts
test('creates an employee successfully', async ({ page }) => {
    await page.goto('/employees/new');
    const form = new EmployeeFormPage(page);

    await form.submit('Alice Smith', 'alice@example.com');

    await expect(page.getByText('Alice Smith')).toBeVisible();
});
```

When an action navigates to a new view, the page object should return the next
page's page object so tests can chain.
