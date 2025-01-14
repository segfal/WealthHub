from Transactions.Transactions import Transaction



class CoffeeTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount=5.75, merchant="Starbucks"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),
            category="Dining",
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