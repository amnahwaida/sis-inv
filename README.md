# SIS-INV (Sistem Inventaris Sekolah)

SIS-INV adalah aplikasi manajemen inventaris sekolah berbasis web yang dirancang untuk melacak, memantau, dan mengelola aset fisik sekolah secara efisien. Proyek ini menggunakan arsitektur **Go (Gin)** di backend, **Vue 3 (Vite)** di frontend, dan **PostgreSQL 15** sebagai database.

## 🚀 Fitur Utama (Phase 1)
- **Autentikasi JWT**: Login, Logout, dan Manajemen Sesi.
- **RBAC (Role-Based Access Control)**: Izin akses berdasarkan role (ADMIN, TEACHER, HEAD).
- **Manajemen User**: Admin dapat mengelola akun staff dan guru.
- **Modern UI**: Antarmuka responsif dengan TailwindCSS dan desain glassmorphism.

## 🛠️ Prasyarat
Pastikan Anda telah menginstal komponen berikut di sistem Anda:
- **Docker & Docker Compose**
- **Go 1.25+**
- **Node.js 22.12+**
- **npm** (biasanya dibundel dengan Node.js)

## 🏗️ Cara Menjalankan

### 1. Persiapan Database
Jalankan container PostgreSQL menggunakan Docker Compose:
```bash
# Di root direktori proyek
sudo docker compose up -d
```
Database akan berjalan di port `5432` dengan nama `inventaris_db`.

### 2. Menjalankan Backend (API)
Gunakan Go untuk menjalankan server API:
```bash
cd backend
# Install dependensi
go mod tidy
# Jalankan server
go run cmd/server/main.go
```
API akan berjalan di `http://localhost:8080`. Migration dan Seeding admin default akan dijalankan secara otomatis saat startup pertama kali.

### 3. Menjalankan Frontend
Gunakan npm untuk menjalankan server development Vite:
```bash
cd frontend
# Install dependensi
npm install
# Jalankan development server
npm run dev
```
Frontend akan berjalan di `http://localhost:5173` (atau port berikutnya jika 5173 sedang digunakan).

## 🔑 Kredensial Default
Gunakan akun administrator berikut untuk login pertama kali:
- **Username:** `admin`
- **Password:** `admin123`

## 📂 Struktur Proyek
- `/backend`: Source code API (Golang Gin, PostgreSQL, JWT).
- `/frontend`: Source code antarmuka (Vue 3, Vite, TailwindCSS).
- `docker-compose.yml`: Defini infrastructure (PostgreSQL).
- `prd.md`: Dokumentasi Product Requirements Document.

## 📝 Catatan Pengembangan
- **CORS:** Jika port frontend berubah, pastikan port tersebut terdaftar di `backend/internal/middleware/cors.go`.
- **Environment:** Variabel environment dikonfigurasi di file `.env` di root direktori.

---
🚀 *SIS-INV - Membangun Transparansi Aset Sekolah.*
