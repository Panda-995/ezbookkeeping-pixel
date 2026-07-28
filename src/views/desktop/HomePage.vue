<template>
    <main class="pixel-register-dashboard">
        <header class="pixel-page-heading">
            <div>
                <div class="pixel-kicker">
                    <span class="pixel-status-dot" aria-hidden="true"></span>
                    {{ tt("Overview") }} ·
                    {{ displayDateRange?.thisMonth?.displayTime }}
                </div>
                <h2>{{ tt("Asset Summary") }}</h2>
                <p>
                    {{
                        tt("format.misc.youHaveAccounts", {
                            count: displayAccountCount,
                        })
                    }}
                </p>
            </div>
            <div class="pixel-page-actions">
                <v-btn
                    variant="outlined"
                    :loading="loadingOverview"
                    @click="reload(true)"
                >
                    <v-icon :icon="mdiRefresh" start />
                    {{ tt("Refresh") }}
                </v-btn>
                <v-btn color="primary" @click="addTransaction">
                    {{ tt("Add Transaction") }}
                </v-btn>
            </div>
        </header>

        <section class="pixel-balance-board" aria-label="Asset summary">
            <article
                class="pixel-balance-primary"
                :class="{ 'is-loading': loadingOverview }"
            >
                <div class="pixel-panel-index">REGISTER / 001</div>
                <div class="pixel-balance-label">
                    {{ tt("Net assets") }}
                    <button
                        class="pixel-icon-button"
                        type="button"
                        :aria-label="
                            showAmountInHomePage
                                ? tt('Hide Amount')
                                : tt('Show Amount')
                        "
                        @click="showAmountInHomePage = !showAmountInHomePage"
                    >
                        <v-icon
                            :icon="
                                showAmountInHomePage
                                    ? mdiEyeOffOutline
                                    : mdiEyeOutline
                            "
                            size="19"
                        />
                    </button>
                </div>
                <div
                    class="pixel-balance-value"
                    v-if="
                        !loadingOverview || (allAccounts && allAccounts.length)
                    "
                >
                    {{ netAssets }}
                </div>
                <v-skeleton-loader width="240" type="text" v-else />
                <div class="pixel-balance-grid">
                    <div>
                        <span>{{ tt("Total assets") }}</span>
                        <strong>{{ totalAssets }}</strong>
                    </div>
                    <div>
                        <span>{{ tt("Total liabilities") }}</span>
                        <strong class="text-expense">{{
                            totalLiabilities
                        }}</strong>
                    </div>
                </div>
                <div class="pixel-board-marks" aria-hidden="true">
                    <i></i><i></i><i></i><i></i><i></i><i></i><i></i><i></i>
                </div>
            </article>

            <article
                class="pixel-cashflow-ticket"
                :class="{ 'is-loading': loadingOverview }"
            >
                <div class="pixel-ticket-header">
                    <div>
                        <span>{{ tt("This Month") }}</span>
                        <strong
                            >{{ displayDateRange?.thisMonth?.startTime }} —
                            {{ displayDateRange?.thisMonth?.endTime }}</strong
                        >
                    </div>
                    <span class="pixel-ticket-number"
                        >NO. {{ displayAccountCount.padStart(2, "0") }}</span
                    >
                </div>
                <div class="pixel-ticket-total">
                    <span>{{ tt("Expense") }}</span>
                    <strong class="text-expense">
                        {{
                            transactionOverview.thisMonth &&
                            transactionOverview.thisMonth.valid
                                ? getDisplayExpenseAmount(
                                      transactionOverview.thisMonth,
                                  )
                                : "-"
                        }}
                    </strong>
                </div>
                <div class="pixel-ticket-row">
                    <span>{{ tt("Monthly income") }}</span>
                    <strong class="text-income">
                        {{
                            transactionOverview.thisMonth &&
                            transactionOverview.thisMonth.valid
                                ? getDisplayIncomeAmount(
                                      transactionOverview.thisMonth,
                                  )
                                : "-"
                        }}
                    </strong>
                </div>
                <router-link
                    class="pixel-text-link"
                    :to="`/transaction/list?${overviewStore.getTransactionListPageParams({ dateType: DateRange.ThisMonth.type })}`"
                >
                    {{ tt("View Details") }} →
                </router-link>
            </article>
        </section>

        <section class="pixel-register-section">
            <div class="pixel-section-heading">
                <div>
                    <span class="pixel-panel-index">CASHFLOW / PERIODS</span>
                    <h3>{{ tt("Transaction Data") }}</h3>
                </div>
                <router-link
                    class="pixel-text-link"
                    to="/transaction/list?pageType=0&dateType=7"
                >
                    {{ tt("Transaction Details") }} →
                </router-link>
            </div>

            <div class="pixel-period-register">
                <router-link
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
                    class="pixel-period-row"
                    :to="`/transaction/list?${overviewStore.getTransactionListPageParams({ dateType: period.dateType })}`"
                >
                    <span class="pixel-period-glyph" aria-hidden="true"></span>
                    <span class="pixel-period-name">
                        <strong>{{ period.label }}</strong>
                        <small>{{ period.datetime }}</small>
                    </span>
                    <span class="pixel-period-amount text-income">
                        <small>{{ tt("Income") }}</small>
                        <strong>{{
                            period.data && period.data.valid
                                ? getDisplayIncomeAmount(period.data)
                                : "-"
                        }}</strong>
                    </span>
                    <span class="pixel-period-amount text-expense">
                        <small>{{ tt("Expense") }}</small>
                        <strong>{{
                            period.data && period.data.valid
                                ? getDisplayExpenseAmount(period.data)
                                : "-"
                        }}</strong>
                    </span>
                    <span class="pixel-period-arrow" aria-hidden="true">→</span>
                </router-link>
            </div>
        </section>

        <section class="pixel-register-section pixel-chart-section">
            <div class="pixel-section-heading">
                <div>
                    <span class="pixel-panel-index">TREND / 12 MONTHS</span>
                    <h3>{{ tt("Statistics & Analysis") }}</h3>
                </div>
            </div>
            <monthly-income-and-expense-card
                :data="monthlyIncomeAndExpenseData"
                :is-dark-mode="isDarkMode"
                :loading="loadingOverview"
                :disabled="loadingOverview"
                :enable-click-item="true"
                @click="clickMonthlyIncomeOrExpense"
            />
        </section>
    </main>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from "@/components/desktop/SnackBar.vue";
