<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Dashboard</h1>
          <p class="text-primary-100/70 text-sm font-medium">Ringkasan status inventaris hari ini</p>
        </div>
        
        <div class="flex items-center gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10">
          <button @click="dashboardStore.fetchSummary" 
                  class="bg-white/20 hover:bg-white/30 text-white px-5 py-2.5 rounded-xl text-xs font-black transition-all flex items-center gap-2 active:scale-95">
            <svg class="w-4 h-4" :class="{'animate-spin': dashboardStore.loading}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            SEGARKAN DATA
          </button>
        </div>
      </div>
    </div>

    <!-- Error State -->
    <div v-if="dashboardStore.error" class="p-4 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 rounded-2xl border border-red-100 dark:border-red-800 text-sm flex items-center gap-3 animate-shake">
      <svg class="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
      </svg>
      <span class="font-bold">{{ dashboardStore.error }}</span>
    </div>

    <!-- Stats Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- Total Items -->
      <div class="group relative overflow-hidden bg-white dark:bg-gray-800 p-6 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-xl transition-all duration-500">
        <div class="absolute top-0 right-0 w-24 h-24 -mt-8 -mr-8 bg-indigo-500 rounded-full blur-2xl opacity-10 group-hover:opacity-20 transition-opacity"></div>
        <div class="relative flex items-center gap-4">
          <div class="w-14 h-14 rounded-2xl flex items-center justify-center text-white shadow-lg bg-gradient-to-br from-indigo-500 to-indigo-700 shadow-indigo-200 dark:shadow-none">
            <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" /></svg>
          </div>
          <div>
            <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest leading-none mb-1">Total Barang</p>
            <p class="text-3xl font-black text-gray-900 dark:text-white leading-none">{{ dashboardStore.summary?.total_items || 0 }}</p>
          </div>
        </div>
      </div>

      <!-- Borrowed -->
      <div class="group relative overflow-hidden bg-white dark:bg-gray-800 p-6 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-xl transition-all duration-500">
        <div class="absolute top-0 right-0 w-24 h-24 -mt-8 -mr-8 bg-blue-500 rounded-full blur-2xl opacity-10 group-hover:opacity-20 transition-opacity"></div>
        <div class="relative flex items-center gap-4">
          <div class="w-14 h-14 rounded-2xl flex items-center justify-center text-white shadow-lg bg-gradient-to-br from-blue-500 to-blue-700 shadow-blue-200 dark:shadow-none">
            <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
          </div>
          <div>
            <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest leading-none mb-1">Dipinjam</p>
            <p class="text-3xl font-black text-gray-900 dark:text-white leading-none">{{ dashboardStore.summary?.items_by_status?.BORROWED || 0 }}</p>
          </div>
        </div>
        <div v-if="dashboardStore.summary?.overdue_count" class="absolute bottom-3 left-1/2 -translate-x-1/2 px-3 py-1 rounded-full bg-red-50 dark:bg-red-900/30 text-red-600 dark:text-red-400 text-[9px] font-black uppercase tracking-tighter">
          ⚠️ {{ dashboardStore.summary.overdue_count }} Terlambat
        </div>
      </div>

      <!-- Maintenance -->
      <div class="group relative overflow-hidden bg-white dark:bg-gray-800 p-6 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-xl transition-all duration-500">
        <div class="absolute top-0 right-0 w-24 h-24 -mt-8 -mr-8 bg-amber-500 rounded-full blur-2xl opacity-10 group-hover:opacity-20 transition-opacity"></div>
        <div class="relative flex items-center gap-4">
          <div class="w-14 h-14 rounded-2xl flex items-center justify-center text-white shadow-lg bg-gradient-to-br from-amber-500 to-amber-700 shadow-amber-200 dark:shadow-none">
            <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
          </div>
          <div>
            <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest leading-none mb-1">Perbaikan</p>
            <p class="text-3xl font-black text-gray-900 dark:text-white leading-none">{{ (dashboardStore.summary?.items_by_condition?.DAMAGED || 0) + (dashboardStore.summary?.items_by_condition?.IN_REPAIR || 0) }}</p>
          </div>
        </div>
      </div>

      <!-- Lost -->
      <div class="group relative overflow-hidden bg-white dark:bg-gray-800 p-6 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-xl transition-all duration-500">
        <div class="absolute top-0 right-0 w-24 h-24 -mt-8 -mr-8 bg-rose-500 rounded-full blur-2xl opacity-10 group-hover:opacity-20 transition-opacity"></div>
        <div class="relative flex items-center gap-4">
          <div class="w-14 h-14 rounded-2xl flex items-center justify-center text-white shadow-lg bg-gradient-to-br from-rose-500 to-rose-700 shadow-rose-200 dark:shadow-none">
            <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3.L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>
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
        <router-link to="/items" class="group bg-white dark:bg-gray-800 p-8 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-2xl hover:border-indigo-500 dark:hover:border-indigo-500 transition-all duration-500 text-center">
          <div class="w-20 h-20 mx-auto bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 rounded-3xl flex items-center justify-center mb-6 group-hover:scale-110 group-hover:bg-indigo-600 group-hover:text-white transition-all duration-500 shadow-inner">
            <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/></svg>
          </div>
          <span class="block text-lg font-black text-gray-900 dark:text-white tracking-tight">Tambah Barang</span>
          <p class="text-xs text-gray-400 mt-2 font-medium">Input aset baru ke sistem</p>
        </router-link>
        
        <router-link to="/borrow" class="group bg-white dark:bg-gray-800 p-8 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-2xl hover:border-blue-500 dark:hover:border-blue-500 transition-all duration-500 text-center">
          <div class="w-20 h-20 mx-auto bg-blue-50 dark:bg-blue-900/30 text-blue-600 rounded-3xl flex items-center justify-center mb-6 group-hover:scale-110 group-hover:bg-blue-600 group-hover:text-white transition-all duration-500 shadow-inner">
            <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4"/></svg>
          </div>
          <span class="block text-lg font-black text-gray-900 dark:text-white tracking-tight">Pinjam Barang</span>
          <p class="text-xs text-gray-400 mt-2 font-medium">Proses peminjaman inventaris</p>
        </router-link>

        <router-link to="/return" class="group bg-white dark:bg-gray-800 p-8 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-2xl hover:border-teal-500 dark:hover:border-teal-500 transition-all duration-500 text-center">
          <div class="w-20 h-20 mx-auto bg-teal-50 dark:bg-teal-900/30 text-teal-600 rounded-3xl flex items-center justify-center mb-6 group-hover:scale-110 group-hover:bg-teal-600 group-hover:text-white transition-all duration-500 shadow-inner">
            <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/></svg>
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
