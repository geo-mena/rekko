<script setup lang="ts">
import { onMounted } from 'vue'
import { useStockStore } from '@/stores/stockStore'
import RecommendationCard from '@/components/recommendation/RecommendationCard.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import ErrorAlert from '@/components/common/ErrorAlert.vue'

const stockStore = useStockStore()

onMounted(async () => {
  await stockStore.fetchRecommendations(10)
})
</script>

<template>
  <div>
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-gray-900">Stock Recommendations</h1>
      <p class="text-gray-500 mt-1">
        Our algorithm analyzes rating changes, price targets, and brokerage consensus to identify the best investment opportunities.
      </p>
    </div>

    <ErrorAlert
      v-if="stockStore.error"
      :message="stockStore.error"
      @dismiss="stockStore.error = null"
      class="mb-6"
    />

    <LoadingSpinner v-if="stockStore.isLoading" size="lg" class="py-12" />

    <template v-else>
      <div v-if="stockStore.recommendations.length > 0" class="space-y-6">
        <RecommendationCard
          v-for="(rec, index) in stockStore.recommendations"
          :key="rec.stock.id"
          :recommendation="rec"
          :featured="index === 0"
        />
      </div>

      <div v-else class="card text-center py-12">
        <svg class="h-12 w-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900 mb-2">No Recommendations Available</h3>
        <p class="text-gray-500">
          Sync stock data to generate recommendations based on analyst ratings and price targets.
        </p>
      </div>
    </template>
  </div>
</template>
