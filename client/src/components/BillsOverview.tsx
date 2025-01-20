import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { AlertTriangle, Receipt, Calendar, Loader2 } from "lucide-react";

// Categories that are considered bills/rent
const BILL_CATEGORIES = new Set([
  'Rent',
  'Utilities',
  'Insurance',
  'Phone Bill',
  'Internet',
  'Mortgage',
  'Water Bill',
  'Electric Bill',
  'Gas Bill'
]);

interface Bill {
  category: string;
  amount: number;
  dueDate: string;
  merchant: string;
  status: 'paid' | 'upcoming' | 'overdue';
}

interface Transaction {
  category: string;
  amount: number;
  date: string;
  merchant: string;
}

interface SpendingResponse {
  transactions: Transaction[];
}

const BillsOverview = () => {
  const [bills, setBills] = useState<Bill[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [totalMonthlyBills, setTotalMonthlyBills] = useState(0);

  useEffect(() => {
    setLoading(true);
    fetch('http://localhost:8080/api/analytics/spending?accountId=1234567891&timeRange=1%20month')
      .then(res => {
        if (!res.ok) {
          throw new Error(`Failed to fetch bills: ${res.status} ${res.statusText}`);
        }
        return res.json();
      })
      .then((data: SpendingResponse) => {
        if (!data || !data.transactions || !Array.isArray(data.transactions)) {
          throw new Error('Invalid data format received from server');
        }

        // Filter transactions to only include bills and recurring payments
        const billTransactions = data.transactions.filter((t: Transaction) => 
          t && t.category && BILL_CATEGORIES.has(t.category)
        );

        // Process transactions into bills
        const processedBills = billTransactions.map((t: Transaction) => ({
          category: t.category,
          amount: Math.abs(t.amount),
          dueDate: t.date,
          merchant: t.merchant,
          status: new Date(t.date) > new Date() ? 'upcoming' : 'paid'
        }));

        setBills(processedBills);
        setTotalMonthlyBills(processedBills.reduce((sum, bill) => sum + bill.amount, 0));
        setError(null);
      })
      .catch(err => {
        console.error('Error fetching bills:', err);
        setError(err.message);
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error) {
    return (
      <Card className="p-6">
        <div className="text-center text-destructive">
          <AlertTriangle className="w-8 h-8 mx-auto mb-2" />
          <p>{error}</p>
        </div>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* Monthly Summary */}
      <motion.div
        initial={{ scale: 0.95, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3 }}
      >
        <Card className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-medium flex items-center">
              <Receipt className="w-5 h-5 mr-2" />
              Monthly Bills & Recurring Payments
            </h3>
            <span className="text-2xl font-bold">${totalMonthlyBills.toFixed(2)}</span>
          </div>
          <p className="text-sm text-gray-500">
            You have {bills.length} recurring payments this month
          </p>
        </Card>
      </motion.div>

      {/* Bills List */}
      <div className="space-y-4">
        {bills.map((bill, index) => (
          <motion.div
            key={`${bill.category}-${bill.dueDate}`}
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ duration: 0.3, delay: index * 0.1 }}
          >
            <Card className="p-4">
              <div className="flex justify-between items-center">
                <div>
                  <h4 className="font-medium">{bill.category}</h4>
                  <div className="flex items-center text-sm text-gray-500 mt-1">
                    <Calendar className="w-4 h-4 mr-1" />
                    <span>Due {new Date(bill.dueDate).toLocaleDateString()}</span>
                  </div>
                  <div className="text-sm text-gray-500">{bill.merchant}</div>
                </div>
                <div className="text-right">
                  <div className="text-lg font-medium">${bill.amount.toFixed(2)}</div>
                  <div className={`text-sm mt-1 ${
                    bill.status === 'paid' ? 'text-green-500' :
                    bill.status === 'upcoming' ? 'text-yellow-500' :
                    'text-red-500'
                  }`}>
                    {bill.status.charAt(0).toUpperCase() + bill.status.slice(1)}
                  </div>
                </div>
              </div>
            </Card>
          </motion.div>
        ))}
      </div>
    </div>
  );
};

export default BillsOverview; 