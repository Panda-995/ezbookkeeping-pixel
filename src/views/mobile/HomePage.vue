<template>
    <f7-page
        class="mobile-pixel-home"
        ptr
        @ptr:refresh="reload"
        @page:afterin="onPageAfterIn"
    >
        <div class="mobile-pixel-masthead">
            <div class="mobile-pixel-brandline">
                <div>
                    <span class="mobile-pixel-signal" aria-hidden="true"></span>
                    <small>{{
                        tt("Personal finance, clearly organized")
                    }}</small>
                </div>
                <f7-link icon-only @click="reload()">
                    <f7-icon f7="arrow_clockwise"></f7-icon>
                </f7-link>
            </div>

            <div class="mobile-pixel-month">
                <span>{{
                    loading
                        ? "---- / --"
                        : displayDateRange?.thisMonth?.displayTime
                }}</span>
                <button
                    type="button"
                    @click="showAmountInHomePage = !showAmountInHomePage"
                >
                    <f7-icon
                        :f7="showAmountInHomePage ? 'eye_slash' : 'eye'"
                    ></f7-icon>
                </button>
            </div>
            <span class="mobile-pixel-balance-label">{{ tt("Expense") }}</span>
            <strong class="mobile-pixel-balance text-expense">
                {{
                    loading
                        ? "—"
                        : transactionOverview.thisMonth
                          ? getDisplayExpenseAmount(
                                transactionOverview.thisMonth,
                            )
                          : "-"
                }}
            </strong>

            <div class="mobile-pixel-income">
                <span>{{ tt("Monthly income") }}</span>
                <strong class="text-income">
                    {{
                        loading
                            ? "—"
                            : transactionOverview.thisMonth
                              ? getDisplayIncomeAmount(
                                    transactionOverview.thisMonth,
                                )
                              : "-"
                    }}
                </strong>
            </div>

            <div class="mobile-pixel-dots" aria-hidden="true">
                <i v-for="index in 18" :key="index"></i>
            </div>
        </div>

        <div class="mobile-pixel-quick">
            <f7-link href="/transaction/list">
                <f7-icon f7="list_bullet_rectangle"></f7-icon>
                <span>{{ tt("Details") }}</span>
            </f7-link>
            <f7-link href="/transaction/add">
                <f7-icon f7="plus"></f7-icon>
                <span>{{ tt("Add Transaction") }}</span>
            </f7-link>
            <f7-link href="/account/list">
                <f7-icon f7="creditcard"></f7-icon>
                <span>{{ tt("Accounts") }}</span>
            </f7-link>
        </div>

        <section
            class="mobile-pixel-register"
            :class="{ 'skeleton-text': loading }"
        >
            <header>
                <div>
                    <h1>{{ tt("Overview") }}</h1>
                </div>
                <f7-link href="/transaction/list"
                    >{{ tt("View Details") }} →</f7-link
                >
            </header>

            <f7-link
                v-for="period in [
                    {
                        key: 'today',
                        label: tt('Today'),
                        datetime: displayDateRange?.today?.displayTime,
                        data: transactionOverview.today,
                        dateType: DateRange.Today.type,
                    },
                    {
                        key: 'week',
                        label: tt('This Week'),
                        datetime: `${displayDateRange?.thisWeek?.startTime || ''} — ${displayDateRange?.thisWeek?.endTime || ''}`,
                        data: transactionOverview.thisWeek,
                        dateType: DateRange.ThisWeek.type,
                    },
                    {
                        key: 'month',
                        label: tt('This Month'),
                        datetime: `${displayDateRange?.thisMonth?.startTime || ''} — ${displayDateRange?.thisMonth?.endTime || ''}`,
                        data: transactionOverview.thisMonth,
                        dateType: DateRange.ThisMonth.type,
                    },
                    {
                        key: 'year',
                        label: tt('This Year'),
                        datetime: displayDateRange?.thisYear?.displayTime,
                        data: transactionOverview.thisYear,
                        dateType: DateRange.ThisYear.type,
                    },
                ]"
                :key="period.key"
                class="mobile-pixel-period"
                :href="`/transaction/list?${overviewStore.getTransactionListPageParams({ dateType: period.dateType })}`"
            >
                <span
                    class="mobile-pixel-period-mark"
                    aria-hidden="true"
                ></span>
                <span class="mobile-pixel-period-name">
                    <strong>{{ period.label }}</strong>
                    <small>{{ period.datetime }}</small>
                </span>
                <span class="mobile-pixel-period-values">
                    <strong class="text-income">
                        {{
                            period.data && period.data.valid
                                ? getDisplayIncomeAmount(period.data)
                                : "-"
                        }}
                    </strong>
                    <strong class="text-expense">
                        {{
                            period.data && period.data.valid
                                ? getDisplayExpenseAmount(period.data)
                                : "-"
                        }}
                    </strong>
                </span>
                <span class="mobile-pixel-period-arrow">→</span>
            </f7-link>
        </section>

        <div class="mobile-pixel-spacer"></div>

        <f7-toolbar tabbar icons bottom class="mobile-pixel-tabbar">
            <f7-link class="link" href="/transaction/list">
                <f7-icon f7="square_list"></f7-icon>
                <span class="tabbar-label">{{ tt("Details") }}</span>
            </f7-link>
            <f7-link class="link" href="/account/list">
                <f7-icon f7="creditcard"></f7-icon>
                <span class="tabbar-label">{{ tt("Accounts") }}</span>
            </f7-link>
            <f7-link
                id="homepage-add-button"
                class="link dragenabled mobile-pixel-add"
                href="/transaction/add"
                @taphold="openTransactionTemplatePopover"
            >
                <span><f7-icon f7="plus"></f7-icon></span>
            </f7-link>
            <f7-link class="link" href="/statistic/transaction">
                <f7-icon f7="chart_pie"></f7-icon>
                <span class="tabbar-label">{{ tt("Statistics") }}</span>
            </f7-link>
            <f7-link class="link" href="/settings">
                <f7-icon f7="gear_alt"></f7-icon>
                <span class="tabbar-label">{{ tt("Settings") }}</span>
            </f7-link>
        </f7-toolbar>

        <f7-popover
            class="template-popover-menu"
            target-el="#homepage-add-button"
            v-model:opened="showTransactionTemplatePopover"
        >
            <f7-list
                dividers
                v-if="
                    isTransactionFromAITextRecognitionEnabled() ||
                    isTransactionFromAIImageRecognitionEnabled() ||
                    (allTransactionTemplates && allTransactionTemplates.length)
                "
            >
                <f7-list-item
                    key="AIClipboardTextRecognition"
                    link="#"
                    no-chevron
                    popover-close
                    :title="tt('AI Clipboard Text Recognition')"
                    @click="addByRecognizingClipboardText"
                    v-if="isTransactionFromAITextRecognitionEnabled()"
                >
                    <template #media
                        ><f7-icon f7="wand_stars"></f7-icon
                    ></template>
                </f7-list-item>
                <f7-list-item
                    key="AIImageRecognition"
                    link="#"
                    no-chevron
                    popover-close
                    :title="tt('AI Image Recognition')"
                    @click="showAIReceiptImageRecognitionSheet = true"
                    v-if="isTransactionFromAIImageRecognitionEnabled()"
                >
                    <template #media
                        ><f7-icon f7="wand_stars"></f7-icon
                    ></template>
                </f7-list-item>
                <f7-list-item
                    popover-close
                    :key="template.id"
                    :title="template.name"
                    :link="'/transaction/add?templateId=' + template.id"
                    v-for="template in allTransactionTemplates"
                >
                    <template #media
                        ><f7-icon f7="doc_plaintext"></f7-icon
                    ></template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <a-i-image-recognition-sheet
            ref="aiImageRecognitionSheet"
            v-model:show="showAIReceiptImageRecognitionSheet"
            @recognition:change="onReceiptRecognitionChanged"
        />
    </f7-page>
