<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Kelola User</h1>
          <p class="text-primary-100/70 text-sm font-medium">Manajemen hak akses dan akun staf sekolah</p>
        </div>
        
        <div class="flex items-center gap-3 backdrop-blur-md bg-white/10 p-2 rounded-2xl border border-white/10">
          <button @click="showModal = true" 
                  class="bg-white text-primary-900 hover:bg-primary-50 px-6 py-2.5 rounded-xl text-[10px] font-black transition-all flex items-center gap-2 shadow-xl active:scale-95">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 4v16m8-8H4" /></svg>
            TAMBAH USER BARU
          </button>
        </div>
      </div>
    </div>

    <!-- Main Table Card -->
    <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden transition-all duration-300">
      <div class="overflow-x-auto relative">
        <!-- Loading Overlay -->
        <div v-if="loading" class="absolute inset-0 bg-white/60 dark:bg-gray-900/60 backdrop-blur-sm z-10 flex flex-col items-center justify-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
          <span class="mt-4 text-xs font-black text-primary-600 uppercase tracking-widest">Memuat Pengguna...</span>
        </div>

        <table class="w-full">
          <thead>
            <tr class="bg-gray-50/50 dark:bg-gray-700/30">
              <th v-for="h in ['Nama Lengkap', 'Username', 'Role', 'Status Keaktifan', '']" :key="h"
                  class="text-left py-5 px-8 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">{{ h }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50 dark:divide-gray-700">
            <template v-if="users.length > 0">
              <tr v-for="u in users" :key="u.id" class="group hover:bg-primary-50/50 dark:hover:bg-primary-900/10 transition-all duration-300">
                <td class="px-8 py-6">
                  <div class="flex items-center gap-4">
                    <div class="w-12 h-12 rounded-2xl bg-gradient-to-br from-primary-500 to-primary-700 text-white flex items-center justify-center text-sm font-black shadow-lg shadow-primary-500/20 uppercase border border-white/20">
                      {{ u.full_name?.charAt(0) }}
                    </div>
                    <div>
                      <p class="text-sm font-black text-gray-900 dark:text-white leading-none uppercase tracking-tight">{{ u.full_name }}</p>
                      <p class="text-[10px] text-gray-400 mt-1 font-bold font-mono tracking-tighter">{{ u.id }}</p>
                    </div>
                  </div>
                </td>
                <td class="px-8 py-6">
                  <span class="text-sm font-bold text-gray-600 dark:text-gray-300 tracking-tight">{{ u.username }}</span>
                </td>
                <td class="px-8 py-6">
                  <span class="px-3 py-1 rounded-lg font-black uppercase text-[9px] tracking-widest shadow-sm inline-block" :class="roleBadgeClass(u.role)">
                    {{ roleLabel(u.role) }}
                  </span>
                </td>
                <td class="px-8 py-6">
                  <span class="inline-flex items-center gap-2 px-3 py-1.5 rounded-xl text-[10px] font-black uppercase tracking-widest" 
                        :class="u.is_active ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600' : 'bg-gray-100 dark:bg-gray-700 text-gray-400'">
                    <span class="w-1.5 h-1.5 rounded-full animate-pulse" :class="u.is_active ? 'bg-emerald-500' : 'bg-gray-400'"></span>
                    {{ u.is_active ? 'Aktif' : 'Nonaktif' }}
                  </span>
                </td>
                <td class="px-8 py-6 text-right">
                  <div class="flex items-center justify-end gap-2 transition-opacity">
                    <button @click="openEditModal(u)" class="p-2 text-blue-500 bg-blue-50 dark:bg-blue-900/30 rounded-xl hover:bg-blue-500 hover:text-white transition-all transform hover:scale-110 active:scale-95 relative z-10" title="Edit">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                    </button>
                    <button @click="deleteUser(u)" class="p-2 text-red-500 bg-red-50 dark:bg-red-900/30 rounded-xl hover:bg-red-500 hover:text-white transition-all transform hover:scale-110 active:scale-95 relative z-10" title="Hapus">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                    </button>
                  </div>
                </td>
              </tr>
            </template>
            <tr v-else-if="!loading" class="text-center">
              <td colspan="5" class="px-8 py-24 italic text-gray-400 font-medium tracking-widest text-xs uppercase">Belum ada data user tersimpan</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- User Modal (Create/Edit) -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/80 backdrop-blur-xl animate-fade-in" @click.self="closeModal">
      <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] shadow-3xl w-full max-w-md animate-scale-up overflow-hidden border border-white/20 dark:border-gray-700">
        <div class="p-8 border-b border-gray-100 dark:border-gray-700 flex items-center justify-between">
          <div>
            <h3 class="text-2xl font-black text-gray-900 dark:text-white leading-none capitalize">{{ editingUser ? 'Edit User' : 'User Baru' }}</h3>
            <p class="text-[10px] font-black text-primary-500 uppercase tracking-widest mt-2">{{ editingUser ? 'Perbarui Hak Akses' : 'Atur Akun Baru' }}</p>
          </div>
          <button @click="closeModal" class="text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 p-3 rounded-2xl transition-all">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>
        
        <form @submit.prevent="saveUser" class="p-8 space-y-6">
          <div class="space-y-2">
            <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Nama Lengkap</label>
            <input v-model="form.full_name" class="input-field rounded-2xl h-14" placeholder="Contoh: Budi Santoso" required />
          </div>
          <div class="space-y-2">
            <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Username</label>
            <input v-model="form.username" class="input-field rounded-2xl h-14" placeholder="budi123" required :disabled="editingUser" />
          </div>
          <div class="space-y-2">
            <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Password Baru</label>
            <input v-model="form.password" type="password" class="input-field rounded-2xl h-14" :placeholder="editingUser ? 'Kosongkan jika tidak diganti' : 'Minimal 6 karakter'" :required="!editingUser" minlength="6" />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Role Utama</label>
              <select v-model="form.role" class="input-field rounded-2xl h-14" required>
                <option value="TEACHER">Guru/Staff</option>
                <option value="HEAD">Kepasek</option>
                <option value="ADMIN">Administrator</option>
              </select>
            </div>
            <div class="space-y-2">
              <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Status</label>
              <select v-model="form.is_active" class="input-field rounded-2xl h-14" required>
                <option :value="true">Aktif</option>
                <option :value="false">Nonaktif</option>
              </select>
            </div>
          </div>
          
          <div class="flex flex-col sm:flex-row gap-3 pt-6 border-t border-gray-100 dark:border-gray-700">
            <button type="button" @click="closeModal" class="btn-secondary flex-1 py-4 rounded-xl font-black text-[10px] tracking-widest">BATAL</button>
            <button type="submit" :disabled="submitting" 
                    class="bg-primary-600 text-white flex-[2] py-4 rounded-xl font-black text-[10px] tracking-[0.3em] shadow-xl shadow-primary-500/20 active:scale-95 disabled:opacity-30 transition-all uppercase">
              {{ submitting ? 'MEMPROSES...' : (editingUser ? 'PERBARUI DATA' : 'DAFTARKAN USER') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../utils/api'

const users = ref([])
const loading = ref(true)
const showModal = ref(false)
const submitting = ref(false)
const editingUser = ref(null)

const form = ref({ full_name: '', username: '', password: '', role: 'TEACHER', is_active: true })
const roleLabels = { ADMIN: 'Admin', TEACHER: 'Guru/Staff', HEAD: 'Kepasek' }
const roleLabel = (role) => roleLabels[role] || role

const roleBadgeClass = (role) => ({
  ADMIN: 'bg-purple-50 dark:bg-purple-900/40 text-purple-600 dark:text-purple-400',
  TEACHER: 'bg-blue-50 dark:bg-blue-900/40 text-blue-600 dark:text-blue-400',
  HEAD: 'bg-amber-50 dark:bg-amber-900/40 text-amber-600 dark:text-amber-400',
}[role] || 'bg-gray-50 text-gray-800')

async function fetchUsers() {
  loading.value = true
  try { const { data } = await api.get('/users'); users.value = data.data || [] } 
  catch (e) { console.error('Gagal ambil data user:', e) } 
  finally { loading.value = false }
}

function openEditModal(user) {
  editingUser.value = user
  form.value = { full_name: user.full_name, username: user.username, password: '', role: user.role, is_active: user.is_active }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingUser.value = null
  form.value = { full_name: '', username: '', password: '', role: 'TEACHER', is_active: true }
}

async function saveUser() {
  submitting.value = true
  try {
    if (editingUser.value) {
      const payload = { ...form.value }; if (!payload.password) delete payload.password
      await api.put(`/users/${editingUser.value.id}`, payload)
    } else { await api.post('/users', form.value) }
    closeModal(); await fetchUsers()
  } catch (e) { alert(e.response?.data?.error || 'Gagal menyimpan user') } 
  finally { submitting.value = false }
}

async function deleteUser(user) {
  if (confirm(`Apakah Anda yakin ingin menghapus user ${user.full_name}?`)) {
    try { 
      await api.delete(`/users/${user.id}`)
      await fetchUsers() 
    } catch (e) { 
      alert(e.response?.data?.error || 'Gagal menghapus user') 
    }
  }
}

onMounted(fetchUsers)
</script>
