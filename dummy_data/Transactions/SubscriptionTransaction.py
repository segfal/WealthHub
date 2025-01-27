from Transactions.Transactions import Transaction



class SubscriptionTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount=15.99, merchant="Netflix"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),
            category="Subscription",
            merchant=merchant,
            location="Online",
            type="Recurring Debit"
        )
        self.payment_method = "Credit Card"
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


