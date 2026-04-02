<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col lg:flex-row lg:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Laporan & Monitoring</h1>
          <p class="text-primary-100/70 text-sm font-medium">Pemantauan transaksi dan aset secara real-time</p>
        </div>
        
        <div class="flex flex-col sm:flex-row items-stretch gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10 w-full lg:w-auto">
          <button @click="exportTransactions" class="bg-white/10 hover:bg-white/20 text-white border border-white/10 px-6 py-3 rounded-xl text-[10px] font-black transition-all flex items-center justify-center gap-2 active:scale-95 w-full">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
            EXPORT LAPORAN (EXCEL)
          </button>
        </div>
      </div>
    </div>

    <!-- Tabs & Filter Section -->
    <div class="flex flex-col md:flex-row items-center justify-between gap-6">
      <div class="flex flex-wrap items-center justify-center gap-2 bg-white dark:bg-gray-800 p-2 rounded-[1.5rem] border border-gray-100 dark:border-gray-700 shadow-sm w-full md:w-fit group">
        <button v-for="tab in [{id:'active', label:'Peminjaman Aktif', icon:'⚡'}, {id:'overdue', label:'Terlambat', icon:'⚠️'}, {id:'history', label:'Riwayat Lengkap', icon:'📜'}]" 
                :key="tab.id"
                @click="activeTab = tab.id" 
                :class="activeTab === tab.id ? 'bg-primary-600 text-white shadow-lg' : 'text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'" 
                class="px-5 py-2.5 rounded-xl text-xs font-black transition-all duration-300 flex items-center gap-2 whitespace-nowrap active:scale-95">
          <span>{{ tab.icon }}</span>
          <span>{{ tab.label }}</span>
          <span v-if="tab.id === 'active' && activeBorrowings.length" class="ml-1 px-2 py-0.5 bg-white/20 rounded-full text-[10px]">{{ activeBorrowings.length }}</span>
          <span v-if="tab.id === 'overdue' && overdueItems.length" class="ml-1 px-2 py-0.5 bg-white/20 rounded-full text-[10px]">{{ overdueItems.length }}</span>
        </button>
      </div>

      <!-- Filters Row -->
      <div class="flex flex-col sm:flex-row items-center gap-3 w-full md:w-auto">
        <!-- Search Input -->
        <div class="relative w-full sm:w-64">
          <input type="text" v-model="searchQuery" placeholder="Cari nama / kode..." 
                 class="input-field pl-10 h-11 rounded-2xl text-sm w-full" />
          <svg class="w-4 h-4 absolute left-3.5 top-3.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
        </div>

        <!-- Class Filter -->
        <div class="flex items-center gap-2 bg-white dark:bg-gray-800 p-2 px-4 rounded-2xl border border-gray-100 dark:border-gray-700 shadow-sm w-full sm:w-auto">
          <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest whitespace-nowrap">Kelas:</label>
          <select v-model="classFilter" @change="handleFilterChange" class="bg-transparent border-none focus:ring-0 text-sm font-black text-gray-900 dark:text-white py-0 w-full sm:w-auto">
            <option value="">Semua</option>
            <option v-for="cls in uniqueClasses" :key="cls" :value="cls">{{ cls }}</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Tables Container -->
    <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden transition-all duration-300">
      <div class="overflow-x-auto relative">
        <!-- Loading Overlay -->
        <div v-if="loading" class="absolute inset-0 bg-white/60 dark:bg-gray-900/60 backdrop-blur-sm z-10 flex flex-col items-center justify-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
          <span class="mt-4 text-xs font-black text-primary-600 uppercase tracking-widest">Sinkronisasi Data...</span>
        </div>

        <!-- Active Borrowings Table -->
        <table v-if="activeTab === 'active'" class="w-full">
          <thead>
            <tr class="bg-gray-50/50 dark:bg-gray-700/30">
              <th v-for="h in ['Barang / Aset', 'Peminjam', 'Kelas', 'Petugas', 'Waktu Pinjam', 'Sisa Waktu']" :key="h"
                  class="text-left py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] whitespace-nowrap">{{ h }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50 dark:divide-gray-700">
            <tr v-if="paginatedData.length === 0" class="text-center">
              <td colspan="6" class="px-8 py-24">
                <div class="w-20 h-20 bg-emerald-50 dark:bg-emerald-900/20 rounded-full flex items-center justify-center mx-auto mb-6 text-emerald-500">
                  <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                </div>
                <p class="text-xl font-black text-gray-900 dark:text-white uppercase tracking-tight">{{ searchQuery ? 'Tidak Ditemukan' : 'Semua Aset Tersedia' }}</p>
                <p class="text-xs text-gray-400 font-medium mt-1">{{ searchQuery ? 'Coba kata kunci lain.' : 'Tidak ada peminjaman aktif saat ini.' }}</p>
              </td>
            </tr>
            <tr v-for="trx in paginatedData" :key="trx.id" class="hover:bg-primary-50/50 dark:hover:bg-primary-900/10 transition-all duration-300">
              <td class="px-8 py-6">
                <div class="font-black text-gray-900 dark:text-white text-sm leading-none">{{ trx.item_name }}</div>
                <div class="text-[10px] text-gray-400 font-bold font-mono mt-1 tracking-tighter">{{ trx.item_code }}</div>
              </td>
              <td class="px-8 py-6">
                <div class="text-sm font-black text-gray-900 dark:text-white tracking-tight">{{ trx.student_name || trx.teacher_name }}</div>
                <span class="text-[9px] font-black uppercase inline-block mt-1 px-1.5 py-0.5 rounded-md" :class="trx.borrower_type === 'STUDENT' ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600' : 'bg-gray-100 dark:bg-gray-700 text-gray-600'">
                  {{ trx.borrower_type === 'STUDENT' ? 'Siswa' : 'Staf' }}
                </span>
              </td>
              <td class="px-8 py-6 text-sm font-bold text-gray-500">{{ trx.student_class || '-' }}</td>
              <td class="px-8 py-6">
                <div class="text-[10px] font-black text-gray-400 uppercase tracking-widest">{{ trx.teacher_name }}</div>
                <div v-if="trx.borrow_photo_url" class="mt-2 text-left">
                  <a :href="trx.borrow_photo_url" target="_blank" @click.stop class="inline-flex items-center gap-1 text-[9px] bg-primary-50 dark:bg-primary-900/30 text-primary-600 px-2 py-1 rounded-md hover:bg-primary-100 transition-colors cursor-pointer">
                    <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg> Foto Pinjam
                  </a>
                </div>
              </td>
              <td class="px-8 py-6">
                <div class="text-xs font-bold text-gray-700 dark:text-gray-300">{{ formatDate(trx.borrowed_at) }}</div>
              </td>
              <td class="px-8 py-6">
                <div v-if="trx.is_overdue" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-red-50 dark:bg-red-900/30 text-red-600 font-black text-[9px] uppercase tracking-widest animate-pulse">
                  <span class="w-1.5 h-1.5 bg-red-600 rounded-full"></span>
                  Terlambat
                </div>
                <div v-else class="text-[10px] font-black text-emerald-500 uppercase tracking-widest">
                  Sampai {{ formatDate(trx.due_date) }}
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- Overdue / History Table -->
        <table v-if="activeTab === 'overdue' || activeTab === 'history'" class="w-full">
          <thead>
            <tr class="bg-gray-50/50 dark:bg-gray-700/30">
              <th v-for="h in ['Barang', 'Peminjam', 'Tanggal Pinjam', activeTab === 'overdue' ? 'Jatuh Tempo' : 'Tgl Kembali', 'Status']" :key="h"
                  class="text-left py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] whitespace-nowrap">{{ h }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50 dark:divide-gray-700">
            <template v-if="paginatedData.length > 0">
              <tr v-for="trx in paginatedData" :key="trx.id" class="hover:bg-primary-50/50 dark:hover:bg-primary-900/10 transition-all duration-300">
                <td class="px-8 py-6">
                  <div class="font-black text-gray-900 dark:text-white text-sm tracking-tight">{{ trx.item_name }}</div>
                  <div class="text-[10px] text-gray-400 font-mono tracking-tighter">{{ trx.item_code }}</div>
                </td>
                <td class="px-8 py-6">
                  <div class="text-xs font-black text-gray-900 dark:text-white">{{ trx.student_name || trx.teacher_name }}</div>
                  <div class="text-[10px] text-gray-400">{{ trx.student_class || 'Staff' }}</div>
                  <div v-if="trx.borrow_photo_url" class="mt-2 text-left">
                    <a :href="trx.borrow_photo_url" target="_blank" @click.stop class="inline-flex items-center gap-1 text-[9px] bg-primary-50 dark:bg-primary-900/30 text-primary-600 px-2 py-1 rounded-md hover:bg-primary-100 transition-colors cursor-pointer">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg> Foto Pinjam
                    </a>
                  </div>
                </td>
                <td class="px-8 py-6 text-xs font-bold text-gray-500">{{ formatDate(trx.borrowed_at) }}</td>
                <td class="px-8 py-6 text-xs" :class="activeTab === 'overdue' ? 'text-red-500 font-black' : 'text-gray-500 font-bold'">
                  {{ activeTab === 'overdue' ? formatDate(trx.due_date) : formatDate(trx.returned_at) }}
                </td>
                <td class="px-8 py-6">
                  <span class="px-3 py-1 rounded-lg text-[9px] font-black uppercase tracking-widest shadow-sm"
                        :class="trx.status === 'RETURNED' ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600' : 'bg-red-50 dark:bg-red-900/30 text-red-600'">
                    {{ trx.status || (trx.is_overdue ? 'OVERDUE' : 'ACTIVE') }}
                  </span>
                </td>
              </tr>
            </template>
            <tr v-else class="text-center">
              <td colspan="5" class="px-8 py-24 italic text-gray-400 font-bold tracking-[0.2em] text-[10px] uppercase">{{ searchQuery ? 'Tidak Ditemukan' : 'Tidak Ada Data Ditemukan' }}</td>
            </tr>
          </tbody>
        </table>

      </div>

      <!-- Pagination Bar -->
      <div v-if="totalPages > 1" class="px-8 py-5 bg-gray-50/50 dark:bg-gray-700/20 border-t border-gray-100 dark:border-gray-700 flex flex-col sm:flex-row items-center justify-between gap-4">
        <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">
          Menampilkan <span class="text-primary-600">{{ startRow }}-{{ endRow }}</span> dari <span class="text-gray-900 dark:text-white">{{ filteredData.length }}</span> data
        </span>
        <div class="flex gap-2">
          <button @click="currentPage--" :disabled="currentPage === 1" class="pagination-btn-standard">
            Kembali
          </button>
          <button v-for="p in visiblePages" :key="p" @click="currentPage = p"
                  class="w-10 h-10 rounded-xl text-[11px] font-black transition-all shadow-sm active:scale-95 border"
                  :class="p === currentPage ? 'bg-primary-600 text-white border-primary-600' : 'bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:bg-primary-50 dark:hover:bg-gray-700'">
            {{ p }}
          </button>
          <button @click="currentPage++" :disabled="currentPage === totalPages" class="pagination-btn-standard">
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

const activeTab = ref('active')
const activeBorrowings = ref([])
const overdueItems = ref([])
const histories = ref([])
const uniqueClasses = ref([])
const classFilter = ref('')
const searchQuery = ref('')
const loading = ref(false)
const currentPage = ref(1)
const perPage = 15

// Current tab's raw data
const currentRawData = computed(() => {
  if (activeTab.value === 'active') return activeBorrowings.value
  if (activeTab.value === 'overdue') return overdueItems.value
  return histories.value
})

// Filtered data (search applied)
const filteredData = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return currentRawData.value
  return currentRawData.value.filter(trx => {
    const fields = [
      trx.item_name, trx.item_code, trx.student_name, 
      trx.teacher_name, trx.student_class
    ].filter(Boolean).join(' ').toLowerCase()
    return fields.includes(q)
  })
})

