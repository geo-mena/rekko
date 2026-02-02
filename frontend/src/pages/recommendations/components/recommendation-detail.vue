<script setup lang="ts">
import { Check } from 'lucide-vue-next'

import { useModal } from '@/composables/use-modal'

import type { StockRecommendation } from '../data/schema'

interface Props {
  recommendation: StockRecommendation
}
const props = defineProps<Props>()
const emit = defineEmits<{ close: [] }>()

const { Modal } = useModal()

function formatPrice(value: number): string {
  return value ? `$${value.toFixed(2)}` : '-'
}
</script>

<template>
  <div>
    <component :is="Modal.Header">
      <component :is="Modal.Title">
        <span class="font-semibold">{{ recommendation.stock.ticker }}</span>
        <span class="ml-2 text-muted-foreground">{{ recommendation.stock.company }}</span>
      </component>
      <component :is="Modal.Description">
        Score: {{ recommendation.score }}/10
      </component>
    </component>

    <div class="grid grid-cols-2 gap-4 py-4">
      <div class="space-y-1">
        <p class="text-sm text-muted-foreground">Upside Potential</p>
        <p
          class="text-sm font-medium"
          :class="{
            'text-emerald-600': recommendation.upsidePotential > 0,
            'text-destructive': recommendation.upsidePotential < 0,
            'text-muted-foreground': recommendation.upsidePotential === 0,
          }"
        >
          {{ recommendation.upsidePotential > 0 ? '+' : '' }}{{ recommendation.upsidePotential.toFixed(1) }}%
        </p>
      </div>

      <div class="space-y-1">
        <p class="text-sm text-muted-foreground">Target Price</p>
        <p class="text-sm font-medium">
          {{ formatPrice(recommendation.stock.targetFrom) }} → {{ formatPrice(recommendation.stock.targetTo) }}
        </p>
      </div>

      <div class="space-y-1">
        <p class="text-sm text-muted-foreground">Rating</p>
        <p class="text-sm font-medium">
          {{ recommendation.stock.ratingFrom || '-' }} → {{ recommendation.stock.ratingTo || '-' }}
        </p>
      </div>

      <div class="space-y-1">
        <p class="text-sm text-muted-foreground">Brokerage</p>
        <p class="text-sm font-medium">{{ recommendation.stock.brokerage }}</p>
      </div>
    </div>

    <div v-if="recommendation.reasons.length" class="space-y-2 py-2">
      <p class="text-sm font-medium">Reasons</p>
      <ul class="space-y-1.5">
        <li
          v-for="(reason, index) in recommendation.reasons"
          :key="index"
          class="flex items-start gap-2 text-sm text-muted-foreground"
        >
          <Check class="size-4 mt-0.5 shrink-0 text-emerald-600" />
          <span>{{ reason }}</span>
        </li>
      </ul>
    </div>

    <component :is="Modal.Footer">
      <UiButton variant="outline" @click="emit('close')">
        Close
      </UiButton>
    </component>
  </div>
</template>
