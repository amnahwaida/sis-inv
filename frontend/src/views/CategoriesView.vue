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
          <button @click="showModal = true" class="btn-premium-primary">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
            </svg>
            TAMBAH KATEGORI
          </button>
        </div>
      </div>
    </div>

    <!-- Search & Filters Container -->
    <div class="bg-white dark:bg-gray-800 p-3 rounded-[2rem] border border-gray-100 dark:border-gray-700 shadow-sm flex flex-col sm:flex-row items-center gap-3 w-full">
      <div class="relative w-full sm:w-80">
        <input type="text" v-model="searchQuery" placeholder="Cari nama kategori..." 
               class="input-field pl-10 h-12 rounded-2xl text-sm w-full" />
        <svg class="w-5 h-5 absolute left-3 top-3.5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
        </svg>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="py-20 flex flex-col items-center justify-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
    </div>
    
    <!-- Empty State (No Search Results) -->
    <div v-else-if="categories.length === 0" class="flex flex-col items-center justify-center py-20 bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm">
      <div class="w-20 h-20 bg-gray-50 dark:bg-gray-900 rounded-full flex items-center justify-center text-gray-300 mb-4">
        <svg class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-5.25v9" />
        </svg>
      </div>
      <p class="text-lg font-black text-gray-400 capitalize">{{ searchQuery ? `Kategori '${searchQuery}' Tidak Ditemukan` : 'Belum Ada Kategori' }}</p>
    </div>

    <!-- Categories Grid -->
    <div v-else class="space-y-8">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div v-for="cat in categories" :key="cat.id" class="card-premium group">
          <div class="card-decoration"></div>
          
          <div class="relative z-10 flex flex-col h-full">
            <div class="w-14 h-14 bg-primary-50 dark:bg-primary-900/40 text-primary-600 dark:text-primary-400 rounded-2xl flex items-center justify-center mb-6 shadow-sm group-hover:bg-primary-600 group-hover:text-white transition-all duration-500 translate-y-0 group-hover:-translate-y-2">
              <svg class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9.568 3H5.25A2.25 2.25 0 003 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581a2.25 2.25 0 003.182 0l5.132-5.132a2.25 2.25 0 000-3.182L11.159 3.659A2.25 2.25 0 009.568 3z" />
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 6h.008v.008H6V6z" />
              </svg>
            </div>
            
            <div class="flex-1 space-y-2">
              <h3 class="text-xl font-black text-gray-900 dark:text-white uppercase tracking-tight">{{ cat.name }}</h3>
              <p class="text-xs text-gray-400 font-bold leading-relaxed line-clamp-2 italic">{{ cat.description || 'Tidak ada deskripsi' }}</p>
            </div>

            <div class="mt-8 pt-6 border-t border-gray-50 dark:border-gray-700 flex items-center justify-between">
              <span class="text-[10px] font-black text-primary-500 uppercase tracking-widest">{{ cat.item_count || 0 }} BARANG</span>
              <div class="flex gap-2">
                <button @click="openEditModal(cat)" class="btn-action-edit">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897L16.863 4.487zm0 0L19.5 7.125" />
                  </svg>
                </button>
                <button @click="deleteCategory(cat)" class="btn-action-delete">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Add Card (only shows on last page if no search) -->
        <button v-if="!searchQuery && currentPage === totalPages" @click="showModal = true" 
                class="border-4 border-dashed border-gray-100 dark:border-gray-700 rounded-[2.5rem] p-8 flex flex-col items-center justify-center group hover:border-primary-200 dark:hover:border-primary-900 transition-all cursor-pointer min-h-[250px]">
          <div class="w-16 h-16 bg-gray-50 dark:bg-gray-800 rounded-full flex items-center justify-center text-gray-300 group-hover:scale-110 group-hover:bg-primary-50 group-hover:text-primary-500 transition-all duration-500">
            <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
            </svg>
          </div>
          <p class="mt-4 text-[10px] font-black text-gray-400 group-hover:text-primary-600 dark:group-hover:text-primary-400 uppercase tracking-widest">Klik Untuk Menambah</p>
        </button>
      </div>

      <!-- Pagination Bar -->
      <div class="px-8 py-5 bg-white dark:bg-gray-800 rounded-3xl border border-gray-100 dark:border-gray-700 flex flex-col sm:flex-row items-center justify-between gap-4 shadow-sm">
        <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">
          Menampilkan <span class="text-primary-600">{{ categories.length ? (meta.page - 1) * meta.page_size + 1 : 0 }}-{{ Math.min(meta.page * meta.page_size, meta.total) }}</span> dari <span class="text-gray-900 dark:text-white">{{ meta.total }}</span> kategori
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

    <!-- Modal (Premium Glassmorphic) -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/80 backdrop-blur-xl animate-fade-in" @click.self="closeModal">
      <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] shadow-3xl w-full max-w-md animate-scale-up overflow-hidden border border-white/20 dark:border-gray-700">
        <div class="p-8 border-b border-gray-100 dark:border-gray-700 flex items-center justify-between">
          <div>
            <h3 class="text-2xl font-black text-gray-900 dark:text-white leading-none capitalize">{{ editingData ? 'Ubah Kategori' : 'Kategori Baru' }}</h3>
            <p class="text-[10px] font-black text-primary-500 uppercase tracking-widest mt-2">{{ editingData ? 'Perbarui Atribut Kelompok' : 'Buat Klasifikasi Baru' }}</p>
          </div>
          <button @click="closeModal" class="text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 p-3 rounded-2xl transition-all">
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
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
            <button type="submit" :disabled="submitting" class="btn-premium-action flex-[2]">
              {{ submitting ? 'MEMPROSES...' : (editingData ? 'SIMPAN PERUBAHAN' : 'BUAT KATEGORI') }}
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

const categories = ref([])
const loading = ref(true)
const showModal = ref(false)
const submitting = ref(false)
const editingData = ref(null)

const searchQuery = ref('')
const currentPage = ref(1)
const perPage = 11 
const meta = ref({
  total: 0,
  total_pages: 1,
  page: 1,
  page_size: 11
})

const form = ref({ name: '', description: '' })

async function fetchCategories() {
  loading.value = true
  try {
    const { data } = await api.get('/categories', {
      params: {
        page: currentPage.value,
        page_size: perPage,
        search: searchQuery.value
      }
    })
    if (data.success) { 
      categories.value = data.data.items || []
      meta.value = data.data.meta || { total: 0, total_pages: 1, page: 1, page_size: 11 }
    }
  } catch (e) { 
    console.error(e) 
    categories.value = []
  } finally { 
    loading.value = false 
  }
}

// Debounced search
let searchTimeout
watch(searchQuery, () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
    fetchCategories()
  }, 400)
})

watch(currentPage, fetchCategories)

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
