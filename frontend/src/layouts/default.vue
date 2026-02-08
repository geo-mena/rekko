<script setup lang="ts">
import { useCookies } from '@vueuse/integrations/useCookies'
import { Braces, Loader, RefreshCw } from 'lucide-vue-next'
import { storeToRefs } from 'pinia'
import { toast } from 'vue-sonner'

import AppSidebar from '@/components/app-sidebar/index.vue'
import CommandMenuPanel from '@/components/command-menu-panel/index.vue'
import ThemePopover from '@/components/custom-theme/theme-popover.vue'
import ToggleTheme from '@/components/toggle-theme.vue'
import { SIDEBAR_COOKIE_NAME } from '@/components/ui/sidebar/utils'
import { cn } from '@/lib/utils'
import { useSyncStocksMutation } from '@/services/api/stocks.api'
import { useThemeStore } from '@/stores/theme'

const defaultOpen = useCookies([SIDEBAR_COOKIE_NAME])
const themeStore = useThemeStore()
const { contentLayout } = storeToRefs(themeStore)

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
    <UiSidebarProvider :default-open="defaultOpen.get(SIDEBAR_COOKIE_NAME)">
        <AppSidebar />
        <UiSidebarInset class="w-full max-w-full peer-data-[state=collapsed]:w-[calc(100%-var(--sidebar-width-icon)-1rem)] peer-data-[state=expanded]:w-[calc(100%-var(--sidebar-width))]">
            <header
                class="flex items-center gap-3 sm:gap-4 h-16 p-4 shrink-0 transition-[width,height] ease-linear"
            >
                <UiSidebarTrigger class="-ml-1" />
                <UiSeparator orientation="vertical" class="h-6" />
                <CommandMenuPanel />
                <div class="flex-1" />
                <div class="ml-auto flex items-center space-x-4">
                    <UiButton variant="outline" class="border-[#ffa500] text-[#ffa500] hover:bg-[#ffa500] hover:text-white dark:border-[#ffa500] dark:text-[#ffa500] dark:hover:bg-[#ffa500] dark:hover:text-slate-900" :disabled="isSyncing" @click="handleSync">
                        <Loader v-if="isSyncing" class="animate-spin" />
                        <RefreshCw v-else />
                        Sync Data
                    </UiButton>
                    <UiButton variant="outline" as="a" href="/swagger" target="_blank">
                        <Braces />
                        API
                    </UiButton>
                    <ToggleTheme />
                    <ThemePopover />
                </div>
            </header>
            <div
                :class="cn(
                    'p-4 grow',
                    contentLayout === 'centered' ? 'container mx-auto ' : '',
                )"
            >
                <router-view />
            </div>
        </UiSidebarInset>
    </UiSidebarProvider>
</template>
