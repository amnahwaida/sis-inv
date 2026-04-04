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
            <svg class="w-5 h-5 absolute left-4 top-1/2 -translate-y-1/2 text-white/40" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
            </svg>
            
            <!-- Auto-suggestion Dropdown -->
            <div v-if="suggestions.length > 0" class="absolute z-50 left-0 right-0 top-full mt-2 bg-white dark:bg-gray-800 rounded-2xl border border-gray-200 dark:border-gray-700 shadow-2xl overflow-hidden animate-scale-up">
              <div v-for="s in suggestions" :key="s.nis" @click="selectSuggestion(s)"
                   class="p-4 hover:bg-primary-50 dark:hover:bg-primary-900/20 cursor-pointer border-b border-gray-50 dark:border-gray-700 last:border-0 transition-colors">
                <div class="font-black text-gray-900 dark:text-white text-sm">{{ s.full_name }}</div>
                <div class="text-[10px] font-bold text-gray-400 uppercase tracking-widest">NIS: {{ s.nis }} • {{ s.class }}</div>
              </div>
            </div>
          </div>
          <button @click="handleSearch" class="btn-premium-primary whitespace-nowrap">
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
      <div class="card-premium group flex flex-col md:flex-row items-center gap-8">
        <div class="card-decoration"></div>
        <div class="w-24 h-24 bg-gradient-to-br from-primary-500 to-indigo-600 text-white rounded-3xl flex items-center justify-center text-3xl font-black shadow-2xl shadow-primary-500/20 border-4 border-white dark:border-gray-700">
          {{ studentInfo.full_name?.charAt(0) }}
        </div>
        <div class="flex-1 text-center md:text-left space-y-2">
          <div class="flex flex-col md:flex-row md:items-center gap-2 md:gap-4">
            <h2 class="text-3xl font-black text-gray-900 dark:text-white uppercase tracking-tight">{{ studentInfo.full_name }}</h2>
            <span v-if="studentInfo.is_active !== false" class="px-3 py-1 bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 text-[10px] font-black uppercase tracking-widest rounded-xl w-fit mx-auto md:mx-0">Aktif</span>
            <span v-else class="px-3 py-1 bg-red-50 dark:bg-red-900/30 text-red-500 text-[10px] font-black uppercase tracking-widest rounded-xl w-fit mx-auto md:mx-0">Non-Aktif</span>
          </div>
          <div class="flex items-center justify-center md:justify-start gap-x-6 gap-y-2 text-gray-400 font-bold text-xs uppercase tracking-widest flex-wrap">
            <span class="flex items-center gap-2">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 9h3.75M15 12h3.75M15 15h3.75M4.5 19.5h15a2.25 2.25 0 002.25-2.25V6.75A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25v10.5A2.25 2.25 0 004.5 19.5zm6-10.125a1.875 1.875 0 11-3.75 0 1.875 1.875 0 013.75 0zm1.294 6.336a6.721 6.721 0 01-3.17.789 6.721 6.721 0 01-3.168-.789 3.376 3.376 0 016.338 0z" />
              </svg> NIS: {{ studentInfo.nis }}
            </span>
            <span class="flex items-center gap-2">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4.26 10.147a60.436 60.436 0 00-.491 6.347A48.627 48.627 0 0112 20.904a48.627 48.627 0 018.232-4.41 60.46 60.46 0 00-.491-6.347m-15.482 0a50.57 50.57 0 00-2.658-.813A59.905 59.905 0 0112 3.493a59.902 59.902 0 0110.399 5.84c-.896.248-1.783.52-2.658.814m-15.482 0A50.697 50.697 0 0112 13.489a50.702 50.702 0 017.74-3.342M6.75 15a.75.75 0 100-1.5.75.75 0 000 1.5zm0 0v-3.675A55.378 55.378 0 0112 8.443m-7.007 11.55A5.981 5.981 0 006.75 15.75v-1.5" />
              </svg> KELAS: {{ studentInfo.class }}
            </span>
          </div>
        </div>
        <div class="flex flex-col items-center gap-4 w-full md:w-auto">
          <div class="grid grid-cols-2 gap-4 w-full">
            <div class="p-4 bg-gray-50 dark:bg-gray-900/40 rounded-3xl text-center">
              <div class="text-2xl font-black text-primary-600">{{ meta.total }}</div>
              <div class="text-[9px] font-black text-gray-400 uppercase tracking-widest">Total Pinjam</div>
            </div>
            <div class="p-4 bg-gray-50 dark:bg-gray-900/40 rounded-3xl text-center">
              <div class="text-2xl font-black text-emerald-600">{{ (histories || []).filter(h => h.returned_at).length }}</div>
              <div class="text-[9px] font-black text-gray-400 uppercase tracking-widest">Berhasil</div>
            </div>
          </div>
          <!-- EXPORT BUTTON -->
          <button @click="exportToExcel" class="btn-premium-action w-full bg-emerald-500 shadow-emerald-500/20">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3" />
            </svg>
            EKSPOR KE EXCEL
          </button>
        </div>
      </div>

      <!-- History Timeline Area -->
      <div class="space-y-6">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <h3 class="text-lg font-black text-gray-900 dark:text-white uppercase tracking-tight">Timeline Transaksi</h3>
          
          <!-- Inner Search -->
          <div class="relative w-full md:w-72 group">
            <input type="text" v-model="innerSearch" @input="debouncedInnerSearch"
                   class="input-field pl-10 h-11 rounded-2xl text-xs bg-white dark:bg-gray-800 border-gray-100 dark:border-gray-700 shadow-sm" 
                   placeholder="Cari aset dalam riwayat..." />
            <svg class="w-4 h-4 absolute left-3.5 top-3.5 text-gray-400 group-focus-within:text-primary-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
            </svg>
          </div>
        </div>

        <div class="relative">
          <!-- Loading Overlay -->
          <div v-if="loadingHistory" class="absolute inset-0 bg-white/60 dark:bg-gray-900/60 backdrop-blur-sm z-10 flex flex-col items-center justify-center rounded-3xl">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
          </div>

          <!-- Vertical Timeline Layout -->
          <div v-if="histories.length > 0" class="relative pl-8 space-y-8 before:absolute before:inset-0 before:left-[11px] before:w-0.5 before:bg-gray-100 dark:before:bg-gray-700">
            <div v-for="trx in histories" :key="trx.id" class="relative group animate-fade-in">
              <!-- Timeline Dot -->
              <div class="absolute -left-[31px] top-6 w-4 h-4 rounded-full border-4 border-white dark:border-gray-900 shadow-sm z-10 transition-transform group-hover:scale-125"
                   :class="trx.status === 'RETURNED' ? 'bg-emerald-500' : 'bg-primary-500 animate-pulse'"></div>
              
              <!-- Timeline Card -->
              <div class="card-premium !p-6 !rounded-[2rem] hover:ring-2 hover:ring-primary-500/20 transition-all">
                <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-6">
                  <div class="flex items-center gap-5">
                    <div class="w-12 h-12 bg-gray-50 dark:bg-gray-900 rounded-2xl flex items-center justify-center text-gray-400 group-hover:text-primary-500 transition-colors">
                      <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-10.5v10.5" />
                      </svg>
                    </div>
                    <div>
                      <div class="text-[10px] font-black text-primary-500 uppercase tracking-widest mb-1">{{ formatDateOnly(trx.borrowed_at) }}</div>
                      <h4 class="text-base font-black text-gray-900 dark:text-white uppercase tracking-tight">{{ trx.item_name }}</h4>
                      <p class="text-[10px] font-mono font-bold text-gray-400 tracking-tighter">{{ trx.item_code }} • PETUGAS: {{ trx.teacher_name }}</p>
                    </div>
                  </div>

                  <div class="flex flex-wrap items-center gap-4 border-t lg:border-t-0 pt-4 lg:pt-0 border-gray-50 dark:border-gray-700">
                    <div class="space-y-1">
                      <p class="text-[9px] font-black text-gray-400 uppercase tracking-widest">Peminjaman</p>
                      <p class="text-xs font-bold text-gray-700 dark:text-gray-300">{{ formatTimeOnly(trx.borrowed_at) }}</p>
                    </div>
                    <div class="w-px h-8 bg-gray-100 dark:bg-gray-700 hidden sm:block"></div>
                    <div class="space-y-1">
                      <p class="text-[9px] font-black text-gray-400 uppercase tracking-widest">Pengembalian</p>
                      <p class="text-xs font-bold" :class="trx.returned_at ? 'text-gray-700 dark:text-gray-300' : 'text-primary-500 italic'">
                        {{ trx.returned_at ? formatTimeOnly(trx.returned_at) : 'Siswa meminjam aset' }}
                      </p>
                    </div>
                    <div class="lg:ml-4">
                      <span v-if="trx.status === 'RETURNED'" class="px-4 py-1.5 bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 rounded-xl text-[9px] font-black uppercase tracking-widest border border-emerald-100 dark:border-emerald-800">Dikembalikan</span>
                      <span v-else class="px-4 py-1.5 bg-primary-50 dark:bg-primary-900/30 text-primary-600 rounded-xl text-[9px] font-black uppercase tracking-widest border border-primary-100 dark:border-primary-800 animate-pulse">Dipinjam</span>
                    </div>
                  </div>
                </div>
                
                <div v-if="trx.borrow_photo_url" class="mt-4 pt-4 border-t border-gray-50 dark:border-gray-700">
                   <a :href="trx.borrow_photo_url" target="_blank" class="inline-flex items-center gap-2 text-[10px] font-black text-primary-600 hover:underline">
                     <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                       <path stroke-linecap="round" stroke-linejoin="round" d="M6.827 6.175A2.31 2.31 0 015.186 7.23c-.38.054-.757.112-1.134.175C2.999 7.58 2.25 8.507 2.25 9.574V18a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9.574c0-1.067-.75-1.994-1.802-2.169a47.865 47.865 0 00-1.134-.175 2.31 2.31 0 01-1.64-1.055l-.822-1.316a2.192 2.192 0 00-1.736-1.039 48.774 48.774 0 00-5.232 0 2.192 2.192 0 00-1.736 1.039l-.821 1.316z" />
                       <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 12.75a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0zM18.75 10.5h.008v.008h-.008V10.5z" />
                     </svg>
                     LIHAT BUKTI FOTO PINJAMAN
                   </a>
                </div>
              </div>
            </div>
          </div>

          <!-- Empty History -->
          <div v-else class="py-20 bg-gray-50/50 dark:bg-gray-800/30 rounded-[3rem] border border-dashed border-gray-200 dark:border-gray-700 text-center space-y-4">
            <div class="w-16 h-16 bg-white dark:bg-gray-800 rounded-2xl shadow-sm mx-auto flex items-center justify-center text-gray-300">
              <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <p class="text-xs font-black text-gray-400 uppercase tracking-widest">{{ innerSearch ? 'Pencarian tidak ditemukan' : 'Belum ada riwayat tercatat' }}</p>
          </div>
        </div>

        <!-- Pagination Bar -->
        <div class="px-8 py-6 bg-white dark:bg-gray-800 rounded-[2rem] border border-gray-100 dark:border-gray-700 shadow-sm flex flex-col sm:flex-row items-center justify-between gap-4">
          <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">
            Menampilkan <span class="text-primary-600">{{ histories.length ? (meta.page - 1) * meta.page_size + 1 : 0 }}-{{ Math.min(meta.page * meta.page_size, meta.total) }}</span> dari <span class="text-gray-900 dark:text-white">{{ meta.total }}</span> data
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

    <!-- Empty State / Search Prompt -->
    <div v-else class="py-24 text-center space-y-6">
      <div class="w-32 h-32 bg-gray-50 dark:bg-gray-800 rounded-full flex items-center justify-center mx-auto text-gray-300 transition-all group-hover:scale-110 shadow-inner">
        <svg class="w-16 h-16" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
        </svg>
      </div>
      <div class="space-y-2">
        <h2 class="text-2xl font-black text-gray-900 dark:text-white uppercase tracking-tight">Cek Riwayat Aset Siswa</h2>
        <p class="text-sm text-gray-400 max-w-sm mx-auto font-medium">Ketik NIS atau nama siswa untuk melihat seluruh riwayat peminjaman mereka secara lengkap.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import api from '../utils/api'

