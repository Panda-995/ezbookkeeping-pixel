import { describe, expect, it } from 'vitest';

import { isSuccessfulApiResponse } from '@/core/api_result.ts';

describe('isSuccessfulApiResponse', () => {
    it('accepts zero and null as valid successful API results', () => {
        expect(isSuccessfulApiResponse({ success: true, result: 0 })).toBe(true);
        expect(isSuccessfulApiResponse({ success: true, result: null })).toBe(true);
    });

    it('rejects missing and unsuccessful API responses', () => {
        expect(isSuccessfulApiResponse(null)).toBe(false);
        expect(isSuccessfulApiResponse(undefined)).toBe(false);
        expect(isSuccessfulApiResponse({ success: false, result: 0 })).toBe(false);
    });
});
