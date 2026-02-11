-- Таблица спортивных событий
CREATE TABLE IF NOT EXISTS sport_events (
    id              UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    sport_type      VARCHAR(50) NOT NULL CHECK (sport_type IN ('football', 'basketball', 'tennis', 'hockey', 'esports')),
    league          VARCHAR(100) NOT NULL,
    home_team       VARCHAR(100) NOT NULL,
    away_team       VARCHAR(100) NOT NULL,
    home_odds       NUMERIC(5,2) NOT NULL CHECK (home_odds > 0),
    draw_odds       NUMERIC(5,2) CHECK (draw_odds IS NULL OR draw_odds > 0),
    away_odds       NUMERIC(5,2) NOT NULL CHECK (away_odds > 0),
    event_time      TIMESTAMPTZ NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'live', 'finished', 'cancelled')),
    home_score      INT DEFAULT 0,
    away_score      INT DEFAULT 0,
    created_by      UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    settled_at      TIMESTAMPTZ
);

-- Таблица ставок на спорт
CREATE TABLE IF NOT EXISTS sport_bets (
    id              UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    event_id        UUID        NOT NULL REFERENCES sport_events(id) ON DELETE CASCADE,
    user_id         UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bet_type        VARCHAR(20) NOT NULL CHECK (bet_type IN ('home', 'draw', 'away')),
    odds            NUMERIC(5,2) NOT NULL,
    amount          NUMERIC(10,2) NOT NULL CHECK (amount > 0),
    potential_win   NUMERIC(10,2) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'won', 'lost', 'cancelled', 'refunded')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    settled_at      TIMESTAMPTZ
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_sport_events_status ON sport_events(status);
CREATE INDEX IF NOT EXISTS idx_sport_events_time ON sport_events(event_time);
CREATE INDEX IF NOT EXISTS idx_sport_events_sport_type ON sport_events(sport_type);
CREATE INDEX IF NOT EXISTS idx_sport_bets_event_id ON sport_bets(event_id);
CREATE INDEX IF NOT EXISTS idx_sport_bets_user_id ON sport_bets(user_id);
CREATE INDEX IF NOT EXISTS idx_sport_bets_status ON sport_bets(status);
