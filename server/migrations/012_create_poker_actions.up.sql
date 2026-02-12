CREATE TABLE poker_actions (
    id           UUID          NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    hand_id      UUID          NOT NULL REFERENCES poker_hands(id) ON DELETE CASCADE,
    player_id    UUID          NOT NULL REFERENCES poker_players(id) ON DELETE CASCADE,
    action       VARCHAR(20)   NOT NULL CHECK (action IN ('fold', 'check', 'call', 'raise', 'all_in', 'bet')),
    amount       DECIMAL(15,4) NOT NULL DEFAULT 0,
    stage        VARCHAR(20)   NOT NULL,
    action_order INT           NOT NULL,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_poker_actions_hand_id ON poker_actions(hand_id);
