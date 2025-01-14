from Transactions.Transactions import Transaction



class AmazonTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, item_category="Shopping"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),
            category=item_category,
            merchant="Amazon.com",
            location="Online",
            type="Debit"
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
