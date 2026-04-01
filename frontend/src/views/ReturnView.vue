<template>
  <div class="animate-fade-in space-y-10 max-w-3xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
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
          <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 4v1m0 11v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"/></svg>
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
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 15v-1a4 4 0 00-4-4H8m0 0l3 3m-3-3l3-3m9 14V5a2 2 0 00-2-2H6a2 2 0 00-2 2v16l4-2 4 2 4-2 4 2z" /></svg>
            </div>
            <div>
              <p class="text-[10px] font-black text-primary-200 uppercase tracking-widest leading-none mb-1">Aset Terdeteksi</p>
              <p class="text-3xl font-black font-mono tracking-tighter leading-none">{{ scannedItemCode }}</p>
            </div>
          </div>
          <button @click="resetScan" class="bg-white/20 hover:bg-white/30 text-white px-6 py-3 rounded-xl text-[10px] font-black tracking-widest transition-all active:scale-95">
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
                { val: 'GOOD', label: 'BAIK', icon: '✨' },
                { val: 'DAMAGED', label: 'RUSAK', icon: '🤕' },
                { val: 'IN_REPAIR', label: 'SERVIS', icon: '🛠️' },
                { val: 'LOST', label: 'HILANG', icon: '👻' }
              ]" :key="c.val" class="relative group cursor-pointer">
                <input type="radio" v-model="form.condition" :value="c.val" class="peer sr-only">
                <div class="flex flex-col items-center justify-center p-6 rounded-[2rem] border-2 border-gray-100 dark:border-gray-700 bg-white dark:bg-gray-800 transition-all peer-checked:scale-105 peer-checked:border-primary-500 dark:peer-checked:border-primary-400 peer-checked:bg-primary-50 dark:peer-checked:bg-primary-600/30 peer-checked:shadow-xl peer-checked:shadow-primary-500/20 group-hover:border-gray-200 dark:group-hover:border-gray-600">
                    <div class="absolute top-4 right-4 text-primary-600 opacity-0 transform scale-0 peer-checked:opacity-100 peer-checked:scale-100 transition-all duration-300">
                      <svg class="w-6 h-6 shadow-sm" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" /></svg>
                    </div>
                    <span class="text-3xl mb-2">{{ c.icon }}</span>
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
                  <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
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
                    class="bg-primary-600 text-white py-5 flex-[2] rounded-2xl font-black text-xs tracking-[0.3em] shadow-xl hover:bg-primary-700 active:scale-95 transition-all">
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
