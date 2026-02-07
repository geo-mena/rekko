import { z } from 'zod'

import { stockSchema } from '@/pages/stocks/data/schema'

export const stockRecommendationSchema = z.object({
    stock: stockSchema,
    score: z.number(),
    reasons: z.array(z.string()),
    upsidePotential: z.number(),
})
export type StockRecommendation = z.infer<typeof stockRecommendationSchema>

export const stockRecommendationListSchema = z.array(stockRecommendationSchema)
