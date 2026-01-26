<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useStockStore } from '@/stores/stockStore'
import StockSearch from '@/components/stock/StockSearch.vue'
import StockFilters from '@/components/stock/StockFilters.vue'
import StockTable from '@/components/stock/StockTable.vue'
import RecommendationCard from '@/components/recommendation/RecommendationCard.vue'
import Pagination from '@/components/common/Pagination.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import ErrorAlert from '@/components/common/ErrorAlert.vue'

const stockStore = useStockStore()

onMounted(async () => {
  await Promise.all([
    stockStore.fetchStocks(),
    stockStore.fetchActions(),
    stockStore.fetchTopRecommendation()
  ])
})

watch(
  () => stockStore.filter,
  () => stockStore.fetchStocks(),
  { deep: true }
)

function handlePageChange(page: number) {
  stockStore.setPage(page)
}
</script>

<template>
  <div class="space-y-6">
    <div v-if="stockStore.topRecommendation" class="mb-8">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">Today's Top Pick</h2>
      <RecommendationCard :recommendation="stockStore.topRecommendation" featured />
    </div>

    <div class="card">
      <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 mb-6">
        <h2 class="text-lg font-semibold text-gray-900">Stock Analysis</h2>
        <div class="w-full lg:w-80">
          <StockSearch v-model="stockStore.filter.search" @update:model-value="stockStore.setSearch" />
        </div>
      </div>

      <StockFilters
        :sort-by="stockStore.filter.sortBy"
        :sort-order="stockStore.filter.sortOrder"
        :action="stockStore.filter.action"
        :actions="stockStore.actions"
        @sort-change="stockStore.setSort"
        @action-change="stockStore.setAction"
        @reset="stockStore.resetFilters"
        class="mb-6"
      />

      <ErrorAlert
        v-if="stockStore.error"
        :message="stockStore.error"
        @dismiss="stockStore.error = null"
        class="mb-6"
      />

      <LoadingSpinner v-if="stockStore.isLoading" size="lg" class="py-12" />

      <template v-else>
        <StockTable :stocks="stockStore.stocks" />

        <div v-if="stockStore.pagination.totalPages > 1" class="mt-6 pt-6 border-t border-gray-200">
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <p class="text-sm text-gray-500">
              Showing {{ (stockStore.pagination.page - 1) * stockStore.pagination.limit + 1 }} to
              {{ Math.min(stockStore.pagination.page * stockStore.pagination.limit, stockStore.pagination.totalCount) }}
              of {{ stockStore.pagination.totalCount }} results
            </p>
            <Pagination
              :current-page="stockStore.pagination.page"
              :total-pages="stockStore.pagination.totalPages"
              :has-next="stockStore.pagination.hasNext"
              :has-prev="stockStore.pagination.hasPrev"
              @page-change="handlePageChange"
            />
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
