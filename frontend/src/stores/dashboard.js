import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../utils/api'

export const useDashboardStore = defineStore('dashboard', () => {
  const summary = ref({
    total_items: 0,
    items_by_status: { AVAILABLE: 0, BORROWED: 0, MAINTENANCE: 0, LOST: 0 },
    items_by_condition: { GOOD: 0, DAMAGED: 0, IN_REPAIR: 0, LOST: 0 },
    overdue_count: 0
  })
  
  const loading = ref(false)
  const error = ref(null)

  async function fetchSummary() {
    loading.value = true
    error.value = null
    try {
      const { data } = await api.get('/dashboard/summary')
      if (data.success) {
        summary.value = data.data
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal mengambil data dashboard'
    } finally {
      loading.value = false
    }
  }

  return {
    summary,
    loading,
    error,
    fetchSummary
  }
})
