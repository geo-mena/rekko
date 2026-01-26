<script setup lang="ts">
import { useRouter } from 'vue-router'
import type { StockRecommendation } from '@/types/stock'

const props = defineProps<{
  recommendation: StockRecommendation
  featured?: boolean
}>()

const router = useRouter()

function viewStock() {
  router.push({ name: 'stock-detail', params: { id: props.recommendation.stock.id } })
}

function formatPrice(price: number): string {
  if (!price) return '-'
  return `$${price.toFixed(2)}`
}

function formatUpside(upside: number): string {
  if (!upside) return '-'
  return `+${upside.toFixed(1)}%`
}
</script>

<template>
  <div
    @click="viewStock"
    :class="[
      'card cursor-pointer transition-all hover:shadow-md',
      featured ? 'border-blue-500 border-2' : ''
    ]"
  >
    <div v-if="featured" class="mb-4">
      <span class="badge badge-blue">Top Pick</span>
    </div>

    <div class="flex justify-between items-start mb-4">
      <div>
        <h3 class="text-xl font-bold text-gray-900">{{ recommendation.stock.ticker }}</h3>
        <p class="text-sm text-gray-500">{{ recommendation.stock.company }}</p>
      </div>
      <div class="text-right">
        <div class="text-2xl font-bold text-green-600">{{ formatUpside(recommendation.upsidePotential) }}</div>
        <p class="text-xs text-gray-500">Upside potential</p>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4 mb-4">
      <div>
        <p class="text-xs text-gray-500 uppercase">Target Price</p>
        <p class="font-semibold text-gray-900">{{ formatPrice(recommendation.stock.targetTo) }}</p>
      </div>
      <div>
        <p class="text-xs text-gray-500 uppercase">Rating</p>
        <p class="font-semibold text-gray-900">{{ recommendation.stock.ratingTo || '-' }}</p>
      </div>
      <div>
        <p class="text-xs text-gray-500 uppercase">Brokerage</p>
        <p class="font-semibold text-gray-900">{{ recommendation.stock.brokerage }}</p>
      </div>
      <div>
        <p class="text-xs text-gray-500 uppercase">Score</p>
        <p class="font-semibold text-gray-900">{{ recommendation.score.toFixed(1) }}</p>
      </div>
    </div>

    <div v-if="recommendation.reasons.length > 0" class="border-t pt-4">
      <p class="text-xs text-gray-500 uppercase mb-2">Reasons</p>
      <ul class="space-y-1">
        <li
          v-for="(reason, index) in recommendation.reasons"
          :key="index"
          class="text-sm text-gray-600 flex items-start gap-2"
        >
          <svg class="h-4 w-4 text-green-500 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
          </svg>
          {{ reason }}
        </li>
      </ul>
    </div>
  </div>
</template>