// Pagination
const totalPages = computed(() => Math.max(1, Math.ceil(filteredData.value.length / perPage)))
const startRow = computed(() => (currentPage.value - 1) * perPage + 1)
const endRow = computed(() => Math.min(currentPage.value * perPage, filteredData.value.length))
const paginatedData = computed(() => filteredData.value.slice((currentPage.value - 1) * perPage, currentPage.value * perPage))

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

// Reset page on search or tab change
watch([searchQuery, activeTab], () => { currentPage.value = 1 })

async function fetchActive() {
  loading.value = true
  try {
    const params = classFilter.value ? `?class=${encodeURIComponent(classFilter.value)}` : ''
    const { data } = await api.get(`/reports/active-borrowings${params}`)
    if (data.success) { activeBorrowings.value = data.data }
  } catch (e) { console.error(e) } 
  finally { loading.value = false }
}

async function fetchOverdue() {
  try {
    const params = classFilter.value ? `?class=${encodeURIComponent(classFilter.value)}` : ''
    const { data } = await api.get(`/reports/overdue${params}`)
    if (data.success) { overdueItems.value = data.data }
  } catch (e) { console.error(e) }
}

async function fetchHistory() {
  loading.value = true
  try {
    const params = classFilter.value ? `?class=${encodeURIComponent(classFilter.value)}` : ''
    const { data } = await api.get(`/reports/history${params}`)
    if (data.success) { histories.value = data.data || [] }
  } catch (e) { console.error(e) } 
  finally { loading.value = false }
}

