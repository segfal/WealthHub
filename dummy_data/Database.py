import os,json
from dotenv import load_dotenv
import psycopg2 as pg

load_dotenv()



DB_URL = os.getenv("DB_URL")
DB_NAME = os.getenv("DB_NAME")
DB_USER = os.getenv("DB_USER")
DB_PASSWORD = os.getenv("DB_PASSWORD")
DB_HOST = os.getenv("DB_HOST")
DB_PORT = os.getenv("DB_PORT")


def connect_to_db():
    conn = pg.connect(DB_URL)
    return conn

def close_db_connection(conn):
    conn.close()

'''
schema for users table:
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

schema for transactions table:
CREATE TABLE transactions (
    transaction_id VARCHAR(20) PRIMARY KEY,
    account_id VARCHAR(20),
    transaction_type VARCHAR(20),
    transaction_amount DECIMAL(10, 2),
    transaction_currency VARCHAR(3),
    transaction_date TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES users(account_id)
);


tables are already created in the database so all we need to do is insert the data
transaction and user will be loaded from the json files in the data directory and inserted into the database this will do a batch insert

Example:
JillDoe.json
{
    "account": {
        "account_id": 1234567893,
        "account_name": "Chase Total Checking",
        "account_type": "Checking",
        "account_number": "****7893",
        "balance": {
            "current": 15975.65,
            "available": 15975.65,
            "currency": "USD"
        },
        "owner_name": "Jill Doe",
        "bank_details": {
            "bank_name": "Chase Bank",
            "routing_number": "021000021",
            "branch": "Manhattan Main Branch, New York, NY"
        },
        "transactions": [
            {
                "transaction_id": "TXN00001",
                "account_id": 1234567893,
                "date": "2025-01-01",
                "amount": -2216,
                "category": "Rent",
                "merchant": "Park Avenue Apartments",
                "location": "Manhattan, New York, NY",
                "type": "ACH Transfer",
                "status": "Completed",
                "timestamp": "2025-01-22T13:00:52",
                "payment_method": "ACH"
            },
            ]
    }

'''

class UserInput:
    def __init__(self, json_file):
        self.json_file = json_file
        self.data = self.load_data()
        self.insert_user()
        self.insert_transactions()

    def load_data(self):
        with open(self.json_file, 'r') as file:
            return json.load(file)
        
    def insert_user(self):
        conn = connect_to_db()
        cursor = conn.cursor()
        cursor.execute("INSERT INTO users (account_id, account_name, account_type, account_number, balance_current, balance_available, balance_currency, owner_name) VALUES (%s, %s, %s, %s, %s, %s, %s, %s)", (self.data['account']['account_id'], self.data['account']['account_name'], self.data['account']['account_type'], self.data['account']['account_number'], self.data['account']['balance']['current'], self.data['account']['balance']['available'], self.data['account']['balance']['currency'], self.data['account']['owner_name']))
        conn.commit()
        cursor.close()
        close_db_connection(conn)
    

    def insert_transactions(self):
        # buik insert all transactions in a batch using upsert
        for transaction in self.data['account']['transactions']:
            conn = connect_to_db()
            cursor = conn.cursor()
            cursor.execute("INSERT INTO transactions (account_id, transaction_type, transaction_amount, transaction_currency, transaction_date) VALUES (%s, %s, %s, %s, %s)", (self.data['account']['account_id'], transaction['type'], transaction['amount'], transaction['currency'], transaction['date']))
            conn.commit()
        cursor.close()
        close_db_connection(conn)
        

    
def store_data(json_file):
    user_input = UserInput(json_file)
    user_input.insert_user()
    user_input.insert_transactions()

