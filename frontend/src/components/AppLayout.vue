<template>
  <div class="flex h-screen bg-gray-50">
    <!-- Sidebar -->
    <aside 
      :class="[sidebarOpen ? 'translate-x-0' : '-translate-x-full', 'lg:translate-x-0']"
      class="fixed inset-y-0 left-0 z-30 w-64 bg-gray-900 transition-transform duration-300 ease-in-out lg:static lg:inset-0"
    >
      <!-- Logo -->
      <div class="flex items-center gap-3 px-6 py-5 border-b border-gray-800">
        <div class="w-9 h-9 bg-primary-500 rounded-lg flex items-center justify-center">
          <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
        </div>
        <div>
          <h1 class="text-white font-bold text-lg leading-tight">SIS-INV</h1>
          <p class="text-gray-500 text-xs">Inventaris Sekolah</p>
        </div>
      </div>

      <!-- Nav -->
      <nav class="mt-4 px-3 space-y-1">
        <router-link 
          v-for="item in filteredNav"
          :key="item.to"
          :to="item.to" 
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors"
          :class="$route.path === item.to ? 'bg-primary-600 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white'"
          @click="sidebarOpen = false"
        >
          <component :is="item.icon" class="w-5 h-5" />
          <span>{{ item.label }}</span>
        </router-link>
      </nav>

      <!-- User info at bottom -->
      <div class="absolute bottom-0 left-0 right-0 p-4 border-t border-gray-800">
        <div class="flex items-center gap-3">
          <div class="w-9 h-9 rounded-full bg-primary-600 flex items-center justify-center text-white font-semibold text-sm">
            {{ userInitials }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-white truncate">{{ authStore.user?.full_name }}</p>
            <p class="text-xs text-gray-500">{{ roleBadge }}</p>
          </div>
          <button @click="handleLogout" class="p-1.5 text-gray-500 hover:text-red-400 transition-colors" title="Logout">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
          </button>
        </div>
      </div>
    </aside>

    <!-- Overlay -->
    <div 
      v-if="sidebarOpen" 
      class="fixed inset-0 bg-black/50 z-20 lg:hidden"
      @click="sidebarOpen = false"
    ></div>

    <!-- Main -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Top bar -->
      <header class="bg-white border-b border-gray-200 px-4 py-3 flex items-center gap-4 lg:px-6">
        <button @click="sidebarOpen = !sidebarOpen" class="lg:hidden p-1.5 text-gray-500 hover:text-gray-700 rounded-lg hover:bg-gray-100">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
        <h2 class="text-lg font-semibold text-gray-800">{{ $route.name }}</h2>
      </header>

      <!-- Page content -->
      <main class="flex-1 overflow-y-auto p-4 lg:p-6">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, h } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const sidebarOpen = ref(false)

const userInitials = computed(() => {
  const name = authStore.user?.full_name || ''
  return name.split(' ').map(w => w[0]).join('').slice(0, 2).toUpperCase()
})

const roleLabels = { ADMIN: 'Administrator', TEACHER: 'Guru/Staff', HEAD: 'Kepala Sekolah' }
const roleBadge = computed(() => roleLabels[authStore.userRole] || authStore.userRole)

const IconDashboard = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z' })]) }
const IconItems = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4' })]) }
const IconUsers = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z' })]) }
const IconBorrow = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 4v16m8-8H4' })]) }
const IconReturn = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15' })]) }
const IconMyBorrowings = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2' })]) }

const navItems = [
  { to: '/', label: 'Dashboard', icon: IconDashboard, roles: ['ADMIN', 'TEACHER', 'HEAD'] },
  { to: '/items', label: 'Barang', icon: IconItems, roles: ['ADMIN', 'TEACHER', 'HEAD'] },
  { to: '/borrow', label: 'Pinjam', icon: IconBorrow, roles: ['ADMIN', 'TEACHER', 'HEAD'] },
  { to: '/return', label: 'Pengembalian', icon: IconReturn, roles: ['ADMIN', 'TEACHER', 'HEAD'] },
  { to: '/my-borrowings', label: 'Peminjaman Saya', icon: IconMyBorrowings, roles: ['ADMIN', 'TEACHER', 'HEAD'] },
  { to: '/users', label: 'Kelola User', icon: IconUsers, roles: ['ADMIN'] },
]

const filteredNav = computed(() =>
  navItems.filter(item => item.roles.includes(authStore.userRole))
)

async function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
