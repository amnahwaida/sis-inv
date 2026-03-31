<template>
  <div class="animate-fade-in space-y-6 max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-gray-900 tracking-tight">Laporan & Monitoring</h1>
        <p class="text-sm text-gray-500 mt-1">Peminjaman aktif, terlambat, dan riwayat transaksi.</p>
      </div>
      <div class="flex items-center gap-2">
        <button @click="exportTransactions" class="px-4 py-2 bg-white border border-gray-200 text-gray-700 font-bold rounded-xl hover:bg-gray-50 transition-all flex items-center gap-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/></svg>
          Export Transaksi (Excel)
        </button>
      </div>
    </div>

    <!-- Tabs -->
    <div class="flex items-center gap-1 bg-gray-100 p-1 rounded-xl w-fit">
      <button @click="activeTab = 'active'" :class="activeTab === 'active' ? 'bg-white shadow text-gray-900' : 'text-gray-500 hover:text-gray-700'" class="px-4 py-2 rounded-lg text-sm font-bold transition-all">
        Sedang Dipinjam
        <span v-if="activeBorrowings.length > 0" class="ml-1 px-1.5 py-0.5 bg-blue-100 text-blue-700 rounded-full text-[10px] font-black">{{ activeBorrowings.length }}</span>
      </button>
      <button @click="activeTab = 'overdue'" :class="activeTab === 'overdue' ? 'bg-white shadow text-gray-900' : 'text-gray-500 hover:text-gray-700'" class="px-4 py-2 rounded-lg text-sm font-bold transition-all">
        Terlambat
        <span v-if="overdueItems.length > 0" class="ml-1 px-1.5 py-0.5 bg-red-100 text-red-700 rounded-full text-[10px] font-black">{{ overdueItems.length }}</span>
      </button>
      <button @click="activeTab = 'history'" :class="activeTab === 'history' ? 'bg-white shadow text-gray-900' : 'text-gray-500 hover:text-gray-700'" class="px-4 py-2 rounded-lg text-sm font-bold transition-all">
        Semua Riwayat
      </button>
    </div>

    <!-- Class Filter (only for active tab) -->
    <div v-if="activeTab === 'active'" class="flex items-center gap-3 animate-fade-in">
      <label class="text-xs font-black text-gray-500 uppercase tracking-widest">Filter Kelas:</label>
      <select v-model="classFilter" @change="fetchActive" class="input-field rounded-xl border-gray-200 w-48 text-sm">
        <option value="">Semua Kelas</option>
        <option v-for="cls in uniqueClasses" :key="cls" :value="cls">{{ cls }}</option>
      </select>
    </div>

    <!-- Tables Container -->
    <div class="card p-0 overflow-hidden border-gray-100 shadow-sm bg-white rounded-3xl animate-scale-up">
      <div class="overflow-x-auto">
        <!-- Active Borrowings Table -->
        <table v-if="activeTab === 'active'" class="w-full">
          <thead class="bg-gray-50/50 border-b border-gray-100">
            <tr>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Barang</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Peminjam</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Kelas</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Guru Proses</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Waktu</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-if="loading" class="text-center">
              <td colspan="6" class="px-6 py-12 text-gray-400 italic">Memuat data...</td>
            </tr>
            <tr v-else-if="activeBorrowings.length === 0" class="text-center">
              <td colspan="6" class="px-6 py-12 text-gray-400">
                <div class="text-3xl mb-2">🎉</div>
                <p class="font-bold text-gray-900">Tidak ada peminjaman aktif</p>
              </td>
            </tr>
            <tr v-for="trx in activeBorrowings" :key="trx.id" class="hover:bg-gray-50/50 transition-colors">
              <td class="px-6 py-4">
                <div class="font-bold text-gray-900">{{ trx.item_name }}</div>
                <div class="text-xs text-gray-400 font-mono">{{ trx.item_code }}</div>
              </td>
              <td class="px-6 py-4">
                <div class="text-sm font-bold text-gray-900">{{ trx.student_name || trx.teacher_name }}</div>
                <span class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-md" :class="trx.borrower_type === 'STUDENT' ? 'bg-indigo-50 text-indigo-600' : 'bg-gray-100 text-gray-600'">
                  {{ trx.borrower_type === 'STUDENT' ? 'Siswa' : 'Staff' }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ trx.student_class || '-' }}</td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ trx.teacher_name }}</td>
              <td class="px-6 py-4">
                <div class="text-sm text-gray-700">{{ formatDate(trx.borrowed_at) }}</div>
                <div class="text-[10px] font-bold mt-0.5" :class="trx.is_overdue ? 'text-red-500' : 'text-gray-400'">
                  Jatuh Tempo: {{ formatDate(trx.due_date) }}
                </div>
              </td>
              <td class="px-6 py-4">
                <span v-if="trx.is_overdue" class="px-3 py-1 bg-red-100 text-red-700 rounded-lg text-[10px] font-black uppercase tracking-wider">Terlambat</span>
                <span v-else class="px-3 py-1 bg-blue-100 text-blue-700 rounded-lg text-[10px] font-black uppercase tracking-wider">Aktif</span>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- Overdue Table -->
        <table v-if="activeTab === 'overdue'" class="w-full">
          <thead class="bg-red-50/50 border-b border-red-100">
            <tr>
              <th class="text-left text-[10px] font-black text-red-400 uppercase tracking-widest px-6 py-4">Barang</th>
              <th class="text-left text-[10px] font-black text-red-400 uppercase tracking-widest px-6 py-4">Peminjam</th>
              <th class="text-left text-[10px] font-black text-red-400 uppercase tracking-widest px-6 py-4">Kelas</th>
              <th class="text-left text-[10px] font-black text-red-400 uppercase tracking-widest px-6 py-4">Guru</th>
              <th class="text-left text-[10px] font-black text-red-400 uppercase tracking-widest px-6 py-4">Jatuh Tempo</th>
              <th class="text-left text-[10px] font-black text-red-400 uppercase tracking-widest px-6 py-4">Tujuan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-red-50">
            <tr v-if="overdueItems.length === 0" class="text-center">
              <td colspan="6" class="px-6 py-12 text-gray-400 italic">Tidak ada peminjaman terlambat</td>
            </tr>
            <tr v-for="trx in overdueItems" :key="trx.id" class="hover:bg-red-50/30 transition-colors">
              <td class="px-6 py-4">
                <div class="font-bold text-gray-900">{{ trx.item_name }}</div>
                <div class="text-[10px] text-gray-400 font-mono">{{ trx.item_code }}</div>
              </td>
              <td class="px-6 py-4 text-sm font-bold text-gray-900">{{ trx.student_name || trx.teacher_name }}</td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ trx.student_class || '-' }}</td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ trx.teacher_name }}</td>
              <td class="px-6 py-4">
                <div class="text-sm font-bold text-red-600">{{ formatDate(trx.due_date) }}</div>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ trx.purpose || '-' }}</td>
            </tr>
          </tbody>
        </table>

        <!-- History Table -->
        <table v-if="activeTab === 'history'" class="w-full">
          <thead class="bg-indigo-50/50 border-b border-indigo-100">
            <tr>
              <th class="text-left text-[10px] font-black text-indigo-400 uppercase tracking-widest px-6 py-4">Barang</th>
              <th class="text-left text-[10px] font-black text-indigo-400 uppercase tracking-widest px-6 py-4">Peminjam</th>
              <th class="text-left text-[10px] font-black text-indigo-400 uppercase tracking-widest px-6 py-4">Tgl Pinjam</th>
              <th class="text-left text-[10px] font-black text-indigo-400 uppercase tracking-widest px-6 py-4">Tgl Kembali</th>
              <th class="text-left text-[10px] font-black text-indigo-400 uppercase tracking-widest px-6 py-4">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-if="histories.length === 0" class="text-center">
              <td colspan="5" class="px-6 py-12 text-gray-400 italic">Belum ada riwayat transaksi</td>
            </tr>
            <tr v-for="(trx, idx) in histories" :key="idx" class="hover:bg-indigo-50/10 transition-colors">
              <td class="px-6 py-4">
                <div class="font-bold text-gray-900 text-sm">{{ trx.item_name }}</div>
                <div class="text-[10px] text-gray-400 font-mono">{{ trx.item_code }}</div>
              </td>
              <td class="px-6 py-4">
                <div class="text-xs font-bold text-gray-900">{{ trx.student_name || trx.teacher_name }}</div>
                <div v-if="trx.student_class" class="text-[10px] text-gray-400">{{ trx.student_class }}</div>
              </td>
              <td class="px-6 py-4 text-xs text-gray-500">{{ formatDate(trx.borrowed_at) }}</td>
              <td class="px-6 py-4 text-xs text-gray-500">{{ formatDate(trx.returned_at) }}</td>
              <td class="px-6 py-4">
                <span class="px-2 py-0.5 rounded text-[10px] font-black uppercase tracking-widest" :class="trx.status === 'RETURNED' ? 'bg-green-100 text-green-700' : 'bg-blue-100 text-blue-700'">
                  {{ trx.status }}
                </span>
              </td>
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
const uniqueClasses = ref([])
const classFilter = ref('')
const loading = ref(false)