import MonthlyIncomeAndExpenseCard, {
    type MonthlyIncomeAndExpenseCardClickEvent,
} from "./overview/cards/MonthlyIncomeAndExpenseCard.vue";

import { ref, computed, useTemplateRef } from "vue";
import { useRouter } from "vue-router";
import { useTheme } from "vuetify";

import { useI18n } from "@/locales/helpers.ts";
import { useHomePageBase } from "@/views/base/HomePageBase.ts";

import { useAccountsStore } from "@/stores/account.ts";
import { useTransactionCategoriesStore } from "@/stores/transactionCategory.ts";
import { useOverviewStore } from "@/stores/overview.ts";
import { useDesktopPageStore } from "@/stores/desktopPage.ts";

import { DateRange } from "@/core/datetime.ts";
import { ThemeType } from "@/core/theme.ts";
import {
    type TransactionMonthlyIncomeAndExpenseData,
    LATEST_12MONTHS_TRANSACTION_AMOUNTS_REQUEST_TYPES,
} from "@/models/transaction.ts";

import { BIG_DECIMAL_ZERO } from "@/lib/numeral.ts";
import {
    getUnixTimeBeforeUnixTime,
    getUnixTimeAfterUnixTime,
} from "@/lib/datetime.ts";
import { isUserLogined, isUserUnlocked } from "@/lib/userstate.ts";

import { mdiRefresh, mdiEyeOutline, mdiEyeOffOutline } from "@mdi/js";

type SnackBarType = InstanceType<typeof SnackBar>;

const router = useRouter();
const theme = useTheme();

const { tt, formatNumberToLocalizedNumerals } = useI18n();
const {
    showAmountInHomePage,
    allAccounts,
    netAssets,
    totalAssets,
    totalLiabilities,
    displayDateRange,
    transactionOverview,
    getDisplayIncomeAmount,
    getDisplayExpenseAmount,
} = useHomePageBase();

const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const overviewStore = useOverviewStore();
const desktopPageStore = useDesktopPageStore();

const snackbar = useTemplateRef<SnackBarType>("snackbar");

const loadingOverview = ref<boolean>(true);

const isDarkMode = computed<boolean>(
    () => theme.global.name.value === ThemeType.Dark,
);

const displayAccountCount = computed<string>(() =>
    formatNumberToLocalizedNumerals(allAccounts.value?.length ?? 0),
);

function clickMonthlyIncomeOrExpense(
    e: MonthlyIncomeAndExpenseCardClickEvent,
): void {
    const minTime = e.monthStartTime;
    const maxTime = getUnixTimeBeforeUnixTime(
        getUnixTimeAfterUnixTime(minTime, 1, "months"),
        1,
        "seconds",
    );
    const type = e.transactionType;

    router.push(
        `/transaction/list?${overviewStore.getTransactionListPageParams({
            type: type,
            dateType: DateRange.Custom.type,
            minTime: minTime,
            maxTime: maxTime,
        })}`,
    );
}

function addTransaction(): void {
    router.push("/transaction/list?pageType=0&dateType=7").then(() => {
        desktopPageStore.setShowAddTransactionDialogInTransactionList();
    });
}

const monthlyIncomeAndExpenseData = computed<
    TransactionMonthlyIncomeAndExpenseData[]
>(() => {
    const data: TransactionMonthlyIncomeAndExpenseData[] = [];

    if (
        !transactionOverview.value ||
        !transactionOverview.value.thisMonth ||
        !transactionOverview.value.thisMonth.valid
    ) {
        return data;
    }

    for (const amountRequestType of LATEST_12MONTHS_TRANSACTION_AMOUNTS_REQUEST_TYPES) {
        const dateRange = overviewStore.transactionDataRange[amountRequestType];

        if (!dateRange) {
            continue;
        }

        const item = transactionOverview.value[amountRequestType];

        data.push({
            monthStartTime: dateRange.startTime,
            incomeAmount: item?.incomeAmount || BIG_DECIMAL_ZERO,
            expenseAmount: item?.expenseAmount || BIG_DECIMAL_ZERO,
            incompleteIncomeAmount: item ? item.incompleteIncomeAmount : true,
            incompleteExpenseAmount: item ? item.incompleteExpenseAmount : true,
        });
    }

    return data;
});

function reload(force: boolean): void {
    loadingOverview.value = true;

    const promises = [
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false }),
        overviewStore.loadTransactionOverview({
            force: force,
            loadLast11Months: true,
        }),
    ];

    Promise.all(promises)
        .then(() => {
            loadingOverview.value = false;

            if (force) {
                snackbar.value?.showMessage("Data has been updated");
            }
        })
        .catch((error) => {
            loadingOverview.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
}

if (isUserLogined() && isUserUnlocked()) {
    reload(false);
}
</script>
