import json
from datetime import datetime, timedelta
import random
from Transactions import *

example_account = {
    "account": {
      "account_id": "1234567890",
      "account_name": "Chase Total Checking",
      "account_type": "Checking",
      "account_number": "****0000",
      "balance": {
        "current": 2435.67,
        "available": 2200.45,
        "currency": "USD"
      },
      "owner_name": "Jane Doe",
      "bank_details": {
        "bank_name": "Chase Bank",
        "routing_number": "021000021",
        "branch": "Manhattan Main Branch, New York, NY"
      },
      "transactions": [
        {
          "transaction_id": "TXN10001",
          "date": "2025-01-01",
          "amount": -12.99,
          "category": "Dining",
          "merchant": "Chipotle Mexican Grill",
          "location": "Manhattan, New York, NY",
          "timestamp": "2025-01-01T12:34:56",
          "type": "Debit",
          "status": "Completed",
          "payment_method": "Debit Card"
        },
        {
          "transaction_id": "TXN10002",
          "date": "2025-01-02",
          "amount": -150.00,
          "category": "Groceries",
          "merchant": "Trader Joe's",
          "location": "Manhattan, New York, NY",
          "timestamp": "2025-01-02T17:45:00",
          "type": "Debit",
          "status": "Completed",
          "payment_method": "Debit Card"
        },
        {
          "transaction_id": "TXN10003",
          "date": "2025-01-03",
          "amount": -2200.00,
          "category": "Rent",
          "merchant": "Park Avenue Apartments",
          "location": "Manhattan, New York, NY",
          "timestamp": "2025-01-03T09:00:00",
          "type": "ACH Transfer",
          "status": "Completed",
          "payment_method": "ACH"
        },
        {
          "transaction_id": "TXN10004",
          "date": "2025-01-04",
          "amount": -12.99,
          "category": "Entertainment",
          "merchant": "Netflix",
          "location": "Online",
          "timestamp": "2025-01-04T21:00:00",
          "type": "Recurring Debit",
          "status": "Completed",
          "payment_method": "Subscription"
        },
        {
          "transaction_id": "TXN10005",
          "date": "2025-01-05",
          "amount": 200.00,
          "category": "Income",
          "merchant": "Freelance Work",
          "location": "Manhattan, New York, NY",
          "timestamp": "2025-01-05T14:30:00",
          "type": "Credit",
          "status": "Pending",
          "payment_method": "Bank Transfer"
        },
        {
          "transaction_id": "TXN10006",
          "date": "2025-01-06",
          "amount": 1000.00,
          "category": "Income",
          "merchant": "Freelance Work",
          "location": "Manhattan, New York, NY",
          "timestamp": "2025-01-06T14:30:00",
          "type": "Credit",
          "status": "Pending",
          "payment_method": "Bank Transfer"
        },
        {
            "transaction_id": "TXN10007",
            "date": "2025-01-07",
            "amount": -12.99,
            "category": "Dining",
            "merchant": "Chipotle Mexican Grill",
            "location": "Manhattan, New York, NY",
            "timestamp": "2025-01-07T14:30:00",
            "type": "Debit",
            "status": "Completed",
            "payment_method": "Debit Card"
        },
        {
            "transaction_id": "TXN10008",
            "date": "2025-01-08",
            "amount": -12.99,
            "category": "Dining",
            "merchant": "Chipotle Mexican Grill",
            "location": "Manhattan, New York, NY",
            "timestamp": "2025-01-08T14:30:00",
            "type": "Debit",
            "status": "Completed",
            "payment_method": "Debit Card"
        }
      ]
    }
  }



















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



# create a user named john his balance currently is 2,000 and his account id is 1234567890
john = User(user_id=1234567890, account_id=1234567890, name="John Doe", email="john.doe@example.com", phone_number="1234567890")
john.set_initial_balance(2000.00)

