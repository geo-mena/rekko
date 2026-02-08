<script setup lang="ts">
import { computed } from 'vue'

import type { StockRecommendation } from '@/pages/recommendations/data/schema'

const props = defineProps<{
    data: StockRecommendation[]
    loading: boolean
}>()

const sortedData = computed(() => {
    return [...props.data]
        .sort((a, b) => b.upsidePotential - a.upsidePotential)
        .slice(0, 5)
})

const maxUpside = computed(() => {
    if (sortedData.value.length === 0)
        return 100
    return Math.max(...sortedData.value.map(d => d.upsidePotential))
})

function getBarWidth(upside: number): string {
    return `${(upside / maxUpside.value) * 100}%`
}

function getBarGradient(upside: number): string {
    return upside >= 0
        ? 'linear-gradient(to right, oklch(0.696 0.17 162.48), oklch(0.74 0.15 162))'
        : 'linear-gradient(to right, oklch(0.58 0.18 25), oklch(0.64 0.15 25))'
}

function getLegendColor(upside: number): string {
    return upside >= 0 ? 'oklch(0.696 0.17 162.48)' : 'oklch(0.55 0.15 25)'
}

function formatUpside(value: number): string {
    const sign = value >= 0 ? '+' : ''
    return `${sign}${value.toFixed(1)}%`
}


</script>

<template>
    <UiCard class="h-full">
        <UiCardHeader>
            <UiCardTitle>Top Picks by Upside</UiCardTitle>
            <UiCardDescription>Stocks with highest upside potential</UiCardDescription>
        </UiCardHeader>
        <UiCardContent>
            <div v-if="loading" class="h-[300px] flex items-center justify-center">
                <UiSkeleton class="h-full w-full" />
            </div>
            <div v-else-if="data.length === 0" class="h-[300px] flex items-center justify-center text-muted-foreground">
                No recommendation data available
            </div>
            <div v-else class="h-[300px] flex flex-col">
                <div class="flex flex-1 flex-col justify-center space-y-3">
                    <div
                        v-for="item in sortedData"
                        :key="item.stock.id"
                        class="flex items-center gap-3"
                    >
                        <div class="w-16 shrink-0 text-sm font-medium">
                            {{ item.stock.ticker }}
                        </div>
                        <div class="relative h-8 flex-1 overflow-hidden rounded bg-muted">
                            <div
                                class="absolute inset-y-0 left-0 flex items-center rounded transition-all duration-500"
                                :style="{
                                    width: getBarWidth(item.upsidePotential),
                                    background: getBarGradient(item.upsidePotential),
                                    boxShadow: `0 0 12px ${getLegendColor(item.upsidePotential)}40`,
                                }"
                            >
                                <span
                                    class="absolute right-2 text-xs font-semibold text-white"
                                    :class="{ 'text-foreground right-auto left-2': item.upsidePotential < maxUpside * 0.2 }"
                                >
                                    {{ formatUpside(item.upsidePotential) }}
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="flex flex-wrap items-center gap-4 pt-4 border-t mt-4">
                    <div
                        v-for="item in sortedData"
                        :key="item.stock.id"
                        class="flex items-center gap-2 text-xs"
                    >
                        <div
                            class="w-3 h-3 rounded-sm shrink-0"
                            :style="{ background: getBarGradient(item.upsidePotential) }"
                        />
                        <span class="font-medium">{{ item.stock.ticker }}</span>
                        <span class="text-muted-foreground">{{ formatUpside(item.upsidePotential) }}</span>
                    </div>
                </div>
            </div>
        </UiCardContent>
    </UiCard>
</template>
