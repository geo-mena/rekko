<script setup lang="ts">
import { ChevronRight, TrendingDown, TrendingUp } from 'lucide-vue-next'

import { useModal } from '@/composables/use-modal'

import type { StockRecommendation } from '../data/schema'

interface Props {
    recommendation: StockRecommendation
}
const _props = defineProps<Props>()

const { Modal } = useModal()

const md = computed(() => _props.recommendation.marketData)

function formatPrice(value: number): string {
    return value ? `$${value.toFixed(2)}` : '-'
}

function formatMarketCap(millions: number): string {
    if (millions >= 1000) {
        return `$${(millions / 1000).toFixed(1)}B`
    }
    return `$${millions.toFixed(0)}M`
}
</script>

<template>
    <div>
        <component :is="Modal.Header">
            <component :is="Modal.Title">
                <span class="font-semibold">{{ recommendation.stock.ticker }}</span>
                <span class="ml-2 text-muted-foreground">{{ recommendation.stock.company }}</span>
            </component>
            <component :is="Modal.Description">
                Score: {{ recommendation.score }}/10
                <span v-if="md?.industry" class="ml-2 text-xs text-muted-foreground">
                    &middot; {{ md.industry }}
                </span>
            </component>
        </component>

        <div v-if="!md" class="flex items-center gap-2 rounded-lg border border-dashed p-3 mt-3 text-sm text-muted-foreground">
            <span>Market data not available for this ticker. Score is based on analyst data only.</span>
        </div>
        <div v-if="md" class="flex items-center gap-4 rounded-lg border bg-muted/30 p-3 mt-3">
            <div class="flex-1">
                <p class="text-xs text-muted-foreground">
                    Current Price
                </p>
                <p class="text-lg font-semibold">
                    {{ formatPrice(md.currentPrice) }}
                </p>
            </div>
            <div class="flex-1">
                <p class="text-xs text-muted-foreground">
                    Day Change
                </p>
                <div
                    class="flex items-center gap-1"
                    :class="md.dayChangePercent >= 0 ? 'text-emerald-600' : 'text-red-600'"
                >
                    <component :is="md.dayChangePercent >= 0 ? TrendingUp : TrendingDown" class="size-4" />
                    <span class="font-medium">{{ md.dayChangePercent >= 0 ? '+' : '' }}{{ md.dayChangePercent.toFixed(2) }}%</span>
                </div>
            </div>
            <div class="flex-1">
                <p class="text-xs text-muted-foreground">
                    Day Range
                </p>
                <p class="text-sm font-medium">
                    {{ formatPrice(md.dayLow) }} – {{ formatPrice(md.dayHigh) }}
                </p>
            </div>
            <div v-if="md.marketCap > 0" class="flex-1">
                <p class="text-xs text-muted-foreground">
                    Market Cap
                </p>
                <p class="text-sm font-medium">
                    {{ formatMarketCap(md.marketCap) }}
                </p>
            </div>
        </div>

        <div class="grid grid-cols-2 gap-4 py-4">
            <div class="space-y-1">
                <p class="text-sm text-muted-foreground">
                    Upside Potential
                </p>
                <p
                    class="text-sm font-medium"
                    :class="{
                        'text-emerald-600': recommendation.upsidePotential > 0,
                        'text-destructive': recommendation.upsidePotential < 0,
                        'text-muted-foreground': recommendation.upsidePotential === 0,
                    }"
                >
                    {{ recommendation.upsidePotential > 0 ? '+' : '' }}{{ recommendation.upsidePotential.toFixed(1) }}%
                    <span class="text-xs text-muted-foreground ml-1">{{ md ? '(vs market price)' : '(vs analyst target change)' }}</span>
                </p>
            </div>

            <div class="space-y-1">
                <p class="text-sm text-muted-foreground">
                    Target Price
                </p>
                <p class="text-sm font-medium">
                    {{ formatPrice(recommendation.stock.targetFrom) }} → {{ formatPrice(recommendation.stock.targetTo) }}
                </p>
            </div>

            <div class="space-y-1">
                <p class="text-sm text-muted-foreground">
                    Rating
                </p>
                <p class="text-sm font-medium">
                    {{ recommendation.stock.ratingFrom || '-' }} → {{ recommendation.stock.ratingTo || '-' }}
                </p>
            </div>

            <div class="space-y-1">
                <p class="text-sm text-muted-foreground">
                    Brokerage
                </p>
                <p class="text-sm font-medium">
                    {{ recommendation.stock.brokerage }}
                </p>
            </div>

            <div class="space-y-1">
                <p class="text-sm text-muted-foreground">
                    Analysts
                </p>
                <p class="text-sm font-medium">
                    {{ recommendation.analystCount }} covering
                </p>
            </div>

            <div v-if="md" class="space-y-1">
                <p class="text-sm text-muted-foreground">
                    Prev. Close
                </p>
                <p class="text-sm font-medium">
                    {{ formatPrice(md.previousClose) }}
                </p>
            </div>
        </div>

        <div v-if="recommendation.reasons.length" class="space-y-2 py-2">
            <p class="text-sm font-medium">
                Reasons
            </p>
            <ul class="space-y-1.5">
                <li
                    v-for="(reason, index) in recommendation.reasons"
                    :key="index"
                    class="flex items-start gap-2 text-sm text-muted-foreground"
                >
                    <ChevronRight class="size-4 mt-0.5 shrink-0 text-emerald-600" />
                    <span>{{ reason }}</span>
                </li>
            </ul>
        </div>
    </div>
</template>
