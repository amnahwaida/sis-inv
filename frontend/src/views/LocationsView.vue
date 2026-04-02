<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Lokasi Ruangan</h1>
          <p class="text-primary-100/70 text-sm font-medium">Manajemen tata letak penempatan aset sekolah</p>
        </div>
        
        <div class="flex flex-col sm:flex-row items-center gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10 w-full md:w-auto">
          <!-- Search Input -->
          <div class="relative w-full sm:w-64">
            <input type="text" v-model="searchQuery" placeholder="Cari lokasi..." 
                   class="w-full bg-white/10 text-white placeholder-white/50 border-none rounded-xl h-10 pl-10 focus:ring-2 focus:ring-white/50 text-xs font-medium" />
            <svg class="w-4 h-4 absolute left-3.5 top-3 text-white/50" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
          </div>

          <button @click="showModal = true" class="btn-premium-primary">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 4v16m8-8H4" /></svg>
            TAMBAH LOKASI
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="py-20 flex flex-col items-center justify-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
    </div>
    
    <!-- Empty State (No Search Results) -->
    <div v-else-if="paginatedData.length === 0" class="flex flex-col items-center justify-center py-20 bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm">
      <div class="w-20 h-20 bg-gray-50 dark:bg-gray-900 rounded-full flex items-center justify-center text-gray-300 mb-4">
        <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
      </div>
      <p class="text-lg font-black text-gray-400 capitalize">{{ searchQuery ? `Lokasi '${searchQuery}' Tidak Ditemukan` : 'Belum Ada Lokasi' }}</p>
    </div>

    <!-- Locations Grid -->
    <div v-else class="space-y-8">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div v-for="loc in paginatedData" :key="loc.id" class="card-premium group">
          <div class="card-decoration"></div>
          
          <div class="relative z-10 flex flex-col h-full">
            <div class="w-14 h-14 bg-indigo-50 dark:bg-indigo-900/40 text-indigo-600 rounded-2xl flex items-center justify-center mb-6 shadow-sm group-hover:bg-indigo-600 group-hover:text-white transition-all duration-500 translate-y-0 group-hover:-translate-y-2">
              <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
            </div>
            
            <div class="flex-1 space-y-2">
              <h3 class="text-xl font-black text-gray-900 dark:text-white uppercase tracking-tight">{{ loc.name }}</h3>
              <p class="text-xs text-gray-400 font-bold leading-relaxed line-clamp-2 italic">{{ loc.description || 'Lokasi umum penyimpanan aset' }}</p>
            </div>

            <div class="mt-8 pt-6 border-t border-gray-50 dark:border-gray-700 flex items-center justify-between">
              <span class="text-[10px] font-black text-indigo-500 uppercase tracking-widest">{{ loc.item_count || 0 }} ASET TERSIMPAN</span>
              <div class="flex gap-2">
                <button @click="openEditModal(loc)" class="btn-action-edit"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg></button>
                <button @click="deleteLocation(loc)" class="btn-action-delete"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg></button>
              </div>
            </div>
          </div>
        </div>

        <!-- Add Card (only shows on last page if no search) -->
        <button v-if="!searchQuery && currentPage === totalPages" @click="showModal = true" 
                class="border-4 border-dashed border-gray-100 dark:border-gray-700 rounded-[2.5rem] p-8 flex flex-col items-center justify-center group hover:border-indigo-200 dark:hover:border-indigo-900 transition-all cursor-pointer min-h-[250px]">
          <div class="w-16 h-16 bg-gray-50 dark:bg-gray-800 rounded-full flex items-center justify-center text-gray-300 group-hover:scale-110 group-hover:bg-indigo-50 group-hover:text-indigo-500 transition-all duration-500">
            <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 4v16m8-8H4" /></svg>
          </div>
          <p class="mt-4 text-[10px] font-black text-gray-400 group-hover:text-indigo-600 dark:group-hover:text-indigo-400 uppercase tracking-widest">Tambah Ruangan Baru</p>
        </button>
      </div>

      <!-- Pagination Bar -->
      <div v-if="totalPages > 1" class="px-8 py-5 bg-white dark:bg-gray-800 rounded-3xl border border-gray-100 dark:border-gray-700 flex flex-col sm:flex-row items-center justify-between gap-4 shadow-sm">
        <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">
          Menampilkan <span class="text-primary-600">{{ startRow }}-{{ endRow }}</span> dari <span class="text-gray-900 dark:text-white">{{ filteredData.length }}</span> ruangan
        </span>
        <div class="flex gap-2">
          <button @click="currentPage--" :disabled="currentPage === 1" class="pagination-btn-standard">
            Kembali
          </button>
          <button v-for="p in visiblePages" :key="p" @click="currentPage = p"
                  class="w-10 h-10 rounded-xl text-[11px] font-black transition-all shadow-sm active:scale-95 border"
                  :class="p === currentPage ? 'bg-primary-600 text-white border-primary-600' : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:bg-primary-50 dark:hover:bg-gray-700'">
            {{ p }}
          </button>
          <button @click="currentPage++" :disabled="currentPage === totalPages" class="pagination-btn-standard">
            Lanjut
          </button>
        </div>
      </div>
    </div>

    <!-- Modal (Premium Glassmorphic) -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/80 backdrop-blur-xl animate-fade-in" @click.self="closeModal">
      <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] shadow-3xl w-full max-w-md animate-scale-up overflow-hidden border border-white/20 dark:border-gray-700">
        <div class="p-8 border-b border-gray-100 dark:border-gray-700 flex items-center justify-between">
          <div>
            <h3 class="text-2xl font-black text-gray-900 dark:text-white leading-none capitalize">{{ editingData ? 'Ubah Lokasi' : 'Lokasi Baru' }}</h3>
            <p class="text-[10px] font-black text-primary-500 uppercase tracking-widest mt-2">{{ editingData ? 'Perbarui Atribut Ruangan' : 'Tambah Area Penyimpanan' }}</p>
          </div>
          <button @click="closeModal" class="text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 p-3 rounded-2xl transition-all">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>
        
        <form @submit.prevent="saveLocation" class="p-8 space-y-6">
          <div class="space-y-1">
            <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Nama Lokasi / Ruangan</label>
            <input v-model="form.name" class="input-field rounded-2xl h-14" placeholder="Contoh: Lab Komputer 1, GUDANG-A, dsb" required />
          </div>
          <div class="space-y-1">
            <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Deskripsi Ruangan</label>
            <textarea v-model="form.description" class="input-field rounded-2xl p-4" rows="3" placeholder="Jelaskan detail posisi atau fungsi ruangan ini..."></textarea>
          </div>
          
          <div class="flex flex-col sm:flex-row gap-3 pt-6 border-t border-gray-100 dark:border-gray-700">
            <button type="button" @click="closeModal" class="btn-secondary flex-1 py-4 rounded-xl font-black text-[10px] tracking-widest">BATAL</button>
            <button type="submit" :disabled="submitting" class="btn-premium-action flex-[2]">
              {{ submitting ? 'MEMPROSES...' : (editingData ? 'SIMPAN LOKASI' : 'BUAT LOKASI') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import api from '../utils/api'

const locations = ref([])
const loading = ref(true)
const showModal = ref(false)
const submitting = ref(false)
const editingData = ref(null)

const searchQuery = ref('')
const currentPage = ref(1)
const perPage = 11 // 11 items per page so the 'Tambah' card makes it an even 12 grid items on the last page.

const form = ref({ name: '', description: '' })

// Search logic
const filteredData = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return locations.value
  return locations.value.filter(loc => 
    loc.name?.toLowerCase().includes(q) || 
    loc.description?.toLowerCase().includes(q)
  )
})

