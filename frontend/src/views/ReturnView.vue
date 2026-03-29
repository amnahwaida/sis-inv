<template>
  <div class="animate-fade-in max-w-2xl mx-auto space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">Kembalikan Barang</h1>
    </div>

    <div class="card space-y-6">
      <div v-if="!scannedItemCode">
        <p class="text-gray-600 mb-4 text-center">Scan QR Code pada label barang untuk mengembalikan, atau masukkan kode barang secara manual.</p>
        <QrScanner :active="true" @scanned="handleScan" />
      </div>

      <div v-else class="space-y-6 animate-fade-in">
        <div class="bg-indigo-50 border border-indigo-100 rounded-lg p-5 flex items-center justify-between">
          <div>
            <p class="text-xs font-semibold text-indigo-800 uppercase tracking-wider mb-1">Barang Terpilih</p>
            <p class="text-xl font-bold text-indigo-900 font-mono">{{ scannedItemCode }}</p>
          </div>
          <button @click="resetScan" class="text-sm font-medium text-indigo-600 hover:text-indigo-800 underline">Ganti Barang</button>
        </div>

        <form @submit.prevent="submitReturn" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Kondisi Saat Dikembalikan *</label>
            <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 mt-2">
              <label class="border rounded-lg p-3 text-center cursor-pointer flex flex-col items-center justify-center gap-1 transition-colors"
                     :class="form.condition === 'GOOD' ? 'bg-indigo-50 border-indigo-500 text-indigo-700' : 'border-gray-200 hover:bg-gray-50 text-gray-700'">
                <input type="radio" v-model="form.condition" value="GOOD" class="sr-only">
                <span class="text-2xl">✅</span>
                <span class="text-xs font-medium">Baik</span>
              </label>
              
              <label class="border rounded-lg p-3 text-center cursor-pointer flex flex-col items-center justify-center gap-1 transition-colors"
                     :class="form.condition === 'DAMAGED' ? 'bg-orange-50 border-orange-500 text-orange-700' : 'border-gray-200 hover:bg-gray-50 text-gray-700'">
                <input type="radio" v-model="form.condition" value="DAMAGED" class="sr-only">
                <span class="text-2xl">⚠️</span>
                <span class="text-xs font-medium">Rusak</span>
              </label>

              <label class="border rounded-lg p-3 text-center cursor-pointer flex flex-col items-center justify-center gap-1 transition-colors"
                     :class="form.condition === 'IN_REPAIR' ? 'bg-yellow-50 border-yellow-500 text-yellow-700' : 'border-gray-200 hover:bg-gray-50 text-gray-700'">
                <input type="radio" v-model="form.condition" value="IN_REPAIR" class="sr-only">
                <span class="text-2xl">🔧</span>
                <span class="text-xs font-medium">Perbaikan</span>
              </label>

              <label class="border rounded-lg p-3 text-center cursor-pointer flex flex-col items-center justify-center gap-1 transition-colors"
                     :class="form.condition === 'LOST' ? 'bg-red-50 border-red-500 text-red-700' : 'border-gray-200 hover:bg-gray-50 text-gray-700'">
                <input type="radio" v-model="form.condition" value="LOST" class="sr-only">
                <span class="text-2xl">❌</span>
                <span class="text-xs font-medium">Hilang</span>
              </label>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Catatan Tambahan</label>
            <textarea v-model="form.notes" rows="3" class="input-field w-full" placeholder="Kondisi fisik barang saat dikembalikan..."></textarea>
          </div>
          
          <div v-if="transactionStore.error" class="p-3 bg-red-50 text-red-700 text-sm rounded-lg flex items-center gap-2">
            <svg class="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
            </svg>
            {{ transactionStore.error }}
          </div>

          <div v-if="successMsg" class="p-4 bg-green-50 justify-between items-center text-green-800 rounded-lg flex gap-3 border border-green-200">
            <div>
              <p class="font-bold">✅ Berhasil Dikembalikan!</p>
              <p class="text-sm mt-1">Status dan kondisi barang telah diperbarui.</p>
            </div>
            <router-link to="/my-borrowings" class="btn-primary text-sm whitespace-nowrap">Lihat Peminjaman</router-link>
          </div>

          <div class="pt-4 border-t border-gray-100 flex gap-3">
            <button type="submit" :disabled="transactionStore.loading || !!successMsg" class="flex-1 btn-primary justify-center">
              {{ transactionStore.loading ? 'Memproses...' : 'Selesaikan Pengembalian' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useTransactionStore } from '../stores/transaction'
import QrScanner from '../components/QrScanner.vue'

const transactionStore = useTransactionStore()
const router = useRouter()

const scannedItemCode = ref('')
const successMsg = ref(false)

const form = ref({
  condition: 'GOOD',
  notes: ''
})

const handleScan = (code) => {
  scannedItemCode.value = code
  successMsg.value = false
  transactionStore.error = null
}

const resetScan = () => {
  scannedItemCode.value = ''
  successMsg.value = false
  form.value.condition = 'GOOD'
  form.value.notes = ''
}

const submitReturn = async () => {
  try {
    const res = await transactionStore.returnItem(
      scannedItemCode.value, 
      form.value.condition, 
      form.value.notes
    )
    successMsg.value = true
  } catch (err) {
    // Error is handled in store
  }
}
</script>
