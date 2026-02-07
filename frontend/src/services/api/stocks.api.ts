import type { AxiosError } from 'axios'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'

import { useAxios } from '@/composables/use-axios'
import type { PaginatedStocks, Stock } from '@/pages/stocks/data/schema'
import type { ApiPaginatedResponse, ApiResponse } from '@/services/types/response.type'

export interface StockFilter {
  search?: string
  sortBy?: string
  sortOrder?: string
  action?: string
  page?: number
  limit?: number
}

export function useGetStocksQuery(filter: Ref<StockFilter>) {
  const { axiosInstance } = useAxios()

  return useQuery<PaginatedStocks, AxiosError>({
    queryKey: ['stocks', filter],
    queryFn: async () => {
      const params = new URLSearchParams()
      const f = filter.value
      if (f.page) params.append('page', String(f.page))
      if (f.limit) params.append('limit', String(f.limit))
      if (f.search) params.append('search', f.search)
      if (f.sortBy) params.append('sortBy', f.sortBy)
      if (f.sortOrder) params.append('sortOrder', f.sortOrder)
      if (f.action) params.append('action', f.action)

      const response = await axiosInstance.get<ApiPaginatedResponse<Stock[]>>(`/stocks?${params.toString()}`)
      const { data, meta } = response.data
      const { pagination } = meta

      return {
        data,
        page: pagination.current_page,
        limit: pagination.per_page,
        totalCount: pagination.total_items,
        totalPages: pagination.total_pages,
        hasNext: pagination.has_next,
        hasPrev: pagination.current_page > 1,
      }
    },
  })
}

export function useGetStockByIdQuery(id: Ref<string>) {
  const { axiosInstance } = useAxios()

  return useQuery<Stock, AxiosError>({
    queryKey: ['stock', id],
    queryFn: async () => {
      const response = await axiosInstance.get<ApiResponse<Stock>>(`/stocks/${id.value}`)
      return response.data.data
    },
    enabled: () => !!id.value,
  })
}

export function useGetActionsQuery() {
  const { axiosInstance } = useAxios()

  return useQuery<string[], AxiosError>({
    queryKey: ['stockActions'],
    queryFn: async () => {
      const response = await axiosInstance.get<ApiResponse<string[]>>('/stocks/actions')
      return response.data.data
    },
  })
}

export function useSyncStocksMutation() {
  const { axiosInstance } = useAxios()
  const queryClient = useQueryClient()

  return useMutation<{ message: string, count: number }, AxiosError>({
    mutationKey: ['syncStocks'],
    mutationFn: async () => {
      const response = await axiosInstance.post<ApiResponse<{ count: number }>>('/sync')
      const { message, data } = response.data
      return { message, count: data.count }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['stocks'] })
      queryClient.invalidateQueries({ queryKey: ['stockActions'] })
      queryClient.invalidateQueries({ queryKey: ['recommendations'] })
      queryClient.invalidateQueries({ queryKey: ['topRecommendation'] })
      queryClient.invalidateQueries({ queryKey: ['dashboardStats'] })
    },
  })
}
