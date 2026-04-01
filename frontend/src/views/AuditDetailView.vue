<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Summary -->
    <div v-if="auditStore.currentSession" class="relative overflow-hidden bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm p-10 transition-all duration-300">
        <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
            <div class="space-y-1">
                <router-link to="/stock-opname" class="inline-flex items-center gap-2 text-[10px] font-black text-primary-600 uppercase tracking-widest mb-4 hover:translate-x-[-1px] transition-transform">
                    <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M10 19l-7-7m0 0l7-7m-7 7h18"/></svg>
                    KEMBALI KE DAFTAR
                </router-link>
                <div class="flex items-center gap-3">
                    <h1 class="text-3xl font-black text-gray-900 dark:text-white uppercase tracking-tight">Audit: {{ auditStore.currentSession.session?.location_name }}</h1>
                    <span :class="auditStore.currentSession.session?.status === 'OPEN' ? 'bg-amber-50 text-amber-600 dark:bg-amber-900/20' : 'bg-emerald-50 text-emerald-600 dark:bg-emerald-900/20'"
                      class="px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest">
                      {{ auditStore.currentSession.session?.status }}
                    </span>
                </div>
                <p class="text-[11px] font-bold text-gray-400 uppercase tracking-widest">{{ formatDate(auditStore.currentSession.session?.started_at) }} • Petugas: {{ auditStore.currentSession.session?.user_name }}</p>
            </div>

            <div class="flex flex-wrap gap-4">
                <button v-if="auditStore.currentSession.session?.status === 'CLOSED'" 
                        @click="exportAudit" 
                        class="bg-emerald-500 text-white px-8 py-4 rounded-2xl text-[10px] font-black hover:bg-emerald-600 active:scale-95 transition-all uppercase tracking-widest shadow-xl shadow-emerald-500/20 flex items-center gap-2">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2-2z" /></svg>
                    EKSPOR EXCEL
                </button>
                <template v-if="auditStore.currentSession.session?.status === 'OPEN'">
                    <button @click="showScanModal = true" class="bg-primary-600 text-white px-8 py-4 rounded-2xl text-[10px] font-black hover:scale-105 active:scale-95 transition-all uppercase tracking-widest shadow-xl shadow-primary-500/20">
                        SCAN BARANG
                    </button>
                    <button @click="handleFinishAudit" class="bg-gray-900 text-white dark:bg-gray-700 px-8 py-4 rounded-2xl text-[10px] font-black hover:bg-black active:scale-95 transition-all uppercase tracking-widest">
                        SELESAIKAN AUDIT
                    </button>
                </template>
            </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mt-10 p-8 bg-gray-50 dark:bg-gray-900/50 rounded-[2.5rem]">
            <div class="text-center md:text-left space-y-1">
                <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest">Total Aset di Lokasi</p>
                <p class="text-4xl font-black text-gray-900 dark:text-white leading-none">{{ (auditStore.currentSession.found_items?.length || 0) + (auditStore.currentSession.missing_items?.length || 0) }}</p>
            </div>
            <div class="text-center md:text-left space-y-1 border-y md:border-y-0 md:border-x border-gray-100 dark:border-gray-700 md:px-8">
                <p class="text-[10px] font-black text-emerald-500 uppercase tracking-widest">Ditemukan (Sesuai)</p>
                <p class="text-4xl font-black text-emerald-600 leading-none">{{ auditStore.currentSession.found_items?.length || 0 }}</p>
            </div>
            <div class="text-center md:text-left space-y-1">
                <p class="text-[10px] font-black text-rose-500 uppercase tracking-widest">Tidak Ditemukan (Missing)</p>
                <p class="text-4xl font-black text-rose-600 leading-none">{{ auditStore.currentSession.missing_items?.length || 0 }}</p>
            </div>
        </div>
    </div>

    <!-- Details Tabs (Found vs Missing) -->
    <div v-if="auditStore.currentSession" class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden min-h-[400px]">
        <div class="flex border-b border-gray-50 dark:border-gray-700">
            <button v-for="tab in ['FOUND', 'MISSING']" :key="tab"
                    @click="activeTab = tab"
                    class="flex-1 py-6 text-[10px] font-black uppercase tracking-[0.3em] transition-all relative overflow-hidden"
                    :class="activeTab === tab ? 'text-primary-600' : 'text-gray-400 hover:text-gray-600 bg-gray-50/50 dark:bg-gray-900/30'">
                {{ tab === 'FOUND' ? 'ASET DITEMUKAN' : 'ASET TIDAK DITEMUKAN' }}
                <div v-if="activeTab === tab" class="absolute bottom-0 left-0 w-full h-1 bg-primary-600"></div>
            </button>
        </div>

        <!-- FOUND LIST -->
        <div v-if="activeTab === 'FOUND'" class="animate-scale-up p-4">
            <table class="w-full">
                <thead>
                    <tr class="bg-emerald-50/50 dark:bg-emerald-900/10">
                        <th class="px-8 py-4 text-left text-[9px] font-black text-emerald-600 uppercase tracking-widest">Kode</th>
                        <th class="px-8 py-4 text-left text-[9px] font-black text-emerald-600 uppercase tracking-widest">Nama Barang</th>
                        <th class="px-8 py-4 text-left text-[9px] font-black text-emerald-600 uppercase tracking-widest">Kondisi Fisik</th>
                        <th class="px-8 py-4 text-left text-[9px] font-black text-emerald-600 uppercase tracking-widest">Waktu Scan</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-gray-50 dark:divide-gray-700">
                    <tr v-for="item in auditStore.currentSession.found_items" :key="item.id" class="hover:bg-gray-50 dark:hover:bg-gray-900/20">
                        <td class="px-8 py-5 text-[11px] font-black text-gray-900 dark:text-white font-mono">{{ item.item_code }}</td>
                        <td class="px-8 py-5 text-[11px] font-black text-gray-900 dark:text-white uppercase">{{ item.item_name }}</td>
                        <td class="px-8 py-5">
                            <span :class="getConditionColor(item.found_condition)" class="px-2 py-1 rounded-lg text-[9px] font-black uppercase tracking-widest bg-gray-100 dark:bg-gray-700">
                                {{ item.found_condition }}
                            </span>
                        </td>
                        <td class="px-8 py-5 text-[10px] font-bold text-gray-400">{{ formatTime(item.scanned_at) }}</td>
                    </tr>
                </tbody>
            </table>
            <div v-if="!auditStore.currentSession.found_items?.length" class="py-20 text-center uppercase text-[10px] font-black text-gray-300 tracking-widest">
                Belum ada aset yang di-scan.
            </div>
        </div>

        <!-- MISSING LIST -->
        <div v-else class="animate-scale-up p-4">
            <table class="w-full">
                <thead>
                    <tr class="bg-rose-50/50 dark:bg-rose-900/10">
                        <th class="px-8 py-4 text-left text-[9px] font-black text-rose-600 uppercase tracking-widest">Kode</th>
                        <th class="px-8 py-4 text-left text-[9px] font-black text-rose-600 uppercase tracking-widest">Nama Barang</th>
                        <th class="px-8 py-4 text-left text-[9px] font-black text-rose-600 uppercase tracking-widest">Status Sistem</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-gray-50 dark:divide-gray-700">
                    <tr v-for="item in auditStore.currentSession.missing_items" :key="item.id" class="hover:bg-gray-50 dark:hover:bg-gray-900/20">
                        <td class="px-8 py-5 text-[11px] font-black text-gray-900 dark:text-white font-mono">{{ item.code }}</td>
                        <td class="px-8 py-5 text-[11px] font-black text-gray-900 dark:text-white uppercase">{{ item.name }}</td>
                        <td class="px-8 py-5">
                            <span class="px-2 py-1 rounded-lg text-[9px] font-black uppercase tracking-widest bg-gray-100 dark:bg-gray-700">
                                {{ item.condition }}
                            </span>
                        </td>
                    </tr>
                </tbody>
            </table>
            <div v-if="!auditStore.currentSession.missing_items?.length" class="py-20 text-center uppercase text-[10px] font-black text-emerald-400 tracking-[.5em]">
                SEMUA ASET DITEMUKAN SITUASI AMAN!
            </div>
        </div>
    </div>

    <!-- SCAN MODAL -->
    <div v-if="showScanModal" class="fixed inset-0 z-50 flex items-center justify-center p-6 bg-gray-950/40 backdrop-blur-sm animate-fade-in transition-all duration-300">
        <div class="bg-white dark:bg-gray-800 w-full max-w-lg rounded-[2.5rem] shadow-2xl overflow-hidden animate-scale-up border border-gray-100 dark:border-gray-700">
            <div class="bg-primary-600 p-8 text-white relative">
                <button @click="showScanModal = false" class="absolute top-6 right-6 text-white/50 hover:text-white transition-colors">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12"/></svg>
                </button>
                <h3 class="text-2xl font-black uppercase tracking-tight">Scan Verifikasi</h3>
            </div>
            
            <div class="p-8 space-y-6">
                <!-- Status Scan -->
                <div v-if="!scannedCode" class="space-y-4">
                    <div class="p-4 bg-gray-900 rounded-2xl overflow-hidden">
                        <QrScanner @scanned="handleItemScanned" />
                    </div>
                    <p class="text-center text-[10px] font-black text-gray-400 uppercase tracking-widest">Arahkan kamera ke label QR Barang</p>
                </div>

                <!-- Input Detail Form -->
                <div v-else class="space-y-6 animate-scale-up">
                    <div class="p-5 bg-primary-50 dark:bg-primary-900/20 rounded-2xl border border-primary-100 dark:border-primary-800 flex items-center justify-between">
                        <div>
                            <p class="text-[9px] font-black uppercase text-primary-400 tracking-widest leading-none mb-1">Aset Terdeteksi</p>
                            <h4 class="text-lg font-black text-gray-900 dark:text-white leading-none">{{ scannedCode }}</h4>
                        </div>
                        <button @click="scannedCode = ''" class="bg-white dark:bg-gray-800 p-2 rounded-xl text-primary-600 shadow-sm">GANTI</button>
                    </div>

                    <div>
                        <label class="block text-[10px] font-black text-gray-400 uppercase tracking-widest mb-3 pl-1">Kondisi Fisik Saat Ini</label>
                        <div class="grid grid-cols-2 gap-3">
                            <button v-for="c in ['GOOD', 'DAMAGED', 'IN_REPAIR', 'LOST']" :key="c"
                                    @click="auditForm.condition = c"
                                    class="p-4 rounded-2xl border-2 text-[10px] font-black uppercase tracking-widest transition-all"
                                    :class="auditForm.condition === c ? 'border-primary-600 bg-primary-50 dark:bg-primary-900/20 text-primary-600' : 'border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/30 text-gray-400'">
                                {{ c === 'GOOD' ? 'BAIK' : (c === 'IN_REPAIR' ? 'SERVIS' : (c === 'DAMAGED' ? 'RUSAK' : 'HILANG')) }}
                            </button>
                        </div>
                    </div>

                    <button @click="submitScanDetail" :disabled="loading" 
                            class="w-full py-5 bg-primary-600 text-white rounded-2xl font-black text-xs uppercase tracking-[0.3em] shadow-xl shadow-primary-500/20 transition-all active:scale-95">
                        KONFIRMASI TEMUAN
                    </button>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuditStore } from '../stores/audit'
