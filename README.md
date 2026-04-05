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

## 🚀 Getting Started

### 1. Prerequisites
Ensure you have the following installed:
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Go 1.22+](https://go.dev/dl/)
- [Node.js 22.12+](https://nodejs.org/)

### 2. Environment Configuration
Copy the example environment file and update the variables:
```bash
cp .env.example .env
```

**Key variables to update in `.env`:**
- `JWT_SECRET`: Change this to a random long string for security.
- `TUNNEL_TOKEN`: Your Cloudflare Tunnel token (required for `docker-compose` tunnel service).
- `DB_PASSWORD`: Set a strong password for your PostgreSQL database.
- `ADMIN_DEFAULT_PASSWORD`: Set the default password used for administrative resets.

### 3. Infrastructure Setup
Spin up the PostgreSQL database and Cloudflare Tunnel:
```bash
# From the project root
docker compose up -d
```
*The database runs on port `5432` by default. If the `TUNNEL_TOKEN` is provided, the tunnel will start automatically.*

### 4. Backend (API) Setup
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```
*API serves at `http://localhost:8080`.*

### 5. Frontend Setup
```bash
cd frontend
npm install
npm run dev
```
*Frontend serves at `http://localhost:5173`. CORS is pre-configured for this port.*

---

## 📦 Build & Deploy

### Production Build
To create highly optimized production bundles:

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

### 1. Database Backup
```bash
./scripts/backup.sh
```
- **What it does**: Creates a full, compressed SQL dump of the database using `pg_dump` inside the Docker container.
- **Result**: A new compressed file appears in the `/backups` directory (e.g., `backups/db_20240405_120000.sql.gz`).
- **Retention**: It automatically purges backups older than 30 days to save disk space.

### 2. Database Restore
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
- **Credentials**: Uses the `ADMIN_DEFAULT_USERNAME` and `ADMIN_DEFAULT_PASSWORD` defined in your `.env` file.
- **Usage**: Use this as a "Safety Valve" if you accidentally demote yourself or lose admin access.

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
