<template>
  <div class="animate-fade-in space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">Kelola User</h1>
      <button @click="showCreateModal = true" class="btn-primary flex items-center gap-2">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Tambah User
      </button>
    </div>

    <!-- Users Table -->
    <div class="card p-0 overflow-hidden shadow-sm border border-gray-100">
      <div v-if="loading" class="flex items-center justify-center py-12">
        <svg class="animate-spin w-8 h-8 text-primary-500" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
        </svg>
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50/80 border-b border-gray-100">
            <tr>
              <th class="text-left text-xs font-bold text-gray-500 uppercase tracking-wider px-6 py-4">Nama Lengkap</th>
              <th class="text-left text-xs font-bold text-gray-500 uppercase tracking-wider px-6 py-4">Username</th>
              <th class="text-left text-xs font-bold text-gray-500 uppercase tracking-wider px-6 py-4">Role</th>
              <th class="text-left text-xs font-bold text-gray-500 uppercase tracking-wider px-6 py-4">Status</th>
              <th class="text-right text-xs font-bold text-gray-500 uppercase tracking-wider px-6 py-4">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 bg-white">
            <tr v-for="u in users" :key="u.id" class="hover:bg-gray-50/50 transition-colors group">
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-xl bg-primary-50 text-primary-600 flex items-center justify-center text-sm font-bold border border-primary-100 uppercase">
                    {{ u.full_name?.charAt(0) }}
                  </div>
                  <div>
                    <p class="text-sm font-semibold text-gray-900 line-clamp-1">{{ u.full_name }}</p>
                    <p class="text-[10px] text-gray-400 font-mono">{{ u.id.split('-').shift() }}...</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 text-sm text-gray-600 font-medium">{{ u.username }}</td>
              <td class="px-6 py-4 text-xs">
                <span class="px-2.5 py-1 rounded-lg font-bold uppercase tracking-wider" :class="roleBadgeClass(u.role)">{{ roleLabel(u.role) }}</span>
              </td>
              <td class="px-6 py-4">
                <span class="inline-flex items-center gap-1.5 px-2 py-1 rounded-lg text-xs font-bold" :class="u.is_active ? 'bg-green-50 text-green-700' : 'bg-gray-100 text-gray-500'">
                  <span class="w-1.5 h-1.5 rounded-full" :class="u.is_active ? 'bg-green-500' : 'bg-gray-400'"></span>
                  {{ u.is_active ? 'Status: Aktif' : 'Status: Nonaktif' }}
                </span>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button @click="openEditModal(u)" class="p-2 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors" title="Edit User">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button @click="deleteUser(u)" class="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors" title="Hapus User">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="users.length === 0">
              <td colspan="5" class="px-6 py-16 text-center text-gray-400 text-sm italic">Belum ada data pengguna dalam sistem.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- User Modal (Create/Edit) -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/60 backdrop-blur-sm" @click.self="closeModal">
      <div class="bg-white rounded-2xl shadow-2xl w-full max-w-md animate-scale-up overflow-hidden">
        <div class="flex items-center justify-between px-6 py-5 border-b border-gray-100 bg-gray-50/50">
          <div>
            <h3 class="text-xl font-bold text-gray-900">{{ editingUser ? 'Edit Pengguna' : 'Tambah Pengguna' }}</h3>
            <p class="text-xs text-gray-500 mt-1 uppercase font-bold tracking-widest">{{ editingUser ? 'Perbarui Akses' : 'Buat Akun Baru' }}</p>
          </div>
          <button @click="closeModal" class="text-gray-400 hover:bg-gray-100 p-2 rounded-xl transition-all">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <form @submit.prevent="saveUser" class="p-6 space-y-5">
          <div>
            <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-1.5 ml-1">Nama Lengkap</label>
            <input v-model="form.full_name" class="input-field rounded-xl border-gray-200 focus:ring-primary-500" placeholder="Contoh: Budi Santoso" required />
          </div>
          <div>
            <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-1.5 ml-1">Username</label>
            <input v-model="form.username" class="input-field rounded-xl border-gray-200" placeholder="budi123" required :disabled="editingUser" />
          </div>
          <div>
            <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-1.5 ml-1">
              {{ editingUser ? 'Ganti Password (Opsional)' : 'Password' }}
            </label>
            <input v-model="form.password" type="password" class="input-field rounded-xl border-gray-200" :placeholder="editingUser ? 'Biarkan kosong jika tidak diganti' : 'Minimal 6 karakter'" :required="!editingUser" minlength="6" />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-1.5 ml-1">Role</label>
              <select v-model="form.role" class="input-field rounded-xl border-gray-200" required>
                <option value="TEACHER">Guru/Staff</option>
                <option value="HEAD">Kepasek</option>
                <option value="ADMIN">Administrator</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-1.5 ml-1">Status Akun</label>
              <select v-model="form.is_active" class="input-field rounded-xl border-gray-200" required>
                <option :value="true">Aktif</option>
                <option :value="false">Nonaktif</option>
              </select>
            </div>
          </div>
          
          <div class="flex gap-3 pt-4">
            <button type="button" @click="closeModal" class="btn-secondary flex-1 py-3 rounded-xl font-bold">Batal</button>
            <button type="submit" :disabled="submitting" class="btn-primary flex-1 py-3 rounded-xl font-bold">
              {{ submitting ? 'Menyimpan...' : (editingUser ? 'Perbarui' : 'Daftarkan') }}
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

const form = ref({ 
  full_name: '', 
  username: '', 
  password: '', 
  role: 'TEACHER',
  is_active: true
})

const roleLabels = { ADMIN: 'Admin', TEACHER: 'Guru', HEAD: 'Kepsek' }
const roleLabel = (role) => roleLabels[role] || role

const roleBadgeClass = (role) => ({
  ADMIN: 'bg-purple-100 text-purple-700',
  TEACHER: 'bg-blue-100 text-blue-700',
  HEAD: 'bg-amber-100 text-amber-700',
}[role] || 'bg-gray-100 text-gray-800')

async function fetchUsers() {
  loading.value = true
  try {
    const { data } = await api.get('/users')
    users.value = data.data || []
  } catch (e) {
    console.error('Failed to fetch users:', e)
  } finally {
    loading.value = false
  }
}

function openEditModal(user) {
  editingUser.value = user
  form.value = {
    full_name: user.full_name,
    username: user.username,
    password: '', // Password always empty on load
    role: user.role,
    is_active: user.is_active
  }
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
      // Update existing user
      const payload = { ...form.value }
      if (!payload.password) delete payload.password // Remove if blank
      await api.put(`/users/${editingUser.value.id}`, payload)
    } else {
      // Create new user
      await api.post('/users', form.value)
    }
    closeModal()
    await fetchUsers()
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal menyimpan user')
  } finally {
    submitting.value = false
  }
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
