-- Migration 003: Add locations table and formalize location management

CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Extract unique locations from existing items and populate the new table
INSERT INTO locations (name)
SELECT DISTINCT location 
FROM items 
WHERE location IS NOT NULL AND location <> ''
ON CONFLICT (name) DO NOTHING;

-- Add location_id column to items
ALTER TABLE items ADD COLUMN IF NOT EXISTS location_id INT REFERENCES locations(id);

-- Link items to the new locations table based on their string location names
UPDATE items i
SET location_id = l.id
FROM locations l
WHERE i.location = l.name;

-- Note: We keep the 'location' string column for now to prevent data loss 
-- if something goes wrong, but we will primarily use location_id from now on.
CREATE INDEX IF NOT EXISTS idx_items_location_id ON items(location_id);