function handleFilterChange() {
  currentPage.value = 1
  if (activeTab.value === 'active') fetchActive()
  else if (activeTab.value === 'overdue') fetchOverdue()
  else if (activeTab.value === 'history') fetchHistory()
}

async function fetchClasses() {
  try {
    const { data } = await api.get('/students/classes')
    if (data.success) { uniqueClasses.value = data.data || [] }
  } catch (e) { console.error('Gagal mengambil daftar kelas dari master data:', e) }
}

async function exportTransactions() {
  try {
    const params = new URLSearchParams()
    if (classFilter.value) params.append('class', classFilter.value)
    params.append('type', activeTab.value)
    
    const response = await api.get(`/reports/export/transactions?${params.toString()}`, { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url; link.setAttribute('download', `SIS-INV_Laporan_${new Date().toISOString().slice(0,10)}.xlsx`)
    document.body.appendChild(link); link.click(); document.body.removeChild(link)
  } catch (e) { alert('Gagal mendownload laporan') }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }).format(new Date(dateStr))
}

watch(activeTab, (newTab) => {
  if (newTab === 'active') fetchActive()
  else if (newTab === 'overdue') fetchOverdue()
  else if (newTab === 'history') fetchHistory()
})

onMounted(() => {
  fetchActive(); fetchOverdue(); fetchHistory(); fetchClasses()
})
</script>
