CREATE TABLE IF NOT EXISTS sport_events (
    id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    sport VARCHAR(50) NOT NULL,
    league VARCHAR(100) NOT NULL,
    home_team VARCHAR(100) NOT NULL,
    away_team VARCHAR(100) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    home_odds NUMERIC(10,2) NOT NULL,
    draw_odds NUMERIC(10,2),
    away_odds NUMERIC(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'upcoming' CHECK (status IN ('upcoming', 'live', 'finished', 'cancelled')),
    result VARCHAR(10) CHECK (result IN ('home', 'draw', 'away')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    settled_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_sport_events_status ON sport_events(status);
CREATE INDEX IF NOT EXISTS idx_sport_events_start_time ON sport_events(start_time);
