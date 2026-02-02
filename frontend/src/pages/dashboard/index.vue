<script setup lang="ts">
import { Loader, RefreshCw } from 'lucide-vue-next'
import { toast } from 'vue-sonner'

import { BasicPage } from '@/components/global-layout'
import { useSyncStocksMutation } from '@/services/api/stocks.api'

import OverviewContent from './components/overview-content.vue'

const { mutate: syncStocks, isPending: isSyncing } = useSyncStocksMutation()

function handleSync() {
  syncStocks(undefined, {
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
  <BasicPage title="Dashboard" description="Stock market overview and recommendations" sticky>
    <template #actions>
      <UiButton variant="outline" :disabled="isSyncing" @click="handleSync">
        <Loader v-if="isSyncing" class="animate-spin" />
        <RefreshCw v-else />
        Sync Data
      </UiButton>
    </template>
    <OverviewContent />
  </BasicPage>
</template>
