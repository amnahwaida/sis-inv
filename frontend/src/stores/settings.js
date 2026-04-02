import { defineStore } from 'pinia'
import api from '../utils/api'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    settings: {
      app_name: 'SIS-INV',
      app_subtitle: 'Inventaris Sekolah',
      app_description: 'School Inventory System',
      app_footer: 'SIS-INV • AMANAH & TERTIB',
      app_security_notice: 'Terlindungi oleh Enkripsi End-to-End'
    },
    loading: false
  }),
  actions: {
    async fetchSettings() {
      this.loading = true
      try {
        const { data } = await api.get('/settings')
        if (data.success) {
          // Merge fetched settings into defaults
          this.settings = { ...this.settings, ...data.data }
        }
      } catch (e) {
        console.error('Failed to fetch settings:', e)
      } finally {
        this.loading = false
      }
    },
    async updateSettings(newSettings) {
      try {
        const { data } = await api.put('/settings', newSettings)
        if (data.success) {
          this.settings = { ...this.settings, ...newSettings }
        }
        return data.success
      } catch (e) {
        console.error('Failed to update settings:', e)
        throw e
      }
    }
  }
})
