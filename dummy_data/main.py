import os,sys,json,random
from datetime import datetime, timedelta
from Transactions import *
from User import User
from transaction_generator import transaction_generator
from Database import store_data
# create a user named john his balance currently is 2,000 and his account id is 1234567890
john = User(
    user_id=1234567890,
    account_id=1234567890,
    name="John Doe",
    email="john.doe@example.com",
    phone_number="1234567890",
)
jane = User(
    user_id=1234567891,
    account_id=1234567891,
    name="Jane Doe",
    email="jane.doe@example.com",
    phone_number="1234567891",
)

jake = User(
    user_id=1234567892,
    account_id=1234567892,
    name="Jake Doe",
    email="jake.doe@example.com",
    phone_number="1234567892",
)

jill = User(
    user_id=1234567893,
    account_id=1234567893,
    name="Jill Doe",
    email="jill.doe@example.com",
    phone_number="1234567893",
)

john.set_initial_balance(2000.00)
jane.set_initial_balance(1000.00)
jake.set_initial_balance(1500.00)
jill.set_initial_balance(1000.00)


transaction_generator(john)
transaction_generator(jane)
transaction_generator(jake)
transaction_generator(jill)

john.save_to_json("john.json")
jane.save_to_json("jane.json")
jake.save_to_json("jake.json")
jill.save_to_json("jill.json")


# send it to server and client using cp
os.system("cp *.json client/src/data")
os.system("cp *.json server/src/data")

# store the current data into a json directory
os.system("mkdir json_data")
os.system("mv *.json json_data")

store_data("json_data/jill.json")
store_data("json_data/john.json")
# store_data("json_data/jane.json")
store_data("json_data/jake.json")
