-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets(
    id serial PRIMARY KEY,
    address text,
    balance int,
    constraint balance_nonnegative CHECK (balance >= 0),
    constraint address_unique UNIQUE (address)
);
CREATE TABLE transactions(
    id serial PRIMARY KEY,
    sender text REFERENCES wallets (address),
    receiver text REFERENCES wallets (address),
    amount int,
    constraint amount_nonnegative check (amount >= 0)
);
CREATE PROCEDURE transfer(
   sender TEXT,
   receiver TEXT, 
   amount INT
)
LANGUAGE plpgsql
AS $$
DECLARE
    sender_balance INT;
    sender_exists BOOLEAN;
    receiver_exists BOOLEAN;
BEGIN
    SELECT EXISTS (SELECT 1 FROM wallets WHERE address = sender) INTO sender_exists;
    SELECT EXISTS (SELECT 1 FROM wallets WHERE address = receiver) INTO receiver_exists;

    IF NOT sender_exists THEN
        RAISE EXCEPTION 'Invalid wallet' USING ERRCODE = 'P0010';
    END IF;

    IF NOT receiver_exists THEN
        RAISE EXCEPTION 'Invalid wallet' USING ERRCODE = 'P0010';
    END IF;

    SELECT balance INTO sender_balance
    FROM wallets
    WHERE address = sender
    FOR UPDATE;

    IF sender_balance < amount THEN
        RAISE EXCEPTION 'Insufficient funds' USING ERRCODE = 'P0011';
    END IF;

    UPDATE wallets 
    SET balance = balance - amount 
    WHERE address = sender;

    UPDATE wallets 
    SET balance = balance + amount 
    WHERE address = receiver;

    INSERT INTO transactions (sender, receiver, amount)
    VALUES (sender, receiver, amount);
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
DROP TABLE wallets;
DROP PROCEDURE transfer;
-- +goose StatementEnd
