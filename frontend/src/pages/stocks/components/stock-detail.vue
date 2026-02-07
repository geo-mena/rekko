<script setup lang="ts">
import { ArrowDown, ArrowUp, Minus } from 'lucide-vue-next'

import { useModal } from '@/composables/use-modal'

import type { Stock } from '../data/schema'

interface Props {
    stock: Stock
}
const props = defineProps<Props>()
const _emit = defineEmits<{ close: [] }>()

const { Modal } = useModal()

const priceChange = computed(() => {
    if (!props.stock.targetFrom || !props.stock.targetTo)
        return null
    return ((props.stock.targetTo - props.stock.targetFrom) / props.stock.targetFrom) * 100
})

function formatPrice(value: number): string {
    return value ? `$${value.toFixed(2)}` : '-'
}

function formatDate(value: string): string {
    if (!value)
        return '-'
    return new Date(value).toLocaleDateString('en-US', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric',
    })
}
</script>

<template>
    <div>
        <component :is="Modal.Header">
            <component :is="Modal.Title">
                <span class="font-semibold">{{ stock.ticker }}</span>
                <span class="ml-2 text-muted-foreground">{{ stock.company }}</span>
            </component>
            <component :is="Modal.Description">
                Analysis by {{ stock.brokerage }}
            </component>
        </component>

        <div class="grid grid-cols-2 gap-4 py-4">
            <div class="space-y-1">
                <p class="text-sm text-muted-foreground">
                    Rating Change
                </p>
                <p class="text-sm font-medium">
                    {{ stock.ratingFrom || '-' }} → {{ stock.ratingTo || '-' }}
                </p>
            </div>

            <div class="space-y-1">
                <p class="text-sm text-muted-foreground">
                    Target Price
                </p>
                <div class="flex items-center gap-1">
                    <span class="text-sm font-medium">
                        {{ formatPrice(stock.targetFrom) }} → {{ formatPrice(stock.targetTo) }}
                    </span>
                    <template v-if="priceChange !== null">
                        <ArrowUp v-if="priceChange > 0" class="size-4 text-emerald-600" />
                        <ArrowDown v-else-if="priceChange < 0" class="size-4 text-destructive" />
                        <Minus v-else class="size-4 text-muted-foreground" />
                        <span
                            class="text-xs font-medium"
                            :class="{
                                'text-emerald-600': priceChange > 0,
                                'text-destructive': priceChange < 0,
                                'text-muted-foreground': priceChange === 0,
                            }"
                        >
                            {{ priceChange > 0 ? '+' : '' }}{{ priceChange.toFixed(1) }}%
                        </span>
                    </template>
                </div>
            </div>

            <div class="space-y-1">
                <p class="text-sm text-muted-foreground">
                    Brokerage
                </p>
                <p class="text-sm font-medium">
                    {{ stock.brokerage }}
                </p>
            </div>

            <div class="space-y-1">
                <p class="text-sm text-muted-foreground">
                    Date
                </p>
                <p class="text-sm font-medium">
                    {{ formatDate(stock.createdAt) }}
                </p>
            </div>
        </div>
    </div>
</template>
