import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../utils/api'

export const useAuditStore = defineStore('audit', () => {
  const sessions = ref([])
  const currentSession = ref(null)
  const loading = ref(false)
  const error = ref(null)

  async function fetchSessions() {
    loading.value = true
    try {
      const { data } = await api.get('/reports/audit-sessions')
      sessions.value = data.data || []
    } catch (err) {
      error.value = 'Gagal mengambil daftar audit'
    } finally {
      loading.value = false
    }
  }

  async function startSession(locationId, notes) {
    loading.value = true
    try {
      const { data } = await api.post('/reports/audit-sessions', { location_id: locationId, notes })
      await fetchSessions()
      return data.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal memulai audit'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchSessionDetail(id) {
    loading.value = true
    try {
      const { data } = await api.get(`/reports/audit-sessions/${id}`)
      currentSession.value = data.data
    } catch (err) {
      error.value = 'Gagal mengambil detail audit'
    } finally {
      loading.value = false
    }
  }

  async function scanAuditItem(sessionId, itemCode, condition, notes) {
    try {
      const { data } = await api.post(`/reports/audit-sessions/${sessionId}/scan`, {
        item_code: itemCode,
        condition: condition,
        notes: notes
      })
      return data.success
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal scan barang'
      throw err
    }
  }

  async function closeSession(id) {
    loading.value = true
    try {
      await api.post(`/reports/audit-sessions/${id}/close`)
      await fetchSessions()
    } catch (err) {
      error.value = 'Gagal menutup sesi audit'
    } finally {
      loading.value = false
    }
  }

  return {
    sessions,
    currentSession,
    loading,
    error,
    fetchSessions,
    startSession,
    fetchSessionDetail,
    scanAuditItem,
    closeSession
  }
})
