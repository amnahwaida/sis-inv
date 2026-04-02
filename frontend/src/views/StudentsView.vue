<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-primary-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Kelola Siswa</h1>
          <p class="text-primary-100/70 text-sm font-medium">Manajemen database dan akses kartu siswa</p>
        </div>
        
        <div class="flex flex-wrap items-center gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10">
          <button @click="exportExcel" 
                  class="btn-premium-action !bg-emerald-500/10 !text-emerald-500 !shadow-none hover:!bg-emerald-500/20 !border-emerald-500/20">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m.75 12l3 3m0 0l3-3m-3 3v-6m-1.5-9H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
            </svg>
            EKSPOR DATA
          </button>
          <button @click="$refs.fileInput.click()" 
                  class="btn-premium-action !bg-primary-500/10 !text-primary-500 !shadow-none hover:!bg-primary-500/20 !border-primary-500/20">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5" />
            </svg>
            IMPORT EXCEL
            <input type="file" ref="fileInput" @change="handleImportExcel" class="hidden" accept=".xlsx, .xls" />
          </button>
          <button @click="downloadTemplate" 
                  class="btn-premium-action !bg-gray-500/10 !text-gray-400 !shadow-none hover:!bg-gray-500/20 !border-gray-500/20">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 3h5.25m-5.25 3h5.25m-5.25 3h5.25m-10.5-9H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
            </svg>
            TEMPLATE
          </button>
          <button @click="openCreateModal" class="btn-premium-primary">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
            </svg>
            TAMBAH SISWA BARU
          </button>
        </div>
      </div>
    </div>

    <!-- Stats Mini Grid -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
      <div class="bg-white dark:bg-gray-800 p-6 rounded-[2rem] border border-gray-100 dark:border-gray-700 shadow-sm">
        <p class="text-3xl font-black text-primary-600 leading-none">{{ students.length }}</p>
        <p class="text-[9px] font-black text-gray-400 dark:text-gray-500 uppercase tracking-widest mt-2">Total Terdaftar</p>
      </div>
    </div>

    <!-- Search & Filter Area -->
    <div class="bg-white dark:bg-gray-800 p-3 rounded-[2rem] border border-gray-100 dark:border-gray-700 shadow-sm flex flex-col md:flex-row gap-3">
      <div class="relative w-full">
        <input type="text" v-model="filters.search" @input="debouncedFetch"
               class="input-field pl-10 h-12 rounded-2xl text-sm w-full" placeholder="Cari nama atau NIS siswa..." />
        <svg class="w-5 h-5 absolute left-3 top-3.5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
        </svg>
      </div>
      <select v-model="filters.class" @change="fetchStudents" 
              class="input-field h-12 rounded-2xl text-sm w-full md:w-64 appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2020%2020%22%3E%3Cpath%20stroke%3D%22%236b7280%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%221.5%22%20d%3D%22m6%208%204%204%204-4%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[right_0.5rem_center] bg-no-repeat pr-10">
        <option value="">Semua Kelas</option>
        <option v-for="cls in classes" :key="cls" :value="cls">{{ cls }}</option>
      </select>
    </div>

    <!-- Table Container -->
    <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden transition-all duration-300">
      <div class="overflow-x-auto relative">
        <!-- Loading Overlay -->
        <div v-if="loading" class="absolute inset-0 bg-white/60 dark:bg-gray-900/60 backdrop-blur-sm z-10 flex flex-col items-center justify-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
          <span class="mt-4 text-xs font-black text-primary-600 uppercase tracking-widest leading-none">Menyinkronkan Siswa...</span>
        </div>

        <table class="w-full">
          <thead>
            <tr class="bg-gray-50/50 dark:bg-gray-700/30">
              <th v-for="h in ['Siswa', 'Data Identitas', 'Kelas', 'Tgl Terdaftar', '']" :key="h"
                  class="text-left py-5 px-8 text-[10px] font-black text-gray-400 dark:text-gray-300 uppercase tracking-[0.2em]">{{ h }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50 dark:divide-gray-700">
            <template v-if="students.length > 0">
              <tr v-for="s in students" :key="s.id" class="group hover:bg-primary-50/50 dark:hover:bg-primary-900/10 transition-all duration-300">
                <td class="px-8 py-6">
                  <div class="flex items-center gap-4">
                    <div class="w-12 h-12 rounded-2xl bg-gray-100 dark:bg-gray-700 text-gray-400 group-hover:bg-primary-100 dark:group-hover:bg-primary-900/50 group-hover:text-primary-600 flex items-center justify-center text-sm font-black transition-all shadow-sm">
                      {{ s.full_name?.charAt(0) }}
                    </div>
                    <div>
                      <p class="text-sm font-black text-gray-900 dark:text-white leading-none uppercase tracking-tight">{{ s.full_name }}</p>
                      <button @click="goToHistory(s.nisn)" class="text-[10px] text-primary-600 dark:text-primary-400 mt-1 font-bold tracking-widest flex items-center gap-1 hover:underline">
                        LIHAT RIWAYAT 
                        <svg class="w-2.5 h-2.5" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
                        </svg>
                      </button>
                    </div>
                  </div>
                </td>
                <td class="px-8 py-6">
                  <div class="text-xs font-bold text-gray-700 dark:text-gray-300 uppercase leading-none mb-1">NIS: {{ s.nisn }}</div>
                </td>
                <td class="px-8 py-6">
                  <div class="flex items-center gap-2">
                    <span class="px-3 py-1 bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 rounded-lg text-[10px] font-black uppercase tracking-widest">{{ s.class || '-' }}</span>
                    <span v-if="s.is_active" class="px-2 py-1 bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 rounded-lg text-[9px] font-black uppercase tracking-widest border border-emerald-100 dark:border-emerald-800">Aktif</span>
                    <span v-else class="px-2 py-1 bg-red-50 dark:bg-red-900/30 text-red-500 rounded-lg text-[9px] font-black uppercase tracking-widest border border-red-100 dark:border-red-800">Non-Aktif</span>
                  </div>
                </td>
                <td class="px-8 py-6 text-xs font-bold text-gray-400">{{ formatDate(s.created_at) }}</td>
                <td class="px-8 py-6 text-right">
                  <div class="flex items-center justify-end gap-2 transition-opacity">
                    <!-- Toggle Status Button -->
                    <button @click="toggleStatus(s)" 
                            :class="s.is_active ? 'text-amber-500 hover:bg-amber-50 dark:hover:bg-amber-900/30' : 'text-emerald-500 hover:bg-emerald-50 dark:hover:bg-emerald-900/30'"
                            class="p-2 rounded-xl transition-all" 
                            :title="s.is_active ? 'Non-aktifkan' : 'Aktifkan'">
                      <svg v-if="s.is_active" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9.75 9.75l4.5 4.5m0-4.5l-4.5 4.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                    </button>
                    <button @click="openEditModal(s)" class="btn-action-edit" title="Edit">
                      <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897L16.863 4.487zm0 0L19.5 7.125" />
                      </svg>
                    </button>
                    <button @click="deleteStudent(s)" class="btn-action-delete" title="Hapus">
                      <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </template>
            <tr v-else-if="!loading" class="text-center">
              <td colspan="5" class="px-8 py-20 italic text-gray-400 font-bold uppercase tracking-widest text-xs">Belum ada data siswa tersimpan</td>
            </tr>
          </tbody>
        </table>

        <!-- Modern Pagination Bar -->
        <div class="px-8 py-6 bg-gray-50/50 dark:bg-gray-700/20 border-t border-gray-100 dark:border-gray-700 flex flex-col sm:flex-row items-center justify-between gap-4">
          <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest leading-none">
            Menampilkan <span class="text-primary-600">{{ students.length ? (meta.page - 1) * meta.page_size + 1 : 0 }}-{{ Math.min(meta.page * meta.page_size, meta.total) }}</span> dari <span class="text-gray-900 dark:text-white">{{ meta.total }}</span> siswa
          </span>
          <div class="flex gap-2">
            <button @click="changePage(meta.page - 1)" :disabled="meta.page === 1" class="pagination-btn-standard">KEMBALI</button>
            <div class="flex gap-1">
              <button v-for="p in visiblePages" :key="p" @click="changePage(p)"
                      class="w-10 h-10 rounded-xl text-[11px] font-black transition-all shadow-sm active:scale-95 border"
                      :class="p === meta.page ? 'bg-primary-600 text-white border-primary-600' : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:bg-primary-50 dark:hover:bg-gray-700'">
                {{ p }}
              </button>
            </div>
            <button @click="changePage(meta.page + 1)" :disabled="meta.page === meta.total_pages || meta.total_pages === 0" class="pagination-btn-standard">LANJUT</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Student Modal (Glassmorphic) -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/80 backdrop-blur-xl animate-fade-in" @click.self="closeModal">
      <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] shadow-3xl w-full max-w-md animate-scale-up overflow-hidden border border-white/20 dark:border-gray-700">
        <div class="p-8 border-b border-gray-100 dark:border-gray-700 flex items-center justify-between">
          <div>
            <h3 class="text-2xl font-black text-gray-900 dark:text-white leading-none capitalize">{{ editingData ? 'Edit Data Siswa' : 'Siswa Baru' }}</h3>
            <p class="text-[10px] font-black text-primary-500 uppercase tracking-widest mt-2">{{ editingData ? 'Perbarui Profil Siswa' : 'Registrasi Profil Baru' }}</p>
          </div>
          <button @click="closeModal" class="text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 p-3 rounded-2xl transition-all">
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <form @submit.prevent="saveStudent" class="p-8 space-y-6">
          <div class="space-y-1">
            <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Nama Lengkap Siswa</label>
            <input v-model="form.full_name" class="input-field rounded-2xl h-14" placeholder="Contoh: Ahmad Subardjo" required />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1">
              <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">NIS (Sesuai Kartu)</label>
              <input v-model="form.nisn" class="input-field rounded-2xl h-14" placeholder="Nomor Induk" required />
            </div>
            <div class="space-y-1">
              <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Kelas / Angkatan</label>
              <input v-model="form.class" class="input-field rounded-2xl h-14" placeholder="Contoh: XII-MIPA-1" required />
            </div>
          </div>
          
          <div class="flex flex-col sm:flex-row gap-3 pt-6 border-t border-gray-100 dark:border-gray-700">
            <button type="button" @click="closeModal" class="btn-secondary flex-1 py-4 rounded-xl font-black text-[10px] tracking-widest">BATAL</button>
            <button type="submit" :disabled="submitting" class="btn-premium-action flex-[2]">
              {{ submitting ? 'MENYIMPAN...' : (editingData ? 'SIMPAN PERUBAHAN' : 'DAFTARKAN SISWA') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../utils/api'

const router = useRouter()
const students = ref([])
const classes = ref([])
const loading = ref(true)
const showModal = ref(false)
const submitting = ref(false)
const editingData = ref(null)
const meta = ref({ page: 1, total_pages: 1, total: 0, page_size: 15 })

const filters = ref({ search: '', class: '', page: 1, page_size: 15 })
const form = ref({ full_name: '', nisn: '', class: '' })

let searchTimeout
const debouncedFetch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => { filters.value.page = 1; fetchStudents() }, 500)
}

