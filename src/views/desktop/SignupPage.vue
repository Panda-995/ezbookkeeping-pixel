<template>
    <main class="atelier-signup">
        <header class="atelier-auth-topline">
            <router-link class="atelier-brand" to="/login">
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
                <language-select-button :disabled="submitting" />
                <router-link to="/login">{{ tt("Log In") }} ↗</router-link>
            </div>
        </header>

        <section class="atelier-signup-stage">
            <section
                class="atelier-signup-form-wrap"
                aria-labelledby="signup-title"
            >
                <span class="atelier-manifesto-index"
                    >QUICK START · ONE PAGE</span
                >
                <h1 id="signup-title">
                    {{ tt("Start recording without a setup wizard") }}
                </h1>
                <p>{{ tt("Signup streamlined description") }}</p>

                <v-alert
                    v-if="registrationCompleteMessage"
                    class="atelier-success"
                    color="success"
                    role="status"
                    variant="tonal"
                >
                    {{ registrationCompleteMessage }}
                    <div class="mt-4">
                        <router-link class="atelier-inline-link" to="/login">
                            {{ tt("Continue") }} →
                        </router-link>
                    </div>
                </v-alert>

                <v-form
                    v-else
                    class="atelier-signup-form"
                    @submit.prevent="submit"
                >
                    <div class="atelier-signup-grid">
                        <v-text-field
                            id="signup-username"
                            name="username"
                            v-model.trim="user.username"
                            type="text"
                            autocomplete="username"
                            autocapitalize="none"
                            autocorrect="off"
                            spellcheck="false"
                            required
                            variant="outlined"
                            :disabled="submitting"
                            :error="submitted && !user.username"
                            :label="tt('Username')"
                            :hint="tt('This will also be your display name')"
                            persistent-hint
                        />

                        <v-text-field
                            id="signup-email"
                            name="email"
                            v-model.trim="user.email"
                            type="email"
                            autocomplete="email"
                            spellcheck="false"
                            required
                            variant="outlined"
                            :disabled="submitting"
                            :error="submitted && !user.email"
                            :label="tt('E-mail')"
                        />

                        <v-text-field
                            id="signup-password"
                            name="new-password"
                            v-model="user.password"
                            type="password"
                            autocomplete="new-password"
                            required
                            variant="outlined"
                            :disabled="submitting"
                            :error="submitted && !user.password"
                            :label="tt('Password')"
                            :hint="tt('At least 6 characters')"
                            persistent-hint
                        />

                        <v-text-field
                            id="signup-confirm-password"
                            name="confirm-password"
                            v-model="user.confirmPassword"
                            type="password"
                            autocomplete="new-password"
                            required
                            variant="outlined"
                            :disabled="submitting"
                            :label="tt('Confirm Password')"
                            :error-messages="
                                inputInvalidProblemMessage
                                    ? [tt(inputInvalidProblemMessage)]
                                    : []
                            "
                        />
                    </div>

                    <div class="atelier-auto-note">
                        <span class="atelier-auto-icon">
                            <v-icon :icon="mdiTuneVariant" aria-hidden="true" />
                        </span>
                        <span>
                            <strong>{{ tt("No extra setup required") }}</strong>
                            <small>{{ tt("Signup defaults summary") }}</small>
                        </span>
                    </div>

                    <div
                        class="atelier-form-error"
                        role="alert"
                        aria-live="assertive"
                    >
                        {{
                            submitProblemMessage
                        }}
                    </div>

                    <v-btn
                        class="atelier-submit"
                        color="primary"
                        type="submit"
                        block
                        :disabled="submitting"
                        :loading="submitting"
                        :aria-busy="submitting"
                    >
                        {{ tt("Create account and continue") }}
                        <span aria-hidden="true">→</span>
                    </v-btn>
                </v-form>
            </section>

            <aside
                class="atelier-signup-poster"
                aria-labelledby="signup-benefits-title"
            >
                <span class="atelier-poster-index">SETUP / 01—03</span>
                <h2 id="signup-benefits-title">
                    READY<br />FROM THE<br />FIRST ENTRY.
                </h2>
                <ul>
                    <li>
                        <span>01</span>
                        <div>
                            <strong>{{ tt("Ready after signup") }}</strong>
                            <small>{{
                                tt(
                                    "Common categories are created automatically",
                                )
                            }}</small>
                        </div>
                        <v-icon
                            :icon="mdiCheckCircleOutline"
                            aria-hidden="true"
                        />
                    </li>
                    <li>
                        <span>02</span>
                        <div>
                            <strong>{{
                                tt("Defaults that make sense")
                            }}</strong>
                            <small>{{
                                tt(
                                    "Language, currency and week settings follow your current locale",
                                )
                            }}</small>
                        </div>
                        <v-icon
                            :icon="mdiCheckCircleOutline"
                            aria-hidden="true"
                        />
                    </li>
                    <li>
                        <span>03</span>
                        <div>
                            <strong>{{ tt("Change anything later") }}</strong>
                            <small>{{
                                tt(
                                    "Advanced preferences stay available in settings",
                                )
                            }}</small>
                        </div>
                        <v-icon
                            :icon="mdiCheckCircleOutline"
                            aria-hidden="true"
                        />
                    </li>
                </ul>
                <footer>{{ tt("Based on ezBookkeeping") }}</footer>
            </aside>
        </section>

        <snack-bar ref="snackbar" />
    </main>
</template>

<script setup lang="ts">
import SnackBar from "@/components/desktop/SnackBar.vue";

import { ref, computed, useTemplateRef, watch } from "vue";
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
const { tt, te, getAllTransactionDefaultCategories } = useI18n();
const {
    user,
    submitting,
    currentLocale,
    inputEmptyProblemMessage,
    inputInvalidProblemMessage,
    prepareUserForSignup,
    focusFirstInvalidInput,
    doAfterSignupSuccess,
} = useSignupPageBase();

const rootStore = useRootStore();
const snackbar = useTemplateRef<SnackBarType>("snackbar");
const registrationCompleteMessage = ref<string>("");
const submitted = ref<boolean>(false);
const submissionProblemMessage = ref<string>("");

const submitProblemMessage = computed<string>(() => {
    if (!submitted.value) {
        return "";
    }

    const validationProblemMessage =
        inputEmptyProblemMessage.value || inputInvalidProblemMessage.value;

    if (validationProblemMessage) {
        return tt(validationProblemMessage);
    }

    return submissionProblemMessage.value;
});

watch(
    () => [
        user.value.username,
        user.value.email,
        user.value.password,
        user.value.confirmPassword,
    ],
    () => {
        submissionProblemMessage.value = "";
    },
);

function submit(): void {
    submitted.value = true;
    submissionProblemMessage.value = "";
    prepareUserForSignup();

    const problemMessage =
        inputEmptyProblemMessage.value || inputInvalidProblemMessage.value;

    if (problemMessage) {
        snackbar.value?.showMessage(problemMessage);
        focusFirstInvalidInput();
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
            submissionProblemMessage.value =
                te(error) || tt("Unable to sign up");

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
}
</script>
