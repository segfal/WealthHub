import os
import psycopg2 as pg
from dotenv import load_dotenv
import json
load_dotenv()


'''
date is seen as timestamp and is formatted like this
2024-12-31 19:00:00

'''
conn = pg.connect(os.getenv("DB_URL"))

cursor = conn.cursor()

account_id = "1234567900"
category_1 = "Bill Payment"
category_2 = "Subscription"
# from 2025-01-01 to 2025-01-31
cursor.execute("SELECT * FROM transactions WHERE account_id = %s AND category = %s OR category = %s AND date >= '2025-01-01 00:00:00' AND date <= '2025-01-31 23:59:59'", (account_id, category_1, category_2))


rows = cursor.fetchall()

print(rows)
# date time is not json serializable    
# save to json
with open("bills.json", "w") as f:
    json.dump(rows, f, default=str, indent=4)
    
category_3 = "Income"  
cursor.execute("SELECT * FROM transactions WHERE account_id = %s AND category = %s AND date >= '2025-01-01 00:00:00' AND date <= '2025-01-31 23:59:59'", (account_id, category_3))
rows = cursor.fetchall() 

print(rows)

with open("income.json", "w") as f:
    json.dump(rows, f, default=str, indent=4)
