# Product Requirements Document (PRD)
# Sistem Inventaris Sekolah (SIS-INV) - MVP

| Dokumen | Versi | Tanggal | Status |
|---------|-------|---------|--------|
| PRD-MVP | 4.0 | 2024 | **Final - Web Based, No Student Account** |

---

## 1. Ringkasan Eksekutif

**SIS-INV** adalah aplikasi manajemen inventaris sekolah berbasis **web** yang dirancang untuk melacak, memantau, dan mengelola aset fisik sekolah (elektronik, furniture, alat lab, olahraga). Sistem ini menggunakan arsitektur **Local Server + Zero Trust Network** dengan **client web-based**.

### 🔥 Perubahan Strategis Utama (v4.0)
**Siswa TIDAK perlu akun/login.** Semua transaksi peminjaman oleh siswa dilakukan oleh Guru/Staff yang bertanggung jawab.

| Aspek | Versi 3.0 (Dengan Akun Siswa) | Versi 4.0 (Tanpa Akun Siswa) |
|-------|-------------------------------|------------------------------|
| **Student Access** | Login sendiri via browser | Tidak login, diwakili guru |
| **Alur Pinjam Siswa** | Siswa request → Guru approve → Ambil | Guru input NIS → Scan barang → Selesai |
| **Kompleksitas Auth** | 5 roles (termasuk STUDENT) | 4 roles (tanpa STUDENT) |
| **Database** | Tabel users + guardian mapping | Lebih sederhana, student info di transaksi |
| **Dev Time** | ~10 minggu | **~7 minggu** |
| **User Management** | Kelola ratusan akun siswa | Hanya kelola akun guru/staff |
| **Akuntabilitas** | Via sistem approval | Via tanggung jawab guru |

### 1.1 Visi Produk
> "Membuat pengelolaan aset sekolah menjadi transparan, akuntabel, dan efisien dengan teknologi yang mudah diakses tanpa instalasi dan tanpa beban manajemen akun siswa."

### 1.2 Masalah yang Diselesaikan
| Masalah | Solusi SIS-INV |
|---------|----------------|
| Barang sekolah sering hilang tanpa jejak | Tracking real-time dengan QR Code |
| Tidak ada riwayat peminjaman yang jelas | Log transaksi digital dengan timestamp |
| Sulit melakukan audit/stock opname | Mode audit khusus dengan scan massal |
| Data tersimpan di Excel yang rentan | Database terstruktur dengan backup otomatis |
| Instalasi app ribet di banyak device | Akses via browser, tanpa instalasi |
| **Kelola ratusan akun siswa ribet** | **Siswa tidak perlu akun, guru yang handle** |
| Siswa lupa password/username | Tidak ada login untuk siswa |

---

## 2. Tujuan & Sasaran

### 2.1 Business Goals
1. **Mengurangi kehilangan aset** sekolah minimal 30% dalam 1 tahun
2. **Mempercepat proses audit** inventaris dari 2 minggu menjadi 2 hari
3. **Memberikan visibilitas** real-time tentang kondisi dan lokasi aset
4. **Menghemat biaya** dengan tidak perlu subscription cloud (self-hosted)
5. **Memudahkan akses** tanpa perlu instalasi aplikasi di setiap device
6. **Menghilangkan beban admin** untuk mengelola akun siswa
7. **Mengurangi waktu development** dari 10 minggu menjadi 7 minggu

### 2.2 Success Metrics (KPI)
| Metric | Target MVP | Target Tahun 1 |
|--------|------------|----------------|
| Akurasi data inventaris | 95% | 99% |
| Waktu respons scan QR | < 3 detik | < 2 detik |
| Uptime server | 95% | 99% |
| Adopsi user (staff/guru) | 90% | 100% |
| Data backup success rate | 100% | 100% |
| Return rate on-time (via guru) | 95% | 98% |
| Page load time | < 3 detik | < 2 detik |
| Waktu input peminjaman siswa | < 30 detik | < 20 detik |

---

## 3. Target Pengguna

### 3.1 User Personas

| Persona | Deskripsi | Kebutuhan Utama | Device | Akses Level |
|---------|-----------|-----------------|--------|-------------|
| **Admin TU** | Staff Tata Usaha yang mengelola inventaris | Input barang, laporan, manajemen user | Desktop/Laptop | **FULL** |
| **Guru/Staff** | Pengguna yang meminjam barang & menangani peminjaman siswa | Pinjam cepat, input pinjaman siswa, lihat riwayat | Mobile/Desktop | **STANDARD** |
| **Kepala Sekolah** | Pengawas yang butuh laporan | Dashboard, laporan ringkas | Desktop/Tablet | **READ-ONLY** |
| **Teknisi IT** | Maintenance server | Backup, monitoring, troubleshooting | Desktop | **ADMIN** |

### 3.2 User Stories

```
SEBAGAI Admin TU
SAYA INGIN bisa input barang dari komputer TU
SEHINGGA tidak perlu install aplikasi khusus

SEBAGAI Guru
SAYA INGIN bisa pinjam proyektor lewat browser HP
SEHINGGA tidak perlu download app di App Store

SEBAGAI Guru (untuk siswa)
SAYA INGIN bisa input peminjaman atas nama siswa dengan cepat
SEHINGGA siswa tidak perlu ribet login, tapi data tetap tercatat

SEBAGAI Kepala Sekolah
SAYA INGIN lihat laporan dari browser laptop
SEHINGGA bisa akses kapan saja tanpa install

SEBAGAI Teknisi IT
SAYA INGIN backup otomatis setiap malam
SEHINGGA data aman jika server rusak
```

---

## 4. Fitur MVP (Scope)

