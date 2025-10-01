CREATE TABLE IF NOT EXISTS wallets (
--     id SERIAL PRIMARY KEY,
    valletId UUID NOT NULL UNIQUE,
    balance NUMERIC(18, 2) NOT NULL DEFAULT 0 CHECK (balance >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_wallets_valletId ON wallets(valletId);

CREATE TABLE IF NOT EXISTS wallet_transactions (
    id SERIAL PRIMARY KEY,
    valletId UUID NOT NULL,
    operation_type VARCHAR(10) NOT NULL CHECK (operation_type IN ('DEPOSIT', 'WITHDRAW')),
    amount NUMERIC(18, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_wallet
    FOREIGN KEY(valletId) REFERENCES wallets(valletId) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_wallet_transactions_valletId ON wallet_transactions(valletId);



-- Тестовые данные
INSERT INTO wallets (valletId, balance) VALUES
    ('11111111-1111-1111-1111-111111111111', 1000.00),
    ('22222222-2222-2222-2222-222222222222', 500.50),
    ('33333333-3333-3333-3333-333333333333', 0.00),
    ('44444444-4444-4444-4444-444444444444', 250.75);