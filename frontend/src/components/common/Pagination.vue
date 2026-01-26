<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  currentPage: number
  totalPages: number
  hasNext: boolean
  hasPrev: boolean
}>()

const emit = defineEmits<{
  pageChange: [page: number]
}>()

const visiblePages = computed(() => {
  const pages: number[] = []
  const start = Math.max(1, props.currentPage - 2)
  const end = Math.min(props.totalPages, props.currentPage + 2)

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

const firstVisiblePage = computed(() => visiblePages.value[0] ?? 1)
const lastVisiblePage = computed(() => visiblePages.value[visiblePages.value.length - 1] ?? props.totalPages)

function goToPage(page: number) {
  if (page >= 1 && page <= props.totalPages) {
    emit('pageChange', page)
  }
}
</script>

<template>
  <nav class="flex items-center justify-center gap-1">
    <button
      @click="goToPage(currentPage - 1)"
      :disabled="!hasPrev"
      class="px-3 py-2 rounded-lg border border-gray-300 text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
    >
      Previous
    </button>

    <template v-if="firstVisiblePage > 1">
      <button
        @click="goToPage(1)"
        class="px-3 py-2 rounded-lg border border-gray-300 text-sm font-medium hover:bg-gray-50"
      >
        1
      </button>
      <span v-if="firstVisiblePage > 2" class="px-2 text-gray-500">...</span>
    </template>

    <button
      v-for="page in visiblePages"
      :key="page"
      @click="goToPage(page)"
      :class="[
        'px-3 py-2 rounded-lg border text-sm font-medium',
        page === currentPage
          ? 'bg-blue-600 text-white border-blue-600'
          : 'border-gray-300 hover:bg-gray-50'
      ]"
    >
      {{ page }}
    </button>

    <template v-if="lastVisiblePage < totalPages">
      <span v-if="lastVisiblePage < totalPages - 1" class="px-2 text-gray-500">...</span>
      <button
        @click="goToPage(totalPages)"
        class="px-3 py-2 rounded-lg border border-gray-300 text-sm font-medium hover:bg-gray-50"
      >
        {{ totalPages }}
      </button>
    </template>

    <button
      @click="goToPage(currentPage + 1)"
      :disabled="!hasNext"
      class="px-3 py-2 rounded-lg border border-gray-300 text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
    >
      Next
    </button>
  </nav>
</template>
