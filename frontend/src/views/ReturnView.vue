<template>
  <div class="animate-fade-in space-y-10 max-w-3xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-primary-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Kembalikan Barang</h1>
          <p class="text-primary-100/70 text-sm font-medium">Lakukan pengecekan kondisi sebelum disimpan kembali</p>
        </div>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden transition-all duration-300">
      <!-- Scanner State -->
      <div v-if="!scannedItemCode" class="p-12 text-center space-y-8">
        <div class="w-24 h-24 bg-primary-50 dark:bg-primary-900/30 rounded-full flex items-center justify-center mx-auto text-primary-600 animate-pulse transition-all">
          <svg class="w-12 h-12" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M7.5 3.75H6A2.25 2.25 0 003.75 6v1.5M16.5 3.75H18A2.25 2.25 0 0120.25 6v1.5m0 9V18a2.25 2.25 0 01-2.25 2.25h-1.5M7.5 20.25H6A2.25 2.25 0 013.75 18v-1.5M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </div>
        <div class="space-y-2">
          <h2 class="text-xl font-black text-gray-900 dark:text-white uppercase tracking-tight">Siapkan QR Code</h2>
          <p class="text-sm text-gray-400 font-medium">Scan label pada barang untuk memproses pengembalian</p>
        </div>
        <div class="max-w-md mx-auto p-4 bg-gray-50 dark:bg-gray-900/50 rounded-3xl border-2 border-dashed border-gray-200 dark:border-gray-700">
          <QrScanner :active="true" @scanned="handleScan" />
        </div>
      </div>

      <!-- Detail Form State -->
      <div v-else class="animate-scale-up">
        <div class="bg-primary-600 p-8 flex flex-col sm:flex-row items-center justify-between gap-6">
          <div class="flex items-center gap-5 text-white">
            <div class="w-16 h-16 bg-white/20 backdrop-blur-md rounded-2xl flex items-center justify-center shadow-inner">
              <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M10.125 2.25h-4.5c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125v-9M10.125 2.25h.375a9 9 0 019 9v.375M10.125 2.25A3.375 3.375 0 0113.5 5.625v1.5c0 .621.504 1.125 1.125 1.125h1.5a3.375 3.375 0 013.375 3.375M9 15l2.25 2.25L15 12" />
              </svg>
            </div>
            <div>
              <p class="text-[10px] font-black text-primary-200 uppercase tracking-widest leading-none mb-1">Aset Terdeteksi</p>
              <p class="text-3xl font-black font-mono tracking-tighter leading-none">{{ scannedItemCode }}</p>
            </div>
          </div>
          <button @click="resetScan" class="bg-white/10 hover:bg-white/20 text-white border border-white/10 px-6 py-3 rounded-xl text-[10px] font-black tracking-widest transition-all active:scale-95">
            GANTI BARANG
          </button>
        </div>

        <!-- Borrow Comparison Info (NEW) -->
        <div v-if="itemInfo?.last_borrow_photo_url" class="p-8 bg-gray-50 dark:bg-gray-900/50 border-b border-gray-100 dark:border-gray-700 animate-fade-in">
           <div class="flex flex-col md:flex-row gap-8 items-center">
              <div class="w-full md:w-48 h-48 rounded-[2rem] overflow-hidden shadow-lg border-2 border-white dark:border-gray-800">
                 <img :src="itemInfo.last_borrow_photo_url" class="w-full h-full object-cover">
              </div>
              <div class="flex-1 space-y-4">
                 <div>
                    <p class="text-[10px] font-black text-primary-600 uppercase tracking-widest leading-none mb-1">Dipinjam Oleh</p>
                    <p class="text-xl font-black text-gray-900 dark:text-white uppercase leading-none">{{ itemInfo.current_borrower }}</p>
                 </div>
                 <div class="p-4 bg-primary-50 dark:bg-primary-900/10 rounded-2xl border border-primary-100 dark:border-primary-800">
                    <p class="text-[10px] font-bold text-primary-600 italic">"Gunakan foto di samping untuk membandingkan kondisi barang saat dipinjam vs saat dikembalikan saat ini."</p>
                 </div>
              </div>
           </div>
        </div>

        <form @submit.prevent="submitReturn" class="p-10 space-y-10">
          <!-- Condition Radio Grid -->
          <div class="space-y-4">
            <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Kondisi Saat Ini</label>
            <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
              <label v-for="c in [
                { val: 'GOOD', label: 'BAIK', icon: 'M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 00-2.456 2.456zM16.894 20.567L16.5 21.75l-.394-1.183a2.25 2.25 0 00-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 001.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 001.423 1.423l1.183.394-1.183.394a2.25 2.25 0 00-1.423 1.423z' },
                { val: 'DAMAGED', label: 'RUSAK', icon: 'M15.182 15.182a4.5 4.5 0 01-6.364 0M21 12a9 9 0 11-18 0 9 9 0 0118 0zM9.75 9.75c0 .414-.168.75-.375.75S9 10.164 9 9.75 9.168 9 9.375 9s.375.336.375.75zm-.375 0h.008v.015h-.008V9.75zm5.625 0c0 .414-.168.75-.375.75s-.375-.336-.375-.75.168-.75.375-.75.375.336.375.75zm-.375 0h.008v.015h-.008V9.75z' },
                { val: 'IN_REPAIR', label: 'SERVIS', icon: 'M11.423 20.11L12.122 21a.5.5 0 00.756 0l.699-.89a8.49 8.49 0 015.485-3.09l1.177-.214a.5.5 0 00.37-.621l-.663-2.473a.5.5 0 00-.518-.363l-1.19.043a8.49 8.49 0 01-5.44-2.06l-.919-.777a.5.5 0 00-.634 0l-.92.777a8.49 8.49 0 01-5.44 2.06l-1.19-.043a.5.5 0 00-.518.362l-.663 2.474a.5.5 0 00.37.62l1.177.214a8.49 8.49 0 015.485 3.09z' },
                { val: 'LOST', label: 'HILANG', icon: 'M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9 5.25h.008v.008H12v-.008z' }
              ]" :key="c.val" class="relative group cursor-pointer">
                <input type="radio" v-model="form.condition" :value="c.val" class="peer sr-only">
                <div class="flex flex-col items-center justify-center p-6 rounded-[2rem] border-2 border-gray-100 dark:border-gray-700 bg-white dark:bg-gray-800 transition-all peer-checked:scale-105 peer-checked:border-primary-500 dark:peer-checked:border-primary-400 peer-checked:bg-primary-50 dark:peer-checked:bg-primary-600/30 peer-checked:shadow-xl peer-checked:shadow-primary-500/20 group-hover:border-gray-200 dark:group-hover:border-gray-600">
                    <div class="absolute top-4 right-4 text-primary-600 opacity-0 transform scale-0 peer-checked:opacity-100 peer-checked:scale-100 transition-all duration-300">
                      <svg class="w-6 h-6 shadow-sm" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                    </div>
                    <span class="w-10 h-10 mb-2 text-gray-500 peer-checked:text-primary-600 transition-colors">
                      <svg class="w-full h-full" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" :d="c.icon" />
                      </svg>
                    </span>
                    <span class="text-[10px] font-black uppercase text-gray-400 peer-checked:text-primary-600 dark:peer-checked:text-primary-300 tracking-[0.2em] leading-none">{{ c.label }}</span>
                </div>
              </label>
            </div>
          </div>

          <div class="space-y-4">
            <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Keterangan / Catatan</label>
            <textarea v-model="form.notes" rows="4" class="input-field rounded-[2.5rem] p-6" placeholder="Beri catatan tambahan jika terjadi kerusakan atau hal mendesak..."></textarea>
          </div>

          <!-- Photo Upload -->
          <div class="space-y-4">
            <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Unggah Foto Bukti</label>
            <div class="flex flex-col items-center justify-center w-full">
              <label v-if="!form.photo_url" class="group flex flex-col items-center justify-center w-full h-64 border-2 border-dashed border-gray-200 dark:border-gray-700 rounded-[2.5rem] cursor-pointer bg-gray-50/50 dark:bg-gray-900/20 hover:bg-white dark:hover:bg-gray-800 transition-all">
                <div class="w-16 h-16 bg-white dark:bg-gray-800 rounded-3xl shadow-sm flex items-center justify-center text-gray-400 group-hover:text-primary-600 group-hover:scale-110 transition-all mb-4">
                  <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6.827 6.175A2.31 2.31 0 015.186 7.23c-.38.054-.757.112-1.134.175C2.999 7.58 2.25 8.507 2.25 9.574V18a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9.574c0-1.067-.75-1.994-1.802-2.169a47.865 47.865 0 00-1.134-.175 2.31 2.31 0 01-1.64-1.055l-.822-1.316a2.192 2.192 0 00-1.736-1.039 48.774 48.774 0 00-5.232 0 2.192 2.192 0 00-1.736 1.039l-.821 1.316z" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 12.75a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0zM18.75 10.5h.008v.008h-.008V10.5z" />
                  </svg>
                </div>
                <p class="text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">AMBIL GAMBAR / UPLOAD</p>
                <input type="file" @change="uploadPhoto" class="hidden" accept="image/*" capture="environment" />
              </label>
              <div v-else class="relative w-full rounded-[2.5rem] overflow-hidden border-4 border-white dark:border-gray-700 shadow-xl h-64">
                <img :src="getFullImageUrl(form.photo_url)" class="w-full h-full object-cover" />
                <button @click="form.photo_url = ''" type="button" class="absolute inset-0 bg-black/40 text-white font-black text-[10px] tracking-widest opacity-0 hover:opacity-100 transition-all uppercase">Hapus Foto</button>
              </div>
            </div>
          </div>
          
          <div v-if="successMsg" class="p-8 bg-emerald-50 dark:bg-emerald-900/30 rounded-[2.5rem] border border-emerald-100 dark:border-emerald-800 text-center animate-scale-up">
            <p class="text-lg font-black text-emerald-900 dark:text-emerald-100 uppercase tracking-tight">PENGEMBALIAN SUKSES!</p>
            <p class="text-xs font-bold text-emerald-600 dark:text-emerald-400 mt-1 uppercase">Aset telah direkam kembali ke sistem.</p>
            <router-link to="/items" class="inline-block mt-6 bg-emerald-600 text-white px-8 py-3 rounded-2xl font-black text-[10px] tracking-widest shadow-xl">LIHAT DAFTAR BARANG</router-link>
          </div>

          <div v-if="!successMsg" class="pt-8 border-t border-gray-100 dark:border-gray-700 flex gap-4">
            <button @click="resetScan" type="button" class="btn-secondary py-5 flex-1 rounded-2xl font-black text-xs">BATAL</button>
            <button type="submit" :disabled="transactionStore.loading || uploading" 
                    class="btn-premium-action flex-[2]">
              KONFIRMASI SEKARANG
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

