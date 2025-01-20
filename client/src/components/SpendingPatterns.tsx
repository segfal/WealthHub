import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { Clock, TrendingUp, AlertCircle, Loader2 } from "lucide-react";
import { Pattern } from "./types";

const SpendingPatterns = () => {
  const [patterns, setPatterns] = useState<Pattern | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    fetch('http://localhost:8080/api/analytics/patterns?accountId=1234567891')
      .then(res => {
        if (!res.ok) {
          throw new Error(`Failed to fetch patterns: ${res.status} ${res.statusText}`);
        }
        return res.json();
      })
      .then(data => {
        if (!data || !data.timePatterns) {
          throw new Error('Invalid data format received from server');
        }
        // Transform the data to match the Pattern interface
        const transformedData: Pattern = {
          timeOfDay: {
            morning: 0,
            afternoon: 0,
            evening: 0,
            night: 0,
            ...data.timePatterns.reduce((acc: any, curr: any) => {
              const hour = parseInt(curr.hour);
              if (hour >= 5 && hour < 12) acc.morning = (acc.morning || 0) + curr.totalSpent;
              else if (hour >= 12 && hour < 17) acc.afternoon = (acc.afternoon || 0) + curr.totalSpent;
              else if (hour >= 17 && hour < 22) acc.evening = (acc.evening || 0) + curr.totalSpent;
              else acc.night = (acc.night || 0) + curr.totalSpent;
              return acc;
            }, {})
          },
          dayOfWeek: data.timePatterns.reduce((acc: any, curr: any) => {
            acc[curr.dayOfWeek] = (acc[curr.dayOfWeek] || 0) + curr.totalSpent;
            return acc;
          }, {}),
          recurringTransactions: data.recurringTransactions?.map((t: any) => ({
            merchant: t.merchant || 'Unknown',
            amount: t.amount || 0,
            frequency: t.frequency || 'Monthly'
          })) || []
        };
        setPatterns(transformedData);
        setError(null);
      })
      .catch(err => {
        console.error('Error fetching patterns:', err);
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
          <AlertCircle className="w-8 h-8 mx-auto mb-2" />
          <p>{error}</p>
        </div>
      </Card>
    );
  }

  return (
    <div className="grid gap-4">
      {/* Time of Day Analysis */}
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3 }}
      >
        <Card className="p-6">
          <div className="flex items-center mb-4">
            <Clock className="w-5 h-5 mr-2" />
            <h3 className="text-lg font-medium">Spending by Time of Day</h3>
          </div>
          <div className="space-y-4">
            {patterns?.timeOfDay && Object.entries(patterns.timeOfDay)
              .sort(([, a], [, b]) => b - a)
              .map(([time, amount]) => (
                <div key={time} className="flex justify-between items-center">
                  <span className="capitalize">{time}</span>
                  <span className="font-medium">${amount.toFixed(2)}</span>
                </div>
              ))}
          </div>
        </Card>
      </motion.div>

      {/* Day of Week Analysis */}
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3, delay: 0.1 }}
      >
        <Card className="p-6">
          <div className="flex items-center mb-4">
            <TrendingUp className="w-5 h-5 mr-2" />
            <h3 className="text-lg font-medium">Spending by Day of Week</h3>
          </div>
          <div className="space-y-4">
            {patterns?.dayOfWeek && Object.entries(patterns.dayOfWeek)
              .sort(([, a], [, b]) => b - a)
              .map(([day, amount]) => (
                <div key={day} className="flex justify-between items-center">
                  <span>{day}</span>
                  <span className="font-medium">${amount.toFixed(2)}</span>
                </div>
              ))}
          </div>
        </Card>
      </motion.div>

      {/* Recurring Transactions */}
      <motion.div
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.3, delay: 0.2 }}
      >
        <Card className="p-6">
          <div className="flex items-center mb-4">
            <AlertCircle className="w-5 h-5 mr-2" />
            <h3 className="text-lg font-medium">Recurring Transactions</h3>
          </div>
          <div className="space-y-4">
            {patterns?.recurringTransactions?.map((transaction, index) => (
              <div key={index} className="flex justify-between items-center p-4 bg-secondary rounded-lg">
                <div>
                  <p className="font-medium">{transaction.merchant}</p>
                  <p className="text-sm text-muted-foreground">{transaction.frequency}</p>
                </div>
                <p className="font-medium">${transaction.amount.toFixed(2)}</p>
              </div>
            ))}
          </div>
        </Card>
      </motion.div>
    </div>
  );
};

export default SpendingPatterns; 