### 4.1 Fitur Wajib (Must Have - P0)

| ID | Fitur | Prioritas | Estimasi | User Role |
|----|-------|-----------|----------|-----------|
| **F01** | **Autentikasi & Autorisasi** | P0 | 2 hari | All |
| F01.1 | Login dengan username/password | P0 | | All |
| F01.2 | Role-based access (ADMIN, TEACHER, HEAD) | P0 | | All |
| F01.3 | Session management dengan JWT | P0 | | All |
| **F02** | **Manajemen Barang (CRUD)** | P0 | 4 hari | Admin |
| F02.1 | Tambah barang dengan form | P0 | | Admin |
| F02.2 | Generate QR Code otomatis | P0 | | Admin |
| F02.3 | Edit informasi barang | P0 | | Admin |
| F02.4 | Hapus/Non-aktifkan barang | P0 | | Admin |
| F02.5 | Set allowed borrower types per item | P0 | | Admin |
| F02.6 | Print QR Code label (PDF) | P0 | | Admin |
| **F03** | **Transaksi Pinjam-Kembali (Guru untuk Diri Sendiri)** | P0 | 3 hari | Teacher |
| F03.1 | Scan QR via web camera | P0 | | Teacher |
| F03.2 | Manual input kode barang | P0 | | Teacher |
| F03.3 | Cek kondisi saat pengembalian | P0 | | Teacher |
| F03.4 | Upload foto kondisi barang | P1 | | Teacher |
| F03.5 | Riwayat transaksi pribadi | P1 | | Teacher |
| **F04** | **Transaksi Pinjam-Kembali (Guru untuk Siswa)** | P0 | 3 hari | Teacher |
| F04.1 | Input NIS/Nama siswa (dengan auto-suggest) | P0 | | Teacher |
| F04.2 | Scan QR barang untuk siswa | P0 | | Teacher |
| F04.3 | Input tujuan & durasi pinjaman siswa | P0 | | Teacher |
| F04.4 | Time-limited borrowing (max 6 hours untuk siswa) | P0 | | Teacher |
| F04.5 | Quick return scan untuk siswa | P0 | | Teacher |
| F04.6 | Riwayat peminjaman per siswa (view only) | P1 | | Teacher |
| **F05** | **Dashboard & Laporan** | P1 | 3 hari | Admin/Head |
| F05.1 | Dashboard ringkas (total, dipinjam, rusak) | P1 | | Admin/Head |
| F05.2 | Laporan barang dipinjam saat ini | P1 | | Admin |
| F05.3 | Laporan peminjaman per kelas/siswa | P1 | | Admin |
| F05.4 | Export laporan ke Excel/PDF | P1 | | Admin |
| **F06** | **Infrastruktur & Keamanan** | P0 | 4 hari | IT |
| F06.1 | Docker Compose deployment | P0 | | IT |
| F06.2 | Cloudflare Zero Trust integration | P0 | | IT |
| F06.3 | Automated backup script | P0 | | IT |
| F06.4 | Firewall configuration | P0 | | IT |
| F06.5 | HTTPS with self-signed cert | P0 | | IT |
| **F07** | **Progressive Web App (PWA)** | P1 | 2 hari | All |
| F07.1 | Install to home screen | P1 | | All |
| F07.2 | Offline cache for basic pages | P1 | | All |

### 4.2 Fitur Tambahan (Nice to Have - Post MVP)

| ID | Fitur | Prioritas | Fase |
|----|-------|-----------|------|
| F08 | Mode Stock Opname / Audit | P1 | Phase 2 |
| F09 | Notifikasi overdue (email) | P2 | Phase 2 |
| F10 | Maintenance log & history | P2 | Phase 2 |
| F11 | Barang habis pakai (consumables) | P2 | Phase 3 |
| F12 | Multi-location / Multi-gedung | P2 | Phase 3 |
| F13 | Auto-blacklist siswa bermasalah | P3 | Phase 3 |
| F14 | Depreciation calculation | P3 | Phase 4 |
| F15 | Integrasi dengan data siswa (Dapodik) | P3 | Phase 4 |

### 4.3 Fitur yang TIDAK Termasuk MVP

- ❌ Akun/login untuk siswa
- ❌ Sistem perpustakaan/buku
- ❌ Sistem denda otomatis (hanya tracking overdue)
- ❌ Reservasi/pengantrian barang
- ❌ Payment gateway integration
- ❌ AI/ML prediction
- ❌ Mobile app native (Android/iOS)
- ❌ Approval workflow (karena guru langsung input)

---

## 5. Spesifikasi Teknis

### 5.1 Arsitektur Sistem

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT DEVICES                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │  Mobile  │  │ Desktop  │  │  Tablet  │  │  Laptop  │        │
│  │ Browser  │  │ Browser  │  │ Browser  │  │ Browser  │        │
│  │(Chrome)  │  │(Chrome)  │  │ (Safari) │  │(Firefox) │        │
│  │ (Guru)   │  │ (Admin)  │  │ (Guru)   │  │ (Head)   │        │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘        │
│       │             │             │             │               │
│       └─────────────┴──────┬──────┴─────────────┘               │
│                            │                                    │
│                    Cloudflare Zero Trust                         │
│                    (Encrypted Tunnel)                           │
│                    HTTPS (Self-Signed)                          │
└────────────────────────────┼────────────────────────────────────┘
                             │
                    ┌────────▼────────┐
                    │  Router Sekolah │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │   Host Server   │
                    │  (Mini PC/NUC)  │
                    │  + Zero Trust   │
                    │  + UPS          │
                    └────────┬────────┘
                             │
              ┌──────────────┴──────────────┐
              │     Docker Compose          │
              │  ┌─────────────────────┐    │
              │  │   Web Server        │    │
              │  │   (Caddy/Nginx)     │    │
              │  │   Port: 443 (HTTPS) │    │
              │  └──────────┬──────────┘    │
              │             │               │
              │  ┌──────────▼──────────┐    │
              │  │   API Container     │    │
              │  │   (Go + Gin)        │    │
              │  │   Port: 8080        │    │
              │  └──────────┬──────────┘    │
              │             │               │
              │  ┌──────────▼──────────┐    │
              │  │   DB Container      │    │
              │  │   (PostgreSQL 15)   │    │
              │  │   Volume: pgdata    │    │
              │  └─────────────────────┘    │
              └─────────────────────────────┘
                              │
                    ┌────────▼────────┐
                    │  Backup Storage │
                    │  (GDrive/Ext.)  │
                    └─────────────────┘
