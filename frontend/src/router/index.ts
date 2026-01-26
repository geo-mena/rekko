import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue')
    },
    {
      path: '/stock/:id',
      name: 'stock-detail',
      component: () => import('@/views/StockDetailView.vue')
    },
    {
      path: '/recommendations',
      name: 'recommendations',
      component: () => import('@/views/RecommendationView.vue')
    }
  ]
})

export default router
