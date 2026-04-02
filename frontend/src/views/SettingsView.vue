<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Pengaturan Sistem</h1>
          <p class="text-primary-100/70 text-sm font-medium">Atur profil, preferensi, and keamanan akun Anda</p>
        </div>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
      
      <!-- Left Column: Navigation Tabs -->
      <div class="lg:col-span-4 space-y-4">
        <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] p-4 border border-gray-100 dark:border-gray-700 shadow-sm">
          <div class="flex flex-col gap-2">
            <button v-for="tab in [
                      {id:'profile', label:'Profil Pengguna', icon:'👤'}, 
                      {id:'security', label:'Keamanan & Password', icon:'🛡️'}, 
                      {id:'app', label:'Preferensi Aplikasi', icon:'📱'},
                      {id:'branding', label:'Kustomisasi Branding', icon:'🎨', adminOnly: true}
                    ].filter(t => !t.adminOnly || authStore.userRole === 'ADMIN')" 
                    :key="tab.id"
                    @click="activeTab = tab.id"
                    :class="activeTab === tab.id ? 'bg-primary-600 text-white shadow-lg' : 'text-gray-500 hover:bg-gray-50 dark:hover:bg-gray-700 dark:text-gray-400'"
                    class="px-6 py-4 rounded-2xl text-xs font-black transition-all flex items-center justify-between group active:scale-95">
              <span class="flex items-center gap-3">
                <span class="text-lg">{{ tab.icon }}</span>
                <span class="uppercase tracking-widest">{{ tab.label }}</span>
              </span>
              <svg class="w-4 h-4 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M9 5l7 7-7 7" /></svg>
            </button>
          </div>
        </div>

        <!-- Help Info -->
        <div class="p-8 bg-indigo-50 dark:bg-indigo-900/20 rounded-[2.5rem] border border-indigo-100 dark:border-indigo-800 space-y-4">
          <div class="w-10 h-10 bg-indigo-600 text-white rounded-xl flex items-center justify-center font-black">?</div>
          <h4 class="text-xs font-black text-indigo-900 dark:text-indigo-100 uppercase tracking-widest">Bantuan & Dukungan</h4>
          <p class="text-[10px] font-bold text-indigo-400 leading-relaxed uppercase tracking-tight">Hubungi Administrator jika Anda mengalami kesulitan dalam pengaturan akun atau memiliki pertanyaan teknis sistem.</p>
        </div>
      </div>

      <!-- Right Column: Form Area -->
      <div class="lg:col-span-8">
        <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden animate-scale-up">
          
          <!-- Profile Form -->
          <div v-if="activeTab === 'profile'" class="p-10 space-y-8 animate-fade-in">
            <div class="flex items-center gap-4">
              <div class="w-16 h-16 bg-primary-50 dark:bg-primary-900/30 text-primary-600 rounded-3xl flex items-center justify-center text-xl font-black shadow-inner">
                {{ authStore.user?.full_name?.charAt(0) }}
              </div>
              <div>
                <h3 class="text-xl font-black text-gray-900 dark:text-white leading-none capitalize">Informasi Akun</h3>
                <p class="text-[10px] font-black text-primary-500 uppercase tracking-widest mt-2">Data profil personal Anda</p>
              </div>
            </div>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 pt-6 border-t border-gray-50 dark:border-gray-700">
              <div class="space-y-1">
                <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Nama Lengkap</label>
                <input :value="authStore.user?.full_name" class="input-field rounded-2xl h-14 opacity-60" disabled />
              </div>
              <div class="space-y-1">
                <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Username (Login)</label>
                <input :value="authStore.user?.username" class="input-field rounded-2xl h-14 opacity-60" disabled />
              </div>
              <div class="space-y-1">
                <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Hak Akses (Role)</label>
                <span class="inline-block px-4 py-3 bg-gray-50 dark:bg-gray-900/40 rounded-2xl text-[10px] font-black text-primary-600 uppercase tracking-widest">{{ authStore.user?.role }}</span>
              </div>
            </div>
            
            <p class="text-[10px] text-gray-400 italic bg-gray-50 dark:bg-gray-900/30 p-4 rounded-xl">Nama and username dikunci oleh Sistem. Hubungi admin untuk perubahan biodata.</p>
          </div>

          <!-- Password / Security Form -->
          <div v-if="activeTab === 'security'" class="p-10 space-y-10 animate-fade-in">
            <h3 class="text-xl font-black text-gray-900 dark:text-white capitalize">Keamanan & Password</h3>
            
            <form @submit.prevent="handleChangePassword" class="space-y-6">
              <div class="space-y-1">
                <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Password Lama</label>
                <input type="password" v-model="passForm.old_password" class="input-field rounded-2xl h-14" required />
              </div>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-1">
                  <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Password Baru</label>
                  <input type="password" v-model="passForm.new_password" class="input-field rounded-2xl h-14" required minlength="6" />
                </div>
                <div class="space-y-1">
                  <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Konfirmasi Password Baru</label>
                  <input type="password" v-model="passForm.confirm_password" class="input-field rounded-2xl h-14" required minlength="6" />
                </div>
              </div>

              <div class="flex justify-end pt-6">
                <button type="submit" :disabled="submitting"
                        class="btn-premium-action !px-10">
                  {{ submitting ? 'MEMPROSES...' : 'GANTI PASSWORD' }}
                </button>
              </div>
            </form>
          </div>

          <!-- App Settings Form -->
          <div v-if="activeTab === 'app'" class="p-10 space-y-10 animate-fade-in">
            <h3 class="text-xl font-black text-gray-900 dark:text-white capitalize">Preferensi Aplikasi</h3>
            
            <div class="space-y-6">
              <div class="flex items-center justify-between p-6 bg-gray-50 dark:bg-gray-900/30 rounded-[2rem] border border-gray-100 dark:border-gray-700">
                <div class="space-y-1">
                  <p class="text-sm font-black text-gray-900 dark:text-white uppercase tracking-tight">Mode Gelap (Dark Mode)</p>
                  <p class="text-xs text-gray-400 font-bold uppercase tracking-widest">Atur kenyamanan mata Anda</p>
                </div>
                <button @click="toggleMode" 
                        class="w-14 h-8 bg-gray-200 dark:bg-primary-600 rounded-full relative transition-colors duration-500 shadow-inner">
                  <div class="absolute top-1 left-1 dark:left-7 w-6 h-6 bg-white rounded-full shadow-lg transition-all duration-300 flex items-center justify-center">
                    <span class="text-[10px]">{{ isDark ? '🌙' : '☀️' }}</span>
                  </div>
                </button>
              </div>
            </div>
          </div>

          <!-- Branding Customization Form (Admin Only) -->
          <div v-if="activeTab === 'branding' && authStore.userRole === 'ADMIN'" class="p-10 space-y-10 animate-fade-in text-gray-900 dark:text-white">
            <div>
              <h3 class="text-xl font-black capitalize">Kustomisasi Branding</h3>
              <p class="text-[10px] font-black text-primary-500 uppercase tracking-widest mt-2">Ubah identitas visual teks aplikasi</p>
            </div>
            
            <form @submit.prevent="handleSaveBranding" class="space-y-6">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-1">
                  <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Nama Aplikasi</label>
                  <input v-model="brandingForm.app_name" class="input-field rounded-2xl h-14" placeholder="SIS-INV" required />
                </div>
                <div class="space-y-1">
                  <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Subtitle Aplikasi</label>
                  <input v-model="brandingForm.app_subtitle" class="input-field rounded-2xl h-14" placeholder="Inventaris Sekolah" required />
                </div>
              </div>
              
              <div class="space-y-1">
                <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Deskripsi (Login Page)</label>
                <input v-model="brandingForm.app_description" class="input-field rounded-2xl h-14" placeholder="School Inventory System" required />
              </div>

              <div class="space-y-1">
                <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Teks Footer (Login Page)</label>
                <input v-model="brandingForm.app_footer" class="input-field rounded-2xl h-14" placeholder="SIS-INV • AMANAH & TERTIB" required />
              </div>

              <div class="space-y-1">
                <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Catatan Keamanan (Footer)</label>
                <input v-model="brandingForm.app_security_notice" class="input-field rounded-2xl h-14" placeholder="Terlindungi oleh Enkripsi End-to-End" required />
              </div>

              <div class="flex justify-end pt-6">
                <button type="submit" :disabled="submittingBranding"
                        class="btn-premium-action !px-10">
                  {{ submittingBranding ? 'MENYIMPAN...' : 'SIMPAN PERUBAHAN BRANDING' }}
                </button>
              </div>
            </form>
          </div>

        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import api from '../utils/api'

const authStore = useAuthStore()
const settingsStore = useSettingsStore()
const activeTab = ref('profile')
const submitting = ref(false)
const submittingBranding = ref(false)
const isDark = ref(document.documentElement.classList.contains('dark'))

const passForm = ref({ old_password: '', new_password: '', confirm_password: '' })
const brandingForm = ref({ ...settingsStore.settings })

async function handleChangePassword() {
  if (passForm.value.new_password !== passForm.value.confirm_password) {
    alert('Konfirmasi password tidak cocok!')
    return
  }
  submitting.value = true
  try {
    await api.put('/users/me/password', passForm.value)
    alert('Password berhasil diganti!')
    passForm.value = { old_password: '', new_password: '', confirm_password: '' }
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal mengganti password')
  } finally {
    submitting.value = false
  }
}

const toggleMode = () => {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark')
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

async function handleSaveBranding() {
  submittingBranding.value = true
  try {
    const success = await settingsStore.updateSettings(brandingForm.value)
    if (success) {
      alert('Konfigurasi branding berhasil diperbarui!')
    }
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal memperbarui branding')
  } finally {
    submittingBranding.value = false
  }
}
</script>
