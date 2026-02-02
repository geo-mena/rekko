<script setup lang="ts">
import { Loader, RefreshCw } from 'lucide-vue-next'
import { toast } from 'vue-sonner'

import { useSyncStocksMutation } from '@/services/api/stocks.api'

const { mutate, isPending } = useSyncStocksMutation()

function handleSync() {
  mutate(undefined, {
    onSuccess: (data) => {
      toast.success('Sync completed', { description: `${data.count} stocks synced` })
    },
    onError: () => {
      toast.error('Sync failed', { description: 'Could not sync stock data' })
    },
  })
}
</script>

<template>
  <UiButton variant="outline" :disabled="isPending" @click="handleSync">
    <Loader v-if="isPending" class="animate-spin" />
    <RefreshCw v-else />
    Sync Data
  </UiButton>
</template>
