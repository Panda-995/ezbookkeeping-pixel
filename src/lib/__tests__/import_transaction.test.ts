import { describe, expect, it } from 'vitest';

import { TransactionType } from '@/core/transaction.ts';
import {
    collectInvalidImportAccounts,
    getAllBatchCreateItemValues,
    getInvertedBatchCreateItemValues,
    getSelectedBatchCreateItems
} from '@/lib/import_transaction.ts';

describe('collectInvalidImportAccounts', () => {
    it('collects and deduplicates missing source and destination accounts in an empty project', () => {
        const result = collectInvalidImportAccounts([
            {
                type: TransactionType.Expense,
                sourceAccountId: '0',
                originalSourceAccountName: 'Wallet',
                originalSourceAccountCurrency: 'CNY',
                destinationAccountId: '',
                originalDestinationAccountName: undefined,
                originalDestinationAccountCurrency: undefined
            },
            {
                type: TransactionType.Transfer,
                sourceAccountId: '0',
                originalSourceAccountName: 'Wallet',
                originalSourceAccountCurrency: 'CNY',
                destinationAccountId: '0',
                originalDestinationAccountName: 'Bank',
                originalDestinationAccountCurrency: 'USD'
            }
        ], {}, '(Empty)', 'EUR');

        expect(result.items).toEqual([
            { name: 'Wallet', value: 'Wallet' },
            { name: 'Bank', value: 'Bank' }
        ]);
        expect(result.currencies).toEqual({
            Wallet: 'CNY',
            Bank: 'USD'
        });
    });

    it('ignores resolved accounts and uses a safe display name and default currency for empty source data', () => {
        const result = collectInvalidImportAccounts([
            {
                type: TransactionType.Expense,
                sourceAccountId: 'existing-account',
                originalSourceAccountName: 'Existing',
                originalSourceAccountCurrency: 'USD',
                destinationAccountId: '',
                originalDestinationAccountName: undefined,
                originalDestinationAccountCurrency: undefined
            },
            {
                type: TransactionType.Expense,
                sourceAccountId: '0',
                originalSourceAccountName: '',
                originalSourceAccountCurrency: '',
                destinationAccountId: '',
                originalDestinationAccountName: undefined,
                originalDestinationAccountCurrency: undefined
            }
        ], {
            'existing-account': {}
        }, '(Empty)', 'CNY');

        expect(result.items).toEqual([{ name: '(Empty)', value: '' }]);
        expect(result.currencies).toEqual({ '': 'CNY' });
    });

    it('keeps account names that overlap with object prototype fields', () => {
        const result = collectInvalidImportAccounts([{
            type: TransactionType.Expense,
            sourceAccountId: '0',
            originalSourceAccountName: '__proto__',
            originalSourceAccountCurrency: 'USD',
            destinationAccountId: '',
            originalDestinationAccountName: undefined,
            originalDestinationAccountCurrency: undefined
        }], {}, '(Empty)', 'CNY');

        expect(result.items).toEqual([{
            name: '__proto__',
            value: '__proto__'
        }]);
        expect(Object.hasOwn(result.currencies, '__proto__')).toBe(true);
        expect(result.currencies['__proto__']).toBe('USD');
    });
});

describe('batch create selection helpers', () => {
    const items = [
        { name: '(Empty)', value: '' },
        { name: 'Food', value: 'Food' },
        { name: 'Travel', value: 'Travel' }
    ];

    it('uses source values as selection keys instead of localized display names', () => {
        expect(getAllBatchCreateItemValues(items)).toEqual(['', 'Food', 'Travel']);
        expect(getSelectedBatchCreateItems(items, ['', 'Travel'])).toEqual([
            { name: '(Empty)', value: '' },
            { name: 'Travel', value: 'Travel' }
        ]);
    });

    it('inverts selections without losing an empty source value', () => {
        expect(getInvertedBatchCreateItemValues(items, ['', 'Travel'])).toEqual(['Food']);
    });
});
