<script setup lang="ts">
import type { DataTableProps } from '@/components/data-table/types'

import DataTable from '@/components/data-table/data-table.vue'
import { generateVueTable } from '@/components/data-table/use-generate-vue-table'

import type { StockRecommendation } from '../data/schema'

import DataTableToolbar from './data-table-toolbar.vue'

const props = defineProps<DataTableProps<StockRecommendation>>()

const search = defineModel<string>('search', { default: '' })

const table = generateVueTable<StockRecommendation>(props)
</script>

<template>
    <DataTable :columns :data :loading :table>
        <template #toolbar>
            <DataTableToolbar v-model:search="search" :table class="w-full overflow-x-auto" />
        </template>
    </DataTable>
</template>
