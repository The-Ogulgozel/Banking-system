CREATE TABLE IF NOT EXISTS transactions (
    id                BIGSERIAL PRIMARY KEY,
    account_id        BIGINT NOT NULL REFERENCES accounts(id),
    to_account_id     BIGINT NULL REFERENCES accounts(id),
    amount            BIGINT NOT NULL CHECK (amount > 0),
    transaction_type  VARCHAR(20) NOT NULL CHECK (transaction_type IN ('deposit', 'withdraw', 'transfer')),
    created_at        TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_account_created ON transactions (account_id, created_at DESC);