const route = useRoute()
const nisnSearch = ref('')
const loading = ref(false)
const loadingHistory = ref(false)
const innerSearch = ref('')
const studentInfo = ref(null)
const histories = ref([])
const suggestions = ref([])
const meta = ref({
  total: 0,
  total_pages: 1,
  page: 1,
  page_size: 10
})

// Pagination
const currentPage = ref(1)
const perPage = 10

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

onMounted(() => {
  if (route.query.nisn) {
    nisnSearch.value = route.query.nisn
    handleSearch()
  }
})

// Auto-suggestion for student search
let searchTimeout
const debouncedStudentSearch = () => {
  clearTimeout(searchTimeout)
  suggestions.value = []
  if (!nisnSearch.value.trim() || nisnSearch.value.trim().length < 2) return
  searchTimeout = setTimeout(fetchSuggestions, 400)
}

// Inner search for transactions in timeline
let innerSearchTimeout
const debouncedInnerSearch = () => {
  clearTimeout(innerSearchTimeout)
  innerSearchTimeout = setTimeout(() => {
    currentPage.value = 1
    if (studentInfo.value) fetchHistory(studentInfo.value.nis)
  }, 400)
}

const fetchSuggestions = async () => {
  try {
    const { data } = await api.get('/students/search', { params: { q: nisnSearch.value.trim(), include_inactive: 'true' } })
    if (data.success) { suggestions.value = data.data || [] }
  } catch (err) { suggestions.value = [] }
}

