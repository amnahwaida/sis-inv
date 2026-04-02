<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="close"></div>

      <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

      <div class="inline-block align-bottom bg-white dark:bg-gray-800 rounded-xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-md sm:w-full border border-transparent dark:border-gray-700">
        <div class="bg-white dark:bg-gray-800 px-4 pt-5 pb-4 sm:p-6 sm:pb-4 border-b border-gray-100 dark:border-gray-700">
          <div class="sm:flex sm:items-start">
            <div class="mt-3 text-center sm:mt-0 sm:text-left w-full">
              <div class="flex justify-between items-center mb-4">
                <h3 class="text-xl leading-6 font-semibold text-gray-900 dark:text-white" id="modal-title">
                  QR Code Barang
                </h3>
                <button @click="close" class="text-gray-400 hover:text-gray-500">
                  <span class="sr-only">Close</span>
                  <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              
              <div v-if="loading" class="flex flex-col items-center justify-center py-12">
                <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mb-4"></div>
                <p class="text-gray-500">Membuat QR Code...</p>
              </div>

              <div v-else-if="errorMsg" class="p-4 bg-red-50 text-red-700 rounded-lg flex items-center gap-3">
                <svg class="h-6 w-6 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
                </svg>
                {{ errorMsg }}
              </div>

              <div v-else-if="qrCodeData" class="flex flex-col items-center justify-center py-6" id="print-area">
                <div class="bg-white dark:bg-gray-800 p-4 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-xl mb-4 text-center">
                  <p class="text-sm font-bold text-gray-800 dark:text-gray-200 mb-2 uppercase">{{ itemData?.name }}</p>
                  <img :src="qrCodeData" alt="QR Code" class="w-48 h-48 mx-auto" />
                  <p class="text-xs font-mono text-gray-600 mt-2 tracking-widest">{{ itemData?.code }}</p>
                </div>
                
                <p class="text-sm text-gray-500 text-center">
                  Cetak label ini dan tempelkan pada fisik barang untuk memudahkan inventarisasi dan peminjaman.
                </p>
              </div>
            </div>
          </div>
        </div>
        
        <div class="bg-gray-50 dark:bg-gray-800/80 px-6 py-5 flex flex-col-reverse sm:flex-row sm:justify-end gap-3 rounded-b-xl border-t border-gray-100 dark:border-gray-700">
          <button type="button" @click="close"
                  class="btn-secondary flex-1 sm:flex-none">
            Tutup
          </button>
          
          <button type="button" @click="printQR" :disabled="loading || !qrCodeData"
                  class="btn-premium-action flex-1 sm:flex-none !px-8">
            <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
            </svg>
            CETAK LABEL
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import api from '../utils/api'

const props = defineProps({
  isOpen: Boolean,
  itemData: Object
})

const emit = defineEmits(['close'])

const loading = ref(false)
const errorMsg = ref('')
const qrCodeData = ref('')

watch(() => props.isOpen, async (newVal) => {
  if (newVal && props.itemData) {
    if (props.itemData.qr_code_data && props.itemData.qr_code_data.startsWith('data:image')) {
       // Ideally we store the base64, but if not we regenerate
       // For MVP the schema stores the literal code in `qr_code_data`, so we must generate it via API
    }
    await generateQrCode()
  } else {
    qrCodeData.value = ''
    errorMsg.value = ''
  }
})

const generateQrCode = async () => {
  loading.value = true
  errorMsg.value = ''
  qrCodeData.value = ''
  
  try {
    const { data } = await api.post(`/items/${props.itemData.id}/qr`)
    if (data.success) {
      qrCodeData.value = data.data.qr_code
    }
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal membuat QR Code'
  } finally {
    loading.value = false
  }
}

const close = () => {
  emit('close')
}

const printQR = () => {
  const printContent = document.getElementById('print-area').innerHTML
  const originalContents = document.body.innerHTML

  // Create a minimal printable version
  // We apply tailwind classes inline or use basic CSS for printing
  const printWindow = window.open('', '_blank', 'width=800,height=600')
  printWindow.document.write(`
    <html>
      <head>
        <title>Print QR Label - ${props.itemData.code}</title>
        <style>
          body { font-family: system-ui, -apple-system, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
          .label { border: 2px solid #000; padding: 20px; text-align: center; width: 250px; border-radius: 8px; }
          .name { font-weight: bold; font-size: 14px; margin-bottom: 10px; text-transform: uppercase; }
          .img { width: 150px; height: 150px; margin: 0 auto; }
          .code { font-family: monospace; font-size: 12px; margin-top: 10px; letter-spacing: 2px; }
          @media print {
            body { align-items: flex-start; margin-top: 20px; }
          }
        </style>
      </head>
      <body>
        <div class="label">
          <div class="name">${props.itemData.name}</div>
          <img src="${qrCodeData.value}" class="img" />
          <div class="code">${props.itemData.code}</div>
        </div>
        <script>
          window.onload = () => { window.print(); window.close(); }
        <\/script>
      </body>
    </html>
  `)
  printWindow.document.close()
}
</script>
