import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

interface Analytics {
  totalSpent: number;
  averageTransaction: number;
  transactionCount: number;
  spendingTrend: { date: string; amount: number }[];
}

const SpendingOverview = () => {
  const [analytics, setAnalytics] = useState<Analytics | null>(null);

  useEffect(() => {
    // Fetch analytics data
    fetch('http://localhost:8080/api/analytics/1234567890')
      .then(res => res.json())
      .then(data => setAnalytics(data))
      .catch(err => console.error('Error fetching analytics:', err));
  }, []);

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
          <p className="text-3xl font-bold">${analytics?.totalSpent || 0}</p>
        </Card>
      </motion.div>

      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3, delay: 0.1 }}
      >
        <Card className="p-6">
          <h3 className="text-lg font-medium">Average Transaction</h3>
          <p className="text-3xl font-bold">${analytics?.averageTransaction || 0}</p>
        </Card>
      </motion.div>

      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3, delay: 0.2 }}
      >
        <Card className="p-6">
          <h3 className="text-lg font-medium">Transaction Count</h3>
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
          <h3 className="text-lg font-medium mb-4">Spending Trend</h3>
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