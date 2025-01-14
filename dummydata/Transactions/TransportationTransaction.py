from Transactions.Transactions import Transaction

class TransportationTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount=127.00, merchant="MTA"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),
            category="Transportation",
            merchant=merchant,
            location="New York, NY",
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
