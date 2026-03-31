<template>
  <div class="animate-fade-in space-y-6">
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-gray-900 tracking-tight">Kelola Siswa</h1>
        <p class="text-sm text-gray-500 mt-1">Daftar siswa yang terdaftar dalam sistem untuk peminjaman barang.</p>
      </div>
      <div class="flex items-center gap-2">
        <button 
          @click="exportStudents" 
          class="px-4 py-2 bg-white border border-gray-200 text-gray-700 font-bold rounded-xl hover:bg-gray-50 transition-all flex items-center gap-2"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/></svg>
          Export CSV
        </button>
        <button 
          @click="triggerImport" 
          class="px-4 py-2 bg-white border border-gray-200 text-gray-700 font-bold rounded-xl hover:bg-gray-50 transition-all flex items-center gap-2"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/></svg>
          Import
        </button>
        <button 
          @click="openModal()" 
          class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white font-bold rounded-xl shadow-lg shadow-primary-100 transition-all flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
          Tambah Siswa
        </button>
        <input type="file" ref="fileInput" class="hidden" accept=".csv" @change="handleImport">
      </div>
    </div>

    <!-- Table -->
    <div class="card p-0 overflow-hidden border-gray-100 shadow-sm bg-white rounded-3xl">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50/50 border-b border-gray-100">
            <tr>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">NIS</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Nama Lengkap</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Kelas</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Status</th>
              <th class="text-right text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-if="loading && students.length === 0">
              <td colspan="5" class="px-6 py-12 text-center text-gray-400">Memuat data...</td>
            </tr>
            <tr v-else-if="students.length === 0">
              <td colspan="5" class="px-6 py-12 text-center text-gray-400">Belum ada data siswa</td>
            </tr>
            <tr v-for="s in students" :key="s.id" class="hover:bg-gray-50/50 transition-colors">
              <td class="px-6 py-4 font-mono text-sm text-gray-600">{{ s.nis }}</td>
              <td class="px-6 py-4 font-bold text-gray-900">{{ s.full_name }}</td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ s.class || '-' }}</td>
              <td class="px-6 py-4">
                <span :class="s.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'" class="px-2.5 py-1 rounded-lg text-[10px] font-black uppercase tracking-wider">
                  {{ s.is_active ? 'Aktif' : 'Non-aktif' }}
                </span>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button @click="openModal(s)" class="p-2 text-primary-600 hover:bg-primary-50 rounded-lg transition-colors">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
                  </button>
                  <button @click="deleteStudent(s)" class="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
      <div class="bg-white rounded-3xl w-full max-w-md shadow-2xl overflow-hidden animate-scale-up">
        <div class="p-6 border-b border-gray-100 flex items-center justify-between bg-primary-50/50">
          <h3 class="text-xl font-black text-gray-900 tracking-tight">{{ editingStudent ? 'Edit Siswa' : 'Tambah Siswa' }}</h3>
          <button @click="closeModal()" class="text-gray-400 hover:text-gray-600"><svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg></button>
        </div>
        
        <form @submit.prevent="saveStudent" class="p-6 space-y-4">
          <div>
            <label class="block text-xs font-black text-gray-500 uppercase tracking-widest mb-1.5 ml-1">NIS *</label>
            <input v-model="form.nis" required class="input-field rounded-xl" placeholder="Masukkan Nomor Induk Siswa">
          </div>
          <div>
            <label class="block text-xs font-black text-gray-500 uppercase tracking-widest mb-1.5 ml-1">Nama Lengkap *</label>
            <input v-model="form.full_name" required class="input-field rounded-xl" placeholder="Nama lengkap siswa">
          </div>
          <div>
            <label class="block text-xs font-black text-gray-500 uppercase tracking-widest mb-1.5 ml-1">Kelas</label>
            <input v-model="form.class" class="input-field rounded-xl" placeholder="Contoh: 12 IPA 1">
          </div>
          <div v-if="editingStudent" class="flex items-center gap-3 p-3 bg-gray-50 rounded-xl">
             <input type="checkbox" id="is_active" v-model="form.is_active" class="w-4 h-4 rounded text-primary-600 focus:ring-primary-500 border-gray-300">
             <label for="is_active" class="text-sm font-bold text-gray-700 select-none">Status Aktif</label>
          </div>

          <div class="flex gap-3 pt-4">
            <button type="button" @click="closeModal()" class="flex-1 px-4 py-2.5 bg-gray-100 hover:bg-gray-200 text-gray-700 font-bold rounded-xl transition-all">Batal</button>
            <button type="submit" :disabled="submitting" class="flex-[2] px-4 py-2.5 bg-primary-600 hover:bg-primary-700 text-white font-bold rounded-xl shadow-lg shadow-primary-100 transition-all flex items-center justify-center gap-2">
              <svg v-if="submitting" class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
              {{ submitting ? 'Menyimpan...' : 'Simpan Data' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../utils/api'

const students = ref([])
const loading = ref(false)
const showModal = ref(false)
const submitting = ref(false)
const editingStudent = ref(null)
const fileInput = ref(null)

const form = ref({
  nis: '',
  full_name: '',
  class: '',
  is_active: true
})

async function fetchStudents() {
  loading.value = true
  try {
    const { data } = await api.get('/students')
    if (data.success) {
      students.value = data.data
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function exportStudents() {
  try {
    const response = await api.get('/students/export', { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', 'daftar_siswa.csv')
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (e) {
    alert('Gagal mendownload CSV')
  }
}

function triggerImport() {
  fileInput.value.click()
}

async function handleImport(event) {
  const file = event.target.files[0]
  if (!file) return

  const formData = new FormData()
  formData.append('file', file)

  submitting.value = true
  try {
    const { data } = await api.post('/students/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    alert(data.message + ` (${data.data.total} siswa terimpor)`)
    await fetchStudents()
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal mengimpor CSV')
  } finally {
    submitting.value = false
    event.target.value = '' // Reset input
  }
}

function openModal(student = null) {
  if (student) {
    editingStudent.value = student
    form.value = { 
      nis: student.nis, 
      full_name: student.full_name, 
      class: student.class, 
      is_active: student.is_active 
    }
  } else {
    editingStudent.value = null
    form.value = { nis: '', full_name: '', class: '', is_active: true }
  }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingStudent.value = null
}

async function saveStudent() {
  submitting.value = true
  try {
    if (editingStudent.value) {
      await api.put(`/students/${editingStudent.value.id}`, form.value)
    } else {
      await api.post('/students', form.value)
    }
    closeModal()
    await fetchStudents()
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal menyimpan data siswa')
  } finally {
    submitting.value = false
  }
}

async function deleteStudent(student) {
  if (confirm(`Hapus data siswa ${student.full_name}?`)) {
    try {
      await api.delete(`/students/${student.id}`)
      await fetchStudents()
    } catch (e) {
      alert(e.response?.data?.error || 'Gagal menghapus siswa')
    }
  }
}

onMounted(fetchStudents)
</script>
