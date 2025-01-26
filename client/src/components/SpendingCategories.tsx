import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { PieChart } from "lucide-react";
import { getCategoryTotals } from "../lib/api";

const EXCLUDED_CATEGORIES = new Set([
  'Rent', 
  'Income',
  'Utilities',
  'Insurance',
  'Phone Bill',
  'Internet',
  'Mortgage',
  'Water Bill',
  'Electric Bill',
  'Gas Bill'
]);

interface CategoryData {
  [key: string]: number;
}

const SpendingCategories = () => {
  const [categories, setCategories] = useState<CategoryData>({});
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
    getCategoryTotals(accountId)
      .then(data => {
        if (!data || typeof data !== 'object') {
          throw new Error('Invalid data format received from server');
        }
        const categoryData: CategoryData = {};
        // Filter out excluded categories and format the data
        Object.entries(data).forEach(([category, amount]) => {
          if (!EXCLUDED_CATEGORIES.has(category)) {
            categoryData[category] = parseFloat(amount as string) || 0;
          }
        });
        setCategories(categoryData);
        setError(null);
      })
      .catch(err => {
        console.error('Error fetching categories:', err);
        setError(err.message);
      })
      .finally(() => setLoading(false));
  }, []);

  // Calculate total spending for percentages
  const totalSpending = Object.values(categories).reduce((sum, value) => sum + Math.abs(value), 0);

  if (loading) {
    return (
      <Card className="p-6">
        <div className="text-center">Loading categories...</div>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className="p-6">
        <div className="text-center text-red-500">
          Error loading categories: {error}
        </div>
      </Card>
    );
  }

  return (
    <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
      <motion.div
        className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-opacity duration-500 blur-xl"
        initial={false}
      />
      <div className="relative">
        <div className="flex items-center mb-6">
          <PieChart className="w-5 h-5 mr-2 text-[#00C805]" />
          <h3 className="text-lg font-medium text-white">Spending by Category</h3>
        </div>
        <div className="space-y-6">
          {Object.entries(categories)
            .sort(([, a], [, b]) => Math.abs(b) - Math.abs(a))
            .map(([category, amount]) => {
              const percentage = (Math.abs(amount) / totalSpending) * 100;
              return (
                <motion.div
                  key={category}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.3 }}
                  className="space-y-2"
                >
                  <div className="flex justify-between items-center">
                    <span className="font-medium text-white">{category}</span>
                    <div className="text-right">
                      <div className="font-medium text-white">${Math.abs(amount).toFixed(2)}</div>
                      <div className="text-sm text-zinc-400">{percentage.toFixed(1)}%</div>
                    </div>
                  </div>
                  <div className="h-2 bg-black/50 rounded-full overflow-hidden backdrop-blur-xl">
                    <motion.div
                      className="h-full bg-gradient-to-r from-[#00C805] to-emerald-600"
                      initial={{ width: 0 }}
                      animate={{ width: `${percentage}%` }}
                      transition={{ duration: 1, ease: "easeOut" }}
                    />
                  </div>
                </motion.div>
              );
            })}
        </div>
      </div>
    </Card>
  );
};

export default SpendingCategories; 