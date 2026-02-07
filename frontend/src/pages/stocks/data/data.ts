import type { FacetedFilterOption } from '@/components/data-table/types'

export type ActionType = 'raised' | 'lowered' | 'default'

export const actionStyles: Record<ActionType, string> = {
  raised: 'bg-cyan-100 text-cyan-800 dark:bg-cyan-900/30 dark:text-cyan-300 border-transparent',
  lowered: 'bg-fuchsia-100 text-fuchsia-800 dark:bg-fuchsia-900/30 dark:text-fuchsia-300 border-transparent',
  default: 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300 border-transparent',
}

export function getActionType(action: string): ActionType {
  const lowerAction = action.toLowerCase()
  if (lowerAction.includes('target raised by')) return 'raised'
  if (lowerAction.includes('target lowered by')) return 'lowered'
  return 'default'
}

export const actionTypes: FacetedFilterOption[] = [
  { label: 'Target Raised', value: 'target raised by' },
  { label: 'Target Lowered', value: 'target lowered by' },
]
