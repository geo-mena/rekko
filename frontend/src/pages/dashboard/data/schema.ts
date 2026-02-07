import { z } from 'zod'

export const actionDistributionSchema = z.object({
    action: z.string(),
    count: z.number(),
})
export type ActionDistribution = z.infer<typeof actionDistributionSchema>

export const brokerageDistributionSchema = z.object({
    brokerage: z.string(),
    count: z.number(),
})
export type BrokerageDistribution = z.infer<typeof brokerageDistributionSchema>

export const dailyActivitySchema = z.object({
    date: z.string(),
    count: z.number(),
})
export type DailyActivity = z.infer<typeof dailyActivitySchema>

export const dashboardStatsSchema = z.object({
    totalStocks: z.number(),
    actionDistribution: z.array(actionDistributionSchema),
    brokerageDistribution: z.array(brokerageDistributionSchema),
    recentActivity: z.array(dailyActivitySchema),
})
export type DashboardStats = z.infer<typeof dashboardStatsSchema>
