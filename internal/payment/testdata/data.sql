INSERT INTO currencies (ticker, decimals)
VALUES ('USD', 2);
INSERT INTO currencies (ticker, decimals)
VALUES ('EUR', 2);


INSERT INTO accounts (owner, balance, currency_id) SELECT
'Alice', 2000, c.id FROM currencies c where ticker = 'USD';
INSERT INTO accounts (owner, balance, currency_id) SELECT
'Alice', 3000, c.id FROM currencies c where ticker = 'EUR';
INSERT INTO accounts (owner, balance, currency_id) SELECT
'Bob', 300, c.id FROM currencies c where ticker = 'USD';

