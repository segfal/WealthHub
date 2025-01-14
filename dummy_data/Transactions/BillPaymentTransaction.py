from Transactions.Transactions import Transaction

class BillPaymentTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, merchant="Bill Payment"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=-abs(amount),
            category="Bill Payment",
            merchant=merchant,
            location="New York, NY",
            type="Bill Payment"
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

