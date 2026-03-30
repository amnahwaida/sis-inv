<template>
  <div class="animate-fade-in space-y-6 max-w-6xl mx-auto">
    <!-- Header & Search -->
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-gray-900 tracking-tight">Cek Riwayat Siswa</h1>
        <p class="text-sm text-gray-500 mt-1">Lacak barang yang sedang atau pernah dipinjam oleh siswa berdasarkan NIS.</p>
      </div>
      
      <div class="flex items-center gap-2">
        <div class="relative flex-1 lg:w-64">
          <span class="absolute inset-y-0 left-0 pl-3 flex items-center text-gray-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
          </span>
          <input 
            v-model="nis" 
            type="text" 
            placeholder="Cari Nama / NIS..." 
            class="pl-10 input-field rounded-xl border-gray-200 focus:ring-primary-500"
            @input="debouncedSearch"
            @keyup.enter="handleSearch"
          >
          
          <!-- Auto-suggest Dropdown -->
          <div v-if="searchResults.length > 0" class="absolute z-50 w-full mt-2 bg-white rounded-xl shadow-xl border border-gray-100 overflow-hidden">
            <div 
              v-for="s in searchResults" 
              :key="s.nis"
              @click="selectStudent(s)"
              class="px-4 py-3 hover:bg-primary-50 cursor-pointer transition-colors border-b border-gray-50 last:border-0"
            >
              <div class="text-sm font-bold text-gray-900">{{ s.full_name }}</div>
              <div class="text-[10px] text-gray-400 font-mono">NIS: {{ s.nis }} • {{ s.class }}</div>
            </div>
          </div>
        </div>
        <button 
          @click="handleSearch" 
          :disabled="loading || !nis"
          class="px-6 py-2.5 bg-primary-600 hover:bg-primary-700 disabled:opacity-50 text-white font-bold rounded-xl shadow-lg shadow-primary-100 transition-all flex items-center gap-2"
        >
          <svg v-if="loading" class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
          {{ loading ? 'Mencari...' : 'Cek' }}
        </button>
      </div>
    </div>

    <!-- Error/No Result State -->
    <div v-if="error" class="p-4 bg-red-50 border border-red-100 text-red-700 rounded-xl text-sm flex items-center gap-3">
      <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/></svg>
      {{ error }}
    </div>

    <!-- Results -->
    <div v-if="hasSearched && !loading" class="space-y-6">
      <!-- Student Info Summary Header -->
      <div v-if="history.length > 0" class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm flex items-center gap-6">
        <div class="w-16 h-16 rounded-2xl bg-primary-50 text-primary-600 flex items-center justify-center text-2xl font-bold">
          {{ history[0].student_name?.[0] || 'S' }}
        </div>
        <div>
          <h2 class="text-xl font-bold text-gray-900">{{ history[0].student_name }}</h2>
          <div class="flex items-center gap-3 mt-1">
            <span class="px-2.5 py-0.5 bg-gray-100 text-gray-600 rounded-md text-xs font-bold font-mono">NIS: {{ nis }}</span>
            <span class="px-2.5 py-0.5 bg-indigo-50 text-indigo-600 rounded-md text-xs font-bold">{{ history[0].student_class }}</span>
            <span class="text-xs text-gray-400">• Ditemukan {{ history.length }} transaksi</span>
          </div>
        </div>
      </div>

      <!-- History Table -->
      <div class="card p-0 overflow-hidden border-gray-100 shadow-sm">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50/50 border-b border-gray-100">
              <tr>
                <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Barang & Guru</th>
                <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Waktu Pinjam</th>
                <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Status</th>
                <th class="text-left text-[10px] font-black text-gray-400 uppercase tracking-widest px-6 py-4">Tujuan</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-50">
              <tr v-for="trx in history" :key="trx.id" class="hover:bg-gray-50/30 transition-colors">
                <td class="px-6 py-4">
                  <div class="font-bold text-gray-900">{{ trx.item_name }}</div>
                  <div class="text-xs text-gray-400 flex items-center gap-1.5 mt-0.5">
                    <span class="font-mono">{{ trx.item_code }}</span>
                    <span>•</span>
                    <span>Proses by: {{ trx.teacher_name }}</span>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <div class="text-sm font-medium text-gray-700">{{ formatDate(trx.borrowed_at) }}</div>
                  <div v-if="trx.returned_at" class="text-[10px] text-green-500 font-bold flex items-center gap-1 mt-0.5">
                    <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                    Dikembalikan: {{ formatDate(trx.returned_at) }}
                  </div>
                  <div v-else class="text-[10px] text-indigo-500 font-bold flex items-center gap-1 mt-0.5">
                    <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                    Jatuh Tempo: {{ formatDate(trx.due_date) }}
                  </div>
                </td>
                <td class="px-6 py-4">
                  <span class="px-3 py-1 rounded-lg text-[10px] font-black uppercase tracking-wider" :class="{
                    'bg-indigo-100 text-indigo-700': trx.status === 'ACTIVE',
                    'bg-green-100 text-green-700': trx.status === 'RETURNED',
                    'bg-red-100 text-red-700': trx.status === 'OVERDUE'
                  }">
                    {{ trx.status === 'ACTIVE' ? 'SedaNg Dibawa' : (trx.status === 'RETURNED' ? 'Kembali' : 'Terlambat') }}
                  </span>
                </td>
                <td class="px-6 py-4">
                  <p class="text-sm text-gray-500 max-w-xs truncate" :title="trx.purpose">{{ trx.purpose || '-' }}</p>
                </td>
              </tr>
              <tr v-if="history.length === 0" class="text-center">
                <td colspan="4" class="px-6 py-16 text-gray-400">
                  <div class="text-4xl mb-3">🔍</div>
                  <p class="text-base font-bold text-gray-900">Tidak ada riwayat untuk NIS ini</p>
                  <p class="text-sm">Pastikan NIS yang Anda masukkan benar.</p>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Empty Initial State -->
    <div v-if="!hasSearched" class="bg-white rounded-3xl border-2 border-dashed border-gray-200 p-12 text-center">
      <div class="w-20 h-20 bg-primary-50 text-primary-600 rounded-full flex items-center justify-center mx-auto mb-6 text-3xl">
        🎓
      </div>
      <h3 class="text-xl font-bold text-gray-900">Siap Melacak Riwayat Siswa?</h3>
      <p class="text-gray-500 mt-2 max-w-sm mx-auto">Masukkan Nomor Induk Siswa (NIS) pada kolom di atas untuk melihat barang apa saja yang pernah mereka pinjam.</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import api from '../utils/api'

const nis = ref('')
const loading = ref(false)
const history = ref([])
const hasSearched = ref(false)
const error = ref('')
const searchResults = ref([])

let searchTimeout
const debouncedSearch = () => {
  clearTimeout(searchTimeout)
  if (!nis.value) {
    searchResults.value = []
    return
  }
  searchTimeout = setTimeout(async () => {
    try {
      const { data } = await api.get(`/students/search?q=${nis.value}`)
      if (data.success) {
        searchResults.value = data.data
      }
    } catch {
      searchResults.value = []
    }
  }, 300)
}

const selectStudent = (student) => {
  nis.value = student.nis
  searchResults.value = []
  handleSearch()
}

const handleSearch = async () => {
  if (!nis.value) return
  
  searchResults.value = []
  loading.value = true
  error.value = ''
  hasSearched.value = true
  
  try {
    const { data } = await api.get(`/transactions/student/${nis.value}`)
    if (data.success) {
      history.value = data.data
    } else {
      error.value = data.error || 'Gagal mengambil data'
    }
  } catch (err) {
    console.error(err)
    error.value = 'Gagal menghubungi server. Pastikan NIS benar.'
  } finally {
    loading.value = false
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}
</script>
