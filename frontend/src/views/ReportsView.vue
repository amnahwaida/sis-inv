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
         <button v-for="tab in [
                   {id:'active', label:'Peminjaman Aktif', icon:'M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z'}, 
                   {id:'overdue', label:'Terlambat', icon:'M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z'}, 
                   {id:'history', label:'Riwayat Lengkap', icon:'M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z'}
                 ]" 
                 :key="tab.id"
                 @click="activeTab = tab.id" 
                 :class="activeTab === tab.id ? 'bg-primary-600 text-white shadow-lg' : 'text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'" 
                 class="px-5 py-2.5 rounded-xl text-xs font-black transition-all duration-300 flex items-center gap-2 whitespace-nowrap active:scale-95">
           <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
             <path stroke-linecap="round" stroke-linejoin="round" :d="tab.icon" />
           </svg>
           <span>{{ tab.label }}</span>
         </button>
      </div>

      <!-- Filters Row -->
      <div class="bg-white dark:bg-gray-800 p-3 rounded-[2rem] border border-gray-100 dark:border-gray-700 shadow-sm flex flex-col sm:flex-row items-center gap-3 w-full md:w-auto">
        <!-- Search Input -->
        <div class="relative w-full sm:w-64">
          <input type="text" v-model="filters.search" @input="debouncedFetch" placeholder="Cari nama / kode..." 
                 class="input-field pl-10 h-12 rounded-2xl text-sm w-full" />
          <svg class="w-5 h-5 absolute left-3 top-3.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
        </div>

        <!-- Class Filter -->
        <select v-model="filters.class" @change="handleFilterChange" 
                class="input-field h-12 rounded-2xl text-sm w-full sm:w-auto appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2020%2020%22%3E%3Cpath%20stroke%3D%22%236b7280%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%221.5%22%20d%3D%22m6%208%204%204%204-4%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[right_0.5rem_center] bg-no-repeat pr-10">
          <option value="">Semua Kelas</option>
          <option v-for="cls in uniqueClasses" :key="cls" :value="cls">{{ cls }}</option>
        </select>
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
            <tr v-if="currentItems.length === 0" class="text-center">
              <td colspan="6" class="px-8 py-24">
                <div class="w-20 h-20 bg-emerald-50 dark:bg-emerald-900/20 rounded-full flex items-center justify-center mx-auto mb-6 text-emerald-500">
                  <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                </div>
                <p class="text-xl font-black text-gray-900 dark:text-white uppercase tracking-tight">{{ filters.search ? 'Tidak Ditemukan' : 'Semua Aset Tersedia' }}</p>
                <p class="text-xs text-gray-400 font-medium mt-1">{{ filters.search ? 'Coba kata kunci lain.' : 'Tidak ada peminjaman aktif saat ini.' }}</p>
              </td>
            </tr>
            <tr v-for="trx in currentItems" :key="trx.id" class="hover:bg-primary-50/50 dark:hover:bg-primary-900/10 transition-all duration-300">
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
            <template v-if="currentItems.length > 0">
              <tr v-for="trx in currentItems" :key="trx.id" class="hover:bg-primary-50/50 dark:hover:bg-primary-900/10 transition-all duration-300">
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
              <td colspan="5" class="px-8 py-24 italic text-gray-400 font-bold tracking-[0.2em] text-[10px] uppercase">{{ filters.search ? 'Tidak Ditemukan' : 'Tidak Ada Data Ditemukan' }}</td>
            </tr>
          </tbody>
        </table>

      </div>

      <!-- Pagination Bar -->
      <div class="px-8 py-5 border-t border-gray-100 dark:border-gray-700 flex flex-col sm:flex-row items-center justify-between gap-4 bg-gray-50/10">
        <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">
          Menampilkan <span class="text-primary-600">{{ currentItems.length ? (meta.page - 1) * meta.page_size + 1 : 0 }}-{{ Math.min(meta.page * meta.page_size, meta.total) }}</span> dari <span class="text-gray-900 dark:text-gray-300">{{ meta.total }}</span> laporan
        </span>
        <div class="flex gap-2">
          <button @click="changePage(meta.page - 1)" :disabled="meta.page === 1" class="pagination-btn-standard">
            Kembali
          </button>
          <div class="flex gap-1">
            <button v-for="p in visiblePages" :key="p" @click="changePage(p)"
                    class="w-10 h-10 rounded-xl text-[11px] font-black transition-all shadow-sm active:scale-95 border"
                    :class="p === meta.page ? 'bg-primary-600 text-white border-primary-600' : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:bg-primary-50 dark:hover:bg-gray-700'">
              {{ p }}
            </button>
          </div>
          <button @click="changePage(meta.page + 1)" :disabled="meta.page === meta.total_pages" class="pagination-btn-standard">
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
const currentItems = ref([])
const meta = ref({ page: 1, page_size: 15, total: 0, total_pages: 1 })
const uniqueClasses = ref([])
const filters = ref({ search: '', class: '', page: 1, page_size: 15 })
const loading = ref(false)

let searchTimeout
const debouncedFetch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    filters.value.page = 1
    fetchData()
  }, 400)
}

const changePage = (p) => {
  if (p < 1 || p > meta.value.total_pages) return
  filters.value.page = p
  fetchData()
}

const visiblePages = computed(() => {
  const pages = []
  const current = meta.value.page
  const total = meta.value.total_pages
  const start = Math.max(1, current - 2)
  const end = Math.min(total, current + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

const buildParams = () => {
  const p = new URLSearchParams()
  p.append('page', filters.value.page)
  p.append('page_size', filters.value.page_size)
  if (filters.value.search) p.append('search', filters.value.search)
  if (filters.value.class) p.append('class', filters.value.class)
  return p.toString()
}

async function fetchData() {
  loading.value = true
  try {
    let endpoint = ''
    if (activeTab.value === 'active') endpoint = '/reports/active-borrowings'
    else if (activeTab.value === 'overdue') endpoint = '/reports/overdue'
    else if (activeTab.value === 'history') endpoint = '/reports/history'
    
    const { data } = await api.get(`${endpoint}?${buildParams()}`)
    if (data.success) {
      currentItems.value = data.data.items || []
      meta.value = data.data.meta || { page: 1, page_size: 15, total: 0, total_pages: 1 }
    }
  } catch (e) {
    console.error(e)
    currentItems.value = []
  } finally {
    loading.value = false
  }
}

function handleFilterChange() {
  filters.value.page = 1
  fetchData()
}

watch(activeTab, () => {
  filters.value.page = 1
  filters.value.search = ''
  fetchData()
})

async function fetchClasses() {
  try {
    const { data } = await api.get('/students/classes')
    if (data.success) { uniqueClasses.value = data.data || [] }
  } catch (e) { console.error('Gagal mengambil daftar kelas dari master data:', e) }
}

async function exportTransactions() {
  try {
    const params = new URLSearchParams()
    if (filters.value.class) params.append('class', filters.value.class)
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
  return new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }).format(new Date(dateStr))
}

onMounted(() => {
  fetchData()
  fetchClasses()
})
</script>
