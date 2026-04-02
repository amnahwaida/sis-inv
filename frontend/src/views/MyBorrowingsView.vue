<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Peminjaman Saya</h1>
          <p class="text-primary-100/70 text-sm font-medium">Pantau status barang yang sedang Anda gunakan</p>
        </div>

        <div class="flex items-center gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10 w-full md:w-auto">
          <!-- Search Input -->
          <div class="relative w-full sm:w-64 group">
            <input type="text" v-model="searchQuery" placeholder="Cari barang peminjaman..." 
                   class="input-field pl-10 h-10 rounded-xl text-xs w-full bg-white/10 text-white placeholder-white/50 border-white/10 focus:ring-white/20" />
            <svg class="w-4 h-4 absolute left-3.5 top-3 text-white/50" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Active Borrowings Table Card -->
    <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden transition-all duration-300">
      <div class="p-8 border-b border-gray-100 dark:border-gray-700 flex items-center justify-between">
        <h2 class="text-lg font-black text-gray-900 dark:text-white uppercase tracking-tight">Barang Belum Kembali</h2>
        <span class="px-4 py-1 bg-primary-50 dark:bg-primary-900/40 text-primary-600 dark:text-primary-400 rounded-full text-[10px] font-black uppercase tracking-widest shadow-sm flex items-center gap-2">
          <svg v-if="loadingActive" class="animate-spin h-3 w-3" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
          {{ loadingActive ? 'SINKRONISASI...' : (filteredActive.length + ' Aset') }}
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
            <template v-if="paginatedActive.length > 0">
              <tr v-for="trx in paginatedActive" :key="trx.id" 
                  @click="$router.push(`/return?code=${trx.item_code}`)"
                  class="group hover:bg-primary-50/50 dark:hover:bg-primary-900/10 transition-all duration-300 cursor-pointer active:bg-primary-100/50 dark:active:bg-primary-900/20">
                <td class="px-8 py-6">
                  <div class="font-black text-gray-900 dark:text-white text-sm tracking-tight leading-none mb-1">{{ trx.item_name }}</div>
                  <div class="text-[10px] text-gray-400 font-mono tracking-tighter">{{ trx.item_code }}</div>
                </td>
                <td class="px-8 py-6">
                  <div class="text-xs font-bold text-gray-700 dark:text-gray-300">{{ trx.borrower_display }}</div>
                  <div v-if="trx.borrow_photo_url" class="mt-2 text-left">
                    <a :href="trx.borrow_photo_url" target="_blank" @click.stop class="inline-flex items-center gap-1 text-[9px] bg-primary-50 dark:bg-primary-900/30 text-primary-600 px-2 py-1 rounded-md hover:bg-primary-100 transition-colors cursor-pointer">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg> Foto Pinjam
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
                               class="btn-action-view !text-indigo-600 !border-indigo-100 dark:!border-indigo-800">
                    Kembalikan
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M14 5l7 7m0 0l-7 7m7-7H3"/></svg>
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
      <div v-if="totalPagesActive > 1" class="px-8 py-5 border-t border-gray-100 dark:border-gray-700 flex flex-col sm:flex-row items-center justify-between gap-4">
        <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">
          Menampilkan <span class="text-primary-600">{{ startRowActive }}-{{ endRowActive }}</span> dari <span class="text-gray-900 dark:text-gray-300">{{ filteredActive.length }}</span> aktif
        </span>
        <div class="flex gap-2">
          <button @click="currentPageActive--" :disabled="currentPageActive === 1" class="pagination-btn-standard">
            Kembali
          </button>
          <button v-for="p in visiblePagesActive" :key="p" @click="currentPageActive = p"
                  class="w-10 h-10 rounded-xl text-[11px] font-black transition-all shadow-sm active:scale-95 border"
                  :class="p === currentPageActive ? 'bg-primary-600 text-white border-primary-600' : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:bg-primary-50 dark:hover:bg-gray-700'">
            {{ p }}
          </button>
          <button @click="currentPageActive++" :disabled="currentPageActive === totalPagesActive" class="pagination-btn-standard">
            Lanjut
          </button>
        </div>
      </div>
    </div>

    <!-- History Header -->
    <div class="pt-10 flex items-center gap-4">
      <div class="h-px bg-gray-100 dark:bg-gray-700 flex-1"></div>
      <h2 class="text-xl font-black text-gray-400 dark:text-gray-600 uppercase tracking-widest leading-none">Riwayat Peminjaman</h2>
      <div class="h-px bg-gray-100 dark:bg-gray-700 flex-1"></div>
    </div>

    <!-- History Table Card -->
    <div class="bg-gray-50 dark:bg-gray-900/20 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 overflow-hidden transition-all duration-300">
      <div class="overflow-x-auto relative">
        <table class="w-full">
          <thead>
            <tr class="bg-gray-100/50 dark:bg-gray-800/50">
              <th v-for="h in ['Barang', 'Tgl Pinjam', 'Tgl Kembali', 'Catatan', 'Kondisi Akhir']" :key="h"
                  class="text-left py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">{{ h }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 dark:divide-gray-800">
            <template v-if="paginatedHistory.length > 0">
              <tr v-for="trx in paginatedHistory" :key="trx.id" class="group hover:bg-white dark:hover:bg-gray-800 transition-all duration-300">
                <td class="px-8 py-6">
                  <div class="font-bold text-gray-700 dark:text-gray-200 text-sm">{{ trx.item_name }}</div>
                  <div class="text-[10px] text-gray-400 font-mono tracking-tighter mb-2">{{ trx.item_code }}</div>
                  <div v-if="trx.borrow_photo_url" class="text-left mt-1">
                    <a :href="trx.borrow_photo_url" target="_blank" @click.stop class="inline-flex items-center gap-1 text-[9px] bg-primary-50 dark:bg-primary-900/30 text-primary-600 px-2 py-1 rounded-md hover:bg-primary-100 transition-colors cursor-pointer">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg> Foto Pinjam
                    </a>
                  </div>
                </td>
                <td class="px-8 py-6 text-xs text-gray-500 font-medium">{{ formatDate(trx.borrowed_at) }}</td>
                <td class="px-8 py-6 text-xs text-gray-500 font-medium">{{ formatDate(trx.returned_at) }}</td>
                <td class="px-8 py-6">
                  <p class="text-xs text-gray-400 italic leading-snug max-w-xs truncate" :title="trx.return_notes || '-'">{{ trx.return_notes || '-' }}</p>
                </td>
                <td class="px-8 py-6">
                  <span v-if="trx.return_condition" class="px-3 py-1 rounded-lg text-[9px] font-black uppercase tracking-widest shadow-sm inline-block"
                        :class="getConditionBadgeClass(trx.return_condition)">
                    {{ getConditionLabel(trx.return_condition) }}
                  </span>
                </td>
              </tr>
            </template>
            <tr v-else-if="!loadingHistory" class="text-center">
              <td colspan="5" class="px-8 py-16 text-gray-400 italic text-sm">
                {{ searchQuery ? 'Riwayat aset tidak ditemukan' : 'Belum ada riwayat tercatat' }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination Bar for History -->
      <div v-if="totalPagesHistory > 1" class="px-8 py-5 border-t border-gray-100 dark:border-gray-800 flex flex-col sm:flex-row items-center justify-between gap-4">
        <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">
          Menampilkan <span class="text-primary-600">{{ startRowHistory }}-{{ endRowHistory }}</span> dari <span class="text-gray-900 dark:text-gray-300">{{ filteredHistory.length }}</span> riwayat
        </span>
        <div class="flex gap-2">
          <button @click="currentPageHistory--" :disabled="currentPageHistory === 1" class="pagination-btn-standard">
            Kembali
          </button>
          <button v-for="p in visiblePagesHistory" :key="p" @click="currentPageHistory = p"
                  class="w-10 h-10 rounded-xl text-[11px] font-black transition-all shadow-sm active:scale-95 border"
                  :class="p === currentPageHistory ? 'bg-primary-600 text-white border-primary-600' : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:bg-primary-50 dark:hover:bg-gray-700'">
            {{ p }}
          </button>
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
import api from '../utils/api'

const activeBorrowings = ref([])
const history = ref([])
const loadingActive = ref(false)
const loadingHistory = ref(false)

const searchQuery = ref('')
const perPage = 10

// ---------------------------
// ACTIVE BORROWINGS LOGIC
// ---------------------------
const currentPageActive = ref(1)

const filteredActive = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return activeBorrowings.value
  return activeBorrowings.value.filter(trx => 
    trx.item_name?.toLowerCase().includes(q) || 
    trx.item_code?.toLowerCase().includes(q)
  )
})

