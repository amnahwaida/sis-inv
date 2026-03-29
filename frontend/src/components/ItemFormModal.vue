<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="close"></div>

      <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

      <div class="inline-block align-bottom bg-white rounded-xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-2xl sm:w-full">
        <form @submit.prevent="submitForm">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4 border-b border-gray-100">
            <div class="sm:flex sm:items-start">
              <div class="mt-3 text-center sm:mt-0 sm:text-left w-full">
                <h3 class="text-xl leading-6 font-semibold text-gray-900" id="modal-title">
                  {{ isEdit ? 'Edit Barang' : 'Tambah Barang Baru' }}
                </h3>
                
                <div class="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
                  <!-- Kode Barang -->
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Kode Barang *</label>
                    <input type="text" v-model="form.code" required :disabled="isEdit"
                           class="w-full px-4 py-2 bg-gray-50 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors disabled:opacity-60" 
                           placeholder="CTH: INV-001">
                  </div>

                  <!-- Nama Barang -->
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Nama Barang *</label>
                    <input type="text" v-model="form.name" required
                           class="w-full px-4 py-2 bg-gray-50 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors" 
                           placeholder="CTH: Laptop Asus ROG">
                  </div>

                  <!-- Kategori -->
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Kategori (Opsional)</label>
                    <input type="number" v-model="form.category_id"
                           class="w-full px-4 py-2 bg-gray-50 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors" 
                           placeholder="ID Kategori">
                  </div>

                  <!-- Lokasi -->
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Lokasi Penyimpanan</label>
                    <input type="text" v-model="form.location"
                           class="w-full px-4 py-2 bg-gray-50 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors" 
                           placeholder="Ruang Guru Lt 1">
                  </div>

                  <!-- Kondisi -->
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Kondisi</label>
                    <select v-model="form.condition" class="w-full px-4 py-2 bg-gray-50 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors">
                      <option value="GOOD">Baik</option>
                      <option value="DAMAGED">Rusak</option>
                      <option value="IN_REPAIR">Dalam Perbaikan</option>
                      <option value="LOST">Hilang</option>
                    </select>
                  </div>

                  <!-- Borrower Type -->
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Izin Pinjam</label>
                    <select v-model="form.borrower_type" class="w-full px-4 py-2 bg-gray-50 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors">
                      <option value="STAFF_ONLY">Hanya Staff/Guru</option>
                      <option value="STUDENT_ALLOWED">Boleh Dipinjam Siswa (Via Guru)</option>
                    </select>
                  </div>

                  <!-- Tanggal Pembelian -->
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Tanggal Pembelian</label>
                    <input type="date" v-model="form.purchase_date"
                           class="w-full px-4 py-2 bg-gray-50 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors">
                  </div>

                  <!-- Harga -->
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Harga Beli</label>
                    <input type="number" step="0.01" v-model="form.purchase_price"
                           class="w-full px-4 py-2 bg-gray-50 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors" 
                           placeholder="5000000">
                  </div>

                  <!-- Garansi -->
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Garansi Sampai</label>
                    <input type="date" v-model="form.warranty_end_date"
                           class="w-full px-4 py-2 bg-gray-50 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors">
                  </div>
                  
                  <!-- Status (Only shown on Edit) -->
                  <div v-show="isEdit">
                    <label class="block text-sm font-medium text-gray-700 mb-1">Status Barang</label>
                    <select v-model="form.status" class="w-full px-4 py-2 bg-gray-50 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors">
                      <option value="AVAILABLE">Tersedia</option>
                      <option value="BORROWED">Dipinjam</option>
                      <option value="MAINTENANCE">Perbaikan</option>
                      <option value="LOST">Hilang</option>
                    </select>
                  </div>

                  <!-- Catatan (Full Width) -->
                  <div class="md:col-span-2">
                    <label class="block text-sm font-medium text-gray-700 mb-1">Catatan Tambahan</label>
                    <textarea v-model="form.notes" rows="3"
                           class="w-full px-4 py-2 bg-gray-50 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors" 
                           placeholder="Sertakan serial number atau detail lainnya..."></textarea>
                  </div>

                </div>

                <div v-if="errorMsg" class="mt-4 p-3 bg-red-50 text-red-700 text-sm rounded-lg border border-red-100 flex items-center gap-2">
                  <svg class="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
                  </svg>
                  {{ errorMsg }}
                </div>
              </div>
            </div>
          </div>
          
          <div class="bg-gray-50 px-4 py-4 sm:px-6 sm:flex sm:flex-row-reverse rounded-b-xl border-t border-gray-100">
            <button type="submit" :disabled="loading"
                    class="w-full inline-flex justify-center rounded-lg border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 transition-colors">
              {{ loading ? 'Menyimpan...' : 'Simpan' }}
            </button>
            <button type="button" @click="close" :disabled="loading"
                    class="mt-3 w-full inline-flex justify-center rounded-lg border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 transition-colors">
              Batal
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { useItemStore } from '../stores/item'

const props = defineProps({
  isOpen: Boolean,
  itemData: Object // null if creating, object if editing
})

const emit = defineEmits(['close', 'saved'])
const itemStore = useItemStore()
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

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    errorMsg.value = ''
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
    errorMsg.value = itemStore.error || 'Terjadi kesalahan'
  } finally {
    loading.value = false
  }
}
</script>
