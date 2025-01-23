import os
from dotenv import load_dotenv
import psycopg2 as pg

load_dotenv()

database_credentials = {
    "DB_URL": os.getenv("DB_URL"),
    "DB_NAME": os.getenv("DB_NAME"),
    "DB_USER": os.getenv("DB_USER"),
    "DB_PASSWORD": os.getenv("DB_PASSWORD"),
    "DB_HOST": os.getenv("DB_HOST"),
    "DB_PORT": os.getenv("DB_PORT")
}

def connect_to_db():
    conn = pg.connect(database_credentials["DB_URL"])
    return conn

def close_db_connection(conn):
    conn.close()

def get_last_transaction_id():
    ## go to the last row in the transactions table and get the transaction_id the transaction_id is a varchar(20)
    conn = connect_to_db()
    cursor = conn.cursor()
    cursor.execute("SELECT transaction_id FROM transactions ORDER BY transaction_id DESC LIMIT 1")
    last_id = cursor.fetchone()[0]
    close_db_connection(conn)
    return last_id


def split_transaction_key(transaction_id):
    ## split the transaction_id into a number and return it
    return int(transaction_id[3:])

def get_next_transaction_id():
    last_id = get_last_transaction_id()
    next_number = split_transaction_key(last_id) + 1
    return f"TXN{next_number:05d}"



