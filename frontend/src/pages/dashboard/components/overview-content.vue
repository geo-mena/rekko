<script setup lang="ts">
import type { RecommendationFilter } from '@/services/api/recommendations.api'

import { useGetDashboardStatsQuery } from '@/services/api/dashboard.api'
import { useGetRecommendationsQuery, useGetTopRecommendationQuery } from '@/services/api/recommendations.api'

import ActionDistributionChart from './action-distribution-chart.vue'
import KpiCards from './kpi-cards.vue'
import TopPicksChart from './top-picks-chart.vue'

const { data: stats, isLoading: statsLoading } = useGetDashboardStatsQuery()
const { data: topPick, isLoading: topPickLoading } = useGetTopRecommendationQuery()

const recommendationsFilter = ref<RecommendationFilter>({ limit: 5 })
const { data: recommendations, isLoading: recsLoading } = useGetRecommendationsQuery(recommendationsFilter)
</script>

<template>
    <div class="space-y-4">
        <KpiCards
            :stats="stats"
            :top-pick="topPick"
            :recommendations="recommendations"
            :loading="statsLoading || topPickLoading || recsLoading"
        />

        <div class="grid gap-4 md:grid-cols-2">
            <ActionDistributionChart
                :data="stats?.actionDistribution ?? []"
                :loading="statsLoading"
            />
            <TopPicksChart
                :data="recommendations ?? []"
                :loading="recsLoading"
            />
        </div>
    </div>
</template>
