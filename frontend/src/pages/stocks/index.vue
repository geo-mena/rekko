<script setup lang="ts">
import type { ServerPagination } from '@/components/data-table/types'
import type { StockFilter } from '@/services/api/stocks.api'

import { BasicPage } from '@/components/global-layout'
import { DEFAULT_PAGE_SIZE } from '@/constants/pagination'
import { useGetStocksQuery } from '@/services/api/stocks.api'

import { columns } from './components/columns'
import DataTable from './components/data-table.vue'

const filter = ref<StockFilter>({ page: 1, limit: DEFAULT_PAGE_SIZE })

const { data, isLoading } = useGetStocksQuery(filter)

const stocks = computed(() => data.value?.data ?? [])

const serverPagination = computed<ServerPagination>(() => ({
  page: data.value?.page ?? 1,
  pageSize: data.value?.limit ?? DEFAULT_PAGE_SIZE,
  total: data.value?.totalCount ?? 0,
  onPageChange: (page: number) => {
    filter.value = { ...filter.value, page }
  },
  onPageSizeChange: (limit: number) => {
    filter.value = { ...filter.value, limit, page: 1 }
  },
}))
</script>

<template>
  <BasicPage
    title="Stocks"
    description="Stock analyst recommendations and price targets"
    sticky
  >
    <div class="overflow-x-auto">
      <DataTable :loading="isLoading" :data="stocks" :columns="columns" :server-pagination="serverPagination" />
    </div>
  </BasicPage>
</template>