const selectSuggestion = (student) => {
  studentInfo.value = student
  nisnSearch.value = student.nis
  suggestions.value = []
  currentPage.value = 1
  fetchHistory(student.nis)
}

const handleSearch = async () => {
  if (!nisnSearch.value.trim()) return
  suggestions.value = []
  loading.value = true
  studentInfo.value = null
  histories.value = []
  currentPage.value = 1

  const searchTerm = nisnSearch.value.trim()

  // Always try to resolve the student first via search API
  try {
    const { data } = await api.get('/students/search', { params: { q: searchTerm, include_inactive: 'true' } })
    if (data.success && data.data && data.data.length > 0) {
      // Prefer exact NIS match, then fallback to first result
      const exactMatch = data.data.find(s => s.nis === searchTerm)
      studentInfo.value = exactMatch || data.data[0]
    }
  } catch (err) { }

  // Use resolved NIS if we found a student, otherwise try the raw search term as NIS
  const resolvedNis = studentInfo.value ? studentInfo.value.nis : searchTerm
  await fetchHistory(resolvedNis)
}

const fetchHistory = async (nis) => {
  loadingHistory.value = true
  try {
    const { data } = await api.get(`/transactions/student/${nis}`, {
      params: {
        page: currentPage.value,
        page_size: perPage,
        search: innerSearch.value
      }
    })
    if (data.success) {
      histories.value = data.data.items || []
      meta.value = data.data.meta || { total: 0, total_pages: 1, page: 1, page_size: 10 }
      
      // Only populate studentInfo from history if we haven't found them via search
      if (!studentInfo.value && histories.value.length > 0) {
        const first = histories.value[0]
        studentInfo.value = {
          full_name: first.student_name || nis,
          nis: nis,
          class: first.student_class || '-'
        }
      }
      // Do NOT create dummy studentInfo — if student not found, leave null
    }
  } catch (err) {
    if (!innerSearch.value) {
      alert(err.response?.data?.error || 'Siswa tidak ditemukan atau belum memiliki riwayat')
      studentInfo.value = null
    }
    histories.value = []
  } finally {
    loading.value = false
    loadingHistory.value = false
  }
}