const totalPagesActive = computed(() => Math.max(1, Math.ceil(filteredActive.value.length / perPage)))
const startRowActive = computed(() => (currentPageActive.value - 1) * perPage + 1)
const endRowActive = computed(() => Math.min(currentPageActive.value * perPage, filteredActive.value.length))
const paginatedActive = computed(() => filteredActive.value.slice((currentPageActive.value - 1) * perPage, currentPageActive.value * perPage))

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

const filteredHistory = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return history.value
  return history.value.filter(trx => 
    trx.item_name?.toLowerCase().includes(q) || 
    trx.item_code?.toLowerCase().includes(q)
  )
})

const totalPagesHistory = computed(() => Math.max(1, Math.ceil(filteredHistory.value.length / perPage)))
const startRowHistory = computed(() => (currentPageHistory.value - 1) * perPage + 1)
const endRowHistory = computed(() => Math.min(currentPageHistory.value * perPage, filteredHistory.value.length))
const paginatedHistory = computed(() => filteredHistory.value.slice((currentPageHistory.value - 1) * perPage, currentPageHistory.value * perPage))

const visiblePagesHistory = computed(() => {
  const pages = []
  const start = Math.max(1, currentPageHistory.value - 2)
  const end = Math.min(totalPagesHistory.value, currentPageHistory.value + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

// Reset pagination when searching
watch(searchQuery, () => {
  currentPageActive.value = 1
  currentPageHistory.value = 1
})

const fetchActive = async () => {
  console.log('Fetching active borrowings...')
  loadingActive.value = true
  try {
    const { data } = await api.get('/transactions/my', { timeout: 15000 })
    console.log('Active borrowings received:', data)
    if (data.success) { activeBorrowings.value = data.data }
  } catch (e) { 
    console.error('Gagal fetch active borrowings:', e)
    activeBorrowings.value = []
  } finally { 
    loadingActive.value = false 
    console.log('Active borrowings fetch complete. Loading state:', loadingActive.value)
  }
}

const fetchHistory = async () => {
  console.log('Fetching history...')
  loadingHistory.value = true
  try {
    const { data } = await api.get('/transactions/my/history', { timeout: 15000 })
    console.log('History received:', data)
    if (data.success) { history.value = data.data }
  } catch (e) { 
    console.error('Gagal fetch history:', e)
    history.value = []
  } finally { 
    loadingHistory.value = false 
    console.log('History fetch complete. Loading state:', loadingHistory.value)
  }
}

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
