import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { Analytics } from "./types";
import { Loader2 } from "lucide-react";
import { getSpendingOverview } from "../lib/api";
  
const SpendingOverview = () => {
  const [analytics, setAnalytics] = useState<Analytics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const accountId = import.meta.env.VITE_ACCOUNT_ID;
    if (!accountId) {
      setError("No account ID provided");
      setLoading(false);
      return;
    }

    setLoading(true);
    getSpendingOverview(accountId)
      .then(data => {
        if (!data) {
          throw new Error('Invalid data format received from server');
        }
        setAnalytics({
          totalSpent: data.totalSpent || 0,
          averageTransaction: data.monthlyAverage || 0,
          transactionCount: data.topCategories?.length || 0,
          spendingTrend: data.spendingPatterns?.map((pattern: any) => ({
            date: pattern.dayOfWeek || 'Unknown',
            amount: pattern.totalSpent || 0
          })) || []
        });
        setError(null);
      })
      .catch(err => {
        console.error('Error fetching analytics:', err);
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
      <div className="text-center text-red-500 p-4">
        <p>Error: {error}</p>
      </div>
    );
  }

  return (
    <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
      {/* Summary Cards */}
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3 }}
      >
        <Card className="p-6">
          <h3 className="text-lg font-medium">Total Spent</h3>
          <p className="text-3xl font-bold">${analytics?.totalSpent.toFixed(2) || '0.00'}</p>
        </Card>
      </motion.div>

      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3, delay: 0.1 }}
      >
        <Card className="p-6">
          <h3 className="text-lg font-medium">Monthly Average</h3>
          <p className="text-3xl font-bold">${analytics?.averageTransaction.toFixed(2) || '0.00'}</p>
        </Card>
      </motion.div>

      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3, delay: 0.2 }}
      >
        <Card className="p-6">
          <h3 className="text-lg font-medium">Categories</h3>
          <p className="text-3xl font-bold">{analytics?.transactionCount || 0}</p>
        </Card>
      </motion.div>

      {/* Spending Trend Chart */}
      <motion.div
        className="col-span-full"
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.5, delay: 0.3 }}
      >
        <Card className="p-6">
          <h3 className="text-lg font-medium mb-4">Spending by Day</h3>
          <div className="h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={analytics?.spendingTrend || []}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="amount" fill="#8884d8" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </Card>
      </motion.div>
    </div>
  );
};

export default SpendingOverview; 