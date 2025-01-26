import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { Clock, TrendingUp, Loader2 } from "lucide-react";
import { getSpendingPatterns } from "../lib/api";
// Categories that are considered bills/rent and should be excluded
const EXCLUDED_CATEGORIES = new Set([
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

interface TimePattern {
  timeOfDay: string;
  dayOfWeek: string;
  frequency: number;
  averageSpend: number;
  category: string;
}

const SpendingPatterns = () => {
  const [patterns, setPatterns] = useState<TimePattern[]>([]);
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
    getSpendingPatterns(accountId)
      .then(data => {
        if (!Array.isArray(data)) {
          throw new Error('Invalid data format received from server');
        }
        // Filter out bills and rent from patterns
        const filteredPatterns = data.filter(pattern => !EXCLUDED_CATEGORIES.has(pattern.category));
        setPatterns(filteredPatterns);
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
          <p>{error}</p>
        </div>
      </Card>
    );
  }

  // Group patterns by day of week
  const dayPatterns = patterns.reduce<{ [key: string]: TimePattern[] }>((acc, pattern) => {
    if (!acc[pattern.dayOfWeek]) {
      acc[pattern.dayOfWeek] = [];
    }
    acc[pattern.dayOfWeek].push(pattern);
    return acc;
  }, {});

  // Calculate daily totals (excluding bills/rent)
  const dailyTotals = Object.entries(dayPatterns).map(([day, patterns]) => ({
    day,
    totalSpend: patterns.reduce((sum, p) => sum + p.averageSpend * p.frequency, 0),
    totalTransactions: patterns.reduce((sum, p) => sum + p.frequency, 0)
  })).sort((a, b) => b.totalSpend - a.totalSpend);

  return (
    <div className="grid gap-4">
      {/* Day of Week Analysis */}
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3 }}
      >
        <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
          <motion.div
            className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-opacity duration-500 blur-xl"
            initial={false}
          />
          <div className="relative">
            <div className="flex items-center mb-4">
              <TrendingUp className="w-5 h-5 mr-2 text-[#00C805]" />
              <h3 className="text-lg font-medium text-white">Discretionary Spending by Day</h3>
            </div>
            <div className="space-y-4">
              {dailyTotals.map(({ day, totalSpend, totalTransactions }) => (
                <div key={day} className="flex justify-between items-center">
                  <div>
                    <span className="font-medium text-white">{day}</span>
                    <span className="text-sm text-zinc-400 ml-2">({totalTransactions} transactions)</span>
                  </div>
                  <span className="font-medium text-white">${totalSpend.toFixed(2)}</span>
                </div>
              ))}
            </div>
          </div>
        </Card>
      </motion.div>

      {/* Time of Day Analysis */}
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3, delay: 0.1 }}
      >
        <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
          <motion.div
            className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-opacity duration-500 blur-xl"
            initial={false}
          />
          <div className="relative">
            <div className="flex items-center mb-4">
              <Clock className="w-5 h-5 mr-2 text-[#00C805]" />
              <h3 className="text-lg font-medium text-white">Most Active Shopping Times</h3>
            </div>
            <div className="space-y-4">
              {patterns
                .sort((a, b) => b.frequency - a.frequency)
                .slice(0, 5)
                .map((pattern) => (
                  <div key={`${pattern.dayOfWeek}-${pattern.timeOfDay}`} className="flex justify-between items-center">
                    <div>
                      <span className="font-medium text-white">{pattern.dayOfWeek}s at {pattern.timeOfDay}</span>
                      <span className="text-sm text-zinc-400 ml-2">({pattern.frequency} times)</span>
                    </div>
                    <span className="font-medium text-white">${pattern.averageSpend.toFixed(2)} avg</span>
                  </div>
                ))}
            </div>
          </div>
        </Card>
      </motion.div>
    </div>
  );
};

export default SpendingPatterns; 