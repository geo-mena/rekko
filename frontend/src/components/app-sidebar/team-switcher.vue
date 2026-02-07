<script lang="ts" setup>
import { useSidebar } from '@/components/ui/sidebar'

const { open } = useSidebar()
</script>

<template>
    <UiSidebarMenu>
        <UiSidebarMenuItem>
            <UiSidebarMenuButton size="lg" class="cursor-default hover:bg-transparent">
                <img src="/images/logo.svg" alt="Rekko" class="size-8 rounded-lg" />
                <div v-if="open" class="grid flex-1 text-sm leading-tight text-left">
                    <span class="font-semibold truncate text-lg">Rekko</span>
                </div>
            </UiSidebarMenuButton>
        </UiSidebarMenuItem>
    </UiSidebarMenu>
</template>
