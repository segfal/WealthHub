import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { Calendar, AlertTriangle, Loader2, Clock, DollarSign } from "lucide-react";

// Categories that are considered bills
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
  status: string;
}

const BillsOverview = () => {
  const [bills, setBills] = useState<Bill[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [totalMonthly, setTotalMonthly] = useState(0);

  useEffect(() => {
    setLoading(true);
    fetch('http://localhost:8080/api/analytics/spending?accountId=1234567891&timeRange=1%20month')
      .then(res => {
        if (!res.ok) {
          throw new Error(`Failed to fetch bills: ${res.status} ${res.statusText}`);
        }
        return res.json();
      })
      .then(data => {
        // Filter only bill transactions
        const billTransactions = (data.transactions || [])
          .filter((t: any) => BILL_CATEGORIES.has(t.category))
          .map((t: any) => ({
            category: t.category,
            amount: Math.abs(t.amount),
            dueDate: new Date(t.date).toLocaleDateString(),
            merchant: t.merchant,
            status: t.status
          }));

        // Calculate total monthly bills
        const total = billTransactions.reduce((sum: number, bill: Bill) => sum + bill.amount, 0);
        
        setBills(billTransactions);
        setTotalMonthly(total);
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
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold">Monthly Bills 📅</h2>
        <div className="text-right">
          <div className="text-sm text-gray-500">Total Monthly</div>
          <div className="text-2xl font-bold">${totalMonthly.toFixed(2)}</div>
        </div>
      </div>

      <div className="grid gap-4">
        {bills.map((bill, index) => (
          <motion.div
            key={`${bill.category}-${bill.dueDate}`}
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.3, delay: index * 0.1 }}
          >
            <Card className="p-4 hover:shadow-lg transition-shadow">
              <div className="flex justify-between items-center">
                <div className="flex items-center space-x-4">
                  <div className="p-2 bg-primary/10 rounded-full">
                    <DollarSign className="w-5 h-5 text-primary" />
                  </div>
                  <div>
                    <h3 className="font-medium">{bill.category}</h3>
                    <p className="text-sm text-gray-500">{bill.merchant}</p>
                  </div>
                </div>
                <div className="text-right">
                  <div className="font-bold">${bill.amount.toFixed(2)}</div>
                  <div className="flex items-center text-sm text-gray-500">
                    <Calendar className="w-4 h-4 mr-1" />
                    Due: {bill.dueDate}
                  </div>
                </div>
              </div>
              
              <div className="mt-2 flex items-center justify-between">
                <div className="flex items-center">
                  <Clock className="w-4 h-4 mr-1 text-gray-500" />
                  <span className="text-sm text-gray-500">Monthly</span>
                </div>
                <div className={`text-sm ${
                  bill.status === 'Completed' ? 'text-green-500' : 'text-yellow-500'
                }`}>
                  {bill.status}
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