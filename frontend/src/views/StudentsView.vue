<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Kelola Siswa</h1>
          <p class="text-primary-100/70 text-sm font-medium">Manajemen database dan akses kartu siswa</p>
        </div>
        
        <div class="flex flex-wrap items-center gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10">
          <button @click="exportExcel" 
                  class="btn-premium-action !bg-emerald-500/10 !text-emerald-500 !shadow-none hover:!bg-emerald-500/20 !border-emerald-500/20">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2-2z" /></svg>
            EKSPOR DATA
          </button>
          <button @click="$refs.fileInput.click()" 
                  class="btn-premium-action !bg-blue-500/10 !text-blue-500 !shadow-none hover:!bg-blue-500/20 !border-blue-500/20">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" /></svg>
            IMPORT EXCEL
            <input type="file" ref="fileInput" @change="handleImportExcel" class="hidden" accept=".xlsx, .xls" />
          </button>
          <button @click="openCreateModal" class="btn-premium-primary">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 4v16m8-8H4" /></svg>
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
      <div class="flex-1 relative">
        <input type="text" v-model="filters.search" @input="debouncedFetch"
               class="input-field pl-12 h-14 rounded-2xl" placeholder="Cari nama atau NIS siswa..." />
        <svg class="w-6 h-6 absolute left-4 top-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
      </div>
      <select v-model="filters.class" @change="fetchStudents" class="input-field h-14 rounded-2xl md:w-64 font-bold">
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
                        LIHAT RIWAYAT <svg class="w-2 h-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M14 5l7 7m0 0l-7 7m7-7H3"/></svg>
                      </button>
                    </div>
                  </div>
                </td>
                <td class="px-8 py-6">
                  <div class="text-xs font-bold text-gray-700 dark:text-gray-300 uppercase leading-none mb-1">NIS: {{ s.nisn }}</div>
                </td>
                <td class="px-8 py-6">
                  <span class="px-3 py-1 bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 rounded-lg text-[10px] font-black uppercase tracking-widest">{{ s.class || '-' }}</span>
                </td>
                <td class="px-8 py-6 text-xs font-bold text-gray-400">{{ formatDate(s.created_at) }}</td>
                <td class="px-8 py-6 text-right">
                  <div class="flex items-center justify-end gap-2 transition-opacity">
                    <button @click="openEditModal(s)" class="btn-action-edit" title="Edit">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                    </button>
                    <button @click="deleteStudent(s)" class="btn-action-delete" title="Hapus">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
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
        <div v-if="meta.total_pages > 1" class="px-8 py-6 bg-gray-50/50 dark:bg-gray-700/20 border-t border-gray-100 dark:border-gray-700 flex flex-col sm:flex-row items-center justify-between gap-4">
          <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest leading-none">HALAMAN <span class="text-primary-600">{{ meta.page }}</span> DARI <span class="text-gray-900 dark:text-white">{{ meta.total_pages }}</span></span>
          <div class="flex gap-2">
            <button @click="changePage(meta.page - 1)" :disabled="meta.page === 1" class="pagination-btn-standard">KEMBALI</button>
            <button @click="changePage(meta.page + 1)" :disabled="meta.page === meta.total_pages" class="pagination-btn-standard">LANJUT</button>
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
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12" /></svg>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../utils/api'

const router = useRouter()
const students = ref([])
const classes = ref(['X-1', 'X-2', 'X-3', 'XI-1', 'XI-2', 'XI-3', 'XII-1', 'XII-2', 'XII-3']) // Default classes
const loading = ref(true)
const showModal = ref(false)
const submitting = ref(false)
const editingData = ref(null)
const meta = ref({ page: 1, total_pages: 1 })

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
      // Backend returns array for student list directly in data
      students.value = Array.isArray(data.data) ? data.data : (data.data.items || [])
      meta.value = data.meta || { page: 1, total_pages: 1 }
    }
  } catch (e) {
    console.error('Gagal memuat siswa:', e)
    students.value = []
  } finally {
    loading.value = false
  }
}

const changePage = (p) => { filters.value.page = p; fetchStudents() }

function openCreateModal() { editingData.value = null; form.value = { full_name: '', nisn: '', class: '' }; showModal.value = true }
function openEditModal(student) { 
  editingData.value = student; 
  // Map back to form structure, backend returns both nis and nisn now
  form.value = { 
    full_name: student.full_name, 
    nisn: student.nisn || student.nis, 
    class: student.class 
  }; 
  showModal.value = true 
}
function closeModal() { showModal.value = false; editingData.value = null }
function goToHistory(nisn) { router.push({ path: '/student-history', query: { nisn } }) }

async function saveStudent() {
  submitting.value = true
  try {
    if (editingData.value) { await api.put(`/students/${editingData.value.id}`, form.value) } 
    else { await api.post('/students', form.value) }
    closeModal(); fetchStudents()
  } catch (e) { alert(e.response?.data?.error || 'Gagal menyimpan siswa') } 
  finally { submitting.value = false }
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

onMounted(fetchStudents)
</script>
