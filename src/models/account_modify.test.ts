import { describe, expect, it } from 'vitest';

import { AccountCategory } from '@/core/account.ts';
import { Account } from '@/models/account.ts';

describe('account modification request', () => {
    it('includes currency and balance after editing an existing account', () => {
        const account = Account.of({
            id: '1001',
            name: 'Cash',
            parentId: '0',
            category: 1,
            type: 1,
            icon: '1',
            color: '176b5b',
            currency: 'CNY',
            balance: 12345,
            comment: '',
            displayOrder: 1,
            hidden: false
        });
        const originalAccount = account.clone();

        account.currency = 'USD';
        account.balance = 67890;

        expect(account.toModifyRequest('session-id', undefined, undefined, originalAccount)).toMatchObject({
            currency: 'USD',
            balance: 67890
        });
    });

    it('omits unchanged money fields to avoid overwriting concurrent balance changes', () => {
        const account = Account.of({
            id: '1001',
            name: 'Cash',
            parentId: '0',
            category: 1,
            type: 1,
            icon: '1',
            color: '176b5b',
            currency: 'CNY',
            balance: 12345,
            comment: '',
            displayOrder: 1,
            hidden: false
        });

        expect(account.toModifyRequest('session-id', undefined, undefined, account)).toMatchObject({
            currency: undefined,
            balance: undefined
        });
    });

    it('includes edited currency and balance for an existing sub-account', () => {
        const originalAccount = Account.of({
            id: '2000',
            name: 'Cards',
            parentId: '0',
            category: 3,
            type: 2,
            icon: '100',
            color: '176b5b',
            currency: '---',
            balance: 0,
            comment: '',
            displayOrder: 1,
            hidden: false,
            subAccounts: [{
                id: '2001',
                name: 'Daily Card',
                parentId: '2000',
                category: 3,
                type: 1,
                icon: '100',
                color: '176b5b',
                currency: 'CNY',
                balance: -10000,
                comment: '',
                displayOrder: 1,
                hidden: false
            }]
        });
        const editedAccount = originalAccount.clone();
        const editedSubAccount = editedAccount.subAccounts?.[0];

        expect(editedSubAccount).toBeDefined();
        editedSubAccount!.currency = 'USD';
        editedSubAccount!.balance = -25000;

        const request = editedAccount.toModifyRequest(
            'session-id',
            editedAccount.subAccounts,
            undefined,
            originalAccount
        );

        expect(request.subAccounts?.[0]).toMatchObject({
            currency: 'USD',
            balance: -25000
        });
    });

    it('derives liability state from an edited category', () => {
        const account = Account.createNewAccount(AccountCategory.Cash, 'CNY', 0);

        account.category = 3;

        expect(account.isLiability).toBe(true);
    });
});
