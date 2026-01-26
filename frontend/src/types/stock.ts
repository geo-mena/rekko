export interface Stock {
  id: string
  ticker: string
  company: string
  brokerage: string
  action: string
  ratingFrom: string
  ratingTo: string
  targetFrom: number
  targetTo: number
  createdAt: string
  updatedAt: string
}

export interface StockFilter {
  search: string
  sortBy: SortField
  sortOrder: SortOrder
  action: string
  page: number
  limit: number
}

export type SortField = 'ticker' | 'company' | 'action' | 'targetTo' | 'createdAt'
export type SortOrder = 'asc' | 'desc'

export interface StockRecommendation {
  stock: Stock
  score: number
  reasons: string[]
  upsidePotential: number
}

export interface PaginatedResponse<T> {
  data: T[]
  page: number
  limit: number
  totalCount: number
  totalPages: number
  hasNext: boolean
  hasPrev: boolean
}

export const defaultFilter: StockFilter = {
  search: '',
  sortBy: 'createdAt',
  sortOrder: 'desc',
  action: '',
  page: 1,
  limit: 20
}
