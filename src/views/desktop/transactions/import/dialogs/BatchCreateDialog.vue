<template>
    <v-dialog width="600" :persistent="submitting || !!selectedValues.length" :model-value="showState"
              @update:model-value="onShowStateChanged">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <div class="d-flex flex-wrap">
                    <h4 class="text-h4 text-wrap" v-if="type === 'expenseCategory'">{{ tt('Create Nonexistent Expense Categories') }}</h4>
                    <h4 class="text-h4 text-wrap" v-if="type === 'incomeCategory'">{{ tt('Create Nonexistent Income Categories') }}</h4>
                    <h4 class="text-h4 text-wrap" v-if="type === 'transferCategory'">{{ tt('Create Nonexistent Transfer Categories') }}</h4>
                    <h4 class="text-h4 text-wrap" v-if="type === 'account'">{{ tt('Create Nonexistent Accounts') }}</h4>
                    <h4 class="text-h4 text-wrap" v-if="type === 'tag'">{{ tt('Create Nonexistent Transaction Tags') }}</h4>
                    <v-spacer/>
                    <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                           :disabled="submitting || !invalidItems || !invalidItems.length" :icon="true">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent">
                            <v-list>
                                <v-list-item :prepend-icon="mdiSelectAll"
                                             :title="tt('Select All')"
                                             :disabled="!invalidItems || !invalidItems.length"
                                             @click="selectAllItems"></v-list-item>
                                <v-list-item :prepend-icon="mdiSelect"
                                             :title="tt('Select None')"
                                             :disabled="!invalidItems || !invalidItems.length"
                                             @click="selectNoneItems"></v-list-item>
                                <v-list-item :prepend-icon="mdiSelectInverse"
                                             :title="tt('Invert Selection')"
                                             :disabled="!invalidItems || !invalidItems.length"
                                             @click="selectInvertItems"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>
            <v-card-text class="d-flex flex-column flex-md-row flex-grow-1 overflow-y-auto">
                <v-row>
                    <v-col cols="12" class="px-0">
                        <v-list class="py-0" density="compact" select-strategy="classic"
                                :disabled="submitting" v-model:selected="selectedValues">
                            <v-list-item class="mx-1 px-2 py-0"
                                         :key="item.value" :value="item.value" :title="item.name"
                                         v-for="item in invalidItems">
                                <template #prepend="{ isSelected, select }">
                                    <v-list-item-action start>
                                        <v-checkbox-btn :model-value="isSelected" @update:model-value="select"></v-checkbox-btn>
                                    </v-list-item-action>
                                </template>
                            </v-list-item>
                        </v-list>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn :disabled="submitting || !selectedValues || !selectedValues.length" @click="confirm">
                        {{ tt('OK') }}
                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="submitting" @click="cancel">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { type NameValue, values } from '@/core/base.ts';
import type { ErrorResponse } from '@/core/api.ts';
import { AccountCategory } from '@/core/account.ts';
import { CategoryType } from '@/core/category.ts';
import { AUTOMATICALLY_CREATED_CATEGORY_ICON_ID } from '@/consts/icon.ts';
import { DEFAULT_CATEGORY_COLOR } from '@/consts/color.ts';
import { DEFAULT_TAG_GROUP_ID } from '@/consts/tag.ts';

import { Account } from '@/models/account.ts';
import { type TransactionCategoryCreateRequest, type TransactionCategoryCreateWithSubCategories, TransactionCategory } from '@/models/transaction_category.ts';
import { type TransactionTagCreateRequest, TransactionTag } from '@/models/transaction_tag.ts';

import { isDefined } from '@/lib/common.ts';
import { getCurrentUnixTime } from '@/lib/datetime.ts';
import {
    getAllBatchCreateItemValues,
    getInvertedBatchCreateItemValues,
    getSelectedBatchCreateItems
} from '@/lib/import_transaction.ts';
import { generateRandomUUID } from '@/lib/misc.ts';

import {
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiDotsVertical
} from '@mdi/js';

export type BatchCreateDialogDataType = 'expenseCategory' | 'incomeCategory' | 'transferCategory' | 'account' | 'tag';

type SnackBarType = InstanceType<typeof SnackBar>;

type BatchCreateError = ({ message: string } | { error: ErrorResponse }) & { processed?: boolean };

interface BatchCreateDialogResponse {
    sourceTargetMap: Record<string, string>;
}

const { tt } = useI18n();

const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const submitting = ref<boolean>(false);
const type = ref<BatchCreateDialogDataType | ''>('');
const invalidItems = ref<NameValue[] | undefined>([]);
const selectedValues = ref<string[]>([]);
const accountCurrencies = ref<Record<string, string>>({});

