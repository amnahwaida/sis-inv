<template>
  <div class="flex h-screen bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100 transition-colors duration-300">
    <!-- Sidebar -->
    <aside 
      :class="[sidebarOpen ? 'translate-x-0' : '-translate-x-full', 'lg:translate-x-0', 'no-print']"
      class="fixed inset-y-0 left-0 z-30 w-64 bg-gray-900 transition-transform duration-300 ease-in-out lg:static lg:inset-0"
    >
      <!-- Logo -->
      <div class="p-6 flex items-center gap-4">
        <div class="w-10 h-10 bg-primary-600 rounded-xl flex items-center justify-center text-white shadow-lg shadow-primary-900/20 overflow-hidden">
          <img v-if="settingsStore.settings.app_logo_url" :src="settingsStore.settings.app_logo_url" class="w-full h-full object-cover" />
          <svg v-else class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-5.25v9" />
          </svg>
        </div>
        <div class="flex-1 min-w-0">
          <h2 class="text-lg font-black text-white tracking-tighter uppercase truncate leading-none mb-1">{{ settingsStore.settings.app_name }}</h2>
          <p class="text-[10px] font-bold text-gray-500 uppercase tracking-widest truncate">{{ settingsStore.settings.app_subtitle }}</p>
        </div>
      </div>

      <!-- Install Button Section (Android/Chrome) -->
      <div v-if="pwaStore.canInstall" class="px-4 py-3 border-b border-gray-800">
        <button 
          @click="pwaStore.triggerInstall()"
          class="w-full flex items-center justify-center gap-2 px-4 py-2.5 bg-primary-600 hover:bg-primary-500 text-white rounded-xl text-xs font-black transition-all shadow-lg shadow-primary-600/20"
        >
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3" />
          </svg>
          INSTAL APLIKASI
        </button>
      </div>

      <!-- iOS Install Instruction -->
      <div v-if="pwaStore.isIOS && !pwaStore.isStandalone" class="px-4 py-3 border-b border-gray-800">
        <button 
          @click="pwaStore.showIOSGuide = true"
          class="w-full flex items-center justify-center gap-2 px-4 py-2.5 border border-primary-500/30 text-primary-400 hover:text-primary-300 rounded-xl text-[10px] font-black transition-all"
        >
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z" />
          </svg>
          CARA INSTAL (iOS)
        </button>
      </div>

      <!-- Nav -->
      <nav class="mt-4 px-3 pb-20 overflow-y-auto max-h-[calc(100vh-250px)] space-y-6 scrollbar-hide">
        <div v-for="(group, groupName) in groupedNav" :key="groupName" class="space-y-1">
          <!-- Group Header (Collapsible) -->
          <button 
            v-if="groupName !== 'Dashboard'"
            @click="toggleGroup(groupName)"
            class="w-full flex items-center justify-between px-3 py-2.5 my-1 rounded-xl text-[10px] font-black text-gray-500 hover:text-white uppercase tracking-[0.2em] transition-all hover:bg-gray-800/50 group"
          >
            <div class="flex items-center gap-3">
              <component :is="groupIcons[groupName]" class="w-4 h-4 opacity-50 group-hover:opacity-100 transition-opacity" />
              <span>{{ groupName }}</span>
            </div>
            <svg 
              class="w-3.5 h-3.5 transition-transform duration-300 transform-gpu" 
              :class="openGroups[groupName] ? 'rotate-180' : ''"
              fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor"
            >
              <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
            </svg>
          </button>

          <!-- Group Items -->
          <transition 
            name="expand"
          >
            <div v-show="groupName === 'Dashboard' || openGroups[groupName]" class="overflow-hidden space-y-1">
              <router-link 
                v-for="item in group"
                :key="item.to"
                :to="item.to" 
                class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all group/item active:scale-95 ml-1"
                :class="$route.path === item.to ? 'bg-primary-600 text-white shadow-lg shadow-primary-900/20' : 'text-gray-400 hover:bg-gray-800 hover:text-white'"
                @click="sidebarOpen = false"
              >
                <template v-if="settingsStore.settings[item.iconKey]">
                  <div 
                    v-if="settingsStore.settings[item.iconKey].trim().startsWith('<svg')" 
                    v-html="settingsStore.settings[item.iconKey]" 
                    class="w-5 h-5 flex items-center justify-center flex-shrink-0 opacity-80 group-hover/item:opacity-100 [&>svg]:w-full [&>svg]:h-full transition-opacity"
                  />
                  <span v-else class="w-5 h-5 flex items-center justify-center text-lg leading-none flex-shrink-0">
                    {{ settingsStore.settings[item.iconKey] }}
                  </span>
                </template>
                <component v-else :is="item.icon" class="w-5 h-5 flex-shrink-0 opacity-80 group-hover/item:opacity-100 transition-opacity" />
                <span class="truncate whitespace-nowrap">{{ item.label }}</span>
              </router-link>
            </div>
          </transition>
        </div>
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
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" />
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
      <header class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-4 py-3 flex items-center justify-between gap-4 lg:px-6 no-print transition-colors duration-300">
        <div class="flex items-center gap-4">
          <button @click="sidebarOpen = !sidebarOpen" class="lg:hidden p-1.5 text-gray-500 hover:text-gray-700 dark:hover:text-gray-300 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800">
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
            </svg>
          </button>
          <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ $route.name }}</h2>
        </div>
        
        <!-- Right Header Items -->
        <div class="flex items-center gap-3">
          <!-- Dark Mode Toggle -->
          <button @click="toggleDarkMode" class="p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 rounded-xl hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors" title="Toggle Dark Mode">
            <svg v-if="!isDark" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z" />
            </svg>
            <svg v-else class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z" />
            </svg>
          </button>
          
          <!-- Offline Indicator -->
          <div v-if="isOffline" class="flex items-center gap-2 px-3 py-1.5 rounded-full bg-red-50 border border-red-100 text-red-600 text-[10px] font-black uppercase tracking-widest animate-pulse">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636a9 9 0 00-12.728 0M12 18v.01M15.536 8.464a5 5 0 00-7.072 0M12 14v-2" />
            </svg>
            Offline
          </div>
        </div>
      </header>

      <!-- Page content -->
      <main class="flex-1 overflow-y-auto p-4 lg:p-6 print:p-0">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <!-- iOS Install Guide Modal -->
    <div v-if="pwaStore.showIOSGuide" class="fixed inset-0 z-[100] flex items-center justify-center p-6 bg-black/60 backdrop-blur-sm">
      <div class="bg-white dark:bg-gray-900 rounded-[2.5rem] p-8 max-w-sm w-full shadow-2xl border border-gray-100 dark:border-gray-800 animate-fade-in">
        <div class="text-center space-y-6">
          <div class="w-16 h-16 bg-primary-100 dark:bg-primary-900/30 rounded-2xl mx-auto flex items-center justify-center">
            <svg class="w-8 h-8 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M15.536 8.464a5 5 0 00-7.072 0M12 14v-2" /></svg>
          </div>
          <div class="space-y-2">
            <h3 class="text-xl font-black text-gray-900 dark:text-white uppercase tracking-tight">Instal di iPhone</h3>
            <p class="text-sm text-gray-500 dark:text-gray-400">Ikuti langkah berikut untuk memasang aplikasi di layar utama Anda:</p>
          </div>
          <div class="bg-gray-50 dark:bg-gray-800/50 rounded-2xl p-6 text-left space-y-4">
            <div class="flex items-start gap-3">
              <div class="w-6 h-6 bg-white dark:bg-gray-700 rounded-lg flex items-center justify-center text-[10px] font-black border border-gray-200 dark:border-gray-600 flex-shrink-0">1</div>
              <p class="text-xs text-gray-600 dark:text-gray-300">Ketuk ikon <span class="font-bold text-primary-600">Share</span> (kotak dengan panah ke atas) di bagian bawah Safari.</p>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-6 h-6 bg-white dark:bg-gray-700 rounded-lg flex items-center justify-center text-[10px] font-black border border-gray-200 dark:border-gray-600 flex-shrink-0">2</div>
              <p class="text-xs text-gray-600 dark:text-gray-300">Gulir ke bawah dan ketuk <span class="font-bold text-primary-600">"Add to Home Screen"</span>.</p>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-6 h-6 bg-white dark:bg-gray-700 rounded-lg flex items-center justify-center text-[10px] font-black border border-gray-200 dark:border-gray-600 flex-shrink-0">3</div>
              <p class="text-xs text-gray-600 dark:text-gray-300">Ketuk <span class="font-bold text-primary-600">"Add"</span> di pojok kanan atas.</p>
            </div>
          </div>
          <button @click="pwaStore.showIOSGuide = false" class="w-full py-4 bg-gray-900 dark:bg-white text-white dark:text-gray-900 rounded-xl text-xs font-black tracking-widest uppercase transition-transform active:scale-95">MENGERTI</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, h, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import { usePwaStore } from '../stores/pwa'

