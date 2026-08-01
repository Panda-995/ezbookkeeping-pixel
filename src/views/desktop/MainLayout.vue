<template>
    <div
        class="layout-wrapper studio-shell layout-content-width-fluid"
        :class="{ 'layout-overlay-nav': mdAndDown }"
    >
        <aside
            class="layout-vertical-nav studio-dock"
            :class="{
                visible: showVerticalOverlayMenu,
                scrolled: isVerticalNavScrolled,
                'overlay-nav': mdAndDown,
            }"
        >
            <div class="studio-dock-header">
                <router-link to="/" class="studio-brand">
                    <span class="studio-brand-mark">
                        <img
                            alt=""
                            width="26"
                            height="26"
                            :src="APPLICATION_LOGO_PATH"
                        />
                    </span>
                    <span class="studio-brand-copy">
                        <small>LEDGER / STUDIO</small>
                        <strong>{{ tt("global.app.title") }}</strong>
                    </span>
                </router-link>

                <button
                    class="studio-dock-close d-lg-none"
                    type="button"
                    :aria-label="tt('Close')"
                    @click="showVerticalOverlayMenu = false"
                >
                    <span aria-hidden="true">×</span>
                </button>
            </div>

            <button
                class="studio-new-entry"
                type="button"
                @click="showAddDialogInTransactionListPage"
            >
                <span class="studio-new-entry-icon" aria-hidden="true">＋</span>
                <span>
                    <strong>{{ tt("Add Transaction") }}</strong>
                    <small>QUICK ENTRY · 01</small>
                </span>
            </button>

            <perfect-scrollbar
                tag="nav"
                class="studio-dock-nav"
                :options="{ wheelPropagation: false }"
                :aria-label="tt('Main Menu')"
                @ps-scroll-y="handleNavScroll"
            >
                <ul class="nav-items">
                    <li class="nav-link home-link">
                        <router-link to="/">
                            <v-icon :icon="mdiHomeOutline" aria-hidden="true" />
                            <span>{{ tt("Overview") }}</span>
                            <small>01</small>
                        </router-link>
                    </li>

                    <li class="nav-section-title">
                        <span>{{ tt("Transaction Data") }}</span>
                    </li>
                    <li class="nav-link">
                        <router-link
                            to="/transaction/list?pageType=0&dateType=7"
                        >
                            <v-icon
                                :icon="mdiListBoxOutline"
                                aria-hidden="true"
                            />
                            <span>{{ tt("Transaction Details") }}</span>
                            <small>02</small>
                        </router-link>
                    </li>
                    <li class="nav-link">
                        <router-link to="/statistics/transaction">
                            <v-icon
                                :icon="mdiChartPieOutline"
                                aria-hidden="true"
                            />
                            <span>{{ tt("Statistics & Analysis") }}</span>
                            <small>03</small>
                        </router-link>
                    </li>
                    <li class="nav-link">
                        <router-link to="/insights/explorer">
                            <v-icon
                                :icon="mdiCompassOutline"
                                aria-hidden="true"
                            />
                            <span>{{ tt("Insights Explorer") }}</span>
                            <small>04</small>
                        </router-link>
                    </li>

                    <li class="nav-section-title">
                        <span>{{ tt("Basis Data") }}</span>
                    </li>
                    <li class="nav-link">
                        <router-link to="/account/list">
                            <v-icon
                                :icon="mdiCreditCardOutline"
                                aria-hidden="true"
                            />
                            <span>{{ tt("Accounts") }}</span>
                            <small>05</small>
                        </router-link>
                    </li>
                    <li class="nav-link">
                        <router-link to="/category/list">
                            <v-icon
                                :icon="mdiViewDashboardOutline"
                                aria-hidden="true"
                            />
                            <span>{{ tt("Transaction Categories") }}</span>
                            <small>06</small>
                        </router-link>
                    </li>
                    <li class="nav-link">
                        <router-link to="/tag/list">
                            <v-icon :icon="mdiTagOutline" aria-hidden="true" />
                            <span>{{ tt("Transaction Tags") }}</span>
                            <small>07</small>
                        </router-link>
                    </li>
                    <li class="nav-link">
                        <router-link to="/template/list">
                            <v-icon
                                :icon="mdiClipboardTextOutline"
                                aria-hidden="true"
                            />
                            <span>{{ tt("Transaction Templates") }}</span>
                            <small>08</small>
                        </router-link>
                    </li>
                    <li
                        class="nav-link"
                        v-if="isUserScheduledTransactionEnabled()"
                    >
                        <router-link to="/schedule/list">
                            <v-icon
                                :icon="mdiClipboardTextClockOutline"
                                aria-hidden="true"
                            />
                            <span>{{ tt("Scheduled Transactions") }}</span>
                            <small>09</small>
                        </router-link>
                    </li>

                    <li class="nav-section-title">
                        <span>{{ tt("Miscellaneous") }}</span>
                    </li>
                    <li class="nav-link">
                        <router-link to="/exchange_rates">
                            <v-icon
                                :icon="mdiSwapHorizontal"
                                aria-hidden="true"
                            />
                            <span>{{ tt("Exchange Rates Data") }}</span>
                        </router-link>
                    </li>
                    <li class="nav-link">
                        <button type="button" @click="showMobileQrCode = true">
                            <v-icon :icon="mdiCellphone" aria-hidden="true" />
                            <span>{{ tt("Use on Mobile Device") }}</span>
                        </button>
                    </li>
                    <li class="nav-link">
                        <router-link to="/about">
                            <v-icon
                                :icon="mdiInformationOutline"
                                aria-hidden="true"
                            />
                            <span>{{ tt("About") }}</span>
                        </router-link>
                    </li>
                </ul>
            </perfect-scrollbar>

            <div class="studio-dock-foot">
                <span class="studio-presence" aria-hidden="true"></span>
                <span>
                    <small>{{ tt("User") }}</small>
                    <strong>{{ currentNickName }}</strong>
                </span>
            </div>
        </aside>

        <div class="layout-content-wrapper studio-workspace">
            <header class="layout-navbar studio-topbar">
                <div class="navbar-content-container">
                    <button
                        class="studio-menu-trigger d-lg-none"
                        type="button"
                        :aria-label="tt('Open Menu')"
                        @click="showVerticalOverlayMenu = true"
                    >
                        <v-icon :icon="mdiMenu" aria-hidden="true" />
                    </button>

                    <div class="studio-page-context">
                        <small>PERSONAL LEDGER</small>
                        <h1>{{ currentPageTitle }}</h1>
                    </div>

                    <div class="studio-topbar-actions">
                        <v-btn
                            class="studio-topbar-add"
                            color="primary"
                            @click="showAddDialogInTransactionListPage"
                        >
                            <v-icon
                                :icon="mdiPlusCircle"
                                start
                                aria-hidden="true"
                            />
                            {{ tt("Add Transaction") }}
                        </v-btn>

                        <v-btn
                            class="studio-round-action"
                            variant="text"
                            :icon="true"
                            :aria-label="tt('Theme')"
                            @click="
                                currentTheme === 'light'
                                    ? (currentTheme = 'dark')
                                    : currentTheme === 'dark'
                                      ? (currentTheme = 'auto')
                                      : (currentTheme = 'light')
                            "
                        >
                            <v-icon
                                :icon="
                                    currentTheme === 'light'
                                        ? mdiWeatherSunny
                                        : currentTheme === 'dark'
                                          ? mdiWeatherNight
                                          : mdiThemeLightDark
                                "
                                aria-hidden="true"
                            />
                        </v-btn>

                        <v-btn
                            class="studio-user-menu"
                            variant="text"
                            :icon="true"
                            :aria-label="tt('User Profile')"
                        >
                            <v-avatar size="42" color="primary">
                                <v-img
                                    :src="currentUserAvatar"
                                    alt=""
                                    width="42"
                                    height="42"
                                    v-if="currentUserAvatar"
                                />
                                <v-icon
                                    :icon="mdiAccount"
                                    v-else
                                    aria-hidden="true"
                                />
                            </v-avatar>
                            <v-menu
                                activator="parent"
                                width="248"
                                location="bottom end"
                                offset="12px"
                            >
                                <v-list class="studio-profile-popover">
                                    <v-list-item
                                        :title="currentNickName"
                                        :subtitle="tt('User Profile')"
                                    >
                                        <template #prepend>
                                            <v-avatar size="38" color="primary">
                                                <v-img
                                                    :src="currentUserAvatar"
                                                    alt=""
                                                    width="38"
                                                    height="38"
                                                    v-if="currentUserAvatar"
                                                />
                                                <v-icon
                                                    :icon="mdiAccount"
                                                    v-else
                                                    aria-hidden="true"
                                                />
                                            </v-avatar>
                                        </template>
                                    </v-list-item>
                                    <v-divider class="my-2" />
                                    <v-list-item
                                        :prepend-icon="mdiAccountCogOutline"
                                        :title="tt('User Settings')"
                                        to="/user/settings"
                                    />
                                    <v-list-item
                                        :prepend-icon="mdiCogOutline"
                                        :title="tt('Application Settings')"
                                        to="/app/settings"
                                    />
                                    <v-divider class="my-2" />
                                    <v-list-item
                                        v-if="isEnableApplicationLock"
                                        :prepend-icon="mdiLockOutline"
                                        :title="tt('Lock Application')"
                                        @click="lock"
                                    />
                                    <v-list-item
                                        :disabled="logouting"
                                        :prepend-icon="mdiLogout"
                                        :title="tt('Log Out')"
                                        @click="logout"
                                    />
                                </v-list>
                            </v-menu>
                        </v-btn>
                    </div>
                </div>
            </header>

            <main
                id="studio-main"
                class="layout-page-content studio-page-stage"
                tabindex="-1"
            >
                <div class="page-content-container">
                    <router-view v-slot="{ Component }">
                        <transition name="studio-route" mode="out-in">
                            <component
                                :is="Component"
                                :key="currentRoutePath"
                            />
                        </transition>
                    </router-view>
                </div>
            </main>
        </div>

        <switch-to-mobile-dialog v-model:show="showMobileQrCode" />

        <button
            class="layout-overlay"
            :class="{ visible: showVerticalOverlayMenu }"
            type="button"
            :aria-label="tt('Close')"
            @click="showVerticalOverlayMenu = false"
        ></button>

        <v-overlay
            class="justify-center align-center"
            :persistent="true"
            v-model="showLoading"
        >
            <v-progress-circular indeterminate />
        </v-overlay>

        <snack-bar ref="snackbar" />
    </div>
