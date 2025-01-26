import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { AlertTriangle, TrendingDown, TrendingUp, ShoppingBag, Coffee, Utensils, Plane, Book, Monitor, Gift, ShoppingCart, CreditCard, Loader2 } from "lucide-react";
import { getSpendingInsights } from "../lib/api";
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

// Emoji mappings for different merchants/categories
const MERCHANT_EMOJIS: { [key: string]: string } = {
  // Shopping
  'Amazon.com': '📦',
  'Target': '🎯',
  'Walmart': '🛒',
  'Best Buy': '🔌',
  'Home Depot': '🏠',
  'IKEA': '🪑',
  
  // Entertainment
  'Netflix': '🎬',
  'Spotify': '🎵',
  'Apple Music': '🎧',
  'Steam': '🎮',
  'PlayStation': '🕹️',
  'Xbox': '🎯',
  'AMC Theaters': '🍿',
  'Regal Cinemas': '🎦',
  
  // Transportation
  'Uber': '🚗',
  'Lyft': '🚙',
  'Shell': '⛽',
  'Chevron': '⛽',
  'Delta Airlines': '✈️',
  'United Airlines': '✈️',
  'American Airlines': '✈️',
  
  // Food & Dining
  'Uber Eats': '🥡',
  'DoorDash': '🛵',
  'GrubHub': '🍽️',
  'Starbucks': '☕',
  'Dunkin': '🍩',
  'McDonalds': '🍔',
  'Chipotle': '🌯',
  'Subway': '🥖',
  'Pizza Hut': '🍕',
  'Dominos': '🍕',
  
  // Travel & Hotels
  'Airbnb': '🏡',
  'Hotels.com': '🏨',
  'Marriott': '🏨',
  'Hilton': '🏨',
  'Expedia': '🌎',
  
  // General Categories
  'Restaurant': '🍽️',
  'Bar': '🍺',
  'Grocery': '🛒',
  'Clothing': '👕',
  'Entertainment': '🎭',
  'Books': '📚',
  'Online': '💻',
  'Pharmacy': '💊',
  'Health': '🏥',
  'Fitness': '🏋️',
  'Sports': '⚽',
  'Education': '📚',
  'Pet Supplies': '🐾',
  'Beauty': '💄',
  'Gaming': '🎮',
  'Music': '🎵',
  'Movies': '🎬',
  'Coffee Shop': '☕'
};

interface Transaction {
  category: string;
  merchant: string;
  amount: number;
  location: string;
  date: string;
}

interface CategoryData {
  category: string;
  totalSpent: string;
  percentage: string;
  trend?: 'increasing' | 'decreasing' | 'stable';
  frequency?: number;
  merchants?: Array<{
    name: string;
    amount: number;
    count: number;
  }>;
}

interface InsightData {
  type: string;
  title: string;
  description: string;
  data: CategoryData[];
}

interface InsightResponse {
  insights: InsightData[];
  totalSpent: number;
  monthlyAverage: number;
}

const SpendingInsights = () => {
  const [insights, setInsights] = useState<CategoryData[]>([]);
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
    getSpendingInsights(accountId)
      .then((response: InsightResponse) => {
        if (!response || !response.insights || !response.insights[0]?.data) {
          throw new Error('Invalid data format received from server');
        }
        
        // Get the top categories data and filter out excluded categories
        const categoryData = response.insights[0].data.filter(
          item => !EXCLUDED_CATEGORIES.has(item.category)
        );
        
        setInsights(categoryData);
        setError(null);
      })
      .catch(err => {
        console.error('Error fetching insights:', err);
        setError(err.message);
      })
      .finally(() => setLoading(false));
  }, []);

  const getCategoryIcon = (category: string) => {
    switch (category.toLowerCase()) {
      case 'groceries': return <ShoppingCart className="w-6 h-6 text-[#4ADE80]" />;
      case 'dining': 
      case 'restaurants': return <Utensils className="w-6 h-6 text-[#4ADE80]" />;
      case 'shopping': return <ShoppingBag className="w-6 h-6 text-[#4ADE80]" />;
      case 'coffee': return <Coffee className="w-6 h-6 text-[#4ADE80]" />;
      case 'travel': return <Plane className="w-6 h-6 text-[#4ADE80]" />;
      case 'entertainment': return <Monitor className="w-6 h-6 text-[#4ADE80]" />;
      case 'subscription': return <CreditCard className="w-6 h-6 text-[#4ADE80]" />;
      default: return <CreditCard className="w-6 h-6 text-[#4ADE80]" />;
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="w-8 h-8 animate-spin text-[#4ADE80]" />
      </div>
    );
  }

  if (error) {
    return (
      <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#4ADE80]/20 relative group shadow-lg shadow-[#4ADE80]/5">
        <div className="text-center text-red-400">
          <AlertTriangle className="w-8 h-8 mx-auto mb-2" />
          <p>{error}</p>
        </div>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-white flex items-center gap-2">
        Spending Insights <span className="text-[#4ADE80]">🔍</span>
      </h2>
      <p className="text-zinc-400">Your highest spending areas (excluding bills)</p>

      <div className="space-y-4">
        {insights.map((item, index) => (
          <motion.div
            key={item.category}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: index * 0.1 }}
          >
            <Card className="p-4 bg-[#1a1d21] backdrop-blur-xl border border-[#4ADE80]/20 relative group shadow-lg shadow-[#4ADE80]/5">
              <motion.div
                className="absolute inset-0 bg-gradient-radial from-[#4ADE80]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-opacity duration-500 blur-xl"
                initial={false}
              />
              <div className="relative">
                <div className="flex justify-between items-center">
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 rounded-full bg-gradient-to-br from-[#4ADE80]/20 to-[#4ADE80]/5 flex items-center justify-center backdrop-blur-xl">
                      {getCategoryIcon(item.category)}
                    </div>
                    <div>
                      <h3 className="font-medium text-white">{item.category}</h3>
                      <div className="flex items-center text-sm text-zinc-400 mt-1">
                        {item.trend === 'increasing' ? (
                          <TrendingUp className="w-4 h-4 text-[#4ADE80] mr-1" />
                        ) : (
                          <TrendingDown className="w-4 h-4 text-[#4ADE80] mr-1" />
                        )}
                        <span>
                          {item.trend === 'increasing' ? 'Increasing' : 'Decreasing'} trend
                        </span>
                      </div>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-lg font-medium text-white">
                      ${parseFloat(item.totalSpent).toFixed(2)}
                    </div>
                    <div className="text-sm text-zinc-400">
                      {parseFloat(item.percentage).toFixed(1)}% of total
                    </div>
                  </div>
                </div>
                <div className="mt-2 h-2 bg-black/50 rounded-full overflow-hidden backdrop-blur-xl">
                  <motion.div
                    className="h-full bg-gradient-to-r from-[#4ADE80] to-emerald-600"
                    initial={{ width: 0 }}
                    animate={{ width: `${parseFloat(item.percentage)}%` }}
                    transition={{ duration: 0.5, delay: 0.6 + index * 0.1 }}
                  />
                </div>
              </div>
            </Card>
          </motion.div>
        ))}
      </div>
    </div>
  );
};

export default SpendingInsights; 