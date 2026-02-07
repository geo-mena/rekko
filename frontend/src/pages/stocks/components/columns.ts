import type { ColumnDef } from '@tanstack/vue-table'

import { h } from 'vue'

import DataTableColumnHeader from '@/components/data-table/column-header.vue'
import Badge from '@/components/ui/badge/Badge.vue'

import type { Stock } from '../data/schema'

import { actionTypes } from '../data/data'
import DataTableRowActions from './data-table-row-actions.vue'

function formatPrice(value: number): string {
  return value ? `$${value.toFixed(2)}` : '-'
}

function formatDate(value: string): string {
  if (!value) return '-'
  const date = new Date(value)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

export const columns: ColumnDef<Stock>[] = [
  {
    accessorKey: 'ticker',
    header: ({ column }) => h(DataTableColumnHeader<Stock>, { column, title: 'Ticker' }),
    cell: ({ row }) => h('div', { class: 'font-semibold' }, row.getValue('ticker')),
    enableSorting: true,
    enableHiding: false,
  },
  {
    accessorKey: 'company',
    header: ({ column }) => h(DataTableColumnHeader<Stock>, { column, title: 'Company' }),
    cell: ({ row }) => h('div', { class: 'max-w-[200px] truncate' }, row.getValue('company')),
    enableSorting: true,
  },
  {
    accessorKey: 'brokerage',
    header: ({ column }) => h(DataTableColumnHeader<Stock>, { column, title: 'Brokerage' }),
    cell: ({ row }) => h('div', { class: 'max-w-[180px] truncate' }, row.getValue('brokerage')),
    enableSorting: false,
  },
  {
    accessorKey: 'action',
    header: ({ column }) => h(DataTableColumnHeader<Stock>, { column, title: 'Action' }),
    cell: ({ row }) => {
      const action = actionTypes.find(a => a.value === (row.getValue('action') as string).toLowerCase())
      const label = row.getValue('action') as string
      if (!action) return h(Badge, { variant: 'outline' }, () => label)
      return h(Badge, { class: action.style, variant: 'outline' }, () => label)
    },
    filterFn: (row, id, value) => value.includes(row.getValue(id)),
  },
  {
    id: 'rating',
    header: ({ column }) => h(DataTableColumnHeader<Stock>, { column, title: 'Rating' }),
    cell: ({ row }) => {
      const from = row.original.ratingFrom || '-'
      const to = row.original.ratingTo || '-'
      return h('div', { class: 'text-sm' }, `${from} → ${to}`)
    },
    enableSorting: false,
  },
  {
    id: 'targetPrice',
    header: ({ column }) => h(DataTableColumnHeader<Stock>, { column, title: 'Target Price' }),
    cell: ({ row }) => {
      const from = formatPrice(row.original.targetFrom)
      const to = formatPrice(row.original.targetTo)
      return h('div', { class: 'text-sm' }, `${from} → ${to}`)
    },
    enableSorting: false,
  },
  {
    accessorKey: 'createdAt',
    header: ({ column }) => h(DataTableColumnHeader<Stock>, { column, title: 'Date' }),
    cell: ({ row }) => h('div', { class: 'text-sm text-muted-foreground' }, formatDate(row.getValue('createdAt'))),
    enableSorting: true,
  },
  {
    id: 'actions',
    cell: ({ row }) => h(DataTableRowActions, { row }),
  },
]
