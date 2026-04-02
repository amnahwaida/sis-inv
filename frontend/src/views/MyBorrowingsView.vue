<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-primary-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Peminjaman Saya</h1>
          <p class="text-primary-100/70 text-sm font-medium">Pantau status barang yang sedang Anda gunakan</p>
        </div>
      </div>
    </div>

    <!-- Filters Row -->
    <div class="bg-white dark:bg-gray-800 p-3 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm flex flex-col md:flex-row gap-3">
      <div class="flex-1 relative group">
        <input type="text" v-model="searchQuery" @input="debouncedSearch"
               class="input-field pl-12 h-14 rounded-2xl" placeholder="Cari nama barang atau kode aset dalam peminjaman Anda..." />
        <svg class="w-6 h-6 absolute left-4 top-4 text-gray-300 group-focus-within:text-primary-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
        </svg>
      </div>
    </div>

    <!-- Active Borrowings Table Card -->
    <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden transition-all duration-300">
      <div class="p-8 border-b border-gray-100 dark:border-gray-700 flex items-center justify-between">
        <h2 class="text-lg font-black text-gray-900 dark:text-white uppercase tracking-tight">Barang Belum Kembali</h2>
        <span class="px-4 py-1 bg-primary-50 dark:bg-primary-900/40 text-primary-600 dark:text-primary-400 rounded-full text-[10px] font-black uppercase tracking-widest shadow-sm flex items-center gap-2">
          <svg v-if="loadingActive" class="animate-spin h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
          </svg>
          {{ loadingActive ? 'SINKRONISASI...' : (metaActive.total + ' Aset') }}
        </span>
      </div>

      <div class="overflow-x-auto relative">
        <table class="w-full">
          <thead>
            <tr class="bg-gray-50/50 dark:bg-gray-700/30">
              <th v-for="h in ['Aset', 'Peminjam', 'Tgl Pinjam', 'Due Date', 'Status', 'Aksi']" :key="h"
                  class="text-left py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">{{ h }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50 dark:divide-gray-700">
            <template v-if="activeBorrowings.length > 0">
              <tr v-for="trx in activeBorrowings" :key="trx.id" 
                  @click="router.push(`/return?code=${trx.item_code}`)"
                  class="group hover:bg-primary-50/50 dark:hover:bg-primary-900/10 transition-all duration-300 cursor-pointer active:bg-primary-100/50 dark:active:bg-primary-900/20">
                <td class="px-8 py-6">
                  <div class="font-black text-gray-900 dark:text-white text-sm tracking-tight leading-none mb-1">{{ trx.item_name }}</div>
                  <div class="text-[10px] text-gray-400 font-mono tracking-tighter">{{ trx.item_code }}</div>
                </td>
                <td class="px-8 py-6">
                  <div class="text-xs font-bold text-gray-700 dark:text-gray-300">{{ trx.borrower_display }}</div>
                  <div v-if="trx.borrow_photo_url" class="mt-2 text-left">
                    <a :href="trx.borrow_photo_url" target="_blank" @click.stop class="inline-flex items-center gap-1 text-[9px] bg-primary-50 dark:bg-primary-900/30 text-primary-600 px-2 py-1 rounded-md hover:bg-primary-100 transition-colors cursor-pointer">
                      <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5zm10.5-11.25h.008v.008h-.008V8.25zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
                      </svg> Foto Pinjam
                    </a>
                  </div>
                </td>
                <td class="px-8 py-6 text-xs font-medium text-gray-500">{{ formatDate(trx.borrow_date) }}</td>
                <td class="px-8 py-6 text-xs font-black" :class="trx.is_overdue ? 'text-red-500 animate-pulse' : 'text-primary-600 dark:text-primary-400'">
                  {{ formatDate(trx.due_date) }}
                </td>
                <td class="px-8 py-6">
                  <div v-if="trx.is_overdue" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-red-50 dark:bg-red-900/30 text-red-600 font-black text-[9px] uppercase tracking-widest">
                    <span class="w-1.5 h-1.5 bg-red-600 rounded-full"></span>
                    TELAT
                  </div>
                  <div v-else class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 font-black text-[9px] uppercase tracking-widest">
                    <span class="w-1.5 h-1.5 bg-emerald-500 rounded-full"></span>
                    AKTIF
                  </div>
                </td>
                <td class="px-8 py-6">
                  <router-link :to="`/return?code=${trx.item_code}`" 
                               @click.stop
                               class="inline-flex items-center gap-2 px-5 py-2.5 bg-primary-50 dark:bg-primary-900/40 text-primary-600 dark:text-primary-400 rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-primary-600 hover:text-white transition-all shadow-sm active:scale-95 group/btn">
                    <span>Kembalikan</span>
                    <svg class="w-3.5 h-3.5 transform group-hover/btn:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
                    </svg>
                  </router-link>
                </td>
              </tr>
            </template>
            <tr v-else-if="!loadingActive" class="text-center">
              <td colspan="6" class="px-8 py-24 italic text-gray-400 font-bold uppercase tracking-[0.2em] text-xs">
                {{ searchQuery ? 'Aset aktif tidak ditemukan' : 'Semua barang telah Anda kembalikan' }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination Bar for Active -->
      <div class="px-8 py-5 border-t border-gray-100 dark:border-gray-700 flex flex-col sm:flex-row items-center justify-between gap-4 bg-gray-50/10">
        <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">
          Menampilkan <span class="text-primary-600">{{ activeBorrowings.length ? (metaActive.page - 1) * metaActive.page_size + 1 : 0 }}-{{ Math.min(metaActive.page * metaActive.page_size, metaActive.total) }}</span> dari <span class="text-gray-900 dark:text-gray-300">{{ metaActive.total }}</span> aktif
        </span>
        <div class="flex gap-2">
          <button @click="currentPageActive--" :disabled="currentPageActive === 1" class="pagination-btn-standard">
            Kembali
          </button>
          <div class="flex gap-1">
            <button v-for="p in visiblePagesActive" :key="p" @click="currentPageActive = p"
                    class="w-10 h-10 rounded-xl text-[11px] font-black transition-all shadow-sm active:scale-95 border"
                    :class="p === currentPageActive ? 'bg-primary-600 text-white border-primary-600' : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:bg-primary-50 dark:hover:bg-gray-700'">
              {{ p }}
            </button>
          </div>
          <button @click="currentPageActive++" :disabled="currentPageActive === totalPagesActive" class="pagination-btn-standard">
            Lanjut
          </button>
        </div>
      </div>
    </div>

    <!-- History Timeline Area -->
    <div class="space-y-6">
      <div class="flex items-center gap-4">
        <div class="h-px bg-gray-100 dark:bg-gray-700 flex-1"></div>
        <h2 class="text-xl font-black text-gray-400 dark:text-gray-600 uppercase tracking-widest leading-none">Timeline Pengembalian</h2>
        <div class="h-px bg-gray-100 dark:bg-gray-700 flex-1"></div>
      </div>

      <div class="relative">
        <div v-if="loadingHistory" class="absolute inset-0 bg-white/60 dark:bg-gray-900/60 backdrop-blur-sm z-10 flex flex-col items-center justify-center rounded-3xl">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        </div>

        <!-- History Timeline UI -->
        <div v-if="history.length > 0" class="relative pl-8 space-y-8 before:absolute before:inset-0 before:left-[11px] before:w-0.5 before:bg-gray-100 dark:before:bg-gray-700">
          <div v-for="trx in history" :key="trx.id" class="relative group animate-fade-in">
            <!-- Timeline Dot -->
            <div class="absolute -left-[31px] top-6 w-4 h-4 rounded-full border-4 border-white dark:border-gray-900 shadow-sm z-10 transition-transform group-hover:scale-125 bg-emerald-500"></div>
            
            <!-- Timeline Card -->
            <div class="card-premium !p-6 !rounded-[2.5rem] hover:ring-2 hover:ring-primary-500/20 transition-all border-gray-100 dark:border-gray-700">
              <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-6">
                <div class="flex items-center gap-5">
                  <div class="w-12 h-12 bg-gray-50 dark:bg-gray-900 rounded-2xl flex items-center justify-center text-gray-400 group-hover:text-primary-500 transition-colors">
                    <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-10.5v10.5" />
                    </svg>
                  </div>
                  <div>
                    <div class="text-[10px] font-black text-primary-500 uppercase tracking-widest mb-1">{{ formatDate(trx.returned_at) }}</div>
                    <h4 class="text-base font-black text-gray-900 dark:text-white uppercase tracking-tight">{{ trx.item_name }}</h4>
                    <p class="text-[10px] font-mono font-bold text-gray-400 tracking-tighter">{{ trx.item_code }}</p>
                  </div>
                </div>

                <div class="flex flex-wrap items-center gap-6 border-t lg:border-t-0 pt-4 lg:pt-0 border-gray-50 dark:border-gray-800">
                  <div class="space-y-1">
                    <p class="text-[9px] font-black text-gray-400 uppercase tracking-widest">Durasi Pinjam</p>
                    <p class="text-xs font-bold text-gray-700 dark:text-gray-300">{{ formatDate(trx.borrowed_at) }} s/d {{ formatDate(trx.returned_at) }}</p>
                  </div>
                  <div class="w-px h-8 bg-gray-100 dark:bg-gray-800 hidden lg:block"></div>
                  <div class="space-y-1">
                    <p class="text-[9px] font-black text-gray-400 uppercase tracking-widest">Kondisi Akhir</p>
                    <span class="px-3 py-1 rounded-lg text-[9px] font-black uppercase tracking-widest shadow-sm inline-block"
                          :class="getConditionBadgeClass(trx.return_condition)">
                      {{ getConditionLabel(trx.return_condition) }}
                    </span>
                  </div>
                </div>
              </div>
              
              <div v-if="trx.return_notes || trx.borrow_photo_url" class="mt-4 pt-4 border-t border-gray-50 dark:border-gray-800 flex items-center justify-between">
                 <p class="text-[10px] text-gray-400 italic font-medium max-w-sm truncate">{{ trx.return_notes || 'Tidak ada catatan pengembalian' }}</p>
                 <a v-if="trx.borrow_photo_url" :href="trx.borrow_photo_url" target="_blank" class="inline-flex items-center gap-1.5 text-[9px] font-black text-primary-600 hover:underline">
                    <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5zm10.5-11.25h.008v.008h-.008V8.25zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
                    </svg> Lihat Foto Pinjam
                 </a>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty History State -->
        <div v-else-if="!loadingHistory" class="py-20 text-center space-y-4 bg-gray-50/50 dark:bg-gray-800/30 rounded-[3rem] border border-dashed border-gray-100 dark:border-gray-700">
           <p class="text-xs font-black text-gray-400 uppercase tracking-[0.2em]">{{ searchQuery ? 'Riwayat riwayat tidak ditemukan' : 'Belum ada riwayat tercatat' }}</p>
        </div>
      </div>

    <!-- Pagination Bar for History -->
    <div class="px-8 py-6 bg-white dark:bg-gray-800 rounded-[2rem] border border-gray-100 dark:border-gray-700 shadow-sm flex flex-col sm:flex-row items-center justify-between gap-4">
      <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">
        Menampilkan <span class="text-primary-600">{{ history.length ? (metaHistory.page - 1) * metaHistory.page_size + 1 : 0 }}-{{ Math.min(metaHistory.page * metaHistory.page_size, metaHistory.total) }}</span> dari <span class="text-gray-900 dark:text-gray-300">{{ metaHistory.total }}</span> riwayat
      </span>
      <div class="flex gap-2">
        <button @click="currentPageHistory--" :disabled="currentPageHistory === 1" class="pagination-btn-standard">
          Kembali
        </button>
        <div class="flex gap-1">
          <button v-for="p in visiblePagesHistory" :key="p" @click="currentPageHistory = p"
                  class="w-10 h-10 rounded-xl text-[11px] font-black transition-all shadow-sm active:scale-95 border"
                  :class="p === currentPageHistory ? 'bg-primary-600 text-white border-primary-600' : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:bg-primary-50 dark:hover:bg-gray-700'">
            {{ p }}
          </button>
        </div>
        <button @click="currentPageHistory++" :disabled="currentPageHistory === totalPagesHistory" class="pagination-btn-standard">
          Lanjut
        </button>
      </div>
    </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import api from '../utils/api'

const router = useRouter()

const activeBorrowings = ref([])
const history = ref([])
const loadingActive = ref(false)
const loadingHistory = ref(false)

const searchQuery = ref('')
const perPage = 10

const metaActive = ref({ total: 0, total_pages: 1, page: 1, page_size: 10 })
const metaHistory = ref({ total: 0, total_pages: 1, page: 1, page_size: 10 })

// ---------------------------
// ACTIVE BORROWINGS LOGIC
// ---------------------------
const currentPageActive = ref(1)

const fetchActive = async () => {
  loadingActive.value = true
  try {
    const { data } = await api.get('/transactions/my', {
      params: {
        page: currentPageActive.value,
        page_size: perPage,
        search: searchQuery.value
      }
    })
    if (data.success) { 
      activeBorrowings.value = data.data.items || []
      metaActive.value = data.data.meta || { total: 0, total_pages: 1 }
    }
  } catch (e) { 
    console.error('Gagal fetch active borrowings:', e)
    activeBorrowings.value = []
  } finally { 
    loadingActive.value = false 
  }
}

const totalPagesActive = computed(() => metaActive.value.total_pages)
const startRowActive = computed(() => (currentPageActive.value - 1) * perPage + 1)
const endRowActive = computed(() => Math.min(currentPageActive.value * perPage, metaActive.value.total))

const visiblePagesActive = computed(() => {
  const pages = []
  const start = Math.max(1, currentPageActive.value - 2)
  const end = Math.min(totalPagesActive.value, currentPageActive.value + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

// ---------------------------
// HISTORY LOGIC
// ---------------------------
const currentPageHistory = ref(1)

const fetchHistory = async () => {
  loadingHistory.value = true
  try {
    const { data } = await api.get('/transactions/my/history', {
      params: {
        page: currentPageHistory.value,
        page_size: perPage,
        search: searchQuery.value
      }
    })
    if (data.success) { 
      history.value = data.data.items || []
      metaHistory.value = data.data.meta || { total: 0, total_pages: 1, page: 1, page_size: 10 }
    }
  } catch (e) { 
    console.error('Gagal fetch history:', e)
    history.value = []
  } finally { 
    loadingHistory.value = false 
  }
}

const totalPagesHistory = computed(() => metaHistory.value.total_pages)
const startRowHistory = computed(() => (currentPageHistory.value - 1) * perPage + 1)
const endRowHistory = computed(() => Math.min(currentPageHistory.value * perPage, metaHistory.value.total))

const visiblePagesHistory = computed(() => {
  const pages = []
  const start = Math.max(1, currentPageHistory.value - 2)
  const end = Math.min(totalPagesHistory.value, currentPageHistory.value + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

// Debounced search
let searchTimeout
const debouncedSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    currentPageActive.value = 1
    currentPageHistory.value = 1
    fetchActive()
    fetchHistory()
  }, 400)
}

watch(currentPageActive, fetchActive)
watch(currentPageHistory, fetchHistory)

const formatDate = (d) => {
  if (!d) return '-'
  return new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }).format(new Date(d))
}

const getConditionBadgeClass = (c) => ({
  'GOOD': 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600',
  'DAMAGED': 'bg-rose-50 dark:bg-rose-900/30 text-rose-600',
  'IN_REPAIR': 'bg-amber-50 dark:bg-amber-900/30 text-amber-600'
}[c] || 'bg-gray-50 text-gray-600')

const getConditionLabel = (c) => ({ 'GOOD': 'Baik', 'DAMAGED': 'Rusak', 'IN_REPAIR': 'Perbaikan' }[c] || c)

onMounted(() => { fetchActive(); fetchHistory() })
</script>
