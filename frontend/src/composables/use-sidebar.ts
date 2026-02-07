import { CreditCard, LayoutDashboard, Lightbulb, PictureInPicture2, Settings, User } from 'lucide-vue-next'

import type { NavGroup } from '@/components/app-sidebar/types'
import type { TwoColAsideNavItem } from '@/components/global-layout/types'

export function useSidebar() {
    const navData = ref<NavGroup[]> ([
        {
            title: 'General',
            items: [
                { title: 'Dashboard', url: '/dashboard', icon: LayoutDashboard },
                { title: 'Stocks', url: '/stocks', icon: PictureInPicture2 },
                { title: 'Recommendations', url: '/recommendations', icon: Lightbulb },
            ],
        },
    ])

    const otherPages = ref<NavGroup[]>([
        {
            title: 'Other',
            items: [
                {
                    title: 'Plans & Pricing',
                    icon: CreditCard,
                    url: '/billing',
                },
            ],
        },
    ])

    const settingsNavItems: TwoColAsideNavItem[] = [
        { title: 'Profile', url: '/settings', icon: User },
        { title: 'Account', url: '/settings/account', icon: Settings },
    ]

    return {
        navData,
        otherPages,
        settingsNavItems,
    }
}
