import os,sys,json,random
from datetime import datetime, timedelta
from Transactions import *
from User import User
from transaction_generator import (
    transaction_generator,
    poor_spending_habits_generator,
    high_payment_habits_generator,
    student_spending_habits_generator
)
from Database import store_data, UserInput
from data_utils.util import database_credentials

def create_user_json(user, filename):
    """Create initial JSON file for a user"""
    user.save_to_json(filename)
    # Store user in database
    store_data(filename)

# Create users
bob = User(
    user_id=1234567894,
    account_id=1234567894,
    name="Bob Doe",
    email="bob.doe@example.com",
    phone_number="1234567894",
)

becky = User(
    user_id=1234567895,
    account_id=1234567895,
    name="Becky Doe",
    email="becky.doe@example.com",
    phone_number="1234567895",
)

# Create Abdul's profile
abdul = User(
    user_id=1234567900,
    account_id=1234567900,
    name="Abdul Mohammed",
    email="abdul.mohammed@example.com",
    phone_number="1234567900"
)

# Set initial balances
print("Setting initial balances...")
bob.set_initial_balance(10000.00)
becky.set_initial_balance(15000.00)
abdul.set_initial_balance(500.00)

# Create initial JSON files and store users in database
print("Creating initial user data...")
os.system("mkdir -p json_data")
create_user_json(bob, "json_data/bob.json")
create_user_json(becky, "json_data/becky.json")
create_user_json(abdul, "json_data/abdul.json")

# Generate transactions
print("Generating transactions...")
poor_spending_habits_generator(bob)
high_payment_habits_generator(becky)
student_spending_habits_generator(abdul)

print("Saving final transaction data...")
bob.save_to_json("json_data/bob.json")
becky.save_to_json("json_data/becky.json")
abdul.save_to_json("json_data/abdul.json")

print("Copying data to client and server...")
os.system("cp json_data/*.json client/src/data/")
os.system("cp json_data/*.json server/src/data/")

print("Updating database with final transaction data...")
store_data("json_data/bob.json")
store_data("json_data/becky.json")
store_data("json_data/abdul.json")

print("Done!")


