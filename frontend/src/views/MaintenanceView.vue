<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Log Perbaikan</h1>
          <p class="text-primary-100/70 text-sm font-medium">Manajemen pemeliharaan & perbaikan aset sekolah</p>
        </div>
        
        <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-2 sm:gap-3 backdrop-blur-md bg-white/10 p-2 md:p-1.5 rounded-2xl border border-white/10 w-full md:w-auto mt-4 md:mt-0">
          <button @click="exportExcel" class="bg-white/10 hover:bg-white/20 text-white border border-white/10 px-4 py-2.5 rounded-xl text-[10px] sm:text-xs font-black transition-all flex items-center gap-2 shadow-sm active:scale-95 whitespace-nowrap">
            <svg class="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2-2z"/></svg>
            EXCEL
          </button>
          <button @click="showCreateModal = true" class="btn-premium-primary">
            <svg class="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 4v16m8-8H4"/></svg>
            LAPOR KERUSAKAN
          </button>
        </div>
      </div>
    </div>

    <!-- Stats Cards Grid -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-6">
      <div v-for="(count, key) in statsMap" :key="key" 
           class="group relative overflow-hidden bg-white dark:bg-gray-800 p-6 rounded-[2rem] border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-xl transition-all duration-500">
        <div :class="['absolute top-0 right-0 w-24 h-24 -mt-8 -mr-8 rounded-full blur-2xl opacity-10 group-hover:opacity-20 transition-opacity', count.color]"></div>
        <div class="relative flex items-center gap-4">
          <div :class="['w-12 h-12 rounded-2xl flex items-center justify-center text-white shadow-lg', count.gradient]">
            <component :is="count.icon" class="w-6 h-6" />
          </div>
          <div>
            <p class="text-3xl font-black text-gray-900 dark:text-white leading-none">{{ count.value }}</p>
            <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest mt-1">{{ count.label }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Filters Row -->
    <div class="flex flex-col sm:flex-row items-center gap-3 w-full">
      <!-- Search Input -->
      <div class="relative w-full sm:w-80">
        <input type="text" v-model="searchQuery" placeholder="Cari barang, vendor, keluhan..." 
               class="input-field pl-10 h-11 rounded-2xl text-sm w-full" />
        <svg class="w-4 h-4 absolute left-3.5 top-3.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
      </div>

      <!-- Status Filter -->
      <div class="flex items-center gap-2 bg-white dark:bg-gray-800 p-2 px-4 rounded-2xl border border-gray-100 dark:border-gray-700 shadow-sm w-full sm:w-auto">
        <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest whitespace-nowrap">Status:</label>
        <select v-model="filterStatus" class="bg-transparent border-none focus:ring-0 text-sm font-black text-gray-900 dark:text-white py-0 w-full sm:w-auto">
          <option value="">Semua Status</option>
          <option value="PENDING">Menunggu</option>
          <option value="IN_PROGRESS">Dikerjakan</option>
          <option value="DONE">Selesai</option>
          <option value="CANCELLED">Dibatalkan</option>
        </select>
      </div>
    </div>

    <!-- Main Content Card -->
    <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden transition-colors duration-300">
      <div class="overflow-x-auto relative">
        <!-- Loading Overlay -->
        <div v-if="loading" class="absolute inset-0 bg-white/60 dark:bg-gray-900/60 backdrop-blur-sm z-10 flex flex-col items-center justify-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
          <span class="mt-4 text-xs font-black text-primary-600 uppercase tracking-widest">Memuat Log...</span>
        </div>

        <table class="w-full">
          <thead>
            <tr class="bg-gray-50/50 dark:bg-gray-700/50">
              <th v-for="h in ['Barang', 'Masalah', 'Catatan', 'Vendor', 'Biaya', 'Status', 'Tanggal', '']" :key="h"
                  class="text-left py-5 px-8 text-[10px] font-black text-gray-400 dark:text-gray-300 uppercase tracking-[0.2em]">{{ h }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50 dark:divide-gray-700">
            <template v-if="paginatedData.length > 0">
              <tr v-for="log in paginatedData" :key="log.id" class="group hover:bg-gray-50/80 dark:hover:bg-gray-700/50 transition-all duration-300">
                <td class="px-8 py-6">
                  <div class="flex items-center gap-4">
                    <div class="w-10 h-10 bg-gray-100 dark:bg-gray-700 rounded-xl flex items-center justify-center text-gray-400 dark:text-gray-500 group-hover:bg-primary-100 group-hover:text-primary-600 transition-colors">
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>
                    </div>
                    <div>
                      <div class="font-black text-gray-900 dark:text-gray-100 text-sm">{{ log.item_name }}</div>
                      <div class="text-[10px] text-gray-400 dark:text-gray-500 font-bold font-mono tracking-tighter">{{ log.item_code }}</div>
                    </div>
                  </div>
                </td>
                <td class="px-8 py-6">
                  <p class="text-xs text-gray-600 dark:text-gray-300 leading-relaxed font-medium max-w-[200px] truncate" :title="log.issue_description">{{ log.issue_description }}</p>
                </td>
                <td class="px-8 py-6">
                  <p class="text-xs text-gray-500 dark:text-gray-400 leading-relaxed italic max-w-[150px] truncate" :title="log.notes || '-'">{{ log.notes || '-' }}</p>
                </td>
                <td class="px-8 py-6 text-xs font-bold text-gray-700 dark:text-gray-300">{{ log.vendor || '---' }}</td>
                <td class="px-8 py-6">
                  <span class="text-sm font-black text-gray-900 dark:text-gray-100">{{ log.cost ? formatCurrency(log.cost) : 'Rp 0' }}</span>
                </td>
                <td class="px-8 py-6">
                  <span class="px-3 py-1 rounded-lg font-black uppercase text-[9px] tracking-widest shadow-sm inline-block" :class="getStatusBadge(log.status)">
                    {{ statusLabels[log.status] }}
                  </span>
                </td>
                <td class="px-8 py-6 whitespace-nowrap text-[11px] font-bold text-gray-400">{{ formatDate(log.reported_at) }}</td>
                <td class="px-8 py-6 text-right">
                  <div class="flex items-center justify-end gap-2 transition-all duration-300">
                    <button v-if="log.status === 'PENDING'" @click="updateLog(log, 'IN_PROGRESS')"
                      class="btn-action-view !text-blue-600 !bg-blue-50 hover:!bg-blue-600 hover:!text-white" title="Kerjakan">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/></svg>
                    </button>
                    <button v-if="log.status === 'IN_PROGRESS'" @click="openFinishModal(log)"
                      class="btn-action-view !text-green-600 !bg-green-50 hover:!bg-green-600 hover:!text-white" title="Selesaikan">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                    </button>
                    <button @click="deleteLog(log)" class="btn-action-delete" title="Hapus">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                    </button>
                  </div>
                </td>
              </tr>
            </template>
            <tr v-else-if="!loading" class="text-center">
              <td colspan="8" class="px-8 py-24 italic text-gray-400 font-medium tracking-widest text-xs uppercase">{{ searchQuery || filterStatus ? 'Pencarian Tidak Ditemukan' : 'Belum Ada Catatan Perbaikan' }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination Bar -->
      <div v-if="totalPages > 1" class="px-8 py-5 bg-gray-50/50 dark:bg-gray-700/20 border-t border-gray-100 dark:border-gray-700 flex flex-col sm:flex-row items-center justify-between gap-4">
        <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">
          Menampilkan <span class="text-primary-600">{{ startRow }}-{{ endRow }}</span> dari <span class="text-gray-900 dark:text-white">{{ filteredData.length }}</span> data
        </span>
        <div class="flex gap-2">
          <button @click="currentPage--" :disabled="currentPage === 1" class="pagination-btn-standard">
            Kembali
          </button>
          <button v-for="p in visiblePages" :key="p" @click="currentPage = p"
                  class="w-10 h-10 rounded-xl text-[11px] font-black transition-all shadow-sm active:scale-95 border"
                  :class="p === currentPage ? 'bg-primary-600 text-white border-primary-600' : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:bg-primary-50 dark:hover:bg-gray-700'">
            {{ p }}
          </button>
          <button @click="currentPage++" :disabled="currentPage === totalPages" class="pagination-btn-standard">
            Lanjut
          </button>
        </div>
      </div>
    </div>

    <!-- Create Maintenance Modal -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showCreateModal" class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4" @click.self="showCreateModal = false">
          <div class="bg-white dark:bg-gray-800 rounded-[3rem] shadow-2xl w-full max-w-xl overflow-hidden animate-scale-up border border-white/20 dark:border-gray-700">
            <!-- Modal Header -->
            <div class="px-10 py-8 bg-gradient-to-br from-primary-900 to-black text-white relative">
              <div class="absolute top-0 right-0 w-32 h-32 bg-primary-500/10 rounded-full blur-2xl"></div>
              <h2 class="text-2xl font-black tracking-tight leading-none">Lapor Kerusakan</h2>
              <p class="text-primary-100/60 text-xs font-bold mt-2 tracking-widest uppercase italic">Laporan Pemeliharaan Aset</p>
              <button @click="showCreateModal = false" class="absolute top-8 right-8 text-white/50 hover:text-white transition-colors">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12"/></svg>
              </button>
            </div>

            <div class="p-10 space-y-8 max-h-[70vh] overflow-y-auto custom-scrollbar">
              <!-- Item Search Group -->
              <div class="space-y-4">
                <div class="flex items-center gap-2 mb-2">
                  <div class="w-1.5 h-6 bg-primary-600 rounded-full"></div>
                  <h3 class="text-xs font-black text-gray-900 uppercase tracking-widest">Langkah 1: Pilih Barang</h3>
                </div>
                
                <div class="relative group">
                  <div class="absolute inset-y-0 left-4 flex items-center text-gray-400 group-focus-within:text-primary-600 transition-colors">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
                  </div>
                  <input 
                    v-model="form.search_query" 
                    class="w-full pl-12 pr-4 py-4 rounded-2xl bg-gray-50 border-none focus:bg-white focus:ring-4 focus:ring-primary-100 transition-all font-bold text-sm" 
                    placeholder="Ketik nama atau kode barang..." 
                    @input="debounceItemSearch"
                  />
                  
                  <!-- Floating Suggestion Results -->
                  <div v-if="suggestedItems.length > 0" class="absolute z-50 w-full mt-2 bg-white rounded-3xl shadow-2xl border border-gray-100 overflow-hidden divide-y divide-gray-50 animate-fade-in ring-1 ring-black/5">
                    <div 
                      v-for="item in suggestedItems" :key="item.id"
                      @click="selectItem(item)"
                      class="px-6 py-4 hover:bg-primary-50 cursor-pointer transition-all flex items-center justify-between group/item"
                    >
                      <div class="flex items-center gap-4">
                        <div class="w-10 h-10 bg-gray-100 rounded-xl flex items-center justify-center text-gray-400 group-hover/item:bg-white transition-colors">
                          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>
                        </div>
                        <div>
                          <div class="text-sm font-black text-gray-900">{{ item.name }}</div>
                          <div class="text-[10px] text-gray-400 font-bold font-mono tracking-tighter">{{ item.code }} • {{ item.location }}</div>
                        </div>
                      </div>
                      <div class="text-[10px] bg-primary-100 text-primary-700 px-2 py-1 rounded-lg font-black opacity-0 group-hover/item:opacity-100 transition-opacity">PILIH</div>
                    </div>
                  </div>
                </div>

                <!-- Active Selection Card -->
                <div v-if="form.item_code" class="bg-primary-50 rounded-2xl p-4 flex items-center justify-between border border-primary-200 animate-scale-up">
                  <div class="flex items-center gap-4">
                    <div class="w-10 h-10 bg-primary-600 rounded-xl flex items-center justify-center text-white shadow-lg shadow-primary-200">
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                    </div>
                    <div>
                      <p class="text-[10px] font-black text-primary-400 uppercase tracking-widest leading-none mb-1">Barang Terpilih</p>
                      <p class="text-sm font-black text-primary-900 leading-none">{{ form.search_query }}</p>
                    </div>
                  </div>
                  <button @click="form.item_code = ''; form.search_query = ''" class="p-2 text-primary-400 hover:text-primary-800 transition-colors">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12"/></svg>
                  </button>
                </div>
              </div>

              <!-- Problem Details Group -->
              <div class="space-y-4">
                <div class="flex items-center gap-2 mb-2">
                  <div class="w-1.5 h-6 bg-red-500 rounded-full"></div>
                  <h3 class="text-xs font-black text-gray-900 uppercase tracking-widest">Langkah 2: Detail Kerusakan</h3>
                </div>
                
                <textarea 
                  v-model="form.issue_description" 
                  rows="3" 
                  class="w-full px-6 py-4 rounded-2xl bg-gray-50 border-none focus:bg-white focus:ring-4 focus:ring-red-100 transition-all font-bold text-sm" 
                  placeholder="Ceritakan apa masalahnya... (Misal: Layar berkedip, tidak bisa nyala, dsb)"
                ></textarea>

                <div class="grid grid-cols-2 gap-4">
                  <div class="space-y-1">
                    <label class="text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-2">Vendor / Teknisi</label>
                    <input v-model="form.vendor" class="w-full px-6 py-4 rounded-2xl bg-gray-50 border-none focus:bg-white focus:ring-4 focus:ring-primary-100 transition-all font-bold text-sm" placeholder="Opsional" />
                  </div>
                  <div class="space-y-1">
                    <label class="text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-2">Estimasi Biaya</label>
                    <div class="relative">
                      <div class="absolute inset-y-0 left-4 flex items-center text-gray-400 font-black text-xs">Rp</div>
                      <input v-model.number="form.cost" type="number" class="w-full pl-10 pr-6 py-4 rounded-2xl bg-gray-50 border-none focus:bg-white focus:ring-4 focus:ring-primary-100 transition-all font-bold text-sm" placeholder="0" />
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Modal Footer -->
            <div class="p-10 pt-0 flex gap-4">
              <button @click="showCreateModal = false" class="flex-1 px-8 py-4 rounded-2xl text-xs font-black text-gray-500 hover:bg-gray-50 transition-all uppercase tracking-widest">Batal</button>
              <button 
                @click="submitCreate" 
                :disabled="submitting || !form.item_code || !form.issue_description" 
                class="btn-premium-action flex-[2]"
              >
                {{ submitting ? 'MEMPROSES...' : 'KIRIM LAPORAN' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Finish / Done Modal -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showFinishModal" class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4" @click.self="showFinishModal = false">
          <div class="bg-white rounded-[3rem] shadow-2xl w-full max-w-lg overflow-hidden animate-scale-up">
            <div class="px-10 py-8 bg-gradient-to-br from-green-600 to-green-900 text-white relative">
              <h2 class="text-2xl font-black tracking-tight leading-none">Selesaikan Perbaikan</h2>
              <div class="bg-white/20 backdrop-blur-sm rounded-2xl p-4 mt-6 flex items-center gap-4">
                <div class="w-12 h-12 bg-white rounded-xl flex items-center justify-center text-green-700 shadow-lg">
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                </div>
                <div>
                  <div class="text-xs font-black text-white uppercase tracking-widest opacity-60">Barang</div>
                  <div class="text-lg font-black leading-none">{{ finishTarget?.item_name }}</div>
                </div>
              </div>
            </div>

            <div class="p-10 space-y-6">
              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-1">
                  <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest ml-2">Final Vendor</label>
                  <input v-model="finishForm.vendor" class="w-full px-6 py-4 rounded-2xl bg-gray-50 border-none focus:ring-4 focus:ring-green-100 transition-all font-bold text-sm" />
                </div>
                <div class="space-y-1">
                  <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest ml-2">Total Biaya Akhir</label>
                  <div class="relative">
                    <div class="absolute inset-y-0 left-4 flex items-center text-gray-400 font-bold text-xs">Rp</div>
                    <input v-model.number="finishForm.cost" type="number" class="w-full pl-10 pr-6 py-4 rounded-2xl bg-gray-50 border-none focus:ring-4 focus:ring-green-100 transition-all font-bold text-sm" />
                  </div>
                </div>
              </div>
              <div class="space-y-1">
                <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest ml-2">Hasil Akhir / Catatan</label>
                <textarea v-model="finishForm.notes" rows="2" class="w-full px-6 py-4 rounded-2xl bg-gray-50 border-none focus:ring-4 focus:ring-green-100 transition-all font-bold text-sm" placeholder="CTH: Komponen IC diganti, layar sudah normal..."></textarea>
              </div>
            </div>

            <div class="px-10 pb-10 flex gap-4">
              <button @click="showFinishModal = false" class="flex-1 px-8 py-4 rounded-2xl text-xs font-black text-gray-400 hover:bg-gray-50 transition-all uppercase tracking-widest">Batal</button>
              <button 
                @click="submitFinish" 
                :disabled="submitting"
                class="btn-premium-action flex-[2] !bg-green-600 !shadow-green-200"
              >
                {{ submitting ? 'MENYIMPAN...' : 'KONFIRMASI SELESAI' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, h } from 'vue'
import api from '../utils/api'

// Icons for stats cards
const IconPending = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-6 h-6' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z' })]) }
const IconProgress = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-6 h-6' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z' })]) }
const IconDone = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-6 h-6' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M5 13l4 4L19 7' })]) }
const IconCancel = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-6 h-6' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M6 18L18 6M6 6l12 12' })]) }

