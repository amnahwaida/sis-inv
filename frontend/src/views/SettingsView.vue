<template>
  <div class="max-w-4xl mx-auto space-y-6">
    <!-- Profile Card -->
    <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
      <div class="p-6 sm:p-8 flex flex-col sm:flex-row items-center gap-6 bg-gradient-to-br from-primary-600 to-indigo-700 text-white">
        <div class="w-24 h-24 rounded-2xl bg-white/20 backdrop-blur-md flex items-center justify-center text-3xl font-bold border border-white/30 shadow-xl">
          {{ userInitials }}
        </div>
        <div class="text-center sm:text-left">
          <h2 class="text-2xl font-black tracking-tight">{{ authStore.user?.full_name }}</h2>
          <p class="text-primary-100 font-medium opacity-90">{{ roleLabel }}</p>
          <div class="mt-3 flex flex-wrap justify-center sm:justify-start gap-2">
            <span class="px-3 py-1 bg-white/10 rounded-full text-[10px] font-bold uppercase tracking-wider border border-white/20">
              User ID: {{ authStore.user?.username }}
            </span>
            <span v-if="authStore.user?.nip" class="px-3 py-1 bg-white/10 rounded-full text-[10px] font-bold uppercase tracking-wider border border-white/20">
              NIP: {{ authStore.user?.nip }}
            </span>
          </div>
        </div>
      </div>
      
      <div class="p-6 grid grid-cols-1 md:grid-cols-2 gap-6 text-sm">
        <div class="space-y-1">
          <p class="text-xs font-bold text-gray-400 uppercase tracking-widest">Email Terdaftar</p>
          <p class="text-gray-900 font-medium">{{ authStore.user?.email || '-' }}</p>
        </div>
        <div class="space-y-1">
          <p class="text-xs font-bold text-gray-400 uppercase tracking-widest">Nomor Telepon</p>
          <p class="text-gray-900 font-medium">{{ authStore.user?.phone || '-' }}</p>
        </div>
      </div>
    </div>

    <!-- Security Card -->
    <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
      <div class="px-6 py-5 border-b border-gray-100 flex items-center justify-between bg-gray-50/30">
        <h3 class="text-lg font-bold text-gray-900 flex items-center gap-2">
          <svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
          Keamanan Akun
        </h3>
      </div>

      <div class="p-6 sm:p-8">
        <form @submit.prevent="handleChangePassword" class="max-w-md space-y-6">
          <div class="space-y-4">
            <div>
              <label class="block text-xs font-bold text-gray-500 uppercase tracking-widest mb-2">Password Saat Ini</label>
              <input 
                v-model="passwordForm.old" 
                type="password" 
                required
                class="input-field rounded-xl border-gray-200 focus:ring-primary-500"
                placeholder="••••••••"
              >
            </div>
            
            <div class="grid grid-cols-1 gap-4">
              <div>
                <label class="block text-xs font-bold text-gray-500 uppercase tracking-widest mb-2">Password Baru</label>
                <input 
                  v-model="passwordForm.new" 
                  type="password" 
                  required
                  minlength="6"
                  class="input-field rounded-xl border-gray-200 focus:ring-primary-500"
                  placeholder="Min. 6 karakter"
                >
              </div>
              <div>
                <label class="block text-xs font-bold text-gray-500 uppercase tracking-widest mb-2">Konfirmasi Password Baru</label>
                <input 
                  v-model="passwordForm.confirm" 
                  type="password" 
                  required
                  minlength="6"
                  class="input-field rounded-xl border-gray-200 focus:ring-primary-500"
                  placeholder="Ulangi password baru"
                >
              </div>
            </div>
          </div>

          <div v-if="status.msg" :class="[status.type === 'success' ? 'bg-green-50 text-green-700 border-green-100' : 'bg-red-50 text-red-700 border-red-100']" class="p-4 rounded-xl border text-xs font-bold flex items-center gap-3">
            <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path v-if="status.type === 'success'" fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              <path v-else fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
            {{ status.msg }}
          </div>

          <button 
            type="submit" 
            :disabled="loading"
            class="w-full sm:w-auto px-8 py-3 bg-primary-600 hover:bg-primary-700 text-white text-sm font-black uppercase tracking-widest rounded-xl shadow-lg shadow-primary-100 active:scale-95 transition-all disabled:opacity-50"
          >
            {{ loading ? 'Memproses...' : 'Perbarui Password' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const loading = ref(false)
const status = ref({ type: '', msg: '' })

const passwordForm = ref({
  old: '',
  new: '',
  confirm: ''
})

const userInitials = computed(() => {
  const name = authStore.user?.full_name || ''
  return name.split(' ').map(w => w[0]).join('').slice(0, 2).toUpperCase()
})

const roleLabels = { ADMIN: 'Administrator', TEACHER: 'Guru/Staff', HEAD: 'Kepala Sekolah' }
const roleLabel = computed(() => roleLabels[authStore.userRole] || authStore.userRole)

const handleChangePassword = async () => {
  if (passwordForm.value.new !== passwordForm.value.confirm) {
    status.value = { type: 'error', msg: 'Konfirmasi password baru tidak cocok' }
    return
  }

  if (passwordForm.value.old === passwordForm.value.new) {
    status.value = { type: 'error', msg: 'Password baru tidak boleh sama dengan password lama' }
    return
  }

  loading.value = true
  status.value = { type: '', msg: '' }

  try {
    await authStore.changePassword(passwordForm.value.old, passwordForm.value.new)
    status.value = { type: 'success', msg: 'Password berhasil diubah secara aman.' }
    passwordForm.value = { old: '', new: '', confirm: '' }
  } catch (err) {
    status.value = { type: 'error', msg: err.message }
  } finally {
    loading.value = false
  }
}
</script>
