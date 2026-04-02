<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header -->
    <div class="relative overflow-hidden bg-gradient-to-br from-indigo-900 to-primary-900 rounded-[2.5rem] p-10 text-white shadow-2xl shadow-primary-900/40 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-80 h-80 bg-primary-500/10 rounded-full blur-3xl animate-pulse"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-64 h-64 bg-indigo-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-8">
        <div class="space-y-2">
          <div class="inline-flex items-center gap-2 px-3 py-1 bg-white/10 rounded-full text-[10px] font-black uppercase tracking-[0.2em] mb-2 backdrop-blur-md border border-white/10">
            <span class="w-2 h-2 bg-indigo-400 rounded-full animate-ping"></span>
            MANAJEMEN ASET
          </div>
          <h1 class="text-4xl font-black tracking-tight leading-none uppercase">Stock Opname</h1>
          <p class="text-primary-100/70 text-sm font-medium">Audit fisik vs sistem untuk memastikan akurasi data aset sekolah.</p>
        </div>
        
        <button @click="showStartModal = true" class="btn-premium-primary">
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          MULAI AUDIT BARU
        </button>
      </div>
    </div>

    <!-- Filters Row -->
    <div class="bg-white dark:bg-gray-800 p-3 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm flex flex-col md:flex-row gap-3">
      <div class="flex-1 relative group">
        <input type="text" v-model="filters.search" @input="debouncedFetch"
               class="input-field pl-12 h-14 rounded-2xl" placeholder="Cari lokasi, petugas, atau catatan audit..." />
        <svg class="w-6 h-6 absolute left-4 top-4 text-gray-300 group-focus-within:text-primary-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
        </svg>
      </div>
    </div>

    <!-- Active Sessions / History Table -->
    <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden transition-all duration-300">
      <div class="p-8 flex items-center justify-between border-b border-gray-50 dark:border-gray-700">
        <h2 class="text-xs font-black text-gray-400 uppercase tracking-[0.4em]">Riwayat & Sesi Audit</h2>
        <div class="flex items-center gap-2">
            <div v-if="auditStore.loading" class="animate-spin w-4 h-4 text-primary-600">
                <svg class="w-full h-full" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
                </svg>
            </div>
            <button @click="auditStore.fetchSessions" class="text-xs font-black text-primary-600 hover:text-primary-700 p-2 uppercase tracking-widest">Refresh</button>
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="bg-gray-50 dark:bg-gray-900/50">
              <th class="px-8 py-5 text-left text-[10px] font-black text-gray-400 uppercase tracking-widest">Status</th>
              <th class="px-8 py-5 text-left text-[10px] font-black text-gray-400 uppercase tracking-widest">Waktu Mulai</th>
              <th class="px-8 py-5 text-left text-[10px] font-black text-gray-400 uppercase tracking-widest">Lokasi</th>
              <th class="px-8 py-5 text-left text-[10px] font-black text-gray-400 uppercase tracking-widest">Petugas</th>
              <th class="px-8 py-5 text-left text-[10px] font-black text-gray-400 uppercase tracking-widest">Progres</th>
              <th class="px-8 py-5 text-right text-[10px] font-black text-gray-400 uppercase tracking-widest">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50 dark:divide-gray-700">
            <tr v-for="s in auditStore.sessions" :key="s.id" class="hover:bg-gray-50/50 dark:hover:bg-gray-900/20 transition-colors group">
              <td class="px-8 py-6">
                <span :class="s.status === 'OPEN' ? 'bg-amber-50 text-amber-600 dark:bg-amber-900/20' : 'bg-emerald-50 text-emerald-600 dark:bg-emerald-900/20'"
                      class="px-3 py-1 rounded-full text-[9px] font-black uppercase tracking-widest">
                  {{ s.status }}
                </span>
              </td>
              <td class="px-8 py-6">
                <div class="text-[11px] font-black text-gray-900 dark:text-white uppercase">{{ formatDate(s.started_at) }}</div>
                <div class="text-[9px] font-medium text-gray-400 mt-0.5">{{ formatTime(s.started_at) }}</div>
              </td>
              <td class="px-8 py-6">
                <div class="text-[11px] font-black text-gray-900 dark:text-white uppercase">{{ s.location_name }}</div>
                <div v-if="s.notes" class="text-[9px] font-medium text-gray-400 mt-1 max-w-[150px] truncate" :title="s.notes">{{ s.notes }}</div>
              </td>
              <td class="px-8 py-6">
                <div class="text-[11px] font-black text-gray-900 dark:text-white uppercase">{{ s.user_name }}</div>
              </td>
              <td class="px-8 py-6">
                <div class="w-full max-w-[100px] bg-gray-100 dark:bg-gray-700 h-1.5 rounded-full overflow-hidden">
                  <div class="h-full bg-primary-600 rounded-full" 
                       :style="`width: ${calculateProgress(s)}%`"
                       :class="{'animate-pulse bg-amber-500': s.status === 'OPEN'}"></div>
                </div>
                <div class="text-[8px] font-black text-gray-400 mt-1.5 uppercase tracking-widest">
                  {{ s.total_found }} / {{ s.total_expected }} TERDATA
                </div>
              </td>
              <td class="px-8 py-6 text-right">
                <router-link :to="`/stock-opname/${s.id}`" 
                             class="inline-flex items-center gap-2 px-4 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-primary-600 hover:text-white dark:hover:bg-primary-600 text-[10px] font-black uppercase tracking-widest rounded-xl transition-all group-hover:shadow-lg group-hover:shadow-primary-600/10 active:scale-95">
                  {{ s.status === 'OPEN' ? 'LANJUTKAN' : 'LIHAT HASIL' }}
                  <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
                  </svg>
                </router-link>
              </td>
            </tr>
            <tr v-if="auditStore.sessions.length === 0 && !auditStore.loading">
              <td colspan="6" class="px-8 py-20 text-center">
                <div class="w-16 h-16 bg-gray-50 dark:bg-gray-900 rounded-2xl flex items-center justify-center mx-auto mb-4 text-gray-300">
                  <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h3.75M9 15h3.75M9 18h3.75m3 .75H18a2.25 2.25 0 002.25-2.25V6.108c0-1.135-.845-2.098-1.976-2.192a48.424 48.424 0 00-1.123-.08m-5.801 0c-.065.21-.1.433-.1.664 0 .414.336.75.75.75h4.5a.75.75 0 00.75-.75 2.25 2.25 0 00-.1-.664m-5.8 0A2.251 2.251 0 0113.5 2.25H15c1.012 0 1.867.668 2.15 1.586m-5.8 0c-.376.023-.75.05-1.124.08C9.095 4.01 8.25 4.973 8.25 6.108V18.25m0 0c0 1.242 1.008 2.25 2.25 2.25h.75" />
                  </svg>
                </div>
                <p class="text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">Belum ada riwayat audit</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination Bar -->
      <div class="px-8 py-5 border-t border-gray-100 dark:border-gray-700 flex flex-col sm:flex-row items-center justify-between gap-4 bg-gray-50/10">
        <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">
          Menampilkan <span class="text-primary-600">{{ auditStore.sessions.length ? (auditStore.meta.page - 1) * auditStore.meta.page_size + 1 : 0 }}-{{ Math.min(auditStore.meta.page * auditStore.meta.page_size, auditStore.meta.total) }}</span> dari <span class="text-gray-900 dark:text-gray-300">{{ auditStore.meta.total }}</span> sesi
        </span>
        <div class="flex gap-2">
          <button @click="changePage(auditStore.meta.page - 1)" :disabled="auditStore.meta.page === 1" class="pagination-btn-standard">
            Kembali
          </button>
          <div class="flex gap-1">
            <button v-for="p in visiblePages" :key="p" @click="changePage(p)"
                    class="w-10 h-10 rounded-xl text-[11px] font-black transition-all shadow-sm active:scale-95 border"
                    :class="p === auditStore.meta.page ? 'bg-primary-600 text-white border-primary-600' : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:bg-primary-50 dark:hover:bg-gray-700'">
              {{ p }}
            </button>
          </div>
          <button @click="changePage(auditStore.meta.page + 1)" :disabled="auditStore.meta.page === auditStore.meta.total_pages" class="pagination-btn-standard">
            Lanjut
          </button>
        </div>
      </div>
    </div>

    <!-- START AUDIT MODAL -->
    <div v-if="showStartModal" class="fixed inset-0 z-50 flex items-center justify-center p-6 bg-gray-950/40 backdrop-blur-sm transition-all duration-500 animate-fade-in">
        <div class="bg-white dark:bg-gray-800 w-full max-w-lg rounded-[2.5rem] shadow-2xl overflow-hidden animate-scale-up border border-gray-100 dark:border-gray-700">
            <div class="bg-indigo-600 p-8 text-white relative">
                <div class="absolute top-0 right-0 -mt-10 -mr-10 w-40 h-40 bg-white/10 rounded-full blur-2xl"></div>
                <button @click="showStartModal = false" class="absolute top-6 right-6 text-white/50 hover:text-white transition-colors">
                    <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
                <h3 class="text-2xl font-black uppercase tracking-tight relative">Mulai Audit Baru</h3>
                <p class="text-indigo-100 text-[10px] font-bold uppercase tracking-widest mt-1 relative">Pilih lokasi yang akan diperiksa</p>
            </div>
            
            <div class="p-8 space-y-8">
                <div>
                    <label class="block text-[10px] font-black text-gray-400 uppercase tracking-widest mb-3 pl-1">Target Lokasi</label>
                    <div class="grid grid-cols-2 gap-3">
                        <button v-for="loc in locations" :key="loc.id"
                                @click="newAudit.location_id = loc.id"
                                class="p-4 rounded-2xl border-2 text-left transition-all relative overflow-hidden group"
                                :class="newAudit.location_id === loc.id ? 'border-primary-600 bg-primary-50 dark:bg-primary-900/20' : 'border-gray-100 dark:border-gray-700 hover:border-gray-200 dark:hover:border-gray-600'">
                            <span class="block text-[11px] font-black text-gray-900 dark:text-white uppercase leading-none group-hover:scale-105 transition-transform">{{ loc.name }}</span>
                            <div v-if="newAudit.location_id === loc.id" class="absolute bottom-1 right-2 text-primary-600">
                                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor">
                                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                </svg>
                            </div>
                        </button>
                    </div>
                </div>

                <div>
                    <label class="block text-[10px] font-black text-gray-400 uppercase tracking-widest mb-3 pl-1">Catatan Tambahan (Opsional)</label>
                    <textarea v-model="newAudit.notes" rows="3" class="w-full bg-gray-50 dark:bg-gray-900 border border-gray-100 dark:border-gray-700 rounded-2xl p-4 text-sm font-bold focus:ring-4 focus:ring-primary-500/10 focus:border-primary-500 outline-none transition-all" placeholder="Contoh: Audit Rutin Akhir Semester 1"></textarea>
                </div>

                <div class="flex gap-4 pt-4">
                    <button @click="showStartModal = false" class="flex-1 py-4 text-xs font-black text-gray-400 uppercase tracking-widest hover:bg-gray-50 dark:hover:bg-gray-900 rounded-2xl transition-all">BATAL</button>
                    <button @click="handleStartAudit" :disabled="!newAudit.location_id" class="btn-premium-action flex-[2]">
                        MULAI SEKARANG
                    </button>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuditStore } from '../stores/audit'
