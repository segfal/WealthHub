from datetime import datetime

class Transaction:
    def __init__(self, transaction_id, account_id, date, amount, category, merchant, location, type):
        self.transaction_id = transaction_id
        self.account_id = account_id
        self.date = date
        self.amount = round(amount, 2)
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
