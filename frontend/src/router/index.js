import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/LoginView.vue'),
    meta: { guest: true },
  },
  {
    path: '/',
    component: () => import('../components/AppLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('../views/DashboardView.vue'),
      },
      {
        path: 'items',
        name: 'Items',
        component: () => import('../views/ItemsView.vue'),
      },
      {
        path: 'borrow',
        name: 'Borrow',
        component: () => import('../views/BorrowView.vue'),
      },
      {
        path: 'return',
        name: 'Return',
        component: () => import('../views/ReturnView.vue'),
      },
      {
        path: 'my-borrowings',
        name: 'My Borrowings',
        component: () => import('../views/MyBorrowingsView.vue'),
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('../views/UsersView.vue'),
        meta: { roles: ['ADMIN'] },
      },
      {
        path: 'students',
        name: 'Students',
        component: () => import('../views/StudentsView.vue'),
        meta: { roles: ['ADMIN'] },
      },
      {
        path: 'categories',
        name: 'Categories',
        component: () => import('../views/CategoriesView.vue'),
        meta: { roles: ['ADMIN'] },
      },
      {
        path: 'locations',
        name: 'Locations',
        component: () => import('../views/LocationsView.vue'),
        meta: { roles: ['ADMIN'] },
      },
      {
        path: 'print-labels',
        name: 'Cetak Label',
        component: () => import('../views/PrintQRView.vue'),
        meta: { roles: ['ADMIN'] },
      },
      {
        path: 'student-history',
        name: 'Riwayat Siswa',
        component: () => import('../views/StudentHistoryView.vue'),
      },
      {
        path: 'reports',
        name: 'Laporan',
        component: () => import('../views/ReportsView.vue'),
        meta: { roles: ['ADMIN'] },
      },
      {
        path: 'maintenance',
        name: 'Perbaikan',
        component: () => import('../views/MaintenanceView.vue'),
        meta: { roles: ['ADMIN'] },
      },
      {
        path: 'audit-logs',
        name: 'Log Audit',
        component: () => import('../views/AuditLogsView.vue'),
        meta: { roles: ['ADMIN'] },
      },
      {
        path: 'stock-opname',
        name: 'Stock Opname',
        component: () => import('../views/StockOpnameView.vue'),
        meta: { roles: ['ADMIN'] },
      },
      {
        path: 'stock-opname/:id',
        name: 'Audit Detail',
        component: () => import('../views/AuditDetailView.vue'),
        meta: { roles: ['ADMIN'] },
      },
      {
        path: 'settings',
        name: 'Pengaturan',
        component: () => import('../views/SettingsView.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return next('/login')
  }

  if (to.meta.guest && authStore.isAuthenticated) {
    return next('/')
  }

  if (to.meta.roles && !to.meta.roles.includes(authStore.userRole)) {
    return next('/')
  }

  next()
})

export default router
