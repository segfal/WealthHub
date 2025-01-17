-- Drop tables if they exist
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS bank_details;
DROP TABLE IF EXISTS users;

-- Create users table
CREATE TABLE users (
    account_id VARCHAR(20) PRIMARY KEY,
    account_name VARCHAR(50),
    account_type VARCHAR(20),
    account_number VARCHAR(20),
    balance_current DECIMAL(10, 2),
    balance_available DECIMAL(10, 2),
    balance_currency VARCHAR(3),
    owner_name VARCHAR(50)
);

-- Create bank_details table
CREATE TABLE bank_details (
    account_id VARCHAR(20) PRIMARY KEY REFERENCES users(account_id),
    bank_name VARCHAR(50),
    routing_number VARCHAR(20),
    branch VARCHAR(100)
);

-- Create transactions table
CREATE TABLE transactions (
    transaction_id VARCHAR(20) PRIMARY KEY,
    account_id VARCHAR(20) REFERENCES users(account_id),
    date TIMESTAMP,
    amount DECIMAL(10, 2),
    category VARCHAR(50),
    merchant VARCHAR(50),
    location VARCHAR(100)
);