</template>

<script setup lang="ts">
import SnackBar from "@/components/desktop/SnackBar.vue";

import { ref, computed, useTemplateRef } from "vue";

import { useDisplay, useTheme } from "vuetify";
import { useRoute, useRouter } from "vue-router";

import { useI18n } from "@/locales/helpers.ts";

import { useRootStore } from "@/stores/index.ts";
import { useSettingsStore } from "@/stores/setting.ts";
import { useUserStore } from "@/stores/user.ts";
import { useDesktopPageStore } from "@/stores/desktopPage.ts";

import { APPLICATION_LOGO_PATH } from "@/consts/asset.ts";
import { ThemeType } from "@/core/theme.ts";

import { getShareCacheImageBlob } from "@/lib/cache.ts";
import { isUserScheduledTransactionEnabled } from "@/lib/server_settings.ts";
import {
    getSystemTheme,
    setExpenseAndIncomeAmountColor,
} from "@/lib/ui/common.ts";
import logger from "@/lib/logger.ts";

import {
    mdiMenu,
    mdiHomeOutline,
    mdiListBoxOutline,
    mdiPlusCircle,
    mdiCreditCardOutline,
    mdiViewDashboardOutline,
    mdiTagOutline,
    mdiClipboardTextOutline,
    mdiClipboardTextClockOutline,
    mdiChartPieOutline,
    mdiCompassOutline,
    mdiSwapHorizontal,
    mdiCogOutline,
    mdiCellphone,
    mdiInformationOutline,
    mdiThemeLightDark,
    mdiWeatherSunny,
    mdiWeatherNight,
    mdiAccount,
    mdiAccountCogOutline,
    mdiLockOutline,
    mdiLogout,
} from "@mdi/js";

