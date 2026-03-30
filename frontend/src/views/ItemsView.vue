<template>
  <div class="animate-fade-in space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">Daftar Barang</h1>
      <div class="flex gap-3">
        <button 
          v-if="authStore.isAdmin" 
          @click="goToPrint" 
          :disabled="selectedIds.length === 0"
          class="btn-secondary flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
          </svg>
          Cetak Label ({{ selectedIds.length }})
        </button>
        <button v-if="authStore.isAdmin" @click="handleExportExcel" class="btn-secondary bg-green-50 text-green-700 border-green-200 hover:bg-green-100 flex items-center gap-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Export Excel
        </button>
        <button v-if="authStore.isAdmin" @click="openAddModal" class="btn-primary flex items-center gap-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Tambah Barang
        </button>
      </div>
    </div>

    <!-- Search & Filter -->
    <div class="card">
      <div class="flex flex-col md:flex-row gap-3">
        <div class="flex-1">
          <input type="text" class="input-field" placeholder="Cari nama / kode barang..." v-model="filters.search" @input="debouncedFetch" />
        </div>
        <select class="input-field md:w-36" v-model="filters.status" @change="fetchData">
          <option value="">Semua Status</option>
          <option value="AVAILABLE">Tersedia</option>
          <option value="BORROWED">Dipinjam</option>
          <option value="MAINTENANCE">Perbaikan</option>
          <option value="LOST">Hilang</option>
        </select>
        <select class="input-field md:w-36" v-model="filters.category_id" @change="fetchData">
          <option value="">Semua Kategori</option>
          <option v-for="cat in categoryStore.categories" :key="cat.id" :value="cat.id">
            {{ cat.name }}
          </option>
        </select>
        <select class="input-field md:w-36" v-model="filters.location_id" @change="fetchData">
          <option value="">Semua Lokasi</option>
          <option v-for="loc in locationStore.locations" :key="loc.id" :value="loc.id">
            {{ loc.name }}
          </option>
        </select>
        <select class="input-field md:w-36" v-model="filters.condition" @change="fetchData">
          <option value="">Semua Kondisi</option>
          <option value="GOOD">Baik</option>
          <option value="DAMAGED">Rusak</option>
          <option value="IN_REPAIR">Dalam Perbaikan</option>
        </select>
      </div>
    </div>

    <!-- Items Table -->
    <div class="card p-0 overflow-hidden relative min-h-[300px]">
      <!-- Loading Overlay -->
      <div v-if="itemStore.loading" class="absolute inset-0 bg-white/80 z-10 flex flex-col items-center justify-center">
        <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-indigo-600"></div>
        <span class="mt-2 text-sm text-gray-600">Memuat data...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th v-if="authStore.isAdmin" class="w-10 px-6 py-3">
                <input type="checkbox" @change="toggleSelectAll" :checked="isAllSelected" class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500" />
              </th>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Kode / Nama</th>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Kategori & Lokasi</th>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Status</th>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Kondisi</th>
              <th class="text-right text-xs font-semibold text-gray-600 uppercase px-6 py-3">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <template v-if="itemStore.items.length > 0">
              <tr v-for="item in itemStore.items" :key="item.id" class="hover:bg-gray-50 transition-colors" :class="{'bg-indigo-50/50': selectedIds.includes(item.id)}">
                <td v-if="authStore.isAdmin" class="px-6 py-4">
                  <input type="checkbox" v-model="selectedIds" :value="item.id" class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500" />
                </td>
                <td class="px-6 py-4">
                  <div class="font-medium text-gray-900">{{ item.name }}</div>
                  <div class="text-sm text-gray-500 font-mono">{{ item.code }}</div>
                </td>
                <td class="px-6 py-4">
                  <div class="text-sm text-gray-900">{{ item.category_name || '-' }}</div>
                  <div class="text-sm text-gray-500">{{ item.location_name || item.location || '-' }}</div>
                </td>
                <td class="px-6 py-4">
                  <span class="px-2.5 py-1 rounded-full text-xs font-medium" :class="{
                    'bg-green-100 text-green-800': item.status === 'AVAILABLE',
                    'bg-blue-100 text-blue-800': item.status === 'BORROWED',
                    'bg-purple-100 text-purple-800': item.status === 'MAINTENANCE',
                    'bg-red-100 text-red-800': item.status === 'LOST'
                  }">
                    {{ getStatusLabel(item.status) }}
                  </span>
                </td>
                <td class="px-6 py-4">
                  <span class="px-2.5 py-1 rounded-full text-xs font-medium" :class="{
                    'bg-gray-100 text-gray-800': item.condition === 'GOOD',
                    'bg-orange-100 text-orange-800': item.condition === 'DAMAGED',
                    'bg-yellow-100 text-yellow-800': item.condition === 'IN_REPAIR'
                  }">
                    {{ getConditionLabel(item.condition) }}
                  </span>
                </td>
                <td class="px-6 py-4 text-right border-l border-gray-50">
                  <div class="flex items-center justify-end gap-3 flex-wrap">
                    <!-- Standard Actions -->
                    <button @click="openDetailModal(item)" class="text-indigo-600 hover:text-indigo-900 text-sm font-medium">Detail</button>
                    
                    <!-- Admin Actions -->
                    <template v-if="authStore.isAdmin">
                      <button @click="openQrModal(item)" class="text-emerald-600 hover:text-emerald-900 text-sm font-medium">QR Code</button>
                      <button @click="openEditModal(item)" class="text-blue-600 hover:text-blue-900 text-sm font-medium">Edit</button>
                      <button @click="deleteItem(item)" class="text-red-600 hover:text-red-900 text-sm font-medium">Hapus</button>
                    </template>
                  </div>
                </td>
              </tr>
            </template>
            <tr v-else-if="!itemStore.loading" class="text-center">
              <td colspan="5" class="px-6 py-12 text-gray-500">
                <p class="text-3xl mb-2">📦</p>
                <p class="text-base font-medium text-gray-900">Belum ada data barang</p>
                <p class="text-sm text-gray-500 mt-1">Coba sesuaikan filter atau tambah barang baru.</p>
              </td>
            </tr>
          </tbody>
        </table>
        
        <!-- Pagination -->
        <div v-if="itemStore.meta.total_pages > 1" class="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
          <span class="text-sm text-gray-700">
            Menampilkan <span class="font-medium">{{ (itemStore.meta.page - 1) * itemStore.meta.page_size + 1 }}</span> sampai <span class="font-medium">{{ Math.min(itemStore.meta.page * itemStore.meta.page_size, itemStore.meta.total) }}</span> dari <span class="font-medium">{{ itemStore.meta.total }}</span> hasil
          </span>
          <div class="flex gap-2">
            <button 
              @click="changePage(itemStore.meta.page - 1)" 
              :disabled="itemStore.meta.page === 1"
              class="px-3 py-1 border border-gray-300 rounded-md text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
            >
              Sebelumnya
            </button>
            <button 
              @click="changePage(itemStore.meta.page + 1)" 
              :disabled="itemStore.meta.page === itemStore.meta.total_pages"
              class="px-3 py-1 border border-gray-300 rounded-md text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
            >
              Selanjutnya
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
    // Add only if not already there
    const uniqueBatch = currentPageIds.filter(id => !selectedIds.value.includes(id))
    selectedIds.value = [...selectedIds.value, ...uniqueBatch]
  } else {
    // Remove only items from current page
    selectedIds.value = selectedIds.value.filter(id => !currentPageIds.includes(id))
  }
}

