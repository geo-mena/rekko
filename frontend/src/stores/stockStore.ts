import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { stockService } from '@/services/stockService'
import type { Stock, StockFilter, StockRecommendation, SortField, SortOrder } from '@/types/stock'
import { defaultFilter } from '@/types/stock'

export const useStockStore = defineStore('stock', () => {
  const stocks = ref<Stock[]>([])
  const selectedStock = ref<Stock | null>(null)
  const recommendations = ref<StockRecommendation[]>([])
  const topRecommendation = ref<StockRecommendation | null>(null)
  const actions = ref<string[]>([])

  const filter = ref<StockFilter>({ ...defaultFilter })
  const pagination = ref({
    page: 1,
    limit: 20,
    totalCount: 0,
    totalPages: 0,
    hasNext: false,
    hasPrev: false
  })

  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const isSyncing = ref(false)

  const hasStocks = computed(() => stocks.value.length > 0)

  async function fetchStocks() {
    isLoading.value = true
    error.value = null

    try {
      const response = await stockService.getStocks(filter.value)
      stocks.value = response.data ?? []
      pagination.value = {
        page: response.page,
        limit: response.limit,
        totalCount: response.totalCount,
        totalPages: response.totalPages,
        hasNext: response.hasNext,
        hasPrev: response.hasPrev
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch stocks'
      stocks.value = []
    } finally {
      isLoading.value = false
    }
  }

  async function fetchStockById(id: string) {
    isLoading.value = true
    error.value = null

    try {
      selectedStock.value = await stockService.getStockById(id)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch stock'
      selectedStock.value = null
    } finally {
      isLoading.value = false
    }
  }

  async function fetchRecommendations(limit = 10) {
    isLoading.value = true
    error.value = null

    try {
      recommendations.value = (await stockService.getRecommendations(limit)) ?? []
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch recommendations'
      recommendations.value = []
    } finally {
      isLoading.value = false
    }
  }

  async function fetchTopRecommendation() {
    try {
      topRecommendation.value = await stockService.getTopRecommendation()
    } catch (e) {
      topRecommendation.value = null
    }
  }

  async function fetchActions() {
    try {
      actions.value = (await stockService.getActions()) ?? []
    } catch (e) {
      actions.value = []
    }
  }

  async function syncStocks() {
    isSyncing.value = true
    error.value = null

    try {
      await stockService.triggerSync()
      await fetchStocks()
      await fetchRecommendations()
      await fetchTopRecommendation()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to sync stocks'
    } finally {
      isSyncing.value = false
    }
  }

  function setSearch(search: string) {
    filter.value.search = search
    filter.value.page = 1
  }

  function setSort(sortBy: SortField, sortOrder: SortOrder) {
    filter.value.sortBy = sortBy
    filter.value.sortOrder = sortOrder
    filter.value.page = 1
  }

  function setAction(action: string) {
    filter.value.action = action
    filter.value.page = 1
  }

  function setPage(page: number) {
    filter.value.page = page
  }

  function resetFilters() {
    filter.value = { ...defaultFilter }
  }

  function clearSelectedStock() {
    selectedStock.value = null
  }

  return {
    stocks,
    selectedStock,
    recommendations,
    topRecommendation,
    actions,
    filter,
    pagination,
    isLoading,
    error,
    isSyncing,
    hasStocks,
    fetchStocks,
    fetchStockById,
    fetchRecommendations,
    fetchTopRecommendation,
    fetchActions,
    syncStocks,
    setSearch,
    setSort,
    setAction,
    setPage,
    resetFilters,
    clearSelectedStock
  }
})
