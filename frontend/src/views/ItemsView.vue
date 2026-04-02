<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col lg:flex-row lg:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Daftar Barang</h1>
          <p class="text-primary-100/70 text-sm font-medium">Manajemen seluruh aset dan inventaris sekolah</p>
        </div>
        
        <div class="flex flex-wrap items-center gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10">
          <button v-if="authStore.isAdmin" @click="goToPrint" :disabled="selectedIds.length === 0"
                  class="bg-indigo-500/20 hover:bg-indigo-500/30 text-indigo-400 px-4 py-3 rounded-xl text-[10px] font-black transition-all flex items-center gap-2 active:scale-95 disabled:opacity-30">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" /></svg>
            CETAK ({{ selectedIds.length }})
          </button>
          <button v-if="authStore.isAdmin" @click="handleExportExcel" 
                  class="bg-emerald-500/20 hover:bg-emerald-500/30 text-emerald-400 px-4 py-3 rounded-xl text-[10px] font-black transition-all flex items-center gap-2 active:scale-95">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
            EXCEL
          </button>
          <button v-if="authStore.isAdmin" @click="$refs.fileInput.click()" 
                  class="bg-blue-500/20 hover:bg-blue-500/30 text-blue-400 px-4 py-3 rounded-xl text-[10px] font-black transition-all flex items-center gap-2 active:scale-95">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" /></svg>
            IMPORT
            <input type="file" ref="fileInput" @change="handleImportExcel" class="hidden" accept=".xlsx, .xls" />
          </button>
          <button v-if="authStore.isAdmin" @click="downloadTemplate" 
                  class="bg-gray-500/20 hover:bg-gray-500/30 text-gray-400 px-4 py-3 rounded-xl text-[10px] font-black transition-all flex items-center gap-2 active:scale-95">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
            TEMPLATE
          </button>
          <button v-if="authStore.isAdmin" @click="openAddModal" class="btn-premium-primary">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 4v16m8-8H4" /></svg>
            TAMBAH BARANG
          </button>
        </div>
      </div>
    </div>

    <!-- Search & Filters Container -->
    <div class="bg-white dark:bg-gray-800 p-3 rounded-[2rem] border border-gray-100 dark:border-gray-700 shadow-sm">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-3">
        <div class="lg:col-span-1 relative">
          <input type="text" class="input-field pl-10 h-12 rounded-2xl" placeholder="Cari nama / kode..." v-model="filters.search" @input="debouncedFetch" />
          <svg class="w-5 h-5 absolute left-3 top-3.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
        </div>
        <select class="input-field h-12 rounded-2xl" v-model="filters.status" @change="fetchData">
          <option value="">Semua Status</option>
          <option value="AVAILABLE">Tersedia</option>
          <option value="BORROWED">Dipinjam</option>
          <option value="MAINTENANCE">Perbaikan</option>
          <option value="LOST">Hilang</option>
        </select>
        <select class="input-field h-12 rounded-2xl" v-model="filters.category_id" @change="fetchData">
          <option value="">Semua Kategori</option>
          <option v-for="cat in categoryStore.categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
        </select>
        <select class="input-field h-12 rounded-2xl" v-model="filters.location" @change="fetchData">
          <option value="">Semua Lokasi</option>
          <option v-for="loc in locationStore.locations" :key="loc.id" :value="loc.name">{{ loc.name }}</option>
        </select>
        <select class="input-field h-12 rounded-2xl" v-model="filters.condition" @change="fetchData">
          <option value="">Semua Kondisi</option>
          <option value="GOOD">Baik</option>
          <option value="DAMAGED">Rusak</option>
          <option value="IN_REPAIR">Diperbaiki</option>
        </select>
      </div>
    </div>

    <!-- Main Table Card -->
    <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden transition-all duration-300">
      <div class="overflow-x-auto relative">
        <!-- Loading Overlay -->
        <div v-if="itemStore.loading" class="absolute inset-0 bg-white/60 dark:bg-gray-900/60 backdrop-blur-sm z-10 flex flex-col items-center justify-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
          <span class="mt-4 text-xs font-black text-primary-600 uppercase tracking-widest">Sinkronisasi Aset...</span>
        </div>

        <table class="w-full">
          <thead>
            <tr class="bg-gray-50/50 dark:bg-gray-700/30">
              <th v-if="authStore.isAdmin" class="w-10 pl-8 py-5">
                <input type="checkbox" @change="toggleSelectAll" :checked="isAllSelected" class="rounded-lg border-gray-300 text-primary-600 focus:ring-primary-500 w-5 h-5" />
              </th>
              <th v-for="h in ['Kode / Nama', 'Kategori & Lokasi', 'Status', 'Kondisi', '']" :key="h"
                  class="text-left py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">{{ h }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50 dark:divide-gray-700">
            <template v-if="itemStore.items.length > 0">
              <tr v-for="item in itemStore.items" :key="item.id" 
                  class="group hover:bg-primary-50/50 dark:hover:bg-primary-900/10 transition-all duration-300"
                  :class="{'bg-primary-50/30 dark:bg-primary-900/5': selectedIds.includes(item.id)}">
                <td v-if="authStore.isAdmin" class="pl-8 py-6">
                  <input type="checkbox" v-model="selectedIds" :value="item.id" class="rounded-lg border-gray-300 text-primary-600 focus:ring-primary-500 w-5 h-5" />
                </td>
                <td class="px-8 py-6">
                  <div class="flex items-center gap-4">
                    <div class="w-12 h-12 bg-gray-100 dark:bg-gray-700 rounded-2xl flex items-center justify-center text-gray-400 dark:text-gray-500 group-hover:bg-primary-100 dark:group-hover:bg-primary-900/50 group-hover:text-primary-600 transition-colors">
                      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>
                    </div>
                    <div>
                      <div class="font-black text-gray-900 dark:text-white text-sm">{{ item.name }}</div>
                      <div class="text-[10px] text-gray-400 font-bold font-mono tracking-tighter">{{ item.code }}</div>
                    </div>
                  </div>
                </td>
                <td class="px-8 py-6">
                  <div class="text-xs font-bold text-gray-700 dark:text-gray-300">{{ item.category_name || '-' }}</div>
                  <div class="text-[10px] text-gray-400 mt-0.5">{{ item.location_name || item.location || '-' }}</div>
                </td>
                <td class="px-8 py-6">
                  <span class="px-3 py-1 rounded-lg font-black uppercase text-[9px] tracking-widest shadow-sm inline-block" 
                        :class="getStatusBadgeClass(item.status)">
                    {{ getStatusLabel(item.status) }}
                  </span>
                </td>
                <td class="px-8 py-6">
                  <span class="px-3 py-1 rounded-lg font-black uppercase text-[9px] tracking-widest shadow-sm inline-block"
                        :class="getConditionBadgeClass(item.condition)">
                    {{ getConditionLabel(item.condition) }}
                  </span>
                </td>
                <td class="px-8 py-6 text-right">
                  <div class="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                    <button @click="openDetailModal(item)" class="btn-action-view" title="Detail">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
                    </button>
                    <template v-if="authStore.isAdmin">
                      <button @click="openQrModal(item)" class="p-2 text-emerald-500 bg-emerald-50 dark:bg-emerald-900/30 rounded-xl hover:bg-emerald-500 hover:text-white transition-all transform hover:scale-110 active:scale-95" title="QR Code">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1|0 11v1|9-9h-1|4 12H3|15.364 6.364l-.707-.707|6.343 6.343l-.707-.707|12.728 0l-.707.707|6.343 17.657l-.707.707|16 12a4 4 0 11-8 0 4 4 0 018 0z"/></svg>
                      </button>
                      <button @click="openEditModal(item)" class="btn-action-edit" title="Edit">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
                      </button>
                      <button @click="deleteItem(item)" class="btn-action-delete" title="Hapus">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                      </button>
                    </template>
                  </div>
                </td>
              </tr>
            </template>
            <tr v-else-if="!itemStore.loading" class="text-center">
              <td colspan="6" class="px-8 py-24">
                <div class="bg-gray-50 dark:bg-gray-700/50 w-24 h-24 rounded-full flex items-center justify-center mx-auto mb-6 text-gray-300 dark:text-gray-600">
                  <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>
                </div>
                <p class="text-gray-900 dark:text-white font-black text-xl">Aset Tidak Ditemukan</p>
                <p class="text-gray-400 text-sm mt-2 max-w-xs mx-auto">Coba sesuaikan kata kunci pencarian atau filter kategori Anda.</p>
              </td>
            </tr>
          </tbody>
        </table>
        
        <!-- Pagination Premium -->
        <div class="px-8 py-6 bg-gray-50/50 dark:bg-gray-700/20 border-t border-gray-100 dark:border-gray-700 flex flex-col sm:flex-row items-center justify-between gap-4">
          <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">
            Menampilkan <span class="text-primary-600">{{ itemStore.items.length ? (itemStore.meta.page - 1) * itemStore.meta.page_size + 1 : 0 }}-{{ Math.min(itemStore.meta.page * itemStore.meta.page_size, itemStore.meta.total) }}</span> dari <span class="text-gray-900 dark:text-white">{{ itemStore.meta.total }}</span> aset
          </span>
          <div class="flex gap-2">
            <button @click="changePage(itemStore.meta.page - 1)" :disabled="itemStore.meta.page === 1" class="pagination-btn-standard">
              Kembali
            </button>
            <button @click="changePage(itemStore.meta.page + 1)" :disabled="itemStore.meta.page === itemStore.meta.total_pages" class="pagination-btn-standard">
              Lanjut
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Modals -->
    <ItemDetailModal :is-open="isDetailModalOpen" :item-data="selectedItem" @close="closeDetailModal" />
    <ItemFormModal :is-open="isModalOpen" :item-data="selectedItem" @close="closeModal" @saved="fetchData" />
    <ItemQrModal :is-open="isQrModalOpen" :item-data="selectedItem" @close="closeQrModal" />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import api from '../utils/api'
import { useItemStore } from '../stores/item'
import { useCategoryStore } from '../stores/category'
import { useLocationStore } from '../stores/location'
import ItemFormModal from '../components/ItemFormModal.vue'
import ItemQrModal from '../components/ItemQrModal.vue'
import ItemDetailModal from '../components/ItemDetailModal.vue'

const router = useRouter()
const authStore = useAuthStore()
const itemStore = useItemStore()
const categoryStore = useCategoryStore()
const locationStore = useLocationStore()

const isModalOpen = ref(false)
const isQrModalOpen = ref(false)
const isDetailModalOpen = ref(false)
const selectedItem = ref(null)
const selectedIds = ref([])

const isAllSelected = computed(() => {
  if (itemStore.items.length === 0) return false
  return itemStore.items.every(item => selectedIds.value.includes(item.id))
})

const toggleSelectAll = (e) => {
  const currentPageIds = itemStore.items.map(i => i.id)
  if (e.target.checked) {
    const uniqueBatch = currentPageIds.filter(id => !selectedIds.value.includes(id))
    selectedIds.value = [...selectedIds.value, ...uniqueBatch]
  } else {
    selectedIds.value = selectedIds.value.filter(id => !currentPageIds.includes(id))
  }
}

const goToPrint = () => {
  if (selectedIds.value.length === 0) return
  router.push({ path: '/print-labels', query: { ids: selectedIds.value.join(',') } })
}

const handleExportExcel = async () => {
  try {
    const response = await api.get('/reports/export/items', { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `SIS-INV_Laporan-Barang_${new Date().toISOString().slice(0, 10)}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (err) { alert('Gagal mendownload laporan Excel') }
}

const downloadTemplate = async () => {
  try {
    const response = await api.get('/items/import/template', { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', 'SIS-INV_Template_Barang.xlsx')
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (err) { alert('Gagal mendownload template Excel') }
}

const handleImportExcel = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  const formData = new FormData()
  formData.append('file', file)
  
  itemStore.loading = true
  try {
    const { data } = await api.post('/items/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    if (data.success) {
      alert(`Berhasil mengimpor ${data.data.total} barang`)
      itemStore.fetchItems()
    }
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal mengimpor file excel. Pastikan format kolom sesuai: Kode, Nama, Kategori, Lokasi, Kondisi, Tipe Peminjam...')
  } finally {
    itemStore.loading = false
    event.target.value = ''
  }
}

const getStatusBadgeClass = (s) => {
  const m = {
    'AVAILABLE': 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400',
    'BORROWED': 'bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400',
    'MAINTENANCE': 'bg-amber-50 dark:bg-amber-900/30 text-amber-600 dark:text-amber-400',
    'LOST': 'bg-rose-50 dark:bg-rose-900/30 text-rose-600 dark:text-rose-400'
  }
  return m[s] || 'bg-gray-50 text-gray-600'
}

const getConditionBadgeClass = (c) => {
  const m = {
    'GOOD': 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400',
    'DAMAGED': 'bg-rose-50 dark:bg-rose-900/30 text-rose-600 dark:text-rose-400',
    'IN_REPAIR': 'bg-amber-50 dark:bg-amber-900/30 text-amber-600 dark:text-amber-400'
  }
  return m[c] || 'bg-gray-50 text-gray-600'
}

const filters = ref({ search: '', status: '', category_id: '', location: '', condition: '', page: 1 })
let searchTimeout
const debouncedFetch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => { filters.value.page = 1; fetchData() }, 500)
}
const fetchData = async () => { await itemStore.fetchItems(filters.value) }
const changePage = (newPage) => { filters.value.page = newPage; fetchData() }

const openAddModal = () => { selectedItem.value = null; isModalOpen.value = true }
const openEditModal = (item) => { selectedItem.value = { ...item }; isModalOpen.value = true }
const openQrModal = (item) => { selectedItem.value = { ...item }; isQrModalOpen.value = true }
const openDetailModal = (item) => { selectedItem.value = { ...item }; isDetailModalOpen.value = true }
const closeModal = () => { isModalOpen.value = false; selectedItem.value = null }
const closeQrModal = () => { isQrModalOpen.value = false; selectedItem.value = null }
const closeDetailModal = () => { isDetailModalOpen.value = false; selectedItem.value = null }

const deleteItem = async (item) => {
  if (confirm(`Hapus ${item.name}?`)) {
    try { await itemStore.deleteItem(item.id); fetchData() } 
    catch (e) { alert(itemStore.error || 'Gagal menghapus barang') }
  }
}

const getStatusLabel = (s) => ({ 'AVAILABLE': 'Tersedia', 'BORROWED': 'Dipinjam', 'MAINTENANCE': 'Perbaikan', 'LOST': 'Hilang' }[s] || s)
const getConditionLabel = (c) => ({ 'GOOD': 'Baik', 'DAMAGED': 'Rusak', 'IN_REPAIR': 'Diperbaiki' }[c] || c)

onMounted(() => {
  categoryStore.fetchCategories()
  locationStore.fetchLocations()
  fetchData()
})
</script>
