# Changelog: SIS-INV Phase 2

Detailed summary of key improvements and final premium overhaul features.

---

## [Phase 2: Premium & Branding Improvements] — 2026-04-03

### ✨ Core Improvements
- **Dynamic Branding System**:
  - Implemented administrative controls for dynamic Logo and Favicon uploads.
  - Refactored all application icons to use **Heroicons v2 Outline** for consistency and performance.
  - Enabled dynamic system-wide application titles and branding-specific default icons.
- **Advanced Dynamic PWA**:
  - Transitioned from static `manifest.json` to dynamic backend-served `/manifest.webmanifest`.
  - Added support for administrator-controlled PWA icons (192x192 & 512x512).
  - PWA theme color and metadata now automatically sync with system branding settings.
- **Multi-Theme Engine**:
  - Implemented 5 premium color themes: **Sapphire**, **Emerald**, **Royal**, **Crimson**, and **Amber**.
  - Refactored Tailwind configuration to use CSS variables and RGB-based primary colors for flexible opacity modifiers.
- **Sidebar UX Overhaul**:
  - Reorganized entire navigation into collapsible categories with refined grouping: **Dashboard**, **Transaksi**, **Inventaris**, **Pemeliharaan**, **Admin**, **Laporan**, **Sistem**.
  - Added specific category icons to group headers for better visual hierarchy.
  - Implemented persistent sidebar state using local storage.
  - Optimized expand/collapse animations with `cubic-bezier`.
  - Enforced `whitespace-nowrap` on all menu items to eliminate inconsistent text wrapping.

### 🐛 Bug Fixes & Refinements
- **UI Consistency**:
  - Standardized all 18 views with consistent typography, spacing, and glassmorphism styling.
  - Fixed color mapping issues in `BorrowView.vue` choice-box components.
  - Optimized sidebar labels for cleaner, single-line presentation.
- **Performance**:
  - Implemented high-performance, server-side pagination for all data-intensive modules (Items, Students, Users).
  - Reduced frontend footprint by replacing heavy dependencies with lightweight Heroicons SVGs.

---
🚀 *SIS-INV: Driving School Asset Management to New Heights.*
