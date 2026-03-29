<template>
  <div class="animate-fade-in space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Peminjaman Saya</h1>
        <p class="text-sm text-gray-500 mt-1 uppercase tracking-wider font-bold text-[10px]">Daftar barang yang Anda kelola / pinjam</p>
      </div>
      <router-link to="/borrow" class="btn-primary flex items-center gap-2 px-6 py-2.5 rounded-xl shadow-lg shadow-primary-100">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.1" d="M12 4v16m8-8H4" />
        </svg>
        Pinjam Baru
      </router-link>
    </div>

    <!-- Borrowings Table -->
    <div class="card p-0 overflow-hidden relative shadow-sm border border-gray-100 min-h-[400px]">
      <!-- Loading Overlay -->
      <div v-if="transactionStore.loading" class="absolute inset-0 bg-white/60 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center">
        <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-primary-600"></div>
        <span class="mt-4 text-xs font-bold text-gray-500 uppercase tracking-widest">Sinkronisasi Data...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50/80 border-b border-gray-100">
            <tr>
              <th class="text-left text-xs font-bold text-gray-500 uppercase tracking-widest px-6 py-4">Barang & Kode</th>
              <th class="text-left text-xs font-bold text-gray-500 uppercase tracking-widest px-6 py-4">Peminjam</th>
              <th class="text-left text-xs font-bold text-gray-500 uppercase tracking-widest px-6 py-4">Waktu Pinjam</th>
              <th class="text-left text-xs font-bold text-gray-500 uppercase tracking-widest px-6 py-4">Jatuh Tempo</th>
              <th class="text-right text-xs font-bold text-gray-500 uppercase tracking-widest px-6 py-4">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <template v-if="transactionStore.myBorrowings.length > 0">
              <tr v-for="item in transactionStore.myBorrowings" :key="item.id" class="hover:bg-gray-50/50 transition-colors group">
                <td class="px-6 py-4">
                  <div class="font-bold text-gray-900 group-hover:text-primary-700 transition-colors">{{ item.item_name }}</div>
                  <div class="text-[10px] text-gray-400 font-mono mt-0.5 tracking-tighter">{{ item.item_code }}</div>
                </td>
                <td class="px-6 py-4">
                  <div class="flex items-center gap-2">
                    <span 
                      class="px-2 py-0.5 rounded text-[10px] font-black uppercase tracking-tighter"
                      :class="item.borrower_type === 'STAFF' ? 'bg-indigo-50 text-indigo-700 border border-indigo-100' : 'bg-amber-50 text-amber-700 border border-amber-100'"
                    >
                      {{ item.borrower_type }}
                    </span>
                    <span class="text-sm font-semibold text-gray-700">{{ item.borrower_display }}</span>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <div class="text-sm text-gray-600 font-medium">{{ formatDate(item.borrow_date) }}</div>
                </td>
                <td class="px-6 py-4">
                  <div class="flex flex-col">
                    <span class="text-sm font-bold" :class="isOverdue(item.due_date) ? 'text-red-600' : 'text-gray-700'">
                      {{ formatDate(item.due_date) }}
                    </span>
                    <span v-if="isOverdue(item.due_date)" class="text-[9px] font-black text-white bg-red-500 px-1.5 py-0.5 rounded w-max mt-1 uppercase">Terlambat</span>
                  </div>
                </td>
                <td class="px-6 py-4 text-right">
                  <router-link 
                    :to="{ path: '/return', query: { code: item.item_code } }" 
                    class="inline-flex items-center gap-1.5 bg-white border border-gray-200 text-gray-700 hover:border-primary-500 hover:text-primary-600 px-4 py-2 rounded-xl text-xs font-bold transition-all shadow-sm hover:shadow-md active:scale-95"
                  >
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M11 15l-3-3m0 0l3-3m-3 3h8M3 12a9 9 0 1118 0 8.959 8.959 0 01-4.5 7.723" />
                    </svg>
                    Kembalikan
                  </router-link>
                </td>
              </tr>
            </template>
            <tr v-else-if="!transactionStore.loading" class="text-center">
              <td colspan="5" class="px-6 py-24 text-gray-500">
                <div class="flex flex-col items-center">
                  <div class="w-16 h-16 bg-gray-50 rounded-full flex items-center justify-center mb-4">
                    <svg class="w-8 h-8 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
                    </svg>
                  </div>
                  <p class="text-base font-bold text-gray-900 tracking-tight">Tidak Ada Pinjaman Aktif</p>
                  <p class="text-xs text-gray-400 mt-1">Anda tidak sedang meminjam atau mengelola barang apapun.</p>
                  <router-link to="/borrow" class="mt-6 text-primary-600 font-bold text-sm hover:underline flex items-center gap-1">
                    Cari barang untuk dipinjam
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
                    </svg>
                  </router-link>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useTransactionStore } from '../stores/transaction'

const transactionStore = useTransactionStore()

onMounted(() => {
  transactionStore.fetchMyBorrowings()
})

const formatDate = (isoString) => {
  if (!isoString) return '-'
  const date = new Date(isoString)
  return new Intl.DateTimeFormat('id-ID', { 
    day: 'numeric', month: 'short', year: 'numeric' 
  }).format(date)
}

const isOverdue = (dateString) => {
  if (!dateString) return false
  const due = new Date(dateString)
  const now = new Date()
  return now > due
}
</script>
