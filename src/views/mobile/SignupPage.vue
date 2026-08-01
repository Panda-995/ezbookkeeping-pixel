<template>
    <f7-page class="ledger-mobile-signup" no-swipeback>
        <f7-navbar class="ledger-mobile-navbar">
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Create an account')"></f7-nav-title>
        </f7-navbar>

        <main class="ledger-mobile-auth-content">
            <header class="ledger-mobile-auth-heading">
                <span class="ledger-brand-mark">
                    <img alt="" :src="APPLICATION_LOGO_PATH" />
                </span>
                <div>
                    <span class="ledger-eyebrow">{{ tt("Quick setup") }}</span>
                    <h1>{{ tt("Start recording without a setup wizard") }}</h1>
                    <p>{{ tt("Signup streamlined description") }}</p>
                </div>
            </header>

            <form class="ledger-mobile-form" @submit.prevent="submit">
                <f7-list
                    form
                    strong
                    inset
                    dividers
                    class="ledger-mobile-form-list"
                >
                    <f7-list-input
                        type="text"
                        autocomplete="username"
                        autocapitalize="none"
                        autocorrect="off"
                        spellcheck="false"
                        clear-button
                        :disabled="submitting"
                        :label="tt('Username')"
                        :placeholder="tt('Your username')"
                        v-model:value="user.username"
                    />

                    <f7-list-input
                        type="email"
                        autocomplete="email"
                        clear-button
                        :disabled="submitting"
                        :label="tt('E-mail')"
                        :placeholder="tt('Your email address')"
                        v-model:value="user.email"
                    />

                    <f7-list-input
                        type="password"
                        autocomplete="new-password"
                        clear-button
                        :disabled="submitting"
                        :label="tt('Password')"
                        :placeholder="
                            tt('Your password, at least 6 characters')
                        "
                        v-model:value="user.password"
                    />

                    <f7-list-input
                        type="password"
                        autocomplete="new-password"
                        clear-button
                        :disabled="submitting"
                        :label="tt('Confirm Password')"
                        :placeholder="tt('Re-enter the password')"
                        v-model:value="user.confirmPassword"
                    />

                    <f7-list-item
                        class="ebk-list-item-error-info"
                        v-if="
                            submitted &&
                            (inputEmptyProblemMessage ||
                                inputInvalidProblemMessage)
                        "
                        :footer="
                            tt(
                                inputEmptyProblemMessage ||
                                    inputInvalidProblemMessage,
                            )
                        "
                    />
                </f7-list>

                <div class="ledger-mobile-auto-note">
                    <f7-icon f7="slider_horizontal_3"></f7-icon>
                    <span>
                        <strong>{{ tt("No extra setup required") }}</strong>
                        <small>{{ tt("Signup defaults summary") }}</small>
                    </span>
                </div>

                <f7-button
                    large
                    fill
                    type="submit"
                    class="ledger-mobile-submit"
                    :class="{
                        disabled: inputIsEmpty || inputIsInvalid || submitting,
                    }"
                    :text="
                        submitting
                            ? tt('Loading...')
                            : tt('Create account and continue')
                    "
                />
            </form>

            <p class="ledger-mobile-auth-alternate">
                {{ tt("Already have an account?") }}
                <f7-link href="/login">{{ tt("Log In") }}</f7-link>
            </p>
        </main>
    </f7-page>
</template>

<script setup lang="ts">
import { ref } from "vue";
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

const { tt, getAllTransactionDefaultCategories } = useI18n();
const { showAlert, showToast } = useI18nUIComponents();
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
const submitted = ref<boolean>(false);

function submit(): void {
    submitted.value = true;
    prepareUserForSignup();

    const problemMessage =
        inputEmptyProblemMessage.value || inputInvalidProblemMessage.value;

    if (problemMessage) {
        showAlert(problemMessage);
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

            if (!error.processed) {
                showToast(error.message || error);
            }
        });
}
</script>