def transaction_generator(user):
    minutesinaday = 60 * 24
    timestamp = datetime(2025, 1, 1, 0, 0, 0)
    end_timestamp = datetime(2026, 1, 1, 0, 0, 0)
    transaction_counter = 1  # Add counter for unique transaction IDs
    
    while timestamp < end_timestamp:
        # Monthly transactions (1st of each month)
        if timestamp.day == 1:
            # Rent payment
            rent = RentalTransaction(
                transaction_id=f"TXN{transaction_counter:05d}",
                account_id=user.account["account_id"],
                date=timestamp.strftime('%Y-%m-%d'),
                amount=round(random.randint(2200, 2300), 2)
            )
            user.add_transaction(rent)
            transaction_counter += 1

            # Bill payments
            bill_payment = BillPaymentTransaction(
                transaction_id=f"TXN{transaction_counter:05d}",
                account_id=user.account["account_id"],
                date=timestamp.strftime('%Y-%m-%d'),
                amount=round(random.randint(200, 300), 2)
            )
            user.add_transaction(bill_payment)
            transaction_counter += 1

            # Netflix subscription
            netflix = NetflixTransaction(
                transaction_id=f"TXN{transaction_counter:05d}",
                account_id=user.account["account_id"],
                date=timestamp.strftime('%Y-%m-%d')
            )
            user.add_transaction(netflix)
            transaction_counter += 1

        # Bi-monthly salary (1st and 15th)
        if timestamp.day in [1, 15]:
            salary = SalaryTransaction(
                transaction_id=f"TXN{transaction_counter:05d}",
                account_id=user.account["account_id"],
                date=timestamp.strftime('%Y-%m-%d'),
                amount=2500
            )
            user.add_transaction(salary)
            transaction_counter += 1

        # Weekly groceries (every 7 days)
        if timestamp.day % 7 == 0:
            groceries = GroceryTransaction(
                transaction_id=f"TXN{transaction_counter:05d}",
                account_id=user.account["account_id"],
                date=timestamp.strftime('%Y-%m-%d'),
                amount=round(random.randint(50, 150), 2)
            )
            user.add_transaction(groceries)
            transaction_counter += 1

        # Daily transactions
        # Random chance for each type of daily transaction
        if random.random() < 0.7:  # 70% chance of buying coffee
            coffee = CoffeeTransaction(
                transaction_id=f"TXN{transaction_counter:05d}",
                account_id=user.account["account_id"],
                date=timestamp.strftime('%Y-%m-%d'),
                amount=round(random.uniform(5.75, 7.75), 2)
            )
            user.add_transaction(coffee)
            transaction_counter += 1

        if random.random() < 0.4:  # 40% chance of Chipotle
            chipotle = ChipotleTransaction(
                amount=round(random.uniform(12.99, 15.99), 2),
                transaction_id=f"TXN{transaction_counter:05d}",
                account_id=user.account["account_id"],
                date=timestamp.strftime('%Y-%m-%d')
            )
            user.add_transaction(chipotle)
            transaction_counter += 1

        if random.random() < 0.3:  # 30% chance of Uber
            uber = UberTransaction(
                transaction_id=f"TXN{transaction_counter:05d}",
                account_id=user.account["account_id"],
                date=timestamp.strftime('%Y-%m-%d'),
                amount=round(random.randint(10, 40), 2),
                service_type=random.choice(["Ride", "Food"])
            )
            user.add_transaction(uber)
            transaction_counter += 1

        if random.random() < 0.2:  # 20% chance of Amazon purchase
            amazon = AmazonTransaction(
                transaction_id=f"TXN{transaction_counter:05d}",
                account_id=user.account["account_id"],
                date=timestamp.strftime('%Y-%m-%d'),
                amount=round(random.randint(10, 100), 2),
                item_category=random.choice(["Electronics", "Books", "Clothing", "Home", "Other"])
            )
            user.add_transaction(amazon)
            transaction_counter += 1

        timestamp += timedelta(days=1)  # Increment by days instead of minutes

    return transaction_counter





transaction_generator(john)

john.save_to_json("john.json")

if __name__ == "__main__":
    try:
        # Create user
        john = User(
            user_id="USER123", 
            account_id="1234567890", 
            name="John Doe", 
            email="john.doe@example.com", 
            phone_number="1234567890"
        )
        
        # Set initial balance
        john.set_initial_balance(2000.00)
        
        # Generate transactions
        num_transactions = transaction_generator(john)
        print(f"Generated {num_transactions} transactions")
        
        # Save to file
        john.save_to_json("john.json")
        
        # Print final balance
        print(f"Final balance: ${john.get_balance()['current']:.2f}")
        
    except Exception as e:
        print(f"Error in main execution: {str(e)}")