```

### 5.2 Tech Stack

| Layer | Technology | Version | Alasan |
|-------|------------|---------|--------|
| **Frontend** | HTML/CSS/JavaScript | ES6+ | Universal browser support |
| **Frontend Framework** | Vue.js 3 atau React | Latest | Reactive UI, component-based |
| **UI Library** | TailwindCSS + HeadlessUI | Latest | Fast styling, responsive |
| **Backend** | Go (Golang) | 1.21+ | Performa tinggi, compiled |
| **Framework** | Gin | Latest | Lightweight HTTP framework |
| **Database** | PostgreSQL | 15 | Robust, relasional, support JSONB |
| **Driver** | pgx | v5 | Native PostgreSQL driver untuk Go |
| **Web Server** | Caddy | Latest | Auto HTTPS, simple config |
| **Container** | Docker | Latest | Konsistensi environment |
| **Orchestration** | Docker Compose | v3.8 | Multi-service management |
| **Network** | Cloudflare Tunnel | Latest | Zero Trust, SSL, Domain |
| **Auth** | JWT | - | Stateless authentication |
| **PWA** | Service Worker | - | Offline caching, installable |
| **QR Library** | html5-qrcode | Latest | Web-based QR scanning |
| **QR Generate** | qrcode.js | Latest | Client-side QR generation |

### 5.3 Database Schema (MVP)

```sql
-- ============================================================================
-- TABEL USERS (Hanya Admin, Guru, Kepala Sekolah - TANPA SISWA)
-- ============================================================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('ADMIN', 'TEACHER', 'HEAD')),
    nip VARCHAR(50), -- NIP untuk guru/staff
    email VARCHAR(100),
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk performa
CREATE INDEX idx_users_role ON users(role);

-- ============================================================================
-- TABEL CATEGORIES (Kategori barang)
-- ============================================================================
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    color_code VARCHAR(7), -- Untuk UI color coding
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- TABEL ITEMS (Barang inventaris)
-- ============================================================================
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL, -- Kode aset unik
    qr_code_data VARCHAR(255), -- Data QR Code
    name VARCHAR(255) NOT NULL,
    category_id INT REFERENCES categories(id),
    location VARCHAR(100), -- Lokasi penyimpanan
    condition VARCHAR(20) DEFAULT 'GOOD' CHECK (condition IN ('GOOD', 'DAMAGED', 'LOST', 'IN_REPAIR')),
    status VARCHAR(20) DEFAULT 'AVAILABLE' CHECK (status IN ('AVAILABLE', 'BORROWED', 'MAINTENANCE', 'LOST')),
    borrower_type VARCHAR(20) DEFAULT 'STAFF_ONLY' CHECK (borrower_type IN ('STAFF_ONLY', 'STUDENT_ALLOWED')),
    purchase_date DATE,
    purchase_price DECIMAL(15,2),
    warranty_end_date DATE,
    notes TEXT,
    photo_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk performa
CREATE INDEX idx_items_status ON items(status);
CREATE INDEX idx_items_category ON items(category_id);
CREATE INDEX idx_items_location ON items(location);
CREATE INDEX idx_items_borrower_type ON items(borrower_type);

-- ============================================================================
-- TABEL STUDENTS (Referensi data siswa - HANYA UNTUK LOOKUP, BUKAN UNTUK LOGIN)
-- ============================================================================
-- Tabel ini opsional. Bisa juga hanya input manual NIS+Nama saat transaksi.
-- Jika ada, bisa di-import dari Dapodik/Excel sekolah.
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    nis VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    class VARCHAR(20), -- "12 IPA 1"
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk performa lookup
CREATE INDEX idx_students_nis ON students(nis);
CREATE INDEX idx_students_class ON students(class);
CREATE INDEX idx_students_name ON students(full_name);

-- ============================================================================
-- TABEL TRANSACTIONS (Transaksi pinjam-kembali)
-- ============================================================================
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID REFERENCES items(id),
    
    -- Borrower bisa guru (user_id) ATAU siswa (student_nis)
    borrower_type VARCHAR(10) NOT NULL CHECK (borrower_type IN ('STAFF', 'STUDENT')),
    user_id UUID REFERENCES users(id), -- NULL jika borrower_type = 'STUDENT'
    student_nis VARCHAR(50), -- NULL jika borrower_type = 'STAFF'
    student_name VARCHAR(100), -- Snapshot nama siswa saat transaksi
    student_class VARCHAR(20), -- Snapshot kelas siswa saat transaksi
    
    borrowed_by UUID REFERENCES users(id) NOT NULL, -- Guru yang memproses transaksi
    borrowed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMP NOT NULL,
    returned_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'RETURNED', 'OVERDUE')),
    purpose TEXT, -- Tujuan peminjaman
    return_condition VARCHAR(20), -- Kondisi saat return
    return_notes TEXT,
    return_photo_url VARCHAR(500), -- Foto kondisi return
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk performa
CREATE INDEX idx_transactions_item ON transactions(item_id);
CREATE INDEX idx_transactions_user ON transactions(user_id);
CREATE INDEX idx_transactions_student ON transactions(student_nis);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_borrowed_by ON transactions(borrowed_by);
CREATE INDEX idx_transactions_due_date ON transactions(due_date);

