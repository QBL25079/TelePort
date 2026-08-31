-- locations в платеже + paid_at
ALTER TABLE payments
    ADD COLUMN IF NOT EXISTS locations TEXT[] NOT NULL DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS paid_at TIMESTAMPTZ;