CREATE TABLE poker_tables (
    id          UUID          NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    name        VARCHAR(100)  NOT NULL,
    small_blind DECIMAL(15,4) NOT NULL CHECK (small_blind > 0),
    big_blind   DECIMAL(15,4) NOT NULL CHECK (big_blind > 0),
    min_buy_in  DECIMAL(15,4) NOT NULL CHECK (min_buy_in > 0),
    max_buy_in  DECIMAL(15,4) NOT NULL CHECK (max_buy_in > 0),
    max_players INT           NOT NULL CHECK (max_players >= 2 AND max_players <= 9),
    status      VARCHAR(20)   NOT NULL DEFAULT 'waiting' CHECK (status IN ('waiting', 'active', 'closed')),
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_poker_tables_status ON poker_tables(status);
