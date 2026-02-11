CREATE TABLE IF NOT EXISTS sport_bets (
    id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES sport_events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bet_type VARCHAR(20) NOT NULL CHECK (bet_type IN ('home', 'draw', 'away')),
    amount NUMERIC(10,2) NOT NULL CHECK (amount > 0),
    odds NUMERIC(10,2) NOT NULL,
    payout NUMERIC(10,2) DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'won', 'lost', 'cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sport_bets_event_id ON sport_bets(event_id);
CREATE INDEX IF NOT EXISTS idx_sport_bets_user_id ON sport_bets(user_id);
CREATE INDEX IF NOT EXISTS idx_sport_bets_status ON sport_bets(status);
