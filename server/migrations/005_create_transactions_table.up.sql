CREATE TABLE transactions (
    id             UUID          NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    wallet_id      UUID          NOT NULL REFERENCES wallets(user_id) ON DELETE CASCADE,
    amount         DECIMAL(15,4) NOT NULL,
    balance_after  DECIMAL(15,4) NOT NULL,
    reference_type VARCHAR(50)   NOT NULL,
    reference_id   UUID,
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_wallet_id ON transactions(wallet_id);
CREATE INDEX idx_transactions_reference ON transactions(reference_type, reference_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
