# SIS-INV (School Inventory System)

SIS-INV is a premium school inventory management system designed for transparency, efficiency, and high-performance asset tracking. Built with a modern tech stack (Go/Gin & Vue 3/Vite), it offers a seamless experience for administrators and teachers to manage school facilities and equipment.

---

## ✨ Features

### Phase 1 (Core Functionality)
- **JWT Authentication**: Secure login/logout and session management.
- **RBAC (Role-Based Access Control)**: Granular access for `ADMIN` and `TEACHER` roles.
- **Data Management**: Complete CRUD for Items, Categories, Locations, Students, and Users.
- **Borrowing Workflow**: End-to-end tracking of items borrowed by students/staff.
- **Audit Logs**: Deep 5W1H audit trail (Who, What, Where, When, Why, How) for every critical action.
- **QR Code Support**: Generate and print QR codes for effortless item identification.

### Phase 2 (Premium Overhaul)
- **Dynamic Branding**: Fully white-label your app directly from the UI.
  - Custom Logos & Favicons (with auto-sync to PWA icons).
  - 15+ Customizable Menu Icons (Heroicons based).
  - Dynamic Custom App Titles.
- **Advanced Theme System**: Switch between premium themes:
  - 🟦 **Sapphire** (Deep Blue)
  - 🟩 **Emerald** (Forest Green)
  - 🟪 **Royal** (Deep Purple)
  - 🟥 **Crimson** (Ruby Red)
  - 🟧 **Amber** (Golden Sunset)
- **Dynamic PWA Support**: 
  - Install the app on desktop and mobile devices.
  - Manifest data (name, theme color, icons) updates dynamically via the backend.
- **Reorganized Sidebar**:
  - Collapsible, grouped navigation for better UI clarity.
  - Smooth expansion animations and persistent state.
- **Premium UI/UX**:
  - Glassmorphism effects and modern dark mode.
  - High-performance server-side pagination for all lists.

---

## 🛠️ Technology Stack

- **Backend**: [Go](https://go.dev/) (Gin Gonic, GORM, JWT)
- **Frontend**: [Vue 3](https://vuejs.org/) (Vite, Pinia, Vue Router)
- **Styling**: [Tailwind CSS](https://tailwindcss.com/)
- **Database**: [PostgreSQL 15+](https://www.postgresql.org/)
- **Infrastructure**: [Docker](https://www.docker.com/) & Docker Compose

---


## 🐳 Docker Deployment (Recommended)

The easiest way to deploy the complete SIS-INV stack (Backend, Frontend, Database, and Backups) is using Docker Compose.

### 1. Configure the Environment
Ensure your `.env` file is ready with all necessary variables:
```bash
cp .env.example .env
# Edit .env and set:
# DB_PASSWORD, JWT_SECRET, and TUNNEL_TOKEN
```

### 2. Launch the Full Stack
Run the following command from the project root to build and start all services in the background:
```bash
docker compose up -d --build
```

**This will spin up 5 active containers:**
- `sisinv_db`: PostgreSQL 15 database.
- `sisinv_backend`: The Go API server.
- `sisinv_frontend`: The Vue 3 application (served via Nginx).
- `sisinv_backup`: Automated daily backup service.
- `cloudflared_tunnel`: Local-to-Internet tunnel via Cloudflare.

### 3. Verification
Check the status of your services:
```bash
docker compose ps
```
The application will be accessible at `http://localhost` (port 80) or via your custom Cloudflare domain.

### 4. Updating the Application
When you pull new changes from GitHub, simply rebuild and restart the containers:
```bash
git pull origin main
docker compose up -d --build
```

---

## 🛠️ Manual Development Setup

If you prefer to run services individually for development:

### 1. Prerequisites
- [Go 1.25+](https://go.dev/dl/)
- [Node.js 22.12+](https://nodejs.org/)
- [PostgreSQL 15+](https://www.postgresql.org/)

### 2. Infrastructure (Database Only)
You can still use Docker just for the database service:
```bash
docker compose up -d db
```

### 3. Backend (API) Setup
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### 4. Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

---

## 📦 Production Builds (Manual)

To manually create optimized production bundles without Docker:

**Frontend Build:**
```bash
cd frontend
npm run build
```
*Produces static files in `/frontend/dist`.*

**Backend Build:**
```bash
cd backend
go build -o server cmd/server/main.go
```

---

## 🛠️ Maintenance & Safety

Powerful terminal-based scripts for system maintenance, located in the `/scripts` directory. **Run these scripts from the project root.**

### 1. Database Backup (Standard/Rotating)
```bash
./scripts/backup.sh
```
- **What it does**: Creates a full, compressed SQL dump using `pg_dump` inside the Docker container.
- **Result**: Files are saved in `/backups/` (e.g., `backups/db_20240405_120000.sql.gz`).
- **Retention**: **Automatic 30-day rotation**. Older files are purged daily to save space. Best for automated daily cron jobs.

### 2. Database Backup (Persistent/Archive)
```bash
./scripts/backup_full.sh
```
- **What it does**: Performs a complete database dump intended for permanent storage.
- **Result**: Files are saved in a separate `/backups/archives/` directory.
- **Retention**: **NO automatic deletion**. These files will remain on your disk indefinitely until you manually remove them. Ideal for backing up before major system changes or for historical records.

### 3. Database Restore
```bash
./scripts/restore.sh backups/your_backup_file.sql.gz
```
- **What it does**: Streams the contents of a backup file back into the active database container.
- **Requirement**: Requires the path to a valid `.sql` or `.sql.gz` file.
- **Safety**: Prompts for manual confirmation before the process begins. **Caution: This overwrites all current data.**

### 3. Emergency Admin Reset
```bash
./scripts/reset_admin.sh
```
- **What it does**: Forcefully promotes the default user to `ADMIN` and resets their access.
- **Docker Support**: **Fully Docker-native.** If your containers are running, it will execute the reset directly inside the `sisinv_backend` container without requiring Go to be installed on your computer.
- **Credentials**: Uses the `ADMIN_DEFAULT_USERNAME` and `ADMIN_DEFAULT_PASSWORD` defined in your `.env` file.

---

## 📂 Project Structure
- `/backend`: Go Gin source code, migration logic, and internal API packages.
- `/frontend`: Vue 3 source code, Pinia stores, and PWA assets.
- `/scripts`: Shell and Go-based maintenance tools.
- `/backups`: Local storage for automated and manual database snapshots.
- `docker-compose.yml`: Database and system infrastructure configuration.

---

## 🔑 Default Credentials
Default access is determined by the `ADMIN_DEFAULT_USERNAME` and `ADMIN_DEFAULT_PASSWORD` keys in your **`.env`** file. Please refer to that file for your initial login details or after running a reset.

---
🚀 *SIS-INV - Empowering School Inventory Transparency.*
