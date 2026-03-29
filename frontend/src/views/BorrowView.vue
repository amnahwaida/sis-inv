<template>
  <div class="animate-fade-in max-w-2xl mx-auto space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">Pinjam Barang</h1>
    </div>

    <div class="card space-y-6">
      <div v-if="!scannedItemCode">
        <p class="text-gray-600 mb-4 text-center">Scan QR Code pada label barang untuk meminjam, atau masukkan kode barang secara manual.</p>
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

        <form @submit.prevent="submitBorrow" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Lama Peminjaman (Hari) *</label>
            <select v-model="form.expected_return_days" class="input-field w-full">
              <option :value="1">1 Hari (Besok)</option>
              <option :value="3">3 Hari</option>
              <option :value="7">7 Hari (1 Minggu)</option>
              <option :value="14">14 Hari (2 Minggu)</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Tujuan / Keperluan *</label>
            <textarea v-model="form.notes" required rows="3" class="input-field w-full" placeholder="CTH: Digunakan untuk praktikum kelas 10A"></textarea>
          </div>
          
          <div v-if="transactionStore.error" class="p-3 bg-red-50 text-red-700 text-sm rounded-lg flex items-center gap-2">
            <svg class="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
            </svg>
            {{ transactionStore.error }}
          </div>

          <div v-if="successMsg" class="p-4 bg-green-50 justify-between items-center text-green-800 rounded-lg flex gap-3 border border-green-200">
            <div>
              <p class="font-bold">✅ Berhasil Dipinjam!</p>
              <p class="text-sm mt-1">{{ successMsg }}</p>
            </div>
            <router-link to="/my-borrowings" class="btn-primary text-sm whitespace-nowrap">Lihat Peminjaman</router-link>
          </div>

          <div class="pt-4 border-t border-gray-100 flex gap-3">
            <button type="submit" :disabled="transactionStore.loading || !!successMsg" class="flex-1 btn-primary justify-center">
              {{ transactionStore.loading ? 'Memproses...' : 'Pinjam Sekarang' }}
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
const successMsg = ref('')

const form = ref({
  expected_return_days: 7,
  notes: ''
})

const handleScan = (code) => {
  scannedItemCode.value = code
  successMsg.value = ''
  transactionStore.error = null
}

const resetScan = () => {
  scannedItemCode.value = ''
  successMsg.value = ''
  form.value.notes = ''
}

const submitBorrow = async () => {
  try {
    const res = await transactionStore.borrowItem(
      scannedItemCode.value, 
      form.value.expected_return_days, 
      form.value.notes
    )
    successMsg.value = `Barang ${res.item_name} telah ditambahkan ke daftar pinjaman Anda.`
    // Don't auto redirect immediately, let user read success message
  } catch (err) {
    // Error is handled in store
  }
}
</script>
