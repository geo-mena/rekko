<script setup lang="ts">
import type { SortField, SortOrder } from '@/types/stock'

const props = defineProps<{
  sortBy: SortField
  sortOrder: SortOrder
  action: string
  actions: string[]
}>()

const emit = defineEmits<{
  sortChange: [sortBy: SortField, sortOrder: SortOrder]
  actionChange: [action: string]
  reset: []
}>()

const sortOptions: { value: SortField; label: string }[] = [
  { value: 'createdAt', label: 'Date' },
  { value: 'ticker', label: 'Ticker' },
  { value: 'company', label: 'Company' },
  { value: 'targetTo', label: 'Target Price' }
]

function handleSortChange(event: Event) {
  const target = event.target as HTMLSelectElement
  emit('sortChange', target.value as SortField, 'desc')
}

function toggleSortOrder() {
  emit('sortChange', props.sortBy, props.sortOrder === 'asc' ? 'desc' : 'asc')
}

function handleActionChange(event: Event) {
  const target = event.target as HTMLSelectElement
  emit('actionChange', target.value)
}
</script>

<template>
  <div class="flex flex-wrap gap-4 items-center">
    <div class="flex items-center gap-2">
      <label class="text-sm font-medium text-gray-700">Sort by:</label>
      <select
        :value="sortBy"
        @change="handleSortChange"
        class="input w-auto py-1.5"
      >
        <option v-for="option in sortOptions" :key="option.value" :value="option.value">
          {{ option.label }}
        </option>
      </select>
      <button
        @click="toggleSortOrder"
        class="p-2 rounded-lg border border-gray-300 hover:bg-gray-50"
        :title="sortOrder === 'asc' ? 'Ascending' : 'Descending'"
      >
        <svg
          class="h-4 w-4 text-gray-600 transition-transform"
          :class="{ 'rotate-180': sortOrder === 'asc' }"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </button>
    </div>

    <div class="flex items-center gap-2">
      <label class="text-sm font-medium text-gray-700">Action:</label>
      <select
        :value="action"
        @change="handleActionChange"
        class="input w-auto py-1.5"
      >
        <option value="">All Actions</option>
        <option v-for="act in actions" :key="act" :value="act">
          {{ act }}
        </option>
      </select>
    </div>

    <button
      @click="emit('reset')"
      class="text-sm text-blue-600 hover:text-blue-800"
    >
      Reset filters
    </button>
  </div>
</template>
