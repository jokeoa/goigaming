-- Таблица игровых сессий (для покера и других игр)
CREATE TABLE IF NOT EXISTS game_sessions (
    id              UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    game_type       VARCHAR(20) NOT NULL CHECK (game_type IN ('poker', 'roulette', 'sports_betting')),
    table_name      VARCHAR(100),
    config          JSONB       NOT NULL DEFAULT '{}',
    state           JSONB       NOT NULL DEFAULT '{}',
    status          VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('waiting', 'active', 'finished', 'cancelled')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at       TIMESTAMPTZ
);

-- Таблица участников сессий
CREATE TABLE IF NOT EXISTS game_session_participants (
    id              UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    session_id      UUID        NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    user_id         UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seat_position   INT,
    buy_in_amount   NUMERIC(10,2) NOT NULL CHECK (buy_in_amount >= 0),
    cash_out_amount NUMERIC(10,2) DEFAULT 0,
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    left_at         TIMESTAMPTZ,
    UNIQUE(session_id, user_id),
    UNIQUE(session_id, seat_position)
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_game_sessions_game_type ON game_sessions(game_type);
CREATE INDEX IF NOT EXISTS idx_game_sessions_status ON game_sessions(status);
CREATE INDEX IF NOT EXISTS idx_game_session_participants_session_id ON game_session_participants(session_id);
CREATE INDEX IF NOT EXISTS idx_game_session_participants_user_id ON game_session_participants(user_id);
