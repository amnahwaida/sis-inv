<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      
      <div class="fixed inset-0 bg-gray-900/40 backdrop-blur-sm transition-opacity" aria-hidden="true" @click="close"></div>

      <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

      <div class="inline-block align-bottom bg-white rounded-2xl text-left overflow-hidden shadow-2xl transform transition-all sm:my-8 sm:align-middle sm:max-w-2xl sm:w-full animate-scale-up">
        <form @submit.prevent="submitForm">
          <!-- Header -->
          <div class="bg-gray-50/50 px-6 py-5 border-b border-gray-100 flex items-center justify-between">
            <div>
              <h3 class="text-xl font-extrabold text-gray-900" id="modal-title">
                {{ isEdit ? 'Perbarui Data Barang' : 'Registrasi Barang Baru' }}
              </h3>
              <p class="text-[10px] text-gray-400 font-bold uppercase tracking-widest mt-1">Master Data Inventaris Sekolah</p>
            </div>
            <button @click="close" type="button" class="text-gray-400 hover:bg-gray-100 p-2 rounded-xl transition-all">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Body -->
          <div class="bg-white px-6 py-6 pb-8">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Kode Barang -->
              <div>
                <label class="block text-xs font-bold text-gray-500 uppercase tracking-widest mb-2 ml-1">Kode Barang *</label>
                <input type="text" v-model="form.code" required :disabled="isEdit"
                       class="input-field rounded-xl border-gray-200 focus:ring-primary-500 font-mono tracking-tight disabled:bg-gray-50 disabled:text-gray-400" 
                       placeholder="CTH: INV-LAPTOP-001">
              </div>

              <!-- Nama Barang -->
              <div>
                <label class="block text-xs font-bold text-gray-500 uppercase tracking-widest mb-2 ml-1">Nama Barang *</label>
                <input type="text" v-model="form.name" required
                       class="input-field rounded-xl border-gray-200 focus:ring-primary-500" 
                       placeholder="Contoh: Asus Vivobook Pro">
              </div>

              <!-- Kategori (REPLACED WITH SELECT) -->
              <div>
                <label class="block text-xs font-bold text-gray-500 uppercase tracking-widest mb-2 ml-1">Kategori Barang</label>
                <select v-model="form.category_id" 
                        class="input-field rounded-xl border-gray-200 focus:ring-primary-500 appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2020%2020%22%3E%3Cpath%20stroke%3D%22%236b7280%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%221.5%22%20d%3D%22m6%208%204%204%204-4%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[right_0.5rem_center] bg-no-repeat pr-10">
                  <option :value="''">-- Pilih Kategori --</option>
                  <option v-for="cat in categoryStore.categories" :key="cat.id" :value="cat.id">
                    {{ cat.name }}
                  </option>
                </select>
                <p class="mt-1.5 ml-1 text-[10px] text-gray-400">Pilih kategori yang sesuai atau buat baru di menu Kategori.</p>
              </div>

              <!-- Lokasi -->
              <div>
                <label class="block text-xs font-bold text-gray-500 uppercase tracking-widest mb-2 ml-1">Lokasi Penyimpanan</label>
                <input type="text" v-model="form.location"
                       class="input-field rounded-xl border-gray-200 focus:ring-primary-500" 
                       placeholder="Contoh: Lab Komputer 1">
              </div>

              <!-- Kondisi -->
              <div>
                <label class="block text-xs font-bold text-gray-500 uppercase tracking-widest mb-2 ml-1">Kondisi Awal</label>
                <select v-model="form.condition" class="input-field rounded-xl border-gray-200 focus:ring-primary-500">
                  <option value="GOOD">Sangat Baik (Normal)</option>
                  <option value="DAMAGED">Rusak / Tidak Fungsi</option>
                  <option value="IN_REPAIR">Dalam Perbaikan</option>
                  <option value="LOST">Barang Hilang</option>
                </select>
              </div>

              <!-- Borrower Type -->
              <div>
                <label class="block text-xs font-bold text-gray-500 uppercase tracking-widest mb-2 ml-1">Aturan Peminjaman</label>
                <select v-model="form.borrower_type" class="input-field rounded-xl border-gray-200 focus:ring-primary-500">
                  <option value="STAFF_ONLY">Hanya Staff/Guru (Internal)</option>
                  <option value="STUDENT_ALLOWED">Boleh Dipinjam Siswa (Via Guru)</option>
                </select>
              </div>

              <!-- Tanggal Pembelian -->
              <div>
                <label class="block text-xs font-bold text-gray-500 uppercase tracking-widest mb-2 ml-1">Tanggal Pembelian</label>
                <input type="date" v-model="form.purchase_date"
                       class="input-field rounded-xl border-gray-200 focus:ring-primary-500">
              </div>

              <!-- Harga -->
              <div>
                <label class="block text-xs font-bold text-gray-500 uppercase tracking-widest mb-2 ml-1">Harga Beli (Rp)</label>
                <input type="number" step="0.01" v-model="form.purchase_price"
                       class="input-field rounded-xl border-gray-200 focus:ring-primary-500" 
                       placeholder="Misal: 5000000">
              </div>

              <!-- Status (Only shown on Edit) -->
              <div v-show="isEdit" class="md:col-span-2">
                <label class="block text-xs font-bold text-gray-500 uppercase tracking-widest mb-2 ml-1">Ubah Status Aset</label>
                <select v-model="form.status" class="input-field rounded-xl border-gray-200 focus:ring-primary-500">
                  <option value="AVAILABLE">Tersedia di Rak</option>
                  <option value="BORROWED">Sedang Dipinjam</option>
                  <option value="MAINTENANCE">Dalam Pemeliharaan</option>
                  <option value="LOST">Status Hilang</option>
                </select>
              </div>

              <!-- Catatan (Full Width) -->
              <div class="md:col-span-2">
                <label class="block text-xs font-bold text-gray-500 uppercase tracking-widest mb-2 ml-1">Spesifikasi / Catatan</label>
                <textarea v-model="form.notes" rows="3"
                       class="input-field rounded-xl border-gray-200 focus:ring-primary-500" 
                       placeholder="Sertakan Serial Number, Model, atau keteragan fisik lainnya..."></textarea>
              </div>
            </div>

            <!-- Error Message -->
            <div v-if="errorMsg" class="mt-8 p-4 bg-red-50 text-red-700 text-xs font-bold rounded-2xl border border-red-100 flex items-center gap-3 animate-shake">
              <svg class="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
              </svg>
              <span>{{ errorMsg }}</span>
            </div>
          </div>
          
          <!-- Footer -->
          <div class="bg-gray-50/80 px-6 py-5 flex gap-3 border-t border-gray-100">
            <button type="button" @click="close" :disabled="loading"
                    class="btn-secondary flex-1 py-3 text-sm font-black uppercase tracking-widest rounded-xl hover:bg-gray-100 transition-all">
              Batal
            </button>
            <button type="submit" :disabled="loading"
                    class="btn-primary flex-1 py-3 text-sm font-black uppercase tracking-widest rounded-xl shadow-xl shadow-primary-200 active:scale-95 transition-all">
              {{ loading ? 'Sinkronisasi...' : 'Simpan Perubahan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed, onMounted } from 'vue'
import { useItemStore } from '../stores/item'
import { useCategoryStore } from '../stores/category'

const props = defineProps({
  isOpen: Boolean,
  itemData: Object // null if creating, object if editing
})

const emit = defineEmits(['close', 'saved'])
const itemStore = useItemStore()
const categoryStore = useCategoryStore()

const loading = ref(false)
const errorMsg = ref('')

const defaultForm = {
  code: '',
  name: '',
  category_id: '',
  location: '',
  condition: 'GOOD',
  borrower_type: 'STAFF_ONLY',
  purchase_date: '',
  purchase_price: '',
  warranty_end_date: '',
  notes: '',
  status: 'AVAILABLE'
}

const form = ref({ ...defaultForm })
const isEdit = computed(() => !!props.itemData)

// Load categories on mount
onMounted(() => {
  categoryStore.fetchCategories()
})

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    errorMsg.value = ''
    categoryStore.fetchCategories() // Refresh categories list
    if (props.itemData) {
      // Map existing item
      form.value = {
        code: props.itemData.code || '',
        name: props.itemData.name || '',
        category_id: props.itemData.category_id || '',
        location: props.itemData.location || '',
        condition: props.itemData.condition || 'GOOD',
        status: props.itemData.status || 'AVAILABLE',
        borrower_type: props.itemData.borrower_type || 'STAFF_ONLY',
        purchase_date: (props.itemData.purchase_date || '').split('T')[0],
        purchase_price: props.itemData.purchase_price || '',
        warranty_end_date: (props.itemData.warranty_end_date || '').split('T')[0],
        notes: props.itemData.notes || ''
      }
    } else {
      // Create new
      form.value = { ...defaultForm }
    }
  }
})

const close = () => {
  if (!loading.value) {
    emit('close')
  }
}

const submitForm = async () => {
  loading.value = true
  errorMsg.value = ''
  
  try {
    const payload = { ...form.value }
    
    // Format dates to ISO-8601 string if not empty, otherwise remove them
    if (payload.purchase_date) {
      payload.purchase_date = new Date(payload.purchase_date).toISOString()
    } else {
      delete payload.purchase_date
    }
    
    if (payload.warranty_end_date) {
      payload.warranty_end_date = new Date(payload.warranty_end_date).toISOString()
    } else {
      delete payload.warranty_end_date
    }

    if (isEdit.value) {
      // Can't update code
      delete payload.code
      await itemStore.updateItem(props.itemData.id, payload)
    } else {
      delete payload.status // Default for create
      await itemStore.createItem(payload)
    }
    
    emit('saved')
    close()
  } catch (e) {
    errorMsg.value = itemStore.error || 'Terjadi kesalahan sistem'
  } finally {
    loading.value = false
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