watch(currentPage, () => {
  if (studentInfo.value) {
    fetchHistory(studentInfo.value.nis)
  }
})

const exportToExcel = async () => {
  if (!studentInfo.value) return
  try {
    const nis = studentInfo.value.nis
    console.log(`Starting export for student: ${nis}`)
    const response = await api.get(`/transactions/student/${nis}/export`, { 
      responseType: 'blob',
      timeout: 30000 // 30s timeout for export
    })
    
    // Sanitize filename: remove characters that might cause issues in file systems
    const safeName = (studentInfo.value.full_name || 'Siswa')
      .replace(/[/\\?%*:|"<>]/g, '-')
      .replace(/\s+/g, '_')
    
    const filename = `SIS-INV_Riwayat_${safeName}_${new Date().toISOString().slice(0,10)}.xlsx`
    
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', filename)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (e) {
    console.error('Export failed:', e)
    let errorMessage = 'Gagal mendownload riwayat siswa'
    
    if (e.response?.data instanceof Blob && e.response.data.type === 'application/json') {
      try {
        const text = await e.response.data.text()
        const errData = JSON.parse(text)
        errorMessage = errData.error || errorMessage
      } catch (parseErr) {
        console.error('Failed to parse error blob:', parseErr)
      }
    } else if (e.response?.status === 404) {
      errorMessage = 'Fitur ekspor tidak ditemukan di server'
    } else if (e.message) {
      errorMessage = `Kesalahan: ${e.message}`
    }
    
    alert(errorMessage)
  }
}

const formatDate = (d) => {
  if (!d) return '-'
  return new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }).format(new Date(d))
}

const formatDateOnly = (d) => {
  if (!d) return '-'
  return new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'long', year: 'numeric' }).format(new Date(d))
}

const formatTimeOnly = (d) => {
  if (!d) return '-'
  return new Intl.DateTimeFormat('id-ID', { hour: '2-digit', minute: '2-digit' }).format(new Date(d)) + ' WIB'
}
</script>
