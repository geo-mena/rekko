<script setup lang="ts">
import { useRouter } from 'vue-router'
import type { Stock } from '@/types/stock'

defineProps<{
  stocks: Stock[]
}>()

const router = useRouter()

function viewStock(stock: Stock) {
  router.push({ name: 'stock-detail', params: { id: stock.id } })
}

function formatPrice(price: number): string {
  if (!price) return '-'
  return `$${price.toFixed(2)}`
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  })
}

function getActionClass(action: string): string {
  const actionLower = action.toLowerCase()
  if (actionLower.includes('upgraded') || actionLower.includes('raised')) {
    return 'badge-green'
  }
  if (actionLower.includes('downgraded') || actionLower.includes('lowered')) {
    return 'badge-red'
  }
  if (actionLower.includes('initiated')) {
    return 'badge-blue'
  }
  return 'badge-yellow'
}
</script>

<template>
  <div class="overflow-x-auto">
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
        <tr>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Ticker
          </th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Company
          </th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Brokerage
          </th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Action
          </th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Rating
          </th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Target
          </th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Date
          </th>
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        <tr
          v-for="stock in stocks"
          :key="stock.id"
          @click="viewStock(stock)"
          class="hover:bg-gray-50 cursor-pointer transition-colors"
        >
          <td class="px-6 py-4 whitespace-nowrap">
            <span class="font-semibold text-blue-600">{{ stock.ticker }}</span>
          </td>
          <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
            {{ stock.company }}
          </td>
          <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
            {{ stock.brokerage }}
          </td>
          <td class="px-6 py-4 whitespace-nowrap">
            <span :class="['badge', getActionClass(stock.action)]">
              {{ stock.action }}
            </span>
          </td>
          <td class="px-6 py-4 whitespace-nowrap text-sm">
            <span class="text-gray-500">{{ stock.ratingFrom || '-' }}</span>
            <span v-if="stock.ratingFrom && stock.ratingTo" class="mx-1 text-gray-400">→</span>
            <span class="text-gray-900 font-medium">{{ stock.ratingTo || '-' }}</span>
          </td>
          <td class="px-6 py-4 whitespace-nowrap text-sm">
            <span class="text-gray-500">{{ formatPrice(stock.targetFrom) }}</span>
            <span v-if="stock.targetFrom && stock.targetTo" class="mx-1 text-gray-400">→</span>
            <span class="text-gray-900 font-medium">{{ formatPrice(stock.targetTo) }}</span>
          </td>
          <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
            {{ formatDate(stock.createdAt) }}
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="stocks.length === 0" class="text-center py-12 text-gray-500">
      No stocks found. Try adjusting your filters or sync data.
    </div>
  </div>
</template>
