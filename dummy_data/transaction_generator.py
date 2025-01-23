from datetime import datetime, timedelta
import random
import json
from Transactions import *
from data_utils.util import get_last_transaction_id, get_next_transaction_id
from Database import bulk_insert_transactions

class TransactionGenerator:
    def __init__(self, user):
        self.user = user
        self.start_date = datetime(2025, 1, 1)
        self.end_date = datetime(2026, 1, 1)
        self.transaction_counter = 1
        # Get initial transaction ID and pre-generate a batch
        self.current_transaction_id = int(get_next_transaction_id()[3:])
        # Store transactions in memory before bulk insert
        self.pending_transactions = []
        
    def _get_next_transaction_id(self):
        """Generate next transaction ID without database call"""
        self.current_transaction_id += 1
        return f"TXN{self.current_transaction_id:05d}"
    
    def _add_transaction(self, transaction_type, date, amount=None, **kwargs):
        """Helper method to create and add transactions"""
        transaction_args = {
            "transaction_id": self._get_next_transaction_id(),
            "account_id": self.user.account["account_id"],
            "date": date,
            **({"amount": amount} if amount is not None else {}),
            **kwargs
        }
        transaction = transaction_type(**transaction_args)
        self.user.add_transaction(transaction)
        self.pending_transactions.append(transaction.create_transaction())

    def _process_pending_transactions(self):
        """Process pending transactions in weekly batches"""
        if not self.pending_transactions:
            return

        # Group transactions by week
        transactions_by_week = {}
        for transaction in self.pending_transactions:
            date = datetime.strptime(transaction['date'], '%Y-%m-%d')
            week_key = date.strftime('%Y-%W')
            if week_key not in transactions_by_week:
                transactions_by_week[week_key] = []
            transactions_by_week[week_key].append(transaction)

        # Process each week's transactions
        for week, transactions in transactions_by_week.items():
            print(f"Processing batch of {len(transactions)} transactions for week {week}")
            bulk_insert_transactions(transactions)

        # Clear pending transactions
        self.pending_transactions = []

    def generate_regular_transactions(self):
        """Generate regular spending pattern transactions"""
        print(f"Generating regular transactions for {self.user.account['owner_name']}")
        current_date = self.start_date
        last_week = None
        
        while current_date < self.end_date:
            date_str = current_date.strftime('%Y-%m-%d')
            current_week = current_date.strftime('%Y-%W')
            
            # Process pending transactions when we move to a new week
            if last_week and current_week != last_week:
                self._process_pending_transactions()
            
            # Monthly transactions
            if current_date.day == 1:
                self._add_transaction(RentalTransaction, date_str, random.randint(2200, 2300))
                self._add_transaction(BillPaymentTransaction, date_str, random.randint(200, 300))
                self._add_transaction(NetflixTransaction, date_str)
            
            # Bi-monthly salary
            if current_date.day in [1, 15]:
                self._add_transaction(SalaryTransaction, date_str, 2500)
            
            # Weekly groceries
            if current_date.day % 7 == 0:
                self._add_transaction(GroceryTransaction, date_str, random.randint(50, 150))
            
            # Daily transactions
            if random.random() < 0.7:
                self._add_transaction(CoffeeTransaction, date_str, random.uniform(5.75, 7.75))
            
            if random.random() < 0.4:
                self._add_transaction(ChipotleTransaction, date_str, random.uniform(12.99, 15.99))
            
            if random.random() < 0.3:
                self._add_transaction(UberTransaction, date_str, random.randint(10, 40),
                                   service_type=random.choice(["Ride", "Food"]))
            
            if random.random() < 0.2:
                self._add_transaction(AmazonTransaction, date_str, random.randint(10, 100))
            
            last_week = current_week
            current_date += timedelta(days=1)
        
        # Process any remaining transactions
        self._process_pending_transactions()
        print(f"Finished generating transactions for {self.user.account['owner_name']}")

    def generate_poor_spending_habits(self):
        """Generate transactions for poor spending habits"""
        print(f"Generating poor spending transactions for {self.user.account['owner_name']}")
        current_date = self.start_date
        last_week = None
        
        while current_date < self.end_date:
            date_str = current_date.strftime('%Y-%m-%d')
            current_week = current_date.strftime('%Y-%W')
            
            # Process pending transactions when we move to a new week
            if last_week and current_week != last_week:
                self._process_pending_transactions()
            
            # Bi-weekly salary
            if current_date.day in [1, 15]:
                self._add_transaction(SalaryTransaction, date_str, 5000.00)
            
            # Multiple daily transactions
            for _ in range(3):
                if random.random() < 0.9:
                    self._add_transaction(CoffeeTransaction, date_str, random.uniform(5.00, 8.00))
                if random.random() < 0.8:
                    self._add_transaction(ChipotleTransaction, date_str, random.uniform(15.00, 25.00))
                if random.random() < 0.7:
                    self._add_transaction(UberTransaction, date_str, random.uniform(20.00, 50.00),
                                       service_type=random.choice(["Ride", "Food"]))
            
            # Daily Amazon shopping
            if random.random() < 0.6:
                self._add_transaction(AmazonTransaction, date_str, random.uniform(30.00, 200.00))
            
            last_week = current_week
            current_date += timedelta(days=1)
        
        # Process any remaining transactions
        self._process_pending_transactions()
        print(f"Finished generating transactions for {self.user.account['owner_name']}")

    def generate_high_payment_habits(self):
        """Generate transactions for high payment habits"""
        print(f"Generating high payment transactions for {self.user.account['owner_name']}")
        current_date = self.start_date
        last_week = None
        
        while current_date < self.end_date:
            date_str = current_date.strftime('%Y-%m-%d')
            current_week = current_date.strftime('%Y-%W')
            
            # Process pending transactions when we move to a new week
            if last_week and current_week != last_week:
                self._process_pending_transactions()
            
            # Bi-weekly salary
            if current_date.day in [1, 15]:
                self._add_transaction(SalaryTransaction, date_str, 7000.00)
            
            # Monthly bills
            if current_date.day == 1:
                self._add_transaction(NetflixTransaction, date_str)
                for _ in range(3):
                    self._add_transaction(UtilityTransaction, date_str, random.uniform(100.00, 300.00))
                self._add_transaction(RentalTransaction, date_str, random.uniform(2500.00, 3500.00))
            
            # Daily bill payments
            if random.random() < 0.4:
                self._add_transaction(BillPaymentTransaction, date_str, random.uniform(50.00, 200.00))
            
            last_week = current_week
            current_date += timedelta(days=1)
        
        # Process any remaining transactions
        self._process_pending_transactions()
        print(f"Finished generating transactions for {self.user.account['owner_name']}")

    def generate_frugal_spending_habits(self):
        """Generate transactions for a frugal spender who saves money"""
        print(f"Generating frugal spending transactions for {self.user.account['owner_name']}")
        current_date = self.start_date
        last_week = None
        
        while current_date < self.end_date:
            date_str = current_date.strftime('%Y-%m-%d')
            current_week = current_date.strftime('%Y-%W')
            
            if last_week and current_week != last_week:
                self._process_pending_transactions()
            
            # Bi-weekly salary (same as others but saves more)
            if current_date.day in [1, 15]:
                self._add_transaction(SalaryTransaction, date_str, 3000.00)
            
            # Monthly expenses (minimal)
            if current_date.day == 1:
                self._add_transaction(RentalTransaction, date_str, 1200.00)  # Cheaper rent
                self._add_transaction(UtilityTransaction, date_str, random.uniform(50.00, 100.00))
            
            # Occasional coffee (less frequent)
            if random.random() < 0.3:  # Only 30% chance
                self._add_transaction(CoffeeTransaction, date_str, random.uniform(2.50, 4.00))
            
            # Grocery shopping (bulk buying, less frequent but larger amounts)
            if current_date.day % 14 == 0:  # Every two weeks
                self._add_transaction(GroceryTransaction, date_str, random.uniform(150.00, 200.00))
            
            last_week = current_week
            current_date += timedelta(days=1)
        
        self._process_pending_transactions()
        print(f"Finished generating transactions for {self.user.account['owner_name']}")

    def generate_tech_savvy_spending_habits(self):
        """Generate transactions for a tech-savvy person with online shopping habits"""
        print(f"Generating tech-savvy spending transactions for {self.user.account['owner_name']}")
        current_date = self.start_date
        last_week = None
        
        while current_date < self.end_date:
            date_str = current_date.strftime('%Y-%m-%d')
            current_week = current_date.strftime('%Y-%W')
            
            if last_week and current_week != last_week:
                self._process_pending_transactions()
            
            # Bi-weekly salary
            if current_date.day in [1, 15]:
                self._add_transaction(SalaryTransaction, date_str, 4500.00)
            
            # Monthly subscriptions
            if current_date.day == 1:
                self._add_transaction(NetflixTransaction, date_str)
                self._add_transaction(SubscriptionTransaction, date_str, 14.99, merchant="Spotify")
                self._add_transaction(SubscriptionTransaction, date_str, 9.99, merchant="iCloud")
                self._add_transaction(SubscriptionTransaction, date_str, 12.99, merchant="Amazon Prime")
                self._add_transaction(RentalTransaction, date_str, 2000.00)
                self._add_transaction(UtilityTransaction, date_str, random.uniform(100.00, 200.00))
            
            # Frequent online shopping
            if random.random() < 0.4:
                self._add_transaction(AmazonTransaction, date_str, random.uniform(20.00, 300.00))
            
            # Food delivery
            if random.random() < 0.6:
                self._add_transaction(UberTransaction, date_str, random.uniform(15.00, 40.00), service_type="Food")
            
            last_week = current_week
            current_date += timedelta(days=1)
        
        self._process_pending_transactions()
        print(f"Finished generating transactions for {self.user.account['owner_name']}")

    def generate_foodie_spending_habits(self):
        """Generate transactions for a food enthusiast"""
        print(f"Generating foodie spending transactions for {self.user.account['owner_name']}")
        current_date = self.start_date
        last_week = None
        
        while current_date < self.end_date:
            date_str = current_date.strftime('%Y-%m-%d')
            current_week = current_date.strftime('%Y-%W')
            
            if last_week and current_week != last_week:
                self._process_pending_transactions()
            
            # Bi-weekly salary
            if current_date.day in [1, 15]:
                self._add_transaction(SalaryTransaction, date_str, 4000.00)
            
            # Monthly fixed expenses
            if current_date.day == 1:
                self._add_transaction(RentalTransaction, date_str, 1800.00)
                self._add_transaction(UtilityTransaction, date_str, random.uniform(100.00, 200.00))
            
            # High-end restaurants on weekends
            if current_date.weekday() >= 5:  # Saturday and Sunday
                self._add_transaction(BillPaymentTransaction, date_str, random.uniform(100.00, 300.00), 
                                    merchant="Fine Dining")
            
            # Regular food expenses
            if random.random() < 0.8:  # 80% chance of eating out
                self._add_transaction(ChipotleTransaction, date_str, random.uniform(15.00, 30.00))
            
            # Coffee enthusiast
            if random.random() < 0.9:  # 90% chance of specialty coffee
                self._add_transaction(CoffeeTransaction, date_str, random.uniform(6.00, 12.00))
            
            # Grocery shopping for cooking
            if current_date.day % 4 == 0:  # Every 4 days
                self._add_transaction(GroceryTransaction, date_str, random.uniform(80.00, 150.00))
            
            last_week = current_week
            current_date += timedelta(days=1)
        
        self._process_pending_transactions()
        print(f"Finished generating transactions for {self.user.account['owner_name']}")

    def generate_student_spending_habits(self):
        """Generate transactions for a student with specific spending habits (Abdul's profile)"""
        print(f"Generating student spending transactions for {self.user.account['owner_name']}")
        current_date = self.start_date
        last_week = None
        
        while current_date < self.end_date:
            date_str = current_date.strftime('%Y-%m-%d')
            current_week = current_date.strftime('%Y-%W')
            
            if last_week and current_week != last_week:
                self._process_pending_transactions()
            
            # Bi-weekly salary ($385)
            if current_date.day in [1, 15]:
                self._add_transaction(SalaryTransaction, date_str, 385.00)
            
            # Monthly fixed expenses
            if current_date.day == 1:
                # Car insurance
                self._add_transaction(BillPaymentTransaction, date_str, 297.00, merchant="Car Insurance")
                # Phone bill
                self._add_transaction(BillPaymentTransaction, date_str, 30.00, merchant="Phone Bill")
                # ChatGPT subscription
                self._add_transaction(SubscriptionTransaction, date_str, 20.00, merchant="ChatGPT")
            
            # Weekday purchases (Monday-Friday)
            if current_date.weekday() < 5:
                # Daily McDonald's around 1 PM
                self._add_transaction(McdonaldsTransaction, date_str, random.uniform(8.00, 15.00))
                # Daily Dunkin' coffee
                self._add_transaction(CoffeeTransaction, date_str, 4.00, merchant="Dunkin'")
            
            # Weekly groceries
            if current_date.day % 7 == 0:
                self._add_transaction(GroceryTransaction, date_str, 30.00)
            
            # Gas every 4 days
            if current_date.day % 4 == 0:
                self._add_transaction(BillPaymentTransaction, date_str, 20.00, merchant="Gas Station")
            
            last_week = current_week
            current_date += timedelta(days=1)
        
        self._process_pending_transactions()
        print(f"Finished generating transactions for {self.user.account['owner_name']}")


# Helper functions to maintain backward compatibility
def transaction_generator(user):
    generator = TransactionGenerator(user)
    generator.generate_regular_transactions()

def poor_spending_habits_generator(user):
    generator = TransactionGenerator(user)
    generator.generate_poor_spending_habits()

def high_payment_habits_generator(user):
    generator = TransactionGenerator(user)
    generator.generate_high_payment_habits()

def frugal_spending_habits_generator(user):
    generator = TransactionGenerator(user)
    generator.generate_frugal_spending_habits()

def tech_savvy_spending_habits_generator(user):
    generator = TransactionGenerator(user)
    generator.generate_tech_savvy_spending_habits()

def foodie_spending_habits_generator(user):
    generator = TransactionGenerator(user)
    generator.generate_foodie_spending_habits()

def student_spending_habits_generator(user):
    generator = TransactionGenerator(user)
    generator.generate_student_spending_habits()


