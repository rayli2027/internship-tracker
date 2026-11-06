import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from './auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: () => import('./views/LoginView.vue'), meta: { guest: true } },
    { path: '/register', component: () => import('./views/RegisterView.vue'), meta: { guest: true } },
    { path: '/', component: () => import('./views/BoardView.vue'), meta: { auth: true } },
    { path: '/timeline', component: () => import('./views/TimelineView.vue'), meta: { auth: true } },
    { path: '/calendar', component: () => import('./views/CalendarView.vue'), meta: { auth: true } },
    { path: '/resume', component: () => import('./views/ResumeView.vue'), meta: { auth: true } },
    { path: '/settings', component: () => import('./views/SettingsView.vue'), meta: { auth: true } },
    {
      path: '/admin',
      component: () => import('./views/admin/AdminLayout.vue'),
      meta: { auth: true, admin: true },
      children: [
        { path: '', redirect: '/admin/overview' },
        { path: 'overview', component: () => import('./views/admin/OverviewView.vue') },
        { path: 'users', component: () => import('./views/admin/UsersView.vue') },
        { path: 'applications', component: () => import('./views/admin/ApplicationsView.vue') },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const { isLoggedIn, isAdmin } = useAuth()

  if (to.meta.auth && !isLoggedIn.value) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  if (to.meta.guest && isLoggedIn.value) {
    return { path: '/' }
  }
  if (to.meta.admin && !isAdmin.value) {
    return { path: '/' }
  }
  return true
})

export default router