let resolveFunc: ((response: BatchCreateDialogResponse | null) => void) | null = null;

function buildBatchCreateCategoryResponse(createdCategories: Record<number, TransactionCategory[]>): BatchCreateDialogResponse {
    const displayNameSourceItemMap: Record<string, string> = Object.create(null) as Record<string, string>;
    const sourceTargetMap: Record<string, string> = Object.create(null) as Record<string, string>;

    for (const item of (invalidItems.value || [])) {
        displayNameSourceItemMap[item.name] = item.value;
    }

    for (const categories of values(createdCategories)) {
        for (const category of categories) {
            if (!category.subCategories || category.subCategories.length < 1) {
                continue;
            }

            for (const subCategory of category.subCategories) {
                const sourceItem = displayNameSourceItemMap[subCategory.name];

                if (!isDefined(sourceItem)) {
                    continue;
                }

                sourceTargetMap[sourceItem] = subCategory.id;
            }
        }
    }

    const response: BatchCreateDialogResponse = {
        sourceTargetMap: sourceTargetMap
    };

    return response;
}

function buildBatchCreateTagResponse(createdTags: TransactionTag[]): BatchCreateDialogResponse {
    const displayNameSourceItemMap: Record<string, string> = Object.create(null) as Record<string, string>;
    const sourceTargetMap: Record<string, string> = Object.create(null) as Record<string, string>;

    for (const item of (invalidItems.value || [])) {
        displayNameSourceItemMap[item.name] = item.value;
    }

    for (const tag of createdTags) {
        const sourceItem = displayNameSourceItemMap[tag.name];

        if (!isDefined(sourceItem)) {
            continue;
        }

        sourceTargetMap[sourceItem] = tag.id;
    }

    const response: BatchCreateDialogResponse = {
        sourceTargetMap: sourceTargetMap
    };

    return response;
}

function open(options: { type: BatchCreateDialogDataType, invalidItems?: NameValue[], accountCurrencies?: Record<string, string> }): Promise<BatchCreateDialogResponse | null> {
    type.value = options.type;
    invalidItems.value = options.invalidItems;
    accountCurrencies.value = options.accountCurrencies || {};
    selectedValues.value = getAllBatchCreateItemValues(options.invalidItems || []);

    showState.value = true;

    return new Promise((resolve) => {
        resolveFunc = resolve;
    });
}

function selectAllItems(): void {
    selectedValues.value = getAllBatchCreateItemValues(invalidItems.value || []);
}

function selectNoneItems(): void {
    selectedValues.value = [];
}

function selectInvertItems(): void {
    selectedValues.value = getInvertedBatchCreateItemValues(invalidItems.value || [], selectedValues.value);
}

function getCategoryTypeAndPrimaryName(): { categoryType: CategoryType, primaryCategoryName: string } {
    if (type.value === 'incomeCategory') {
        return {
            categoryType: CategoryType.Income,
            primaryCategoryName: tt('Default Income Category')
        };
    } else if (type.value === 'transferCategory') {
        return {
            categoryType: CategoryType.Transfer,
            primaryCategoryName: tt('Default Transfer Category')
        };
    }

    return {
        categoryType: CategoryType.Expense,
        primaryCategoryName: tt('Default Expense Category')
    };
}