</template>

<script setup lang="ts">
import AIImageRecognitionSheet, {
    type AIImageRecognitionResult,
} from "@/components/mobile/AIImageRecognitionSheet.vue";

import { ref, computed, useTemplateRef } from "vue";
import type { Router } from "framework7/types";

import { useI18n } from "@/locales/helpers.ts";
import { useI18nUIComponents, isiOS } from "@/lib/ui/mobile.ts";
import { useHomePageBase } from "@/views/base/HomePageBase.ts";

import { useSettingsStore } from "@/stores/setting.ts";
import { useAccountsStore } from "@/stores/account.ts";
import { useTransactionCategoriesStore } from "@/stores/transactionCategory.ts";
import { useTransactionTemplatesStore } from "@/stores/transactionTemplate.ts";
import { useOverviewStore } from "@/stores/overview.ts";

import { DateRange } from "@/core/datetime.ts";
import { TemplateType } from "@/core/template.ts";
import { TransactionTemplate } from "@/models/transaction_template.ts";

import { isFunction } from "@/lib/common.ts";
import { isUserLogined, isUserUnlocked } from "@/lib/userstate.ts";
import { getShareCacheImageBlob } from "@/lib/cache.ts";
import {
    isTransactionFromAITextRecognitionEnabled,
    isTransactionFromAIImageRecognitionEnabled,
} from "@/lib/server_settings.ts";
import logger from "@/lib/logger.ts";

type AIImageRecognitionSheetType = InstanceType<typeof AIImageRecognitionSheet>;

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt } = useI18n();
const { showToast } = useI18nUIComponents();

