<template>
  <div class="animate-fade-in space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">Daftar Barang</h1>
      <button v-if="authStore.isAdmin" class="btn-primary flex items-center gap-2">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Tambah Barang
      </button>
    </div>

    <!-- Search & Filter -->
    <div class="card">
      <div class="flex flex-col md:flex-row gap-3">
        <div class="flex-1">
          <input type="text" class="input-field" placeholder="Cari nama / kode barang..." v-model="search" />
        </div>
        <select class="input-field md:w-40" v-model="statusFilter">
          <option value="">Semua Status</option>
          <option value="AVAILABLE">Tersedia</option>
          <option value="BORROWED">Dipinjam</option>
          <option value="MAINTENANCE">Perbaikan</option>
          <option value="LOST">Hilang</option>
        </select>
      </div>
    </div>

    <!-- Items Table -->
    <div class="card p-0 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Kode</th>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Nama Barang</th>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Lokasi</th>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Status</th>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Kondisi</th>
              <th class="text-right text-xs font-semibold text-gray-600 uppercase px-6 py-3">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr class="text-center">
              <td colspan="6" class="px-6 py-12 text-gray-500">
                <p class="text-lg">📦</p>
                <p class="mt-2 text-sm">Belum ada data barang</p>
                <p class="text-xs text-gray-400 mt-1">Fitur CRUD barang akan diimplementasikan di Phase 2</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const search = ref('')
const statusFilter = ref('')
</script>
