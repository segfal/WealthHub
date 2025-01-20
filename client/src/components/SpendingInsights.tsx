import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { AlertTriangle, TrendingDown, TrendingUp, ShoppingBag, Coffee, Utensils, Plane, Book, Monitor, Gift } from "lucide-react";

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

interface SpendingInsight {
  category: string;
  totalSpent: number;
  frequency: number;
  merchants: Array<{
    name: string;
    amount: number;
    count: number;
  }>;
  trend: 'increasing' | 'decreasing' | 'stable';
}

const SpendingInsights = () => {
  const [insights, setInsights] = useState<SpendingInsight[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    fetch('http://localhost:8080/api/analytics/spending?accountId=1234567891&timeRange=1%20month')
      .then(res => {
        if (!res.ok) {
          throw new Error(`Failed to fetch insights: ${res.status} ${res.statusText}`);
        }
        return res.json();
      })
      .then(data => {
        // Filter out bills and process spending data
        const processedInsights = processTransactions(data.transactions || []);
        setInsights(processedInsights);
        setError(null);
      })
      .catch(err => {
        console.error('Error fetching insights:', err);
        setError(err.message);
      })
      .finally(() => setLoading(false));
  }, []);

  const processTransactions = (transactions: Transaction[]): SpendingInsight[] => {
    const categoryMap = new Map<string, {
      total: number;
      count: number;
      merchants: Map<string, { amount: number; count: number }>;
      previousTotal: number;
    }>();

    // Group transactions by category
    transactions
      .filter(t => !EXCLUDED_CATEGORIES.has(t.category))
      .forEach(t => {
        if (!categoryMap.has(t.category)) {
          categoryMap.set(t.category, {
            total: 0,
            count: 0,
            merchants: new Map(),
            previousTotal: 0
          });
        }
        const category = categoryMap.get(t.category)!;
        const amount = Math.abs(t.amount);
        category.total += amount;
        category.count++;

        if (!category.merchants.has(t.merchant)) {
          category.merchants.set(t.merchant, { amount: 0, count: 0 });
        }
        const merchant = category.merchants.get(t.merchant)!;
        merchant.amount += amount;
        merchant.count++;
      });

    // Convert to array and sort by total spent
    return Array.from(categoryMap.entries())
      .map(([category, data]) => ({
        category,
        totalSpent: data.total,
        frequency: data.count,
        merchants: Array.from(data.merchants.entries())
          .map(([name, stats]) => ({
            name,
            amount: stats.amount,
            count: stats.count
          }))
          .sort((a, b) => b.amount - a.amount),
        trend: data.total > data.previousTotal ? 'increasing' : 
               data.total < data.previousTotal ? 'decreasing' : 'stable'
      }))
      .sort((a, b) => b.totalSpent - a.totalSpent);
  };

  const getCategoryIcon = (category: string) => {
    switch (category.toLowerCase()) {
      case 'shopping': return <ShoppingBag className="w-5 h-5" />;
      case 'coffee': return <Coffee className="w-5 h-5" />;
      case 'dining': 
      case 'restaurants': 
      case 'food': return <Utensils className="w-5 h-5" />;
      case 'travel': return <Plane className="w-5 h-5" />;
      case 'books': 
      case 'education': return <Book className="w-5 h-5" />;
      case 'electronics': return <Monitor className="w-5 h-5" />;
      case 'entertainment': return <span className="text-xl">🎭</span>;
      case 'groceries': return <span className="text-xl">🛒</span>;
      case 'health': return <span className="text-xl">��</span>;
      case 'fitness': return <span className="text-xl">🏋️</span>;
      case 'pets': return <span className="text-xl">🐾</span>;
      case 'beauty': return <span className="text-xl">💄</span>;
      case 'gaming': return <span className="text-xl">🎮</span>;
      case 'music': return <span className="text-xl">🎵</span>;
      case 'movies': return <span className="text-xl">🎬</span>;
      default: return <Gift className="w-5 h-5" />;
    }
  };

  const getSpendingIndicator = (amount: number, frequency: number) => {
    const monthlyThreshold = 1000; // Adjust these thresholds as needed
    const frequencyThreshold = 20;

    if (amount > monthlyThreshold && frequency > frequencyThreshold) {
      return (
        <div className="flex items-center text-red-500">
          <AlertTriangle className="w-5 h-5 mr-1" />
          <span className="text-sm">High spending alert!</span>
        </div>
      );
    } else if (amount < monthlyThreshold / 2 && frequency < frequencyThreshold / 2) {
      return (
        <div className="flex items-center text-green-500">
          <TrendingDown className="w-5 h-5 mr-1" />
          <span className="text-sm">Good spending habits!</span>
        </div>
      );
    }
    return null;
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="animate-bounce text-4xl">💰</div>
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
      <h2 className="text-2xl font-bold mb-4">Spending Insights 🔍</h2>
      
      <div className="grid gap-6">
        {insights.map((insight, index) => (
          <motion.div
            key={insight.category}
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ duration: 0.3, delay: index * 0.1 }}
          >
            <Card className="p-6 hover:shadow-lg transition-shadow">
              <div className="flex justify-between items-start mb-4">
                <div className="flex items-center">
                  <span className="mr-2">{getCategoryIcon(insight.category)}</span>
                  <div>
                    <h3 className="text-lg font-medium">{insight.category}</h3>
                    <p className="text-sm text-gray-500">
                      {insight.frequency} transactions this month
                    </p>
                  </div>
                </div>
                {getSpendingIndicator(insight.totalSpent, insight.frequency)}
              </div>

              <div className="space-y-4">
                <div>
                  <div className="text-sm text-gray-500 mb-1">Total Spent</div>
                  <div className="text-2xl font-bold">
                    ${insight.totalSpent.toFixed(2)}
                    {insight.trend === 'increasing' && <TrendingUp className="inline ml-2 text-red-500" />}
                    {insight.trend === 'decreasing' && <TrendingDown className="inline ml-2 text-green-500" />}
                  </div>
                </div>

                <div>
                  <div className="text-sm text-gray-500 mb-2">Top Merchants</div>
                  <div className="space-y-2">
                    {insight.merchants.slice(0, 3).map(merchant => (
                      <motion.div
                        key={merchant.name}
                        className="flex justify-between items-center p-2 bg-gray-50 rounded-lg"
                        whileHover={{ scale: 1.02 }}
                      >
                        <div className="flex items-center">
                          <span className="text-xl mr-2">
                            {MERCHANT_EMOJIS[merchant.name] || '🏪'}
                          </span>
                          <span>{merchant.name}</span>
                        </div>
                        <div className="text-right">
                          <div className="font-medium">${merchant.amount.toFixed(2)}</div>
                          <div className="text-sm text-gray-500">{merchant.count} times</div>
                        </div>
                      </motion.div>
                    ))}
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

export default SpendingInsights; 