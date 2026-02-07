<script setup lang="ts">
import { watchDebounced } from '@vueuse/core'

import type { RecommendationFilter } from '@/services/api/recommendations.api'

import { BasicPage } from '@/components/global-layout'
import { useGetRecommendationsQuery } from '@/services/api/recommendations.api'

import { columns } from './components/columns'
import DataTable from './components/data-table.vue'

const filter = ref<RecommendationFilter>({ limit: 50 })
const search = ref('')

watchDebounced(search, (value) => {
    filter.value = { ...filter.value, search: value }
}, { debounce: 300 })

const { data, isLoading } = useGetRecommendationsQuery(filter)
const recommendations = computed(() => data.value ?? [])
</script>

<template>
    <BasicPage
        title="Top Picks"
        description="AI-powered stock recommendations based on analyst data"
        sticky
    >
        <div class="overflow-x-auto">
            <DataTable v-model:search="search" :loading="isLoading" :data="recommendations" :columns="columns" />
        </div>
    </BasicPage>
</template>
