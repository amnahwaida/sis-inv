<template>
  <div class="animate-fade-in space-y-10 max-w-7xl mx-auto pb-20">
    <!-- Header Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-gray-900 to-primary-900 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-primary-900/20 transition-all duration-500">
      <div class="absolute top-0 right-0 -mt-12 -mr-12 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -mb-12 -ml-12 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl"></div>
      
      <div class="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <h1 class="text-3xl font-black tracking-tight leading-none">Pengaturan Sistem</h1>
          <p class="text-primary-100/70 text-sm font-medium">Atur profil, preferensi, and keamanan akun Anda</p>
        </div>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
      
      <!-- Left Column: Navigation Tabs -->
      <div class="lg:col-span-4 space-y-4">
        <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] p-4 border border-gray-100 dark:border-gray-700 shadow-sm">
          <div class="flex flex-col gap-2">
            <button v-for="tab in [
                      {id:'profile', label:'Profil Pengguna', icon: 'M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z'}, 
                      {id:'security', label:'Keamanan & Password', icon: 'M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6.248c-.063.63-.098 1.27-.098 1.918 0 5.617 3.557 10.413 8.5 12.897 4.943-2.484 8.5-7.28 8.5-12.897 0-.648-.035-1.288-.098-1.918A11.959 11.959 0 0112 2.714z'}, 
                      {id:'app', label:'Preferensi Aplikasi', icon: 'M10.5 1.5H8.25A2.25 2.25 0 006 3.75v16.5a2.25 2.25 0 002.25 2.25h7.5A2.25 2.25 0 0018 20.25V3.75a2.25 2.25 0 00-2.25-2.25H13.5m-3 0V3h3V1.5m-3 0h3m-3 18.75h3'},
                      {id:'branding', label:'Kustomisasi Branding', icon: 'M9.53 16.122a3 3 0 00-2.012.324l-1.734.917a.75.75 0 01-1.071-.742V15.61a3 3 0 00-.734-1.939L2.527 12.003a.75.75 0 010-1.002l1.452-1.668a3 3 0 00.734-1.938V2.392a.75.75 0 011.071-.742l1.734.917a3 3 0 002.012.324h7.456a.75.75 0 01.594.288l2.91 3.556a3 3 0 001.077.892l2.368 1.134a.75.75 0 01.354.912l-.916 2.38a3 3 0 000 2.164l.916 2.38a.75.75 0 01-.354.912l-2.368 1.134a3 3 0 00-1.077.892l-2.91 3.556a.75.75 0 01-.594.288H9.53z'},
                      {id:'system', label:'Data & Sistem', icon: 'M4.5 12a7.5 7.5 0 0015 0m-15 0a7.5 7.5 0 1115 0m-15 0H3m16.5 0H21m-1.5 0a5.999 5.999 0 00-5.999-5.999 5.999 5.999 0 00-5.999 5.999 5.999 5.999 0 005.999 5.999 5.999 5.999 0 005.999-5.999z'}
                    ].filter(t => !t.adminOnly || authStore.userRole === 'ADMIN')" 
                    :key="tab.id"
                    @click="activeTab = tab.id"
                    :class="activeTab === tab.id ? 'bg-primary-600 text-white shadow-lg' : 'text-gray-500 hover:bg-gray-50 dark:hover:bg-gray-700 dark:text-gray-400'"
                    class="px-6 py-4 rounded-2xl text-xs font-black transition-all flex items-center justify-between group active:scale-95">
              <span class="flex items-center gap-3">
                <span class="w-5 h-5 flex items-center justify-center">
                  <svg class="w-full h-full" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" :d="tab.icon" />
                  </svg>
                </span>
                <span class="uppercase tracking-widest font-black leading-none">{{ tab.label }}</span>
              </span>
              <svg class="w-4 h-4 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Help Info -->
        <div class="p-8 bg-indigo-50 dark:bg-indigo-900/20 rounded-[2.5rem] border border-indigo-100 dark:border-indigo-800 space-y-4">
          <div class="w-10 h-10 bg-indigo-600 text-white rounded-xl flex items-center justify-center font-black">
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9 5.25h.008v.008H12v-.008z" />
            </svg>
          </div>
          <h4 class="text-xs font-black text-indigo-900 dark:text-indigo-100 uppercase tracking-widest">Bantuan & Dukungan</h4>
          <p class="text-[10px] font-bold text-indigo-400 leading-relaxed uppercase tracking-tight">Hubungi Administrator jika Anda mengalami kesulitan dalam pengaturan akun atau memiliki pertanyaan teknis sistem.</p>
        </div>
      </div>

      <!-- Right Column: Form Area -->
      <div class="lg:col-span-8">
        <div class="bg-white dark:bg-gray-800 rounded-[2.5rem] border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden animate-scale-up">
          
          <!-- Profile Form -->
          <div v-if="activeTab === 'profile'" class="p-10 space-y-8 animate-fade-in">
            <div class="flex items-center gap-4">
              <div class="w-16 h-16 bg-primary-50 dark:bg-primary-900/30 text-primary-600 rounded-3xl flex items-center justify-center text-xl font-black shadow-inner">
                {{ authStore.user?.full_name?.charAt(0) }}
              </div>
              <div>
                <h3 class="text-xl font-black text-gray-900 dark:text-white leading-none capitalize">Informasi Akun</h3>
                <p class="text-[10px] font-black text-primary-500 uppercase tracking-widest mt-2">Data profil personal Anda</p>
              </div>
            </div>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 pt-6 border-t border-gray-50 dark:border-gray-700">
              <div class="space-y-1">
                <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Nama Lengkap</label>
                <input :value="authStore.user?.full_name" class="input-field rounded-2xl h-14 opacity-60" disabled />
              </div>
              <div class="space-y-1">
                <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Username (Login)</label>
                <input :value="authStore.user?.username" class="input-field rounded-2xl h-14 opacity-60" disabled />
              </div>
              <div class="space-y-1">
                <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Hak Akses (Role)</label>
                <span class="inline-block px-4 py-3 bg-gray-50 dark:bg-gray-900/40 rounded-2xl text-[10px] font-black text-primary-600 uppercase tracking-widest">{{ authStore.user?.role }}</span>
              </div>
            </div>
            
            <p class="text-[10px] text-gray-400 italic bg-gray-50 dark:bg-gray-900/30 p-4 rounded-xl">Nama and username dikunci oleh Sistem. Hubungi admin untuk perubahan biodata.</p>
          </div>

          <!-- Password / Security Form -->
          <div v-if="activeTab === 'security'" class="p-10 space-y-10 animate-fade-in">
            <h3 class="text-xl font-black text-gray-900 dark:text-white capitalize">Keamanan & Password</h3>
            
            <form @submit.prevent="handleChangePassword" class="space-y-6">
              <div class="space-y-1">
                <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Password Lama</label>
                <input type="password" v-model="passForm.old_password" class="input-field rounded-2xl h-14" required />
              </div>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-1">
                  <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Password Baru</label>
                  <input type="password" v-model="passForm.new_password" class="input-field rounded-2xl h-14" required minlength="6" />
                </div>
                <div class="space-y-1">
                  <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Konfirmasi Password Baru</label>
                  <input type="password" v-model="passForm.confirm_password" class="input-field rounded-2xl h-14" required minlength="6" />
                </div>
              </div>

              <div class="flex justify-end pt-6">
                <button type="submit" :disabled="submitting"
                        class="btn-premium-action !px-10">
                  {{ submitting ? 'MEMPROSES...' : 'GANTI PASSWORD' }}
                </button>
              </div>
            </form>
          </div>

          <!-- App Settings Form -->
          <div v-if="activeTab === 'app'" class="p-10 space-y-10 animate-fade-in">
            <h3 class="text-xl font-black text-gray-900 dark:text-white capitalize">Preferensi Aplikasi</h3>
            
            <div class="space-y-6">
              <div class="flex items-center justify-between p-6 bg-gray-50 dark:bg-gray-900/30 rounded-[2rem] border border-gray-100 dark:border-gray-700">
                <div class="space-y-1">
                  <p class="text-sm font-black text-gray-900 dark:text-white uppercase tracking-tight">Mode Gelap (Dark Mode)</p>
                  <p class="text-xs text-gray-400 font-bold uppercase tracking-widest">Atur kenyamanan mata Anda</p>
                </div>
                <button @click="toggleMode" 
                        class="w-14 h-8 bg-gray-200 dark:bg-primary-600 rounded-full relative transition-colors duration-500 shadow-inner">
                  <div class="absolute top-1 left-1 dark:left-7 w-6 h-6 bg-white rounded-full shadow-lg transition-all duration-300 flex items-center justify-center text-primary-600">
                    <svg v-if="!isDark" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M3 12h2.25m.386-6.364l1.591 1.591M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z" />
                    </svg>
                    <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z" />
                    </svg>
                  </div>
                </button>
              </div>
            </div>
          </div>

          <!-- Branding Customization Form (Admin Only) -->
          <div v-if="activeTab === 'branding' && authStore.userRole === 'ADMIN'" class="p-10 space-y-10 animate-fade-in text-gray-900 dark:text-white">
            <div>
              <h3 class="text-xl font-black capitalize">Kustomisasi Branding</h3>
              <p class="text-[10px] font-black text-primary-500 uppercase tracking-widest mt-2">Ubah identitas visual teks aplikasi</p>
            </div>
            
            <form @submit.prevent="handleSaveBranding" class="space-y-8">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-1">
                  <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Nama Aplikasi</label>
                  <input v-model="brandingForm.app_name" class="input-field rounded-2xl h-14" placeholder="SIS-INV" required />
                </div>
                <div class="space-y-1">
                  <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Subtitle Aplikasi</label>
                  <input v-model="brandingForm.app_subtitle" class="input-field rounded-2xl h-14" placeholder="Inventaris Sekolah" required />
                </div>
              </div>
              
              <div class="space-y-1">
                <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Deskripsi (Login Page)</label>
                <input v-model="brandingForm.app_description" class="input-field rounded-2xl h-14" placeholder="School Inventory System" required />
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-1">
                  <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Teks Footer (Login Page)</label>
                  <input v-model="brandingForm.app_footer" class="input-field rounded-2xl h-14" placeholder="SIS-INV • AMANAH & TERTIB" required />
                </div>
                <div class="space-y-1">
                  <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] ml-1">Credit</label>
                  <input v-model="brandingForm.app_security_notice" class="input-field rounded-2xl h-14" placeholder="Crafted with ❤️ by vannyezhaa..." required />
                </div>
              </div>

              <!-- Theme Selection -->
              <div class="pt-6 border-t border-gray-100 dark:border-gray-700 space-y-6">
                <div class="flex items-center justify-between">
                  <h4 class="text-xs font-black uppercase tracking-widest text-primary-600">Tone Warna Sistem</h4>
                  <span class="text-[10px] text-gray-400 font-bold uppercase tracking-tighter">Pilih nuansa warna aplikasi Anda</span>
                </div>
                
                <div class="grid grid-cols-2 md:grid-cols-5 gap-4">
                  <button v-for="theme in [
                            {id: 'blue', label: 'Sapphire Blue', color: 'bg-blue-500'},
                            {id: 'green', label: 'Emerald Green', color: 'bg-emerald-600'},
                            {id: 'purple', label: 'Royal Purple', color: 'bg-purple-600'},
                            {id: 'red', label: 'Crimson Red', color: 'bg-red-600'},
                            {id: 'orange', label: 'Amber Gold', color: 'bg-amber-500'}
                          ]" 
                          :key="theme.id"
                          type="button"
                          @click="brandingForm.app_theme = theme.id"
                          :class="brandingForm.app_theme === theme.id ? 'border-primary-500 ring-4 ring-primary-500/20 bg-gray-50 dark:bg-gray-700' : 'border-gray-100 dark:border-gray-700 hover:border-primary-300 dark:hover:border-gray-600'"
                          class="flex items-center gap-3 p-4 rounded-3xl border-2 transition-all group active:scale-95">
                    <div :class="theme.color" class="w-8 h-8 rounded-2xl shadow-lg transform group-hover:scale-110 transition-transform"></div>
                    <div class="text-left">
                      <p class="text-[10px] font-black uppercase tracking-tight" :class="brandingForm.app_theme === theme.id ? 'text-primary-600' : 'text-gray-500'">{{ theme.label }}</p>
                    </div>
                  </button>
                </div>
              </div>

              <!-- Logo Customization -->
              <div class="pt-6 border-t border-gray-100 dark:border-gray-700 space-y-6">
                <div class="flex items-center justify-between">
                  <h4 class="text-xs font-black uppercase tracking-widest text-primary-600">Logo Utama Aplikasi</h4>
                  <span class="text-[10px] text-gray-400 font-bold uppercase tracking-tighter">Support PNG, SVG, JPG (Max 2MB)</span>
                </div>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
                  <!-- Sidebar Logo -->
                  <div class="space-y-4">
                    <label class="block text-[10px] font-black text-gray-400 uppercase tracking-widest ml-1">Logo Sidebar</label>
                    <div class="flex flex-col gap-4 p-6 bg-gray-50 dark:bg-gray-800/50 rounded-3xl border-2 border-dashed border-gray-200 dark:border-gray-700 hover:border-primary-500 transition-colors group">
                      <div class="w-20 h-20 bg-white dark:bg-gray-800 rounded-2xl flex items-center justify-center shadow-xl overflow-hidden border border-gray-100 dark:border-gray-700 mx-auto text-primary-500">
                        <img v-if="brandingForm.app_logo_url" :src="brandingForm.app_logo_url" class="w-full h-full object-cover" />
                        <svg v-else class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 21h19.5m-18-18v18m10.5-18v18m6-13.5V21M6.75 6.75h.75m-.75 3h.75m-.75 3h.75m3-3h.75m-.75 3h.75m-.75 3h.75M6.75 21v-3.375c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21M3 3h12m-.75 4.5H21m-3.75 3.75h.008v.008h-.008v-.008zm0 3h.008v.008h-.008v-.008zm0 3h.008v.008h-.008v-.008z" />
                        </svg>
                      </div>
                      <div class="space-y-2">
                        <input type="file" ref="appLogoInput" class="hidden" accept="image/*" @change="e => handleLogoUpload(e, 'app_logo_url')" />
                        <button type="button" @click="appLogoInput.click()" 
                                class="w-full py-2.5 px-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl text-[10px] font-black uppercase tracking-widest hover:border-primary-500 transition-all active:scale-95">
                          Pilih Logo
                        </button>
                        <button v-if="brandingForm.app_logo_url" type="button" @click="brandingForm.app_logo_url = ''" 
                                class="w-full text-[9px] font-black text-red-500 uppercase tracking-tight">Hapus</button>
                      </div>
                    </div>
                  </div>
                  <!-- Login Logo -->
                  <div class="space-y-4">
                    <label class="block text-[10px] font-black text-gray-400 uppercase tracking-widest ml-1">Logo Login</label>
                    <div class="flex flex-col gap-4 p-6 bg-gray-50 dark:bg-gray-800/50 rounded-3xl border-2 border-dashed border-gray-200 dark:border-gray-700 hover:border-primary-500 transition-colors group">
                      <div class="w-20 h-20 bg-white dark:bg-gray-800 rounded-2xl flex items-center justify-center shadow-xl overflow-hidden border border-gray-100 dark:border-gray-700 mx-auto text-primary-500">
                        <img v-if="brandingForm.login_logo_url" :src="brandingForm.login_logo_url" class="w-full h-full object-cover" />
                        <svg v-else class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z" />
                        </svg>
                      </div>
                      <div class="space-y-2">
                        <input type="file" ref="loginLogoInput" class="hidden" accept="image/*" @change="e => handleLogoUpload(e, 'login_logo_url')" />
                        <button type="button" @click="loginLogoInput.click()" 
                                class="w-full py-2.5 px-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl text-[10px] font-black uppercase tracking-widest hover:border-primary-500 transition-all active:scale-95">
                          Pilih Logo
                        </button>
                        <button v-if="brandingForm.login_logo_url" type="button" @click="brandingForm.login_logo_url = ''" 
                                class="w-full text-[9px] font-black text-red-500 uppercase tracking-tight">Hapus</button>
                      </div>
                    </div>
                  </div>
                  <!-- Favicon -->
                  <div class="space-y-4">
                    <label class="block text-[10px] font-black text-gray-400 uppercase tracking-widest ml-1">Favicon (Tab Browser)</label>
                    <div class="flex flex-col gap-4 p-6 bg-gray-50 dark:bg-gray-800/50 rounded-3xl border-2 border-dashed border-gray-200 dark:border-gray-700 hover:border-primary-500 transition-colors group">
                      <div class="w-20 h-20 bg-white dark:bg-gray-800 rounded-2xl flex items-center justify-center shadow-xl overflow-hidden border border-gray-100 dark:border-gray-700 mx-auto text-primary-500">
                        <img v-if="brandingForm.app_favicon_url" :src="brandingForm.app_favicon_url" class="w-full h-full object-contain p-2" />
                        <svg v-else class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M12 21a9.004 9.004 0 008.716-6.747M12 21a9.004 9.004 0 01-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 017.843 4.582M12 3a8.997 8.997 0 00-7.843 4.582m15.686 0A11.953 11.953 0 0112 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0121 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0112 16.5c-3.162 0-6.133-.815-8.716-2.247m0 0A9.015 9.015 0 013 12c0-.856.12-1.683.342-2.463" />
                        </svg>
                      </div>
                      <div class="space-y-2">
                        <input type="file" ref="faviconInput" class="hidden" accept="image/*" @change="e => handleLogoUpload(e, 'app_favicon_url')" />
                        <button type="button" @click="faviconInput.click()" 
                                class="w-full py-2.5 px-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl text-[10px] font-black uppercase tracking-widest hover:border-primary-500 transition-all active:scale-95">
                          Pilih Ikon
                        </button>
                        <button v-if="brandingForm.app_favicon_url" type="button" @click="brandingForm.app_favicon_url = ''" 
                                class="w-full text-[9px] font-black text-red-500 uppercase tracking-tight">Hapus</button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Menu Icons Customization -->
              <div class="pt-6 border-t border-gray-100 dark:border-gray-700 space-y-6">
                <div class="flex items-center justify-between">
                  <h4 class="text-xs font-black uppercase tracking-widest text-primary-600">Ikon Menu (Heroicons / Emoji)</h4>
                  <span class="text-[10px] text-gray-400 font-bold uppercase tracking-tighter">Gunakan kode SVG atau Emoji</span>
                </div>

                <!-- Customization Note -->
                <div class="bg-primary-50 dark:bg-primary-900/10 rounded-2xl p-6 border border-primary-100 dark:border-primary-900/30 space-y-3">
                  <div class="flex items-center gap-2 text-primary-600">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z" />
                    </svg>
                    <span class="text-[10px] font-black uppercase tracking-widest">Cara Kustom dengan Heroicons:</span>
                  </div>
                  <ol class="text-[11px] text-gray-600 dark:text-gray-400 space-y-1.5 list-decimal ml-4">
                    <li>Buka <a href="https://heroicons.com" target="_blank" class="text-primary-600 font-bold border-b border-primary-200 hover:border-primary-600 transition-colors">heroicons.com</a></li>
                    <li>Cari ikon yang sesuai, klik kanan dan pilih <strong>"Copy SVG"</strong>.</li>
                    <li><strong>Tempelkan (Paste)</strong> kode yang sudah di-copy tadi ke kolom input menu pilihan Anda.</li>
                  </ol>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  <div v-for="icon in [
                    {key: 'icon_dashboard', label: 'Dashboard'},
                    {key: 'icon_items', label: 'Daftar Barang'},
                    {key: 'icon_borrow', label: 'Pinjaman Baru'},
                    {key: 'icon_return', label: 'Pengembalian'},
                    {key: 'icon_my_borrowings', label: 'Pinjaman Saya'},
                    {key: 'icon_student_history', label: 'Cek Riwayat'},
                    {key: 'icon_reports', label: 'Laporan Aset'},
                    {key: 'icon_users', label: 'Kelola User'},
                    {key: 'icon_students', label: 'Kelola Siswa'},
                    {key: 'icon_categories', label: 'Kategori'},
                    {key: 'icon_locations', label: 'Lokasi Stok'},
                    {key: 'icon_maintenance', label: 'Perbaikan'},
                    {key: 'icon_stock_opname', label: 'Stock Opname'},
                    {key: 'icon_audit_logs', label: 'Log Audit'},
                    {key: 'icon_settings', label: 'Pengaturan'}
                  ]" :key="icon.key" class="space-y-1.5">
                    <label class="block text-[9px] font-black text-gray-400 uppercase tracking-tight ml-1">{{ icon.label }}</label>
                    <div class="relative group">
                      <!-- Live Preview (SVG or Emoji) -->
                      <div class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 flex items-center justify-center pointer-events-none transition-transform group-focus-within:scale-110">
                        <template v-if="brandingForm[icon.key]">
                          <div v-if="brandingForm[icon.key].trim().startsWith('<svg')" 
                               v-html="brandingForm[icon.key]" 
                               class="w-full h-full text-primary-600 [&>svg]:w-full [&>svg]:h-full" />
                          <span v-else class="text-lg leading-none">{{ brandingForm[icon.key] }}</span>
                        </template>
                        <span v-else class="text-gray-300">➖</span>
                      </div>
                      <input v-model="brandingForm[icon.key]" 
                             class="input-field rounded-xl h-11 text-xs pl-10 border-gray-200 dark:border-gray-700 focus:border-primary-500 transition-all" 
                             :placeholder="'Kode SVG atau Emoji...'" />
                    </div>
                  </div>
                </div>
              </div>

              <!-- PWA Customization -->
              <div class="pt-6 border-t border-gray-100 dark:border-gray-700 space-y-6">
                <div class="flex items-center justify-between">
                  <h4 class="text-xs font-black uppercase tracking-widest text-primary-600">Kustomisasi PWA (Progressive Web App)</h4>
                  <span class="text-[10px] text-gray-400 font-bold uppercase tracking-tighter">Ikon untuk instalasi di Home Screen</span>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                  <!-- PWA Icon 192 -->
                  <div class="space-y-4">
                    <label class="block text-[10px] font-black text-gray-400 uppercase tracking-widest ml-1">Ikon PWA (192x192)</label>
                    <div class="flex flex-col gap-4 p-6 bg-gray-50 dark:bg-gray-800/50 rounded-3xl border-2 border-dashed border-gray-200 dark:border-gray-700 hover:border-primary-500 transition-colors group">
                      <div class="w-20 h-20 bg-white dark:bg-gray-800 rounded-2xl flex items-center justify-center shadow-xl overflow-hidden border border-gray-100 dark:border-gray-700 mx-auto text-primary-500">
                        <img v-if="brandingForm.app_pwa_icon_192" :src="brandingForm.app_pwa_icon_192" class="w-full h-full object-contain p-2" />
                        <svg v-else class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 1.5H8.25A2.25 2.25 0 006 3.75v16.5a2.25 2.25 0 002.25 2.25h7.5A2.25 2.25 0 0018 20.25V3.75a2.25 2.25 0 00-2.25-2.25H13.5m-3 0V3h3V1.5m-3 0h3m-3 18.75h3" />
                        </svg>
                      </div>
                      <div class="space-y-2">
                        <input type="file" ref="pwa192Input" class="hidden" accept="image/png" @change="e => handleLogoUpload(e, 'app_pwa_icon_192')" />
                        <button type="button" @click="pwa192Input.click()" 
                                class="w-full py-2.5 px-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl text-[10px] font-black uppercase tracking-widest hover:border-primary-500 transition-all active:scale-95">
                          Pilih Ikon 192px
                        </button>
                      </div>
                    </div>
                  </div>
                  <!-- PWA Icon 512 -->
                  <div class="space-y-4">
                    <label class="block text-[10px] font-black text-gray-400 uppercase tracking-widest ml-1">Ikon PWA (512x512)</label>
                    <div class="flex flex-col gap-4 p-6 bg-gray-50 dark:bg-gray-800/50 rounded-3xl border-2 border-dashed border-gray-200 dark:border-gray-700 hover:border-primary-500 transition-colors group">
                      <div class="w-20 h-20 bg-white dark:bg-gray-800 rounded-2xl flex items-center justify-center shadow-xl overflow-hidden border border-gray-100 dark:border-gray-700 mx-auto text-primary-500">
                        <img v-if="brandingForm.app_pwa_icon_512" :src="brandingForm.app_pwa_icon_512" class="w-full h-full object-contain p-2" />
                        <svg v-else class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 1.5H8.25A2.25 2.25 0 006 3.75v16.5a2.25 2.25 0 002.25 2.25h7.5A2.25 2.25 0 0018 20.25V3.75a2.25 2.25 0 00-2.25-2.25H13.5m-3 0V3h3V1.5m-3 0h3m-3 18.75h3" />
                        </svg>
                      </div>
                      <div class="space-y-2">
                        <input type="file" ref="pwa512Input" class="hidden" accept="image/png" @change="e => handleLogoUpload(e, 'app_pwa_icon_512')" />
                        <button type="button" @click="pwa512Input.click()" 
                                class="w-full py-2.5 px-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl text-[10px] font-black uppercase tracking-widest hover:border-primary-500 transition-all active:scale-95">
                          Pilih Ikon 512px
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <div class="flex flex-col md:flex-row gap-4 justify-between pt-10 border-t border-gray-100 dark:border-gray-700">
                <button type="button" @click="handleResetBranding" :disabled="submittingBranding"
                        class="px-6 py-4 rounded-2xl text-[10px] font-black uppercase tracking-widest border border-red-100 dark:border-red-900/30 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-all active:scale-95 disabled:opacity-50">
                   RESET BRANDING KE AWAL
                </button>
                <button type="submit" :disabled="submittingBranding"
                        class="btn-premium-action !px-10">
                  {{ submittingBranding ? 'MENYIMPAN...' : 'SIMPAN SEMUA PERUBAHAN' }}
                </button>
              </div>
            </form>
          </div>

          <!-- Data & System Management (Admin Only) -->
          <div v-if="activeTab === 'system' && authStore.userRole === 'ADMIN'" class="p-10 space-y-10 animate-fade-in text-gray-900 dark:text-white">
            <div>
              <h3 class="text-xl font-black capitalize">Manajemen Data & Sistem</h3>
              <p class="text-[10px] font-black text-primary-500 uppercase tracking-widest mt-2">Backup, Restore, dan Reset Seluruh Data</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Backup Card -->
              <div class="p-8 bg-blue-50 dark:bg-blue-900/20 rounded-[2.5rem] border border-blue-100 dark:border-blue-800 space-y-4">
                <div class="w-12 h-12 bg-blue-600 text-white rounded-2xl flex items-center justify-center text-xl shadow-lg">
                  <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M17.5 14.5v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0111.5 7.375v-1.5a3.375 3.375 0 00-3.375-3.375H5.75m0 11v6.75c0 .414.336.75.75.75h10.5a.75.75 0 00.75-.75v-6.75m-12 0l-3-3m3 3l3-3m9 3l3-3m-3 3l-3-3" />
                  </svg>
                </div>
                <h4 class="text-sm font-black uppercase tracking-widest">Backup Seluruh Data</h4>
                <p class="text-[10px] font-bold text-blue-400 leading-relaxed uppercase tracking-tight">Unduh seluruh database dan file upload dalam format ZIP untuk cadangan eksternal.</p>
                <button @click="handleBackup" :disabled="isProcessing" class="btn-premium-action !bg-blue-600 !px-6 w-full">
                  {{ isProcessing ? 'MEMPROSES...' : 'DOWNLOAD BACKUP .ZIP' }}
                </button>
              </div>

              <!-- Restore Card -->
              <div class="p-8 bg-emerald-50 dark:bg-emerald-900/20 rounded-[2.5rem] border border-emerald-100 dark:border-emerald-800 space-y-4">
                <div class="w-12 h-12 bg-emerald-600 text-white rounded-2xl flex items-center justify-center text-xl shadow-lg">
                  <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 9.75v6.75m0 0l-3-3m3 3l3-3m-8.25 6a4.5 4.5 0 01-1.41-8.775 5.25 5.25 0 0110.233-2.33 3 3 0 013.758 3.848A3.752 3.752 0 0118 19.5H6.75z" />
                  </svg>
                </div>
                <h4 class="text-sm font-black uppercase tracking-widest">Restore dari Backup</h4>
                <p class="text-[10px] font-bold text-emerald-400 leading-relaxed uppercase tracking-tight">Unggah file ZIP backup untuk memulihkan seluruh kondisi sistem. Data saat ini akan terhapus!</p>
                <input type="file" ref="restoreFile" class="hidden" accept=".zip" @change="onFileSelected" />
                <button @click="$refs.restoreFile.click()" :disabled="isProcessing" class="btn-premium-action !bg-emerald-600 !px-6 w-full">
                  {{ isProcessing ? 'SEDANG RESTORE...' : 'PILIH FILE & RESTORE' }}
                </button>
              </div>
            </div>

            <!-- Danger Zone: Reset -->
            <div class="p-8 bg-red-50 dark:bg-red-900/10 rounded-[2.5rem] border border-red-100 dark:border-red-900/30 space-y-6">
              <div class="flex items-center gap-4">
                <div class="w-12 h-12 bg-red-600 text-white rounded-2xl flex items-center justify-center shadow-lg">
                  <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
                  </svg>
                </div>
                <div>
                  <h4 class="text-sm font-black uppercase tracking-widest text-red-600">Danger Zone: Reset Sistem</h4>
                  <p class="text-[10px] font-bold text-red-400 leading-relaxed uppercase tracking-tight">MENGHAPUS SELURUH ASSET, SISWA, RIWAYAT, DAN UPLOADS. TIDAK DAPAT DIBATALKAN!</p>
                </div>
              </div>
              <div class="flex flex-col md:flex-row gap-4 items-end">
                <div class="flex-1 space-y-2">
                  <label class="text-[10px] font-black text-red-400 uppercase tracking-widest">Ketik "HAPUS" untuk konfirmasi</label>
                  <input v-model="resetConfirmText" class="input-field !border-red-200 !text-red-600 rounded-2xl h-14" placeholder="HAPUS" />
                </div>
                <button @click="handleReset" :disabled="isProcessing || resetConfirmText !== 'HAPUS'" 
                        class="btn-premium-action !bg-red-600 !px-10 h-14 disabled:opacity-30 disabled:cursor-not-allowed">
                  {{ isProcessing ? 'MENGHAPUS...' : 'RESET TOTAL SISTEM' }}
                </button>
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import api from '../utils/api'