type SnackBarType = InstanceType<typeof SnackBar>;

const display = useDisplay();
const theme = useTheme();
const route = useRoute();
const router = useRouter();

const { tt, initLocale } = useI18n();

const rootStore = useRootStore();
const settingsStore = useSettingsStore();
const userStore = useUserStore();
const desktopPageStore = useDesktopPageStore();

const snackbar = useTemplateRef<SnackBarType>("snackbar");

const logouting = ref<boolean>(false);
const isVerticalNavScrolled = ref<boolean>(false);
const showVerticalOverlayMenu = ref<boolean>(false);
const showLoading = ref<boolean>(false);
const showMobileQrCode = ref<boolean>(false);

const mdAndDown = computed<boolean>(() => display.mdAndDown.value);
const currentRoutePath = computed<string>(() => route.path);
const currentPageTitle = computed<string>(() => {
    if (route.path === "/") {
        return tt("Overview");
    } else if (route.path.startsWith("/transaction")) {
        return tt("Transaction Details");
    } else if (route.path.startsWith("/account")) {
        return tt("Accounts");
    } else if (route.path.startsWith("/category")) {
        return tt("Transaction Categories");
    } else if (route.path.startsWith("/tag")) {
        return tt("Transaction Tags");
    } else if (route.path.startsWith("/statistics")) {
        return tt("Statistics & Analysis");
    }

    return tt("global.app.title");
});

