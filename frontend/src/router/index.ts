import { createRouter, createWebHistory, type NavigationGuardNext, type RouteLocationNormalized, type RouteRecordNormalized } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('../pages/Login.vue'),
      meta: { guest: true },
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('../pages/Register.vue'),
      meta: { guest: true },
    },
    {
      path: '/',
      component: () => import('../layouts/MainLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'Dashboard',
          component: () => import('../pages/Dashboard.vue'),
        },
        {
          path: 'tokens',
          name: 'Tokens',
          component: () => import('../pages/TokenManage.vue'),
        },
        {
          path: 'profile',
          name: 'Profile',
          component: () => import('../pages/Profile.vue'),
        },
        {
          path: 'users',
          name: 'Users',
          component: () => import('../pages/UserManage.vue'),
          meta: { requiresAdmin: true },
        },
        {
          path: 'settings',
          name: 'Settings',
          component: () => import('../pages/Settings.vue'),
          meta: { requiresRoot: true },
        },
        {
          path: 'channels',
          name: 'Channels',
          component: () => import('../pages/ChannelManage.vue'),
          meta: { requiresRoot: true },
        },
        {
          path: 'logs',
          name: 'Logs',
          component: () => import('../pages/AuditLogs.vue'),
        },
      ],
    },
  ],
})

router.beforeEach(async (to: RouteLocationNormalized, _from: RouteLocationNormalized, next: NavigationGuardNext) => {
  const userStore = useUserStore()

  // 检查路由是否需要认证（包括父路由）
  const requiresAuth = to.matched.some((record: RouteRecordNormalized) => record.meta.requiresAuth)
  const requiresAdmin = to.matched.some((record: RouteRecordNormalized) => record.meta.requiresAdmin)
  const requiresRoot = to.matched.some((record: RouteRecordNormalized) => record.meta.requiresRoot)
  const isGuestPage = to.meta.guest

  // 如果未登录且不是游客页面，尝试获取用户信息
  if (!userStore.isLoggedIn && !isGuestPage) {
    await userStore.fetchUser()
  }

  // 需要认证的页面
  if (requiresAuth && !userStore.isLoggedIn) {
    return next('/login')
  }

  // 需要管理员权限
  if (requiresAdmin && !userStore.isAdmin) {
    return next('/')
  }

  // 需要Root权限
  if (requiresRoot && !userStore.isRoot) {
    return next('/')
  }

  // 已登录用户访问登录/注册页面
  if (isGuestPage && userStore.isLoggedIn) {
    return next('/')
  }

  next()
})

export default router
