import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../utils/api'

export const useCategoryStore = defineStore('categories', () => {
  const categories = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function fetchCategories() {
    loading.value = true
    error.value = null
    try {
      const { data } = await api.get('/categories?page_size=1000')
      if (data.success) {
        categories.value = data.data?.items || data.data || []
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal mengambil data kategori'
      categories.value = []
    } finally {
      loading.value = false
    }
  }

  async function createCategory(catData) {
    loading.value = true
    error.value = null
    try {
      const { data } = await api.post('/categories', catData)
      if (data.success) {
        await fetchCategories()
        return data.data
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal menambah kategori'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteCategory(id) {
    loading.value = true
    error.value = null
    try {
      const { data } = await api.delete(`/categories/${id}`)
      if (data.success) {
        await fetchCategories()
        return true
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal menghapus kategori'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    categories,
    loading,
    error,
    fetchCategories,
    createCategory,
    deleteCategory
  }
})