-- ============================================================================
-- TABEL MAINTENANCE_LOGS (Riwayat perbaikan)
-- ============================================================================
CREATE TABLE maintenance_logs (
    id SERIAL PRIMARY KEY,
    item_id UUID REFERENCES items(id),
    reported_by UUID REFERENCES users(id),
    reported_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    issue_description TEXT NOT NULL,
    cost DECIMAL(15,2),
    vendor VARCHAR(100),
    status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'IN_PROGRESS', 'DONE', 'CANCELLED')),
    completed_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk performa
CREATE INDEX idx_maintenance_item ON maintenance_logs(item_id);
CREATE INDEX idx_maintenance_status ON maintenance_logs(status);

-- ============================================================================
-- TABEL AUDIT_LOGS (Log aktivitas sistem)
-- ============================================================================
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50),
    entity_id UUID,
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(45),
    user_agent VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk performa
CREATE INDEX idx_audit_user ON audit_logs(user_id);
CREATE INDEX idx_audit_action ON audit_logs(action);
CREATE INDEX idx_audit_created ON audit_logs(created_at);
```

### 5.4 API Endpoints (MVP)

#### Authentication
| Method | Endpoint | Auth | Role | Deskripsi |
|--------|----------|------|------|-----------|
| POST | `/api/v1/auth/login` | ❌ | All | Login user |
| POST | `/api/v1/auth/logout` | ✅ | All | Logout user |
| GET | `/api/v1/auth/me` | ✅ | All | Get current user info |
| PUT | `/api/v1/auth/password` | ✅ | All | Change password |

#### Items Management (Admin Only)
| Method | Endpoint | Auth | Role | Deskripsi |
|--------|----------|------|------|-----------|
| GET | `/api/v1/items` | ✅ | All | List all items (with filter) |
| POST | `/api/v1/items` | ✅ | ADMIN | Create new item |
| GET | `/api/v1/items/:id` | ✅ | All | Get item detail |
| PUT | `/api/v1/items/:id` | ✅ | ADMIN | Update item |
| DELETE | `/api/v1/items/:id` | ✅ | ADMIN | Delete item |
| POST | `/api/v1/items/:id/qr` | ✅ | ADMIN | Generate QR code |
| GET | `/api/v1/items/print-labels` | ✅ | ADMIN | Print QR labels (PDF) |

#### Transactions - Staff Borrowing (Untuk Diri Sendiri)
| Method | Endpoint | Auth | Role | Deskripsi |
|--------|----------|------|------|-----------|
| POST | `/api/v1/transactions/borrow/staff` | ✅ | TEACHER, ADMIN | Borrow for self |
| POST | `/api/v1/transactions/return` | ✅ | TEACHER, ADMIN | Return item (self or student) |
| GET | `/api/v1/transactions/my-borrowings` | ✅ | TEACHER, ADMIN | View own active borrowings |

#### Transactions - Student Borrowing (Diproses Guru)
| Method | Endpoint | Auth | Role | Deskripsi |
|--------|----------|------|------|-----------|
| GET | `/api/v1/students/search?q=` | ✅ | TEACHER, ADMIN | Search student by NIS/name |
| POST | `/api/v1/transactions/borrow/student` | ✅ | TEACHER, ADMIN | Borrow item for student |
| GET | `/api/v1/transactions/student/:nis` | ✅ | TEACHER, ADMIN | View borrowing history for student |
| GET | `/api/v1/transactions/class/:class` | ✅ | TEACHER, ADMIN | View active borrowings by class |

#### Dashboard & Reports
| Method | Endpoint | Auth | Role | Deskripsi |
|--------|----------|------|------|-----------|
| GET | `/api/v1/dashboard/summary` | ✅ | ADMIN, HEAD | Dashboard statistics |
| GET | `/api/v1/reports/overdue` | ✅ | ADMIN | Report overdue items |
| GET | `/api/v1/reports/active-borrowings` | ✅ | ADMIN | Currently borrowed items |
| GET | `/api/v1/reports/by-student` | ✅ | ADMIN | Borrowing report per student |
| POST | `/api/v1/reports/export` | ✅ | ADMIN | Export to Excel/PDF |

#### Maintenance
| Method | Endpoint | Auth | Role | Deskripsi |
|--------|----------|------|------|-----------|
| POST | `/api/v1/maintenance` | ✅ | ADMIN, TEACHER | Create maintenance log |
| GET | `/api/v1/maintenance/:item_id` | ✅ | ADMIN | Get maintenance history |
| PUT | `/api/v1/maintenance/:id` | ✅ | ADMIN | Update maintenance status |

#### System (IT Admin)
| Method | Endpoint | Auth | Role | Deskripsi |
|--------|----------|------|------|-----------|
| GET | `/api/v1/system/health` | ❌ | Public | Health check endpoint |
| GET | `/api/v1/system/backup-status` | ✅ | ADMIN | Check last backup status |
| POST | `/api/v1/system/trigger-backup` | ✅ | ADMIN | Manual backup trigger |

---

## 6. Keamanan & Compliance

### 6.1 Security Requirements

| Requirement | Implementation | Priority |
|-------------|----------------|----------|
| **Network Security** | Cloudflare Zero Trust, Tunnel, no port forwarding | P0 |
| **HTTPS** | Self-signed certificate via Caddy | P0 |
| **Authentication** | JWT with expiry (24 hours access, 7 days refresh) | P0 |
| **Password** | Bcrypt hashing (cost 12) | P0 |
| **Authorization** | Role-based access control (RBAC) | P0 |
| **Data Encryption** | TLS via Cloudflare Tunnel | P0 |
| **Firewall** | Cloudflare Tunnel only, local loopback | P0 |
| **Session** | Refresh token mechanism | P1 |
| **Audit Log** | Log semua aktivitas penting + browser info | P1 |
| **Rate Limiting** | Max 100 requests/minute per IP | P1 |
| **CORS** | Restrict to school domain | P1 |
| **XSS Protection** | Content Security Policy headers | P1 |
| **CSRF Protection** | CSRF tokens for state-changing operations | P1 |

### 6.2 Student Borrowing Policy (Via Guru)

```yaml
Student Borrowing Rules (Handled by Teacher):
  - Teachers can only borrow student-eligible items for students
  - Items must have borrower_type = 'STUDENT_ALLOWED'
  - Maximum borrowing duration for students: 6 hours (same day return)
  - Maximum concurrent borrowings per student: 2 items
  - Teacher is responsible for verifying student identity (NIS)
  - Teacher is responsible for item condition check on return
  - Overdue/damaged items: Teacher receives notification
  - Repeat offenders: Teacher can flag student for temporary suspension