const goToPrint = () => {
  if (selectedIds.value.length === 0) return
  router.push({
    path: '/print-labels',
    query: { ids: selectedIds.value.join(',') }
  })
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
  } catch (err) {
    alert('Gagal mendownload laporan Excel')
  }
}

const filters = ref({
  search: '',
  status: '',
  category_id: '',
  location_id: '',
  condition: '',
  page: 1
})

let searchTimeout
const debouncedFetch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    filters.value.page = 1
    fetchData()
  }, 500)
}

const fetchData = async () => {
  await itemStore.fetchItems(filters.value)
}

const changePage = (newPage) => {
  filters.value.page = newPage
  fetchData()
}

const openAddModal = () => {
  selectedItem.value = null
  isModalOpen.value = true
}

const openEditModal = (item) => {
  selectedItem.value = { ...item }
  isModalOpen.value = true
}

const openQrModal = (item) => {
  selectedItem.value = { ...item }
  isQrModalOpen.value = true
}

const openDetailModal = (item) => {
  selectedItem.value = { ...item }
  isDetailModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
  selectedItem.value = null
}

const closeQrModal = () => {
  isQrModalOpen.value = false
  selectedItem.value = null
}

const closeDetailModal = () => {
  isDetailModalOpen.value = false
  selectedItem.value = null
}

const deleteItem = async (item) => {
  if (confirm(`Apakah Anda yakin ingin menghapus barang ${item.name} (${item.code})?`)) {
    try {
      await itemStore.deleteItem(item.id)
      fetchData() // Refresh list after delete
    } catch (e) {
      alert(itemStore.error || 'Gagal menghapus barang')
    }
  }
}

const getStatusLabel = (status) => {
  const map = {
    'AVAILABLE': 'Tersedia',
    'BORROWED': 'Dipinjam',
    'MAINTENANCE': 'Perbaikan',
    'LOST': 'Hilang'
  }
  return map[status] || status
}

const getConditionLabel = (condition) => {
  const map = {
    'GOOD': 'Baik',
    'DAMAGED': 'Rusak',
    'IN_REPAIR': 'Diperbaiki'
  }
  return map[condition] || condition
}

onMounted(() => {
  categoryStore.fetchCategories()
  locationStore.fetchLocations()
  fetchData()
})
</script>
