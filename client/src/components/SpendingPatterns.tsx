import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { Clock, TrendingUp, AlertCircle, Loader2 } from "lucide-react";

interface Pattern {
  timeOfDay: {
    morning: number;
    afternoon: number;
    evening: number;
    night: number;
  };
  dayOfWeek: {
    [key: string]: number;
  };
  recurringTransactions: {
    merchant: string;
    amount: number;
    frequency: string;
  }[];
}

const SpendingPatterns = () => {
  const [patterns, setPatterns] = useState<Pattern | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    fetch('http://localhost:8080/api/patterns/1234567890')
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch patterns');
        return res.json();
      })
      .then(data => {
        setPatterns(data);
        setError(null);
      })
      .catch(err => {
        console.error('Error fetching patterns:', err);
        setError(err.message);
        // Set dummy data for development
        setPatterns({
          timeOfDay: {
            morning: 250,
            afternoon: 350,
            evening: 450,
            night: 150
          },
          dayOfWeek: {
            Monday: 300,
            Tuesday: 400,
            Wednesday: 350,
            Thursday: 500,
            Friday: 600,
            Saturday: 450,
            Sunday: 300
          },
          recurringTransactions: [
            {
              merchant: "Netflix",
              amount: 14.99,
              frequency: "Monthly"
            },
            {
              merchant: "Gym Membership",
              amount: 49.99,
              frequency: "Monthly"
            },
            {
              merchant: "Spotify",
              amount: 9.99,
              frequency: "Monthly"
            }
          ]
        });
      })
      .finally(() => setLoading(false));
  }, []);

  const timeData = patterns?.timeOfDay ? [
    { name: 'Morning', value: patterns.timeOfDay.morning },
    { name: 'Afternoon', value: patterns.timeOfDay.afternoon },
    { name: 'Evening', value: patterns.timeOfDay.evening },
    { name: 'Night', value: patterns.timeOfDay.night },
  ] : [];

  const weekData = patterns?.dayOfWeek ? 
    Object.entries(patterns.dayOfWeek).map(([day, value]) => ({
      name: day,
      value: value
    })) : [];

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
    <div className="grid gap-4 grid-cols-1 lg:grid-cols-2">
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
          <div className="h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={timeData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="value" fill="#8884d8" />
              </BarChart>
            </ResponsiveContainer>
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
          <div className="h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={weekData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="value" fill="#82ca9d" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </Card>
      </motion.div>

      {/* Recurring Transactions */}
      <motion.div
        className="lg:col-span-2"
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