async function createSelectedCategories(): Promise<BatchCreateDialogResponse> {
    await transactionCategoriesStore.loadAllCategories({ force: false });

    const { categoryType, primaryCategoryName } = getCategoryTypeAndPrimaryName();
    const selectedItems = getSelectedBatchCreateItems(invalidItems.value || [], selectedValues.value);
    const sourceTargetMap: Record<string, string> = Object.create(null) as Record<string, string>;
    const pendingItems: NameValue[] = [];
    const categories = transactionCategoriesStore.allTransactionCategories[categoryType] || [];
    let primaryCategory: TransactionCategory | undefined;

    for (const category of categories) {
        if ((!category.parentId || category.parentId === '0') && category.name === primaryCategoryName) {
            primaryCategory = category;
        }

        for (const subCategory of (category.subCategories || [])) {
            const matchingItem = selectedItems.find(item => item.name === subCategory.name);

            if (matchingItem) {
                sourceTargetMap[matchingItem.value] = subCategory.id;
            }
        }
    }

    if (!primaryCategory) {
        primaryCategory = categories.find(category => category.icon === AUTOMATICALLY_CREATED_CATEGORY_ICON_ID && category.color === DEFAULT_CATEGORY_COLOR);
    }

    for (const item of selectedItems) {
        if (!sourceTargetMap[item.value]) {
            pendingItems.push(item);
        }
    }

    if (pendingItems.length < 1) {
        return { sourceTargetMap };
    }

    if (primaryCategory) {
        for (const item of pendingItems) {
            const category = TransactionCategory.createNewCategory(categoryType, primaryCategory.id);
            category.name = item.name;
            category.icon = AUTOMATICALLY_CREATED_CATEGORY_ICON_ID;
            const createdCategory = await transactionCategoriesStore.saveCategory({
                category,
                isEdit: false,
                clientSessionId: generateRandomUUID()
            });
            sourceTargetMap[item.value] = createdCategory.id;
        }
    } else {
        const subCategories: TransactionCategoryCreateRequest[] = [];

        for (const item of pendingItems) {
            const category = TransactionCategory.createNewCategory(categoryType);
            category.name = item.name;
            category.icon = AUTOMATICALLY_CREATED_CATEGORY_ICON_ID;
            subCategories.push(category.toCreateRequest(''));
        }

        const submitCategories: TransactionCategoryCreateWithSubCategories[] = [{
            name: primaryCategoryName,
            type: categoryType,
            icon: AUTOMATICALLY_CREATED_CATEGORY_ICON_ID,
            color: DEFAULT_CATEGORY_COLOR,
            subCategories
        }];
        const response = await transactionCategoriesStore.addCategories({ categories: submitCategories });
        Object.assign(sourceTargetMap, buildBatchCreateCategoryResponse(response).sourceTargetMap);
    }

    await transactionCategoriesStore.loadAllCategories({ force: false });
    return { sourceTargetMap };
}

async function createSelectedAccounts(): Promise<BatchCreateDialogResponse> {
    const sourceTargetMap: Record<string, string> = Object.create(null) as Record<string, string>;
    const selectedItems = getSelectedBatchCreateItems(invalidItems.value || [], selectedValues.value);

    for (const item of selectedItems) {
        const currency = accountCurrencies.value[item.value];

        if (!currency) {
            throw { message: 'Unable to add account' };
        }

        const existingAccount = accountsStore.allPlainAccounts.find(account => account.name === item.name && account.currency === currency);

        if (existingAccount) {
            sourceTargetMap[item.value] = existingAccount.id;
            continue;
        }

        const account = Account.createNewAccount(AccountCategory.Default, currency, getCurrentUnixTime());
        account.name = item.name;
        const createdAccount = await accountsStore.saveAccount({
            account,
            subAccounts: [],
            isEdit: false,
            clientSessionId: generateRandomUUID()
        });
        sourceTargetMap[item.value] = createdAccount.id;
    }

    return { sourceTargetMap };
}

async function createSelectedTags(): Promise<BatchCreateDialogResponse> {
    const submitTags: TransactionTagCreateRequest[] = [];
    const selectedItems = getSelectedBatchCreateItems(invalidItems.value || [], selectedValues.value);

    for (const item of selectedItems) {
        const tag = TransactionTag.createNewTag(item.name, DEFAULT_TAG_GROUP_ID);
        submitTags.push(tag.toCreateRequest());
    }

    const response = await transactionTagsStore.addTags({
        tags: submitTags,
        groupId: DEFAULT_TAG_GROUP_ID,
        skipExists: true
    });
    await transactionTagsStore.loadAllTags({ force: false });
    return buildBatchCreateTagResponse(response);
}

async function confirm(): Promise<void> {
    if (submitting.value || selectedValues.value.length < 1) {
        return;
    }

    submitting.value = true;

    try {
        let response: BatchCreateDialogResponse;

        if (type.value === 'expenseCategory' || type.value === 'incomeCategory' || type.value === 'transferCategory') {
            response = await createSelectedCategories();
        } else if (type.value === 'account') {
            response = await createSelectedAccounts();
        } else if (type.value === 'tag') {
            response = await createSelectedTags();
        } else {
            submitting.value = false;
            return;
        }

        submitting.value = false;
        showState.value = false;
        resolveFunc?.(response);
        resolveFunc = null;
    } catch (error) {
        submitting.value = false;
        const batchCreateError = error as BatchCreateError;

        if (!batchCreateError.processed) {
            snackbar.value?.showError(batchCreateError);
        }
    }
}

function onShowStateChanged(value: boolean): void {
    if (value) {
        showState.value = true;
    } else if (!submitting.value) {
        cancel();
    }
}

function cancel(): void {
    resolveFunc?.(null);
    resolveFunc = null;
    showState.value = false;
}

defineExpose({
    open
});
</script>
