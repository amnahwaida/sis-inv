<template>
  <div class="animate-fade-in space-y-6 max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-gray-900 tracking-tight">Log Audit Sistem</h1>
        <p class="text-sm text-gray-500 mt-1">Rekam jejak aktivitas krusial dalam sistem (Terbatas untuk Admin).</p>
      </div>
      <div class="flex items-center gap-2">
        <button @click="fetchLogs" class="p-2 text-primary-600 hover:bg-primary-50 rounded-xl transition-all" title="Refresh">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/></svg>
        </button>
      </div>
    </div>

    <!-- Logs Table -->
    <div class="card p-0 overflow-hidden border-gray-100 shadow-sm bg-white rounded-3xl">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50/50 border-b border-gray-100">
            <tr>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4 w-12">No</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Waktu</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Pelaku</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Aksi</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Deskripsi</th>
              <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Alamat IP</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-if="loading" class="text-center">
              <td colspan="6" class="px-6 py-12 text-gray-400 italic">Memuat data log...</td>
            </tr>
            <tr v-else-if="logs.length === 0" class="text-center">
              <td colspan="6" class="px-6 py-12 text-gray-400">
                <p class="font-bold text-gray-900 italic">Belum ada riwayat audit yang tertangkap.</p>
              </td>
            </tr>
            <tr v-for="(log, idx) in logs" :key="log.id" class="hover:bg-gray-50/50 transition-colors text-xs">
              <td class="px-6 py-4 font-mono text-gray-400">{{ idx + 1 }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="font-bold text-gray-900">{{ formatDate(log.created_at) }}</div>
                <div class="text-[10px] text-gray-400">{{ formatTime(log.created_at) }}</div>
              </td>
              <td class="px-6 py-4">
                <div class="font-bold text-gray-800">{{ log.user_name || 'System' }}</div>
              </td>
              <td class="px-6 py-4">
                <span class="px-2 py-0.5 rounded-md font-black uppercase text-[9px] tracking-widest border border-gray-100 shadow-sm" :class="getActionBadge(log.action)">
                  {{ log.action }}
                </span>
              </td>
              <td class="px-6 py-4 text-gray-600 font-medium max-w-xs truncate" :title="log.description">
                {{ log.description }}
              </td>
              <td class="px-6 py-4 font-mono text-[10px] text-gray-500">
                {{ log.ip_address || '127.0.0.1' }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../utils/api'

const logs = ref([])
const loading = ref(false)

async function fetchLogs() {
  loading.value = true
  try {
    const { data } = await api.get('/reports/audit')
    if (data.success) {
      logs.value = data.data
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Intl.DateTimeFormat('id-ID', {
    day: '2-digit', month: 'short', year: 'numeric'
  }).format(new Date(dateStr))
}

const formatTime = (dateStr) => {
  if (!dateStr) return '-'
  return new Intl.DateTimeFormat('id-ID', {
    hour: '2-digit', minute: '2-digit', second: '2-digit'
  }).format(new Date(dateStr))
}

const getActionBadge = (action) => {
  if (action.includes('CREATE')) return 'bg-green-50 text-green-700 border-green-200'
  if (action.includes('DELETE')) return 'bg-red-50 text-red-700 border-red-200'
  if (action.includes('UPDATE')) return 'bg-blue-50 text-blue-700 border-blue-200'
  return 'bg-gray-50 text-gray-700 border-gray-200'
}

onMounted(fetchLogs)
</script>
