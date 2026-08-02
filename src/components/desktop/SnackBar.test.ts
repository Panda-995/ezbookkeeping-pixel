import { readFileSync } from 'node:fs';

import { describe, expect, it } from 'vitest';

describe('desktop snackbar theme', () => {
    const source = readFileSync(new URL('./SnackBar.vue', import.meta.url), 'utf8');

    it('uses the notification color pair so message text contrasts with its background', () => {
        expect(source).toContain('<v-snackbar color="notification-background" v-model="showState">');
    });
});
