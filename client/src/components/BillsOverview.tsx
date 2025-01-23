import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { AlertTriangle, Receipt, Calendar, Loader2 } from "lucide-react";
import { getBillsOverview } from "../lib/api";

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
  status: 'paid' | 'upcoming' | 'overdue';
}

interface BillsResponse {
  bills: Array<{
    category: string;
    amount: string;
    dueDate: string;
    status: string;
  }>;
  monthlyTotal: number;
  totalBills: number;
}

const BillsOverview = () => {
  const [bills, setBills] = useState<Bill[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [totalMonthlyBills, setTotalMonthlyBills] = useState(0);

  useEffect(() => {
    const accountId = import.meta.env.VITE_ACCOUNT_ID;
    if (!accountId) {
      setError("No account ID provided");
      setLoading(false);
      return;
    }

    setLoading(true);
    getBillsOverview(accountId)
      .then((data: BillsResponse) => {
        if (!data || !data.bills) {
          throw new Error('Invalid data format received from server');
        }

        // Process bills data
        const processedBills = data.bills.map(bill => ({
          category: bill.category,
          amount: parseFloat(bill.amount),
          dueDate: bill.dueDate,
          status: bill.status as 'paid' | 'upcoming' | 'overdue'
        }));

        setBills(processedBills);
        setTotalMonthlyBills(data.monthlyTotal);
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