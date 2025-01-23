from Transactions.Transactions import Transaction

class McdonaldsTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount=12.99, merchant="McDonald's"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),  # Ensure amount is negative
            category="Fast Food",
            merchant=merchant,
            location="Queens, New York, NY",
            type="Debit"
        )
        self.payment_method = "Debit Card"
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