```

### 6.3 Role Permission Matrix

| Feature | ADMIN | TEACHER | HEAD |
|---------|-------|---------|------|
| Create/Edit Items | ✅ | ❌ | ❌ |
| Delete Items | ✅ | ❌ | ❌ |
| Borrow (Self) | ✅ | ✅ | ✅ |
| Borrow (For Student) | ✅ | ✅ | ❌ |
| View All Transactions | ✅ | ❌ | ✅ |
| View Own Transactions | ✅ | ✅ | ✅ |
| View Student Borrowings | ✅ | ✅ (by class) | ✅ |
| Generate Reports | ✅ | ❌ | ✅ |
| System Settings | ✅ | ❌ | ❌ |
| Maintenance Log | ✅ | ✅ | ✅ |
| Print QR Labels | ✅ | ❌ | ❌ |
| Manage Users | ✅ | ❌ | ❌ |

### 6.4 Backup Strategy

```bash
#!/bin/bash
# /opt/inventaris/backup.sh

# Configuration
BACKUP_DIR="/backup/inventaris"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30
DB_CONTAINER="inventaris_db"
DB_USER="sekolah_admin"
DB_NAME="inventaris_db"

# Create backup directory
mkdir -p $BACKUP_DIR

# Database backup
echo "[$(date)] Starting database backup..."
docker exec $DB_CONTAINER pg_dump -U $DB_USER $DB_NAME > $BACKUP_DIR/db_$DATE.sql

# Compress backup
gzip $BACKUP_DIR/db_$DATE.sql
echo "[$(date)] Database backup completed: db_$DATE.sql.gz"

# Upload to cloud storage (optional - using rclone)
if [ -f "/root/.config/rclone/rclone.conf" ]; then
    rclone copy $BACKUP_DIR/db_$DATE.sql.gz remote:inventaris-backup/
    echo "[$(date)] Backup uploaded to cloud storage"
fi

# Delete old backups
echo "[$(date)] Cleaning up old backups..."
find $BACKUP_DIR -name "*.sql.gz" -mtime +$RETENTION_DAYS -delete

