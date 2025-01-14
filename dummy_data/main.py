import os,sys,json,random
from datetime import datetime, timedelta
from Transactions import *
from User import User
from transaction_generator import transaction_generator

# create a user named john his balance currently is 2,000 and his account id is 1234567890
john = User(
    user_id=1234567890,
    account_id=1234567890,
    name="John Doe",
    email="john.doe@example.com",
    phone_number="1234567890",
)
john.set_initial_balance(2000.00)


transaction_generator(john)

john.save_to_json("john.json")

print(john.account)
print(john.account["transactions"])