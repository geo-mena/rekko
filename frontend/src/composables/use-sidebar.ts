import { BadgeHelp, BellDot, Boxes, Bug, Component, CreditCard, LayoutDashboard, Lightbulb, ListTodo, Palette, PictureInPicture2, Podcast, Settings, SquareUserRound, TrendingUp, User, Users, Wrench } from 'lucide-vue-next'

import type { NavGroup } from '@/components/app-sidebar/types'

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

  return {
    navData,
    otherPages,
  }
}
