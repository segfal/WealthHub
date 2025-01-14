from Transactions.Transactions import Transaction

class GroceryTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, merchant="Trader Joe's"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),  # Ensure amount is negative
            category="Groceries",
            merchant=merchant,
            location="Manhattan, New York, NY",
            type="Debit"
        )
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