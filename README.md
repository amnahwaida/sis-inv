# SIS-INV (School Inventory System)

SIS-INV is a premium school inventory management system designed for transparency, efficiency, and high-performance asset tracking. Built with a modern tech stack (Go/Gin & Vue 3/Vite), it offers a seamless experience for administrators and teachers to manage school facilities and equipment.

---

## ✨ Features

### Phase 1 (Core Functionality)
- **JWT Authentication**: Secure login/logout and session management.
- **RBAC (Role-Based Access Control)**: Granular access for `ADMIN`, `TEACHER`, and `HEAD` roles.
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

### 2. Infrastructure Setup
Spin up the PostgreSQL database container:
```bash
# From the project root
docker compose up -d
```
*The database runs on port `5432` with the name `inventaris_db`.*

### 3. Backend (API) Setup
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```
*API serves at `http://localhost:8080`. Seeding (default admin) happens automatically on first run.*

### 4. Frontend Setup
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

**Full Environment Build:**
```bash
# Use Docker to containerize the entire stack
docker compose build
```

---

## 📂 Project Structure
- `/backend`: Go Gin source code, migration logic, and internal packages.
- `/frontend`: Vue 3 source code, Pinia stores, and assets.
- `docker-compose.yml`: Database and system infrastructure configuration.
- `prd.md`: Product Requirements Document (Historical Reference).
- `CHANGELOG.md`: Detailed history of Phase 2 improvements.

---

## 🔑 Default Credentials
- **Username:** `admin`
- **Password:** `admin123` (Change this in Settings after login!)

---
🚀 *SIS-INV - Empowering School Inventory Transparency.*
