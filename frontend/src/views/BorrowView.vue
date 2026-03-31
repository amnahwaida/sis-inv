<template>
  <div class="animate-fade-in max-w-2xl mx-auto space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Pinjam Barang</h1>
        <p class="text-sm text-gray-500 mt-1">Gunakan scanner untuk meminjam aset sekolah.</p>
      </div>
    </div>

    <!-- Scanner or Selected Item Card -->
    <div class="card overflow-hidden transition-all duration-300 min-h-[400px] flex flex-col justify-center">
      <div v-if="!scannedItemCode" class="p-4">
        <p class="text-gray-500 mb-6 text-center text-sm">Scan QR Code pada label barang untuk mematikan status ketersediaan.</p>
        <QrScanner :active="true" @scanned="handleScan" />
      </div>

      <div v-else-if="loadingItem" class="flex flex-col items-center justify-center p-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        <p class="mt-4 text-sm font-bold text-gray-500 uppercase tracking-widest">Memuat Detail Barang...</p>
      </div>

      <div v-else-if="itemDetail" class="space-y-6 animate-scale-up p-1">
        <!-- Selected Item Badge -->
        <div class="bg-primary-50 border border-primary-100 rounded-2xl p-6 flex items-center justify-between shadow-sm">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-white rounded-xl shadow-sm flex items-center justify-center text-primary-600">
              <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
            <div>
              <p class="text-xs font-bold text-primary-500 uppercase tracking-widest mb-0.5">{{ itemDetail.name }}</p>
              <p class="text-2xl font-black text-primary-900 font-mono tracking-tight">{{ itemDetail.code }}</p>
            </div>
          </div>
          <button @click="resetScan" class="p-2 hover:bg-white rounded-xl text-primary-600 transition-colors" title="Ganti Barang">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form @submit.prevent="submitBorrow" class="space-y-6 px-4 pb-4">
          <!-- Borrower Type Selector -->
          <div class="space-y-3">
            <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider ml-1">Jenis Peminjam *</label>
            <div class="bg-gray-50 rounded-2xl p-1.5 flex gap-1.5 border border-gray-100">
              <button 
                type="button"
                @click="form.borrower_type = 'STAFF'"
                class="flex-1 py-3 rounded-xl text-sm font-bold transition-all"
                :class="form.borrower_type === 'STAFF' ? 'bg-white shadow-sm text-primary-600 ring-1 ring-gray-200' : 'text-gray-500 hover:text-gray-700'"
              >
                Diri Sendiri (Staff)
              </button>
              <button 
                type="button"
                @click="form.borrower_type = 'STUDENT'"
                :disabled="itemDetail.borrower_type === 'STAFF_ONLY'"
                class="flex-1 py-3 rounded-xl text-sm font-bold transition-all flex items-center justify-center gap-2"
                :class="[
                  form.borrower_type === 'STUDENT' ? 'bg-white shadow-sm text-primary-600 ring-1 ring-gray-200' : 'text-gray-500 hover:text-gray-700',
                  itemDetail.borrower_type === 'STAFF_ONLY' ? 'opacity-50 grayscale cursor-not-allowed' : ''
                ]"
              >
                <svg v-if="itemDetail.borrower_type === 'STAFF_ONLY'" class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
                </svg>
                Wakili Siswa
              </button>
            </div>
            <!-- Staff Only Warning -->
            <div v-if="itemDetail.borrower_type === 'STAFF_ONLY'" class="flex items-start gap-2 px-2 py-1">
               <svg class="w-4 h-4 text-amber-500 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                 <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
               </svg>
               <p class="text-[11px] font-bold text-amber-600 uppercase tracking-tight">Barang ini dibatasi hanya untuk penggunaan Staff/Guru.</p>
            </div>
          </div>

          <!-- Student Fields Conditional -->
          <div v-if="form.borrower_type === 'STUDENT'" class="grid grid-cols-1 md:grid-cols-2 gap-4 animate-fade-in bg-primary-50/30 p-4 rounded-2xl border border-primary-50">
            <div class="md:col-span-2 relative">
              <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 ml-1">Nama Siswa / Cari NIS *</label>
              <input 
                v-model="form.student_name" 
                required 
                class="input-field rounded-xl border-gray-200" 
                placeholder="Ketik nama atau NIS untuk mencari..." 
                @input="debouncedSearch"
              />
              
              <!-- Auto-suggest Dropdown -->
              <div v-if="searchResults.length > 0" class="absolute z-50 w-full mt-1 bg-white rounded-xl shadow-xl border border-gray-100 overflow-hidden divide-y divide-gray-50">
                <div 
                  v-for="s in searchResults" 
                  :key="s.nis"
                  @click="selectStudent(s)"
                  class="px-4 py-3 hover:bg-primary-50 cursor-pointer transition-colors"
                >
                  <div class="text-sm font-bold text-primary-900">{{ s.full_name }}</div>
                  <div class="text-[10px] text-gray-400 font-mono">NIS: {{ s.nis }} • {{ s.class }}</div>
                </div>
              </div>
            </div>
            <div>
              <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 ml-1">NIS Siswa</label>
              <input v-model="form.student_nis" class="input-field rounded-xl border-gray-200" placeholder="Nomor Induk Siswa" />
            </div>
            <div>
              <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 ml-1">Kelas</label>
              <input v-model="form.student_class" class="input-field rounded-xl border-gray-200" placeholder="CTH: X-IPA-1" />
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 ml-1">Estimasi Dikembalikan *</label>
              <div v-if="form.borrower_type === 'STUDENT'" class="p-3 bg-indigo-50 border border-indigo-100 rounded-xl text-indigo-700 text-sm font-bold flex items-center gap-2">
                <svg class="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
                </svg>
                Otomatis 6 Jam (Aturan Siswa)
              </div>
              <select v-else v-model="form.expected_return_days" class="input-field rounded-xl border-gray-200" required>
                <option :value="1">1 Hari (Besok)</option>
                <option :value="3">3 Hari</option>
                <option :value="7">1 Minggu</option>
                <option :value="14">2 Minggu</option>
                <option :value="30">1 Bulan</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 ml-1">Tujuan / Keperluan *</label>
            <textarea v-model="form.purpose" required rows="3" class="input-field rounded-xl border-gray-200" placeholder="Contoh: Digunakan untuk praktikum di Laboratorium Kimia"></textarea>
          </div>
          
          <div v-if="transactionStore.error" class="p-4 bg-red-50 border border-red-100 text-red-700 text-xs rounded-xl flex items-center gap-3 animate-shake">
            <div class="w-8 h-8 rounded-full bg-red-100 flex items-center justify-center flex-shrink-0">
               <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                 <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
               </svg>
            </div>
            <p class="font-medium font-bold uppercase tracking-tight">{{ transactionStore.error }}</p>
          </div>

          <div v-if="successMsg" class="p-6 bg-green-50 justify-between items-center text-green-800 rounded-2xl flex gap-6 border border-green-200 animate-scale-up">
            <div>
              <p class="font-black text-lg flex items-center gap-2">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                </svg>
                Peminjaman Berhasil!
              </p>
              <p class="text-sm mt-1 font-medium opacity-90">{{ successMsg }}</p>
            </div>
            <router-link to="/my-borrowings" class="btn-primary px-6 py-3 text-sm whitespace-nowrap shadow-lg shadow-primary-200">Daftar Pinjaman</router-link>
          </div>

          <div class="pt-6 border-t border-gray-100 flex gap-4">
            <button @click="resetScan" type="button" class="btn-secondary px-8 rounded-xl font-bold">Batal Scan</button>
            <button type="submit" :disabled="transactionStore.loading || !!successMsg" class="flex-1 btn-primary justify-center py-4 rounded-xl font-bold text-lg shadow-xl shadow-primary-200 disabled:shadow-none translate-y-0 active:translate-y-1 transition-all">
              {{ transactionStore.loading ? 'Memproses...' : 'Proses Peminjaman' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import api from '../utils/api'
import { useTransactionStore } from '../stores/transaction'
import { useItemStore } from '../stores/item'
import QrScanner from '../components/QrScanner.vue'

const transactionStore = useTransactionStore()
const itemStore = useItemStore()

const scannedItemCode = ref('')
const loadingItem = ref(false)
const itemDetail = ref(null)
const successMsg = ref('')

const form = ref({
  borrower_type: 'STAFF',
  student_name: '',
  student_nis: '',
  student_class: '',
  expected_return_days: 7,
  purpose: ''
})

const searchResults = ref([])
let searchTimeout

const debouncedSearch = () => {
  clearTimeout(searchTimeout)
  if (!form.value.student_name) {
    searchResults.value = []
    return
  }
  searchTimeout = setTimeout(async () => {
    try {
      const { data } = await api.get(`/students/search?q=${form.value.student_name}`)
      if (data.success) {
        searchResults.value = data.data
      }
    } catch {
      searchResults.value = []
    }
  }, 300)
}

const selectStudent = (student) => {
  form.value.student_name = student.full_name
  form.value.student_nis = student.nis
  form.value.student_class = student.class
  searchResults.value = []
}

const handleScan = async (code) => {
  scannedItemCode.value = code
  successMsg.value = ''
  transactionStore.error = null
  loadingItem.value = true
  
  try {
    const detail = await itemStore.fetchItemByCode(code)
    itemDetail.value = detail
    // If restricted, default to STAFF
    if (detail.borrower_type === 'STAFF_ONLY') {
      form.value.borrower_type = 'STAFF'
    }
  } catch (err) {
    scannedItemCode.value = '' // Reset on error
  } finally {
    loadingItem.value = false
  }
}

const resetScan = () => {
  scannedItemCode.value = ''
  itemDetail.value = null
  successMsg.value = ''
  form.value.purpose = ''
  form.value.student_name = ''
  form.value.student_nis = ''
  form.value.student_class = ''
  form.value.borrower_type = 'STAFF'
}

const submitBorrow = async () => {
  try {
    const payload = {
      item_code: scannedItemCode.value,
      borrower_type: form.value.borrower_type,
      purpose: form.value.purpose,
      expected_return_days: form.value.expected_return_days
    }

    if (form.value.borrower_type === 'STUDENT') {
      payload.student_name = form.value.student_name
      payload.student_nis = form.value.student_nis
      payload.student_class = form.value.student_class
    }

    const res = await transactionStore.borrowItem(payload)
    
    const borrowerLabel = form.value.borrower_type === 'STUDENT' ? `Siswa (${form.value.student_name})` : 'Anda'
    successMsg.value = `Barang ${res.item_name} berhasil dipinjam untuk ${borrowerLabel}.`
  } catch (err) {
    // Error is already reactive in transactionStore
  }
}
</script>

<style scoped>
.animate-shake {
  animation: shake 0.6s cubic-bezier(.36,.07,.19,.97) both;
}

@keyframes shake {
  10%, 90% { transform: translate3d(-1px, 0, 0); }
  20%, 80% { transform: translate3d(2px, 0, 0); }
  30%, 50%, 70% { transform: translate3d(-4px, 0, 0); }
  40%, 60% { transform: translate3d(4px, 0, 0); }
}
</style>
