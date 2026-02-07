<script setup lang="ts">
import type { Row } from '@tanstack/vue-table'
import type { Component } from 'vue'

import { Ellipsis } from 'lucide-vue-next'

import { useModal } from '@/composables/use-modal'

import type { Stock } from '../data/schema'

interface DataTableRowActionsProps {
    row: Row<Stock>
}
const props = defineProps<DataTableRowActionsProps>()
const stock = computed(() => props.row.original)
const isOpen = ref(false)

const showComponent = shallowRef<Component | null>(null)

function handleSelect() {
    showComponent.value = defineAsyncComponent(() => import('./stock-detail.vue'))
}

const { contentClass, Modal } = useModal()
</script>

<template>
    <component :is="Modal.Root" v-model:open="isOpen">
        <UiDropdownMenu>
            <UiDropdownMenuTrigger as-child>
                <UiButton
                    variant="ghost"
                    class="flex h-8 w-8 p-0 data-[state=open]:bg-muted"
                >
                    <Ellipsis class="size-4" />
                    <span class="sr-only">Open menu</span>
                </UiButton>
            </UiDropdownMenuTrigger>
            <UiDropdownMenuContent align="end" class="w-[160px]">
                <component :is="Modal.Trigger" as-child>
                    <UiDropdownMenuItem @click.stop="handleSelect()">
                        View Details
                    </UiDropdownMenuItem>
                </component>
            </UiDropdownMenuContent>
        </UiDropdownMenu>

        <component :is="Modal.Content" :class="contentClass">
            <component :is="showComponent" :stock="stock" @close="isOpen = false" />
        </component>
    </component>
</template>
