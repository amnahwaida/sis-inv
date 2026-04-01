<template>
  <div class="print-container bg-white min-h-screen p-8">
    <!-- Screen Header (Hidden on Print) -->
    <div class="no-print mb-8 flex items-center justify-between bg-gray-50 p-4 rounded-xl border border-gray-100">
      <div>
        <h1 class="text-xl font-bold text-gray-900">Preview Cetak Label QR</h1>
        <p class="text-sm text-gray-500">Siapkan printer Anda. Gunakan kertas A4 untuk hasil terbaik.</p>
      </div>
      <div class="flex gap-3">
        <button @click="$router.back()" class="px-5 py-2 text-sm font-bold text-gray-600 hover:bg-gray-200 rounded-lg transition-colors">
          Kembali
        </button>
        <button @click="printLabels" class="px-5 py-2 text-sm font-bold text-white bg-primary-600 hover:bg-primary-700 rounded-lg shadow-lg shadow-primary-200 transition-all flex items-center gap-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
          </svg>
          Cetak Sekarang
        </button>
      </div>
    </div>

    <!-- Labels Grid -->
    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 print:grid-cols-3 print:gap-4 font-sans">
      <div v-for="item in items" :key="item.id" class="label-card border border-gray-200 p-4 rounded-md flex flex-col items-center text-center bg-white print:border-gray-500 print:rounded-none">
        <div class="w-full flex justify-between items-start mb-2 border-b border-gray-100 pb-2 print:border-gray-300">
          <span class="text-[8px] font-black uppercase text-gray-400">SIS-INV System</span>
          <span class="text-[8px] font-bold text-gray-900">{{ item.category_name || 'Aset' }}</span>
        </div>
        
        <div class="qr-placeholder mb-3 flex items-center justify-center bg-white p-1">
           <qrcode-vue :value="item.code" :size="120" level="H" :render-as="'svg'" />
        </div>

        <div class="space-y-0.5">
          <p class="text-[10px] font-black text-gray-900 uppercase leading-tight">{{ item.name }}</p>
          <p class="text-[12px] font-mono font-bold text-primary-700 tracking-tighter">{{ item.code }}</p>
          <div class="mt-2 text-[7px] text-gray-400 italic">
             Lokasi: {{ item.location_name || '-' }}
          </div>
        </div>
      </div>
    </div>

    <div v-if="items.length === 0 && !loading" class="text-center py-24 text-gray-400 italic">
      Tidak ada data barang untuk dicetak.
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useItemStore } from '../stores/item'
import QrcodeVue from 'qrcode.vue'

const route = useRoute()
const itemStore = useItemStore()
const items = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const selectedIds = route.query.ids ? route.query.ids.split(',') : []
    
    // We fetch items, then filter them by the passed IDs
    await itemStore.fetchItems({ page_size: 1000 }) // Load more just for printing
    
    if (selectedIds.length > 0) {
      items.value = itemStore.items.filter(item => selectedIds.includes(item.id.toString()))
    } else {
      items.value = itemStore.items
    }
  } catch (err) {
    console.error('Failed to load items for printing', err)
  } finally {
    loading.value = false
  }
})

const printLabels = () => {
  window.print()
}
</script>

<style scoped>
@media print {
  .no-print {
    display: none !important;
  }
  .print-container {
    padding: 0 !important;
    background: transparent !important;
  }
  .label-card {
    break-inside: avoid;
    page-break-inside: avoid;
    border: 1px solid #000 !important;
    margin-bottom: 0.5rem;
  }
  body {
    background: white;
  }
}

.label-card {
  min-height: 220px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
</style>
