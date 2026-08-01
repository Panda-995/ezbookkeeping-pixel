<template>
    <main class="pixel-auth-terminal">
        <section class="pixel-auth-intro">
            <router-link class="pixel-auth-brand" to="/">
                <span class="pixel-brand-mark">
                    <img alt="" :src="APPLICATION_LOGO_PATH" />
                </span>
                <span>
                    <strong>{{ tt("global.app.title") }}</strong>
                    <small>{{
                        tt("Personal finance, clearly organized")
                    }}</small>
                </span>
            </router-link>

            <div class="pixel-auth-copy">
                <div class="pixel-kicker">
                    <span class="pixel-status-dot" aria-hidden="true"></span>
                    {{ tt("Private, self-hosted and open source") }}
                </div>
                <h1>{{ tt("Welcome to ezBookkeeping") }}</h1>
                <p>
                    {{
                        tt(
                            "Every transaction stays clear, editable and traceable",
                        )
                    }}
                </p>
            </div>

            <div class="pixel-auth-ledger" aria-hidden="true">
                <div class="pixel-auth-ledger-head">
                    <span>MONTHLY REGISTER</span>
                    <span>2026 / 07</span>
                </div>
                <div class="pixel-auth-ledger-balance">
                    <small>NET ASSETS</small>
                    <strong>¥ 12,680.00</strong>
                </div>
                <div class="pixel-auth-ledger-line">
                    <span><i class="is-income"></i> INCOME</span>
                    <strong>+ 8,320.00</strong>
                </div>
                <div class="pixel-auth-ledger-line">
                    <span><i class="is-expense"></i> EXPENSE</span>
                    <strong>− 3,175.40</strong>
                </div>
                <div class="pixel-auth-ledger-cells">
                    <i v-for="index in 20" :key="index"></i>
                </div>
            </div>

            <div class="pixel-auth-intro-footer">
                <span>{{ tt("Based on ezBookkeeping") }}</span>
                <a
                    href="https://github.com/mayswind/ezbookkeeping"
                    target="_blank"
                    rel="noreferrer"
                    >mayswind/ezbookkeeping</a
                >
            </div>
        </section>

        <section class="pixel-auth-gate">
            <div class="pixel-auth-panel">
                <div class="pixel-auth-panel-head">
                    <span>SECURE ACCESS / 01</span>
                    <i aria-hidden="true"></i>
                </div>

                <div class="pixel-auth-panel-title">
                    <span class="pixel-panel-index">{{
                        show2faInput ? "TWO-FACTOR CHECK" : "ACCOUNT LOGIN"
                    }}</span>
                    <h2>{{ show2faInput ? tt("Continue") : tt("Log In") }}</h2>
                    <p v-if="isInternalAuthEnabled()">
                        {{
                            tips ||
                            tt("Please log in with your ezBookkeeping account")
                        }}
                    </p>
                </div>

                <v-form
                    class="pixel-auth-form"
                    @submit.prevent="show2faInput ? verify() : login()"
                >
                    <template v-if="isInternalAuthEnabled() && !show2faInput">
                        <label
                            class="pixel-field-label"
                            for="pixel-login-username"
                            >{{ tt("Username") }}</label
                        >
                        <v-text-field
                            id="pixel-login-username"
                            v-model.trim="username"
                            type="text"
                            autocomplete="username"
                            autocapitalize="none"
                            autocorrect="off"
                            spellcheck="false"
                            inputmode="email"
                            :autofocus="true"
                            :disabled="
                                loggingInByPassword ||
                                loggingInByOAuth2 ||
                                verifying
                            "
                            :placeholder="tt('Your username or email')"
                            hide-details="auto"
                            @input="tempToken = ''"
                            @keyup.enter="passwordInput?.focus()"
                        />

                        <label
                            class="pixel-field-label"
                            for="pixel-login-password"
                            >{{ tt("Password") }}</label
                        >
                        <v-text-field
                            id="pixel-login-password"
                            ref="passwordInput"
                            v-model="password"
                            type="password"
                            autocomplete="current-password"
                            :disabled="
                                loggingInByPassword ||
                                loggingInByOAuth2 ||
                                verifying
                            "
                            :placeholder="tt('Your password')"
                            hide-details="auto"
                            @input="tempToken = ''"
                            @keyup.enter="login"
                        />

                        <div class="pixel-auth-links">
                            <button
                                type="button"
                                @click="showMobileQrCode = true"
                            >
                                {{ tt("Use on Mobile Device") }}
                            </button>
                            <router-link
                                to="/forgetpassword"
                                :class="{
                                    disabled:
                                        !isUserForgetPasswordEnabled() ||
                                        loggingInByPassword ||
                                        loggingInByOAuth2 ||
                                        verifying,
                                }"
                            >
                                {{ tt("Forget Password?") }}
                            </router-link>
                        </div>
                    </template>

                    <template
                        v-else-if="isInternalAuthEnabled() && show2faInput"
                    >
                        <label class="pixel-field-label" for="pixel-login-2fa">
                            {{
                                twoFAVerifyType === "passcode"
                                    ? tt("Passcode")
                                    : tt("Backup Code")
                            }}
                        </label>
                        <v-text-field
                            id="pixel-login-2fa"
                            ref="passcodeInput"
                            v-model="passcode"
                            type="number"
                            autocomplete="one-time-code"
                            :disabled="
                                loggingInByPassword ||
                                loggingInByOAuth2 ||
                                verifying
                            "
                            :placeholder="tt('Passcode')"
                            :append-inner-icon="mdiHelpCircleOutline"
                            hide-details="auto"
                            @click:append-inner="twoFAVerifyType = 'backupcode'"
                            @keyup.enter="verify"
                            v-if="twoFAVerifyType === 'passcode'"
                        />
                        <v-text-field
                            id="pixel-login-2fa"
                            v-model="backupCode"
                            type="text"
                            :disabled="
                                loggingInByPassword ||
                                loggingInByOAuth2 ||
                                verifying
                            "
                            :placeholder="tt('Backup Code')"
                            :append-inner-icon="mdiOnepassword"
                            hide-details="auto"
                            @click:append-inner="twoFAVerifyType = 'passcode'"
                            @keyup.enter="verify"
                            v-else
                        />
                    </template>

                    <v-btn
                        class="pixel-auth-submit"
                        color="primary"
                        type="submit"
                        block
                        :disabled="
                            show2faInput ? twoFAInputIsEmpty : inputIsEmpty
                        "
                        :loading="loggingInByPassword || verifying"
                        v-if="isInternalAuthEnabled()"
                    >
                        {{ show2faInput ? tt("Continue") : tt("Log In") }}
                    </v-btn>

                    <template v-if="isOAuth2Enabled()">
                        <div class="pixel-auth-separator">
                            <span>{{ tt("or") }}</span>
                        </div>
                        <v-btn
                            block
                            variant="outlined"
                            :disabled="
                                show2faInput ||
                                loggingInByPassword ||
                                loggingInByOAuth2 ||
                                verifying
                            "
                            :href="oauth2LoginUrl"
                            :loading="loggingInByOAuth2"
                            @click="loggingInByOAuth2 = true"
                        >
                            {{ oauth2LoginDisplayName }}
                        </v-btn>
                    </template>
                </v-form>

                <div class="pixel-auth-create" v-if="isInternalAuthEnabled()">
                    <span>{{ tt("Don't have an account?") }}</span>
                    <router-link
                        to="/signup"
                        :class="{
                            disabled:
                                !isUserRegistrationEnabled() ||
                                loggingInByPassword ||
                                loggingInByOAuth2 ||
                                verifying,
                        }"
                    >
                        {{ tt("Create an account") }} →
                    </router-link>
                </div>

                <div class="pixel-auth-meta">
                    <language-select-button
                        :disabled="
                            loggingInByPassword ||
                            loggingInByOAuth2 ||
                            verifying
                        "
                    />
                    <span>{{ version }}</span>
                </div>
            </div>
        </section>

        <switch-to-mobile-dialog v-model:show="showMobileQrCode" />
        <snack-bar ref="snackbar" />
    </main>
