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
    <div class="card p-0 overflow-hidden">
      <div v-if="loading" class="flex items-center justify-center py-12">
        <svg class="animate-spin w-8 h-8 text-primary-500" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
        </svg>
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Nama</th>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Username</th>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Role</th>
              <th class="text-left text-xs font-semibold text-gray-600 uppercase px-6 py-3">Status</th>
              <th class="text-right text-xs font-semibold text-gray-600 uppercase px-6 py-3">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="u in users" :key="u.id" class="hover:bg-gray-50 transition-colors">
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-full bg-primary-100 text-primary-700 flex items-center justify-center text-sm font-semibold">
                    {{ u.full_name?.charAt(0)?.toUpperCase() }}
                  </div>
                  <div>
                    <p class="text-sm font-medium text-gray-900">{{ u.full_name }}</p>
                    <p class="text-xs text-gray-500" v-if="u.email">{{ u.email }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 text-sm text-gray-600">{{ u.username }}</td>
              <td class="px-6 py-4">
                <span class="badge" :class="roleBadgeClass(u.role)">{{ roleLabel(u.role) }}</span>
              </td>
              <td class="px-6 py-4">
                <span class="badge" :class="u.is_active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-600'">
                  {{ u.is_active ? 'Aktif' : 'Nonaktif' }}
                </span>
              </td>
              <td class="px-6 py-4 text-right">
                <button @click="() => {}" class="text-gray-400 hover:text-gray-600 p-1">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01" />
                  </svg>
                </button>
              </td>
            </tr>
            <tr v-if="users.length === 0">
              <td colspan="5" class="px-6 py-12 text-center text-gray-500 text-sm">Tidak ada data user</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create User Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50" @click.self="showCreateModal = false">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md animate-fade-in">
        <div class="flex items-center justify-between px-6 py-4 border-b border-gray-100">
          <h3 class="text-lg font-semibold text-gray-900">Tambah User Baru</h3>
          <button @click="showCreateModal = false" class="text-gray-400 hover:text-gray-600">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <form @submit.prevent="createUser" class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama Lengkap</label>
            <input v-model="form.full_name" class="input-field" required />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
            <input v-model="form.username" class="input-field" required />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
            <input v-model="form.password" type="password" class="input-field" required minlength="6" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Role</label>
            <select v-model="form.role" class="input-field" required>
              <option value="TEACHER">Guru/Staff</option>
              <option value="HEAD">Kepala Sekolah</option>
              <option value="ADMIN">Administrator</option>
            </select>
          </div>
          <div class="flex gap-3 pt-2">
            <button type="button" @click="showCreateModal = false" class="btn-secondary flex-1">Batal</button>
            <button type="submit" :disabled="creating" class="btn-primary flex-1">
              {{ creating ? 'Menyimpan...' : 'Simpan' }}
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
const showCreateModal = ref(false)
const creating = ref(false)
const form = ref({ full_name: '', username: '', password: '', role: 'TEACHER' })

const roleLabels = { ADMIN: 'Admin', TEACHER: 'Guru', HEAD: 'Kepsek' }
const roleLabel = (role) => roleLabels[role] || role

const roleBadgeClass = (role) => ({
  ADMIN: 'bg-purple-100 text-purple-800',
  TEACHER: 'bg-blue-100 text-blue-800',
  HEAD: 'bg-amber-100 text-amber-800',
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

async function createUser() {
  creating.value = true
  try {
    await api.post('/users', form.value)
    showCreateModal.value = false
    form.value = { full_name: '', username: '', password: '', role: 'TEACHER' }
    await fetchUsers()
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal membuat user')
  } finally {
    creating.value = false
  }
}

onMounted(fetchUsers)
</script>
