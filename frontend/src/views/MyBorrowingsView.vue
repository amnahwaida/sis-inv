<template>
  <div class="animate-fade-in space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">Peminjaman Saya</h1>
      <router-link to="/borrow" class="btn-primary flex items-center gap-2">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Pinjam Baru
      </router-link>
    </div>

    <!-- Borrowings Table -->
    <div class="card p-0 overflow-hidden relative min-h-[300px]">
      <!-- Loading Overlay -->
      <div v-if="transactionStore.loading" class="absolute inset-0 bg-white/80 z-10 flex flex-col items-center justify-center">
        <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-indigo-600"></div>
        <span class="mt-2 text-sm text-gray-600">Memuat data...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Barang</th>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Tanggal Pinjam</th>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Jatuh Tempo</th>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Status</th>
              <th class="text-right text-xs font-semibold text-gray-600 uppercase px-6 py-3">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <template v-if="transactionStore.myBorrowings.length > 0">
              <tr v-for="item in transactionStore.myBorrowings" :key="item.id" class="hover:bg-gray-50 transition-colors">
                <td class="px-6 py-4">
                  <div class="font-medium text-gray-900">{{ item.item_name }}</div>
                  <div class="text-sm text-gray-500 font-mono">{{ item.item_code }}</div>
                </td>
                <td class="px-6 py-4">
                  <div class="text-sm text-gray-900">{{ formatDate(item.borrow_date) }}</div>
                </td>
                <td class="px-6 py-4">
                  <div class="text-sm" :class="isOverdue(item.expected_return_date) ? 'text-red-600 font-bold' : 'text-gray-900'">
                    {{ formatDate(item.expected_return_date) }}
                  </div>
                  <div v-if="isOverdue(item.expected_return_date)" class="text-xs text-red-500 mt-0.5">TERLAMBAT</div>
                </td>
                <td class="px-6 py-4">
                  <span class="px-2.5 py-1 rounded-full text-xs font-medium" :class="isOverdue(item.expected_return_date) ? 'bg-red-100 text-red-800' : 'bg-blue-100 text-blue-800'">
                    {{ isOverdue(item.expected_return_date) ? 'Overdue' : 'Active' }}
                  </span>
                </td>
                <td class="px-6 py-4 text-right">
                  <router-link to="/return" class="text-indigo-600 hover:text-indigo-900 text-sm font-medium">Kembalikan</router-link>
                </td>
              </tr>
            </template>
            <tr v-else-if="!transactionStore.loading" class="text-center">
              <td colspan="5" class="px-6 py-12 text-gray-500">
                <p class="text-3xl mb-2">🎓</p>
                <p class="text-base font-medium text-gray-900">Tidak ada peminjaman aktif</p>
                <p class="text-sm text-gray-500 mt-1">Anda belum meminjam barang apapun saat ini.</p>
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
