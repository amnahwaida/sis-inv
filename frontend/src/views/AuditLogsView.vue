<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Log Audit</h1>
          <p class="text-primary-100/70 text-sm font-medium">Rekaman jejak aktivitas and keamanan sistem</p>
        </div>
        
        <div class="flex flex-col sm:flex-row items-center gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10 w-full md:w-auto">
          <button @click="exportExcel" class="bg-white/10 hover:bg-white/20 text-white border border-white/10 px-5 py-2.5 rounded-xl text-[10px] font-black transition-all flex items-center justify-center gap-2 active:scale-95">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2-2z" /></svg>
            EXCEL
          </button>
          <button @click="fetchLogs" class="btn-premium-primary">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
            </svg>
            REFRESH LOG
          </button>
        </div>
      </div>
    </div>

    <!-- Search & Filters Container -->
    <div class="bg-white dark:bg-gray-800 p-3 rounded-[2rem] border border-gray-100 dark:border-gray-700 shadow-sm flex flex-col sm:flex-row items-center gap-3 w-full">
      <!-- Search Input -->
      <div class="relative w-full sm:w-80">
        <input type="text" v-model="searchQuery" placeholder="Cari aktor, detail objek..." 
               class="input-field pl-10 h-12 rounded-2xl text-sm w-full" />
        <svg class="w-5 h-5 absolute left-3 top-3.5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
        </svg>
      </div>

      <!-- Action Filter -->
      <select v-model="filterAction" 
              class="input-field h-12 rounded-2xl text-sm w-full sm:w-auto appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2020%2020%22%3E%3Cpath%20stroke%3D%22%236b7280%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%221.5%22%20d%3D%22m6%208%204%204%204-4%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[right_0.5rem_center] bg-no-repeat pr-10">
        <option value="">Semua Aksi</option>
        <option v-for="action in availableActions" :key="action" :value="action">
          {{ getCleanActionLabel(action) }}
        </option>
      </select>
    </div>

    <!-- Audit Logs Table Card -->
    <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden transition-all duration-300">
      <div class="p-8 border-b border-gray-100 dark:border-gray-700 bg-gray-50/30 dark:bg-gray-900/10 flex items-center justify-between">
        <h2 class="text-lg font-black text-gray-900 dark:text-white uppercase tracking-tight">Timeline Aktivitas</h2>
        <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest leading-none">Sinkronisasi Real-time</span>
      </div>

      <div class="overflow-x-auto relative">
        <!-- Loading Overlay -->
        <div v-if="loading" class="absolute inset-0 bg-white/60 dark:bg-gray-900/60 backdrop-blur-sm z-10 flex flex-col items-center justify-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        </div>

        <table class="w-full">
          <thead>
            <tr class="bg-gray-100/50 dark:bg-gray-700/30">
              <th v-for="h in ['Waktu / Tanggal', 'Aktor (User)', 'Aksi / Tindakan', 'Detail Objek', 'IP Address']" :key="h"
                  class="text-left py-5 px-8 text-[10px] font-black text-gray-400 dark:text-gray-300 uppercase tracking-[0.2em] whitespace-nowrap">{{ h }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 dark:divide-gray-700">
            <template v-if="logs.length > 0">
              <tr v-for="log in logs" :key="log.id" class="group hover:bg-primary-50/50 dark:hover:bg-primary-900/10 transition-all duration-300">
                <td class="px-8 py-6">
                  <div class="text-xs font-black text-gray-900 dark:text-white leading-none mb-1">{{ formatFullDate(log.created_at) }}</div>
                  <div class="text-[10px] text-gray-400 font-bold uppercase tracking-widest">{{ formatTime(log.created_at) }}</div>
                </td>
                <td class="px-8 py-6">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center text-[10px] font-black text-gray-500 uppercase">{{ log.user_name?.charAt(0) || 'S' }}</div>
                    <span class="text-xs font-bold text-gray-700 dark:text-gray-300">{{ log.user_name || 'SYSTEM' }}</span>
                  </div>
                </td>
                <td class="px-8 py-6">
                  <span class="px-3 py-1.5 rounded-xl text-[9px] font-black uppercase tracking-widest shadow-sm"
                        :class="getActionBadgeClass(log.action)">
                    {{ log.action }}
                  </span>
                </td>
                <td class="px-8 py-6">
                  <div class="max-w-sm">
                    <p class="text-xs text-gray-700 dark:text-gray-300 font-medium leading-relaxed" 
                       :class="{ 'line-clamp-2': !expandedLogs[log.id] }">
                      {{ log.description || '-' }}
                    </p>
                    <button v-if="log.description && log.description.length > 80" 
                            @click="toggleLog(log.id)"
                            class="mt-2 text-[10px] font-black uppercase tracking-widest text-primary-600 dark:text-primary-400 hover:text-primary-700 transition-colors flex items-center gap-1">
                      {{ expandedLogs[log.id] ? 'Sembunyikan' : 'Selengkapnya' }}
                      <svg class="w-3 h-3 transform transition-transform" :class="{ 'rotate-180': expandedLogs[log.id] }" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
                      </svg>
                    </button>
                  </div>
                </td>
                <td class="px-8 py-6">
                  <span class="text-[10px] font-bold font-mono text-gray-400 bg-gray-50 dark:bg-gray-900 px-2 py-1 rounded-md">{{ log.ip_address || '---' }}</span>
                </td>
              </tr>
            </template>
            <tr v-else-if="!loading" class="text-center">
              <td colspan="5" class="px-8 py-20 italic text-gray-400 font-bold uppercase tracking-widest text-xs">{{ searchQuery || filterAction ? 'Pencarian Tidak Ditemukan' : 'Belum ada rekaman aktivitas audit' }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination Bar -->
      <div class="px-8 py-5 bg-gray-50/50 dark:bg-gray-700/20 border-t border-gray-100 dark:border-gray-700 flex flex-col sm:flex-row items-center justify-between gap-4">
        <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">
          Menampilkan <span class="text-primary-600">{{ logs.length ? (meta.page - 1) * meta.page_size + 1 : 0 }}-{{ Math.min(meta.page * meta.page_size, meta.total) }}</span> dari <span class="text-gray-900 dark:text-white">{{ meta.total }}</span> log
        </span>
        <div class="flex gap-2">
          <button @click="currentPage--" :disabled="currentPage === 1" class="pagination-btn-standard">
            Kembali
          </button>
          <div class="flex gap-1">
            <button v-for="p in visiblePages" :key="p" @click="currentPage = p"
                    class="w-10 h-10 rounded-xl text-[11px] font-black transition-all shadow-sm active:scale-95 border"
                    :class="p === currentPage ? 'bg-primary-600 text-white border-primary-600' : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:bg-primary-50 dark:hover:bg-gray-700'">
              {{ p }}
            </button>
          </div>
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

const logs = ref([])
const availableActions = ref([])
const loading = ref(true)
const expandedLogs = ref({})

const toggleLog = (id) => {
  expandedLogs.value[id] = !expandedLogs.value[id]
}

const searchQuery = ref('')
const filterAction = ref('')
const currentPage = ref(1)
const perPage = 20
const meta = ref({
  total: 0,
  total_pages: 1,
  page: 1,
  page_size: 20
})

async function fetchLogs() {
  loading.value = true
  try {
    const { data } = await api.get('/reports/audit', {
      params: {
        page: currentPage.value,
        page_size: perPage,
        search: searchQuery.value,
        action: filterAction.value
      }
    })
    if (data.success) { 
      logs.value = data.data.items || []
      meta.value = data.data.meta || { total: 0, total_pages: 1, page: 1, page_size: 20 }
    }
  } catch (e) { 
    console.error('Gagal fetch audit logs:', e) 
    logs.value = []
  } finally { 
    loading.value = false 
  }
}

async function fetchActions() {
  try {
    const { data } = await api.get('/reports/audit-actions')
    if (data.success) {
      availableActions.value = data.data
    }
  } catch (e) {
    console.error('Gagal fetch audit actions:', e)
  }
}

const getCleanActionLabel = (action) => {
  const mapping = {
    'LOGIN': 'Login Pengguna',
    'LOGOUT': 'Logout Pengguna',
    'CREATE_ITEM': 'Tambah Barang',
    'UPDATE_ITEM': 'Ubah Barang',
    'DELETE_ITEM': 'Hapus Barang',
    'IMPORT_ITEMS': 'Import Barang (Excel)',
    'START_AUDIT': 'Mulai Opname',
    'CLOSE_AUDIT': 'Selesai Opname',
    'BORROW_ITEM': 'Peminjaman',
    'RETURN_ITEM': 'Pengembalian',
    'CREATE_USER': 'Tambah User',
    'UPDATE_USER': 'Ubah User',
    'DELETE_USER': 'Hapus User',
    'LOGIN_FAILED': 'Gagal Login'
  }
  return mapping[action] || action
}

// Debounced search
let searchTimeout
watch([searchQuery, filterAction], () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
    fetchLogs()
  }, 400)
})

