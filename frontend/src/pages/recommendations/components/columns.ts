import type { ColumnDef } from '@tanstack/vue-table'

import { h } from 'vue'

import DataTableColumnHeader from '@/components/data-table/column-header.vue'
import { SelectColumn } from '@/components/data-table/table-columns'
import Badge from '@/components/ui/badge/Badge.vue'

import type { StockRecommendation } from '../data/schema'

import DataTableRowActions from './data-table-row-actions.vue'

function getScoreStyle(score: number): string {
    if (score >= 8)
        return 'border-teal-500 bg-teal-50 text-teal-700 dark:bg-teal-950 dark:text-teal-300'
    if (score >= 5)
        return 'border-sky-500 bg-sky-50 text-sky-700 dark:bg-sky-950 dark:text-sky-300'
    return ''
}

export const columns: ColumnDef<StockRecommendation>[] = [
    SelectColumn as ColumnDef<StockRecommendation>,
    {
        id: 'ticker',
        header: ({ column }) => h(DataTableColumnHeader<StockRecommendation>, { column, title: 'Ticker' }),
        cell: ({ row }) => h('div', { class: 'font-semibold' }, row.original.stock.ticker),
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
            return h(Badge, { class: getScoreStyle(score), variant: 'outline' }, () => `${score}/10`)
        },
        enableSorting: true,
    },
    {
        accessorKey: 'upsidePotential',
        header: ({ column }) => h(DataTableColumnHeader<StockRecommendation>, { column, title: 'Upside Potential' }),
        cell: ({ row }) => {
            const value = row.getValue('upsidePotential') as number
            const formatted = `${value > 0 ? '+' : ''}${value.toFixed(1)}%`
            return h('div', {
                class: value > 0 ? 'text-emerald-600 font-medium' : value < 0 ? 'text-destructive font-medium' : 'text-muted-foreground',
            }, formatted)
        },
        enableSorting: true,
    },
    {
        accessorKey: 'analystCount',
        header: ({ column }) => h(DataTableColumnHeader<StockRecommendation>, { column, title: 'Analysts' }),
        cell: ({ row }) => {
            const count = row.getValue('analystCount') as number
            return h('div', { class: 'text-sm font-medium' }, `${count}`)
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
