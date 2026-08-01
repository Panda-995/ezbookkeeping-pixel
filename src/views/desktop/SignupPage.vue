<template>
    <main class="ledger-auth-page ledger-auth-signup">
        <section
            class="ledger-auth-brand-panel"
            aria-labelledby="signup-intro-title"
        >
            <router-link class="ledger-auth-brand" to="/login">
                <span class="ledger-brand-mark">
                    <img alt="" :src="APPLICATION_LOGO_PATH" />
                </span>
                <span>
                    <strong>{{ tt("global.app.title") }}</strong>
                    <small>{{
                        tt("Personal finance, clearly organized")
                    }}</small>
                </span>
            </router-link>

            <div class="ledger-auth-intro">
                <span class="ledger-eyebrow">{{
                    tt("Create an account")
                }}</span>
                <h1 id="signup-intro-title">
                    {{ tt("Start recording without a setup wizard") }}
                </h1>
                <p>{{ tt("Signup streamlined description") }}</p>
            </div>

            <ul class="ledger-auth-benefits">
                <li>
                    <v-icon :icon="mdiCheckCircleOutline" aria-hidden="true" />
                    <span>
                        <strong>{{ tt("Ready after signup") }}</strong>
                        <small>{{
                            tt("Common categories are created automatically")
                        }}</small>
                    </span>
                </li>
                <li>
                    <v-icon :icon="mdiCheckCircleOutline" aria-hidden="true" />
                    <span>
                        <strong>{{ tt("Defaults that make sense") }}</strong>
                        <small>{{
                            tt(
                                "Language, currency and week settings follow your current locale",
                            )
                        }}</small>
                    </span>
                </li>
                <li>
                    <v-icon :icon="mdiCheckCircleOutline" aria-hidden="true" />
                    <span>
                        <strong>{{ tt("Change anything later") }}</strong>
                        <small>{{
                            tt(
                                "Advanced preferences stay available in settings",
                            )
                        }}</small>
                    </span>
                </li>
            </ul>
        </section>

        <section class="ledger-auth-form-panel">
            <div class="ledger-auth-form-card">
                <div class="ledger-auth-form-heading">
                    <span class="ledger-eyebrow">{{ tt("Quick setup") }}</span>
                    <h2>{{ tt("Create an account") }}</h2>
                    <p>
                        {{ tt("Already have an account?") }}
                        <router-link to="/login">{{
                            tt("Log In")
                        }}</router-link>
                    </p>
                </div>

                <v-alert
                    v-if="registrationCompleteMessage"
                    class="mb-6"
                    color="success"
                    role="status"
                    variant="tonal"
                >
                    {{ registrationCompleteMessage }}
                    <div class="mt-4">
                        <router-link class="ledger-inline-link" to="/login">{{
                            tt("Continue")
                        }}</router-link>
                    </div>
                </v-alert>

                <v-form
                    v-else
                    class="ledger-signup-form"
                    @submit.prevent="submit"
                >
                    <div class="ledger-form-grid">
                        <v-text-field
                            id="signup-username"
                            v-model.trim="user.username"
                            type="text"
                            autocomplete="username"
                            autocapitalize="none"
                            autocorrect="off"
                            spellcheck="false"
                            :autofocus="true"
                            :disabled="submitting"
                            :label="tt('Username')"
                            :hint="tt('This will also be your display name')"
                            persistent-hint
                        />

                        <v-text-field
                            id="signup-email"
                            v-model.trim="user.email"
                            type="email"
                            autocomplete="email"
                            :disabled="submitting"
                            :label="tt('E-mail')"
                        />

                        <v-text-field
                            id="signup-password"
                            v-model="user.password"
                            type="password"
                            autocomplete="new-password"
                            :disabled="submitting"
                            :label="tt('Password')"
                            :hint="tt('At least 6 characters')"
                            persistent-hint
                        />

                        <v-text-field
                            id="signup-confirm-password"
                            v-model="user.confirmPassword"
                            type="password"
                            autocomplete="new-password"
                            :disabled="submitting"
                            :label="tt('Confirm Password')"
                            :error-messages="
                                inputInvalidProblemMessage
                                    ? [tt(inputInvalidProblemMessage)]
                                    : []
                            "
                        />
                    </div>

                    <div class="ledger-auto-setup-note">
                        <v-icon :icon="mdiTuneVariant" aria-hidden="true" />
                        <span>
                            <strong>{{ tt("No extra setup required") }}</strong>
                            <small>{{ tt("Signup defaults summary") }}</small>
                        </span>
                    </div>

                    <div
                        class="ledger-form-error"
                        role="alert"
                        aria-live="assertive"
                    >
                        {{
                            submitProblemMessage ? tt(submitProblemMessage) : ""
                        }}
                    </div>

                    <v-btn
                        class="ledger-primary-action"
                        color="primary"
                        type="submit"
                        block
                        :disabled="inputIsEmpty || inputIsInvalid || submitting"
                        :loading="submitting"
                    >
                        {{ tt("Create account and continue") }}
                    </v-btn>
                </v-form>

                <footer class="ledger-auth-card-footer">
                    <language-select-button :disabled="submitting" />
                    <span>{{ tt("Based on ezBookkeeping") }}</span>
                </footer>
            </div>
        </section>

        <snack-bar ref="snackbar" />
    </main>
