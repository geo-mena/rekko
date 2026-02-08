import type { ColumnDef } from '@tanstack/vue-table'

import { ArrowDown, ArrowUp, DollarSign, Minus, Star, TrendingDown, TrendingUp, Users } from 'lucide-vue-next'
import { h } from 'vue'

import DataTableColumnHeader from '@/components/data-table/column-header.vue'
import Badge from '@/components/ui/badge/Badge.vue'

import type { StockRecommendation } from '../data/schema'

import DataTableRowActions from './data-table-row-actions.vue'
import DataTableTickerCell from './data-table-ticker-cell.vue'

function getScoreColor(score: number): string {
    if (score >= 7) {
        return '#16a34a'
    }
    if (score >= 5) {
        return '#ffa500'
    }
    return '#dc2626'
}

export const columns: ColumnDef<StockRecommendation>[] = [
    {
        id: 'ticker',
        header: ({ column }) => h(DataTableColumnHeader<StockRecommendation>, { column, title: 'Ticker' }),
        cell: ({ row }) => h(DataTableTickerCell, { row }),
        enableSorting: true,
        enableHiding: false,
    },
    {
        id: 'company',
        header: ({ column }) => h(DataTableColumnHeader<StockRecommendation>, { column, title: 'Company' }),
        cell: ({ row }) => h('div', { class: 'max-w-[200px] truncate' }, row.original.stock.company),
        enableSorting: true,
    },
    {
        accessorKey: 'score',
        header: ({ column }) => h(DataTableColumnHeader<StockRecommendation>, { column, title: 'Score' }),
        cell: ({ row }) => {
            const score = row.getValue('score') as number
            const color = getScoreColor(score)
            const style = { color }

            return h(Badge, {
                class: 'flex items-center',
                style,
                variant: 'secondary',
            }, () => [
                h(Star, { class: 'mr-2 h-4 w-4', style }),
                h('span', `${score}/10`),
            ])
        },
        filterFn: (row, id, value: string[]) => {
            const score = row.getValue(id) as number
            return value.some((v) => {
                if (v === 'high')
                    return score >= 7
                if (v === 'medium')
                    return score >= 5 && score < 7
                if (v === 'low')
                    return score >= 1 && score < 5
                return false
            })
        },
        enableSorting: true,
    },
    {
        id: 'currentPrice',
        header: ({ column }) => h(DataTableColumnHeader<StockRecommendation>, { column, title: 'Price' }),
        cell: ({ row }) => {
            const md = row.original.marketData
            if (!md)
                return h('span', { class: 'text-sm text-muted-foreground italic', title: 'Market data unavailable for this ticker' }, 'not available')

            const trendIcon = md.dayChangePercent >= 0 ? TrendingUp : TrendingDown
            const colorClass = md.dayChangePercent >= 0 ? 'text-emerald-600' : 'text-red-600'
            const changeFormatted = `${md.dayChangePercent >= 0 ? '+' : ''}${md.dayChangePercent.toFixed(2)}%`

            return h('div', { class: 'flex flex-col gap-0.5' }, [
                h('div', { class: 'flex items-center gap-1' }, [
                    h(DollarSign, { class: 'size-3.5 text-muted-foreground' }),
                    h('span', { class: 'font-medium' }, md.currentPrice.toFixed(2)),
                ]),
                h('div', { class: `flex items-center gap-1 text-xs ${colorClass}` }, [
                    h(trendIcon, { class: 'size-3' }),
                    h('span', changeFormatted),
                ]),
            ])
        },
        enableSorting: false,
    },
    {
        accessorKey: 'upsidePotential',
        header: ({ column }) => h(DataTableColumnHeader<StockRecommendation>, { column, title: 'Upside Potential' }),
        cell: ({ row }) => {
            const value = row.getValue('upsidePotential') as number
            const isPositive = value > 0
            const isNegative = value < 0
            const icon = isPositive ? ArrowUp : isNegative ? ArrowDown : Minus
            const colorClass = isPositive
                ? 'text-emerald-600'
                : isNegative
                    ? 'text-red-600'
                    : 'text-muted-foreground'
            const formatted = `${isPositive ? '+' : ''}${value.toFixed(1)}%`

            return h('div', { class: 'flex items-center gap-1' }, [
                h(icon, { class: `size-3.5 ${colorClass}` }),
                h('span', { class: `font-medium ${colorClass}` }, formatted),
            ])
        },
        enableSorting: true,
    },
    {
        accessorKey: 'analystCount',
        header: ({ column }) => h(DataTableColumnHeader<StockRecommendation>, { column, title: 'Analysts' }),
        cell: ({ row }) => {
            const count = row.getValue('analystCount') as number
            return h(Badge, {
                class: 'flex items-center',
                variant: 'outline',
            }, () => [
                h(Users, { class: 'mr-1.5 h-3.5 w-3.5' }),
                h('span', {}, `${count}`),
            ])
        },
        enableSorting: true,
    },
    {
        id: 'reasons',
        header: ({ column }) => h(DataTableColumnHeader<StockRecommendation>, { column, title: 'Reasons' }),
        cell: ({ row }) => {
            const count = row.original.reasons.length
            return h('div', { class: 'text-sm text-muted-foreground' }, `${count} reason${count !== 1 ? 's' : ''}`)
        },
        enableSorting: false,
    },
    {
        id: 'actions',
        cell: ({ row }) => h(DataTableRowActions, { row }),
    },
]
