import type { AxiosError } from 'axios'

import { useQuery } from '@tanstack/vue-query'

import type { StockRecommendation } from '@/pages/recommendations/data/schema'
import type { ApiResponse } from '@/services/types/response.type'

import { useAxios } from '@/composables/use-axios'

export interface RecommendationFilter {
    search?: string
    limit?: number
}

export function useGetRecommendationsQuery(filter: Ref<RecommendationFilter>) {
    const { axiosInstance } = useAxios()

    return useQuery<StockRecommendation[], AxiosError>({
        queryKey: ['recommendations', filter],
        queryFn: async () => {
            const params = new URLSearchParams()
            const f = filter.value
            if (f.limit)
                params.append('limit', String(f.limit))
            if (f.search)
                params.append('search', f.search)

            const response = await axiosInstance.get<ApiResponse<StockRecommendation[]>>(`/recommendations?${params.toString()}`)
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
