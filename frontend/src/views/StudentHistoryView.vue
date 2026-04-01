<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Riwayat Siswa</h1>
          <p class="text-primary-100/70 text-sm font-medium">Lacak keaktifan peminjaman aset oleh siswa</p>
        </div>
        
        <div class="flex flex-col sm:flex-row items-stretch gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10 w-full lg:w-auto">
          <div class="relative group flex-1">
            <input type="text" v-model="nisnSearch" @input="debouncedStudentSearch" @keydown.enter="handleSearch"
                   class="bg-white/10 border-none text-white placeholder-white/40 text-xs font-bold rounded-xl h-12 pl-12 pr-4 w-full focus:bg-white/20 focus:ring-0 transition-all" 
                   placeholder="Ketik NIS atau Nama Siswa..." />
            <svg class="w-5 h-5 absolute left-4 top-1/2 -translate-y-1/2 text-white/40" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
            
            <!-- Auto-suggestion Dropdown -->
            <div v-if="suggestions.length > 0" class="absolute z-50 left-0 right-0 top-full mt-2 bg-white dark:bg-gray-800 rounded-2xl border border-gray-200 dark:border-gray-700 shadow-2xl overflow-hidden animate-scale-up">
              <div v-for="s in suggestions" :key="s.nis" @click="selectSuggestion(s)"
                   class="p-4 hover:bg-primary-50 dark:hover:bg-primary-900/20 cursor-pointer border-b border-gray-50 dark:border-gray-700 last:border-0 transition-colors">
                <div class="font-black text-gray-900 dark:text-white text-sm">{{ s.full_name }}</div>
                <div class="text-[10px] font-bold text-gray-400 uppercase tracking-widest">NIS: {{ s.nis }} • {{ s.class }}</div>
              </div>
            </div>
          </div>
          <button @click="handleSearch" 
                  class="bg-white text-primary-900 hover:bg-primary-50 px-6 py-3 rounded-xl text-[10px] font-black transition-all shadow-xl active:scale-95 whitespace-nowrap">
            CARI SISWA
          </button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="py-20 flex flex-col items-center justify-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      <span class="mt-4 text-xs font-black text-primary-600 uppercase tracking-widest animate-pulse">Menelusuri Database...</span>
    </div>

    <!-- Student Info & History Content -->
    <div v-else-if="studentInfo" class="space-y-10 animate-scale-up">
      <!-- Student Card -->
      <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] p-10 border border-gray-100 dark:border-gray-700 shadow-sm flex flex-col md:flex-row items-center gap-8">
        <div class="w-24 h-24 bg-gradient-to-br from-primary-500 to-indigo-600 text-white rounded-3xl flex items-center justify-center text-3xl font-black shadow-2xl shadow-primary-500/20 border-4 border-white dark:border-gray-700">
          {{ studentInfo.full_name?.charAt(0) }}
        </div>
        <div class="flex-1 text-center md:text-left space-y-2">
          <div class="flex flex-col md:flex-row md:items-center gap-2 md:gap-4">
            <h2 class="text-3xl font-black text-gray-900 dark:text-white uppercase tracking-tight">{{ studentInfo.full_name }}</h2>
            <span class="px-3 py-1 bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 text-[10px] font-black uppercase tracking-widest rounded-xl w-fit mx-auto md:mx-0">Aktif</span>
          </div>
          <div class="flex flex-center flex-wrap md:justify-start gap-x-6 gap-y-2 text-gray-400 font-bold text-xs uppercase tracking-widest">
            <span class="flex items-center gap-2"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10 6H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-5m-4 0V5a2 2 0 012-2h2a2 2 0 012 2v1m-4 0h4"/></svg> NIS: {{ studentInfo.nis }}</span>
            <span class="flex items-center gap-2"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5"/></svg> KELAS: {{ studentInfo.class }}</span>
          </div>
        </div>
        <div class="grid grid-cols-2 gap-4 w-full md:w-auto">
          <div class="p-6 bg-gray-50 dark:bg-gray-900/40 rounded-3xl text-center">
            <div class="text-2xl font-black text-primary-600">{{ histories.length }}</div>
            <div class="text-[9px] font-black text-gray-400 uppercase tracking-widest">Total Pinjam</div>
          </div>
          <div class="p-6 bg-gray-50 dark:bg-gray-900/40 rounded-3xl text-center">
            <div class="text-2xl font-black text-emerald-600">{{ histories.filter(h => h.returned_at).length }}</div>
            <div class="text-[9px] font-black text-gray-400 uppercase tracking-widest">Berhasil</div>
          </div>
        </div>
      </div>

      <!-- History Table -->
      <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden transition-all duration-300">
        <div class="p-8 border-b border-gray-100 dark:border-gray-700">
          <h3 class="text-lg font-black text-gray-900 dark:text-white uppercase tracking-tight">Timeline Transaksi</h3>
        </div>
        <div class="overflow-x-auto relative">
          <table class="w-full">
            <thead>
              <tr class="bg-gray-50/50 dark:bg-gray-700/30">
                <th v-for="h in ['Aset', 'Guru Piket', 'Tgl Pinjam', 'Tgl Kembali', 'Status']" :key="h"
                    class="text-left py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">{{ h }}</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-50 dark:divide-gray-700">
              <template v-if="histories.length > 0">
                <tr v-for="trx in histories" :key="trx.id" class="group hover:bg-primary-50/50 dark:hover:bg-primary-900/10 transition-all duration-300">
                  <td class="px-8 py-6">
                    <div class="font-black text-gray-900 dark:text-white text-sm tracking-tight leading-none mb-1">{{ trx.item_name }}</div>
                    <div class="text-[10px] text-gray-400 font-mono font-bold tracking-tighter">{{ trx.item_code }}</div>
                  </td>
                  <td class="px-8 py-6 text-xs font-black text-gray-600 dark:text-gray-400 uppercase tracking-widest">{{ trx.teacher_name }}</td>
                  <td class="px-8 py-6 text-xs font-medium text-gray-500">{{ formatDate(trx.borrowed_at) }}</td>
                  <td class="px-8 py-6 text-xs font-medium text-gray-500">{{ formatDate(trx.returned_at) }}</td>
                  <td class="px-8 py-6">
                    <span v-if="trx.status === 'RETURNED'" class="px-3 py-1 rounded-lg text-[9px] font-black uppercase tracking-widest shadow-sm inline-block bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600">
                      Dikembalikan
                    </span>
                    <span v-else class="text-[9px] font-black text-primary-600 animate-pulse uppercase tracking-[0.2em]">Masih Dipinjam</span>
                  </td>
                </tr>
              </template>
              <tr v-else class="text-center">
                <td colspan="5" class="px-8 py-16 text-gray-400 italic text-sm">Siswa ini belum pernah melakukan peminjaman.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Empty State / Search Prompt -->
    <div v-else class="py-24 text-center space-y-6">
      <div class="w-32 h-32 bg-gray-50 dark:bg-gray-800 rounded-full flex items-center justify-center mx-auto text-gray-300 transition-all group-hover:scale-110 shadow-inner">
        <svg class="w-16 h-16" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
      </div>
      <div class="space-y-2">
        <h2 class="text-2xl font-black text-gray-900 dark:text-white uppercase tracking-tight">Cek Riwayat Aset Siswa</h2>
        <p class="text-sm text-gray-400 max-w-sm mx-auto font-medium">Ketik NIS atau nama siswa untuk melihat seluruh riwayat peminjaman mereka secara lengkap.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '../utils/api'

