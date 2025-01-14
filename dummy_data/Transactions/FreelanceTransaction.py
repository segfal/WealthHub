from Transactions.Transactions import Transaction

class FreelanceTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, merchant="Freelance Work"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=abs(amount),  # Ensure amount is positive
            category="Income",
            merchant="Freelance Work",
            location="Manhattan, New York, NY",
            type="Credit"
        )
        self.status = "Pending"
        self.payment_method = "Bank Transfer"
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