const {
    showAmountInHomePage,
    displayDateRange,
    transactionOverview,
    getDisplayIncomeAmount,
    getDisplayExpenseAmount,
} = useHomePageBase();

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTemplatesStore = useTransactionTemplatesStore();
const overviewStore = useOverviewStore();

const aiImageRecognitionSheet = useTemplateRef<AIImageRecognitionSheetType>(
    "aiImageRecognitionSheet",
);

const loading = ref<boolean>(true);
const showTransactionTemplatePopover = ref<boolean>(false);
const showAIReceiptImageRecognitionSheet = ref<boolean>(false);

const allTransactionTemplates = computed<TransactionTemplate[]>(() => {
    const allTemplates = transactionTemplatesStore.allVisibleTemplates;
    return allTemplates[TemplateType.Normal.type] || [];
});

function openTransactionTemplatePopover(): void {
    if (
        isTransactionFromAIImageRecognitionEnabled() ||
        (allTransactionTemplates.value && allTransactionTemplates.value.length)
    ) {
        showTransactionTemplatePopover.value = true;
    }
}

function init(): void {
    if (isUserLogined() && isUserUnlocked()) {
        loading.value = true;

        const promises = [
            getShareCacheImageBlob(),
            accountsStore.loadAllAccounts({ force: false }),
            transactionCategoriesStore.loadAllCategories({ force: false }),
            transactionTemplatesStore.loadAllTemplates({
                templateType: TemplateType.Normal.type,
                force: false,
            }),
            overviewStore.loadTransactionOverview({ force: false }),
        ];

        Promise.all(promises)
            .then((responses) => {
                if (responses[0] && responses[0] instanceof Blob) {
                    aiImageRecognitionSheet.value?.loadImage(responses[0]);
                    showAIReceiptImageRecognitionSheet.value = true;
                }

                loading.value = false;
            })
            .catch((error) => {
                loading.value = false;

                if (!error.processed) {
                    showToast(error.message || error);
                }
            });
    }
}

function reload(done?: () => void): void {
    const force = !!done;

    overviewStore
        .loadTransactionOverview({
            force: force,
        })
        .then(() => {
            done?.();

            if (force) {
                showToast("Data has been updated");
            }
        })
        .catch((error) => {
            done?.();

            if (!error.processed) {
                showToast(error.message || error);
            }
        });
}

function addByRecognizingClipboardText(): void {
    if (
        navigator.clipboard &&
        isFunction(navigator.clipboard.readText) &&
        !isiOS()
    ) {
        navigator.clipboard
            .readText()
            .then((text) => {
                const clipboardText = text && text.trim() ? text.trim() : "";
                props.f7router.navigate("/transaction/add", {
                    props: {
                        autoRecognizeClipboardText: clipboardText,
                    },
                });
            })
            .catch((error) => {
                logger.error("failed to read clipboard", error);
                props.f7router.navigate("/transaction/add", {
                    props: {
                        autoRecognizeClipboardText: "",
                    },
                });
            });
    } else {
        props.f7router.navigate("/transaction/add", {
            props: {
                autoRecognizeClipboardText: "",
            },
        });
    }
}

function onReceiptRecognitionChanged(result: AIImageRecognitionResult): void {
    const recognizedResponse = result.response;
    const autoUploadRecognizedImage =
        settingsStore.appSettings.autoUploadTransactionPictureForAIRecognition;
    const params: string[] = [];

    if (recognizedResponse.type) {
        params.push(`type=${recognizedResponse.type}`);
    }

    if (recognizedResponse.time) {
        params.push(`time=${recognizedResponse.time}`);
    }

    if (recognizedResponse.categoryId) {
        params.push(`categoryId=${recognizedResponse.categoryId}`);
    }

    if (recognizedResponse.sourceAccountId) {
        params.push(`accountId=${recognizedResponse.sourceAccountId}`);
    }

    if (recognizedResponse.destinationAccountId) {
        params.push(
            `destinationAccountId=${recognizedResponse.destinationAccountId}`,
        );
    }

    if (recognizedResponse.sourceAmount) {
        params.push(`amount=${recognizedResponse.sourceAmount}`);
    }

    if (recognizedResponse.destinationAmount) {
        params.push(
            `destinationAmount=${recognizedResponse.destinationAmount}`,
        );
    }

    if (recognizedResponse.tagIds) {
        params.push(`tagIds=${recognizedResponse.tagIds.join(",")}`);
    }

    if (recognizedResponse.comment) {
        params.push(
            `comment=${encodeURIComponent(recognizedResponse.comment)}`,
        );
    }

    params.push(`noTransactionDraft=true`);

    props.f7router.navigate(`/transaction/add?${params.join("&")}`, {
        props: {
            autoUploadPicture: autoUploadRecognizedImage
                ? result.imageFile
                : undefined,
        },
    });
}

function onPageAfterIn(): void {
    if (!loading.value) {
        reload();
    }
}

init();
</script>
