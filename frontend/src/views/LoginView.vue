<template>
  <div class="min-h-screen relative flex items-center justify-center p-6 bg-slate-50 dark:bg-gray-950 overflow-hidden transition-colors duration-700">
    <!-- Animated background decorators -->
    <div class="absolute top-0 right-0 -mt-24 -mr-24 w-[40rem] h-[40rem] bg-primary-500/10 dark:bg-primary-900/10 rounded-full blur-[120px] animate-pulse"></div>
    <div class="absolute bottom-0 left-0 -mb-24 -ml-24 w-[30rem] h-[30rem] bg-indigo-500/10 dark:bg-indigo-900/10 rounded-full blur-[100px] animate-pulse transition-all duration-1000"></div>



    <!-- Main Login Card -->
    <div class="relative z-10 w-full max-w-lg lg:max-w-xl animate-fade-in">
      <div class="relative bg-white/70 dark:bg-gray-900/70 backdrop-blur-[40px] rounded-[3rem] p-10 lg:p-14 shadow-[0_32px_120px_-20px_rgba(0,0,0,0.15)] border border-white/40 dark:border-gray-800/50 overflow-hidden group">
        <!-- Theme Toggle (Inside Card) -->
        <button @click="toggleDarkMode" 
              class="absolute top-6 right-6 lg:top-8 lg:right-8 z-50 p-3 rounded-2xl bg-white/20 dark:bg-gray-800/40 backdrop-blur-2xl border border-white/30 dark:border-gray-700/50 shadow-2xl hover:scale-110 active:scale-95 transition-all duration-500 group text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
              title="Toggle Dark Mode">
          <svg v-if="!isDark" class="w-5 h-5 lg:w-6 lg:h-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z" />
          </svg>
          <svg v-else class="w-5 h-5 lg:w-6 lg:h-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z" />
          </svg>
          <div class="absolute inset-0 rounded-2xl bg-primary-500/10 scale-0 group-hover:scale-100 transition-transform duration-500 -z-10"></div>
        </button>
        <!-- Floating Aura on hover -->
        <div class="absolute -top-24 -right-24 w-48 h-48 bg-primary-500/10 rounded-full blur-3xl opacity-0 group-hover:opacity-100 transition-opacity duration-1000"></div>
        
        <!-- App Identity -->
        <div class="text-center space-y-4 mb-12">
          <div class="w-20 h-20 bg-gradient-to-br from-primary-600 to-indigo-700 rounded-[2rem] mx-auto flex items-center justify-center shadow-2xl shadow-primary-500/30 rotate-3 transition-transform duration-700 hover:rotate-0 overflow-hidden">
            <template v-if="settingsStore.settings.login_logo_url || settingsStore.settings.app_logo_url">
              <img :src="settingsStore.settings.login_logo_url || settingsStore.settings.app_logo_url" class="w-full h-full object-cover" />
            </template>
            <template v-else>
              <svg class="w-10 h-10 text-white" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.744c0 5.561 3.411 10.321 8.244 12.316a11.967 11.967 0 008.158-12.316c0-1.311-.212-2.571-.598-3.744A11.914 11.914 0 0112 2.714z" />
              </svg>
            </template>
          </div>
          <div class="space-y-1">
            <h1 class="text-4xl font-black text-gray-900 dark:text-white tracking-tighter uppercase leading-none">{{ settingsStore.settings.app_name }}</h1>
            <p class="text-[10px] font-black text-primary-600 dark:text-primary-400 uppercase tracking-[0.4em]">{{ settingsStore.settings.app_description }}</p>
          </div>
        </div>

        <!-- Login Form -->
        <form @submit.prevent="handleLogin" class="space-y-8">
          <div class="space-y-6">
            <div class="relative group">
              <label class="block text-[10px] font-black text-gray-400 uppercase tracking-widest mb-2 pl-4">Username</label>
              <div class="relative">
                <input v-model="username" 
                     class="w-full pl-12 pr-6 h-16 bg-white dark:bg-gray-800/50 rounded-2xl border-none ring-1 ring-gray-100 dark:ring-gray-700 focus:ring-4 focus:ring-primary-100 dark:focus:ring-primary-900/30 transition-all font-bold text-sm" 
                     placeholder="Nama pengguna Anda" required />
                <svg class="w-5 h-5 absolute left-4 top-1/2 -translate-y-1/2 text-gray-300 group-focus-within:text-primary-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" />
                </svg>
              </div>
            </div>

            <div class="relative group">
              <label class="block text-[10px] font-black text-gray-400 uppercase tracking-widest mb-2 pl-4">Keamanan (Password)</label>
              <div class="relative">
                <input v-model="password" type="password"
                     class="w-full pl-12 pr-6 h-16 bg-white dark:bg-gray-800/50 rounded-2xl border-none ring-1 ring-gray-100 dark:ring-gray-700 focus:ring-4 focus:ring-primary-100 dark:focus:ring-primary-900/30 transition-all font-bold text-sm" 
                     placeholder="Ketik password ketat" required />
                <svg class="w-5 h-5 absolute left-4 top-1/2 -translate-y-1/2 text-gray-300 group-focus-within:text-primary-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z" />
                </svg>
              </div>
            </div>
          </div>

          <div v-if="error" class="p-5 bg-rose-50 dark:bg-rose-900/20 border border-rose-100 dark:border-rose-900/30 text-rose-600 dark:text-rose-400 text-xs font-black rounded-2xl flex items-center gap-3 animate-shake uppercase tracking-widest">
            <svg class="w-5 h-5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
            </svg>
            {{ error }}
          </div>

          <button :disabled="loading" 
                  class="btn-premium-primary w-full py-5 text-sm">
            <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-full group-hover:translate-x-full transition-transform duration-1000"></div>
            {{ loading ? 'SINKRONISASI...' : 'AUTENTIKASI MASUK' }}
          </button>
        </form>

        <p class="text-center text-[10px] text-gray-400 font-bold uppercase tracking-[0.3em] mt-12 pb-2 border-b border-gray-100 dark:border-gray-800">{{ settingsStore.settings.app_footer }}</p>

        <!-- PWA Install Button (Android/Chrome) -->
        <button v-if="pwaStore.canInstall" @click="pwaStore.triggerInstall()" 
                class="mt-6 w-full flex items-center justify-center gap-2 px-4 py-3.5 bg-primary-600 hover:bg-primary-500 text-white rounded-2xl text-[10px] font-black uppercase tracking-widest transition-all shadow-lg shadow-primary-500/20">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3" />
          </svg>
          INSTAL APLIKASI
        </button>

        <!-- iOS Install Guide -->
        <button v-if="pwaStore.isIOS && !pwaStore.isStandalone" @click="pwaStore.showIOSGuide = true"
                class="mt-6 w-full flex items-center justify-center gap-2 px-4 py-3.5 border border-primary-500/30 text-primary-600 dark:text-primary-400 rounded-2xl text-[10px] font-black uppercase tracking-widest transition-all hover:bg-primary-50 dark:hover:bg-primary-900/10">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z" />
          </svg>
          CARA INSTAL (iOS)
        </button>
      </div>
      
      <!-- Footer Info -->
      <footer class="mt-8 text-center flex flex-col items-center gap-2">
        <p class="text-[10px] font-black text-gray-400 dark:text-gray-500 uppercase tracking-widest">{{ settingsStore.settings.app_security_notice }}</p>
        <div class="w-1.5 h-1.5 bg-emerald-500 rounded-full shadow-lg shadow-emerald-500/50 animate-pulse"></div>
      </footer>
    </div>
  </div>

  <!-- iOS Install Guide Modal -->
  <div v-if="pwaStore.showIOSGuide" class="fixed inset-0 z-[100] flex items-center justify-center p-6 bg-black/60 backdrop-blur-sm">
    <div class="bg-white dark:bg-gray-900 rounded-[2.5rem] p-8 max-w-sm w-full shadow-2xl border border-gray-100 dark:border-gray-800">
      <div class="text-center space-y-6">
        <div class="w-16 h-16 bg-primary-100 dark:bg-primary-900/30 rounded-2xl mx-auto flex items-center justify-center">
          <svg class="w-8 h-8 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
        </div>
        <div class="space-y-2">
          <h3 class="text-xl font-black text-gray-900 dark:text-white uppercase tracking-tight">Instal di iPhone</h3>
          <p class="text-sm text-gray-500 dark:text-gray-400">Ikuti langkah berikut:</p>
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
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import { usePwaStore } from '../stores/pwa'

const router = useRouter()
const authStore = useAuthStore()
const settingsStore = useSettingsStore()
const pwaStore = usePwaStore()

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const isDark = ref(localStorage.getItem('theme') === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches))

const handleLogin = async () => {
  loading.value = true
  error.value = ''
  try {
    const success = await authStore.login(username.value, password.value)
    if (success) { router.push('/') } 
    else { error.value = 'HAK AKSES DITOLAK' }
  } catch (err) {
    error.value = err.response?.data?.error || 'GAGAL TERHUBUNG KE SERVER'
  } finally {
    loading.value = false
  }
}

const toggleDarkMode = () => {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark')
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

onMounted(() => {
  if (isDark.value) document.documentElement.classList.add('dark')
  else document.documentElement.classList.remove('dark')
})
</script>

<style scoped>
.animate-shake {
  animation: shake 0.5s cubic-bezier(.36,.07,.19,.97) both;
}

@keyframes shake {
  10%, 90% { transform: translate3d(-1px, 0, 0); }
  20%, 80% { transform: translate3d(2px, 0, 0); }
  30%, 50%, 70% { transform: translate3d(-4px, 0, 0); }
  40%, 60% { transform: translate3d(4px, 0, 0); }
}
</style>
