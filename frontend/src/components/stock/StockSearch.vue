<script setup lang="ts">
import { ref, watch } from 'vue'
import { useDebounce } from '@/composables/useDebounce'

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const searchInput = ref(props.modelValue)
const debouncedSearch = useDebounce(searchInput, 300)

watch(debouncedSearch, (value) => {
  emit('update:modelValue', value)
})

watch(() => props.modelValue, (value) => {
  if (value !== searchInput.value) {
    searchInput.value = value
  }
})

function clearSearch() {
  searchInput.value = ''
  emit('update:modelValue', '')
}
</script>

<template>
  <div class="relative">
    <svg
      class="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
      />
    </svg>
    <input
      v-model="searchInput"
      type="text"
      placeholder="Search by ticker or company..."
      class="input pl-10 pr-10"
    />
    <button
      v-if="searchInput"
      @click="clearSearch"
      class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
    >
      <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
        <path
          fill-rule="evenodd"
          d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
          clip-rule="evenodd"
        />
      </svg>
    </button>
  </div>
</template>
