import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../utils/api'

export const useItemStore = defineStore('items', () => {
  const items = ref([])
  const currentItem = ref(null)
  const meta = ref({
    page: 1,
    page_size: 20,
    total: 0,
    total_pages: 0
  })
  const loading = ref(false)
  const error = ref(null)

  async function fetchItems(filters = {}) {
    loading.value = true
    error.value = null
    try {
      const params = new URLSearchParams()
      for (const [key, value] of Object.entries(filters)) {
        if (value) params.append(key, value)
      }
      
      const { data } = await api.get(`/items?${params.toString()}`)
      if (data.success) {
        items.value = data.data.items || []
        meta.value = data.data.meta
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal mengambil data barang'
      items.value = []
    } finally {
      loading.value = false
    }
  }

  async function fetchItem(id) {
    loading.value = true
    error.value = null
    try {
      const { data } = await api.get(`/items/${id}`)
      if (data.success) {
        currentItem.value = data.data
        return data.data
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal mengambil detail barang'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createItem(itemData) {
    loading.value = true
    error.value = null
    try {
      // Clean up empty strings to null or omit
      const payload = { ...itemData }
      if (payload.category_id === '') delete payload.category_id
      if (payload.purchase_price === '') delete payload.purchase_price
      
      const { data } = await api.post('/items', payload)
      if (data.success) {
        return data.data
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal menambah barang'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateItem(id, itemData) {
    loading.value = true
    error.value = null
    try {
      const payload = { ...itemData }
      if (payload.category_id === '') payload.category_id = null
      if (payload.purchase_price === '') payload.purchase_price = null
      
      const { data } = await api.put(`/items/${id}`, payload)
      if (data.success) {
        return true
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal mengupdate barang'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchItemHistory(id) {
    loading.value = true
    error.value = null
    try {
      const { data } = await api.get(`/items/${id}/history`)
      if (data.success) {
        return data.data
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal mengambil riwayat barang'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteItem(id) {
    loading.value = true
    error.value = null
    try {
      const { data } = await api.delete(`/items/${id}`)
      if (data.success) {
        // Remove from local array if exists
        items.value = items.value.filter(item => item.id !== id)
        return true
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal menghapus barang'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchItemByCode(code) {
    loading.value = true
    error.value = null
    try {
      const { data } = await api.get(`/items/code/${code}`)
      if (data.success) {
        return data.data
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal mengambil detail barang'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    items,
    currentItem,
    meta,
    loading,
    error,
    fetchItems,
    fetchItem,
    fetchItemByCode,
    fetchItemHistory,
    createItem,
    updateItem,
    deleteItem
  }
})
