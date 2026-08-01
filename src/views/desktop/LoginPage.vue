<template>
    <main class="atelier-auth">
        <header class="atelier-auth-topline">
            <router-link class="atelier-brand" to="/">
                <span class="atelier-brand-mark">
                    <img
                        alt=""
                        width="26"
                        height="26"
                        :src="APPLICATION_LOGO_PATH"
                    />
                </span>
                <span>
                    <small>OPEN LEDGER / SELF-HOSTED</small>
                    <strong>{{ tt("global.app.title") }}</strong>
                </span>
            </router-link>
            <div class="atelier-auth-meta">
                <language-select-button
                    :disabled="
                        loggingInByPassword || loggingInByOAuth2 || verifying
                    "
                />
                <span>{{ version }}</span>
            </div>
        </header>

        <section class="atelier-auth-stage">
            <div
                class="atelier-manifesto"
                aria-labelledby="login-manifesto-title"
            >
                <span class="atelier-manifesto-index"
                    >PRIVATE FINANCE · 2026</span
                >
                <h1 id="login-manifesto-title">
                    OWN YOUR<br />
                    <em>MONEY.</em>
                </h1>
                <p>
                    {{
                        tt(
                            "Every transaction stays clear, editable and traceable",
                        )
                    }}
                </p>

                <div class="atelier-figure" aria-hidden="true">
                    <span>NET / 07</span>
                    <strong>¥ 12,680</strong>
                    <div>
                        <i style="--bar-size: 38%"></i>
                        <i style="--bar-size: 62%"></i>
                        <i style="--bar-size: 48%"></i>
                        <i style="--bar-size: 79%"></i>
                        <i style="--bar-size: 56%"></i>
                        <i style="--bar-size: 92%"></i>
                    </div>
                </div>

                <div class="atelier-source">
                    <span>{{ tt("Based on ezBookkeeping") }}</span>
                    <a
                        href="https://github.com/mayswind/ezbookkeeping"
                        target="_blank"
                        rel="noreferrer"
                    >
                        mayswind / ezbookkeeping ↗
                    </a>
                </div>
            </div>

            <section class="atelier-login-card" aria-labelledby="login-title">
                <div class="atelier-card-tape" aria-hidden="true"></div>
                <div class="atelier-card-heading">
                    <span>{{
                        show2faInput ? "SECURITY / 02" : "SECURE ACCESS / 01"
                    }}</span>
                    <h2 id="login-title">
                        {{ show2faInput ? tt("Continue") : tt("Log In") }}
                    </h2>
                    <p v-if="isInternalAuthEnabled()">
                        {{
                            tips ||
                            tt("Please log in with your ezBookkeeping account")
                        }}
                    </p>
                </div>

                <v-form
                    class="atelier-login-form"
                    @submit.prevent="show2faInput ? verify() : login()"
                >
                    <template v-if="isInternalAuthEnabled() && !show2faInput">
                        <label
                            class="atelier-field-label"
                            for="atelier-login-username"
                        >
                            <span>01</span>{{ tt("Username") }}
                        </label>
                        <v-text-field
                            id="atelier-login-username"
                            name="username"
                            v-model.trim="username"
                            type="text"
                            autocomplete="username"
                            autocapitalize="none"
                            autocorrect="off"
                            spellcheck="false"
                            inputmode="email"
                            variant="outlined"
                            :disabled="
                                loggingInByPassword ||
                                loggingInByOAuth2 ||
                                verifying
                            "
                            hide-details="auto"
                            @input="tempToken = ''"
                            @keyup.enter="passwordInput?.focus()"
                        />

                        <label
                            class="atelier-field-label"
                            for="atelier-login-password"
                        >
                            <span>02</span>{{ tt("Password") }}
                        </label>
                        <v-text-field
                            id="atelier-login-password"
                            ref="passwordInput"
                            name="password"
                            v-model="password"
                            type="password"
                            autocomplete="current-password"
                            variant="outlined"
                            :disabled="
                                loggingInByPassword ||
                                loggingInByOAuth2 ||
                                verifying
                            "
                            hide-details="auto"
                            @input="tempToken = ''"
                            @keyup.enter="login"
                        />

                        <div class="atelier-form-links">
                            <button
                                type="button"
                                @click="showMobileQrCode = true"
                            >
                                {{ tt("Use on Mobile Device") }}
                            </button>
                            <router-link
                                to="/forgetpassword"
                                :aria-disabled="
                                    !isUserForgetPasswordEnabled() ||
                                    loggingInByPassword ||
                                    loggingInByOAuth2 ||
                                    verifying
                                "
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
                        <label
                            class="atelier-field-label"
                            for="atelier-login-2fa"
                        >
                            <span>03</span>
                            {{
                                twoFAVerifyType === "passcode"
                                    ? tt("Passcode")
                                    : tt("Backup Code")
                            }}
                        </label>
                        <v-text-field
                            id="atelier-login-2fa"
                            ref="passcodeInput"
                            name="passcode"
                            v-model="passcode"
                            type="number"
                            inputmode="numeric"
                            autocomplete="one-time-code"
                            variant="outlined"
                            :disabled="
                                loggingInByPassword ||
                                loggingInByOAuth2 ||
                                verifying
                            "
                            :append-inner-icon="mdiHelpCircleOutline"
                            :aria-label="tt('Use Backup Code')"
                            hide-details="auto"
                            @click:append-inner="twoFAVerifyType = 'backupcode'"
                            @keyup.enter="verify"
                            v-if="twoFAVerifyType === 'passcode'"
                        />
                        <v-text-field
                            id="atelier-login-2fa"
                            name="backup-code"
                            v-model="backupCode"
                            type="text"
                            autocomplete="off"
                            variant="outlined"
                            :disabled="
                                loggingInByPassword ||
                                loggingInByOAuth2 ||
                                verifying
                            "
                            :append-inner-icon="mdiOnepassword"
                            :aria-label="tt('Use Passcode')"
                            hide-details="auto"
                            @click:append-inner="twoFAVerifyType = 'passcode'"
                            @keyup.enter="verify"
                            v-else
                        />
                    </template>

                    <v-btn
                        v-if="isInternalAuthEnabled()"
                        class="atelier-submit"
                        color="primary"
                        type="submit"
                        block
                        :disabled="
                            show2faInput ? twoFAInputIsEmpty : inputIsEmpty
                        "
                        :loading="loggingInByPassword || verifying"
                    >
                        {{ show2faInput ? tt("Continue") : tt("Log In") }}
                        <span aria-hidden="true">→</span>
                    </v-btn>

                    <template v-if="isOAuth2Enabled()">
                        <div class="atelier-separator">
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

                <div class="atelier-create" v-if="isInternalAuthEnabled()">
                    <span>{{ tt("Don't have an account?") }}</span>
                    <router-link
                        to="/signup"
                        :aria-disabled="
                            !isUserRegistrationEnabled() ||
                            loggingInByPassword ||
                            loggingInByOAuth2 ||
                            verifying
                        "
                        :class="{
                            disabled:
                                !isUserRegistrationEnabled() ||
                                loggingInByPassword ||
                                loggingInByOAuth2 ||
                                verifying,
                        }"
                    >
                        {{ tt("Create an account") }} ↗
                    </router-link>
                </div>
            </section>
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
