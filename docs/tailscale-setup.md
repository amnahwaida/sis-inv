# Panduan Keamanan & Akses Jarak Jauh (Tailscale & HTTPS)

Dokumen ini berisi panduan untuk Fase 4 (Security) SIS-INV untuk mengamankan koneksi server lokal ke internet melalui VPN Zero-Trust dan HTTPS tanpa perlu repot konfigurasi Port Forwarding pada router sekolah.

## 1. Menggunakan Tailscale (Zero-Trust VPN)
Tailscale adalah Mesh VPN berbasis WireGuard yang paling mudah digunakan. Semua perangkat yang terhubung ke jaringan Tailscale yang sama bisa saling berkomunikasi lokal (seperti LAN), meski berbeda kota atau provider.

### Cara Setup di Server Lokal SIS-INV:
1. Daftar atau Login di [Tailscale](https://tailscale.com/)
2. Di server (Linux Ubuntu/Debian tempat SIS-INV berjalan), jalankan:
   ```bash
   curl -fsSL https://tailscale.com/install.sh | sh
   ```
3. Aktifkan Tailscale:
   ```bash
   sudo tailscale up
   ```
4. Copy link otentikasi yang muncul di terminal ke browser, login menggunakan akun Tailscale Anda.
5. Server akan mendapatkan IP khusus (biasanya berawalan `100.x.y.z`). Cek dengan perintah: `tailscale ip -4`

### Mengakses dari Handphone/Komputer Guru
1. Download aplikasi Tailscale di Google Play Store (Android), App Store (iOS), atau Windows/Mac.
2. Login menggunakan akun/tim yang sama.
3. Buka browser di HP/Laptop tersebut, dan ketikkan alamat IP Tailscale dari Server beserta port SIS-INV (contoh: `http://100.x.y.z:5173`).

---

## 2. HTTPS Mudah dengan Local SSL / Tailscale Funnel (Opsi Lanjutan)

Karena aplikasi SIS-INV kini merupakan **PWA (Progressive Web App)**, kadang browser modern mengharuskan akses berjalan di HTTPS (terutama fitur Kamera & Scanner QR). Jika diakses lewat IP Tailscale biasa, SSL Certificate mungkin tidak valid.

Berikut cara menggunakan fitur **Tailscale Funnel** untuk mem-publish server lokal menggunakan HTTPS dengan domain `.ts.net`:

1. Pastikan server sudah join tailscale (`tailscale up`).
2. Jalankan perintah `tailscale funnel` untuk melayani traffic ke port Frontend (contoh di atas menggunakan Vite di `5173`):
   ```bash
   tailscale serve https / http://127.0.0.1:5173
   tailscale funnel 443 on
   ```
3. Setelah berhasil di-funnel, Tailscale akan memberikan public URL berformat `https://[nama-server].[nama-tailnet].ts.net`.
4. Buka URL ini di browser, maka aplikasi akan menggunakan **HTTPS valid**. PWA dapat diinstall dan kamera barcode scanner akan berfungsi tanpa error dari browser.

---

## 3. Log Audit Sistem
Aksi krusial di sistem seperti:
1. Pembuatan Akun Baru (User)
2. Pembuatan Data Barang
3. Penghapusan Data Barang (Hard & Soft delete)

Sudah mulai dicatat secara asinkron di dalam tabel `audit_logs` (Fitur Log Minimalis). Log ini dapat digunakan oleh tim IT jika terjadi perselisihan data krusial di sekolah.
