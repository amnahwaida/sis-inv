<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Pinjam Barang</h1>
          <p class="text-primary-100/70 text-sm font-medium">Proses pencatatan peminjaman aset baru</p>
        </div>
        
        <div class="flex items-center gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10">
          <router-link to="/my-borrowings" 
                    class="bg-white/20 hover:bg-white/30 text-white px-5 py-2.5 rounded-xl text-xs font-black transition-all flex items-center gap-2 active:scale-95">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" /></svg>
            RIWAYAT SAYA
          </router-link>
        </div>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
      
      <!-- Left Column: Selection & Details -->
      <div class="lg:col-span-7 space-y-8">
        <!-- Step 1: Barang -->
        <section class="bg-white dark:bg-gray-800 rounded-[2.5rem] p-8 border border-gray-100 dark:border-gray-700 shadow-sm">
          <div class="flex items-center gap-4 mb-8">
            <div class="w-10 h-10 bg-indigo-600 text-white rounded-xl flex items-center justify-center font-black">1</div>
            <h2 class="text-xl font-black text-gray-900 dark:text-white uppercase tracking-tight">Pilih Barang</h2>
          </div>
          
          <div v-if="!selectedItem" class="space-y-6">
            <div class="p-6 bg-gray-50 dark:bg-gray-900/50 rounded-3xl border border-dashed border-gray-200 dark:border-gray-700">
              <QrScanner @scanned="handleScanned" />
            </div>
            <div class="flex items-center gap-4">
              <div class="h-px bg-gray-100 dark:bg-gray-700 flex-1"></div>
              <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">Atau Cari Barang</span>
              <div class="h-px bg-gray-100 dark:bg-gray-700 flex-1"></div>
            </div>
            <div class="relative">
              <input type="text" v-model="searchQuery" @input="debouncedSearch" 
                     class="input-field pl-12 h-14 rounded-2xl text-base shadow-sm" placeholder="Ketik nama atau kode barang..." />
              <svg class="w-6 h-6 absolute left-4 top-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
            </div>
            
            <!-- Search Results Dropdown -->
            <div v-if="searchResults.length > 0" class="bg-white dark:bg-gray-800 rounded-2xl border border-gray-200 dark:border-gray-700 shadow-xl overflow-hidden animate-slide-in">
              <div v-for="item in searchResults" :key="item.id" @click="selectItem(item)"
                   class="p-4 hover:bg-primary-50 dark:hover:bg-primary-900/20 cursor-pointer border-b border-gray-50 dark:border-gray-700 last:border-0 flex items-center justify-between transition-colors">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 bg-gray-100 dark:bg-gray-700 rounded-lg flex items-center justify-center text-gray-400">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>
                  </div>
                  <div>
                    <div class="font-black text-gray-900 dark:text-white text-sm">{{ item.name }}</div>
                    <div class="text-[10px] font-bold text-gray-400 font-mono">{{ item.code }}</div>
                  </div>
                </div>
                <span class="text-[9px] font-black uppercase" 
                      :class="item.status === 'AVAILABLE' ? 'text-emerald-600 bg-emerald-50' : 'text-amber-600 bg-amber-50'">
                  {{ item.status }}
                </span>
              </div>
            </div>
          </div>

          <!-- Selected Item Card -->
          <div v-else class="animate-scale-up p-6 bg-primary-50 dark:bg-primary-900/20 rounded-[2rem] border border-primary-100 dark:border-primary-800 relative">
            <button @click="selectedItem = null" class="absolute top-4 right-4 p-2 text-primary-400 hover:text-primary-600 transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12"/></svg>
            </button>
            <div class="flex items-center gap-4">
              <div class="w-16 h-16 bg-white dark:bg-gray-800 text-primary-600 rounded-2xl flex items-center justify-center shadow-lg">
                <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>
              </div>
              <div>
                <p class="text-[10px] font-black text-primary-400 uppercase tracking-widest leading-none mb-1">Barang Terpilih</p>
                <h3 class="text-xl font-black text-gray-900 dark:text-white leading-tight">{{ selectedItem.name }}</h3>
                <p class="text-xs font-bold text-primary-600 font-mono mt-1">{{ selectedItem.code }}</p>
              </div>
            </div>
          </div>
        </section>

        <!-- Step 2: Form Peminjam -->
        <section class="bg-white dark:bg-gray-800 rounded-[2.5rem] p-8 border border-gray-100 dark:border-gray-700 shadow-sm">
          <div class="flex items-center gap-4 mb-8">
            <div class="w-10 h-10 bg-indigo-600 text-white rounded-xl flex items-center justify-center font-black">2</div>
            <h2 class="text-xl font-black text-gray-900 dark:text-white uppercase tracking-tight">Informasi Peminjam</h2>
          </div>
          
          <div class="space-y-8">
            <!-- Borrower Type Toggle -->
            <div class="space-y-3">
              <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] mb-2 pl-1">Siapa yang meminjam?</label>
              <div class="flex p-1.5 bg-gray-100 dark:bg-gray-900/50 rounded-[1.5rem] w-full max-w-md border border-gray-200 dark:border-gray-700">
                <button v-for="type in [{id:'STUDENT', label:'Untuk Siswa', icon:'🎓'}, {id:'STAFF', label:'Internal (Guru/Staff)', icon:'👤'}]" :key="type.id"
                        @click="borrowerType = type.id; selectedStudent = null; studentSearch = ''"
                        class="flex-1 py-3.5 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all flex items-center justify-center gap-2"
                        :class="borrowerType === type.id ? 'bg-indigo-600 text-white shadow-xl shadow-indigo-500/20 scale-[1.02]' : 'text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'">
                  <span>{{ type.icon }}</span>
                  {{ type.label }}
                </button>
              </div>
            </div>

            <!-- Student Search (Only if type is STUDENT) -->
            <div v-if="borrowerType === 'STUDENT'" class="space-y-4">
              <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] mb-2 pl-1">Data Siswa (NIS/Nama)</label>
              <div class="relative">
                <input type="text" v-model="studentSearch" @input="debouncedStudentSearch" 
                       class="input-field pl-12 h-14 rounded-2xl" placeholder="Ketik nama atau NIS siswa..." />
                <svg class="w-6 h-6 absolute left-4 top-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/></svg>
                
                <!-- Suggestions Dropdown -->
                <div v-if="studentSuggestions.length > 0 && !selectedStudent" 
                     class="absolute z-50 w-full mt-2 bg-white dark:bg-gray-800 rounded-2xl border border-gray-200 dark:border-gray-700 shadow-2xl max-h-60 overflow-y-auto">
                  <div v-for="s in studentSuggestions" :key="s.nis" @click="selectStudent(s)"
                       class="p-4 hover:bg-primary-50 dark:hover:bg-primary-900/20 cursor-pointer border-b border-gray-50 dark:border-gray-700 last:border-0 transition-colors">
                    <div class="font-black text-gray-900 dark:text-white text-sm">{{ s.full_name }}</div>
                    <div class="text-[10px] font-bold text-gray-400 uppercase tracking-widest">{{ s.nis }} • {{ s.class }}</div>
                  </div>
                </div>
              </div>

              <!-- Manual Entry if not selected / New Student -->
              <div v-if="selectedStudent" class="animate-scale-up p-5 bg-emerald-50 dark:bg-emerald-900/20 rounded-2xl border border-emerald-100 dark:border-emerald-800 flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 bg-white dark:bg-gray-800 text-emerald-600 rounded-xl flex items-center justify-center shadow-md">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/></svg>
                  </div>
                  <div>
                    <div class="font-black text-gray-900 dark:text-white text-sm">{{ selectedStudent.full_name }}</div>
                    <div class="text-[10px] font-bold text-emerald-600 uppercase tracking-widest">{{ selectedStudent.nis }} • {{ selectedStudent.class }}</div>
                  </div>
                </div>
                <button @click="selectedStudent = null; studentSearch = ''" class="text-emerald-400 hover:text-emerald-600">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12"/></svg>
                </button>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] mb-2 pl-1">Durasi Pinjam (Hari)</label>
                <input type="number" v-model="durationDays" min="1" class="input-field h-14 rounded-2xl" />
              </div>
              <div>
                <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] mb-2 pl-1">Tujuan Peminjaman</label>
                <input type="text" v-model="purpose" class="input-field h-14 rounded-2xl" placeholder="Contoh: Mengajar Kelas 12" />
              </div>
            </div>

            <!-- PHOTO UPLOAD (NEW FEATURE) -->
            <div class="space-y-4">
              <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] mb-2 pl-1">Foto Kondisi Awal</label>
              <div class="flex flex-col items-center justify-center w-full">
                <label v-if="!photoUrl" class="group flex flex-col items-center justify-center w-full h-48 border-2 border-dashed border-gray-200 dark:border-gray-700 rounded-[2rem] cursor-pointer bg-gray-50/50 dark:bg-gray-900/20 hover:bg-white dark:hover:bg-gray-800 transition-all">
                  <div class="w-12 h-12 bg-white dark:bg-gray-800 rounded-2xl shadow-sm flex items-center justify-center text-gray-400 group-hover:text-primary-600 group-hover:scale-110 transition-all mb-3">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
                  </div>
                  <p class="text-[9px] font-black text-gray-400 uppercase tracking-[0.2em]">AMBIL GAMBAR / UPLOAD</p>
                  <input type="file" @change="uploadPhoto" class="hidden" accept="image/*" capture="environment" />
                </label>
                <div v-else class="relative w-full rounded-[2rem] overflow-hidden border-2 border-primary-100 dark:border-gray-700 shadow-xl h-48">
                  <img :src="photoUrl" class="w-full h-full object-cover" />
                  <button @click="photoUrl = ''" type="button" class="absolute inset-0 bg-black/40 text-white font-black text-[10px] tracking-widest opacity-0 hover:opacity-100 transition-all uppercase">Hapus Foto</button>
                </div>
              </div>
              <p v-if="uploading" class="text-[9px] font-bold text-primary-500 animate-pulse uppercase tracking-widest text-center">Sedang mengunggah foto...</p>
            </div>

          </div>
        </section>
      </div>

      <!-- Right Column: Summary & Submit -->
      <div class="lg:col-span-5 sticky top-8">
        <div class="bg-gray-900 text-white rounded-[2.5rem] p-10 shadow-2xl relative overflow-hidden">
          <div class="absolute top-0 right-0 w-48 h-48 bg-primary-600/20 rounded-full blur-3xl"></div>
          
          <h3 class="text-2xl font-black mb-8 relative uppercase tracking-tight">Ringkasan</h3>
          
          <div class="space-y-6 relative mb-10">
            <div class="flex justify-between items-start pb-4 border-b border-white/10">
              <span class="text-gray-400 text-[10px] font-black uppercase tracking-widest">Barang</span>
              <span class="text-right font-black text-sm" :class="selectedItem ? 'text-white' : 'text-white/20 italic'">{{ selectedItem?.name || 'Belum dipilih' }}</span>
            </div>
            <div class="flex justify-between items-start pb-4 border-b border-white/10">
              <span class="text-gray-400 text-[10px] font-black uppercase tracking-widest">Penerima</span>
              <span class="text-right font-black text-sm" :class="borrowerType === 'STAFF' ? 'text-white' : (selectedStudent ? 'text-white' : 'text-white/20 italic')">
                {{ borrowerType === 'STAFF' ? 'Diri Sendiri' : (selectedStudent?.full_name || 'Cari siswa...') }}
              </span>
            </div>
            <div class="flex justify-between items-start pb-4 border-b border-white/10">
              <span class="text-gray-400 text-[10px] font-black uppercase tracking-widest">Durasi</span>
              <span class="text-right font-black text-sm">{{ durationDays }} Hari</span>
            </div>
          </div>
          
          <button @click="handleBorrow" :disabled="loading || uploading || !selectedItem || (borrowerType === 'STUDENT' && !selectedStudent)" 
                  class="w-full bg-white text-gray-900 hover:bg-primary-50 py-5 rounded-2xl font-black text-xs uppercase tracking-[0.3em] transition-all active:scale-95 disabled:opacity-20 disabled:cursor-not-allowed shadow-xl">
            {{ loading ? 'MENYIMPAN...' : 'KONFIRMASI SEKARANG' }}
          </button>
          
          <p class="text-center text-[10px] text-gray-500 font-bold uppercase tracking-widest mt-6">SIS-INV • AMANAH & TERTIB</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../utils/api'
