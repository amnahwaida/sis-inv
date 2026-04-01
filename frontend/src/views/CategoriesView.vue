<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Kategori Aset</h1>
          <p class="text-primary-100/70 text-sm font-medium">Klasifikasi pengelompokan jenis barang inventaris</p>
        </div>
        
        <div class="flex items-center gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10">
          <button @click="showModal = true" 
                  class="bg-white text-primary-900 hover:bg-primary-50 px-6 py-3 rounded-xl text-[10px] font-black transition-all flex items-center gap-2 shadow-xl active:scale-95">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 4v16m8-8H4" /></svg>
            TAMBAH KATEGORI
          </button>
        </div>
      </div>
    </div>

    <!-- Categories Grid -->
    <div v-if="loading" class="py-20 flex flex-col items-center justify-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
    </div>
    
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
      <div v-for="cat in categories" :key="cat.id" 
           class="group relative bg-white dark:bg-gray-800 rounded-[2.5rem] p-8 border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-2xl hover:-translate-y-2 transition-all duration-500 overflow-hidden">
        <div class="absolute top-0 right-0 w-32 h-32 bg-primary-50 dark:bg-primary-900/10 rounded-full -mt-16 -mr-16 group-hover:scale-150 transition-transform duration-700"></div>
        
        <div class="relative z-10 flex flex-col h-full">
          <div class="w-14 h-14 bg-primary-50 dark:bg-primary-900/40 text-primary-600 dark:text-primary-400 rounded-2xl flex items-center justify-center mb-6 shadow-sm group-hover:bg-primary-600 group-hover:text-white transition-all duration-500 translate-y-0 group-hover:-translate-y-2">
            <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" /></svg>
          </div>
          
          <div class="flex-1 space-y-2">
            <h3 class="text-xl font-black text-gray-900 dark:text-white uppercase tracking-tight">{{ cat.name }}</h3>
            <p class="text-xs text-gray-400 font-bold leading-relaxed line-clamp-2 italic">{{ cat.description || 'Tidak ada deskripsi' }}</p>
          </div>

          <div class="mt-8 pt-6 border-t border-gray-50 dark:border-gray-700 flex items-center justify-between">
            <span class="text-[10px] font-black text-primary-500 uppercase tracking-widest">{{ cat.item_count || 0 }} BARANG</span>
            <div class="flex gap-2">
              <button @click="openEditModal(cat)" class="p-2 text-blue-500 hover:bg-blue-50 dark:hover:bg-blue-900/40 rounded-xl transition-all active:scale-95 relative z-10"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg></button>
              <button @click="deleteCategory(cat)" class="p-2 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/40 rounded-xl transition-all active:scale-95 relative z-10"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg></button>
            </div>
          </div>
        </div>
      </div>

      <!-- Add Card -->
      <button @click="showModal = true" 
              class="border-4 border-dashed border-gray-100 dark:border-gray-700 rounded-[2.5rem] p-8 flex flex-col items-center justify-center group hover:border-primary-200 dark:hover:border-primary-900 transition-all cursor-pointer min-h-[250px]">
        <div class="w-16 h-16 bg-gray-50 dark:bg-gray-800 rounded-full flex items-center justify-center text-gray-300 group-hover:scale-110 group-hover:bg-primary-50 group-hover:text-primary-500 transition-all duration-500">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 4v16m8-8H4" /></svg>
        </div>
        <p class="mt-4 text-[10px] font-black text-gray-400 group-hover:text-primary-600 dark:group-hover:text-primary-400 uppercase tracking-widest">Klik Untuk Menambah</p>
      </button>
    </div>

    <!-- Modal (Premium Glassmorphic) -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/80 backdrop-blur-xl animate-fade-in" @click.self="closeModal">
      <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] shadow-3xl w-full max-w-md animate-scale-up overflow-hidden border border-white/20 dark:border-gray-700">
        <div class="p-8 border-b border-gray-100 dark:border-gray-700 flex items-center justify-between">
          <div>
            <h3 class="text-2xl font-black text-gray-900 dark:text-white leading-none capitalize">{{ editingData ? 'Ubah Kategori' : 'Kategori Baru' }}</h3>
            <p class="text-[10px] font-black text-primary-500 uppercase tracking-widest mt-2">{{ editingData ? 'Perbarui Atribut Kelompok' : 'Buat Klasifikasi Baru' }}</p>
          </div>
          <button @click="closeModal" class="text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 p-3 rounded-2xl transition-all">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>
        
        <form @submit.prevent="saveCategory" class="p-8 space-y-6">
          <div class="space-y-1">
            <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Nama Kategori</label>
            <input v-model="form.name" class="input-field rounded-2xl h-14" placeholder="Contoh: Elektronik, Perabot, dsb" required />
          </div>
          <div class="space-y-1">
            <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Keterangan Singkat</label>
            <textarea v-model="form.description" class="input-field rounded-2xl p-4" rows="3" placeholder="Jelaskan jenis barang dalam kategori ini..."></textarea>
          </div>
          
          <div class="flex flex-col sm:flex-row gap-3 pt-6 border-t border-gray-100 dark:border-gray-700">
            <button type="button" @click="closeModal" class="btn-secondary flex-1 py-4 rounded-xl font-black text-[10px] tracking-widest">BATAL</button>
            <button type="submit" :disabled="submitting" 
                    class="bg-primary-600 text-white flex-[2] py-4 rounded-xl font-black text-[10px] tracking-[0.3em] shadow-xl shadow-primary-500/20 active:scale-95 disabled:opacity-30 transition-all uppercase">
              {{ submitting ? 'MEMPROSES...' : (editingData ? 'SIMPAN PERUBAHAN' : 'BUAT KATEGORI') }}
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

const categories = ref([])
const loading = ref(true)
const showModal = ref(false)
const submitting = ref(false)
const editingData = ref(null)

const form = ref({ name: '', description: '' })

async function fetchCategories() {
  loading.value = true
  try {
    const { data } = await api.get('/categories')
    if (data.success) { categories.value = data.data }
  } catch (e) { console.error(e) } 
  finally { loading.value = false }
}

function openEditModal(cat) {
  editingData.value = cat
  form.value = { name: cat.name, description: cat.description }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingData.value = null
  form.value = { name: '', description: '' }
}

async function saveCategory() {
  submitting.value = true
  try {
    if (editingData.value) { await api.put(`/categories/${editingData.value.id}`, form.value) } 
    else { await api.post('/categories', form.value) }
    closeModal(); fetchCategories()
  } catch (e) { alert(e.response?.data?.error || 'Gagal menyimpan kategori') } 
  finally { submitting.value = false }
}

async function deleteCategory(cat) {
  if (confirm(`Hapus kategori "${cat.name}"? Barang yang menggunakan kategori ini akan kehilangan relasinya.`)) {
    try { 
      await api.delete(`/categories/${cat.id}`)
      await fetchCategories() 
    } catch (e) { 
      alert(e.response?.data?.error || 'Gagal menghapus') 
    }
  }
}

onMounted(fetchCategories)
</script>
