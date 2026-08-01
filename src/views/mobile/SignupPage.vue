<template>
    <f7-page class="mobile-atelier-signup" no-navbar no-swipeback>
        <header class="mobile-atelier-signup-topline">
            <f7-link href="/login" :aria-label="tt('Back')"
                >← {{ tt("Back") }}</f7-link
            >
            <span class="mobile-atelier-signup-brand">
                <img
                    alt=""
                    width="24"
                    height="24"
                    :src="APPLICATION_LOGO_PATH"
                />
                <strong>{{ tt("global.app.title") }}</strong>
            </span>
            <small>01 / 01</small>
        </header>

        <main class="mobile-atelier-signup-content">
            <header class="mobile-atelier-signup-hero">
                <span>QUICK START · ONE PAGE</span>
                <h1>{{ tt("Start recording without a setup wizard") }}</h1>
                <p>{{ tt("Signup streamlined description") }}</p>
                <i aria-hidden="true">＋</i>
            </header>

            <form class="mobile-atelier-signup-form" @submit.prevent="submit">
                <div class="mobile-atelier-form-head">
                    <span>ACCOUNT / DETAILS</span>
                    <h2>{{ tt("Create an account") }}</h2>
                </div>

                <label class="mobile-atelier-sr-only" for="signup-username">{{ tt("Username") }}</label>
                <label class="mobile-atelier-sr-only" for="signup-email">{{ tt("E-mail") }}</label>
                <label class="mobile-atelier-sr-only" for="signup-password">{{ tt("Password") }}</label>
                <label class="mobile-atelier-sr-only" for="signup-confirm-password">{{ tt("Confirm Password") }}</label>

                <f7-list form strong inset class="mobile-atelier-form-list">
                    <f7-list-input
                        input-id="signup-username"
                        name="username"
                        type="text"
                        autocomplete="username"
                        autocapitalize="none"
                        autocorrect="off"
                        spellcheck="false"
                        required
                        clear-button
                        outline
                        :disabled="submitting"
                        :label="tt('Username')"
                        v-model:value="user.username"
                    />

                    <f7-list-input
                        input-id="signup-email"
                        name="email"
                        type="email"
                        autocomplete="email"
                        spellcheck="false"
                        required
                        clear-button
                        outline
                        :disabled="submitting"
                        :label="tt('E-mail')"
                        v-model:value="user.email"
                    />

                    <f7-list-input
                        input-id="signup-password"
                        name="new-password"
                        type="password"
                        autocomplete="new-password"
                        required
                        clear-button
                        outline
                        :disabled="submitting"
                        :label="tt('Password')"
                        v-model:value="user.password"
                    />

                    <f7-list-input
                        input-id="signup-confirm-password"
                        name="confirm-password"
                        type="password"
                        autocomplete="new-password"
                        required
                        clear-button
                        outline
                        :disabled="submitting"
                        :label="tt('Confirm Password')"
                        v-model:value="user.confirmPassword"
                    />
                </f7-list>

                <p
                    class="mobile-atelier-form-error"
                    role="alert"
                    aria-live="assertive"
                >
                    {{
                        submitted ? formProblemMessage : ""
                    }}
                </p>

                <div class="mobile-atelier-auto-note">
                    <f7-icon
                        f7="slider_horizontal_3"
                        aria-hidden="true"
                    ></f7-icon>
                    <span>
                        <strong>{{ tt("No extra setup required") }}</strong>
                        <small>{{ tt("Signup defaults summary") }}</small>
                    </span>
                </div>

                <f7-button
                    large
                    fill
                    type="submit"
                    class="mobile-atelier-submit"
                    :class="{ disabled: submitting }"
                    :disabled="submitting"
                    :aria-busy="submitting"
                    :text="
                        submitting
                            ? tt('Loading...')
                            : `${tt('Create account and continue')} →`
                    "
                />
            </form>

            <p class="mobile-atelier-signup-alternate">
                {{ tt("Already have an account?") }}
                <f7-link href="/login">{{ tt("Log In") }} ↗</f7-link>
            </p>
        </main>
    </f7-page>
</template>

<script setup lang="ts">
import { computed, ref, watch } from "vue";
import type { Router } from "framework7/types";

import { useI18n } from "@/locales/helpers.ts";
import {
    useI18nUIComponents,
    showLoading,
    hideLoading,
} from "@/lib/ui/mobile.ts";
import { useSignupPageBase } from "@/views/base/SignupPageBase.ts";
import { useRootStore } from "@/stores/index.ts";

import type { LocalizedPresetCategory } from "@/core/category.ts";
import { APPLICATION_LOGO_PATH } from "@/consts/asset.ts";

import { categorizedArrayToPlainArray } from "@/lib/common.ts";
import { isUserLogined } from "@/lib/userstate.ts";

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt, te, getAllTransactionDefaultCategories } = useI18n();
const { showAlert, showToast } = useI18nUIComponents();
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
const submitted = ref<boolean>(false);
const submissionProblemMessage = ref<string>("");
const formProblemMessage = computed<string>(() => {
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
        showAlert(problemMessage);
        focusFirstInvalidInput();
        return;
    }

    submitting.value = true;
    showLoading(() => submitting.value);

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
            hideLoading();

            if (!isUserLogined()) {
                showToast(
                    response.needVerifyEmail
                        ? "You have been successfully registered. An account activation link has been sent to your email address, please activate your account first."
                        : "You have been successfully registered",
                    5000,
                );
                props.f7router.navigate("/login");
                return;
            }

            doAfterSignupSuccess(response);

            if (!response.presetCategoriesSaved) {
                showToast(
                    "You have been successfully registered, but there was an failure when adding preset categories. You can re-add preset categories in settings page anytime.",
                    5000,
                );
            } else {
                showToast("You have been successfully registered");
            }

            props.f7router.navigate("/");
        })
        .catch((error) => {
            submitting.value = false;
            hideLoading();
            submissionProblemMessage.value =
                te(error) || tt("Unable to sign up");

            if (!error.processed) {
                showToast(error.message || error);
            }
        });
}
</script>
