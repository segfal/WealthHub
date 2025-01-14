from Transactions.Transactions import Transaction
    
class NetflixTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-12.99,
            category="Entertainment",
            merchant="Netflix",
            location="Online",
            type="Recurring Debit"
        )
        self.payment_method = "Subscription"
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