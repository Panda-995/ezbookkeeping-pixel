<template>
    <div class="studio-dashboard">
        <header class="studio-dashboard-heading">
            <div>
                <span class="studio-overline">
                    <i aria-hidden="true"></i>
                    {{ displayDateRange?.thisMonth?.displayTime }} ·
                    {{ tt("Overview") }}
                </span>
                <h2>{{ tt("Asset Summary") }}</h2>
                <p>
                    {{
                        tt("format.misc.youHaveAccounts", {
                            count: displayAccountCount,
                        })
                    }}
                </p>
            </div>
            <v-btn
                class="studio-refresh"
                variant="outlined"
                :loading="loadingOverview"
                @click="reload(true)"
            >
                <v-icon :icon="mdiRefresh" start aria-hidden="true" />
                {{ tt("Refresh") }}
            </v-btn>
        </header>

        <section class="studio-bento" :aria-label="tt('Asset Summary')">
            <article
                class="studio-net-card"
                :class="{ 'is-loading': loadingOverview }"
            >
                <div class="studio-card-index">
                    <span>01 / BALANCE</span>
                    <button
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
                            size="20"
                            aria-hidden="true"
                        />
                    </button>
                </div>

                <div class="studio-net-value">
                    <span>{{ tt("Net assets") }}</span>
                    <strong
                        v-if="
                            !loadingOverview ||
                            (allAccounts && allAccounts.length)
                        "
                    >
                        {{ netAssets }}
                    </strong>
                    <v-skeleton-loader v-else width="62%" type="text" />
                </div>

                <div class="studio-balance-split">
                    <div>
                        <small>{{ tt("Total assets") }}</small>
                        <strong>{{ totalAssets }}</strong>
                    </div>
                    <div>
                        <small>{{ tt("Total liabilities") }}</small>
                        <strong>{{ totalLiabilities }}</strong>
                    </div>
                </div>
                <div class="studio-pixel-ruler" aria-hidden="true">
                    <i v-for="index in 16" :key="index"></i>
                </div>
            </article>

            <article
                class="studio-month-card"
                :class="{ 'is-loading': loadingOverview }"
            >
                <div class="studio-card-index">
                    <span>02 / {{ tt("This Month") }}</span>
                    <span aria-hidden="true">↗</span>
                </div>
                <div class="studio-month-expense">
                    <span>{{ tt("Expense") }}</span>
                    <strong>
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
                <div class="studio-month-income">
                    <span>{{ tt("Monthly income") }}</span>
                    <strong>
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
                    class="studio-arrow-link"
                    :to="`/transaction/list?${overviewStore.getTransactionListPageParams({ dateType: DateRange.ThisMonth.type })}`"
                >
                    {{ tt("View Details") }}
                    <span aria-hidden="true">→</span>
                </router-link>
            </article>

            <aside class="studio-bento-note" aria-hidden="true">
                <span>KEEP THE</span>
                <strong>NUMBERS</strong>
                <span>IN FOCUS.</span>
                <i></i>
            </aside>
        </section>

        <section class="studio-period-section">
            <div class="studio-section-heading">
                <div>
                    <span class="studio-overline">03 / CASH FLOW</span>
                    <h3>{{ tt("Transaction Data") }}</h3>
                </div>
                <router-link
                    class="studio-arrow-link"
                    to="/transaction/list?pageType=0&dateType=7"
                >
                    {{ tt("Transaction Details") }}
                    <span aria-hidden="true">→</span>
                </router-link>
            </div>

            <div class="studio-period-grid">
                <router-link
                    v-for="(period, index) in [
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
                    class="studio-period-card"
                    :to="`/transaction/list?${overviewStore.getTransactionListPageParams({ dateType: period.dateType })}`"
                >
                    <span class="studio-period-number">0{{ index + 1 }}</span>
                    <div class="studio-period-title">
                        <strong>{{ period.label }}</strong>
                        <small>{{ period.datetime }}</small>
                    </div>
                    <dl>
                        <div>
                            <dt>{{ tt("Income") }}</dt>
                            <dd class="text-income">
                                {{
                                    period.data && period.data.valid
                                        ? getDisplayIncomeAmount(period.data)
                                        : "-"
                                }}
                            </dd>
                        </div>
                        <div>
                            <dt>{{ tt("Expense") }}</dt>
                            <dd class="text-expense">
                                {{
                                    period.data && period.data.valid
                                        ? getDisplayExpenseAmount(period.data)
                                        : "-"
                                }}
                            </dd>
                        </div>
                    </dl>
                    <span class="studio-period-arrow" aria-hidden="true"
                        >↗</span
                    >
                </router-link>
            </div>
        </section>

        <section class="studio-chart-panel">
            <div class="studio-section-heading">
                <div>
                    <span class="studio-overline">04 / 12 MONTHS</span>
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

        <snack-bar ref="snackbar" />
    </div>
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