const router = useRouter()
const authStore = useAuthStore()
const settingsStore = useSettingsStore()
const pwaStore = usePwaStore()
const sidebarOpen = ref(false)
const isOffline = ref(!navigator.onLine)

const userThemeKey = computed(() => `theme_${authStore.user?.username || 'default'}`)

const isDark = ref(
  localStorage.getItem(userThemeKey.value) === 'dark' || 
  (!localStorage.getItem(userThemeKey.value) && localStorage.getItem('theme') === 'dark') || 
  (!localStorage.getItem(userThemeKey.value) && !('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)
)

const toggleDarkMode = () => {
  isDark.value = !isDark.value
  if (isDark.value) {
    document.documentElement.classList.add('dark')
    localStorage.setItem(userThemeKey.value, 'dark')
  } else {
    document.documentElement.classList.remove('dark')
    localStorage.setItem(userThemeKey.value, 'light')
  }
}

const updateOnlineStatus = () => {
  isOffline.value = !navigator.onLine
}

onMounted(() => {
  window.addEventListener('online', updateOnlineStatus)
  window.addEventListener('offline', updateOnlineStatus)
  
  if (isDark.value) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
})

onUnmounted(() => {
  window.removeEventListener('online', updateOnlineStatus)
  window.removeEventListener('offline', updateOnlineStatus)
})

const userInitials = computed(() => {
  const name = authStore.user?.full_name || ''
  return name.split(' ').map(w => w[0]).join('').slice(0, 2).toUpperCase()
})

const roleLabels = { ADMIN: 'Administrator', TEACHER: 'Guru/Staff' }
const roleBadge = computed(() => roleLabels[authStore.userRole] || authStore.userRole)

const IconDashboard = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z' })]) }
const IconItems = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4' })]) }
const IconUsers = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z' })]) }
const IconCategories = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z' })]) }
const IconLocations = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z' }), h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M15 11a3 3 0 11-6 0 3 3 0 016 0z' })]) }
const IconBorrow = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 4v16m8-8H4' })]) }
const IconReturn = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15' })]) }
const IconMyBorrowings = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2' })]) }
const IconReports = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z' })]) }
const IconStudent = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 14l9-5-9-5-9 5 9 5z' }), h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 14l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14z' }), h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 14l9-5-9-5-9 5 9 5zm0 0l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14zm-4 6v-7.5l4-2.222' })]) }
const IconSettings = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z' }), h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M15 12a3 3 0 11-6 0 3 3 0 016 0z' })]) }
const IconAuditLog = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z' })]) }
const IconStockOpname = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01' })]) }
const IconMaintenance = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-5 h-5' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z' }), h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M15 12a3 3 0 11-6 0 3 3 0 016 0z' })]) }

