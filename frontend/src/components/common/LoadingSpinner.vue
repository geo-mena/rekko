<script setup lang="ts">
defineProps<{
  size?: 'sm' | 'md' | 'lg'
}>()
</script>

<template>
  <div class="flex justify-center items-center">
    <svg
      :class="[
        'animate-spin text-blue-600',
        size === 'sm' ? 'h-5 w-5' : size === 'lg' ? 'h-12 w-12' : 'h-8 w-8'
      ]"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle
        class="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        stroke-width="4"
      />
      <path
        class="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
      />
    </svg>
  </div>
</template>
