import api from './api'
import type { Stock, StockFilter, StockRecommendation, PaginatedResponse } from '@/types/stock'

export const stockService = {
  async getStocks(filter: Partial<StockFilter>): Promise<PaginatedResponse<Stock>> {
    const params = new URLSearchParams()

    if (filter.page) params.append('page', String(filter.page))
    if (filter.limit) params.append('limit', String(filter.limit))
    if (filter.search) params.append('search', filter.search)
    if (filter.sortBy) params.append('sortBy', filter.sortBy)
    if (filter.sortOrder) params.append('sortOrder', filter.sortOrder)
    if (filter.action) params.append('action', filter.action)

    const response = await api.get<PaginatedResponse<Stock>>(`/stocks?${params.toString()}`)
    return response.data
  },

  async getStockById(id: string): Promise<Stock> {
    const response = await api.get<Stock>(`/stocks/${id}`)
    return response.data
  },

  async getStocksByTicker(ticker: string): Promise<Stock[]> {
    const response = await api.get<{ data: Stock[] }>(`/stocks/ticker/${ticker}`)
    return response.data.data
  },

  async getActions(): Promise<string[]> {
    const response = await api.get<{ data: string[] }>('/stocks/actions')
    return response.data.data
  },

  async getRecommendations(limit = 10): Promise<StockRecommendation[]> {
    const response = await api.get<{ data: StockRecommendation[] }>(`/recommendations?limit=${limit}`)
    return response.data.data
  },

  async getTopRecommendation(): Promise<StockRecommendation> {
    const response = await api.get<StockRecommendation>('/recommendations/top')
    return response.data
  },

  async triggerSync(): Promise<{ message: string; count: number }> {
    const response = await api.post<{ message: string; count: number }>('/sync')
    return response.data
  }
}
