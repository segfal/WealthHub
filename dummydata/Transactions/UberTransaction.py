from Transactions.Transactions import Transaction



class UberTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, service_type="Ride"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),
            category="Transportation" if service_type == "Ride" else "Food Delivery",
            merchant=f"Uber{' Eats' if service_type == 'Food' else ''}",
            location="New York, NY",
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