const route = useRoute()
const nisnSearch = ref('')
const loading = ref(false)
const studentInfo = ref(null)
const histories = ref([])
const suggestions = ref([])

onMounted(() => {
  if (route.query.nisn) {
    nisnSearch.value = route.query.nisn
    handleSearch()
  }
})

// Auto-suggestion: search students by NIS or name
let searchTimeout
const debouncedStudentSearch = () => {
  clearTimeout(searchTimeout)
  suggestions.value = []
  if (!nisnSearch.value.trim() || nisnSearch.value.trim().length < 2) return
  searchTimeout = setTimeout(fetchSuggestions, 400)
}

const fetchSuggestions = async () => {
  try {
    const { data } = await api.get('/students/search', { params: { q: nisnSearch.value.trim() } })
    if (data.success) {
      suggestions.value = data.data || []
    }
  } catch (err) {
    suggestions.value = []
  }
}

const selectSuggestion = (student) => {
  // Store the student info for display, use NIS for the search
  studentInfo.value = student
  nisnSearch.value = student.nis
  suggestions.value = []
  fetchHistory(student.nis)
}

const handleSearch = async () => {
  if (!nisnSearch.value.trim()) return
  suggestions.value = []
  loading.value = true
  studentInfo.value = null
  histories.value = []

  const searchTerm = nisnSearch.value.trim()

  // First, try to find student info from students/search if we don't have it yet
  if (!studentInfo.value) {
    try {
      const { data } = await api.get('/students/search', { params: { q: searchTerm } })
      if (data.success && data.data && data.data.length > 0) {
        // Pick the best match (exact NIS match first, then first result)
        const exactMatch = data.data.find(s => s.nis === searchTerm)
        studentInfo.value = exactMatch || data.data[0]
      }
    } catch (err) { /* ignore, will still try history */ }
  }

  await fetchHistory(searchTerm)
}

const fetchHistory = async (nis) => {
  loading.value = true
  try {
    // Use the correct backend endpoint: /transactions/student/:nis
    const { data } = await api.get(`/transactions/student/${nis}`)
    if (data.success) {
      histories.value = data.data || []
      
      // If we have history data but no studentInfo yet, build it from the first transaction
      if (!studentInfo.value && histories.value.length > 0) {
        const first = histories.value[0]
        studentInfo.value = {
          full_name: first.student_name || nis,
          nis: nis,
          class: first.student_class || '-'
        }
      }
      
      // If still no studentInfo and no history, show an error
      if (!studentInfo.value) {
        studentInfo.value = { full_name: nis, nis: nis, class: '-' }
      }
    }
  } catch (err) {
    alert(err.response?.data?.error || 'Siswa tidak ditemukan atau belum memiliki riwayat')
    studentInfo.value = null
  } finally {
    loading.value = false
  }
}

const formatDate = (d) => {
  if (!d) return '-'
  return new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }).format(new Date(d))
}
</script>