const authStore = useAuthStore()
const settingsStore = useSettingsStore()
const activeTab = ref('profile')
const submitting = ref(false)
const submittingBranding = ref(false)
const isProcessing = ref(false)
const resetConfirmText = ref('')
const restoreFile = ref(null)
const appLogoInput = ref(null)
const loginLogoInput = ref(null)
const faviconInput = ref(null)
const pwa192Input = ref(null)
const pwa512Input = ref(null)
const isDark = ref(document.documentElement.classList.contains('dark'))

const passForm = ref({ old_password: '', new_password: '', confirm_password: '' })
const brandingForm = ref({ ...settingsStore.settings })

// Sync form with store when settings are loaded
watch(() => settingsStore.settings, (newVal) => {
  brandingForm.value = { ...newVal }
}, { deep: true })

async function handleLogoUpload(event, key) {
  const file = event.target.files[0]
  if (!file) return

  // Basic validation
  if (file.size > 2 * 1024 * 1024) {
    alert('File terlalu besar! Maksimal 2MB.')
    return
  }

  const formData = new FormData()
  formData.append('file', file)

  try {
    submittingBranding.value = true
    const { data } = await api.post('/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    if (data.success) {
      brandingForm.value[key] = data.data.url
    }
  } catch (e) {
    alert('Gagal mengunggah logo: ' + (e.response?.data?.error || e.message))
  } finally {
    submittingBranding.value = false
    event.target.value = ''
  }
}

