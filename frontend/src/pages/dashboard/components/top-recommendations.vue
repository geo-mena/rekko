<script setup lang="ts">
import type { StockRecommendation } from '@/pages/recommendations/data/schema'

defineProps<{
  recommendations: StockRecommendation[]
  loading: boolean
}>()

function scoreBadgeClass(score: number): string {
  if (score >= 70) return 'bg-teal-100/30 text-teal-900 dark:text-teal-200 border-teal-200'
  if (score >= 40) return 'bg-sky-200/40 text-sky-900 dark:text-sky-100 border-sky-300'
  return 'bg-neutral-300/40 border-neutral-300'
}
</script>

<template>
  <UiCard>
    <UiCardHeader>
      <UiCardTitle>Top Recommendations</UiCardTitle>
      <UiCardDescription>Highest-scored stocks based on analyst ratings</UiCardDescription>
    </UiCardHeader>
    <UiCardContent>
      <div v-if="loading" class="space-y-4">
        <UiSkeleton v-for="i in 5" :key="i" class="h-12 w-full" />
      </div>
      <div v-else-if="recommendations.length === 0" class="flex items-center justify-center py-8 text-muted-foreground">
        No recommendations available
      </div>
      <div v-else class="space-y-4">
        <div
          v-for="rec in recommendations"
          :key="rec.stock.id"
          class="flex items-center justify-between gap-4"
        >
          <div class="flex flex-col min-w-0">
            <span class="font-semibold text-sm">{{ rec.stock.ticker }}</span>
            <span class="text-xs text-muted-foreground truncate">{{ rec.stock.company }}</span>
          </div>
          <div class="flex items-center gap-3 shrink-0">
            <UiBadge variant="outline" :class="scoreBadgeClass(rec.score)">
              {{ rec.score.toFixed(0) }}
            </UiBadge>
            <span class="text-sm font-medium tabular-nums" :class="rec.upsidePotential >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-destructive'">
              {{ rec.upsidePotential >= 0 ? '+' : '' }}{{ rec.upsidePotential.toFixed(1) }}%
            </span>
          </div>
        </div>
      </div>
    </UiCardContent>
  </UiCard>
</template>
