from Transactions.Transactions import Transaction


class ChipotleTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, location="Manhattan, New York, NY"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-12.99,  # Standard Chipotle meal price
            category="Dining",
            merchant="Chipotle Mexican Grill",
            location=location,
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