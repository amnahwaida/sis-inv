<template>
  <div class="animate-fade-in space-y-6 max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-gray-900 tracking-tight">Log Perbaikan Barang</h1>
        <p class="text-sm text-gray-500 mt-1">Catat, pantau, dan selesaikan perbaikan aset sekolah.</p>
      </div>
      <div class="flex items-center gap-2">
        <!-- Filter Status -->
        <select v-model="filterStatus" @change="fetchLogs" class="input text-xs py-2 px-3 rounded-xl">
          <option value="">Semua Status</option>
          <option value="PENDING">Menunggu</option>
          <option value="IN_PROGRESS">Sedang Dikerjakan</option>
          <option value="DONE">Selesai</option>
          <option value="CANCELLED">Dibatalkan</option>
        </select>
        <button @click="showCreateModal = true" class="btn-primary text-xs flex items-center gap-2 px-4 py-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
          Lapor Kerusakan
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="card p-4 text-center">
        <p class="text-2xl font-black text-yellow-500">{{ stats.pending }}</p>
        <p class="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-1">Menunggu</p>
      </div>
      <div class="card p-4 text-center">
        <p class="text-2xl font-black text-blue-500">{{ stats.inProgress }}</p>
        <p class="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-1">Dikerjakan</p>
      </div>
      <div class="card p-4 text-center">
        <p class="text-2xl font-black text-green-500">{{ stats.done }}</p>
        <p class="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-1">Selesai</p>
      </div>
      <div class="card p-4 text-center">
        <p class="text-2xl font-black text-gray-400">{{ stats.cancelled }}</p>
        <p class="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-1">Dibatalkan</p>
      </div>
    </div>

    <!-- Logs Table -->
    <div class="card p-0 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50/50 border-b border-gray-100">
            <tr>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-5 py-4">Barang</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-5 py-4">Masalah</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-5 py-4">Vendor</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-5 py-4">Biaya</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-5 py-4">Status</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-5 py-4">Tanggal</th>
              <th class="text-right text-[10px] font-black text-gray-400 uppercase tracking-widest px-5 py-4">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-if="loading" class="text-center">
              <td colspan="7" class="px-6 py-12 text-gray-400 italic">Memuat data...</td>
            </tr>
            <tr v-else-if="logs.length === 0" class="text-center">
              <td colspan="7" class="px-6 py-12">
                <p class="text-gray-900 font-bold">Tidak ada log perbaikan.</p>
                <p class="text-gray-400 text-xs mt-1">Klik "Lapor Kerusakan" untuk menambahkan.</p>
              </td>
            </tr>
            <tr v-for="log in logs" :key="log.id" class="hover:bg-gray-50/50 transition-colors text-xs">
              <td class="px-5 py-4">
                <div class="font-bold text-gray-900">{{ log.item_name }}</div>
                <div class="text-[10px] text-gray-400 font-mono">{{ log.item_code }}</div>
              </td>
              <td class="px-5 py-4 text-gray-700 max-w-[200px] truncate" :title="log.issue_description">
                {{ log.issue_description }}
              </td>
              <td class="px-5 py-4 text-gray-600">{{ log.vendor || '-' }}</td>
              <td class="px-5 py-4 font-bold text-gray-900">{{ log.cost ? formatCurrency(log.cost) : '-' }}</td>
              <td class="px-5 py-4">
                <span class="px-2 py-0.5 rounded-md font-black uppercase text-[9px] tracking-widest" :class="getStatusBadge(log.status)">
                  {{ statusLabels[log.status] }}
                </span>
              </td>
              <td class="px-5 py-4 whitespace-nowrap text-gray-500">{{ formatDate(log.reported_at) }}</td>
              <td class="px-5 py-4 text-right">
                <div class="flex items-center justify-end gap-1">
                  <button v-if="log.status === 'PENDING'" @click="updateLog(log, 'IN_PROGRESS')"
                    class="px-2 py-1 text-blue-600 bg-blue-50 rounded-lg hover:bg-blue-100 transition-colors text-[10px] font-bold">
                    Kerjakan
                  </button>
                  <button v-if="log.status === 'IN_PROGRESS'" @click="openFinishModal(log)"
                    class="px-2 py-1 text-green-600 bg-green-50 rounded-lg hover:bg-green-100 transition-colors text-[10px] font-bold">
                    Selesai
                  </button>
                  <button v-if="log.status === 'PENDING' || log.status === 'IN_PROGRESS'" @click="updateLog(log, 'CANCELLED')"
                    class="px-2 py-1 text-gray-500 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors text-[10px] font-bold">
                    Batal
                  </button>
                  <button @click="deleteLog(log)" class="px-2 py-1 text-red-500 bg-red-50 rounded-lg hover:bg-red-100 transition-colors text-[10px] font-bold">
                    Hapus
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create Modal -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showCreateModal" class="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4" @click.self="showCreateModal = false">
          <div class="bg-white rounded-3xl shadow-2xl w-full max-w-md p-6 space-y-5">
            <h2 class="text-lg font-black text-gray-900">Lapor Kerusakan Barang</h2>

            <div>
              <label class="label">Kode Barang</label>
              <input v-model="form.item_code" class="input" placeholder="Contoh: PC-LAB-001" />
            </div>
            <div>
              <label class="label">Deskripsi Kerusakan</label>
              <textarea v-model="form.issue_description" class="input" rows="3" placeholder="Jelaskan kerusakannya..."></textarea>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="label">Vendor / Tukang</label>
                <input v-model="form.vendor" class="input" placeholder="Nama tukang/toko" />
              </div>
              <div>
                <label class="label">Estimasi Biaya (Rp)</label>
                <input v-model.number="form.cost" class="input" type="number" placeholder="0" />
              </div>
            </div>
            <div>
              <label class="label">Catatan Tambahan</label>
              <input v-model="form.notes" class="input" placeholder="Opsional" />
            </div>

            <div v-if="formError" class="text-red-500 text-xs font-bold">{{ formError }}</div>

            <div class="flex gap-3">
              <button @click="showCreateModal = false" class="flex-1 btn-secondary text-sm py-2.5">Batal</button>
              <button @click="submitCreate" :disabled="submitting" class="flex-1 btn-primary text-sm py-2.5">
                {{ submitting ? 'Menyimpan...' : 'Simpan' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Finish Modal (for entering final cost and notes) -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showFinishModal" class="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4" @click.self="showFinishModal = false">
          <div class="bg-white rounded-3xl shadow-2xl w-full max-w-md p-6 space-y-5">
            <h2 class="text-lg font-black text-gray-900">Selesaikan Perbaikan</h2>
            <p class="text-sm text-gray-500">{{ finishTarget?.item_name }} ({{ finishTarget?.item_code }})</p>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="label">Vendor / Tukang</label>
                <input v-model="finishForm.vendor" class="input" />
              </div>
              <div>
                <label class="label">Total Biaya (Rp)</label>
                <input v-model.number="finishForm.cost" class="input" type="number" />
              </div>
            </div>
            <div>
              <label class="label">Catatan Hasil Perbaikan</label>
              <textarea v-model="finishForm.notes" class="input" rows="2" placeholder="Komponen diganti, dll..."></textarea>
            </div>

            <div class="flex gap-3">
              <button @click="showFinishModal = false" class="flex-1 btn-secondary text-sm py-2.5">Batal</button>
              <button @click="submitFinish" :disabled="submitting" class="flex-1 btn-primary text-sm py-2.5">
                {{ submitting ? 'Menyimpan...' : 'Tandai Selesai' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../utils/api'

const logs = ref([])
const loading = ref(false)
const filterStatus = ref('')
const showCreateModal = ref(false)
const showFinishModal = ref(false)
const submitting = ref(false)
const formError = ref('')
const finishTarget = ref(null)

const form = ref({ item_code: '', issue_description: '', cost: 0, vendor: '', notes: '' })
const finishForm = ref({ cost: 0, vendor: '', notes: '' })

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

async function fetchLogs() {
  loading.value = true
  try {
    const params = filterStatus.value ? { status: filterStatus.value } : {}
    const { data } = await api.get('/maintenance', { params })
    if (data.success) logs.value = data.data
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

async function submitCreate() {
  if (!form.value.item_code || !form.value.issue_description) {
    formError.value = 'Kode barang dan deskripsi kerusakan wajib diisi.'
    return
  }
  formError.value = ''
  submitting.value = true
  try {
    const { data } = await api.post('/maintenance', form.value)
    if (data.success) {
      showCreateModal.value = false
      form.value = { item_code: '', issue_description: '', cost: 0, vendor: '', notes: '' }
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