import api from '../utils/api'

const router = useRouter()
const auditStore = useAuditStore()
const showStartModal = ref(false)
const locations = ref([])

const newAudit = ref({ location_id: null, notes: '' })
const filters = ref({ search: '', page: 1, page_size: 10 })

let searchTimeout
const debouncedFetch = () => {
    clearTimeout(searchTimeout)
    searchTimeout = setTimeout(() => {
        filters.value.page = 1
        auditStore.fetchSessions(filters.value)
    }, 400)
}

const changePage = (p) => {
    if (p < 1 || p > auditStore.meta.total_pages) return
    filters.value.page = p
    auditStore.fetchSessions(filters.value)
}

const visiblePages = computed(() => {
    const pages = []
    const current = auditStore.meta.page
    const total = auditStore.meta.total_pages
    const start = Math.max(1, current - 2)
    const end = Math.min(total, current + 2)
    for (let i = start; i <= end; i++) pages.push(i)
    return pages
})

const formatDate = (date) => new Date(date).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })
const formatTime = (date) => new Date(date).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })

const calculateProgress = (s) => {
    if (!s.total_expected) return 0
    return Math.min(100, Math.round((s.total_found / s.total_expected) * 100))
}

const handleStartAudit = async () => {
  try {
    const session = await auditStore.startSession(newAudit.value.location_id, newAudit.value.notes)
    showStartModal.value = false
    router.push(`/stock-opname/${session.id}`)
  } catch (err) {
    alert(err)
  }
}

onMounted(async () => {
  auditStore.fetchSessions(filters.value)
  // Fetch locations for modal
  try {
    const { data } = await api.get('/locations')
    locations.value = data.data.items || data.data || []
  } catch (err) {
    console.error('Gagal ambil lokasi', err)
  }
})
</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.5s ease-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-scale-up {
  animation: scaleUp 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
@keyframes scaleUp {
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: scale(1); }
}
</style>