// Pagination logic
const totalPages = computed(() => Math.max(1, Math.ceil(filteredData.value.length / perPage)))
const startRow = computed(() => (currentPage.value - 1) * perPage + 1)
const endRow = computed(() => Math.min(currentPage.value * perPage, filteredData.value.length))
const paginatedData = computed(() => filteredData.value.slice((currentPage.value - 1) * perPage, currentPage.value * perPage))

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

watch(searchQuery, () => {
  currentPage.value = 1
})

async function fetchLocations() {
  loading.value = true
  try {
    const { data } = await api.get('/locations')
    if (data.success) { locations.value = data.data }
  } catch (e) { console.error(e) } 
  finally { loading.value = false }
}

function openEditModal(loc) {
  editingData.value = loc
  form.value = { name: loc.name, description: loc.description }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingData.value = null
  form.value = { name: '', description: '' }
}

async function saveLocation() {
  submitting.value = true
  try {
    if (editingData.value) { await api.put(`/locations/${editingData.value.id}`, form.value) } 
    else { await api.post('/locations', form.value) }
    closeModal(); fetchLocations()
  } catch (e) { alert(e.response?.data?.error || 'Gagal menyimpan lokasi') } 
  finally { submitting.value = false }
}

async function deleteLocation(loc) {
  if (confirm(`Hapus lokasi "${loc.name}"? Barang yang menggunakan lokasi ini akan tetap tersimpan namun datanya perlu diperbarui.`)) {
    try { 
      await api.delete(`/locations/${loc.id}`)
      await fetchLocations() 
    } catch (e) { 
      alert(e.response?.data?.error || 'Gagal menghapus') 
    }
  }
}

onMounted(fetchLocations)
</script>
