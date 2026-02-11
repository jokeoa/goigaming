CREATE TABLE poker_hands (
    id              UUID          NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    table_id        UUID          NOT NULL REFERENCES poker_tables(id) ON DELETE CASCADE,
    hand_number     INT           NOT NULL,
    pot             DECIMAL(15,4) NOT NULL DEFAULT 0,
    community_cards TEXT          NOT NULL DEFAULT '',
    stage           VARCHAR(20)   NOT NULL DEFAULT 'waiting' CHECK (stage IN ('waiting', 'preflop', 'flop', 'turn', 'river', 'showdown', 'complete')),
    winner_id       UUID          REFERENCES poker_players(id) ON DELETE SET NULL,
    started_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    ended_at        TIMESTAMPTZ,
    UNIQUE(table_id, hand_number)
);

CREATE INDEX idx_poker_hands_table_id ON poker_hands(table_id);
