import { beforeEach, describe, expect, it, vi } from 'vitest';
import { createPinia, setActivePinia } from 'pinia';

const { getImportTransactionsProcessMock } = vi.hoisted(() => ({
    getImportTransactionsProcessMock: vi.fn()
}));

vi.mock('@/lib/services.ts', () => ({
    default: {
        getImportTransactionsProcess: getImportTransactionsProcessMock
    }
}));

vi.mock('@/lib/logger.ts', () => ({
    default: {
        error: vi.fn()
    }
}));

vi.mock('@/lib/userstate.ts', () => ({
    getUserTransactionDraft: vi.fn(() => null),
    updateUserTransactionDraft: vi.fn(),
    clearUserTransactionDraft: vi.fn()
}));

vi.mock('@/stores/setting.ts', () => ({ useSettingsStore: vi.fn(() => ({})) }));
vi.mock('@/stores/user.ts', () => ({ useUserStore: vi.fn(() => ({})) }));
vi.mock('@/stores/account.ts', () => ({ useAccountsStore: vi.fn(() => ({})) }));
vi.mock('@/stores/transactionCategory.ts', () => ({ useTransactionCategoriesStore: vi.fn(() => ({})) }));
vi.mock('@/stores/overview.ts', () => ({ useOverviewStore: vi.fn(() => ({})) }));
vi.mock('@/stores/statistics.ts', () => ({ useStatisticsStore: vi.fn(() => ({})) }));
vi.mock('@/stores/explorer.ts', () => ({ useExplorersStore: vi.fn(() => ({})) }));
vi.mock('@/stores/exchangeRates.ts', () => ({ useExchangeRatesStore: vi.fn(() => ({})) }));

import { useTransactionsStore } from '@/stores/transaction.ts';

describe('transaction import progress', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('returns zero while an import has just started', async () => {
        getImportTransactionsProcessMock.mockResolvedValueOnce({
            data: {
                success: true,
                result: 0
            }
        });

        const store = useTransactionsStore();
        await expect(store.getImportTransactionsProcess({ clientSessionId: 'session' })).resolves.toBe(0);
    });

    it('returns null when the server has no progress record yet', async () => {
        getImportTransactionsProcessMock.mockResolvedValueOnce({
            data: {
                success: true,
                result: null
            }
        });

        const store = useTransactionsStore();
        await expect(store.getImportTransactionsProcess({ clientSessionId: 'session' })).resolves.toBeNull();
    });
});