</template>

<script setup lang="ts">
import { VTextField } from "vuetify/components/VTextField";
import SnackBar from "@/components/desktop/SnackBar.vue";

import { ref, useTemplateRef, nextTick } from "vue";
import { useRouter } from "vue-router";

import { useI18n } from "@/locales/helpers.ts";
import { useLoginPageBase } from "@/views/base/LoginPageBase.ts";

import { useRootStore } from "@/stores/index.ts";

import { APPLICATION_LOGO_PATH } from "@/consts/asset.ts";
import { KnownErrorCode } from "@/consts/api.ts";

import { generateRandomUUID } from "@/lib/misc.ts";
import {
    isUserRegistrationEnabled,
    isUserForgetPasswordEnabled,
    isUserVerifyEmailEnabled,
    isInternalAuthEnabled,
    isOAuth2Enabled,
} from "@/lib/server_settings.ts";

import { mdiOnepassword, mdiHelpCircleOutline } from "@mdi/js";

type SnackBarType = InstanceType<typeof SnackBar>;

const router = useRouter();

const { tt } = useI18n();

const rootStore = useRootStore();

const {
    version,
    username,
    password,
    passcode,
    backupCode,
    tempToken,
    twoFAVerifyType,
    oauth2ClientSessionId,
    loggingInByPassword,
    loggingInByOAuth2,
    verifying,
    inputIsEmpty,
    twoFAInputIsEmpty,
    oauth2LoginUrl,
    oauth2LoginDisplayName,
    tips,
    doAfterLogin,
} = useLoginPageBase("desktop");

