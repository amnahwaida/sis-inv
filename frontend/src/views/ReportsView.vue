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
          <button @click="exportTransactions" 
                  class="bg-white text-primary-900 hover:bg-primary-50 px-6 py-3 rounded-xl text-[10px] font-black transition-all flex items-center justify-center gap-2 shadow-xl active:scale-95 w-full">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
            EXPORT LAPORAN (EXCEL)
          </button>
        </div>
      </div>
    </div>

    <!-- Tabs & Filter Section -->
    <div class="flex flex-col md:flex-row items-center justify-between gap-6">
      <div class="flex flex-wrap items-center justify-center gap-2 bg-white dark:bg-gray-800 p-2 rounded-[1.5rem] border border-gray-100 dark:border-gray-700 shadow-sm w-full md:w-fit">
        <button v-for="tab in [{id:'active', label:'Peminjaman Aktif', icon:'⚡'}, {id:'overdue', label:'Terlambat', icon:'⚠️'}, {id:'history', label:'Riwayat Lengkap', icon:'📜'}, {id:'students', label:'Daftar Siswa', icon:'🎓'}]" 
                :key="tab.id"
                @click="activeTab = tab.id" 
                :class="activeTab === tab.id ? 'bg-primary-600 text-white shadow-lg' : 'text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'" 
                class="px-5 py-2.5 rounded-xl text-xs font-black transition-all duration-300 flex items-center gap-2 whitespace-nowrap">
          <span>{{ tab.icon }}</span>
          <span>{{ tab.label }}</span>
          <span v-if="tab.id === 'active' && activeBorrowings.length" class="ml-1 px-2 py-0.5 bg-white/20 rounded-full text-[10px]">{{ activeBorrowings.length }}</span>
          <span v-if="tab.id === 'overdue' && overdueItems.length" class="ml-1 px-2 py-0.5 bg-white/20 rounded-full text-[10px]">{{ overdueItems.length }}</span>
        </button>
      </div>

      <!-- Class Filter Overlay -->
      <div v-if="activeTab === 'active'" class="flex items-center justify-between gap-3 bg-white dark:bg-gray-800 p-2 px-4 rounded-[1.5rem] border border-gray-100 dark:border-gray-700 shadow-sm w-full md:w-auto animate-fade-in group">
        <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest pl-1">Filter Kelas:</label>
        <select v-model="classFilter" @change="fetchActive" class="bg-transparent border-none focus:ring-0 text-sm font-black text-gray-900 dark:text-white py-1 w-full md:w-auto text-right md:text-left">
          <option value="">Semua Kelas</option>
          <option v-for="cls in uniqueClasses" :key="cls" :value="cls">{{ cls }}</option>
        </select>
      </div>
    </div>

    <!-- Tables Container with Ultra Rounded style -->
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
            <tr v-if="activeBorrowings.length === 0" class="text-center">
              <td colspan="6" class="px-8 py-24">
                <div class="w-20 h-20 bg-emerald-50 dark:bg-emerald-900/20 rounded-full flex items-center justify-center mx-auto mb-6 text-emerald-500">
                  <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                </div>
                <p class="text-xl font-black text-gray-900 dark:text-white uppercase tracking-tight">Semua Aset Tersedia</p>
                <p class="text-xs text-gray-400 font-medium mt-1">Tidak ada peminjaman aktif saat ini.</p>
              </td>
            </tr>
            <tr v-for="trx in activeBorrowings" :key="trx.id" class="hover:bg-primary-50/50 dark:hover:bg-primary-900/10 transition-all duration-300">
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
              <td class="px-8 py-6 text-[10px] font-black text-gray-400 uppercase tracking-widest">{{ trx.teacher_name }}</td>
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

        <!-- Other tables (Overdue, History) would follow the same structure -->
        <table v-if="activeTab === 'overdue' || activeTab === 'history'" class="w-full">
          <thead>
            <tr class="bg-gray-50/50 dark:bg-gray-700/30">
              <th v-for="h in ['Barang', 'Peminjam', 'Tanggal Pinjam', activeTab === 'overdue' ? 'Jatuh Tempo' : 'Tgl Kembali', 'Status']" :key="h"
                  class="text-left py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] whitespace-nowrap">{{ h }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50 dark:divide-gray-700">
            <template v-if="(activeTab === 'overdue' ? overdueItems : histories).length > 0">
              <tr v-for="trx in (activeTab === 'overdue' ? overdueItems : histories)" :key="trx.id" class="hover:bg-primary-50/50 dark:hover:bg-primary-900/10 transition-all duration-300">
                <td class="px-8 py-6">
                  <div class="font-black text-gray-900 dark:text-white text-sm tracking-tight">{{ trx.item_name }}</div>
                  <div class="text-[10px] text-gray-400 font-mono tracking-tighter">{{ trx.item_code }}</div>
                </td>
                <td class="px-8 py-6">
                  <div class="text-xs font-black text-gray-900 dark:text-white">{{ trx.student_name || trx.teacher_name }}</div>
                  <div class="text-[10px] text-gray-400">{{ trx.student_class || 'Staff' }}</div>
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
              <td colspan="5" class="px-8 py-24 italic text-gray-400 font-bold tracking-[0.2em] text-[10px] uppercase">Tidak Ada Data Ditemukan</td>
            </tr>
          </tbody>
        </table>

        <!-- Students Table -->
        <table v-if="activeTab === 'students'" class="w-full">
          <thead>
            <tr class="bg-gray-50/50 dark:bg-gray-700/30">
              <th class="text-left py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">Siswa</th>
              <th class="text-left py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">NIS</th>
              <th class="text-left py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">Kelas</th>
              <th class="text-right py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">
                <button @click="exportStudents" class="bg-emerald-500 text-white px-4 py-1.5 rounded-lg text-[9px] font-black hover:bg-emerald-600 transition-all uppercase tracking-widest shadow-lg shadow-emerald-500/20">
                  EXCEL
                </button>
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50 dark:divide-gray-700">
            <tr v-for="s in studentsList" :key="s.id" class="group hover:bg-primary-50/30 dark:hover:bg-primary-900/10 transition-colors">
              <td class="px-8 py-5">
                <span class="text-sm font-black text-gray-900 dark:text-white uppercase">{{ s.full_name }}</span>
              </td>
              <td class="px-8 py-5 text-sm font-bold text-gray-500 dark:text-gray-400">{{ s.nis }}</td>
              <td class="px-8 py-5">
                 <span class="px-2 py-1 bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 rounded text-[10px] font-black">{{ s.class }}</span>
              </td>
              <td class="px-8 py-5 text-right"></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import api from '../utils/api'

const activeTab = ref('active')
const activeBorrowings = ref([])
const overdueItems = ref([])
const histories = ref([])
const studentsList = ref([])
const uniqueClasses = ref([])
const classFilter = ref('')
const loading = ref(false)

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
    const { data } = await api.get('/reports/overdue')
    if (data.success) { overdueItems.value = data.data }
  } catch (e) { console.error(e) }
}

