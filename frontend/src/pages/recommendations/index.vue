<script setup lang="ts">
import { BasicPage } from '@/components/global-layout'
import { useGetRecommendationsQuery } from '@/services/api/recommendations.api'

import { columns } from './components/columns'
import DataTable from './components/data-table.vue'

const limit = ref(50)
const { data, isLoading } = useGetRecommendationsQuery(limit)
const recommendations = computed(() => data.value ?? [])
</script>

<template>
  <BasicPage
    title="Recommendations"
    description="AI-powered stock recommendations based on analyst data"
    sticky
  >
    <div class="overflow-x-auto">
      <DataTable :loading="isLoading" :data="recommendations" :columns="columns" />
    </div>
  </BasicPage>
</template>
