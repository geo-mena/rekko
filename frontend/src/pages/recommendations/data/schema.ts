import { z } from 'zod'

import { stockSchema } from '@/pages/stocks/data/schema'

export const marketDataSchema = z.object({
    currentPrice: z.number(),
    dayChange: z.number(),
    dayChangePercent: z.number(),
    dayHigh: z.number(),
    dayLow: z.number(),
    previousClose: z.number(),
    marketCap: z.number(),
    industry: z.string(),
})
export type MarketData = z.infer<typeof marketDataSchema>

export const stockRecommendationSchema = z.object({
    stock: stockSchema,
    score: z.number(),
    reasons: z.array(z.string()),
    upsidePotential: z.number(),
    analystCount: z.number(),
    marketData: marketDataSchema.nullable().optional(),
})
export type StockRecommendation = z.infer<typeof stockRecommendationSchema>

export const stockRecommendationListSchema = z.array(stockRecommendationSchema)
