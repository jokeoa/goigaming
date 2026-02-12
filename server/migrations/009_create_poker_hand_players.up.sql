CREATE TABLE poker_hand_players (
    id          UUID          NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    hand_id     UUID          NOT NULL REFERENCES poker_hands(id) ON DELETE CASCADE,
    player_id   UUID          NOT NULL REFERENCES poker_players(id) ON DELETE CASCADE,
    hole_cards  TEXT          NOT NULL DEFAULT '',
    bet_amount  DECIMAL(15,4) NOT NULL DEFAULT 0,
    last_action VARCHAR(20)   NOT NULL DEFAULT '',
    is_active   BOOLEAN       NOT NULL DEFAULT TRUE,
    UNIQUE(hand_id, player_id)
);

CREATE INDEX idx_poker_hand_players_hand_id ON poker_hand_players(hand_id);
