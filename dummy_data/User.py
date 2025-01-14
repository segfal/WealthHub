
from Transactions.Transactions import Transaction
import json
from datetime import datetime

from Transactions import *




class User:
    def __init__(self, user_id, account_id, name, email, phone_number):
        self.user_id = user_id
        self.account = {
            "account_id": account_id,
            "account_name": "Chase Total Checking",
            "account_type": "Checking",
            "account_number": f"****{str(account_id)[-4:]}",
            "balance": {
                "current": 0.0,
                "available": 0.0,
                "currency": "USD"
            },
            "owner_name": name,
            "bank_details": {
                "bank_name": "Chase Bank",
                "routing_number": "021000021",
                "branch": "Manhattan Main Branch, New York, NY"
            },
            "transactions": []
        }
        self.email = email
        self.phone_number = phone_number

    def validate_transaction(self, transaction):
        """Validate transaction before adding it"""
        if not transaction or not hasattr(transaction, 'create_transaction'):
            raise ValueError("Invalid transaction object")
        
        tx_dict = transaction.create_transaction()
        required_fields = ['transaction_id', 'amount', 'date']
        
        for field in required_fields:
            if field not in tx_dict:
                raise ValueError(f"Transaction missing required field: {field}")
        
        # Validate amount is numeric
        if not isinstance(tx_dict['amount'], (int, float)):
            raise ValueError("Transaction amount must be numeric")
        
        return tx_dict

    def add_transaction(self, transaction):
        """Add transaction with validation"""
        try:
            tx_dict = self.validate_transaction(transaction)
            
            # Round the amount to 2 decimal places
            tx_dict["amount"] = round(tx_dict["amount"], 2)
            
            self.account["transactions"].append(tx_dict)
            
            # Update balance based on transaction amount
            amount = tx_dict["amount"]
            self.account["balance"]["current"] = round(self.account["balance"]["current"] + amount, 2)
            self.account["balance"]["available"] = round(self.account["balance"]["available"] + amount, 2)
            
            # Ensure balance doesn't go below zero (optional)
            if self.account["balance"]["available"] < 0:
                print(f"Warning: Available balance is negative: ${self.account['balance']['available']:.2f}")
                
        except Exception as e:
            print(f"Error adding transaction: {str(e)}")
            raise

    def get_transactions(self):
        return self.account["transactions"]

    def get_transaction_by_id(self, transaction_id):
        return next((t for t in self.account["transactions"] if t["transaction_id"] == transaction_id), None)
    
    def get_balance(self):
        return self.account["balance"]
    
    def set_initial_balance(self, current_balance, available_balance=None):
        self.account["balance"]["current"] = round(current_balance, 2)
        self.account["balance"]["available"] = round(available_balance if available_balance is not None else current_balance, 2)

    def save_to_json(self, filename):
        with open(filename, 'w') as f:
            json.dump({"account": self.account}, f, indent=4)

    def load_from_json(self, filename):
        with open(filename, 'r') as f:
            data = json.load(f)
            self.account = data["account"]

    def get_transactions_by_date_range(self, start_date, end_date):
        """Get transactions within a date range"""
        try:
            start = datetime.strptime(start_date, '%Y-%m-%d')
            end = datetime.strptime(end_date, '%Y-%m-%d')
            
            return [
                tx for tx in self.account["transactions"]
                if start <= datetime.strptime(tx["date"], '%Y-%m-%d') <= end
            ]
        except ValueError as e:
            print(f"Error parsing dates: {str(e)}")
            raise
