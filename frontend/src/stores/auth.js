import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../utils/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const accessToken = ref(localStorage.getItem('access_token') || '')

  const isAuthenticated = computed(() => !!accessToken.value)
  const userRole = computed(() => user.value?.role || '')
  const isAdmin = computed(() => userRole.value === 'ADMIN')
  const isTeacher = computed(() => userRole.value === 'TEACHER')
  const isHead = computed(() => userRole.value === 'HEAD')

  async function login(username, password) {
    const { data } = await api.post('/auth/login', { username, password })
    if (data.success) {
      accessToken.value = data.data.access_token
      user.value = data.data.user
      localStorage.setItem('access_token', data.data.access_token)
      localStorage.setItem('refresh_token', data.data.refresh_token)
      localStorage.setItem('user', JSON.stringify(data.data.user))
      return data.data
    }
    throw new Error(data.error || 'Login gagal')
  }

  async function fetchMe() {
    try {
      const { data } = await api.get('/auth/me')
      if (data.success) {
        user.value = data.data
        localStorage.setItem('user', JSON.stringify(data.data))
      }
    } catch {
      logout()
    }
  }

  function logout() {
    accessToken.value = ''
    user.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
  }

  return {
    user,
    accessToken,
    isAuthenticated,
    userRole,
    isAdmin,
    isTeacher,
    isHead,
    login,
    fetchMe,
    logout,
  }
})