watch(currentPage, fetchLogs)

const totalPages = computed(() => meta.value.total_pages)
const startRow = computed(() => (currentPage.value - 1) * perPage + 1)
const endRow = computed(() => Math.min(currentPage.value * perPage, meta.value.total))

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

async function exportExcel() {
  try {
    const response = await api.get('/reports/export/audit', { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `SIS-INV_Audit-Log_${new Date().toISOString().slice(0,10)}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (e) { alert('Gagal mendownload audit log') }
}

const getActionBadgeClass = (a) => {
  if (a.includes('LOGIN')) return 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600'
  if (a.includes('CREATE')) return 'bg-blue-50 dark:bg-blue-900/30 text-blue-600'
  if (a.includes('UPDATE')) return 'bg-amber-50 dark:bg-amber-900/30 text-amber-600'
  if (a.includes('DELETE')) return 'bg-rose-50 dark:bg-rose-900/30 text-rose-600'
  return 'bg-gray-100 text-gray-600'
}

const formatFullDate = (d) => d ? new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'long', year: 'numeric' }).format(new Date(d)) : '-'
const formatTime = (d) => d ? new Intl.DateTimeFormat('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' }).format(new Date(d)) : ''

onMounted(() => {
  fetchLogs()
  fetchActions()
})
</script>