async function handleChangePassword() {
  if (passForm.value.new_password !== passForm.value.confirm_password) {
    alert('Konfirmasi password tidak cocok!')
    return
  }
  submitting.value = true
  try {
    await api.put('/users/me/password', passForm.value)
    alert('Password berhasil diganti!')
    passForm.value = { old_password: '', new_password: '', confirm_password: '' }
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal mengganti password')
  } finally {
    submitting.value = false
  }
}

const userThemeKey = computed(() => `theme_${authStore.user?.username || 'default'}`)
const toggleMode = () => {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark')
  localStorage.setItem(userThemeKey.value, isDark.value ? 'dark' : 'light')
}

async function handleSaveBranding() {
  submittingBranding.value = true
  try {
    const success = await settingsStore.updateSettings(brandingForm.value)
    if (success) {
      alert('Konfigurasi branding berhasil diperbarui!')
    }
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal memperbarui branding')
  } finally {
    submittingBranding.value = false
  }
}

async function handleResetBranding() {
  if (!confirm('APAKAH ANDA YAKIN? Semua kustomisasi branding (nama, teks, and ikon) akan dikembalikan ke setelan awal.')) return
  
  submittingBranding.value = true
  try {
    const success = await settingsStore.resetBranding()
    if (success) {
      brandingForm.value = { ...settingsStore.settings }
      alert('Branding telah dikembalikan ke setelan awal!')
    }
  } catch (e) {
    alert('Gagal mereset branding: ' + (e.response?.data?.error || e.message))
  } finally {
    submittingBranding.value = false
  }
}

