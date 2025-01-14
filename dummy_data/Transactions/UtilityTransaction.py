from Transactions.Transactions import Transaction

class UtilityTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, merchant="ConEdison", utility_type="Electricity"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),  # Ensure amount is negative
            category="Utilities",
            merchant=merchant,
            location="New York, NY",
            type="Bill Payment"
        )
        self.payment_method = "ACH"
        self.utility_type = utility_type
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