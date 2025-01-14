from Transactions.Transactions import Transaction


class SalaryTransaction(Transaction):
    def __init__(self, transaction_id, account_id, date, amount, merchant="Tech Corp Inc"):
        super().__init__(
            transaction_id=transaction_id,
            account_id=account_id,
            date=date,
            amount=abs(amount),  # Ensure amount is positive
            category="Income",
            merchant=merchant,
            location="New York, NY",
            type="Direct Deposit"
        )
        self.payment_method = "ACH"
        self.status = "Completed"
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