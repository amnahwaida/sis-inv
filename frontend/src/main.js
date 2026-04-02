import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import './assets/main.css'
import { registerSW } from 'virtual:pwa-register'

registerSW({ immediate: true })

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)
app.use(router)

// Import stores AFTER pinia is installed
import { useSettingsStore } from './stores/settings'
import { usePwaStore } from './stores/pwa'

// Fetch branding settings globally
const settingsStore = useSettingsStore()
settingsStore.fetchSettings()

// Capture beforeinstallprompt globally so it's never lost
const pwaStore = usePwaStore()
window.addEventListener('beforeinstallprompt', (e) => {
  e.preventDefault()
  pwaStore.setInstallPrompt(e)
})

app.mount('#app')
