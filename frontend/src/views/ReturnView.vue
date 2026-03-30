<template>
  <div class="animate-fade-in max-w-2xl mx-auto space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Kembalikan Barang</h1>
        <p class="text-sm text-gray-500 mt-1">Lakukan pengecekan kondisi barang sebelum dikembalikan ke rak.</p>
      </div>
    </div>

    <div class="card overflow-hidden transition-all duration-300">
      <div v-if="!scannedItemCode" class="p-4">
        <p class="text-gray-500 mb-6 text-center text-sm">Scan QR Code pada label barang untuk memproses pengembalian.</p>
        <QrScanner :active="true" @scanned="handleScan" />
      </div>

      <div v-else class="space-y-6 animate-scale-up p-1">
        <!-- Selected Item Badge -->
        <div class="bg-indigo-50 border border-indigo-100 rounded-2xl p-6 flex items-center justify-between shadow-sm">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-white rounded-xl shadow-sm flex items-center justify-center text-indigo-600">
              <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 15v-1a4 4 0 00-4-4H8m0 0l3 3m-3-3l3-3m9 14V5a2 2 0 00-2-2H6a2 2 0 00-2 2v16l4-2 4 2 4-2 4 2z" />
              </svg>
            </div>
            <div>
              <p class="text-xs font-bold text-indigo-500 uppercase tracking-widest mb-0.5">Barang Yang Dikembalikan</p>
              <p class="text-2xl font-black text-indigo-900 font-mono tracking-tight">{{ scannedItemCode }}</p>
            </div>
          </div>
          <button @click="resetScan" class="p-2 hover:bg-white rounded-xl text-indigo-600 transition-colors" title="Ganti Barang">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form @submit.prevent="submitReturn" class="space-y-6 px-4 pb-4">
          <div>
            <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-3 ml-1">Kondisi Saat Dikembalikan *</label>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
              <label class="relative group cursor-pointer">
                <input type="radio" v-model="form.condition" value="GOOD" class="peer sr-only" required>
                <div class="flex flex-col items-center justify-center p-4 rounded-2xl border-2 border-gray-100 bg-white transition-all peer-checked:border-green-500 peer-checked:bg-green-50/30 group-hover:border-gray-200">
                   <span class="text-2xl mb-1">✨</span>
                   <span class="text-xs font-black uppercase text-gray-400 peer-checked:text-green-700 tracking-tighter">Sangat Baik</span>
                </div>
              </label>

              <label class="relative group cursor-pointer">
                <input type="radio" v-model="form.condition" value="DAMAGED" class="peer sr-only">
                <div class="flex flex-col items-center justify-center p-4 rounded-2xl border-2 border-gray-100 bg-white transition-all peer-checked:border-orange-500 peer-checked:bg-orange-50/30 group-hover:border-gray-200">
                   <span class="text-2xl mb-1">🤕</span>
                   <span class="text-xs font-black uppercase text-gray-400 peer-checked:text-orange-700 tracking-tighter">Bermasalah</span>
                </div>
              </label>

              <label class="relative group cursor-pointer">
                <input type="radio" v-model="form.condition" value="IN_REPAIR" class="peer sr-only">
                <div class="flex flex-col items-center justify-center p-4 rounded-2xl border-2 border-gray-100 bg-white transition-all peer-checked:border-yellow-500 peer-checked:bg-yellow-50/30 group-hover:border-gray-200">
                   <span class="text-2xl mb-1">🛠️</span>
                   <span class="text-xs font-black uppercase text-gray-400 peer-checked:text-yellow-700 tracking-tighter">Perbaikan</span>
                </div>
              </label>

              <label class="relative group cursor-pointer">
                <input type="radio" v-model="form.condition" value="LOST" class="peer sr-only">
                <div class="flex flex-col items-center justify-center p-4 rounded-2xl border-2 border-gray-100 bg-white transition-all peer-checked:border-red-500 peer-checked:bg-red-50/30 group-hover:border-gray-200">
                   <span class="text-2xl mb-1">👻</span>
                   <span class="text-xs font-black uppercase text-gray-400 peer-checked:text-red-700 tracking-tighter">Hilang</span>
                </div>
              </label>
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 ml-1">Catatan Kondisi</label>
            <textarea v-model="form.notes" rows="3" class="input-field rounded-xl border-gray-200" placeholder="Jelaskan kondisi fisik barang (Opsional)"></textarea>
          </div>

          <!-- Photo Upload -->
          <div class="space-y-3">
            <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider ml-1">Foto Bukti Kondisi</label>
            <div class="flex items-center justify-center w-full">
              <label v-if="!form.photo_url" class="flex flex-col items-center justify-center w-full h-48 border-2 border-gray-100 border-dashed rounded-2xl cursor-pointer bg-gray-50/50 hover:bg-gray-100/50 transition-all group">
                <div class="flex flex-col items-center justify-center pt-5 pb-6">
                  <div class="w-12 h-12 bg-white rounded-xl shadow-sm flex items-center justify-center text-gray-400 group-hover:text-primary-600 transition-colors mb-4">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                  </div>
                  <p class="mb-1 text-xs font-bold text-gray-500 uppercase tracking-widest">Ambil Foto / Upload</p>
                  <p class="text-[10px] text-gray-400">PNG, JPG atau JPEG (Max. 5MB)</p>
                </div>
                <input type="file" @change="uploadPhoto" class="hidden" accept="image/*" capture="environment" />
              </label>
              
              <div v-else class="relative w-full group overflow-hidden rounded-2xl border-4 border-white shadow-xl">
                <img :src="getFullImageUrl(form.photo_url)" class="w-full h-64 object-cover" />
                <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                  <button @click="form.photo_url = ''" type="button" class="btn-primary bg-red-600 hover:bg-red-700 border-none px-4 py-2">
                    Hapus & Ganti Foto
                  </button>
                </div>
              </div>
            </div>
            <p v-if="uploading" class="text-[10px] font-black text-primary-600 animate-pulse text-center uppercase tracking-widest">Mengunggah file...</p>
          </div>
          
          <div v-if="transactionStore.error" class="p-4 bg-red-50 border border-red-100 text-red-700 text-xs rounded-xl flex items-center gap-3">
            <svg class="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
            </svg>
            <span class="font-medium">{{ transactionStore.error }}</span>
          </div>

          <div v-if="successMsg" class="p-6 bg-green-50 justify-between items-center text-green-800 rounded-2xl flex gap-6 border border-green-200 animate-scale-up">
            <div>
              <p class="font-black text-lg flex items-center gap-2">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Berhasil Dikembalikan
              </p>
              <p class="text-sm mt-1 font-medium opacity-90">Data status barang telah diperbarui ke sistem.</p>
            </div>
            <router-link to="/my-borrowings" class="btn-primary px-6 py-3 text-sm whitespace-nowrap shadow-lg shadow-primary-200">Lihat Daftar</router-link>
          </div>

          <div class="pt-6 border-t border-gray-100 flex gap-4">
            <button @click="resetScan" type="button" class="btn-secondary px-8 rounded-xl font-bold">Batal</button>
            <button type="submit" :disabled="transactionStore.loading || uploading || !!successMsg" class="flex-1 btn-primary justify-center py-4 rounded-xl font-bold text-lg shadow-xl shadow-primary-200 transition-all">
              {{ transactionStore.loading ? 'Memproses...' : (uploading ? 'Menunggu Upload...' : 'Konfirmasi Pengembalian') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useTransactionStore } from '../stores/transaction'
import QrScanner from '../components/QrScanner.vue'
import api from '../utils/api'

const transactionStore = useTransactionStore()
const route = useRoute()

const scannedItemCode = ref('')
const successMsg = ref(false)
const uploading = ref(false)

const form = ref({
  condition: 'GOOD',
  notes: '',
  photo_url: ''
})

onMounted(() => {
  if (route.query.code) {
    scannedItemCode.value = route.query.code
  }
})

const handleScan = (code) => {
  scannedItemCode.value = code
  successMsg.value = false
  transactionStore.error = null
}

const getFullImageUrl = (url) => {
  if (!url) return ''
  return url
}

const uploadPhoto = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  if (file.size > 5 * 1024 * 1024) {
    alert('File terlalu besar! Maksimal 5MB.')
    return
  }

  uploading.value = true
  const formData = new FormData()
  formData.append('file', file)

  try {
    const response = await api.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    if (response.data.success) {
      form.value.photo_url = response.data.data.url
    } else {
      alert('Gagal upload: ' + response.data.message)
    }
  } catch (err) {
    alert('Gagal mengunggah gambar: ' + (err.response?.data?.error || err.message))
  } finally {
    uploading.value = false
  }
}

const resetScan = () => {
  scannedItemCode.value = ''
  successMsg.value = false
  form.value = {
    condition: 'GOOD',
    notes: '',
    photo_url: ''
  }
}

const submitReturn = async () => {
  try {
    await transactionStore.returnItem(
      scannedItemCode.value, 
      form.value.condition, 
      form.value.notes,
      form.value.photo_url
    )
    successMsg.value = true
  } catch (err) {
    // Error is handled in store
  }
}
</script>
