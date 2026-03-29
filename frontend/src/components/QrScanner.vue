<template>
  <div class="qr-scanner-container relative">
    <div v-if="loading" class="absolute inset-0 z-10 bg-white/80 flex flex-col items-center justify-center">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600 mb-2"></div>
      <p class="text-sm text-gray-500">Mempersiapkan kamera...</p>
    </div>
    
    <div v-if="error" class="p-4 bg-red-50 text-red-700 rounded-lg text-sm mb-4">
      {{ error }}
      <button @click="initScanner" class="mt-2 text-xs font-semibold underline block">Coba Lagi</button>
    </div>

    <!-- The element where HTML5 QRCode will render the video -->
    <div id="qr-reader" class="w-full max-w-sm mx-auto overflow-hidden rounded-xl border-2 border-dashed border-gray-300 bg-black aspect-square"></div>
    
    <!-- Manual Input Fallback -->
    <div class="mt-4 text-center">
      <p class="text-xs text-gray-500 mb-2">Atau masukkan kode barang manual:</p>
      <div class="flex max-w-xs mx-auto">
        <input type="text" v-model="manualCode" placeholder="CTH: INV-001" 
               class="flex-1 min-w-0 block w-full px-3 py-2 rounded-l-md border border-gray-300 focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" />
        <button @click="submitManual" 
                class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-r-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
          Cari
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { Html5Qrcode } from 'html5-qrcode'

const props = defineProps({
  active: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['scanned'])

const loading = ref(true)
const error = ref('')
const manualCode = ref('')
let html5QrCode = null

const initScanner = async () => {
  loading.value = true
  error.value = ''
  
  try {
    // Check if permissions exist
    const devices = await Html5Qrcode.getCameras()
    if (devices && devices.length) {
      startScanning()
    } else {
      error.value = 'Kamera tidak ditemukan pada perangkat ini.'
      loading.value = false
    }
  } catch (err) {
    error.value = 'Akses kamera ditolak atau tidak didukung oleh browser Anda. Silakan gunakan input manual.'
    loading.value = false
  }
}

const startScanning = () => {
  if (html5QrCode) {
    html5QrCode.stop().catch(() => {})
  }
  
  html5QrCode = new Html5Qrcode("qr-reader")
  
  const qrCodeSuccessCallback = (decodedText, decodedResult) => {
    // Stop scanning once we got a result, to prevent multiple rapid emits
    stopScanning()
    emit('scanned', decodedText)
  }
  
  const config = { fps: 10, qrbox: { width: 250, height: 250 } }
  
  // Try back camera first
  html5QrCode.start({ facingMode: "environment" }, config, qrCodeSuccessCallback)
    .then(() => {
      loading.value = false
    })
    .catch((err) => {
      // If back camera fails, try front camera or let it fail
      error.value = 'Gagal memulai pemindaian QR Code.'
      loading.value = false
    })
}

const stopScanning = () => {
  if (html5QrCode) {
    try {
      html5QrCode.stop()
    } catch(err) {
      // ignore
    }
  }
}

const submitManual = () => {
  if (manualCode.value.trim()) {
    emit('scanned', manualCode.value.trim())
    manualCode.value = ''
  }
}

onMounted(() => {
  if (props.active) {
    initScanner()
  }
})

onBeforeUnmount(() => {
  stopScanning()
})
</script>

<style>
/* Clean up the default styling from html5-qrcode */
#qr-reader {
  border: none !important;
}
#qr-reader img {
  display: none !important; /* Hide the default scanner logo */
}
#qr-reader__dashboard_section_csr span {
  display: none !important;
}
#qr-reader__dashboard_section_swaplink {
  display: none !important;
}
</style>
