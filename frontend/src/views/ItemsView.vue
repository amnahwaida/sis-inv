<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-primary-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col lg:flex-row lg:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Daftar Barang</h1>
          <p class="text-primary-100/70 text-sm font-medium">Manajemen seluruh aset dan inventaris sekolah</p>
        </div>
        
        <div class="flex flex-wrap items-center gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10">
          <button v-if="authStore.isAdmin" @click="goToPrint" :disabled="selectedIds.length === 0"
                  class="bg-primary-500/20 hover:bg-primary-500/30 text-primary-400 px-4 py-3 rounded-xl text-[10px] font-black transition-all flex items-center gap-2 active:scale-95 disabled:opacity-30">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6.72 13.829c-.24.03-.48.062-.72.096m.72-.096a42.415 42.415 0 0110.56 0m-10.56 0L6.34 18m10.94-4.171c.24.03.48.062.72.096m-.72-.096L17.66 18m0 0l.229 2.523a1.125 1.125 0 01-1.12 1.227H7.231c-.618 0-1.113-.493-1.12-1.112L6.34 18m11.32 0c.068 0 .135-.01.2-.03m-11.51 0a.2.2 0 00.2.03m11.31 0h.375c.621 0 1.125-.504 1.125-1.125v-4.135c0-.621-.504-1.125-1.125-1.125h-1.125M6.34 18h-.375A1.125 1.125 0 014.84 16.875v-4.135c0-.621.504-1.125 1.125-1.125h1.125M16.5 7.5h.375c.621 0 1.125-.504 1.125-1.125V4.125c0-.621-.504-1.125-1.125-1.125H6.375c-.621 0-1.125.504-1.125 1.125v2.25c0 .621.504 1.125 1.125 1.125h1.125m10.5 0h-10.5" />
            </svg>
            CETAK ({{ selectedIds.length }})
          </button>
          <button v-if="authStore.isAdmin" @click="handleExportExcel" 
                  class="bg-emerald-500/20 hover:bg-emerald-500/30 text-emerald-400 px-4 py-3 rounded-xl text-[10px] font-black transition-all flex items-center gap-2 active:scale-95">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m.75 12l3 3m0 0l3-3m-3 3v-6m-1.5-9H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
            </svg>
            EXCEL
          </button>
          <button v-if="authStore.isAdmin" @click="$refs.fileInput.click()" 
                  class="bg-primary-500/20 hover:bg-primary-500/30 text-primary-400 px-4 py-3 rounded-xl text-[10px] font-black transition-all flex items-center gap-2 active:scale-95">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5" />
            </svg>
            IMPORT
            <input type="file" ref="fileInput" @change="handleImportExcel" class="hidden" accept=".xlsx, .xls" />
          </button>
          <button v-if="authStore.isAdmin" @click="downloadTemplate" 
                  class="bg-gray-500/20 hover:bg-gray-500/30 text-gray-400 px-4 py-3 rounded-xl text-[10px] font-black transition-all flex items-center gap-2 active:scale-95">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 3h5.25m-5.25 3h5.25m-5.25 3h5.25m-10.5-9H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
            </svg>
            TEMPLATE
          </button>
          <button v-if="authStore.isAdmin" @click="openAddModal" class="btn-premium-primary">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
            </svg>
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
          <svg class="w-5 h-5 absolute left-3 top-3.5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
          </svg>
        </div>
        <select class="input-field h-12 rounded-2xl appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2020%2020%22%3E%3Cpath%20stroke%3D%22%236b7280%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%221.5%22%20d%3D%22m6%208%204%204%204-4%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[right_0.5rem_center] bg-no-repeat pr-10" v-model="filters.status" @change="fetchData">
          <option value="">Semua Status</option>
          <option value="AVAILABLE">Tersedia</option>
          <option value="BORROWED">Dipinjam</option>
          <option value="MAINTENANCE">Perbaikan</option>
          <option value="LOST">Hilang</option>
        </select>
        <select class="input-field h-12 rounded-2xl appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2020%2020%22%3E%3Cpath%20stroke%3D%22%236b7280%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%221.5%22%20d%3D%22m6%208%204%204%204-4%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[right_0.5rem_center] bg-no-repeat pr-10" v-model="filters.category_id" @change="fetchData">
          <option value="">Semua Kategori</option>
          <option v-for="cat in categoryStore.categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
        </select>
        <select class="input-field h-12 rounded-2xl appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2020%2020%22%3E%3Cpath%20stroke%3D%22%236b7280%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%221.5%22%20d%3D%22m6%208%204%204%204-4%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[right_0.5rem_center] bg-no-repeat pr-10" v-model="filters.location" @change="fetchData">
          <option value="">Semua Lokasi</option>
          <option v-for="loc in locationStore.locations" :key="loc.id" :value="loc.name">{{ loc.name }}</option>
        </select>
        <select class="input-field h-12 rounded-2xl appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2020%2020%22%3E%3Cpath%20stroke%3D%22%236b7280%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%221.5%22%20d%3D%22m6%208%204%204%204-4%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[right_0.5rem_center] bg-no-repeat pr-10" v-model="filters.condition" @change="fetchData">
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
                      <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-5.25v9" />
                      </svg>
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
                      <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z" />
                        <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      </svg>
                    </button>
                    <template v-if="authStore.isAdmin">
                      <button @click="openQrModal(item)" class="p-2 text-emerald-500 bg-emerald-50 dark:bg-emerald-900/30 rounded-xl hover:bg-emerald-500 hover:text-white transition-all transform hover:scale-110 active:scale-95" title="QR Code">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 4.875c0-.621.504-1.125 1.125-1.125h4.5c.621 0 1.125.504 1.125 1.125v4.5c0 .621-.504 1.125-1.125 1.125h-4.5a1.125 1.125 0 01-1.125-1.125v-4.5zM3.75 14.625c0-.621.504-1.125 1.125-1.125h4.5c.621 0 1.125.504 1.125 1.125v4.5c0 .621-.504 1.125-1.125 1.125h-4.5a1.125 1.125 0 01-1.125-1.125v-4.5zM13.5 4.875c0-.621.504-1.125 1.125-1.125h4.5c.621 0 1.125.504 1.125 1.125v4.5c0 .621-.504 1.125-1.125 1.125h-4.5A1.125 1.125 0 0113.5 9.375v-4.5z" />
                          <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 6.75h.75v.75h-.75v-.75zM6.75 16.5h.75v.75h-.75v-.75zM16.5 6.75h.75v.75h-.75v-.75zM13.5 13.5h.75v.75h-.75v-.75zM13.5 19.5h.75v.75h-.75v-.75zM19.5 13.5h.75v.75h-.75v-.75zM19.5 19.5h.75v.75h-.75v-.75zM16.5 16.5h.75v.75h-.75v-.75z" />
                        </svg>
                      </button>
                      <button @click="openEditModal(item)" class="btn-action-edit" title="Edit">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897L16.863 4.487zm0 0L19.5 7.125" />
                        </svg>
                      </button>
                      <button @click="deleteItem(item)" class="btn-action-delete" title="Hapus">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                        </svg>
                      </button>
                    </template>
                  </div>
                </td>
              </tr>
            </template>
            <tr v-else-if="!itemStore.loading" class="text-center">
              <td colspan="6" class="px-8 py-24">
                <div class="bg-gray-50 dark:bg-gray-700/50 w-24 h-24 rounded-full flex items-center justify-center mx-auto mb-6 text-gray-300 dark:text-gray-600">
                  <svg class="w-12 h-12" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-5.25v9" />
                  </svg>
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
            <div class="flex gap-1">
              <button v-for="p in visiblePages" :key="p" @click="changePage(p)"
                      class="w-10 h-10 rounded-xl text-[11px] font-black transition-all shadow-sm active:scale-95 border"
                      :class="p === itemStore.meta.page ? 'bg-primary-600 text-white border-primary-600' : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:bg-primary-50 dark:hover:bg-gray-700'">
                {{ p }}
              </button>
            </div>
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
const changePage = (newPage) => { 
  if (newPage < 1 || newPage > itemStore.meta.total_pages) return
  filters.value.page = newPage; 
  fetchData() 
}

const visiblePages = computed(() => {
  const pages = []
  const current = itemStore.meta.page
  const total = itemStore.meta.total_pages
  const start = Math.max(1, current - 2)
  const end = Math.min(total, current + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

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
