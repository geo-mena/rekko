<script setup lang="ts">
import { useGetDashboardStatsQuery } from '@/services/api/dashboard.api'
import { useGetRecommendationsQuery, useGetTopRecommendationQuery } from '@/services/api/recommendations.api'

import KpiCards from './kpi-cards.vue'
import TopRecommendations from './top-recommendations.vue'

const { data: stats, isLoading: statsLoading } = useGetDashboardStatsQuery()
const { data: topPick, isLoading: topPickLoading } = useGetTopRecommendationQuery()

const recommendationsLimit = ref(5)
const { data: recommendations, isLoading: recsLoading } = useGetRecommendationsQuery(recommendationsLimit)
</script>

<template>
    <div class="space-y-4">
        <KpiCards
            :stats="stats"
            :top-pick="topPick"
            :recommendations="recommendations"
            :loading="statsLoading || topPickLoading || recsLoading"
        />

        <TopRecommendations
            :recommendations="recommendations ?? []"
            :loading="recsLoading"
        />
    </div>
</template>
