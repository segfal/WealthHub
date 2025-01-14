from datetime import datetime, timedelta
import random
from Transactions import *


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