# Log completion
echo "[$(date)] Backup completed successfully" >> $BACKUP_DIR/backup.log
```

**Cron Job Setup:**
```bash
# Edit crontab
crontab -e
# Add this line (backup every day at 02:00 AM)
0 2 * * * /opt/inventaris/backup.sh >> /opt/inventaris/backup_cron.log 2>&1
```

### 6.5 Disaster Recovery

| Scenario | Recovery Time Objective (RTO) | Recovery Point Objective (RPO) | Procedure |
|----------|-------------------------------|--------------------------------|-----------|
| Server crash | 4 jam | 24 jam (backup terakhir) | Restore dari backup ke hardware baru |
| Database corrupt | 2 jam | 24 jam | pg_restore dari file backup |
| Data accidental delete | 1 jam | 24 jam | Restore specific table from backup |
| Full system failure | 8 jam | 24 jam | Rebuild Docker + restore backup |
| Zero Trust tunnel lost | 30 menit | 0 | Check cloudflared service status |
| SSL cert expired | 30 menit | 0 | Caddy auto-renews, manual override if needed |

---

## 7. User Interface Guidelines

### 7.1 Design Principles
- **Simple:** Maksimal 3 klik untuk aksi utama
- **Fast:** Page load time < 3 detik
- **Responsive:** Works on mobile, tablet, desktop
- **Accessible:** WCAG 2.1 AA compliance
- **Role-Adaptive:** UI berbeda berdasarkan role user
- **PWA:** Installable to home screen

### 7.2 Key Pages (Web)

| Page | Elements | Roles | Notes |
|------|----------|-------|-------|
| **Login** | Username, Password, Login Button | All | Remember me option |
| **Dashboard** | Summary cards, Quick actions, Recent activity | All | Role-based view |
| **Scan** | Camera view, Flash toggle, Manual input | All | html5-qrcode library |
| **Item List** | Search, Filter, Table/Card view | All | Color-coded status |
| **Item Detail** | Info, QR, History, Actions | All | Edit (Admin only) |
| **Borrow (Self)** | Item scan, Purpose, Due date | Teacher | Direct borrow |
| **Borrow (Student)** | Student search (NIS), Item scan, Purpose | Teacher | Quick input flow |
| **Return** | Scan QR, Condition check, Photo upload | All | Works for self & student |
| **Student Lookup** | Search by NIS/name/class, borrowing history | Teacher | Read-only |
| **Reports** | Filters, Table, Export buttons | Admin | Excel/PDF export |
| **Profile** | User info, Settings, Logout | All | Change password |

### 7.3 Quick Student Borrow Flow (UX Focus)

```
1. Guru klik "Pinjam untuk Siswa"
2. Input NIS siswa → Auto-suggest nama & kelas (opsional: scan kartu pelajar)
3. Scan QR barang yang dipinjam
4. Input tujuan (opsional) → Due date auto-set (max 6 jam)
5. Klik "Konfirmasi" → Selesai (30 detik total)
```

### 7.4 Feedback Mechanisms
| Event | Visual | Audio | Notes |
|-------|--------|-------|-------|
| Scan Success | Green flash + Checkmark | Optional beep | Via Web Audio API |
| Scan Failed | Red flash + X mark | Optional beep | Via Web Audio API |
| Submit Success | Toast notification | - | 3 seconds auto-dismiss |
| Submit Error | Error dialog | - | With clear message |
| Loading | Skeleton screen + Spinner | - | - |
| Offline Mode | Banner indicator | - | PWA detection |
| Session Timeout | Modal redirect to login | - | 30 minutes inactivity |

### 7.5 Color Coding System

| Status | Color | Hex Code | Usage |
|--------|-------|----------|-------|
| Available | Green | `#4CAF50` | Item ready to borrow |
| Borrowed (Staff) | Blue | `#2196F3` | Item borrowed by teacher |
| Borrowed (Student) | Light Blue | `#64B5F6` | Item borrowed by student |
| Overdue | Orange | `#FF9800` | Past due date |
| Damaged | Red | `#F44336` | Needs repair |
| Lost | Dark Red | `#B71C1C` | Cannot be found |
| Maintenance | Purple | `#9C27B0` | Being repaired |

### 7.6 Responsive Breakpoints

| Breakpoint | Width | Device |
|------------|-------|--------|
| Mobile | < 640px | Phone |
| Tablet | 640px - 1024px | Tablet |
| Desktop | > 1024px | Laptop/Desktop |

---

## 8. Timeline & Roadmap

### 8.1 Development Phases

```
Phase 1: Foundation (Week 1-2)
├── Setup Docker & Database
├── Backend API basic CRUD
├── Authentication system with RBAC
├── Web project setup (Vue/React)
└── Cloudflare Tunnel integration test

Phase 2: Core Features (Week 3-4)
├── Item management (Admin)
├── QR Code generation & web scanning
├── Transaction (borrow/return) - Staff self
├── Basic dashboard
└── Backup automation

Phase 3: Student Borrowing via Teacher (Week 5-6)
├── Student lookup/search API
├── Borrow flow for student (via teacher)
├── Time-limited borrowing logic for students
├── Student borrowing reports
└── Testing with pilot teachers

Phase 4: PWA & Security (Week 7)
├── Service Worker implementation
├── Offline caching
├── HTTPS configuration
├── Security audit
└── Documentation

Phase 5: Pilot & Launch (Week 8-9)
├── Pilot deployment (1-2 departments)
├── Teacher training
├── Feedback collection & bug fixes
├── Full launch
└── Post-launch support
```

### 8.2 Milestone Checklist

| Milestone | Deliverable | Success Criteria | Week |
|-----------|-------------|------------------|------|
| M1 | Backend API Ready | All endpoints tested with Postman | 2 |
| M2 | Web App Alpha | Can login, scan, borrow/return (self) | 4 |
| M3 | Student Borrow Flow | Teacher can borrow for student via NIS | 6 |
| M4 | PWA Ready | Installable, offline cache working | 7 |
| M5 | Server Deployed | Running on school server with Zero Trust + SSL Domain | 7 |
| M6 | Backup Tested | Restore successful on test environment | 7 |
| M7 | Security Audit Passed | No critical vulnerabilities | 8 |
| M8 | Teacher Training | 90% pilot teachers can use independently | 8 |
| M9 | Pilot Complete | 2 weeks successful pilot usage | 9 |
| M10 | Go Live | System in production use | 9 |

---

## 9. Risks & Mitigation

| Risk | Impact | Probability | Mitigation | Owner |
|------|--------|-------------|------------|-------|
| **Server hardware failure** | High | Medium | UPS + Automated backup + Spare hardware | IT |
| **Internet outage** | Medium | High | PWA offline cache for basic pages | Dev |
| **Browser compatibility** | Medium | Medium | Test on Chrome, Firefox, Safari, Edge | Dev |
| **User resistance** | Medium | Medium | Training + Simple UX + Champion users | Admin |
| **QR label damaged** | Medium | High | Laminated labels + Manual input fallback | Admin |
| **Data breach** | High | Low | Zero Trust network + RBAC + Audit logs | IT |
| **Scope creep** | Medium | High | Strict MVP scope + Change request process | PM |
| **Teacher forgets to return** | Medium | Medium | Overdue notifications + Dashboard alerts | System |
| **Wrong student NIS input** | Low | Medium | Auto-suggest from student table + confirmation | Dev |
| **Database corruption** | High | Low | Daily backup + Weekly restore test | IT |
| **Zero Trust tunnel down** | Medium | Low | Local domain fallback | Dev |
| **Web camera permission denied** | Medium | Medium | Manual input fallback + Permission guide | Dev |

