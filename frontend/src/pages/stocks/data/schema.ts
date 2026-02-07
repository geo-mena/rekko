import { z } from 'zod'

export const stockActionSchema = z.enum([
    'upgraded',
    'downgraded',
    'initiated',
    'raised',
    'lowered',
    'reiterated',
    'maintained',
])
export type StockAction = z.infer<typeof stockActionSchema>

export const stockSchema = z.object({
    id: z.string(),
    ticker: z.string(),
    company: z.string(),
    brokerage: z.string(),
    action: z.string(),
    ratingFrom: z.string(),
    ratingTo: z.string(),
    targetFrom: z.number(),
    targetTo: z.number(),
    createdAt: z.string(),
    updatedAt: z.string(),
})
export type Stock = z.infer<typeof stockSchema>

export const stockListSchema = z.array(stockSchema)

export const paginatedStocksSchema = z.object({
    data: stockListSchema,
    page: z.number(),
    limit: z.number(),
    totalCount: z.number(),
    totalPages: z.number(),
    hasNext: z.boolean(),
    hasPrev: z.boolean(),
})
export type PaginatedStocks = z.infer<typeof paginatedStocksSchema>
