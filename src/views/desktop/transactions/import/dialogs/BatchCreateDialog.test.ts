import { readFileSync } from 'node:fs';

import { describe, expect, it } from 'vitest';

describe('import batch create dialog selection', () => {
    const source = readFileSync(new URL('./BatchCreateDialog.vue', import.meta.url), 'utf8');

    it('binds list selection to source values and Vuetify selected state', () => {
        expect(source).toContain(':value="item.value"');
        expect(source).toContain('<template #prepend="{ isSelected, select }">');
        expect(source).toContain('@update:model-value="select"');
        expect(source).not.toContain('updateSelectedNames');
    });

    it('supports creating missing accounts for an empty project', () => {
        expect(source).toContain("type === 'account'");
        expect(source).toContain("tt('Create Nonexistent Accounts')");
    });

    it('resolves cancellation without an unhandled promise rejection', () => {
        expect(source).toContain('resolveFunc?.(null)');
        expect(source).not.toContain('rejectFunc');
    });
});