---

## 10. Acceptance Criteria

### 10.1 Functional Acceptance
- [ ] User dapat login dengan role yang benar
- [ ] Admin dapat tambah/edit/hapus barang
- [ ] QR Code dapat di-generate dan di-scan via web camera
- [ ] Guru dapat pinjam-kembali barang untuk diri sendiri
- [ ] Guru dapat input peminjaman untuk siswa via NIS
- [ ] Guru dapat return barang (self atau siswa) dengan cek kondisi
- [ ] Transaksi tercatat dengan benar di database
- [ ] Dashboard menampilkan data yang akurat
- [ ] Laporan dapat di-export ke Excel/PDF
- [ ] Backup berjalan otomatis setiap hari
- [ ] Restore backup berhasil dalam < 1 jam
- [ ] PWA dapat diinstall ke home screen
- [ ] Basic pages work offline (cached)
- [ ] Student lookup by NIS works with auto-suggest

### 10.2 Non-Functional Acceptance
- [ ] Page load time < 3 detik (local network)
- [ ] PWA works offline for cached pages
- [ ] Tidak ada critical security vulnerability
- [ ] Backup dapat di-restore dalam < 1 jam
- [ ] Aplikasi works on Chrome, Firefox, Safari, Edge (latest 2 versions)
- [ ] Student borrowing max 6 hours enforced
- [ ] Audit logs record all critical actions + browser info
- [ ] HTTPS enforced on all pages
- [ ] Responsive design works on mobile, tablet, desktop
- [ ] Waktu input peminjaman siswa < 30 detik

### 10.3 User Acceptance
- [ ] 90% guru/staff dapat menggunakan tanpa bantuan setelah training
- [ ] User satisfaction score > 4/5
- [ ] Tidak ada complaint critical dalam 2 minggu pertama
- [ ] Return rate on-time > 95%
- [ ] No student account management overhead

---

## 11. Documentation Requirements

| Document | Audience | Format | Owner |
|----------|----------|--------|-------|
| **User Manual** | Teachers/Admin | PDF + Video Tutorial | PM |
| **Quick Start Guide** | Teachers | 1-page PDF | PM |
| **Student Borrowing Guide** | Teachers | 1-page PDF (flow) | PM |
| **Technical Doc** | IT Staff | Markdown + Diagram | Dev |
| **API Documentation** | Developers | Swagger/OpenAPI | Dev |
| **Deployment Guide** | IT Staff | Step-by-step PDF | Dev |
| **Troubleshooting** | Support | FAQ Format | Dev |
| **Backup/Restore Guide** | IT Staff | Step-by-step PDF | IT |
| **PWA Installation Guide** | All Users | PDF + Images | Dev |

---

## 12. Budget Estimate

| Item | Cost (IDR) | Type | Notes |
|------|------------|------|-------|
| Server Hardware (Mini PC) | 3.000.000 | One-time | Intel NUC or equivalent |
| UPS (600VA) | 500.000 | One-time | For safe shutdown |
| SSD 256GB (if not included) | 400.000 | One-time | For database performance |
| QR Labels (1000 pcs, vinyl) | 300.000 | One-time | Laminated |
| QR Labels (100 pcs, aluminum) | 500.000 | One-time | For expensive items |
| Network Cable & Accessories | 200.000 | One-time | For server connection |
| Domain (optional) | 150.000 | Per year | For professional look |
| Cloudflare | 0 | Free | Free Zero Trust tier |
| SSL Certificate | 0 | Free | Self-signed via Caddy |
| Development Cost | - | - | Internal/Outsourced |
| **Total Initial** | **~5.050.000** | | Excluding development |

### Recurring Costs
| Item | Cost (IDR) | Frequency |
|------|------------|-----------|
| Electricity (server) | ~50.000 | Per month |
| Internet | Included | School existing |
| Cloud backup (optional) | 0-100.000 | Per month |
| **Total Recurring** | **~50.000-150.000** | **Per month** |

### Savings vs Previous Versions
| Item | v2.0 (Flutter+Student) | v4.0 (Web+No Student) | Savings |
|------|------------------------|----------------------|---------|
| Development Time | 12 weeks | 7 weeks | **42%** |
| Client Setup | Install app + student accounts | Browser only | **100%** |
| User Management | 100+ student accounts | Only staff accounts | **90%** |
| Maintenance | Multi-platform + student support | Single web app | **80%** |

---

## 13. Approval

| Role | Name | Signature | Date |
|------|------|-----------|------|
| Product Owner | | | |
| Technical Lead | | | |
| School Principal | | | |
| IT Manager | | | |
| Admin TU Representative | | | |

---

## 14. Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2024-01-01 | | Initial Draft (Flutter, no student) |
| 2.0 | 2024-01-15 | | Added student account + guardian approval (Flutter) |
| 3.0 | 2024-01-20 | | Changed to Web-Based PWA (with student account) |
| **4.0** | **2024-01-25** | | **Removed student accounts entirely** |
| | | | Simplified auth to 3 roles (ADMIN, TEACHER, HEAD) |
| | | | Added student borrowing flow via teacher |
| | | | Added students table for lookup only (no login) |
| | | | Updated timeline: 10 weeks → 7 weeks |
| | | | Simplified database schema |
| | | | Updated acceptance criteria & user stories |

---

## 15. Appendix

