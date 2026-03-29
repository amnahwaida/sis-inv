<template>
  <div class="animate-fade-in space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
      <button @click="dashboardStore.fetchSummary" class="text-sm font-medium text-indigo-600 flex items-center gap-1 hover:text-indigo-800 transition-colors">
        <svg class="w-4 h-4" :class="{'animate-spin': dashboardStore.loading}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        Refresh
      </button>
    </div>

    <!-- Error State -->
    <div v-if="dashboardStore.error" class="p-4 bg-red-50 text-red-700 rounded-lg text-sm flex items-center gap-2">
      <svg class="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
      </svg>
      {{ dashboardStore.error }}
    </div>

    <!-- Summary Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      
      <!-- Total Items -->
      <div class="card bg-gradient-to-br from-indigo-500 to-primary-600 text-white shadow-indigo-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-indigo-100 text-sm font-medium">Total Barang</p>
            <h3 class="text-3xl font-bold mt-1">{{ dashboardStore.summary?.total_items || 0 }}</h3>
          </div>
          <div class="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Out / Borrowed -->
      <div class="card border-l-4 border-l-blue-500">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm font-medium">Sedang Dipinjam</p>
            <h3 class="text-3xl font-bold text-gray-900 mt-1">{{ dashboardStore.summary?.items_by_status?.BORROWED || 0 }}</h3>
          </div>
          <div class="w-12 h-12 bg-blue-50 text-blue-500 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
        <div v-if="dashboardStore.summary?.overdue_count" class="mt-4 text-xs font-medium text-red-600 bg-red-50 py-1 px-2 rounded-md inline-block">
          ⚠️ {{ dashboardStore.summary.overdue_count }} peminjaman terlambat
        </div>
      </div>

      <!-- Maintenance / Repair -->
      <div class="card border-l-4 border-l-yellow-500">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm font-medium">Perbaikan / Rusak</p>
            <h3 class="text-3xl font-bold text-gray-900 mt-1">
              {{ (dashboardStore.summary?.items_by_condition?.DAMAGED || 0) + (dashboardStore.summary?.items_by_condition?.IN_REPAIR || 0) }}
            </h3>
          </div>
          <div class="w-12 h-12 bg-yellow-50 text-yellow-500 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Lost -->
      <div class="card border-l-4 border-l-red-500">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm font-medium">Barang Hilang</p>
            <h3 class="text-3xl font-bold text-gray-900 mt-1">{{ dashboardStore.summary?.items_by_condition?.LOST || 0 }}</h3>
          </div>
          <div class="w-12 h-12 bg-red-50 text-red-500 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3.L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
        </div>
      </div>

    </div>

    <!-- Quick Actions -->
    <div class="mt-8">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">Aksi Cepat</h2>
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <router-link to="/items" class="bg-white p-4 rounded-xl border border-gray-100 shadow-sm hover:shadow-md hover:border-indigo-100 transition-all text-center group">
          <div class="w-10 h-10 mx-auto bg-indigo-50 text-indigo-600 rounded-lg flex items-center justify-center mb-3 group-hover:scale-110 transition-transform">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/></svg>
          </div>
          <span class="text-sm font-medium text-gray-700">Tambah Barang</span>
        </router-link>
        
        <router-link to="/borrow" class="bg-white p-4 rounded-xl border border-gray-100 shadow-sm hover:shadow-md hover:border-blue-100 transition-all text-center group">
          <div class="w-10 h-10 mx-auto bg-blue-50 text-blue-600 rounded-lg flex items-center justify-center mb-3 group-hover:scale-110 transition-transform">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
          </div>
          <span class="text-sm font-medium text-gray-700">Pinjam Barang</span>
        </router-link>

        <router-link to="/return" class="bg-white p-4 rounded-xl border border-gray-100 shadow-sm hover:shadow-md hover:border-teal-100 transition-all text-center group">
          <div class="w-10 h-10 mx-auto bg-teal-50 text-teal-600 rounded-lg flex items-center justify-center mb-3 group-hover:scale-110 transition-transform">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/></svg>
          </div>
          <span class="text-sm font-medium text-gray-700">Kembalikan Barang</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useDashboardStore } from '../stores/dashboard'

const dashboardStore = useDashboardStore()

onMounted(() => {
  dashboardStore.fetchSummary()
})
</script>