import QrScanner from '../components/QrScanner.vue'

const route = useRoute()
const auditStore = useAuditStore()
const activeTab = ref('FOUND')
const showScanModal = ref(false)
const scannedCode = ref('')
const loading = ref(false)

const auditForm = ref({ condition: 'GOOD', notes: '' })

const formatDate = (date) => new Date(date).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' })
const formatTime = (date) => new Date(date).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })

const getConditionColor = (cond) => {
    switch(cond) {
        case 'GOOD': return 'text-emerald-500'
        case 'DAMAGED': return 'text-amber-500'
        case 'IN_REPAIR': return 'text-blue-500'
        case 'LOST': return 'text-rose-500'
    }
}

const handleItemScanned = (code) => {
    scannedCode.value = code
}

const submitScanDetail = async () => {
    loading.value = true
    try {
        await auditStore.scanAuditItem(route.params.id, scannedCode.value, auditForm.condition, auditForm.notes)
        scannedCode.value = ''
        auditForm.value = { condition: 'GOOD', notes: '' }
        await auditStore.fetchSessionDetail(route.params.id)
        // Keep scan modal open for next scanning but reset the inner step if we want it to be "rapid scanning"
        // For now, let's keep it simple: close it or reset it.
    } catch (err) {
        alert(err)
    } finally {
        loading.value = false
    }
}

async function handleFinishAudit() {
    if (!confirm('Selesaikan audit? Data tidak bisa diubah lagi setelah status CLOSED.')) return
    try {
        await auditStore.closeSession(route.params.id)
        await auditStore.fetchSessionDetail(route.params.id)
    } catch (e) {
        alert('Gagal menyelesaikan audit')
    }
}

async function exportAudit() {
    try {
        const response = await api.get(`/reports/audit-sessions/${route.params.id}/export`, { responseType: 'blob' })
        const url = window.URL.createObjectURL(new Blob([response.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', `Hasil_Audit_${auditStore.currentSession.session?.location_name}_${new Date().toISOString().slice(0,10)}.xlsx`)
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
    } catch (e) {
        alert('Gagal mendownload hasil audit')
    }
}

onMounted(() => {
    auditStore.fetchSessionDetail(route.params.id)
})
</script>

<style scoped>
.animate-fade-in { animation: fadeIn 0.5s ease-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
.animate-scale-up { animation: scaleUp 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); }
@keyframes scaleUp { from { opacity: 0; transform: scale(0.95); } to { opacity: 1; transform: scale(1); } }
</style>
