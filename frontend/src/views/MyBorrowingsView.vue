<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Peminjaman Saya</h1>
          <p class="text-primary-100/70 text-sm font-medium">Pantau status barang yang sedang Anda gunakan</p>
        </div>
      </div>
    </div>

    <!-- Active Borrowings Table Card -->
    <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden transition-all duration-300">
      <div class="p-8 border-b border-gray-100 dark:border-gray-700 flex items-center justify-between">
        <h2 class="text-lg font-black text-gray-900 dark:text-white uppercase tracking-tight">Barang Belum Kembali</h2>
        <span class="px-4 py-1 bg-primary-50 dark:bg-primary-900/40 text-primary-600 dark:text-primary-400 rounded-full text-[10px] font-black uppercase tracking-widest shadow-sm">
          {{ loadings.active ? 'Membaca...' : (activeBorrowings.length + ' Aset') }}
        </span>
      </div>

      <div class="overflow-x-auto relative">
        <table class="w-full">
          <thead>
            <tr class="bg-gray-50/50 dark:bg-gray-700/30">
              <th v-for="h in ['Aset', 'Peminjam', 'Tgl Pinjam', 'Due Date', 'Status', 'Aksi']" :key="h"
                  class="text-left py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">{{ h }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50 dark:divide-gray-700">
            <template v-if="activeBorrowings.length > 0">
              <tr v-for="trx in activeBorrowings" :key="trx.id" 
                  @click="$router.push(`/return?code=${trx.item_code}`)"
                  class="group hover:bg-primary-50/50 dark:hover:bg-primary-900/10 transition-all duration-300 cursor-pointer active:bg-primary-100/50 dark:active:bg-primary-900/20">
                <td class="px-8 py-6">
                  <div class="font-black text-gray-900 dark:text-white text-sm tracking-tight leading-none mb-1">{{ trx.item_name }}</div>
                  <div class="text-[10px] text-gray-400 font-mono tracking-tighter">{{ trx.item_code }}</div>
                </td>
                <td class="px-8 py-6">
                  <div class="text-xs font-bold text-gray-700 dark:text-gray-300">{{ trx.borrower_display }}</div>
                </td>
                <td class="px-8 py-6 text-xs font-medium text-gray-500">{{ formatDate(trx.borrow_date) }}</td>
                <td class="px-8 py-6 text-xs font-black" :class="trx.is_overdue ? 'text-red-500 animate-pulse' : 'text-primary-600 dark:text-primary-400'">
                  {{ formatDate(trx.due_date) }}
                </td>
                <td class="px-8 py-6">
                  <div v-if="trx.is_overdue" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-red-50 dark:bg-red-900/30 text-red-600 font-black text-[9px] uppercase tracking-widest">
                    <span class="w-1.5 h-1.5 bg-red-600 rounded-full"></span>
                    TELAT
                  </div>
                  <div v-else class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 font-black text-[9px] uppercase tracking-widest">
                    <span class="w-1.5 h-1.5 bg-emerald-500 rounded-full"></span>
                    AKTIF
                  </div>
                </td>
                <td class="px-8 py-6">
                  <router-link :to="`/return?code=${trx.item_code}`" 
                               @click.stop
                               class="inline-flex items-center gap-2 px-4 py-2 bg-white dark:bg-gray-800 hover:bg-indigo-50 dark:hover:bg-indigo-900/30 text-indigo-600 border border-indigo-100 dark:border-indigo-800 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all shadow-sm hover:shadow active:scale-95">
                    Kembalikan
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M14 5l7 7m0 0l-7 7m7-7H3"/></svg>
                  </router-link>
                </td>
              </tr>
            </template>
            <tr v-else-if="!loadings.active" class="text-center">
              <td colspan="6" class="px-8 py-24 italic text-gray-400 font-bold uppercase tracking-[0.2em] text-xs">Semua barang telah Anda kembalikan</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- History Header -->
    <div class="pt-10 flex items-center gap-4">
      <div class="h-px bg-gray-100 dark:bg-gray-700 flex-1"></div>
      <h2 class="text-xl font-black text-gray-400 dark:text-gray-600 uppercase tracking-widest leading-none">Riwayat Peminjaman</h2>
      <div class="h-px bg-gray-100 dark:bg-gray-700 flex-1"></div>
    </div>

    <!-- History Table Card -->
    <div class="bg-gray-50 dark:bg-gray-900/20 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 overflow-hidden transition-all duration-300">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="bg-gray-100/50 dark:bg-gray-800/50">
              <th v-for="h in ['Barang', 'Tgl Pinjam', 'Tgl Kembali', 'Catatan', 'Kondisi Akhir']" :key="h"
                  class="text-left py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">{{ h }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 dark:divide-gray-800">
            <template v-if="history.length > 0">
              <tr v-for="trx in history" :key="trx.id" class="group hover:bg-white dark:hover:bg-gray-800 transition-all duration-300">
                <td class="px-8 py-6">
                  <div class="font-bold text-gray-700 dark:text-gray-200 text-sm">{{ trx.item_name }}</div>
                  <div class="text-[10px] text-gray-400 font-mono tracking-tighter">{{ trx.item_code }}</div>
                </td>
                <td class="px-8 py-6 text-xs text-gray-500 font-medium">{{ formatDate(trx.borrowed_at) }}</td>
                <td class="px-8 py-6 text-xs text-gray-500 font-medium">{{ formatDate(trx.returned_at) }}</td>
                <td class="px-8 py-6">
                  <p class="text-xs text-gray-400 italic leading-snug max-w-xs truncate" :title="trx.return_notes || '-'">{{ trx.return_notes || '-' }}</p>
                </td>
                <td class="px-8 py-6">
                  <span v-if="trx.return_condition" class="px-3 py-1 rounded-lg text-[9px] font-black uppercase tracking-widest shadow-sm inline-block"
                        :class="getConditionBadgeClass(trx.return_condition)">
                    {{ getConditionLabel(trx.return_condition) }}
                  </span>
                </td>
              </tr>
            </template>
            <tr v-else-if="!loadings.history" class="text-center">
              <td colspan="5" class="px-8 py-16 text-gray-400 italic text-sm">Belum ada riwayat tercatat</td>
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

const activeBorrowings = ref([])
const history = ref([])
const loadings = ref({ active: false, history: false })

const fetchActive = async () => {
  loadings.value.active = true
  try {
    const { data } = await api.get('/transactions/my')
    if (data.success) { activeBorrowings.value = data.data }
  } catch (e) { console.error('Gagal fetch active borrowings:', e) } 
  finally { loadings.value.active = false }
}

const fetchHistory = async () => {
  loadings.value.history = true
  try {
    const { data } = await api.get('/transactions/my/history')
    if (data.success) { history.value = data.data }
  } catch (e) { console.error('Gagal fetch history:', e) } 
  finally { loadings.value.history = false }
}

const formatDate = (d) => {
  if (!d) return '-'
  return new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }).format(new Date(d))
}

const getConditionBadgeClass = (c) => ({
  'GOOD': 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600',
  'DAMAGED': 'bg-rose-50 dark:bg-rose-900/30 text-rose-600',
  'IN_REPAIR': 'bg-amber-50 dark:bg-amber-900/30 text-amber-600'
}[c] || 'bg-gray-50 text-gray-600')

const getConditionLabel = (c) => ({ 'GOOD': 'Baik', 'DAMAGED': 'Rusak', 'IN_REPAIR': 'Perbaikan' }[c] || c)

onMounted(() => { fetchActive(); fetchHistory() })
</script>
