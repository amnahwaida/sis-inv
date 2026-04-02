<template>
  <div class="print-container bg-white min-h-screen p-8">
    <!-- Screen Header (Hidden on Print) -->
    <div class="no-print mb-8 flex items-center justify-between bg-gray-50 p-4 rounded-xl border border-gray-100">
      <div>
        <h1 class="text-xl font-bold text-gray-900">Preview Cetak Label QR</h1>
        <p class="text-sm text-gray-500">Siapkan printer Anda. Gunakan kertas A4 untuk hasil terbaik.</p>
      </div>
      <div class="flex gap-3">
        <button @click="$router.back()" class="btn-secondary px-5 py-2 !rounded-lg h-auto">
          Kembali
        </button>
        <button @click="printLabels" class="btn-premium-action !py-2 !px-5 !rounded-lg flex items-center gap-2">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6.72 13.829c-.24.03-.48.062-.72.096m.72-.096a42.415 42.415 0 0110.56 0m-10.56 0L6.34 18m10.94-4.171c.24.03.48.062.72.096m-.72-.096L17.66 18m0 0l.229 2.523a1.125 1.125 0 01-1.12 1.227H7.231c-.618 0-1.11-.51-1.107-1.127L6.34 18m11.32 0h-11.32m0 0l1.107-1.107m11.32 1.107l-1.107-1.107m1.107 0c.34-.034.67-.07.999-.111m-1.339.111l-.39-4.881a3 3 0 013-3.238h.75m-6.75 3.238L12 3m0 0l-3.75 3.238M12 3v10.5" />
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
