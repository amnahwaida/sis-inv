-- ============================================================================
-- SIS-INV Database Schema - Migration 005
-- ============================================================================

CREATE TABLE IF NOT EXISTS settings (
    key VARCHAR(255) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO settings (key, value) VALUES 
('app_name', 'SIS-INV'),
('app_subtitle', 'Inventaris Sekolah'),
('app_description', 'School Inventory System'),
('app_footer', 'SIS-INV • AMANAH & TERTIB'),
('app_security_notice', 'Crafted with ❤️ by <a href="mailto:vannyezhaa@gmail.com" class="hover:underline text-primary-500">vannyezhaa</a> – Business & Collaboration: +6285159982101')
ON CONFLICT (key) DO NOTHING;
