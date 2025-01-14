from Transactions.Transactions import Transaction

class RentalTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, merchant="Park Avenue Apartments"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),  # Ensure amount is negative
            category="Rent",
            merchant=merchant,
            location="Manhattan, New York, NY",
            type="ACH Transfer"
        )
        self.payment_method = "ACH"
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