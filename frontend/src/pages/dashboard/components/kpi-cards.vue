<script setup lang="ts">
import { BarChart3, Building2, Star, TrendingUp } from 'lucide-vue-next'

import type { DashboardStats } from '../data/schema'
import type { StockRecommendation } from '@/pages/recommendations/data/schema'

const props = defineProps<{
  stats?: DashboardStats | null
  topPick?: StockRecommendation | null
  recommendations?: StockRecommendation[] | null
  loading: boolean
}>()

const avgUpside = computed(() => {
  if (!props.recommendations || props.recommendations.length === 0) return '0.0'
  const sum = props.recommendations.reduce((acc, r) => acc + r.upsidePotential, 0)
  return (sum / props.recommendations.length).toFixed(1)
})
</script>

<template>
  <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
    <UiCard>
      <UiCardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
        <UiCardTitle class="text-sm font-medium">Total Stocks</UiCardTitle>
        <BarChart3 class="size-4 text-muted-foreground" />
      </UiCardHeader>
      <UiCardContent>
        <div class="text-2xl font-bold">
          <UiSkeleton v-if="loading" class="h-8 w-20" />
          <span v-else>{{ stats?.totalStocks?.toLocaleString() ?? '—' }}</span>
        </div>
        <p class="text-xs text-muted-foreground">Tracked analyst ratings</p>
      </UiCardContent>
    </UiCard>

    <UiCard>
      <UiCardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
        <UiCardTitle class="text-sm font-medium">Top Pick</UiCardTitle>
        <Star class="size-4 text-muted-foreground" />
      </UiCardHeader>
      <UiCardContent>
        <div class="text-2xl font-bold">
          <UiSkeleton v-if="loading" class="h-8 w-20" />
          <span v-else>{{ topPick?.stock.ticker ?? '—' }}</span>
        </div>
        <p v-if="topPick" class="text-xs text-emerald-600 dark:text-emerald-400">
          +{{ topPick.upsidePotential.toFixed(1) }}% upside potential
        </p>
        <p v-else class="text-xs text-muted-foreground">No data available</p>
      </UiCardContent>
    </UiCard>

    <UiCard>
      <UiCardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
        <UiCardTitle class="text-sm font-medium">Avg Upside</UiCardTitle>
        <TrendingUp class="size-4 text-muted-foreground" />
      </UiCardHeader>
      <UiCardContent>
        <div class="text-2xl font-bold">
          <UiSkeleton v-if="loading" class="h-8 w-20" />
          <span v-else class="text-emerald-600 dark:text-emerald-400">{{ avgUpside }}%</span>
        </div>
        <p class="text-xs text-muted-foreground">Across top recommendations</p>
      </UiCardContent>
    </UiCard>

    <UiCard>
      <UiCardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
        <UiCardTitle class="text-sm font-medium">Brokerages</UiCardTitle>
        <Building2 class="size-4 text-muted-foreground" />
      </UiCardHeader>
      <UiCardContent>
        <div class="text-2xl font-bold">
          <UiSkeleton v-if="loading" class="h-8 w-20" />
          <span v-else>{{ stats?.brokerageDistribution?.length ?? '—' }}</span>
        </div>
        <p class="text-xs text-muted-foreground">Active analyst firms</p>
      </UiCardContent>
    </UiCard>
  </div>
</template>
