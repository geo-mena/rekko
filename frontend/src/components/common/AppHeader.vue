<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { useStockStore } from '@/stores/stockStore'

const stockStore = useStockStore()

async function handleSync() {
  await stockStore.syncStocks()
}
</script>

<template>
  <header class="bg-white border-b border-gray-200 sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <div class="flex items-center gap-8">
          <RouterLink to="/" class="text-xl font-bold text-gray-900">
            StockRecommender
          </RouterLink>
          <nav class="hidden md:flex gap-6">
            <RouterLink
              to="/"
              class="text-gray-600 hover:text-gray-900 transition-colors"
              active-class="text-blue-600 font-medium"
            >
              Stocks
            </RouterLink>
            <RouterLink
              to="/recommendations"
              class="text-gray-600 hover:text-gray-900 transition-colors"
              active-class="text-blue-600 font-medium"
            >
              Recommendations
            </RouterLink>
          </nav>
        </div>
        <button
          @click="handleSync"
          :disabled="stockStore.isSyncing"
          class="btn btn-primary flex items-center gap-2"
        >
          <svg
            v-if="stockStore.isSyncing"
            class="animate-spin h-4 w-4"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          <span>{{ stockStore.isSyncing ? 'Syncing...' : 'Sync Data' }}</span>
        </button>
      </div>
    </div>
  </header>
</template>
