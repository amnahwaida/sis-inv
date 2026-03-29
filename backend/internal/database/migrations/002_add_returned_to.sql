DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM pg_attribute WHERE attrelid = 'transactions'::regclass AND attname = 'returned_to') THEN
        ALTER TABLE transactions ADD COLUMN returned_to UUID REFERENCES users(id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_transactions_returned_to ON transactions(returned_to);
