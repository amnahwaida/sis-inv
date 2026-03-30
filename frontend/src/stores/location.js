import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../utils/api'

export const useLocationStore = defineStore('locations', () => {
  const locations = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function fetchLocations() {
    loading.value = true
    error.value = null
    try {
      const { data } = await api.get('/locations')
      if (data.success) {
        locations.value = data.data || []
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal mengambil data lokasi'
      locations.value = []
    } finally {
      loading.value = false
    }
  }

  async function createLocation(locData) {
    loading.value = true
    error.value = null
    try {
      const { data } = await api.post('/locations', locData)
      if (data.success) {
        await fetchLocations()
        return data.data
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal menambah lokasi'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteLocation(id) {
    loading.value = true
    error.value = null
    try {
      const { data } = await api.delete(`/locations/${id}`)
      if (data.success) {
        await fetchLocations()
        return true
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal menghapus lokasi'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    locations,
    loading,
    error,
    fetchLocations,
    createLocation,
    deleteLocation
  }
})
