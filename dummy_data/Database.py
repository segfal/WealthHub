import psycopg2 as pg
import json
from data_utils.util import database_credentials, connect_to_db, close_db_connection, get_last_transaction_id, split_transaction_key, get_next_transaction_id
from datetime import datetime



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

def bulk_insert_transactions(transactions):
    """Bulk insert transactions into the database"""
    if not transactions:
        return
        
    conn = connect_to_db()
    cursor = conn.cursor()
    
    try:
        # Prepare the values for bulk insert
        values = [(
            t['transaction_id'],
            t['account_id'],
            t['date'],
            t['amount'],
            t['category'],
            t['merchant'],
            t['location']
        ) for t in transactions]
        
        # Execute bulk insert
        cursor.executemany(
            "INSERT INTO transactions (transaction_id, account_id, date, amount, category, merchant, location) "
            "VALUES (%s, %s, %s, %s, %s, %s, %s)",
            values
        )
        conn.commit()
        print(f"Bulk inserted {len(transactions)} transactions")
    except Exception as e:
        print(f"Error in bulk insert: {str(e)}")
        conn.rollback()
    finally:
        cursor.close()
        close_db_connection(conn)

class UserInput:
    def __init__(self, json_file):
        self.json_file = json_file
        self.data = self.load_data()
        self.insert_user()
        self.insert_transactions()

    def get_last_transaction_key(self):
        conn = connect_to_db()
        cursor = conn.cursor()
        cursor.execute("SELECT transaction_id FROM transactions ORDER BY date DESC LIMIT 1")
        result = cursor.fetchone()
        close_db_connection(conn)
        return result[0] if result else None

    def load_data(self):
        with open(self.json_file, 'r') as file:
            return json.load(file)
    def check_user_exists(self):
        ## if user already exists, skip
        conn = connect_to_db()
        cursor = conn.cursor()
        # account_id is an integer
        cursor.execute("SELECT * FROM users WHERE account_id = %s", (str(self.data['account']['account_id']),))
        result = cursor.fetchone()
        close_db_connection(conn)
        return result
    def insert_user(self):
        ## if user already exists, skip
        if self.check_user_exists():
            print(f"User {self.data['account']['account_id']} already exists, skipping")
            return
        conn = connect_to_db()
        cursor = conn.cursor()
        cursor.execute("INSERT INTO users (account_id, account_name, account_type, account_number, balance_current, balance_available, balance_currency, owner_name) VALUES (%s, %s, %s, %s, %s, %s, %s, %s)", (self.data['account']['account_id'], self.data['account']['account_name'], self.data['account']['account_type'], self.data['account']['account_number'], self.data['account']['balance']['current'], self.data['account']['balance']['available'], self.data['account']['balance']['currency'], self.data['account']['owner_name']))
        conn.commit()
        cursor.close()
        close_db_connection(conn)
    
    def split_key(self, key):
        split_key = key.split("TXN")
        last_transaction_key_number = int(split_key[1])
        return last_transaction_key_number
    
    
    def insert_transactions(self):
        """Insert transactions in weekly batches"""
        if not self.data['account']['transactions']:
            return

        # Group transactions by week
        transactions_by_week = {}
        for transaction in self.data['account']['transactions']:
            # Get the week number for this transaction
            date = datetime.strptime(transaction['date'], '%Y-%m-%d')
            week_key = date.strftime('%Y-%W')  # Year and week number
            
            if week_key not in transactions_by_week:
                transactions_by_week[week_key] = []
            transactions_by_week[week_key].append(transaction)

        # Insert each week's transactions as a batch
        for week, transactions in transactions_by_week.items():
            print(f"Inserting batch of {len(transactions)} transactions for week {week}")
            bulk_insert_transactions(transactions)
        

    
def store_data(json_file):
    user_input = UserInput(json_file)
    user_input.insert_user()
    user_input.insert_transactions()

    


