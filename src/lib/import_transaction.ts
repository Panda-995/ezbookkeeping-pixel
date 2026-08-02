import type { NameValue } from '@/core/base.ts';
import { TransactionType } from '@/core/transaction.ts';

export interface ImportAccountReference {
    readonly type: number;
    readonly sourceAccountId: string;
    readonly originalSourceAccountName: string;
    readonly originalSourceAccountCurrency: string;
    readonly destinationAccountId: string;
    readonly originalDestinationAccountName?: string;
    readonly originalDestinationAccountCurrency?: string;
}

export interface InvalidImportAccountCollection {
    readonly items: NameValue[];
    readonly currencies: Record<string, string>;
}

export function collectInvalidImportAccounts(
    transactions: readonly ImportAccountReference[],
    allAccountsMap: Readonly<Record<string, unknown>>,
    emptyDisplayName: string,
    defaultCurrency: string
): InvalidImportAccountCollection {
    const items: NameValue[] = [];
    const currencies: Record<string, string> = Object.create(null) as Record<string, string>;
    const addedNames = new Set<string>();

    const addAccount = (name: string, currency: string | undefined): void => {
        if (addedNames.has(name)) {
            return;
        }

        addedNames.add(name);
        items.push({
            name: name || emptyDisplayName,
            value: name
        });
        currencies[name] = currency || defaultCurrency;
    };

    for (const transaction of transactions) {
        const sourceAccountId = transaction.sourceAccountId;

        if (!sourceAccountId || sourceAccountId === '0' || !allAccountsMap[sourceAccountId]) {
            addAccount(transaction.originalSourceAccountName, transaction.originalSourceAccountCurrency);
        }

        const destinationAccountId = transaction.destinationAccountId;

        if (transaction.type === TransactionType.Transfer &&
            typeof transaction.originalDestinationAccountName === 'string' &&
            (!destinationAccountId || destinationAccountId === '0' || !allAccountsMap[destinationAccountId])) {
            addAccount(transaction.originalDestinationAccountName, transaction.originalDestinationAccountCurrency);
        }
    }

    return {
        items,
        currencies
    };
}

export function getAllBatchCreateItemValues(items: readonly NameValue[]): string[] {
    return items.map(item => item.value);
}

export function getSelectedBatchCreateItems(items: readonly NameValue[], selectedValues: readonly string[]): NameValue[] {
    const selectedValueSet = new Set(selectedValues);
    return items.filter(item => selectedValueSet.has(item.value));
}

export function getInvertedBatchCreateItemValues(items: readonly NameValue[], selectedValues: readonly string[]): string[] {
    const selectedValueSet = new Set(selectedValues);
    return items.filter(item => !selectedValueSet.has(item.value)).map(item => item.value);
}
