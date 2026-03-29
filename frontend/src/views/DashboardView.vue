<template>
  <div class="animate-fade-in space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p class="text-gray-500 text-sm mt-1">Selamat datang, {{ authStore.user?.full_name }} 👋</p>
      </div>
    </div>

    <!-- Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div v-for="card in summaryCards" :key="card.label" class="card hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">{{ card.label }}</p>
            <p class="text-2xl font-bold mt-1" :class="card.color">{{ card.value }}</p>
          </div>
          <div class="w-12 h-12 rounded-xl flex items-center justify-center" :class="card.bg">
            <span class="text-xl">{{ card.icon }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="card" v-if="authStore.isAdmin || authStore.isTeacher">
      <h3 class="text-lg font-semibold text-gray-800 mb-4">Aksi Cepat</h3>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
        <router-link to="/items" class="flex flex-col items-center gap-2 p-4 rounded-xl bg-gray-50 hover:bg-primary-50 hover:text-primary-700 transition-colors text-gray-600 text-sm font-medium">
          <span class="text-2xl">📦</span>
          <span>Lihat Barang</span>
        </router-link>
        <button v-if="authStore.isAdmin" class="flex flex-col items-center gap-2 p-4 rounded-xl bg-gray-50 hover:bg-green-50 hover:text-green-700 transition-colors text-gray-600 text-sm font-medium">
          <span class="text-2xl">➕</span>
          <span>Tambah Barang</span>
        </button>
        <button class="flex flex-col items-center gap-2 p-4 rounded-xl bg-gray-50 hover:bg-blue-50 hover:text-blue-700 transition-colors text-gray-600 text-sm font-medium">
          <span class="text-2xl">📷</span>
          <span>Scan QR</span>
        </button>
        <button v-if="authStore.isTeacher || authStore.isAdmin" class="flex flex-col items-center gap-2 p-4 rounded-xl bg-gray-50 hover:bg-purple-50 hover:text-purple-700 transition-colors text-gray-600 text-sm font-medium">
          <span class="text-2xl">🎓</span>
          <span>Pinjam Siswa</span>
        </button>
      </div>
    </div>

    <!-- Info Banner -->
    <div class="bg-gradient-to-r from-primary-500 to-primary-700 rounded-xl p-6 text-white">
      <h3 class="text-lg font-semibold">🚀 Selamat Datang di SIS-INV</h3>
      <p class="text-primary-100 text-sm mt-2">Sistem Inventaris Sekolah berbasis web. Kelola aset, pantau peminjaman, dan buat laporan dengan mudah.</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()

const summaryCards = ref([
  { label: 'Total Barang', value: '—', icon: '📦', color: 'text-gray-900', bg: 'bg-blue-100' },
  { label: 'Dipinjam', value: '—', icon: '🔄', color: 'text-blue-600', bg: 'bg-blue-50' },
  { label: 'Tersedia', value: '—', icon: '✅', color: 'text-green-600', bg: 'bg-green-50' },
  { label: 'Rusak/Hilang', value: '—', icon: '⚠️', color: 'text-red-600', bg: 'bg-red-50' },
])
</script>