async function fetchHistory() {
  loading.value = true
  try {
    const { data } = await api.get('/reports/history')
    if (data.success) { histories.value = data.data || [] }
  } catch (e) { console.error(e) } 
  finally { loading.value = false }
}

async function fetchStudents() {
  loading.value = true
  try {
    const { data } = await api.get('/students')
    if (data.success) { 
      studentsList.value = Array.isArray(data.data) ? data.data : (data.data.items || [])
    }
  } catch (e) { console.error(e) } 
  finally { loading.value = false }
}

async function fetchClasses() {
  try {
    const { data } = await api.get('/reports/active-borrowings')
    if (data.success) {
      const classes = [...new Set(data.data.filter(t => t.student_class).map(t => t.student_class))]
      uniqueClasses.value = classes.sort()
    }
  } catch (e) { /* ignore */ }
}

async function exportTransactions() {
  try {
    const response = await api.get('/reports/export/transactions', { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url; link.setAttribute('download', `SIS-INV_Laporan_${new Date().toISOString().slice(0,10)}.xlsx`)
    document.body.appendChild(link); link.click(); document.body.removeChild(link)
  } catch (e) { alert('Gagal mendownload laporan') }
}

async function exportStudents() {
  try {
    const response = await api.get('/students/export', { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url; link.setAttribute('download', `SIS-INV_Daftar_Siswa_${new Date().toISOString().slice(0,10)}.xlsx`)
    document.body.appendChild(link); link.click(); document.body.removeChild(link)
  } catch (e) { alert('Gagal mendownload data siswa') }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }).format(new Date(dateStr))
}

watch(activeTab, (newTab) => {
  if (newTab === 'active') fetchActive()
  else if (newTab === 'overdue') fetchOverdue()
  else if (newTab === 'history') fetchHistory()
  else if (newTab === 'students') fetchStudents()
})

onMounted(() => {
  fetchActive(); fetchOverdue(); fetchHistory(); fetchClasses(); fetchStudents()
})
</script>
