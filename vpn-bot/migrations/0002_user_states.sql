CREATE TABLE IF NOT EXISTS user_states(
    telegram_id  BIGINT PRIMARY KEY,
    plan_id      TEXT,
    locations    TEXT[] NOT NULL DEFAULT '{}',
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);