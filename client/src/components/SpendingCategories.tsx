import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { Progress } from "./ui/progress";
import { CategoryData } from "./types";
import { getSpendingCategories } from "../lib/api";

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

interface PredictedSpend {
  category: string;
  likelihood: number;
  predictedDate: string;
  warning: string;
  amount: number;
}

const SpendingCategories = () => {
  const [predictions, setPredictions] = useState<PredictedSpend[]>([]);
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
    getSpendingCategories(accountId)
      .then(data => {
        if (!data || !data.topCategories) {
          throw new Error('Invalid data format received from server');
        }
        const categoryData: CategoryData = {}; 
        data.topCategories.forEach((cat: any) => { 
          if(cat.category != "Income"){ 
            categoryData[cat.category || 'Other'] = parseFloat(cat.totalSpent) || 0; 
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
    <Card className="p-6">
      <h3 className="text-lg font-medium mb-6">Spending by Category</h3>
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
                  <span className="font-medium">{category}</span>
                  <div className="text-right">
                    <div className="font-medium">${Math.abs(amount).toFixed(2)}</div>
                    <div className="text-sm text-gray-500">{percentage.toFixed(1)}%</div>
                  </div>
                </div>
                <Progress value={percentage} className="h-2" />
              </motion.div>
            );
          })}
      </div>
    </Card>
  );
};

export default SpendingCategories; 