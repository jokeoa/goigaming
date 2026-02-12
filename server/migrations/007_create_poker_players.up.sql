CREATE TABLE poker_players (
    id          UUID          NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    table_id    UUID          NOT NULL REFERENCES poker_tables(id) ON DELETE CASCADE,
    user_id     UUID          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    username    VARCHAR(100)  NOT NULL,
    stack       DECIMAL(15,4) NOT NULL DEFAULT 0 CHECK (stack >= 0),
    seat_number INT           NOT NULL CHECK (seat_number >= 1 AND seat_number <= 9),
    status      VARCHAR(20)   NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'sitting_out', 'all_in', 'folded')),
    joined_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    UNIQUE(table_id, seat_number),
    UNIQUE(table_id, user_id)
);

CREATE INDEX idx_poker_players_table_id ON poker_players(table_id);
CREATE INDEX idx_poker_players_user_id ON poker_players(user_id);
