<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-primary-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Pinjam Barang</h1>
          <p class="text-primary-100/70 text-sm font-medium">Proses pencatatan peminjaman aset baru</p>
        </div>
        
        <div class="flex items-center gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10">
          <router-link to="/my-borrowings" class="bg-white/10 hover:bg-white/20 text-white border border-white/10 px-5 py-2.5 rounded-xl text-xs font-black transition-all flex items-center gap-2 active:scale-95">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h3.75M9 15h3.75M9 18h3.75m3 .75H18a2.25 2.25 0 002.25-2.25V6.108c0-1.135-.845-2.098-1.976-2.192a48.424 48.424 0 00-1.123-.08m-5.801 0c-.065.21-.1.433-.1.664 0 .414.336.75.75.75h4.5a.75.75 0 00.75-.75 2.25 2.25 0 00-.1-.664m-5.8 0A2.251 2.251 0 0113.5 2.25H15c1.012 0 1.867.668 2.15 1.586m-5.8 0c-.376.023-.75.05-1.124.08C9.095 4.01 8.25 4.973 8.25 6.108V18.25m0 0c0 1.242 1.008 2.25 2.25 2.25h.75" />
            </svg>
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
        <section class="card-premium group">
          <div class="card-decoration"></div>
          <div class="flex items-center gap-4 mb-8">
            <div class="w-10 h-10 bg-primary-600 text-white rounded-xl flex items-center justify-center font-black">1</div>
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
              <svg class="w-6 h-6 absolute left-4 top-4 text-gray-300" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
              </svg>
            </div>
            
            <!-- Search Results Dropdown -->
            <div v-if="searchResults.length > 0" class="bg-white dark:bg-gray-800 rounded-2xl border border-gray-200 dark:border-gray-700 shadow-xl overflow-hidden animate-slide-in">
              <div v-for="item in searchResults" :key="item.id" @click="selectItem(item)"
                   class="p-4 hover:bg-primary-50 dark:hover:bg-primary-900/20 cursor-pointer border-b border-gray-50 dark:border-gray-700 last:border-0 flex items-center justify-between transition-colors">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 bg-gray-100 dark:bg-gray-700 rounded-lg flex items-center justify-center text-gray-400">
                    <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-10.5v10.5" />
                    </svg>
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
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
            <div class="flex items-center gap-4">
              <div class="w-16 h-16 bg-white dark:bg-gray-800 text-primary-600 rounded-2xl flex items-center justify-center shadow-lg">
                <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-10.5v10.5" />
                </svg>
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
        <section class="card-premium group !overflow-visible">
          <div class="card-decoration"></div>
          <div class="flex items-center gap-4 mb-8">
            <div class="w-10 h-10 bg-primary-600 text-white rounded-xl flex items-center justify-center font-black">2</div>
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
                        :class="borrowerType === type.id ? 'bg-primary-600 text-white shadow-xl shadow-primary-500/20 scale-[1.02]' : 'text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'">
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
                <svg class="w-6 h-6 absolute left-4 top-4 text-gray-300" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" />
                </svg>
                
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
              <div v-if="selectedStudent" class="animate-scale-up p-5 bg-primary-50 dark:bg-primary-900/20 rounded-2xl border border-primary-100 dark:border-primary-800 flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 bg-white dark:bg-gray-800 text-primary-600 rounded-xl flex items-center justify-center shadow-md">
                    <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                    </svg>
                  </div>
                  <div>
                    <div class="font-black text-gray-900 dark:text-white text-sm">{{ selectedStudent.full_name }}</div>
                    <div class="text-[10px] font-bold text-primary-600 uppercase tracking-widest">{{ selectedStudent.nis }} • {{ selectedStudent.class }}</div>
                  </div>
                </div>
                <button @click="selectedStudent = null; studentSearch = ''" class="text-primary-400 hover:text-primary-600 transition-colors">
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                  </svg>
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
                    <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M6.827 6.175A2.31 2.31 0 015.186 7.23c-.38.054-.757.112-1.134.175C2.999 7.58 2.25 8.507 2.25 9.574V18a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9.574c0-1.067-.75-1.994-1.802-2.169a47.865 47.865 0 00-1.134-.175 2.31 2.31 0 01-1.64-1.055l-.822-1.316a2.192 2.192 0 00-1.736-1.039 48.774 48.774 0 00-5.232 0 2.192 2.192 0 00-1.736 1.039l-.821 1.316z" />
                      <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 12.75a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0zM18.75 10.5h.008v.008h-.008V10.5z" />
                    </svg>
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
        <div class="bg-white dark:bg-gray-800 text-gray-900 dark:text-white rounded-[2.5rem] p-10 shadow-2xl relative overflow-hidden border border-gray-100 dark:border-gray-700 transition-all duration-300">
          <div class="absolute top-0 right-0 w-48 h-48 bg-primary-600/10 dark:bg-primary-600/20 rounded-full blur-3xl"></div>
          
          <h3 class="text-2xl font-black mb-8 relative uppercase tracking-tight text-gray-900 dark:text-white">Ringkasan</h3>
          
          <div class="space-y-6 relative mb-10">
            <div class="flex justify-between items-start pb-4 border-b border-gray-100 dark:border-gray-700">
              <span class="text-gray-400 dark:text-gray-500 text-[10px] font-black uppercase tracking-widest">Barang</span>
              <span class="text-right font-black text-sm" :class="selectedItem ? 'text-gray-900 dark:text-white' : 'text-gray-300 dark:text-white/20 italic'">{{ selectedItem?.name || 'Belum dipilih' }}</span>
            </div>
            <div class="flex justify-between items-start pb-4 border-b border-gray-100 dark:border-gray-700">
              <span class="text-gray-400 dark:text-gray-500 text-[10px] font-black uppercase tracking-widest">Penerima</span>
              <span class="text-right font-black text-sm" :class="(borrowerType === 'STAFF' || selectedStudent) ? 'text-gray-900 dark:text-white' : 'text-gray-300 dark:text-white/20 italic'">
                {{ borrowerType === 'STAFF' ? 'Diri Sendiri' : (selectedStudent?.full_name || 'Cari siswa...') }}
              </span>
            </div>
            <div class="flex justify-between items-start pb-4 border-b border-gray-100 dark:border-gray-700">
              <span class="text-gray-400 dark:text-gray-500 text-[10px] font-black uppercase tracking-widest">Durasi</span>
              <span class="text-right font-black text-sm text-gray-900 dark:text-white">{{ durationDays }} Hari</span>
            </div>
          </div>
          
          <button @click="handleBorrow" :disabled="loading || uploading || !selectedItem || (borrowerType === 'STUDENT' && !selectedStudent)" 
                   class="btn-premium-action w-full py-5">
            {{ loading ? 'MENYIMPAN...' : 'KONFIRMASI SEKARANG' }}
          </button>
          
          <p class="text-center text-[10px] text-gray-400 dark:text-gray-500 font-bold uppercase tracking-widest mt-6">SIS-INV • AMANAH & TERTIB</p>
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
