<script lang="ts" setup>
import { sidebarData } from './data/sidebar-data'
import NavTeam from './nav-team.vue'
import TeamSwitcher from './team-switcher.vue'
</script>

<template>
    <UiSidebar variant="floating" collapsible="icon" class="z-50">
        <UiSidebarHeader>
            <TeamSwitcher />
        </UiSidebarHeader>

        <UiSidebarContent>
            <NavTeam :nav-main="sidebarData.navMain" />
        </UiSidebarContent>

        <UiSidebarRail />
    </UiSidebar>
</template>