async function fetchActive() {
  loading.value = true
  try {
    const params = classFilter.value ? `?class=${encodeURIComponent(classFilter.value)}` : ''
    const { data } = await api.get(`/reports/active-borrowings${params}`)
    if (data.success) {
      activeBorrowings.value = data.data
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function fetchOverdue() {
  try {
    const { data } = await api.get('/reports/overdue')
    if (data.success) {
      overdueItems.value = data.data
    }
  } catch (e) {
    console.error(e)
  }
}

async function fetchHistory() {
  loading.value = true
  try {
    const { data } = await api.get('/reports/history')
    if (data.success) {
      histories.value = data.data
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
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
    link.href = url
    link.setAttribute('download', `SIS-INV_Laporan-Transaksi_${new Date().toISOString().slice(0,10)}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (e) {
    alert('Gagal mendownload laporan')
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Intl.DateTimeFormat('id-ID', {
    day: '2-digit', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  }).format(new Date(dateStr))
}

watch(activeTab, (newTab) => {
  if (newTab === 'active') fetchActive()
  if (newTab === 'overdue') fetchOverdue()
  if (newTab === 'history') fetchHistory()
})

onMounted(() => {
  fetchActive()
  fetchOverdue()
  fetchHistory()
  fetchClasses()
})
</script>