const groupIcons = {
  'Transaksi': { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-4 h-4' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5' })]) },
  'Inventaris': { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-4 h-4' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4' })]) },
  'Pemeliharaan': { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-4 h-4' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M11.42 15.17L17.25 21A2.67 2.67 0 1113.5 17.25l-5.83-5.83m.57-2.57l5.41-5.41a2.67 2.67 0 113.78 3.78l-5.41 5.41M11.42 15.17l5.83 5.83M11.42 15.17l-5.83-5.83m-1.5 1.5l1.5-1.5m0 0l1.5 1.5m-1.5-1.5l1.5-1.5' })]) },
  'Admin': { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-4 h-4' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z' })]) },
  'Laporan': { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-4 h-4' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z' })]) },
  'Sistem': { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-4 h-4' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M4.5 12a7.5 7.5 0 0015 0m-15 0a7.5 7.5 0 1115 0m-15 0H3m16.5 0H21m-1.5 0H12m-8.457 3.077l1.41-.513m0 0l-1.41-1.413m1.41 1.413L5.636 5.636m11.144 11.144l-1.413 1.413m0 0l1.413-1.413m-1.413 1.413L18.364 18.364' })]) }
}

const navItems = [
  { to: '/', label: 'Dashboard', icon: IconDashboard, iconKey: 'icon_dashboard', roles: ['ADMIN', 'TEACHER'], group: 'Dashboard' },
  
  { to: '/borrow', label: 'Pinjam', icon: IconBorrow, iconKey: 'icon_borrow', roles: ['ADMIN', 'TEACHER'], group: 'Transaksi' },
  { to: '/return', label: 'Kembali', icon: IconReturn, iconKey: 'icon_return', roles: ['ADMIN', 'TEACHER'], group: 'Transaksi' },
  { to: '/my-borrowings', label: 'Pinjaman Saya', icon: IconMyBorrowings, iconKey: 'icon_my_borrowings', roles: ['ADMIN', 'TEACHER'], group: 'Transaksi' },
  
  { to: '/items', label: 'Barang', icon: IconItems, iconKey: 'icon_items', roles: ['ADMIN', 'TEACHER'], group: 'Inventaris' },
  { to: '/categories', label: 'Kategori', icon: IconCategories, iconKey: 'icon_categories', roles: ['ADMIN'], group: 'Inventaris' },
  { to: '/locations', label: 'Lokasi', icon: IconLocations, iconKey: 'icon_locations', roles: ['ADMIN'], group: 'Inventaris' },
  
  { to: '/stock-opname', label: 'Opname', icon: IconStockOpname, iconKey: 'icon_stock_opname', roles: ['ADMIN'], group: 'Pemeliharaan' },
  { to: '/maintenance', label: 'Perbaikan', icon: IconMaintenance, iconKey: 'icon_maintenance', roles: ['ADMIN'], group: 'Pemeliharaan' },
  
  { to: '/users', label: 'Data User', icon: IconUsers, iconKey: 'icon_users', roles: ['ADMIN'], group: 'Admin' },
  { to: '/students', label: 'Data Siswa', icon: IconStudent, iconKey: 'icon_students', roles: ['ADMIN'], group: 'Admin' },
  { to: '/student-history', label: 'Riwayat Siswa', icon: IconStudent, iconKey: 'icon_student_history', roles: ['ADMIN', 'TEACHER'], group: 'Admin' },
  
  { to: '/reports', label: 'Laporan', icon: IconReports, iconKey: 'icon_reports', roles: ['ADMIN'], group: 'Laporan' },
  { to: '/audit-logs', label: 'Audit Log', icon: IconAuditLog, iconKey: 'icon_audit_logs', roles: ['ADMIN'], group: 'Laporan' },
  
  { to: '/settings', label: 'Pengaturan', icon: IconSettings, iconKey: 'icon_settings', roles: ['ADMIN', 'TEACHER'], group: 'Sistem' },
]

const openGroups = ref(JSON.parse(localStorage.getItem('openGroups') || '{"Transaksi": true, "Inventaris": true}'))

const toggleGroup = (groupName) => {
  openGroups.value[groupName] = !openGroups.value[groupName]
  localStorage.setItem('openGroups', JSON.stringify(openGroups.value))
}

const groupedNav = computed(() => {
  const groups = {}
  navItems.forEach(item => {
    if (item.roles.includes(authStore.userRole)) {
      if (!groups[item.group]) groups[item.group] = []
      groups[item.group].push(item)
    }
  })
  return groups
})

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
.expand-enter-active, .expand-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  max-height: 500px;
}
.expand-enter-from, .expand-leave-to {
  opacity: 0;
  max-height: 0;
  transform: translateY(-8px) scale(0.98);
}
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
@media print {
  .no-print {
    display: none !important;
  }
}
</style>
