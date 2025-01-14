import json
from datetime import datetime, timedelta
import random
from Transactions import *
from User import User
from transaction_generator import transaction_generator


















# create a user named john his balance currently is 2,000 and his account id is 1234567890
john = User(user_id=1234567890, account_id=1234567890, name="John Doe", email="john.doe@example.com", phone_number="1234567890")
john.set_initial_balance(2000.00)




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