const form = ref({ condition: 'GOOD', notes: '', photo_url: '' })

onMounted(() => { 
  if (route.query.code) {
    handleScan(route.query.code)
  }
})

const itemInfo = ref(null)
const handleScan = async (code) => { 
  scannedItemCode.value = code; 
  successMsg.value = false 
  try {
    const { data } = await api.get(`/items/code/${code}`)
    itemInfo.value = data.data
  } catch (err) { console.error(err) }
}
const resetScan = () => { 
  scannedItemCode.value = ''; 
  successMsg.value = false; 
  itemInfo.value = null;
  form.value = { condition: 'GOOD', notes: '', photo_url: '' } 
}
const getFullImageUrl = (url) => url || ''

const uploadPhoto = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  uploading.value = true
  const formData = new FormData()
  formData.append('file', file)
  try {
    const res = await api.post('/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    if (res.data.success) {
      form.value.photo_url = res.data.data.url
    }
  } catch (err) { alert('Gagal mengunggah foto kondisi.') }
  finally { uploading.value = false }
}

const submitReturn = async () => {
  try {
    await transactionStore.returnItem(scannedItemCode.value, form.value.condition, form.value.notes, form.value.photo_url)
    successMsg.value = true
  } catch (err) { 
    alert(err.response?.data?.error || 'Gagal mengembalikan barang. Pastikan data lengkap.')
    console.error('Gagal kembalikan barang', err) 
  }
}
</script>
