import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { Analytics } from "./types";
import { Loader2, TrendingUp, CreditCard, Wallet, ArrowUpRight, Sparkles, DollarSign } from "lucide-react";
import { getSpendingOverview } from "../lib/api";

const cardVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.5,
      ease: [0.6, -0.05, 0.01, 0.99]
    }
  },
  hover: {
    scale: 1.02,
    transition: {
      duration: 0.2,
      ease: "easeInOut"
    }
  }
};

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
        <motion.div
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ 
            scale: [0.8, 1, 0.8],
            opacity: [0, 1, 0.8]
          }}
          transition={{ 
            duration: 2,
            repeat: Infinity,
            ease: "easeInOut"
          }}
          className="relative"
        >
          <Loader2 className="w-12 h-12 text-[#00C805]" />
          <motion.div
            className="absolute inset-0 rounded-full bg-[#00C805]/20 blur-xl"
            animate={{
              scale: [1, 1.5, 1],
              opacity: [0.5, 0.2, 0.5]
            }}
            transition={{
              duration: 2,
              repeat: Infinity,
              ease: "easeInOut"
            }}
          />
        </motion.div>
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

  const getSpendingStatus = () => {
    const monthlyAvg = analytics?.averageTransaction || 0;
    const total = analytics?.totalSpent || 0;
    
    if (total < monthlyAvg * 0.8) {
      return {
        message: "Under Average 🎯",
        color: "text-green-500",
        trend: "down"
      };
    } else if (total > monthlyAvg * 1.2) {
      return {
        message: "Over Average ⚠️",
        color: "text-red-500",
        trend: "up"
      };
    }
    return {
      message: "On Track ✨",
      color: "text-[#00C805]",
      trend: "stable"
    };
  };

  const spendingStatus = getSpendingStatus();

  return (
    <div className="space-y-8">
      {/* Financial Overview Cards */}
      <div className="grid gap-6 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        <motion.div
          variants={cardVariants}
          initial="hidden"
          animate="visible"
          whileHover="hover"
          className="col-span-1"
        >
          <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
            <motion.div
              className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-opacity duration-500 blur-xl"
              initial={false}
            />
            <div className="relative">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h3 className="text-lg font-medium text-zinc-400">Total Spent</h3>
                  <p className="text-3xl font-bold mt-1 text-white">${analytics?.totalSpent.toFixed(2) || '0.00'}</p>
                </div>
                <div className="w-12 h-12 rounded-full bg-gradient-to-br from-[#00C805]/20 to-[#00C805]/5 flex items-center justify-center backdrop-blur-xl">
                  <DollarSign className="w-6 h-6 text-[#00C805]" />
                </div>
              </div>
              <div className={`flex items-center ${spendingStatus.color}`}>
                <span className="text-sm font-medium">{spendingStatus.message}</span>
                {spendingStatus.trend === "up" && <ArrowUpRight className="w-4 h-4 ml-1" />}
              </div>
            </div>
          </Card>
        </motion.div>

        <motion.div
          variants={cardVariants}
          initial="hidden"
          animate="visible"
          whileHover="hover"
          className="col-span-1"
        >
          <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
            <motion.div
              className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-opacity duration-500 blur-xl"
              initial={false}
            />
            <div className="relative">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h3 className="text-lg font-medium text-zinc-400">Monthly Average</h3>
                  <p className="text-3xl font-bold mt-1 text-white">${analytics?.averageTransaction.toFixed(2) || '0.00'}</p>
                </div>
                <div className="w-12 h-12 rounded-full bg-gradient-to-br from-[#00C805]/20 to-[#00C805]/5 flex items-center justify-center backdrop-blur-xl">
                  <CreditCard className="w-6 h-6 text-[#00C805]" />
                </div>
              </div>
              <div className="text-zinc-400 text-sm">
                Based on last 30 days
              </div>
            </div>
          </Card>
        </motion.div>

        <motion.div
          variants={cardVariants}
          initial="hidden"
          animate="visible"
          whileHover="hover"
          className="col-span-1"
        >
          <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
            <motion.div
              className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-opacity duration-500 blur-xl"
              initial={false}
            />
            <div className="relative">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h3 className="text-lg font-medium text-zinc-400">Active Categories</h3>
                  <p className="text-3xl font-bold mt-1 text-white">{analytics?.transactionCount || 0}</p>
                </div>
                <div className="w-12 h-12 rounded-full bg-gradient-to-br from-[#00C805]/20 to-[#00C805]/5 flex items-center justify-center backdrop-blur-xl">
                  <Wallet className="w-6 h-6 text-[#00C805]" />
                </div>
              </div>
              <div className="text-zinc-400 text-sm">
                Spending categories this month
              </div>
            </div>
          </Card>
        </motion.div>
      </div>

      {/* Financial Health Score */}
      <motion.div
        variants={cardVariants}
        initial="hidden"
        animate="visible"
        whileHover="hover"
      >
        <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
          <motion.div
            className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-opacity duration-500 blur-xl"
            initial={false}
          />
          <div className="relative">
            <div className="flex items-center mb-6">
              <TrendingUp className="w-5 h-5 mr-2 text-[#00C805]" />
              <h3 className="text-lg font-medium text-white">Financial Health Overview</h3>
              <motion.div
                className="ml-2"
                animate={{
                  rotate: 360,
                  scale: [1, 1.2, 1]
                }}
                transition={{
                  duration: 3,
                  repeat: Infinity,
                  ease: "linear"
                }}
              >
                <Sparkles className="w-4 h-4 text-[#00C805]" />
              </motion.div>
            </div>
            
            <div className="grid gap-6 grid-cols-1 md:grid-cols-2">
              <div className="space-y-4">
                <div>
                  <div className="flex justify-between items-center mb-2">
                    <span className="text-sm text-zinc-400">Spending Ratio</span>
                    <span className="text-sm font-medium text-white">
                      {((analytics?.totalSpent || 0) / (analytics?.averageTransaction || 1) * 100).toFixed(1)}%
                    </span>
                  </div>
                  <div className="h-2 bg-black/50 rounded-full overflow-hidden backdrop-blur-xl">
                    <motion.div
                      className="h-full bg-gradient-to-r from-[#00C805] to-emerald-600"
                      initial={{ width: 0 }}
                      animate={{ width: `${Math.min(((analytics?.totalSpent || 0) / (analytics?.averageTransaction || 1) * 100), 100)}%` }}
                      transition={{ duration: 1, ease: "easeOut" }}
                    />
                  </div>
                </div>

                <div>
                  <div className="flex justify-between items-center mb-2">
                    <span className="text-sm text-zinc-400">Category Diversity</span>
                    <span className="text-sm font-medium text-white">
                      {Math.min(((analytics?.transactionCount || 0) / 10 * 100), 100).toFixed(1)}%
                    </span>
                  </div>
                  <div className="h-2 bg-black/50 rounded-full overflow-hidden backdrop-blur-xl">
                    <motion.div
                      className="h-full bg-gradient-to-r from-[#00C805] to-emerald-600"
                      initial={{ width: 0 }}
                      animate={{ width: `${Math.min(((analytics?.transactionCount || 0) / 10 * 100), 100)}%` }}
                      transition={{ duration: 1, ease: "easeOut" }}
                    />
                  </div>
                </div>
              </div>

              <div className="space-y-4">
                <div className="p-4 bg-black/50 rounded-lg backdrop-blur-xl">
                  <h4 className="text-sm font-medium mb-2 text-white">Quick Tips 💡</h4>
                  <ul className="text-sm text-zinc-400 space-y-2">
                    <li className="flex items-start">
                      <span className="mr-2">•</span>
                      {spendingStatus.trend === "up" 
                        ? <span className="text-red-400">Consider reviewing your discretionary spending</span>
                        : <span className="text-[#00C805]">Great job maintaining your spending habits!</span>}
                    </li>
                    <li className="flex items-start">
                      <span className="mr-2">•</span>
                      {analytics?.transactionCount && analytics.transactionCount < 5
                        ? <span className="text-yellow-400">Try diversifying your spending categories</span>
                        : <span className="text-[#00C805]">Good category distribution!</span>}
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </Card>
      </motion.div>
    </div>
  );
};

export default SpendingOverview; 