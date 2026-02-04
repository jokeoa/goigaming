-- Таблица игровых столов для руллетки
CREATE TABLE IF NOT EXISTS roulette_tables (
    id         UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    name       VARCHAR(100) NOT NULL UNIQUE,
    min_bet    NUMERIC(10,2) NOT NULL CHECK (min_bet > 0),
    max_bet    NUMERIC(10,2) NOT NULL CHECK (max_bet > 0),
    status     VARCHAR(20)  NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'maintenance')),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Таблица раундов (оборотов колеса)
CREATE TABLE IF NOT EXISTS roulette_rounds (
    id              UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    table_id        UUID        NOT NULL REFERENCES roulette_tables(id) ON DELETE CASCADE,
    round_number    INT         NOT NULL,
    result          INT         CHECK (result IS NULL OR (result >= 0 AND result <= 36)),
    result_color    VARCHAR(10) CHECK (result_color IS NULL OR result_color IN ('red', 'black', 'green')),
    seed_hash       VARCHAR(255),
    seed_revealed   VARCHAR(255),
    betting_ends_at TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    settled_at      TIMESTAMPTZ,
    UNIQUE (table_id, round_number)
);

-- Таблица ставок
CREATE TABLE IF NOT EXISTS roulette_bets (
    id         UUID          NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    round_id   UUID          NOT NULL REFERENCES roulette_rounds(id) ON DELETE CASCADE,
    user_id    UUID          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bet_type   VARCHAR(20)   NOT NULL CHECK (bet_type IN ('straight', 'split', 'street', 'corner', 'line', 'dozen', 'column', 'red', 'black', 'odd', 'even', 'high', 'low')),
    bet_value  VARCHAR(50)   NOT NULL,
    amount     NUMERIC(10,2) NOT NULL CHECK (amount > 0),
    payout     NUMERIC(10,2) DEFAULT 0,
    status     VARCHAR(20)   NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'won', 'lost', 'cancelled')),
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_roulette_rounds_table_id ON roulette_rounds(table_id);
CREATE INDEX IF NOT EXISTS idx_roulette_rounds_status ON roulette_rounds(table_id, settled_at);
CREATE INDEX IF NOT EXISTS idx_roulette_bets_round_id ON roulette_bets(round_id);
CREATE INDEX IF NOT EXISTS idx_roulette_bets_user_id ON roulette_bets(user_id);