async function handleBackup() {
  isProcessing.value = true
  try {
    const response = await api.get('/admin/system/backup', { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `sis_inv_backup_${new Date().toISOString().split('T')[0]}.zip`)
    document.body.appendChild(link)
    link.click()
    link.remove()
  } catch (e) {
    alert('Gagal mengunduh backup')
  } finally {
    isProcessing.value = false
  }
}

async function onFileSelected(event) {
  const file = event.target.files[0]
  if (!file) return

  if (!confirm('PERINGATAN: Seluruh data saat ini akan dihapus dan diganti dengan data dari backup. Lanjutkan?')) {
    event.target.value = ''
    return
  }

  isProcessing.value = true
  const formData = new FormData()
  formData.append('file', file)

  try {
    await api.post('/admin/system/restore', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    alert('Sistem berhasil dipulihkan! Halaman akan dimuat ulang.')
    window.location.reload()
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal melakukan restore')
  } finally {
    isProcessing.value = false
    event.target.value = ''
  }
}

async function handleReset() {
  if (resetConfirmText.value !== 'HAPUS') return
  
  if (!confirm('APAKAH ANDA YAKIN? Seluruh data aset, siswa, dan transaksi akan DIHAPUS PERMANEN.')) return

  isProcessing.value = true
  try {
    await api.post('/admin/system/reset')
    alert('Sistem telah direset ke kondisi awal. Halaman akan dimuat ulang.')
    window.location.reload()
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal melakukan reset')
  } finally {
    isProcessing.value = false
    resetConfirmText.value = ''
  }
}
</script>