import QrScanner from '../components/QrScanner.vue'

const router = useRouter()
const loading = ref(false)
const searchQuery = ref('')
const searchResults = ref([])
const selectedItem = ref(null)

const borrowerType = ref('STUDENT') // DEFAULT to student based on PRD v4 focus
const studentSearch = ref('')
const studentSuggestions = ref([])
const selectedStudent = ref(null)

const durationDays = ref(1)
const purpose = ref('')
const photoUrl = ref('')
const uploading = ref(false)

// Item search
let searchTimeout
const debouncedSearch = () => {
  clearTimeout(searchTimeout)
  if (!searchQuery.value.trim()) { searchResults.value = []; return }
  searchTimeout = setTimeout(async () => {
    try {
      const { data } = await api.get('/items', { params: { search: searchQuery.value, status: 'AVAILABLE' } })
      searchResults.value = data.data?.items || []
    } catch (err) { console.error(err) }
  }, 500)
}

const handleScanned = async (code) => {
  try {
    const { data } = await api.get(`/items/code/${code}`)
    if (data.data.status !== 'AVAILABLE') { 
      alert(`Barang ini sedang ${data.data.status}.`)
      return 
    }
    selectItem(data.data)
  } catch (err) { alert('Kode barang tidak valid.') }
}

const selectItem = (item) => {
  selectedItem.value = item
  searchQuery.value = ''
  searchResults.value = []
}

