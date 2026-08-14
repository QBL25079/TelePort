-- Пользователи
CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    telegram_id     BIGINT UNIQUE NOT NULL,
    username        TEXT,
    first_name      TEXT,
    last_name       TEXT,
    is_admin        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Тарифы (можно будет потом вынести в админку)
CREATE TABLE plans (
    id              TEXT PRIMARY KEY,           -- '1m', '3m', '6m', '12m'
    name            TEXT NOT NULL,              -- '1 месяц'
    duration_days   INT NOT NULL,
    price           INT NOT NULL,               -- в рублях или звёздах
    is_active       BOOLEAN NOT NULL DEFAULT TRUE
);

-- Локации / страны
CREATE TABLE locations (
    id              TEXT PRIMARY KEY,           -- 'nl', 'de', 'us'
    name            TEXT NOT NULL,              -- 'Нидерланды'
    flag            TEXT NOT NULL,              -- '🇳🇱'
    is_active       BOOLEAN NOT NULL DEFAULT TRUE
);

-- Подписки пользователей
CREATE TABLE subscriptions (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_id         TEXT NOT NULL REFERENCES plans(id),
    status          TEXT NOT NULL CHECK (status IN ('active', 'expired', 'cancelled')),
    happ_link       TEXT,                       -- ссылка для Happ
    starts_at       TIMESTAMPTZ NOT NULL,
    expires_at      TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Какие локации входят в подписку (если пользователь выбирал несколько)
CREATE TABLE subscription_locations (
    subscription_id BIGINT NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    location_id     TEXT NOT NULL REFERENCES locations(id),
    PRIMARY KEY (subscription_id, location_id)
);

-- Платежи
CREATE TABLE payments (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    plan_id         TEXT NOT NULL,
    amount          INT NOT NULL,
    currency        TEXT NOT NULL DEFAULT 'RUB',
    status          TEXT NOT NULL,              -- pending / success / failed
    external_id     TEXT,                       -- ID от платёжки
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индексы
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_users_telegram_id ON users(telegram_id);

-- Начальные данные
INSERT INTO plans (id, name, duration_days, price) VALUES
('1m',  '1 месяц',   30,  199),
('3m',  '3 месяца',  90,  499),
('6m',  '6 месяцев', 180, 899),
('12m', '1 год',     365, 1490)
ON CONFLICT (id) DO NOTHING;

INSERT INTO locations (id, name, flag) VALUES
('nl', 'Нидерланды', '🇳🇱'),
('de', 'Германия',   '🇩🇪'),
('fi', 'Финляндия',  '🇫🇮'),
('us', 'США',        '🇺🇸'),
('lv', 'Латвия',     '🇱🇻')
ON CONFLICT (id) DO NOTHING;