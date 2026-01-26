<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStockStore } from '@/stores/stockStore'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import ErrorAlert from '@/components/common/ErrorAlert.vue'

const route = useRoute()
const router = useRouter()
const stockStore = useStockStore()

const stockId = computed(() => route.params.id as string)

onMounted(async () => {
  if (stockId.value) {
    await stockStore.fetchStockById(stockId.value)
  }
})

function goBack() {
  router.back()
}

function formatPrice(price: number): string {
  if (!price) return '-'
  return `$${price.toFixed(2)}`
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    weekday: 'long',
    month: 'long',
    day: 'numeric',
    year: 'numeric'
  })
}

function calculateChange(from: number, to: number): string {
  if (!from || !to) return '-'
  const change = ((to - from) / from) * 100
  const sign = change > 0 ? '+' : ''
  return `${sign}${change.toFixed(1)}%`
}

function getChangeClass(from: number, to: number): string {
  if (!from || !to) return 'text-gray-500'
  return to > from ? 'text-green-600' : to < from ? 'text-red-600' : 'text-gray-500'
}
</script>

<template>
  <div>
    <button
      @click="goBack"
      class="flex items-center gap-2 text-gray-600 hover:text-gray-900 mb-6"
    >
      <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
      </svg>
      Back to Stocks
    </button>

    <LoadingSpinner v-if="stockStore.isLoading" size="lg" class="py-12" />

    <ErrorAlert
      v-else-if="stockStore.error"
      :message="stockStore.error"
      @dismiss="stockStore.error = null"
    />

    <template v-else-if="stockStore.selectedStock">
      <div class="card">
        <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-4 mb-8">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">{{ stockStore.selectedStock.ticker }}</h1>
            <p class="text-lg text-gray-500">{{ stockStore.selectedStock.company }}</p>
          </div>
          <span class="badge badge-blue text-sm">{{ stockStore.selectedStock.action }}</span>
        </div>

        <div class="grid md:grid-cols-2 gap-8">
          <div class="space-y-6">
            <div>
              <h2 class="text-sm font-medium text-gray-500 uppercase mb-2">Rating Change</h2>
              <div class="flex items-center gap-4">
                <div class="text-center">
                  <p class="text-xs text-gray-400 mb-1">From</p>
                  <p class="text-lg font-semibold text-gray-600">
                    {{ stockStore.selectedStock.ratingFrom || '-' }}
                  </p>
                </div>
                <svg class="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
                </svg>
                <div class="text-center">
                  <p class="text-xs text-gray-400 mb-1">To</p>
                  <p class="text-lg font-semibold text-gray-900">
                    {{ stockStore.selectedStock.ratingTo || '-' }}
                  </p>
                </div>
              </div>
            </div>

            <div>
              <h2 class="text-sm font-medium text-gray-500 uppercase mb-2">Target Price</h2>
              <div class="flex items-center gap-4">
                <div class="text-center">
                  <p class="text-xs text-gray-400 mb-1">From</p>
                  <p class="text-lg font-semibold text-gray-600">
                    {{ formatPrice(stockStore.selectedStock.targetFrom) }}
                  </p>
                </div>
                <svg class="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
                </svg>
                <div class="text-center">
                  <p class="text-xs text-gray-400 mb-1">To</p>
                  <p class="text-lg font-semibold text-gray-900">
                    {{ formatPrice(stockStore.selectedStock.targetTo) }}
                  </p>
                </div>
                <div
                  v-if="stockStore.selectedStock.targetFrom && stockStore.selectedStock.targetTo"
                  :class="['text-lg font-bold', getChangeClass(stockStore.selectedStock.targetFrom, stockStore.selectedStock.targetTo)]"
                >
                  {{ calculateChange(stockStore.selectedStock.targetFrom, stockStore.selectedStock.targetTo) }}
                </div>
              </div>
            </div>
          </div>

          <div class="space-y-6">
            <div>
              <h2 class="text-sm font-medium text-gray-500 uppercase mb-2">Brokerage</h2>
              <p class="text-lg font-semibold text-gray-900">{{ stockStore.selectedStock.brokerage }}</p>
            </div>

            <div>
              <h2 class="text-sm font-medium text-gray-500 uppercase mb-2">Analysis Date</h2>
              <p class="text-lg font-semibold text-gray-900">
                {{ formatDate(stockStore.selectedStock.createdAt) }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </template>

    <div v-else class="text-center py-12 text-gray-500">
      Stock not found
    </div>
  </div>
</template>
