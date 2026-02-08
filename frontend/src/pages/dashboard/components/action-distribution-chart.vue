<script setup lang="ts">
import { VisDonut, VisSingleContainer } from '@unovis/vue'
import { ArrowRight } from 'lucide-vue-next'
import { computed } from 'vue'

import type { ChartConfig } from '@/components/ui/chart'

import { ChartContainer } from '@/components/ui/chart'

import type { ActionDistribution } from '../data/schema'

const props = defineProps<{
    data: ActionDistribution[]
    loading: boolean
}>()

const actionColors: Record<string, string> = {
    upgraded: 'oklch(0.696 0.17 162.48)',
    raised: 'oklch(0.74 0.15 162)',
    initiated: 'hsl(221, 83%, 53%)',
    reiterated: 'hsl(262, 83%, 58%)',
    maintained: 'hsl(38, 92%, 50%)',
    downgraded: 'oklch(0.488 0.243 264.376)',
    lowered: 'oklch(0.53 0.2 264)',
    target: 'hsl(215, 14%, 64%)',
}

function getActionColor(action: string): string {
    const normalizedAction = action.toLowerCase()
    for (const [key, color] of Object.entries(actionColors)) {
        if (normalizedAction.includes(key)) {
            return color
        }
    }
    return 'hsl(215, 14%, 64%)'
}

const chartConfig = computed(() => {
    const config: ChartConfig = {}
    props.data.forEach((item) => {
        config[item.action] = {
            label: item.action,
            color: getActionColor(item.action),
        }
    })
    return config
})

const colorAccessor = (d: ActionDistribution) => getActionColor(d.action)
const valueAccessor = (d: ActionDistribution) => d.count

const total = computed(() => props.data.reduce((sum, item) => sum + item.count, 0))
</script>

<template>
    <UiCard class="h-full">
        <UiCardHeader class="flex flex-row items-center justify-between">
            <div>
                <UiCardTitle>Action Distribution</UiCardTitle>
                <UiCardDescription>Distribution of recommendation actions</UiCardDescription>
            </div>
            <RouterLink to="/stocks" class="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors">
                View all
                <ArrowRight class="size-4" />
            </RouterLink>
        </UiCardHeader>
        <UiCardContent>
            <div v-if="loading" class="h-[300px] flex flex-col items-center justify-center">
                <div class="relative flex items-center justify-center">
                    <UiSkeleton class="size-40 rounded-full" />
                    <div class="absolute size-24 rounded-full bg-card" />
                </div>
                <div class="flex flex-wrap items-center justify-center gap-3 pt-6">
                    <div v-for="i in 4" :key="i" class="flex items-center gap-1.5">
                        <UiSkeleton class="size-2 rounded-sm" />
                        <UiSkeleton class="h-3 w-14" />
                        <UiSkeleton class="h-3 w-6" />
                    </div>
                </div>
            </div>
            <div v-else-if="data.length === 0" class="h-[300px] flex items-center justify-center text-muted-foreground">
                No action data available
            </div>
            <ChartContainer v-else :config="chartConfig" class="h-[300px] w-full">
                <VisSingleContainer :data="data">
                    <VisDonut
                        :value="valueAccessor"
                        :color="colorAccessor"
                        :arc-width="50"
                        :corner-radius="4"
                        :central-label="String(total)"
                        central-sub-label="Total"
                    />
                </VisSingleContainer>
                <div class="flex flex-wrap items-center justify-center gap-3 pt-4">
                    <div
                        v-for="item in data"
                        :key="item.action"
                        class="flex items-center gap-1.5 text-xs"
                    >
                        <div
                            class="h-2 w-2 shrink-0 rounded-[2px]"
                            :style="{ backgroundColor: getActionColor(item.action) }"
                        />
                        <span class="text-muted-foreground">{{ item.action }}</span>
                        <span class="font-medium tabular-nums">{{ item.count }}</span>
                    </div>
                </div>
            </ChartContainer>
        </UiCardContent>
    </UiCard>
</template>
