<template>
  <div class="min-h-screen relative flex items-center justify-center p-6 bg-slate-50 dark:bg-gray-950 overflow-hidden transition-colors duration-700">
    <!-- Animated background decorators -->
    <div class="absolute top-0 right-0 -mt-24 -mr-24 w-[40rem] h-[40rem] bg-primary-500/10 dark:bg-primary-900/10 rounded-full blur-[120px] animate-pulse"></div>
    <div class="absolute bottom-0 left-0 -mb-24 -ml-24 w-[30rem] h-[30rem] bg-indigo-500/10 dark:bg-indigo-900/10 rounded-full blur-[100px] animate-pulse transition-all duration-1000"></div>

    <!-- Theme Toggle (Floating Top Right) -->
    <button @click="toggleDarkMode" 
          class="fixed top-8 right-8 z-50 p-3 rounded-2xl bg-white/20 dark:bg-gray-800/40 backdrop-blur-2xl border border-white/30 dark:border-gray-700/50 shadow-2xl hover:scale-110 active:scale-95 transition-all duration-500 group text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
          title="Toggle Dark Mode">
      <svg v-if="!isDark" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path></svg>
      <svg v-else class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"></path></svg>
      <div class="absolute inset-0 rounded-2xl bg-primary-500/10 scale-0 group-hover:scale-100 transition-transform duration-500 -z-10"></div>
    </button>

    <!-- Main Login Card -->
    <div class="relative z-10 w-full max-w-lg lg:max-w-xl animate-fade-in">
      <div class="bg-white/70 dark:bg-gray-900/70 backdrop-blur-[40px] rounded-[3rem] p-10 lg:p-14 shadow-[0_32px_120px_-20px_rgba(0,0,0,0.15)] border border-white/40 dark:border-gray-800/50 overflow-hidden group">
        <!-- Floating Aura on hover -->
        <div class="absolute -top-24 -right-24 w-48 h-48 bg-primary-500/10 rounded-full blur-3xl opacity-0 group-hover:opacity-100 transition-opacity duration-1000"></div>
        
        <!-- App Identity -->
        <div class="text-center space-y-4 mb-12">
          <div class="w-20 h-20 bg-gradient-to-br from-primary-600 to-indigo-700 rounded-[2rem] mx-auto flex items-center justify-center shadow-2xl shadow-primary-500/30 rotate-3 transition-transform duration-700 hover:rotate-0">
            <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" /></svg>
          </div>
          <div class="space-y-1">
            <h1 class="text-4xl font-black text-gray-900 dark:text-white tracking-tighter uppercase leading-none">SIS-INV</h1>
            <p class="text-[10px] font-black text-primary-600 dark:text-primary-400 uppercase tracking-[0.4em]">School Inventory System</p>
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
                <svg class="w-5 h-5 absolute left-4 top-1/2 -translate-y-1/2 text-gray-300 group-focus-within:text-primary-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>
              </div>
            </div>

            <div class="relative group">
              <label class="block text-[10px] font-black text-gray-400 uppercase tracking-widest mb-2 pl-4">Keamanan (Password)</label>
              <div class="relative">
                <input v-model="password" type="password"
                     class="w-full pl-12 pr-6 h-16 bg-white dark:bg-gray-800/50 rounded-2xl border-none ring-1 ring-gray-100 dark:ring-gray-700 focus:ring-4 focus:ring-primary-100 dark:focus:ring-primary-900/30 transition-all font-bold text-sm" 
                     placeholder="Ketik password ketat" required />
                <svg class="w-5 h-5 absolute left-4 top-1/2 -translate-y-1/2 text-gray-300 group-focus-within:text-primary-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" /></svg>
              </div>
            </div>
          </div>

          <div v-if="error" class="p-5 bg-rose-50 dark:bg-rose-900/20 border border-rose-100 dark:border-rose-900/30 text-rose-600 dark:text-rose-400 text-xs font-black rounded-2xl flex items-center gap-3 animate-shake uppercase tracking-widest">
            <svg class="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" /></svg>
            {{ error }}
          </div>

          <button :disabled="loading" 
                  class="btn-premium-primary w-full py-5 text-sm">
            <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-full group-hover:translate-x-full transition-transform duration-1000"></div>
            {{ loading ? 'SINKRONISASI...' : 'AUTENTIKASI MASUK' }}
          </button>
        </form>

        <p class="text-center text-[10px] text-gray-400 font-bold uppercase tracking-[0.3em] mt-12 pb-2 border-b border-gray-100 dark:border-gray-800">SIS-INV • AMANAH & TERTIB</p>
      </div>
      
      <!-- Footer Info -->
      <footer class="mt-8 text-center flex flex-col items-center gap-2">
        <p class="text-[10px] font-black text-gray-400 dark:text-gray-500 uppercase tracking-widest">Terlindungi oleh Enkripsi End-to-End</p>
        <div class="w-1.5 h-1.5 bg-emerald-500 rounded-full shadow-lg shadow-emerald-500/50 animate-pulse"></div>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

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
