import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import './assets/main.css'
import { registerSW } from 'virtual:pwa-register'

registerSW({ immediate: true })

import { useSettingsStore } from './stores/settings'

const app = createApp(App)
app.use(createPinia())
app.use(router)

// Fetch branding settings globally
const settingsStore = useSettingsStore()
settingsStore.fetchSettings()

app.mount('#app')
