<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-0">
    <div class="fixed inset-0 bg-gray-900/50 transition-opacity backdrop-blur-sm" @click="$emit('close')"></div>
    
    <div class="bg-white rounded-xl shadow-xl transform transition-all sm:max-w-2xl sm:w-full max-h-[90vh] flex flex-col relative z-10 overflow-hidden">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between bg-gray-50/50">
        <div>
          <h3 class="text-lg font-bold text-gray-900">Detail Barang</h3>
          <p class="text-sm text-gray-500 font-mono mt-0.5">{{ itemData?.code }}</p>
        </div>
        <button @click="$emit('close')" class="text-gray-400 hover:bg-gray-200 hover:text-gray-600 rounded-lg p-1.5 transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Body -->
      <div class="p-6 overflow-y-auto">
        <div v-if="itemData" class="space-y-6">
          
          <!-- Basic Info Card -->
          <div class="bg-white border text-sm border-gray-100 rounded-xl p-5 shadow-sm space-y-4">
            <div class="flex items-center gap-4 border-b border-gray-50 pb-4">
              <div class="w-12 h-12 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center font-bold text-xl uppercase">
                {{ itemData.name.charAt(0) }}
              </div>
              <div>
                <h4 class="text-xl font-bold text-gray-900">{{ itemData.name }}</h4>
                <span class="inline-flex items-center px-2 py-0.5 mt-1 rounded text-xs font-medium" :class="{
                  'bg-green-100 text-green-800': itemData.status === 'AVAILABLE',
                  'bg-blue-100 text-blue-800': itemData.status === 'BORROWED',
                  'bg-purple-100 text-purple-800': itemData.status === 'MAINTENANCE',
                  'bg-red-100 text-red-800': itemData.status === 'LOST'
                }">
                  Status: {{ getStatusLabel(itemData.status) }}
                </span>
                <span class="inline-flex items-center px-2 py-0.5 mt-1 ml-2 rounded text-xs font-medium" :class="{
                  'bg-gray-100 text-gray-800': itemData.condition === 'GOOD',
                  'bg-orange-100 text-orange-800': itemData.condition === 'DAMAGED',
                  'bg-yellow-100 text-yellow-800': itemData.condition === 'IN_REPAIR'
                }">
                  Kondisi: {{ getConditionLabel(itemData.condition) }}
                </span>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
              <div>
                <span class="block text-xs font-medium text-gray-500 uppercase">Kategori</span>
                <span class="block mt-1 font-medium text-gray-900">{{ itemData.category_name || '-' }}</span>
              </div>
              <div>
                <span class="block text-xs font-medium text-gray-500 uppercase">Lokasi Asal</span>
                <span class="block mt-1 font-medium text-gray-900">{{ itemData.location || '-' }}</span>
              </div>
              <div>
                <span class="block text-xs font-medium text-gray-500 uppercase">Aturan Peminjaman</span>
                <span class="block mt-1 font-medium text-gray-900">{{ itemData.borrower_type === 'STUDENT_ALLOWED' ? 'Bisa untuk Siswa' : 'Hanya Guru/Staff' }}</span>
              </div>
              <div v-if="itemData.purchase_date">
                <span class="block text-xs font-medium text-gray-500 uppercase">Tanggal Pembelian</span>
                <span class="block mt-1 font-medium text-gray-900">{{ formatDate(itemData.purchase_date) }}</span>
              </div>
              <div v-if="itemData.warranty_end_date">
                <span class="block text-xs font-medium text-gray-500 uppercase">Akhir Garansi</span>
                <span class="block mt-1 font-medium text-gray-900">{{ formatDate(itemData.warranty_end_date) }}</span>
              </div>
            </div>
            
            <div v-if="itemData.notes" class="pt-2">
              <span class="block text-xs font-medium text-gray-500 uppercase mb-1">Catatan Tambahan</span>
              <p class="text-gray-700 bg-gray-50 p-3 rounded-lg text-sm">{{ itemData.notes }}</p>
            </div>
          </div>

          <!-- History Section -->
          <div class="space-y-3">
            <h4 class="text-sm font-bold text-gray-900 uppercase tracking-wider flex items-center gap-2">
              <svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Riwayat Peminjaman
            </h4>
            
            <div v-if="loadingHistory" class="flex justify-center py-8">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
            </div>
            
            <div v-else-if="history.length === 0" class="text-center py-8 bg-gray-50 rounded-xl border border-dashed border-gray-200">
              <p class="text-sm text-gray-500">Belum ada riwayat peminjaman untuk barang ini.</p>
            </div>

            <div v-else class="overflow-hidden border border-gray-100 rounded-xl shadow-sm">
              <table class="w-full text-sm text-left table-fixed">
                <thead class="bg-gray-50 text-gray-600 font-medium border-b border-gray-100">
                  <tr>
                    <th class="px-4 py-2.5 w-auto">Peminjam / Ket.</th>
                    <th class="px-4 py-2.5 w-24">Status</th>
                    <th class="px-4 py-2.5 w-36 text-right">Tanggal</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-100 bg-white">
                  <tr v-for="h in history" :key="h.id" class="hover:bg-gray-50/50 transition-colors align-top">
                    <td class="px-4 py-3 min-w-0">
                      <div class="flex flex-col">
                        <div class="font-bold text-gray-900 truncate" :title="h.borrower_name">
                          {{ h.borrower_name }}
                          <span v-if="h.student_class" class="text-[10px] font-normal text-gray-400 ml-1">({{ h.student_class }})</span>
                        </div>
                        <div class="flex items-center gap-1.5 mt-0.5">
                          <span class="text-[9px] font-black uppercase tracking-tighter px-1 rounded border" 
                                :class="h.borrower_type === 'STAFF' ? 'bg-indigo-50 text-indigo-600 border-indigo-100' : 'bg-amber-50 text-amber-600 border-amber-100'">
                            {{ h.borrower_type }}
                          </span>
                          <span v-if="h.borrower_type === 'STUDENT'" class="text-[9px] text-gray-400 flex items-center gap-1">
                            Fac: <span class="font-bold italic">{{ h.staff_name }}</span>
                          </span>
                        </div>
                      </div>
                      <div v-if="h.purpose || h.return_notes" class="mt-1.5 leading-relaxed">
                        <div class="text-xs space-y-1" :class="{'line-clamp-2': !isExpanded(h.id)}">
                          <p v-if="h.purpose" class="text-gray-500 italic break-words">
                            <span class="font-bold not-italic text-[10px] text-gray-400 uppercase mr-1">Pinjam:</span>
                            "{{ h.purpose }}"
                          </p>
                          <p v-if="h.return_notes" class="text-gray-600 italic break-words border-t border-gray-50 pt-1">
                            <span class="font-bold not-italic text-[10px] text-gray-400 uppercase mr-1">Kembali:</span>
                            "{{ h.return_notes }}"
                          </p>
                        </div>
                        <!-- Inline Toggle -->
                        <button 
                          v-if="(h.purpose?.length || 0) + (h.return_notes?.length || 0) > 40" 
                          @click="toggleExpand(h.id)" 
                          class="text-[10px] text-indigo-500 hover:text-indigo-700 font-bold mt-1 block"
                        >
                          {{ isExpanded(h.id) ? 'Sembunyikan' : 'Selengkapnya...' }}
                        </button>
                      </div>
                    </td>
                    <td class="px-4 py-3">
                      <div class="flex flex-col gap-1.5">
                        <span class="inline-flex px-2 py-0.5 rounded-full text-[10px] font-bold uppercase w-fit" :class="{
                          'bg-blue-100 text-blue-700': h.status === 'ACTIVE',
                          'bg-green-100 text-green-700': h.status === 'RETURNED',
                          'bg-red-100 text-red-700': h.status === 'OVERDUE'
                        }">
                          {{ h.status === 'ACTIVE' ? 'DIPINJAM' : h.status === 'RETURNED' ? 'KEMBALI' : 'TERLAMBAT' }}
                        </span>
                        <!-- Return Condition Sub-info -->
                        <span v-if="h.status === 'RETURNED' && h.return_condition" class="text-[9px] font-medium text-gray-500 flex items-center gap-1">
                          <span class="w-1.5 h-1.5 rounded-full" :class="{
                            'bg-green-500': h.return_condition === 'GOOD',
                            'bg-red-500': h.return_condition === 'DAMAGED',
                            'bg-yellow-500': h.return_condition === 'IN_REPAIR',
                            'bg-gray-900': h.return_condition === 'LOST'
                          }"></span>
                          {{ getConditionLabel(h.return_condition) }}
                        </span>
                      </div>
                    </td>
                    <td class="px-4 py-3 text-right tabular-nums whitespace-nowrap">
                      <div class="text-xs font-semibold text-gray-800">{{ formatDateShort(h.borrowed_at) }}</div>
                      <div v-if="h.returned_at" class="mt-1 text-[10px] text-gray-500 flex flex-col">
                        <span class="text-[9px] uppercase tracking-tighter text-gray-400">Kembali pada:</span>
                        {{ formatDateShort(h.returned_at) }}
                      </div>
                      <div v-else class="mt-1 text-[10px] font-medium text-blue-600 animate-pulse">Belum kembali</div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          
        </div>
      </div>

      <!-- Footer Actions -->
      <div class="px-6 py-4 bg-gray-50 border-t border-gray-100 flex justify-end gap-3">
        <button type="button" @click="$emit('close')" class="btn-secondary">Tutup</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useItemStore } from '../stores/item'

