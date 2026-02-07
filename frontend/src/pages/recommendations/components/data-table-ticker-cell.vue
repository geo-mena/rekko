<script setup lang="ts">
import type { Row } from '@tanstack/vue-table'
import type { Component } from 'vue'

import { useModal } from '@/composables/use-modal'

import type { StockRecommendation } from '../data/schema'

interface DataTableTickerCellProps {
    row: Row<StockRecommendation>
}
const props = defineProps<DataTableTickerCellProps>()
const recommendation = computed(() => props.row.original)
const ticker = computed(() => props.row.original.stock.ticker)
const isOpen = ref(false)

const showComponent = shallowRef<Component | null>(null)

function handleClick() {
    showComponent.value = defineAsyncComponent(() => import('./recommendation-detail.vue'))
    isOpen.value = true
}

const { contentClass, Modal } = useModal()
</script>

<template>
    <component :is="Modal.Root" v-model:open="isOpen">
        <button
            type="button"
            class="font-semibold text-primary hover:underline cursor-pointer"
            @click="handleClick"
        >
            {{ ticker }}
        </button>

        <component :is="Modal.Content" :class="contentClass">
            <component :is="showComponent" :recommendation="recommendation" @close="isOpen = false" />
        </component>
    </component>
</template>
