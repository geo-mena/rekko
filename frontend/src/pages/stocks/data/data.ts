import { ArrowDownRight, ArrowUpRight, Minus } from 'lucide-vue-next'
import { h } from 'vue'

import type { FacetedFilterOption } from '@/components/data-table/types'

export type ActionType = 'raised' | 'lowered' | 'default'

export const actionStatuses = [
    { value: 'raised', label: 'Target Raised', icon: h(ArrowUpRight), color: 'green' },
    { value: 'lowered', label: 'Target Lowered', icon: h(ArrowDownRight), color: 'red' },
    { value: 'default', label: 'Other', icon: h(Minus), color: 'gray' },
]

export function getActionType(action: string): ActionType {
    const lowerAction = action.toLowerCase()
    if (lowerAction.includes('target raised by'))
        return 'raised'
    if (lowerAction.includes('target lowered by'))
        return 'lowered'
    return 'default'
}

export const actionTypes: FacetedFilterOption[] = [
    { label: 'Target Raised', value: 'target raised by' },
    { label: 'Target Lowered', value: 'target lowered by' },
]