### 15.1 Glossary
| Term | Definition |
|------|------------|
| **ZTNA** | Zero Trust Network Access - Security model requiring verification for every access |
| **QR Code** | Quick Response Code - 2D barcode for storing information |
| **CRUD** | Create, Read, Update, Delete - Basic database operations |
| **RTO** | Recovery Time Objective - Maximum acceptable downtime |
| **RPO** | Recovery Point Objective - Maximum acceptable data loss |
| **RBAC** | Role-Based Access Control - Permission management by user role |
| **JWT** | JSON Web Token - Secure authentication token format |
| **MVP** | Minimum Viable Product - Basic version with core features |
| **PWA** | Progressive Web App - Web app with native app-like features |
| **NIS** | Nomor Induk Siswa - Student ID number |

### 15.2 References
- [Cloudflare Zero Trust Documentation](https://developers.cloudflare.com/cloudflare-one/)
- [Vue.js Documentation](https://vuejs.org/)
- [React Documentation](https://react.dev/)
- [Go Gin Framework](https://gin-gonic.com/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Compose Specification](https://docs.docker.com/compose/)
- [OWASP Security Guidelines](https://owasp.org/www-project-top-ten/)
- [PWA Documentation](https://web.dev/progressive-web-apps/)
- [html5-qrcode Library](https://scanapp.org/)
- [TailwindCSS](https://tailwindcss.com/)

### 15.3 Browser Compatibility Matrix

| Browser | Version | Support | Notes |
|---------|---------|---------|-------|
| Chrome | 90+ | ✅ Full | Recommended |
| Firefox | 88+ | ✅ Full | - |
| Safari | 14+ | ✅ Full | iOS 14+ |
| Edge | 90+ | ✅ Full | Chromium-based |
| Opera | 76+ | ✅ Full | - |
| Samsung Internet | 14+ | ✅ Full | Android |

### 15.4 Quick Student Borrow Flow Diagram

```
┌─────────────┐     ┌─────────────┐
│    Guru     │     │   System    │
│  (Browser)  │     │   (Web)     │
│             │     │             │
└──────┬──────┘     └──────┬──────┘
       │                   │
       │ 1. Klik "Pinjam untuk Siswa"
       │──────────────────>│
       │                   │
       │ 2. Input NIS siswa
       │──────────────────>│
       │                   │
       │ 3. Auto-suggest nama & kelas
       │<──────────────────│
       │                   │
       │ 4. Scan QR barang
       │──────────────────>│
       │                   │
       │ 5. Input tujuan (opsional)
       │──────────────────>│
       │                   │
       │ 6. Due date auto-set (max 6 jam)
       │<──────────────────│
       │                   │
       │ 7. Klik "Konfirmasi"
       │──────────────────>│
       │                   │
       │ 8. Transaksi aktif + notifikasi
       │<──────────────────│
       │                   │
       │ 9. Serahkan barang ke siswa
       │                   │
       │ ... (saat return) │
       │                   │
       │ 10. Scan QR + cek kondisi
       │──────────────────>│
       │                   │
       │ 11. Selesai + update status
       │<──────────────────│
```

### 15.5 Data Flow: Student Borrowing

```
Input: Guru input NIS "12345"
       ↓
API: GET /api/v1/students/search?q=12345
       ↓
DB: SELECT * FROM students WHERE nis LIKE '12345%'
       ↓
Response: [{nis: "12345", name: "Budi Santoso", class: "10 IPA 2"}]
       ↓
UI: Tampilkan hasil, guru konfirmasi
       ↓
Input: Guru scan QR barang "INV-001"
       ↓
API: POST /api/v1/transactions/borrow/student
{
  "student_nis": "12345",
  "student_name": "Budi Santoso",  // snapshot
  "student_class": "10 IPA 2",     // snapshot
  "item_code": "INV-001",
  "purpose": "Praktikum Fisika",
  "borrowed_by": "uuid-guru"
}
       ↓
DB: INSERT INTO transactions (...)
       ↓
Response: {success: true, transaction_id: "uuid-txn"}
       ↓
UI: Tampilkan sukses + detail transaksi
```

---

**Dokumen ini adalah living document dan akan diupdate seiring perkembangan proyek.**

**Last Updated:** 2024
**Next Review:** After Phase 3 Completion (Week 6)

---

## 🎯 Kesimpulan Final: Mengapa Tanpa Akun Siswa?

| Keuntungan | Penjelasan |
|------------|------------|
| ✅ **Lebih Sederhana** | Tidak perlu kelola ratusan akun siswa, password reset, dll |
| ✅ **Lebih Cepat Develop** | Kurangi 30-40% kompleksitas auth & UI |
| ✅ **Lebih Aman** | Tidak ada risiko akun siswa diretas/disalahgunakan |
| ✅ **Lebih Realistis** | Siswa tidak selalu punya akses device/konsisten |
| ✅ **Akuntabilitas Jelas** | Guru yang input = guru yang bertanggung jawab |
| ✅ **Data Tetap Tercatat** | NIS, nama, kelas siswa tetap tersimpan di transaksi |

**Trade-off yang Diterima:**
- ⚠️ Siswa tidak bisa cek riwayat pinjam sendiri (tapi bisa tanya guru)
- ⚠️ Guru perlu 30 detik ekstra untuk input NIS siswa (masih sangat cepat)
- ⚠️ Perlu data siswa (NIS+Nama) di sistem (tapi hanya untuk lookup, bukan auth)

**Rekomendasi Implementasi:**
1. Siapkan file Excel daftar siswa (NIS, Nama, Kelas) untuk di-import ke tabel `students`
2. Latih guru untuk flow "Pinjam untuk Siswa" yang cepat (<30 detik)
3. Sediakan opsi input manual NIS jika auto-suggest tidak ditemukan
4. Tambahkan fitur "Riwayat per Kelas" agar guru mudah pantau siswanya

**Ini adalah versi paling pragmatis untuk MVP sekolah.** 🚀