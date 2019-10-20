CREATE TABLE IF NOT EXISTS currencies (
  id       serial,
  ticker   text    NOT NULL,
  decimals integer NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS accounts (
  id          serial,
  owner       text    NOT NULL,
  balance     bigint  NOT NULL DEFAULT 0,
  currency_id integer NOT NULL,
  PRIMARY KEY (id)
);

ALTER TABLE accounts DROP CONSTRAINT IF EXISTS accounts_to_currency_foreign_key;
ALTER TABLE accounts
  ADD CONSTRAINT accounts_to_currency_foreign_key FOREIGN KEY (currency_id)
REFERENCES currencies (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE accounts DROP CONSTRAINT IF EXISTS accounts_amount_check;
ALTER TABLE accounts
  ADD CONSTRAINT accounts_amount_check CHECK (
  balance >= 0);

CREATE TABLE IF NOT EXISTS payments (
  id              serial,
  account_from_id integer                  NOT NULL,
  account_to_id   integer                  NOT NULL,
  amount          bigint                   NOT NULL,
  currency_id     integer                  NOT NULL,
  date            timestamp with time zone NOT NULL,
  PRIMARY KEY (id)
);

ALTER TABLE payments
  DROP CONSTRAINT IF EXISTS payments_to_currency_foreign_key;
ALTER TABLE payments
  ADD CONSTRAINT payments_to_currency_foreign_key FOREIGN KEY (currency_id)
REFERENCES currencies (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE payments
  DROP CONSTRAINT IF EXISTS payments_to_account_from_foreign_key;
ALTER TABLE payments
  ADD CONSTRAINT payments_to_account_from_foreign_key FOREIGN KEY (account_from_id)
REFERENCES accounts (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE payments
  DROP CONSTRAINT IF EXISTS payments_to_account_to_foreign_key;
ALTER TABLE payments
  ADD CONSTRAINT payments_to_account_to_foreign_key FOREIGN KEY (account_to_id)
REFERENCES accounts (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE payments
  DROP CONSTRAINT IF EXISTS payment_amount_check;
ALTER TABLE payments
  ADD CONSTRAINT payment_amount_check CHECK (
  amount > 0);