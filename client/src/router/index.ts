import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { useUserStore } from '@/stores/store'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue')
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('../views/admin/AdminView.vue'),
      meta: {
        requireAuth: true
      },
      children: [
        {
          path: 'quiz',
          name: 'quiz',
          component: () => import('../views/quiz/QuizView.vue')
        },
        {
          path: 'quiz/:id',
          name: 'quizById',
          component: () => import('../views/quiz/QuizDataView.vue')
        }
      ]
    },
    {
      path: '/chat',
      name: 'chat',
      component: () => import('../views/chat/ChatView.vue')
    },

    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue')
    }
  ]
})

router.beforeEach((to, from, next) => {
  if (to.meta.requireAuth) {
    const store = useUserStore()
    if (store.isLogged) {
      // user is authenticated
      next()
    } else {
      next('/login')
    }
  } else {
    // non protected route
    next()
  }
})

export default router
