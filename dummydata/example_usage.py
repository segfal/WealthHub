import json

# Create example transactions
account_id = "1234567890"

# Create a Chipotle transaction
chipotle_tx = ChipotleTransaction(
    transaction_id="TXN10001",
    account_id=account_id,
    date="2025-01-01"
)

# Create a Netflix subscription
netflix_tx = NetflixTransaction(
    transaction_id="TXN10004",
    account_id=account_id,
    date="2025-01-04"
)

# Create a rental payment
rent_tx = RentalTransaction(
    transaction_id="TXN10003",
    account_id=account_id,
    date="2025-01-03",
    amount=2200.00
)

# Create a freelance payment
freelance_tx = FreelanceTransaction(
    transaction_id="TXN10005",
    account_id=account_id,
    date="2025-01-05",
    amount=200.00
)

# Create a grocery transaction
grocery_tx = GroceryTransaction(
    transaction_id="TXN10002",
    account_id=account_id,
    date="2025-01-02",
    amount=150.00
)

# Create a utility bill payment
utility_tx = UtilityTransaction(
    transaction_id="TXN10009",
    account_id=account_id,
    date="2025-01-09",
    amount=150.50,
    merchant="ConEdison",
    utility_type="Electricity"
)

# Create a monthly metro card payment
transport_tx = TransportationTransaction(
    transaction_id="TXN10010",
    account_id=account_id,
    date="2025-01-10"
)

# Create a salary deposit
salary_tx = SalaryTransaction(
    transaction_id="TXN10011",
    account_id=account_id,
    date="2025-01-15",
    amount=5000.00
)

# Create a coffee purchase
coffee_tx = CoffeeTransaction(
    transaction_id="TXN10012",
    account_id=account_id,
    date="2025-01-11"
)

# Create an Amazon purchase
amazon_tx = AmazonTransaction(
    transaction_id="TXN10013",
    account_id=account_id,
    date="2025-01-12",
    amount=49.99,
    item_category="Electronics"
)

# Create an Uber ride
uber_tx = UberTransaction(
    transaction_id="TXN10014",
    account_id=account_id,
    date="2025-01-13",
    amount=25.50,
    service_type="Ride"
)

# Create a subscription payment (e.g., Spotify)
spotify_tx = SubscriptionTransaction(
    transaction_id="TXN10015",
    account_id=account_id,
    date="2025-01-14",
    amount=9.99,
    merchant="Spotify",
    category="Entertainment"
)

# Create a new user
user = User(
    user_id="USER123",
    account_id="1234567890",
    name="Jane Doe",
    email="jane.doe@email.com",
    phone_number="212-555-0123"
)

# Set initial balance
user.set_initial_balance(2435.67, 2200.45)

# Add all the transactions
transactions_to_add = [
    chipotle_tx,
    netflix_tx,
    rent_tx,
    freelance_tx,
    grocery_tx,
    utility_tx,
    transport_tx,
    salary_tx,
    coffee_tx,
    amazon_tx,
    uber_tx,
    spotify_tx
]

for tx in transactions_to_add:
    user.add_transaction(tx)

# Save to JSON file
user.save_to_json("user_account.json")

# Print the entire account structure
print(json.dumps({"account": user.account}, indent=2))

# Print current balance
print(f"\nCurrent Balance: ${user.get_balance()['current']:.2f}")
print(f"Available Balance: ${user.get_balance()['available']:.2f}") 