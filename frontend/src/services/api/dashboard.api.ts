import type { AxiosError } from 'axios'

import { useQuery } from '@tanstack/vue-query'

import { useAxios } from '@/composables/use-axios'
import type { DashboardStats } from '@/pages/dashboard/data/schema'
import type { ApiResponse } from '@/services/types/response.type'

export function useGetDashboardStatsQuery() {
  const { axiosInstance } = useAxios()

  return useQuery<DashboardStats, AxiosError>({
    queryKey: ['dashboardStats'],
    queryFn: async () => {
      const response = await axiosInstance.get<ApiResponse<DashboardStats>>('/dashboard/stats')
      return response.data.data
    },
    staleTime: 5 * 60 * 1000,
  })
}
