CREATE TABLE wallets (
    user_id    UUID         NOT NULL PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    balance    DECIMAL(15,4) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    version    INT          NOT NULL DEFAULT 1,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
