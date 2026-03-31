<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Log Perbaikan</h1>
          <p class="text-primary-100/70 text-sm font-medium">Manajemen pemeliharaan & perbaikan aset sekolah</p>
        </div>
        
        <div class="flex items-center gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10">
          <select v-model="filterStatus" @change="fetchLogs" class="bg-transparent border-none text-white text-xs font-bold focus:ring-0 cursor-pointer pr-8">
            <option value="" class="text-gray-900">Semua Status</option>
            <option value="PENDING" class="text-gray-900">Menunggu</option>
            <option value="IN_PROGRESS" class="text-gray-900">Dikerjakan</option>
            <option value="DONE" class="text-gray-900">Selesai</option>
            <option value="CANCELLED" class="text-gray-900">Dibatalkan</option>
          </select>
          <button @click="showCreateModal = true" class="bg-white text-primary-900 hover:bg-primary-50 px-6 py-2.5 rounded-xl text-xs font-black transition-all flex items-center gap-2 shadow-xl shadow-black/10 active:scale-95">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 4v16m8-8H4"/></svg>
            LAPOR KERUSAKAN
          </button>
        </div>
      </div>
    </div>

    <!-- Stats Cards Grid -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-6">
      <div v-for="(count, key) in statsMap" :key="key" 
           class="group relative overflow-hidden bg-white p-6 rounded-[2rem] border border-gray-100 shadow-sm hover:shadow-xl transition-all duration-500">
        <div :class="['absolute top-0 right-0 w-24 h-24 -mt-8 -mr-8 rounded-full blur-2xl opacity-10 group-hover:opacity-20 transition-opacity', count.color]"></div>
        <div class="relative flex items-center gap-4">
          <div :class="['w-12 h-12 rounded-2xl flex items-center justify-center text-white shadow-lg', count.gradient]">
            <component :is="count.icon" class="w-6 h-6" />
          </div>
          <div>
            <p class="text-3xl font-black text-gray-900 leading-none">{{ count.value }}</p>
            <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest mt-1">{{ count.label }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content Card -->
    <div class="bg-white rounded-[2.5rem] border border-gray-100 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="bg-gray-50/50">
              <th v-for="h in ['Barang', 'Masalah', 'Vendor', 'Biaya', 'Status', 'Tanggal', '']" :key="h"
                  class="text-left py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">{{ h }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-if="loading" class="animate-pulse">
              <td colspan="7" class="px-8 py-12 text-center text-gray-400 font-bold uppercase tracking-widest text-xs italic">Menyelaraskan Data...</td>
            </tr>
            <tr v-else-if="logs.length === 0" class="text-center">
              <td colspan="7" class="px-8 py-20">
                <div class="bg-gray-50 w-20 h-20 rounded-full flex items-center justify-center mx-auto mb-4 text-gray-300">
                  <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/></svg>
                </div>
                <p class="text-gray-900 font-black text-lg">Belum Ada Catatan</p>
                <p class="text-gray-400 text-sm mt-1 max-w-xs mx-auto">Semua riwayat perbaikan barang akan muncul di daftar ini.</p>
              </td>
            </tr>
            <tr v-for="log in logs" :key="log.id" class="group hover:bg-gray-50/80 transition-all duration-300">
              <td class="px-8 py-6">
                <div class="flex items-center gap-4">
                  <div class="w-10 h-10 bg-gray-100 rounded-xl flex items-center justify-center text-gray-400 group-hover:bg-primary-100 group-hover:text-primary-600 transition-colors">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>
                  </div>
                  <div>
                    <div class="font-black text-gray-900 text-sm">{{ log.item_name }}</div>
                    <div class="text-[10px] text-gray-400 font-bold font-mono tracking-tighter">{{ log.item_code }}</div>
                  </div>
                </div>
              </td>
              <td class="px-8 py-6">
                <p class="text-xs text-gray-600 leading-relaxed font-medium max-w-[250px]">{{ log.issue_description }}</p>
              </td>
              <td class="px-8 py-6 text-xs font-bold text-gray-700">{{ log.vendor || '---' }}</td>
              <td class="px-8 py-6">
                <span class="text-sm font-black text-gray-900">{{ log.cost ? formatCurrency(log.cost) : 'Rp 0' }}</span>
              </td>
              <td class="px-8 py-6">
                <span class="px-3 py-1 rounded-lg font-black uppercase text-[9px] tracking-widest shadow-sm inline-block" :class="getStatusBadge(log.status)">
                  {{ statusLabels[log.status] }}
                </span>
              </td>
              <td class="px-8 py-6 whitespace-nowrap text-[11px] font-bold text-gray-400">{{ formatDate(log.reported_at) }}</td>
              <td class="px-8 py-6 text-right">
                <div class="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                  <button v-if="log.status === 'PENDING'" @click="updateLog(log, 'IN_PROGRESS')"
                    class="p-2 text-blue-600 bg-blue-50 rounded-xl hover:bg-blue-600 hover:text-white transition-all transform hover:scale-110" title="Kerjakan">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/></svg>
                  </button>
                  <button v-if="log.status === 'IN_PROGRESS'" @click="openFinishModal(log)"
                    class="p-2 text-green-600 bg-green-50 rounded-xl hover:bg-green-600 hover:text-white transition-all transform hover:scale-110" title="Selesaikan">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                  </button>
                  <button @click="deleteLog(log)" class="p-2 text-red-500 bg-red-50 rounded-xl hover:bg-red-600 hover:text-white transition-all transform hover:scale-110" title="Hapus">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create Maintenance Modal -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showCreateModal" class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4" @click.self="showCreateModal = false">
          <div class="bg-white rounded-[3rem] shadow-2xl w-full max-w-xl overflow-hidden animate-scale-up border border-white/20">
            <!-- Modal Header -->
            <div class="px-10 py-8 bg-gradient-to-br from-primary-900 to-black text-white relative">
              <div class="absolute top-0 right-0 w-32 h-32 bg-primary-500/10 rounded-full blur-2xl"></div>
              <h2 class="text-2xl font-black tracking-tight leading-none">Lapor Kerusakan</h2>
              <p class="text-primary-100/60 text-xs font-bold mt-2 tracking-widest uppercase italic">Laporan Pemeliharaan Aset</p>
              <button @click="showCreateModal = false" class="absolute top-8 right-8 text-white/50 hover:text-white transition-colors">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12"/></svg>
              </button>
            </div>

            <div class="p-10 space-y-8 max-h-[70vh] overflow-y-auto custom-scrollbar">
              <!-- Item Search Group -->
              <div class="space-y-4">
                <div class="flex items-center gap-2 mb-2">
                  <div class="w-1.5 h-6 bg-primary-600 rounded-full"></div>
                  <h3 class="text-xs font-black text-gray-900 uppercase tracking-widest">Langkah 1: Pilih Barang</h3>
                </div>
                
                <div class="relative group">
                  <div class="absolute inset-y-0 left-4 flex items-center text-gray-400 group-focus-within:text-primary-600 transition-colors">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
                  </div>
                  <input 
                    v-model="form.search_query" 
                    class="w-full pl-12 pr-4 py-4 rounded-2xl bg-gray-50 border-none focus:bg-white focus:ring-4 focus:ring-primary-100 transition-all font-bold text-sm" 
                    placeholder="Ketik nama atau kode barang..." 
                    @input="debounceItemSearch"
                  />
                  
                  <!-- Floating Suggestion Results -->
                  <div v-if="suggestedItems.length > 0" class="absolute z-50 w-full mt-2 bg-white rounded-3xl shadow-2xl border border-gray-100 overflow-hidden divide-y divide-gray-50 animate-fade-in ring-1 ring-black/5">
                    <div 
                      v-for="item in suggestedItems" :key="item.id"
                      @click="selectItem(item)"
                      class="px-6 py-4 hover:bg-primary-50 cursor-pointer transition-all flex items-center justify-between group/item"
                    >
                      <div class="flex items-center gap-4">
                        <div class="w-10 h-10 bg-gray-100 rounded-xl flex items-center justify-center text-gray-400 group-hover/item:bg-white transition-colors">
                          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>
                        </div>
                        <div>
                          <div class="text-sm font-black text-gray-900">{{ item.name }}</div>
                          <div class="text-[10px] text-gray-400 font-bold font-mono tracking-tighter">{{ item.code }} • {{ item.location }}</div>
                        </div>
                      </div>
                      <div class="text-[10px] bg-primary-100 text-primary-700 px-2 py-1 rounded-lg font-black opacity-0 group-hover/item:opacity-100 transition-opacity">PILIH</div>
                    </div>
                  </div>
                </div>

                <!-- Active Selection Card -->
                <div v-if="form.item_code" class="bg-primary-50 rounded-2xl p-4 flex items-center justify-between border border-primary-200 animate-scale-up">
                  <div class="flex items-center gap-4">
                    <div class="w-10 h-10 bg-primary-600 rounded-xl flex items-center justify-center text-white shadow-lg shadow-primary-200">
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                    </div>
                    <div>
                      <p class="text-[10px] font-black text-primary-400 uppercase tracking-widest leading-none mb-1">Barang Terpilih</p>
                      <p class="text-sm font-black text-primary-900 leading-none">{{ form.search_query }}</p>
                    </div>
                  </div>
                  <button @click="form.item_code = ''; form.search_query = ''" class="p-2 text-primary-400 hover:text-primary-800 transition-colors">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12"/></svg>
                  </button>
                </div>
              </div>

              <!-- Problem Details Group -->
              <div class="space-y-4">
                <div class="flex items-center gap-2 mb-2">
                  <div class="w-1.5 h-6 bg-red-500 rounded-full"></div>
                  <h3 class="text-xs font-black text-gray-900 uppercase tracking-widest">Langkah 2: Detail Kerusakan</h3>
                </div>
                
                <textarea 
                  v-model="form.issue_description" 
                  rows="3" 
                  class="w-full px-6 py-4 rounded-2xl bg-gray-50 border-none focus:bg-white focus:ring-4 focus:ring-red-100 transition-all font-bold text-sm" 
                  placeholder="Ceritakan apa masalahnya... (Misal: Layar berkedip, tidak bisa nyala, dsb)"
                ></textarea>

                <div class="grid grid-cols-2 gap-4">
                  <div class="space-y-1">
                    <label class="text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-2">Vendor / Teknisi</label>
                    <input v-model="form.vendor" class="w-full px-6 py-4 rounded-2xl bg-gray-50 border-none focus:bg-white focus:ring-4 focus:ring-primary-100 transition-all font-bold text-sm" placeholder="Opsional" />
                  </div>
                  <div class="space-y-1">
                    <label class="text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-2">Estimasi Biaya</label>
                    <div class="relative">
                      <div class="absolute inset-y-0 left-4 flex items-center text-gray-400 font-black text-xs">Rp</div>
                      <input v-model.number="form.cost" type="number" class="w-full pl-10 pr-6 py-4 rounded-2xl bg-gray-50 border-none focus:bg-white focus:ring-4 focus:ring-primary-100 transition-all font-bold text-sm" placeholder="0" />
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Modal Footer -->
            <div class="p-10 pt-0 flex gap-4">
              <button @click="showCreateModal = false" class="flex-1 px-8 py-4 rounded-2xl text-xs font-black text-gray-500 hover:bg-gray-50 transition-all uppercase tracking-widest">Batal</button>
              <button 
                @click="submitCreate" 
                :disabled="submitting || !form.item_code || !form.issue_description" 
                class="flex-[2] bg-primary-900 text-white px-8 py-4 rounded-2xl text-xs font-black shadow-2xl shadow-primary-900/40 translate-y-0 active:translate-y-1 transition-all disabled:opacity-50 disabled:translate-y-0 uppercase tracking-widest"
              >
                {{ submitting ? 'MEMPROSES...' : 'KIRIM LAPORAN' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Finish / Done Modal -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showFinishModal" class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4" @click.self="showFinishModal = false">
          <div class="bg-white rounded-[3rem] shadow-2xl w-full max-w-lg overflow-hidden animate-scale-up">
            <div class="px-10 py-8 bg-gradient-to-br from-green-600 to-green-900 text-white relative">
              <h2 class="text-2xl font-black tracking-tight leading-none">Selesaikan Perbaikan</h2>
              <div class="bg-white/20 backdrop-blur-sm rounded-2xl p-4 mt-6 flex items-center gap-4">
                <div class="w-12 h-12 bg-white rounded-xl flex items-center justify-center text-green-700 shadow-lg">
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                </div>
                <div>
                  <div class="text-xs font-black text-white uppercase tracking-widest opacity-60">Barang</div>
                  <div class="text-lg font-black leading-none">{{ finishTarget?.item_name }}</div>
                </div>
              </div>
            </div>

            <div class="p-10 space-y-6">
              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-1">
                  <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest ml-2">Final Vendor</label>
                  <input v-model="finishForm.vendor" class="w-full px-6 py-4 rounded-2xl bg-gray-50 border-none focus:ring-4 focus:ring-green-100 transition-all font-bold text-sm" />
                </div>
                <div class="space-y-1">
                  <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest ml-2">Total Biaya Akhir</label>
                  <div class="relative">
                    <div class="absolute inset-y-0 left-4 flex items-center text-gray-400 font-bold text-xs">Rp</div>
                    <input v-model.number="finishForm.cost" type="number" class="w-full pl-10 pr-6 py-4 rounded-2xl bg-gray-50 border-none focus:ring-4 focus:ring-green-100 transition-all font-bold text-sm" />
                  </div>
                </div>
              </div>
              <div class="space-y-1">
                <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest ml-2">Hasil Akhir / Catatan</label>
                <textarea v-model="finishForm.notes" rows="2" class="w-full px-6 py-4 rounded-2xl bg-gray-50 border-none focus:ring-4 focus:ring-green-100 transition-all font-bold text-sm" placeholder="CTH: Komponen IC diganti, layar sudah normal..."></textarea>
              </div>
            </div>

            <div class="px-10 pb-10 flex gap-4">
              <button @click="showFinishModal = false" class="flex-1 px-8 py-4 rounded-2xl text-xs font-black text-gray-400 hover:bg-gray-50 transition-all uppercase tracking-widest">Batal</button>
              <button 
                @click="submitFinish" 
                :disabled="submitting"
                class="flex-[2] bg-green-600 text-white px-8 py-4 rounded-2xl text-xs font-black shadow-2xl shadow-green-200 transition-all active:scale-95 uppercase tracking-widest"
              >
                {{ submitting ? 'MENYIMPAN...' : 'KONFIRMASI SELESAI' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, h } from 'vue'
import api from '../utils/api'

// Icons for stats cards
const IconPending = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-6 h-6' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z' })]) }
const IconProgress = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-6 h-6' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z' })]) }
const IconDone = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-6 h-6' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M5 13l4 4L19 7' })]) }
const IconCancel = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-6 h-6' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M6 18L18 6M6 6l12 12' })]) }

