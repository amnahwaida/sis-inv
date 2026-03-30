import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../utils/api'

export const useTransactionStore = defineStore('transactions', () => {
  const myBorrowings = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function fetchMyBorrowings() {
    loading.value = true
    error.value = null
    try {
      const { data } = await api.get('/transactions/my')
      if (data.success) {
        myBorrowings.value = data.data || []
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal mengambil data peminjaman'
      myBorrowings.value = []
    } finally {
      loading.value = false
    }
  }

  /**
   * 
   * @param {Object} borrowData
   * {
   *   item_code: string,
   *   borrower_type: 'STAFF' | 'STUDENT',
   *   expected_return_days: number,
   *   purpose: string,
   *   student_name?: string,
   *   student_nis?: string,
   *   student_class?: string
   * }
   */
  async function borrowItem(borrowData) {
    loading.value = true
    error.value = null
    try {
      const { data } = await api.post('/transactions/borrow', borrowData)
      if (data.success) {
        return data.data
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal meminjam barang'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function returnItem(itemCode, condition, notes, photoUrl) {
    loading.value = true
    error.value = null
    try {
      const payload = {
        item_code: itemCode,
        condition: condition || 'GOOD',
        notes: notes,
        photo_url: photoUrl
      }
      const { data } = await api.post('/transactions/return', payload)
      if (data.success) {
        return data.data
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal mengembalikan barang'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    myBorrowings,
    loading,
    error,
    fetchMyBorrowings,
    borrowItem,
    returnItem
  }
})