// Student Suggestion
let studentTimeout
const debouncedStudentSearch = () => {
  clearTimeout(studentTimeout)
  if (studentSearch.value.length < 2) { studentSuggestions.value = []; return }
  studentTimeout = setTimeout(async () => {
    try {
      const { data } = await api.get('/students/search', { params: { q: studentSearch.value } })
      studentSuggestions.value = data.data || []
    } catch (err) { console.error(err) }
  }, 400)
}

const selectStudent = (student) => {
  selectedStudent.value = student
  studentSearch.value = student.full_name
  studentSuggestions.value = []
}

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
      photoUrl.value = res.data.data.url
    }
  } catch (err) { alert('Gagal mengunggah foto kondisi.') }
  finally { uploading.value = false }
}

const handleBorrow = async () => {
  if (!selectedItem.value) return
  loading.value = true
  try {
    const payload = {
      item_code: selectedItem.value.code,
      borrower_type: borrowerType.value,
      expected_return_days: parseFloat(durationDays.value),
      purpose: purpose.value,
      photo_url: photoUrl.value
    }

    if (borrowerType.value === 'STUDENT' && selectedStudent.value) {
      payload.student_nis = selectedStudent.value.nis
      payload.student_name = selectedStudent.value.full_name
      payload.student_class = selectedStudent.value.class
    }

    await api.post('/transactions/borrow', payload)
    router.push('/my-borrowings')
  } catch (err) {
    alert(err.response?.data?.error || 'Gagal memproses peminjaman.')
  } finally { loading.value = false }
}

onMounted(() => {
  // If we come from scanning QR on dashboard / items
  const queryCode = router.currentRoute.value.query.code
  if (queryCode) handleScanned(queryCode)
})
</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.5s ease-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-scale-up {
  animation: scaleUp 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
@keyframes scaleUp {
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: scale(1); }
}
.input-field {
  @apply w-full bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 focus:border-primary-500 focus:ring-4 focus:ring-primary-500/10 transition-all outline-none text-gray-900 dark:text-white font-bold;
}
</style>