const logs = ref([])
const loading = ref(false)
const filterStatus = ref('')
const showCreateModal = ref(false)
const showFinishModal = ref(false)
const submitting = ref(false)
const formError = ref('')
const finishTarget = ref(null)

const form = ref({ item_code: '', search_query: '', issue_description: '', cost: 0, vendor: '', notes: '' })
const finishForm = ref({ cost: 0, vendor: '', notes: '' })

const suggestedItems = ref([])
let searchTimeout = null

const statusLabels = { PENDING: 'Menunggu', IN_PROGRESS: 'Dikerjakan', DONE: 'Selesai', CANCELLED: 'Dibatalkan' }

const stats = computed(() => {
  const all = logs.value
  return {
    pending: all.filter(l => l.status === 'PENDING').length,
    inProgress: all.filter(l => l.status === 'IN_PROGRESS').length,
    done: all.filter(l => l.status === 'DONE').length,
    cancelled: all.filter(l => l.status === 'CANCELLED').length,
  }
})

const statsMap = computed(() => ({
  pending: { label: 'Menunggu', value: stats.value.pending, color: 'bg-yellow-500', gradient: 'bg-gradient-to-br from-yellow-400 to-orange-500', icon: IconPending },
  progress: { label: 'Dikerjakan', value: stats.value.inProgress, color: 'bg-blue-500', gradient: 'bg-gradient-to-br from-blue-400 to-indigo-600', icon: IconProgress },
  done: { label: 'Selesai', value: stats.value.done, color: 'bg-green-500', gradient: 'bg-gradient-to-br from-green-400 to-emerald-600', icon: IconDone },
  cancelled: { label: 'Batal', value: stats.value.cancelled, color: 'bg-gray-500', gradient: 'bg-gradient-to-br from-gray-400 to-slate-600', icon: IconCancel },
}))

