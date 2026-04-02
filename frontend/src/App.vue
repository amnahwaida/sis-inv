<template>
  <router-view />
</template>

<script setup>
import { watch } from 'vue'
import { useSettingsStore } from './stores/settings'

const settingsStore = useSettingsStore()

// Watch for app name changes to update document title
watch(() => settingsStore.settings.app_name, (newName) => {
  if (newName) {
    document.title = `${newName} | Sistem Inventaris Sekolah`
  } else {
    document.title = 'SIS-INV | Sistem Inventaris Sekolah'
  }
}, { immediate: true })

// Watch for favicon changes to update the browser tab icon
watch(() => settingsStore.settings.app_favicon_url, (url) => {
  let link = document.querySelector("link[rel*='icon']")
  if (!link) {
    link = document.createElement('link')
    link.rel = 'icon'
    document.getElementsByTagName('head')[0].appendChild(link)
  }
  
  // If we have a custom URL, use it. Otherwise fallback to default /favicon.svg
  link.href = url || '/favicon.svg'
}, { immediate: true })

// Watch for theme changes
watch(() => settingsStore.settings.app_theme, (theme) => {
  if (theme) {
    document.documentElement.setAttribute('data-theme', theme)
  } else {
    document.documentElement.removeAttribute('data-theme')
  }
}, { immediate: true })
</script>
