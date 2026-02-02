import type { AxiosError } from 'axios'

import { useQuery } from '@tanstack/vue-query'

import { useAxios } from '@/composables/use-axios'
import type { DashboardStats } from '@/pages/dashboard/data/schema'

export function useGetDashboardStatsQuery() {
  const { axiosInstance } = useAxios()

  return useQuery<DashboardStats, AxiosError>({
    queryKey: ['dashboardStats'],
    queryFn: async () => {
      const response = await axiosInstance.get('/dashboard/stats')
      return response.data
    },
    staleTime: 5 * 60 * 1000,
  })
}
