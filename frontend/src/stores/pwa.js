import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const usePwaStore = defineStore('pwa', () => {
  const installPrompt = ref(null)
  const showIOSGuide = ref(false)

  const isIOS = computed(() => {
    return [
      'iPad Simulator',
      'iPhone Simulator',
      'iPod Simulator',
      'iPad',
      'iPhone',
      'iPod'
    ].includes(navigator.platform)
    || (navigator.userAgent.includes("Mac") && "ontouchend" in document)
  })

  const isStandalone = computed(() => {
    return window.matchMedia('(display-mode: standalone)').matches || window.navigator.standalone
  })

  const canInstall = computed(() => {
    return !!installPrompt.value
  })

  function setInstallPrompt(event) {
    installPrompt.value = event
  }

  async function triggerInstall() {
    if (!installPrompt.value) return
    installPrompt.value.prompt()
    const { outcome } = await installPrompt.value.userChoice
    if (outcome === 'accepted') {
      console.log('User accepted the install prompt')
    }
    installPrompt.value = null
  }

  return {
    installPrompt,
    showIOSGuide,
    isIOS,
    isStandalone,
    canInstall,
    setInstallPrompt,
    triggerInstall
  }
})
