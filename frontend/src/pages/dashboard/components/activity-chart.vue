<script setup lang="ts">
import { VisArea, VisAxis, VisLine, VisXYContainer } from '@unovis/vue'

import type { ChartConfig } from '@/components/ui/chart'

import { ChartContainer, ChartCrosshair, ChartTooltip, ChartTooltipContent, componentToString } from '@/components/ui/chart'

import type { DailyActivity } from '../data/schema'

const props = defineProps<{
    data: DailyActivity[]
    loading: boolean
}>()

const chartConfig = {
    count: {
        label: 'Activity',
        color: 'hsl(var(--chart-1))',
    },
} satisfies ChartConfig

function x(_: DailyActivity, i: number) {
    return i
}

function y(d: DailyActivity) {
    return d.count
}

function formatDate(d: DailyActivity) {
    const date = new Date(d.date)
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

function xTickFormat(i: number) {
    if (!props.data[i])
        return ''
    return formatDate(props.data[i])
}

function tooltipLabelFormatter(d: number | Date) {
    if (typeof d !== 'number' || !props.data[d])
        return ''
    return formatDate(props.data[d])
}
</script>

<template>
    <UiCard>
        <UiCardHeader>
            <UiCardTitle>Activity (Last 30 Days)</UiCardTitle>
            <UiCardDescription>Daily recommendation activity over time</UiCardDescription>
        </UiCardHeader>
        <UiCardContent>
            <div v-if="loading" class="h-[300px] flex items-center justify-center">
                <UiSkeleton class="h-full w-full" />
            </div>
            <div v-else-if="data.length === 0" class="h-[300px] flex items-center justify-center text-muted-foreground">
                No activity data available
            </div>
            <ChartContainer v-else :config="chartConfig" class="h-[300px] w-full" :cursor="true">
                <VisXYContainer :data="data" :margin="{ top: 10, right: 10, bottom: 30, left: 40 }">
                    <VisArea
                        :x="x"
                        :y="y"
                        color="hsl(var(--chart-1))"
                        :opacity="0.2"
                        curve-type="linear"
                    />
                    <VisLine
                        :x="x"
                        :y="y"
                        color="hsl(var(--chart-1))"
                        :line-width="2"
                        curve-type="linear"
                    />
                    <VisAxis
                        type="x"
                        :tick-format="xTickFormat"
                        :num-ticks="6"
                        :grid-line="false"
                    />
                    <VisAxis
                        type="y"
                        :num-ticks="5"
                        :grid-line="true"
                    />

                    <ChartCrosshair
                        color="hsl(var(--chart-1))"
                        :template="componentToString(chartConfig, ChartTooltipContent, { labelFormatter: tooltipLabelFormatter })"
                    />
                    <ChartTooltip />
                </VisXYContainer>
            </ChartContainer>
        </UiCardContent>
    </UiCard>
</template>
