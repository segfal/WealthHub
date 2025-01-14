from Transactions.Transactions import Transaction



class SubscriptionTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, merchant, category="Subscription"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),
            category=category,
            merchant=merchant,
            location="Online",
            type="Recurring Debit"
        )
        self.payment_method = "Subscription"
        self.object = {
            "account_id": self.account_id,
            "account_name": self.account_name,
            "account_type": self.account_type,
            "account_number": self.account_number,
            "balance": self.balance,
            "owner_name": self.owner_name,
            "bank_details": self.bank_details,
            "transactions": self.transactions
        }


