CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS merchants (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY ,
    balance BIGINT,
    created TIMESTAMP DEFAULT NOW(),
    updated TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transfers (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY ,
    merchant_id uuid NOT NULL,
    merchant_ref_id TEXT NOT NULL,
    bank_ref_id TEXT DEFAULT NULL UNIQUE,
    bank_code TEXT NOT NULL,
    amount BIGINT NOT NULL,
    status TEXT NOT NULL,
    destination_acc_number TEXT NOT NULL,
    created TIMESTAMP DEFAULT NOW(),
    updated TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_transfer_merchant_id FOREIGN KEY (merchant_id) REFERENCES merchants(id),
    CONSTRAINT transfers_merchant_id_merchant_ref_if UNIQUE (merchant_id, merchant_ref_id)
);

CREATE TABLE IF NOT EXISTS ledgers (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    merchant_id UUID NOT NULL,
    transfer_id UUID NOT NULL,
    credit BIGINT NOT NULL,
    debit BIGINT NOT NULL,
    created TIMESTAMP DEFAULT NOW(),
    updated TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_ledgers_merchant_id FOREIGN KEY (merchant_id) REFERENCES merchants(id),
    CONSTRAINT fk_ledgers_transfer_id FOREIGN KEY (transfer_id) REFERENCES transfers(id)
);

INSERT INTO merchants (id, balance) VALUES ('df83833d-e1b5-4493-af61-e8618b688b67', 10000000);