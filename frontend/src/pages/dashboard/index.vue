<script setup lang="ts">
import { BasicPage } from '@/components/global-layout'

import OverviewContent from './components/overview-content.vue'
</script>

<template>
  <BasicPage title="Dashboard" description="Stock market overview and recommendations" sticky>
    <OverviewContent />
  </BasicPage>
</template>
