CREATE TABLE accounts (
    account_id VARCHAR(20) PRIMARY KEY,
    account_name VARCHAR(50),
    account_type VARCHAR(20),
    account_number VARCHAR(10),
    balance_current DECIMAL(10, 2),
    balance_available DECIMAL(10, 2),
    currency VARCHAR(10),
    owner_name VARCHAR(50),
    bank_name VARCHAR(50),
    routing_number VARCHAR(10),
    branch VARCHAR(100)
);



CREATE TABLE transactions (
    transaction_id VARCHAR(20) PRIMARY KEY,
    account_id VARCHAR(20),
    date DATE,
    amount DECIMAL(10, 2),
    category VARCHAR(50),
    merchant VARCHAR(50),
    location VARCHAR(100)
);