const currentNickName = computed<string>(
    () => userStore.currentUserNickname || tt("User"),
);
const currentUserAvatar = computed<string | null>(() =>
    userStore.getUserAvatarUrl(userStore.currentUserBasicInfo, true),
);

const currentTheme = computed<string>({
    get: () => {
        return settingsStore.appSettings.theme;
    },
    set: (value: string) => {
        if (value !== settingsStore.appSettings.theme) {
            settingsStore.setTheme(value);

            if (value === ThemeType.Light || value === ThemeType.Dark) {
                theme.change(value);
            } else {
                theme.change(getSystemTheme());
            }
        }
    },
});

const isEnableApplicationLock = computed<boolean>(
    () => settingsStore.appSettings.applicationLock,
);

function handleNavScroll(e: Event): void {
    isVerticalNavScrolled.value = (e.target as HTMLElement).scrollTop > 0;
}

function clearShareImageCache(): void {
    getShareCacheImageBlob().then((blob) => {
        if (blob) {
            logger.warn(
                "desktop version does not support receving shared image, the share image cache has been cleared",
            );
        }
    });
}

function lock(): void {
    rootStore.lock();
    router.replace("/unlock");
}

function logout(): void {
    logouting.value = true;
    showLoading.value = true;

    rootStore
        .logout()
        .then(() => {
            logouting.value = false;
            showLoading.value = false;

            settingsStore.clearAppSettings();

            const localeDefaultSettings = initLocale(
                userStore.currentUserLanguage,
                settingsStore.appSettings.timeZone,
            );
            settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

            setExpenseAndIncomeAmountColor(
                userStore.currentUserExpenseAmountColor,
                userStore.currentUserIncomeAmountColor,
            );

            router.replace("/login");
        })
        .catch((error) => {
            logouting.value = false;
            showLoading.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
}

function showAddDialogInTransactionListPage(): void {
    const showDialog = (): void => {
        desktopPageStore.setShowAddTransactionDialogInTransactionList();
    };

    if (route.path.startsWith("/transaction/list")) {
        showDialog();
        return;
    }

    router.push("/transaction/list?pageType=0&dateType=7").then(showDialog);
}

clearShareImageCache();
</script>

<style>
.main-logo {
    width: 1.75rem;
    height: 1.75rem;
}

.nav-link.home-link > a:not(.router-link-exact-active):hover::before {
    opacity: calc(var(--v-hover-opacity) * var(--v-theme-overlay-multiplier));
}
</style>
