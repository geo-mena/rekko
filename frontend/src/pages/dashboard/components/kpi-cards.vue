<script setup lang="ts">
import { Award, ChartNoAxesCombined, Landmark, LayoutGrid } from 'lucide-vue-next'

import type { StockRecommendation } from '@/pages/recommendations/data/schema'

import type { DashboardStats } from '../data/schema'

const props = defineProps<{
    stats?: DashboardStats | null
    topPick?: StockRecommendation | null
    recommendations?: StockRecommendation[] | null
    loading: boolean
}>()

const avgUpside = computed(() => {
    if (!props.recommendations || props.recommendations.length === 0)
        return '0.0'
    const sum = props.recommendations.reduce((acc, r) => acc + r.upsidePotential, 0)
    return (sum / props.recommendations.length).toFixed(1)
})
</script>

<template>
    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <UiCard style="background: linear-gradient(135deg, oklch(0.6 0.15 250 / 0.03) 0%, transparent 50%)">
            <UiCardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
                <UiCardTitle class="text-sm font-medium">
                    Total Stocks
                </UiCardTitle>
                <div class="flex items-center justify-center size-9 rounded-lg bg-blue-500/10">
                    <LayoutGrid class="size-5" style="color: oklch(0.6 0.15 250)" />
                </div>
            </UiCardHeader>
            <UiCardContent>
                <template v-if="loading">
                    <UiSkeleton class="h-8 w-16 mb-1" />
                    <UiSkeleton class="h-3 w-28" />
                </template>
                <template v-else>
                    <div class="text-2xl font-bold">
                        {{ stats?.totalStocks?.toLocaleString() ?? '—' }}
                    </div>
                    <p class="text-xs text-muted-foreground">
                        Tracked analyst ratings
                    </p>
                </template>
            </UiCardContent>
        </UiCard>

        <UiCard style="background: linear-gradient(135deg, oklch(0.75 0.15 85 / 0.03) 0%, transparent 50%)">
            <UiCardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
                <UiCardTitle class="text-sm font-medium">
                    Top Pick
                </UiCardTitle>
                <div class="flex items-center justify-center size-9 rounded-lg bg-amber-500/10">
                    <Award class="size-5" style="color: oklch(0.75 0.15 85)" />
                </div>
            </UiCardHeader>
            <UiCardContent>
                <template v-if="loading">
                    <UiSkeleton class="h-8 w-20 mb-1" />
                    <UiSkeleton class="h-3 w-32" />
                </template>
                <template v-else>
                    <div class="text-2xl font-bold">
                        {{ topPick?.stock.ticker ?? '—' }}
                    </div>
                    <p v-if="topPick" class="text-xs text-emerald-600 dark:text-emerald-400">
                        +{{ topPick.upsidePotential.toFixed(1) }}% upside potential
                    </p>
                    <p v-else class="text-xs text-muted-foreground">
                        No data available
                    </p>
                </template>
            </UiCardContent>
        </UiCard>

        <UiCard style="background: linear-gradient(135deg, oklch(0.696 0.17 162.48 / 0.03) 0%, transparent 50%)">
            <UiCardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
                <UiCardTitle class="text-sm font-medium">
                    Avg Upside
                </UiCardTitle>
                <div class="flex items-center justify-center size-9 rounded-lg" style="background: oklch(0.696 0.17 162.48 / 0.1)">
                    <ChartNoAxesCombined class="size-5" style="color: oklch(0.696 0.17 162.48)" />
                </div>
            </UiCardHeader>
            <UiCardContent>
                <template v-if="loading">
                    <UiSkeleton class="h-8 w-14 mb-1" />
                    <UiSkeleton class="h-3 w-36" />
                </template>
                <template v-else>
                    <div class="text-2xl font-bold text-emerald-600 dark:text-emerald-400">
                        {{ avgUpside }}%
                    </div>
                    <p class="text-xs text-muted-foreground">
                        Across top recommendations
                    </p>
                </template>
            </UiCardContent>
        </UiCard>

        <UiCard style="background: linear-gradient(135deg, oklch(0.6 0.2 290 / 0.03) 0%, transparent 50%)">
            <UiCardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
                <UiCardTitle class="text-sm font-medium">
                    Brokerages
                </UiCardTitle>
                <div class="flex items-center justify-center size-9 rounded-lg bg-violet-500/10">
                    <Landmark class="size-5" style="color: oklch(0.6 0.2 290)" />
                </div>
            </UiCardHeader>
            <UiCardContent>
                <template v-if="loading">
                    <UiSkeleton class="h-8 w-10 mb-1" />
                    <UiSkeleton class="h-3 w-24" />
                </template>
                <template v-else>
                    <div class="text-2xl font-bold">
                        {{ stats?.brokerageDistribution?.length ?? '—' }}
                    </div>
                    <p class="text-xs text-muted-foreground">
                        Active analyst firms
                    </p>
                </template>
            </UiCardContent>
        </UiCard>
    </div>
</template>
