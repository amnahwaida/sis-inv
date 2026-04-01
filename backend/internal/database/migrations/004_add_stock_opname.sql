-- ============================================================================
-- SIS-INV Migration 004: Stock Opname & Borrow Condition Photo
-- ============================================================================

-- 1. Add borrow_photo_url to transactions
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS borrow_photo_url VARCHAR(500);

-- 2. Create audit_sessions table
CREATE TABLE IF NOT EXISTS audit_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id INT REFERENCES locations(id),
    user_id UUID REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CLOSED')),
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Create audit_items table
CREATE TABLE IF NOT EXISTS audit_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES audit_sessions(id) ON DELETE CASCADE,
    item_id UUID REFERENCES items(id),
    found_condition VARCHAR(20) CHECK (found_condition IN ('GOOD', 'DAMAGED', 'LOST', 'IN_REPAIR')),
    scanned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    notes TEXT
);

-- Index for performance
CREATE INDEX IF NOT EXISTS idx_audit_sessions_user ON audit_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_sessions_location ON audit_sessions(location_id);
CREATE INDEX IF NOT EXISTS idx_audit_items_session ON audit_items(session_id);
CREATE INDEX IF NOT EXISTS idx_audit_items_item ON audit_items(item_id);