const fetchStudents = async () => {
  loading.value = true
  try {
    const { data } = await api.get('/students', { params: filters.value })
    if (data.success) {
      students.value = data.data.items || []
      meta.value = data.data.meta || { page: 1, total_pages: 1, total: 0, page_size: 15 }
    }
  } catch (e) {
    console.error('Gagal memuat siswa:', e)
    students.value = []
  } finally {
    loading.value = false
  }
}

async function fetchClasses() {
  try {
    const { data } = await api.get('/students/classes')
    if (data.success) { classes.value = data.data || [] }
  } catch (e) { console.error('Gagal mengambil daftar kelas dari master data:', e) }
}

const changePage = (p) => { 
  if (p < 1 || p > meta.value.total_pages) return
  filters.value.page = p
  fetchStudents() 
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

function openCreateModal() { editingData.value = null; form.value = { full_name: '', nisn: '', class: '', is_active: true }; showModal.value = true }
function openEditModal(student) { 
  editingData.value = student; 
  // Map back to form structure, backend returns both nis and nisn now
  form.value = { 
    full_name: student.full_name, 
    nisn: student.nisn || student.nis, 
    class: student.class,
    is_active: student.is_active
  }; 
  showModal.value = true 
}
function closeModal() { showModal.value = false; editingData.value = null }
function goToHistory(nisn) { router.push({ path: '/student-history', query: { nisn } }) }

const saveStudent = async () => {
  submitting.value = true
  try {
    if (editingData.value) {
      await api.put(`/students/${editingData.value.id}`, {
        full_name: form.value.full_name,
        nisn: form.value.nisn,
        nis: form.value.nisn, // Send as both to satisfy API
        class: form.value.class,
        is_active: form.value.is_active
      })
    } else {
      await api.post('/students', {
        full_name: form.value.full_name,
        nisn: form.value.nisn,
        nis: form.value.nisn,
        class: form.value.class,
        is_active: true
      })
    }
    closeModal()
    fetchStudents()
    fetchClasses()
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal menyimpan data siswa')
  } finally {
    submitting.value = false
  }
}

const toggleStatus = async (student) => {
  if (!confirm(`Apakah Anda yakin ingin ${student.is_active ? 'menonaktifkan' : 'mengaktifkan'} siswa ${student.full_name}?`)) return
  try {
    await api.put(`/students/${student.id}`, {
      full_name: student.full_name,
      nisn: student.nisn || student.nis,
      nis: student.nisn || student.nis,
      class: student.class,
      is_active: !student.is_active
    })
    fetchStudents()
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal mengubah status siswa')
  }
}

async function exportExcel() {
  try {
    const response = await api.get('/students/export', { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `SIS-INV_Daftar_Siswa_${new Date().toISOString().slice(0,10)}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (e) { alert('Gagal mendownload data siswa') }
}

async function downloadTemplate() {
  try {
    const response = await api.get('/students/import/template', { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', 'SIS-INV_Template_Siswa.xlsx')
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (e) { alert('Gagal mendownload template excel') }
}

async function handleImportExcel(event) {
  const file = event.target.files[0]
  if (!file) return
  
  const formData = new FormData()
  formData.append('file', file)
  
  loading.value = true
  try {
    const { data } = await api.post('/students/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    if (data.success) {
      alert(`Berhasil mengimpor ${data.data.total} siswa`)
      fetchStudents()
    }
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal mengimpor file excel. Pastikan format kolom sesuai: NIS, Nama Lengkap, Kelas')
  } finally {
    loading.value = false
    event.target.value = '' // Clear input
  }
}

async function deleteStudent(student) {
  if (confirm(`Apakah Anda yakin ingin menghapus siswa ${student.full_name}?`)) {
    try { 
      await api.delete(`/students/${student.id}`)
      await fetchStudents() 
    } catch (e) { 
      alert(e.response?.data?.error || 'Gagal menghapus siswa') 
    }
  }
}

const formatDate = (d) => d ? new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'short' }).format(new Date(d)) : '-'

onMounted(() => {
  fetchStudents()
  fetchClasses()
})
</script>