const passwordInput = useTemplateRef<VTextField>("passwordInput");
const passcodeInput = useTemplateRef<VTextField>("passcodeInput");
const snackbar = useTemplateRef<SnackBarType>("snackbar");

const show2faInput = ref<boolean>(false);
const showMobileQrCode = ref<boolean>(false);

function login(): void {
    if (!username.value) {
        snackbar.value?.showMessage("Username cannot be blank");
        return;
    }

    if (!password.value) {
        snackbar.value?.showMessage("Password cannot be blank");
        return;
    }

    if (tempToken.value) {
        show2faInput.value = true;
        return;
    }

    if (loggingInByPassword.value) {
        return;
    }

    loggingInByPassword.value = true;

    rootStore
        .authorize({
            loginName: username.value,
            password: password.value,
        })
        .then((authResponse) => {
            loggingInByPassword.value = false;

            if (authResponse.need2FA) {
                tempToken.value = authResponse.token;
                show2faInput.value = true;

                nextTick(() => {
                    if (passcodeInput.value) {
                        passcodeInput.value.focus();
                        passcodeInput.value.select();
                    }
                });

                return;
            }

            doAfterLogin(authResponse);
            router.replace("/");
        })
        .catch((error) => {
            loggingInByPassword.value = false;

            if (
                isUserVerifyEmailEnabled() &&
                error.error &&
                error.error.errorCode === KnownErrorCode.UserEmailNotVerified &&
                error.error.context &&
                error.error.context.email
            ) {
                router.push(
                    `/verify_email?email=${encodeURIComponent(error.error.context.email)}&emailSent=${error.error.context.hasValidEmailVerifyToken || false}`,
                );
                return;
            }

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
}

function verify(): void {
    if (twoFAInputIsEmpty.value || verifying.value) {
        return;
    }

    if (twoFAVerifyType.value === "passcode" && !passcode.value) {
        snackbar.value?.showMessage("Passcode cannot be blank");
        return;
    } else if (twoFAVerifyType.value === "backupcode" && !backupCode.value) {
        snackbar.value?.showMessage("Backup code cannot be blank");
        return;
    }

    verifying.value = true;

    rootStore
        .authorize2FA({
            token: tempToken.value,
            passcode:
                twoFAVerifyType.value === "passcode" ? passcode.value : null,
            recoveryCode:
                twoFAVerifyType.value === "backupcode"
                    ? backupCode.value
                    : null,
        })
        .then((authResponse) => {
            verifying.value = false;

            doAfterLogin(authResponse);
            router.replace("/");
        })
        .catch((error) => {
            verifying.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
}

oauth2ClientSessionId.value = generateRandomUUID();
</script>
