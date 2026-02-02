<script setup lang="ts">
import type { Table } from '@tanstack/vue-table'

import { X } from 'lucide-vue-next'
import { computed } from 'vue'

import DataTableViewOptions from '@/components/data-table/view-options.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

import type { StockRecommendation } from '../data/schema'

interface DataTableToolbarProps {
  table: Table<StockRecommendation>
}

const props = defineProps<DataTableToolbarProps>()

const isFiltered = computed(() => props.table.getState().columnFilters.length > 0)
</script>

<template>
  <div class="flex items-center justify-between">
    <div class="flex items-center flex-1 space-x-2">
      <Input
        placeholder="Filter by ticker..."
        :model-value="(table.getColumn('ticker')?.getFilterValue() as string) ?? ''"
        class="h-8 w-[150px] lg:w-[250px]"
        @input="table.getColumn('ticker')?.setFilterValue($event.target.value)"
      />

      <Button
        v-if="isFiltered"
        variant="ghost"
        class="h-8 px-2 lg:px-3"
        @click="table.resetColumnFilters()"
      >
        Reset
        <X class="size-4 ml-2" />
      </Button>
    </div>
    <DataTableViewOptions :table="table" />
  </div>
</template>
