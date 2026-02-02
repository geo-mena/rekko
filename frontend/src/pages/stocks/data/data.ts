import type { FacetedFilterOption } from '@/components/data-table/types'

export const actionTypes: (FacetedFilterOption & { style: string })[] = [
  {
    label: 'Upgraded',
    value: 'upgraded',
    style: 'bg-teal-100/30 text-teal-900 dark:text-teal-200 border-teal-200',
  },
  {
    label: 'Downgraded',
    value: 'downgraded',
    style: 'bg-destructive/10 dark:bg-destructive/50 text-destructive dark:text-primary border-destructive/10',
  },
  {
    label: 'Initiated',
    value: 'initiated',
    style: 'bg-sky-200/40 text-sky-900 dark:text-sky-100 border-sky-300',
  },
  {
    label: 'Raised',
    value: 'raised',
    style: 'bg-emerald-100/30 text-emerald-900 dark:text-emerald-200 border-emerald-200',
  },
  {
    label: 'Lowered',
    value: 'lowered',
    style: 'bg-orange-100/30 text-orange-900 dark:text-orange-200 border-orange-200',
  },
  {
    label: 'Reiterated',
    value: 'reiterated',
    style: 'bg-violet-100/30 text-violet-900 dark:text-violet-200 border-violet-200',
  },
  {
    label: 'Maintained',
    value: 'maintained',
    style: 'bg-neutral-300/40 border-neutral-300',
  },
]