const logs = ref([])
const loading = ref(false)
const filterStatus = ref('')
const searchQuery = ref('')
const currentPage = ref(1)
const perPage = 10

const showCreateModal = ref(false)
const showFinishModal = ref(false)
const submitting = ref(false)
const formError = ref('')
const finishTarget = ref(null)

const form = ref({ item_code: '', search_query: '', issue_description: '', cost: 0, vendor: '', notes: '' })
const finishForm = ref({ cost: 0, vendor: '', notes: '' })

const suggestedItems = ref([])
let searchTimeout = null

const statusLabels = { PENDING: 'Menunggu', IN_PROGRESS: 'Dikerjakan', DONE: 'Selesai', CANCELLED: 'Dibatalkan' }

const stats = computed(() => {
  const all = logs.value
  return {
    pending: all.filter(l => l.status === 'PENDING').length,
    inProgress: all.filter(l => l.status === 'IN_PROGRESS').length,
    done: all.filter(l => l.status === 'DONE').length,
    cancelled: all.filter(l => l.status === 'CANCELLED').length,
  }
})

const statsMap = computed(() => ({
  pending: { label: 'Menunggu', value: stats.value.pending, color: 'bg-yellow-500', gradient: 'bg-gradient-to-br from-yellow-400 to-orange-500', icon: IconPending },
  progress: { label: 'Dikerjakan', value: stats.value.inProgress, color: 'bg-blue-500', gradient: 'bg-gradient-to-br from-blue-400 to-indigo-600', icon: IconProgress },
  done: { label: 'Selesai', value: stats.value.done, color: 'bg-green-500', gradient: 'bg-gradient-to-br from-green-400 to-emerald-600', icon: IconDone },
  cancelled: { label: 'Batal', value: stats.value.cancelled, color: 'bg-gray-500', gradient: 'bg-gradient-to-br from-gray-400 to-slate-600', icon: IconCancel },
}))