async function fetchLogs() {
  loading.value = true
  try {
    const params = filterStatus.value ? { status: filterStatus.value } : {}
    const { data } = await api.get('/maintenance', { params })
    if (data.success) logs.value = data.data
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

const debounceItemSearch = () => {
  clearTimeout(searchTimeout)
  if (!form.value.search_query) {
    suggestedItems.value = []
    return
  }
  searchTimeout = setTimeout(async () => {
    try {
      const { data } = await api.get(`/items?search=${form.value.search_query}&page_size=10`)
      if (data.success) {
        // Only show items that are NOT currently in maintenance or borrowed (optional, but safer)
        suggestedItems.value = data.data.items
      }
    } catch (e) { console.error(e) }
  }, 300)
}

const selectItem = (item) => {
  form.value.item_code = item.code
  form.value.search_query = item.name
  suggestedItems.value = []
}

async function submitCreate() {
  if (!form.value.item_code || !form.value.issue_description) {
    formError.value = 'Barang dan deskripsi kerusakan wajib diisi.'
    return
  }
  formError.value = ''
  submitting.value = true
  try {
    const { data } = await api.post('/maintenance', {
      item_code: form.value.item_code,
      issue_description: form.value.issue_description,
      cost: form.value.cost,
      vendor: form.value.vendor,
      notes: form.value.notes
    })
    if (data.success) {
      showCreateModal.value = false
      form.value = { item_code: '', search_query: '', issue_description: '', cost: 0, vendor: '', notes: '' }
      fetchLogs()
    }
  } catch (e) {
    formError.value = e.response?.data?.error || 'Gagal menyimpan data'
  } finally { submitting.value = false }
}

async function updateLog(log, newStatus) {
  if (!confirm(`Ubah status "${statusLabels[log.status]}" → "${statusLabels[newStatus]}"?`)) return
  try {
    await api.put(`/maintenance/${log.id}`, { status: newStatus, cost: log.cost || 0, vendor: log.vendor || '', notes: log.notes || '' })
    fetchLogs()
  } catch (e) { alert('Gagal mengubah status') }
}

function openFinishModal(log) {
  finishTarget.value = log
  finishForm.value = { cost: log.cost || 0, vendor: log.vendor || '', notes: log.notes || '' }
  showFinishModal.value = true
}

async function submitFinish() {
  submitting.value = true
  try {
    await api.put(`/maintenance/${finishTarget.value.id}`, { status: 'DONE', ...finishForm.value })
    showFinishModal.value = false
    fetchLogs()
  } catch (e) { alert('Gagal menyelesaikan perbaikan') }
  finally { submitting.value = false }
}

async function deleteLog(log) {
  if (!confirm(`Hapus log perbaikan untuk "${log.item_name}"?`)) return
  try {
    await api.delete(`/maintenance/${log.id}`)
    fetchLogs()
  } catch (e) { alert('Gagal menghapus') }
}

const formatDate = (d) => d ? new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }).format(new Date(d)) : '-'
const formatCurrency = (v) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)
const getStatusBadge = (s) => ({
  PENDING: 'bg-yellow-50 text-yellow-700 border border-yellow-200',
  IN_PROGRESS: 'bg-blue-50 text-blue-700 border border-blue-200',
  DONE: 'bg-green-50 text-green-700 border border-green-200',
  CANCELLED: 'bg-gray-100 text-gray-500 border border-gray-200',
}[s] || 'bg-gray-50 text-gray-700')

onMounted(fetchLogs)
</script>