const props = defineProps({
  isOpen: Boolean,
  itemData: Object
})

defineEmits(['close'])

const itemStore = useItemStore()
const history = ref([])
const loadingHistory = ref(false)
const expandedHistoryIds = ref(new Set())

watch(() => props.isOpen, async (newVal) => {
  if (newVal && props.itemData?.id) {
    expandedHistoryIds.value = new Set() // Reset on open
    loadHistory()
  }
})

const toggleExpand = (id) => {
  if (expandedHistoryIds.value.has(id)) {
    expandedHistoryIds.value.delete(id)
  } else {
    expandedHistoryIds.value.add(id)
  }
}

const isExpanded = (id) => expandedHistoryIds.value.has(id)

const loadHistory = async () => {
  loadingHistory.value = true
  try {
    const data = await itemStore.fetchItemHistory(props.itemData.id)
    history.value = data
  } catch (e) {
    console.error('Failed to load history', e)
  } finally {
    loadingHistory.value = false
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

const formatDate = (isoString) => {
  if (!isoString) return '-'
  const date = new Date(isoString)
  return new Intl.DateTimeFormat('id-ID', { 
    day: 'numeric', month: 'long', year: 'numeric' 
  }).format(date)
}

const formatDateShort = (isoString) => {
  if (!isoString) return '-'
  const date = new Date(isoString)
  return new Intl.DateTimeFormat('id-ID', { 
    day: 'numeric', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  }).format(date)
}
</script>
