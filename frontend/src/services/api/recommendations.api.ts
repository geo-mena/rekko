import type { AxiosError } from 'axios'

import { useQuery } from '@tanstack/vue-query'

import { useAxios } from '@/composables/use-axios'
import type { StockRecommendation } from '@/pages/recommendations/data/schema'
import type { ApiResponse } from '@/services/types/response.type'

export function useGetRecommendationsQuery(limit: Ref<number>) {
  const { axiosInstance } = useAxios()

  return useQuery<StockRecommendation[], AxiosError>({
    queryKey: ['recommendations', limit],
    queryFn: async () => {
      const response = await axiosInstance.get<ApiResponse<StockRecommendation[]>>(`/recommendations?limit=${limit.value}`)
      return response.data.data ?? []
    },
  })
}

export function useGetTopRecommendationQuery() {
  const { axiosInstance } = useAxios()

  return useQuery<StockRecommendation | null, AxiosError>({
    queryKey: ['topRecommendation'],
    queryFn: async () => {
      const response = await axiosInstance.get<ApiResponse<StockRecommendation>>('/recommendations/top')
      return response.data.data
    },
  })
}
