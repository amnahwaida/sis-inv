<template>
  <div class="animate-fade-in space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Kelola Lokasi</h1>
        <p class="text-sm text-gray-500 mt-1 uppercase tracking-widest font-bold text-[10px]">Tentukan lokasi penyimpanan inventaris sekolah</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Create Panel -->
      <div class="lg:col-span-1">
        <div class="card p-6 border border-gray-100 shadow-sm sticky top-24">
          <h3 class="text-lg font-bold text-gray-900 mb-6">Tambah Lokasi</h3>
          <form @submit.prevent="saveLocation" class="space-y-4">
            <div>
              <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-2 ml-1">Nama Lokasi *</label>
              <input v-model="form.name" required class="input-field rounded-xl border-gray-200" placeholder="Contoh: Lab Komputer 1, Ruang Guru, dsb" />
            </div>
            <div>
              <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-2 ml-1">Deskripsi / Detail</label>
              <textarea v-model="form.description" rows="3" class="input-field rounded-xl border-gray-200" placeholder="Detail lokasi (gedung, lantai, no ruangan)..."></textarea>
            </div>
            <div v-if="locationStore.error" class="p-3 bg-red-50 text-red-700 text-xs rounded-xl font-medium border border-red-100">
              {{ locationStore.error }}
            </div>
            <button type="submit" :disabled="locationStore.loading" class="btn-primary w-full py-3 rounded-xl font-black uppercase text-xs tracking-widest shadow-lg shadow-primary-100">
              {{ locationStore.loading ? 'Menyimpan...' : 'Simpan Lokasi' }}
            </button>
          </form>
        </div>
      </div>

      <!-- List Panel -->
      <div class="lg:col-span-2">
        <div class="card p-0 overflow-hidden border border-gray-100 shadow-sm min-h-[500px] relative font-sans">
          <!-- Loading Overlay -->
          <div v-if="locationStore.loading && locationStore.locations.length === 0" class="absolute inset-0 bg-white/60 backdrop-blur-sm z-10 flex flex-col items-center justify-center">
            <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-primary-600"></div>
          </div>

          <div class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-gray-50/80 border-b border-gray-100 font-sans">
                <tr>
                  <th class="text-left text-xs font-bold text-gray-500 uppercase tracking-widest px-6 py-4">Nama & Lokasi</th>
                  <th class="text-right text-xs font-bold text-gray-500 uppercase tracking-widest px-6 py-4 w-24">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100">
                <template v-if="locationStore.locations.length > 0">
                  <tr v-for="loc in locationStore.locations" :key="loc.id" class="hover:bg-gray-50/50 transition-colors group">
                    <td class="px-6 py-4">
                      <div class="flex items-center gap-4">
                        <div class="w-10 h-10 rounded-xl bg-gray-100 flex items-center justify-center text-gray-400 group-hover:bg-primary-50 group-hover:text-primary-600 transition-colors">
                          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                          </svg>
                        </div>
                        <div>
                          <p class="font-bold text-gray-900 group-hover:text-primary-700 transition-colors">{{ loc.name }}</p>
                          <p class="text-[11px] text-gray-400 line-clamp-1 mt-0.5">{{ loc.description || 'Area sekolah ' + loc.name }}</p>
                        </div>
                      </div>
                    </td>
                    <td class="px-6 py-4 text-right">
                      <button @click="deleteLocation(loc)" class="p-2 text-gray-300 hover:text-red-600 hover:bg-red-50 rounded-xl transition-all" title="Hapus Lokasi">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.895-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                      </button>
                    </td>
                  </tr>
                </template>
                <tr v-else-if="!locationStore.loading" class="text-center font-sans">
                  <td colspan="2" class="px-6 py-24 text-gray-400 italic text-sm">Belum ada lokasi terdaftar.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useLocationStore } from '../stores/location'

const locationStore = useLocationStore()

const form = ref({
  name: '',
  description: ''
})

onMounted(() => {
  locationStore.fetchLocations()
})

const saveLocation = async () => {
  try {
    await locationStore.createLocation(form.value)
    form.value.name = ''
    form.value.description = ''
  } catch (e) {
    // Error is in store
  }
}

const deleteLocation = async (loc) => {
  if (confirm(`Hapus lokasi "${loc.name}"? (Hanya bisa dihapus jika tidak ada barang di lokasi ini)`)) {
    try {
      await locationStore.deleteLocation(loc.id)
    } catch (e) {
      alert(locationStore.error || 'Gagal menghapus lokasi')
    }
  }
}
</script>