// Client-side Search and Pagination
const filteredData = computed(() => {
  let result = logs.value
  
  if (filterStatus.value) {
    result = result.filter(l => l.status === filterStatus.value)
  }
  
  const q = searchQuery.value.toLowerCase().trim()
  if (q) {
    result = result.filter(l => 
      l.item_name?.toLowerCase().includes(q) || 
      l.item_code?.toLowerCase().includes(q) ||
      l.issue_description?.toLowerCase().includes(q) ||
      l.vendor?.toLowerCase().includes(q)
    )
  }
  
  return result
})

const totalPages = computed(() => Math.max(1, Math.ceil(filteredData.value.length / perPage)))
const startRow = computed(() => (currentPage.value - 1) * perPage + 1)
const endRow = computed(() => Math.min(currentPage.value * perPage, filteredData.value.length))
const paginatedData = computed(() => filteredData.value.slice((currentPage.value - 1) * perPage, currentPage.value * perPage))

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

watch([searchQuery, filterStatus], () => {
  currentPage.value = 1
})

async function fetchLogs() {
  loading.value = true
  try {
    // We fetch all records locally to allow smooth client-side filtering and stats calculation
    const { data } = await api.get('/maintenance')
    if (data.success) logs.value = data.data
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

async function exportExcel() {
  try {
    const params = filterStatus.value ? `?status=${filterStatus.value}` : ''
    const response = await api.get(`/reports/export/maintenance${params}`, { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    const dateStr = new Date().toISOString().slice(0, 10).replace(/-/g, '')
    link.setAttribute('download', `SIS-INV_Log-Perbaikan_${dateStr}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (e) {
    console.error("Failed to export Excel", e)
    alert("Gagal mengunduh file Excel")
  }
}

const debounceItemSearch = () => {
  clearTimeout(searchTimeout)
  if (!form.value.search_query) {
    suggestedItems.value = []
    return
  }
  searchTimeout = setTimeout(async () => {
    try {
      const { data } = await api.get(`/items?search=${form.value.search_query}&page_size=10`)
      if (data.success) {
        suggestedItems.value = data.data.items || []
      }
    } catch (e) { 
      console.error(e)
      suggestedItems.value = []
    }
  }, 300)
}

const selectItem = (item) => {
  form.value.item_code = item.code
  form.value.search_query = item.name
  suggestedItems.value = []
}

async function submitCreate() {
  if (!form.value.item_code || !form.value.issue_description) {
    formError.value = 'Barang dan deskripsi kerusakan wajib diisi.'
    return
  }
  formError.value = ''
  submitting.value = true
  try {
    const { data } = await api.post('/maintenance', {
      item_code: form.value.item_code,
      issue_description: form.value.issue_description,
      cost: form.value.cost,
      vendor: form.value.vendor,
      notes: form.value.notes
    })
    if (data.success) {
      showCreateModal.value = false
      form.value = { item_code: '', search_query: '', issue_description: '', cost: 0, vendor: '', notes: '' }
      fetchLogs()
    }
  } catch (e) {
    formError.value = e.response?.data?.error || 'Gagal menyimpan data'
  } finally { submitting.value = false }
}

async function updateLog(log, newStatus) {
  if (!confirm(`Ubah status "${statusLabels[log.status]}" → "${statusLabels[newStatus]}"?`)) return
  try {
    await api.put(`/maintenance/${log.id}`, { status: newStatus, cost: log.cost || 0, vendor: log.vendor || '', notes: log.notes || '' })
    fetchLogs()
  } catch (e) { alert('Gagal mengubah status') }
}

function openFinishModal(log) {
  finishTarget.value = log
  finishForm.value = { cost: log.cost || 0, vendor: log.vendor || '', notes: log.notes || '' }
  showFinishModal.value = true
}

async function submitFinish() {
  submitting.value = true
  try {
    await api.put(`/maintenance/${finishTarget.value.id}`, { status: 'DONE', ...finishForm.value })
    showFinishModal.value = false
    fetchLogs()
  } catch (e) { alert('Gagal menyelesaikan perbaikan') }
  finally { submitting.value = false }
}

async function deleteLog(log) {
  if (!confirm(`Hapus log perbaikan untuk "${log.item_name}"?`)) return
  try {
    await api.delete(`/maintenance/${log.id}`)
    await fetchLogs()
  } catch (e) { 
    alert(e.response?.data?.error || 'Gagal menghapus') 
  }
}

const formatDate = (d) => d ? new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }).format(new Date(d)) : '-'
const formatCurrency = (v) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)
const getStatusBadge = (s) => ({
  PENDING: 'bg-yellow-50 text-yellow-700 border border-yellow-200',
  IN_PROGRESS: 'bg-blue-50 text-blue-700 border border-blue-200',
  DONE: 'bg-green-50 text-green-700 border border-green-200',
  CANCELLED: 'bg-gray-100 text-gray-500 border border-gray-200',
}[s] || 'bg-gray-50 text-gray-700')

onMounted(fetchLogs)
</script>
