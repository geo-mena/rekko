<script setup lang="ts">
import { BasicPage } from '@/components/global-layout'
import type { StockFilter } from '@/services/api/stocks.api'
import { useGetStocksQuery } from '@/services/api/stocks.api'

import { columns } from './components/columns'
import DataTable from './components/data-table.vue'
import StockSync from './components/stock-sync.vue'

const filter = ref<StockFilter>({ page: 1, limit: 20 })

const { data, isLoading } = useGetStocksQuery(filter)
</script>

<template>
  <BasicPage
    title="Stocks"
    description="Stock analyst recommendations and price targets"
    sticky
  >
    <template #actions>
      <StockSync />
    </template>
    <div class="overflow-x-auto">
      <DataTable :loading="isLoading" :data="data?.data ?? []" :columns="columns" />
    </div>
  </BasicPage>
</template>
