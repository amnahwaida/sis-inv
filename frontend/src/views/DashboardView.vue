<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-primary-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Dashboard</h1>
          <p class="text-primary-100/70 text-sm font-medium">Ringkasan status inventaris hari ini</p>
        </div>
        
        <div class="flex items-center gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10">
          <button @click="dashboardStore.fetchSummary" class="bg-white/10 hover:bg-white/20 text-white border border-white/10 px-5 py-2.5 rounded-xl text-xs font-black transition-all flex items-center gap-2 active:scale-95">
            <svg class="w-4 h-4" :class="{'animate-spin': dashboardStore.loading}" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
            </svg>
            SEGARKAN DATA
          </button>
        </div>
      </div>
    </div>

    <!-- Error State -->
    <div v-if="dashboardStore.error" class="p-4 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 rounded-2xl border border-red-100 dark:border-red-800 text-sm flex items-center gap-3 animate-shake">
      <svg class="w-5 h-5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" d="M9.75 9.75l4.5 4.5m0-4.5l-4.5 4.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span class="font-bold">{{ dashboardStore.error }}</span>
    </div>

    <!-- Stats Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- Total Items -->
      <div class="card-premium group">
        <div class="card-decoration bg-primary-500"></div>
        <div class="relative flex items-center gap-4">
          <div class="w-14 h-14 rounded-2xl flex items-center justify-center text-white shadow-lg bg-gradient-to-br from-primary-500 to-primary-700 shadow-primary-200 dark:shadow-none">
            <svg class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-5.25v9" />
            </svg>
          </div>
          <div>
            <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest leading-none mb-1">Total Barang</p>
            <p class="text-3xl font-black text-gray-900 dark:text-white leading-none">{{ dashboardStore.summary?.total_items || 0 }}</p>
          </div>
        </div>
      </div>

      <!-- Borrowed -->
      <div class="card-premium group">
        <div class="card-decoration bg-primary-500"></div>
        <div class="relative flex items-center gap-4">
          <div class="w-14 h-14 rounded-2xl flex items-center justify-center text-white shadow-lg bg-gradient-to-br from-primary-500 to-primary-700 shadow-primary-200 dark:shadow-none">
            <svg class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div>
            <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest leading-none mb-1">Dipinjam</p>
            <p class="text-3xl font-black text-gray-900 dark:text-white leading-none">{{ dashboardStore.summary?.items_by_status?.BORROWED || 0 }}</p>
          </div>
        </div>
        <div v-if="dashboardStore.summary?.overdue_count" class="relative mt-4 inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-red-50 dark:bg-red-900/30 text-red-600 dark:text-red-400 text-[10px] font-black uppercase tracking-widest shadow-sm w-fit">
          <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
          </svg>
          {{ dashboardStore.summary.overdue_count }} Terlambat
        </div>
      </div>

      <!-- Maintenance -->
      <div class="card-premium group">
        <div class="card-decoration bg-amber-500"></div>
        <div class="relative flex items-center gap-4">
          <div class="w-14 h-14 rounded-2xl flex items-center justify-center text-white shadow-lg bg-gradient-to-br from-amber-500 to-amber-700 shadow-amber-200 dark:shadow-none">
            <svg class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21.75 6.75a4.5 4.5 0 01-4.884 4.484c-1.076-.091-2.264.071-2.95.904l-7.152 8.655a2.25 2.25 0 01-3.323.14l-1.182-1.182a2.25 2.25 0 01.14-3.323l8.655-7.152c.833-.686.995-1.874.904-2.95a4.5 4.5 0 014.484-4.884h.01a4.5 4.5 0 014.5 4.5v.01z" />
              <path stroke-linecap="round" stroke-linejoin="round" d="M9.448 10.276a2.25 2.25 0 012.276-2.276m-7.5 7.5a2.25 2.25 0 012.276-2.276" />
            </svg>
          </div>
          <div>
            <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest leading-none mb-1">Perbaikan</p>
            <p class="text-3xl font-black text-gray-900 dark:text-white leading-none">{{ (dashboardStore.summary?.items_by_condition?.DAMAGED || 0) + (dashboardStore.summary?.items_by_condition?.IN_REPAIR || 0) }}</p>
          </div>
        </div>
      </div>

      <!-- Lost -->
      <div class="card-premium group">
        <div class="card-decoration bg-rose-500"></div>
        <div class="relative flex items-center gap-4">
          <div class="w-14 h-14 rounded-2xl flex items-center justify-center text-white shadow-lg bg-gradient-to-br from-rose-500 to-rose-700 shadow-rose-200 dark:shadow-none">
            <svg class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
            </svg>
          </div>
          <div>
            <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest leading-none mb-1">Barang Hilang</p>
            <p class="text-3xl font-black text-gray-900 dark:text-white leading-none">{{ dashboardStore.summary?.items_by_condition?.LOST || 0 }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="space-y-6">
      <div class="flex items-center gap-4">
        <div class="h-px bg-gray-100 dark:bg-gray-700 flex-1"></div>
        <h2 class="text-xs font-black text-gray-400 uppercase tracking-[0.3em]">Aksi Cepat</h2>
        <div class="h-px bg-gray-100 dark:bg-gray-700 flex-1"></div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-3 gap-6">
        <router-link to="/items" class="card-premium group text-center">
          <div class="w-20 h-20 mx-auto bg-primary-50 dark:bg-primary-900/30 text-primary-600 rounded-3xl flex items-center justify-center mb-6 group-hover:scale-110 group-hover:bg-primary-600 group-hover:text-white transition-all duration-500 shadow-inner">
            <svg class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
            </svg>
          </div>
          <span class="block text-lg font-black text-gray-900 dark:text-white tracking-tight">Tambah Barang</span>
          <p class="text-xs text-gray-400 mt-2 font-medium">Input aset baru ke sistem</p>
        </router-link>
        
        <router-link to="/borrow" class="card-premium group text-center !hover:border-primary-500">
          <div class="w-20 h-20 mx-auto bg-primary-50 dark:bg-primary-900/30 text-primary-600 rounded-3xl flex items-center justify-center mb-6 group-hover:scale-110 group-hover:bg-primary-600 group-hover:text-white transition-all duration-500 shadow-inner">
            <svg class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <span class="block text-lg font-black text-gray-900 dark:text-white tracking-tight">Pinjam Barang</span>
          <p class="text-xs text-gray-400 mt-2 font-medium">Proses peminjaman inventaris</p>
        </router-link>
 
        <router-link to="/return" class="card-premium group text-center !hover:border-primary-500">
          <div class="w-20 h-20 mx-auto bg-primary-50 dark:bg-primary-900/30 text-primary-600 rounded-3xl flex items-center justify-center mb-6 group-hover:scale-110 group-hover:bg-primary-600 group-hover:text-white transition-all duration-500 shadow-inner">
            <svg class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
            </svg>
          </div>
          <span class="block text-lg font-black text-gray-900 dark:text-white tracking-tight">Kembali Barang</span>
          <p class="text-xs text-gray-400 mt-2 font-medium">Selesaikan masa pinjam aset</p>
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
