import json
from datetime import datetime, timedelta
import random

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


class Transaction:
    def __init__(self, transaction_id, account_id, date, amount, category, merchant, location, type):
        self.transaction_id = transaction_id
        self.account_id = account_id
        self.date = date
        self.amount = amount
        self.category = category
        self.merchant = merchant
        self.location = location
        self.type = type
        self.status = "Completed"  # Default status
        self.timestamp = datetime.now().strftime("%Y-%m-%dT%H:%M:%S")
        self.payment_method = "Debit Card"  # Default payment method
        self.object = {
            "transaction_id": self.transaction_id,
            "account_id": self.account_id,
            "date": self.date,
            "amount": self.amount,
            "category": self.category,
            "merchant": self.merchant,
            "location": self.location,
            "type": self.type
        }

    def create_transaction(self):
        return {
            "transaction_id": self.transaction_id,
            "account_id": self.account_id,
            "date": self.date,
            "amount": self.amount,
            "category": self.category,
            "merchant": self.merchant,
            "location": self.location,
            "type": self.type,
            "status": self.status,
            "timestamp": self.timestamp,
            "payment_method": self.payment_method
        }

class ChipotleTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, location="Manhattan, New York, NY"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-12.99,  # Standard Chipotle meal price
            category="Dining",
            merchant="Chipotle Mexican Grill",
            location=location,
            type="Debit"
        )
        self.object = {
            "transaction_id": self.transaction_id,
            "account_id": self.account_id,
            "date": self.date,
            "amount": self.amount,
            "category": self.category,
            "merchant": self.merchant,
            "location": self.location,
            "type": self.type
        }

class NetflixTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-12.99,
            category="Entertainment",
            merchant="Netflix",
            location="Online",
            type="Recurring Debit"
        )
        self.payment_method = "Subscription"
        self.object = {
            "transaction_id": self.transaction_id,
            "account_id": self.account_id,
            "date": self.date,
            "amount": self.amount,
            "category": self.category,
            "merchant": self.merchant,
            "location": self.location,
            "type": self.type
        }

class RentalTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, merchant="Park Avenue Apartments"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),  # Ensure amount is negative
            category="Rent",
            merchant=merchant,
            location="Manhattan, New York, NY",
            type="ACH Transfer"
        )
        self.payment_method = "ACH"
        self.object = {
            "transaction_id": self.transaction_id,
            "account_id": self.account_id,
            "date": self.date,
            "amount": self.amount,
            "category": self.category,
            "merchant": self.merchant,
            "location": self.location,
            "type": self.type
        }

class BillPaymentTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),
            category="Bill Payment",
            merchant="Bill Payment",
            location="New York, NY",
            type="Bill Payment"
        )

class FreelanceTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=abs(amount),  # Ensure amount is positive
            category="Income",
            merchant="Freelance Work",
            location="Manhattan, New York, NY",
            type="Credit"
        )
        self.status = "Pending"
        self.payment_method = "Bank Transfer"
        self.object = {
            "transaction_id": self.transaction_id,
            "account_id": self.account_id,
            "date": self.date,
            "amount": self.amount,
            "category": self.category,
            "merchant": self.merchant,
            "location": self.location,
            "type": self.type
        }

class GroceryTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, merchant="Trader Joe's"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),  # Ensure amount is negative
            category="Groceries",
            merchant=merchant,
            location="Manhattan, New York, NY",
            type="Debit"
        )
        self.object = {
            "transaction_id": self.transaction_id,
            "account_id": self.account_id,
            "date": self.date,
            "amount": self.amount,
            "category": self.category,
            "merchant": self.merchant,
            "location": self.location,
            "type": self.type
        }

class UtilityTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, merchant="ConEdison", utility_type="Electricity"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),  # Ensure amount is negative
            category="Utilities",
            merchant=merchant,
            location="New York, NY",
            type="Bill Payment"
        )
        self.payment_method = "ACH"
        self.utility_type = utility_type
        self.object = {
            "transaction_id": self.transaction_id,
            "account_id": self.account_id,
            "date": self.date,
            "amount": self.amount,
            "category": self.category,
            "merchant": self.merchant,
            "location": self.location,
            "type": self.type
        }
class TransportationTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount=127.00, merchant="MTA"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),
            category="Transportation",
            merchant=merchant,
            location="New York, NY",
            type="Recurring Debit"
        ) 
        self.payment_method = "Credit Card"
        self.object = {
            "transaction_id": self.transaction_id,
            "account_id": self.account_id,
            "date": self.date,
            "amount": self.amount,
            "category": self.category,
            "merchant": self.merchant,
            "location": self.location,
            "type": self.type
        }

class SalaryTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, merchant="Tech Corp Inc"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=abs(amount),  # Ensure amount is positive
            category="Income",
            merchant=merchant,
            location="New York, NY",
            type="Direct Deposit"
        )
        self.payment_method = "ACH"
        self.status = "Completed"
        self.object = {
            "transaction_id": self.transaction_id,
            "account_id": self.account_id,
            "date": self.date,
            "amount": self.amount,
            "category": self.category,
            "merchant": self.merchant,
            "location": self.location,
            "type": self.type
        }

class CoffeeTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount=5.75, merchant="Starbucks"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),
            category="Dining",
            merchant=merchant,
            location="Manhattan, New York, NY",
            type="Debit"
        )
        self.object = {
            "transaction_id": self.transaction_id,
            "account_id": self.account_id,
            "date": self.date,
            "amount": self.amount,
            "category": self.category,
            "merchant": self.merchant,
            "location": self.location,
            "type": self.type
        }

class AmazonTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, item_category="Shopping"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),
            category=item_category,
            merchant="Amazon.com",
            location="Online",
            type="Debit"
        )
        self.payment_method = "Credit Card"
        self.object = {
            "transaction_id": self.transaction_id,
            "account_id": self.account_id,
            "date": self.date,
            "amount": self.amount,
            "category": self.category,
            "merchant": self.merchant,
            "location": self.location,
            "type": self.type
        }

class UberTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, service_type="Ride"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),
            category="Transportation" if service_type == "Ride" else "Food Delivery",
            merchant=f"Uber{' Eats' if service_type == 'Food' else ''}",
            location="New York, NY",
            type="Debit"
        )
        self.object = {
            "transaction_id": self.transaction_id,
            "account_id": self.account_id,
            "date": self.date,
            "amount": self.amount,
            "category": self.category,
            "merchant": self.merchant,
            "location": self.location,
            "type": self.type
        }

class SubscriptionTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, merchant, category="Subscription"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),
            category=category,
            merchant=merchant,
            location="Online",
            type="Recurring Debit"
        )
        self.payment_method = "Subscription"
        self.object = {
            "account_id": self.account_id,
            "account_name": self.account_name,
            "account_type": self.account_type,
            "account_number": self.account_number,
            "balance": self.balance,
            "owner_name": self.owner_name,
            "bank_details": self.bank_details,
            "transactions": self.transactions
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
            self.account["transactions"].append(tx_dict)
            
            # Update balance based on transaction amount
            amount = tx_dict["amount"]
            self.account["balance"]["current"] += amount
            self.account["balance"]["available"] += amount
            
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
        self.account["balance"]["current"] = current_balance
        self.account["balance"]["available"] = available_balance if available_balance is not None else current_balance

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
                account_id=user.account["account_id"],  # Fix: use account from user object
                date=timestamp.strftime('%Y-%m-%d'),
                amount=random.randint(2200, 2300)
            )
            user.add_transaction(rent)
            transaction_counter += 1

            # Bill payments
            bill_payment = BillPaymentTransaction(
                transaction_id=f"TXN{transaction_counter:05d}",
                account_id=user.account["account_id"],
                date=timestamp.strftime('%Y-%m-%d'),
                amount=random.randint(200, 300)
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
                amount=random.randint(50, 150)  # More realistic grocery amounts
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
                amount=random.uniform(5.75, 7.75)
            )
            user.add_transaction(coffee)
            transaction_counter += 1

        if random.random() < 0.4:  # 40% chance of Chipotle
            chipotle = ChipotleTransaction(
                amount=random.uniform(12.99, 15.99),
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
                amount=random.randint(10, 40),
                service_type=random.choice(["Ride", "Food"])
            )
            user.add_transaction(uber)
            transaction_counter += 1

        if random.random() < 0.2:  # 20% chance of Amazon purchase
            amazon = AmazonTransaction(
                transaction_id=f"TXN{transaction_counter:05d}",
                account_id=user.account["account_id"],
                date=timestamp.strftime('%Y-%m-%d'),
                amount=random.randint(10, 100),
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
