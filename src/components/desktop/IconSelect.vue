<template>
    <v-select
        density="comfortable"
        item-title="icon"
        item-value="id"
        persistent-placeholder
        :disabled="disabled"
        :label="label"
        v-model="icon"
        @update:menu="onMenuStateChanged"
    >
        <template #selection>
            <v-label class="cursor-pointer">
                <ItemIcon :icon-type="iconType" :icon-id="icon" :color="color" />
            </v-label>
        </template>

        <template #no-data>
            <div class="icon-select-dropdown px-2" ref="dropdownMenu">
                <div class="icon-item" :class="{ 'row-has-selected-item': hasSelectedIcon(row) }"
                     :style="`grid-template-columns: repeat(${itemPerRow}, minmax(0, 1fr));`"
                     :key="idx" v-for="(row, idx) in allIconRows">
                    <div class="text-center" :key="iconInfo.id" v-for="iconInfo in row">
                        <button
                            type="button"
                            class="ebk-option-button"
                            :aria-label="`${label || tt('Icon')} ${iconInfo.id}`"
                            :aria-pressed="modelValue === iconInfo.id"
                            @click="icon = iconInfo.id"
                        >
                            <ItemIcon class="ma-2" icon-type="fixed" :icon-id="iconInfo.icon" :color="color" v-if="!modelValue || modelValue !== iconInfo.id" />
                            <v-badge class="right-bottom-icon" color="primary"
                                     offset-x="8" offset-y="10"
                                     :location="`bottom ${textDirection === TextDirection.LTR ? 'right' : 'left'}`"
                                     :icon="mdiCheck"
                                     v-if="modelValue && modelValue === iconInfo.id">
                                <ItemIcon class="ma-2" icon-type="fixed" :icon-id="iconInfo.icon" :color="color" />
                            </v-badge>
                        </button>
                    </div>
                </div>
            </div>
        </template>
    </v-select>
</template>

<script setup lang="ts">
import { ref, computed, useTemplateRef, nextTick } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { TextDirection } from '@/core/text.ts';
import type { ColorValue } from '@/core/color.ts';
import type { IconInfo, IconInfoWithId } from '@/core/icon.ts';
import { arrayContainsFieldValue } from '@/lib/common.ts';
import { getIconsInRows } from '@/lib/icon.ts';
import { scrollToSelectedItem } from '@/lib/ui/common.ts';

import {
    mdiCheck
} from '@mdi/js';

const props = defineProps<{
    modelValue: string;
    disabled?: boolean;
    label?: string;
    iconType: string;
    color: ColorValue;
    columnCount?: number;
    allIconInfos: Record<string, IconInfo>;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
}>();

const { tt, getCurrentLanguageTextDirection } = useI18n();

const dropdownMenu = useTemplateRef<HTMLElement>('dropdownMenu');
const itemPerRow = ref<number>(props.columnCount || 7);

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());
const allIconRows = computed<IconInfoWithId[][]>(() => getIconsInRows(props.allIconInfos, itemPerRow.value));

const icon = computed<string>({
    get: () => props.modelValue,
    set: (value: string) => emit('update:modelValue', value)
});

function hasSelectedIcon(row: IconInfoWithId[]): boolean {
    return arrayContainsFieldValue(row, 'id', props.modelValue);
}

function onMenuStateChanged(state: boolean): void {
    if (state) {
        nextTick(() => {
            if (dropdownMenu.value && dropdownMenu.value.parentElement) {
                scrollToSelectedItem(dropdownMenu.value.parentElement, null, null, '.row-has-selected-item');
            }
        });
    }
}
</script>

<style>
.icon-select-dropdown .icon-item {
    display: grid;
}

.ebk-option-button {
    display: inline-grid;
    min-width: 44px;
    min-height: 44px;
    padding: 0;
    place-items: center;
    border: 0;
    border-radius: 8px;
    color: inherit;
    background: transparent;
}

.ebk-option-button:hover {
    background: rgb(var(--v-theme-on-surface) / 0.06);
}

.ebk-option-button:focus-visible {
    outline: 3px solid rgb(var(--v-theme-primary));
    outline-offset: 2px;
}
</style>
