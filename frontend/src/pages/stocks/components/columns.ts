import type { ColumnDef } from '@tanstack/vue-table'

import { ArrowDown, ArrowUp, Landmark, Minus } from 'lucide-vue-next'
import { h } from 'vue'

import DataTableColumnHeader from '@/components/data-table/column-header.vue'
import { Badge } from '@/components/ui/badge'

import type { Stock } from '../data/schema'

import { actionStatuses, getActionType } from '../data/data'
import DataTableRowActions from './data-table-row-actions.vue'
import DataTableTickerCell from './data-table-ticker-cell.vue'

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
    cell: ({ row }) => h(DataTableTickerCell, { row }),
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
    cell: ({ row }) => {
      const brokerage = row.getValue('brokerage') as string
      const displayText = brokerage || 'Not available'
      return h(Badge, {
        class: 'flex max-w-[180px] items-center',
        variant: 'outline',
      }, () => [
        h(Landmark, { class: 'mr-2 h-4 w-4' }),
        h('span', { class: 'truncate' }, displayText),
      ])
    },
    enableSorting: false,
  },
  {
    accessorKey: 'action',
    header: ({ column }) => h(DataTableColumnHeader<Stock>, { column, title: 'Action' }),
    cell: ({ row }) => {
      const label = row.getValue('action') as string
      const actionType = getActionType(label)
      const status = actionStatuses.find(s => s.value === actionType)
      if (!status)
        return h('div', { }, label)

      const style = {
        color: status.color,
      }

      return h(Badge, {
        class: 'flex items-center',
        style,
        variant: 'secondary',
      }, () => [
        status.icon && h(status.icon, { class: 'mr-2 h-4 w-4', style }),
        h('span', label),
      ])
    },
    filterFn: (row, id, value) => {
      const action = (row.getValue(id) as string).toLowerCase()
      return value.some((v: string) => action.includes(v))
    },
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
      const fromValue = row.original.targetFrom
      const toValue = row.original.targetTo
      const from = formatPrice(fromValue)
      const to = formatPrice(toValue)

      const priceChange = fromValue && toValue
        ? ((toValue - fromValue) / fromValue) * 100
        : null

      const priceText = `${from} → ${to}`

      if (priceChange === null) {
        return h('div', { class: 'text-sm' }, priceText)
      }

      const isPositive = priceChange > 0
      const isNegative = priceChange < 0
      const icon = isPositive ? ArrowUp : isNegative ? ArrowDown : Minus
      const colorClass = isPositive
        ? 'text-emerald-600'
        : isNegative
          ? 'text-red-600'
          : 'text-muted-foreground'

      const percentageText = `${isPositive ? '+' : ''}${priceChange.toFixed(1)}%`

      return h('div', { class: 'flex items-center gap-1.5 text-sm' }, [
        h('span', {}, priceText),
        h(icon, { class: `size-3.5 ${colorClass}` }),
        h('span', { class: `text-xs font-medium ${colorClass}` }, percentageText),
      ])
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