</template>

<script setup lang="ts">
import SnackBar from "@/components/desktop/SnackBar.vue";

import { ref, computed, useTemplateRef } from "vue";
import { useRouter } from "vue-router";

import { useI18n } from "@/locales/helpers.ts";
import { useSignupPageBase } from "@/views/base/SignupPageBase.ts";
import { useRootStore } from "@/stores/index.ts";

import type { LocalizedPresetCategory } from "@/core/category.ts";
import { APPLICATION_LOGO_PATH } from "@/consts/asset.ts";

import { categorizedArrayToPlainArray } from "@/lib/common.ts";
import { isUserLogined } from "@/lib/userstate.ts";

import { mdiCheckCircleOutline, mdiTuneVariant } from "@mdi/js";

type SnackBarType = InstanceType<typeof SnackBar>;

const router = useRouter();
const { tt, getAllTransactionDefaultCategories } = useI18n();
const {
    user,
    submitting,
    currentLocale,
    inputEmptyProblemMessage,
    inputInvalidProblemMessage,
    inputIsEmpty,
    inputIsInvalid,
    prepareUserForSignup,
    doAfterSignupSuccess,
} = useSignupPageBase();

const rootStore = useRootStore();
const snackbar = useTemplateRef<SnackBarType>("snackbar");
const registrationCompleteMessage = ref<string>("");
const submitted = ref<boolean>(false);

const submitProblemMessage = computed<string>(() => {
    if (!submitted.value) {
        return "";
    }

    return inputEmptyProblemMessage.value || inputInvalidProblemMessage.value;
});

function submit(): void {
    submitted.value = true;
    prepareUserForSignup();

    const problemMessage =
        inputEmptyProblemMessage.value || inputInvalidProblemMessage.value;

    if (problemMessage) {
        snackbar.value?.showMessage(problemMessage);
        return;
    }

    submitting.value = true;

    const allPresetCategories = getAllTransactionDefaultCategories(
        0,
        currentLocale.value,
    );
    const presetCategories: LocalizedPresetCategory[] =
        categorizedArrayToPlainArray(allPresetCategories);

    rootStore
        .register({
            user: user.value,
            presetCategories,
        })
        .then((response) => {
            submitting.value = false;

            if (!isUserLogined()) {
                registrationCompleteMessage.value = response.needVerifyEmail
                    ? tt(
                          "You have been successfully registered. An account activation link has been sent to your email address, please activate your account first.",
                      )
                    : tt("You have been successfully registered");
                return;
            }

            doAfterSignupSuccess(response);

            if (!response.presetCategoriesSaved) {
                snackbar.value?.showMessage(
                    "You have been successfully registered, but there was an failure when adding preset categories. You can re-add preset categories in settings page anytime.",
                );
            } else {
                snackbar.value?.showMessage(
                    "You have been successfully registered",
                );
            }

            router.replace("/");
        })
        .catch((error) => {
            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
}
</script>
