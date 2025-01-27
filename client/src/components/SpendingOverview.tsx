import { useEffect, useState, useRef } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { AlertTriangle, DollarSign, CreditCard, LineChart, Loader2 } from "lucide-react";
import { getSpendingOverview } from "../lib/api";

interface SpendingData {
  account: {
    account_id: string;
    account_name: string;
    account_type: string;
    account_number: string;
    balance: {
      current: number;
      available: number;
      currency: string;
    };
  };
  top_categories: Array<{
    category: string;
    total: number;
    percentage: number;
  }>;
  spending_patterns: Array<{
    day: string;
    amount: number;
  }>;
  predicted_spending: Array<{
    month: string;
    amount: number;
  }>;
  total_spent: number;
  monthly_average: number;
  active_categories: number;
  spending_ratio: number;
  category_diversity: number;
  tips: string[];
}

const SpendingOverview = () => {
  const [spendingData, setSpendingData] = useState<SpendingData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const animateValue = (element: HTMLElement, start: number, end: number, duration: number) => {
    const startTime = performance.now();
    
    const updateValue = (currentTime: number) => {
      const elapsed = currentTime - startTime;
      const progress = Math.min(elapsed / duration, 1);
      
      // Easing function for smooth animation
      const easeOutQuart = 1 - Math.pow(1 - progress, 4);
      const current = start + (end - start) * easeOutQuart;
      
      element.textContent = current.toFixed(2);
      
      if (progress < 1) {
        requestAnimationFrame(updateValue);
      }
    };
    
    requestAnimationFrame(updateValue);
  };

  const valueRef = useRef<HTMLSpanElement>(null);
  const averageRef = useRef<HTMLSpanElement>(null);
  const categoriesRef = useRef<HTMLSpanElement>(null);

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
        if (!data || typeof data.total_spent === 'undefined' || !data.account) {
          throw new Error('Invalid spending data format');
        }
        
        // Calculate active categories from top_categories
        const activeCategories = data.top_categories?.length || 0;
        
        // Calculate spending ratio (total spent vs monthly average)
        const spendingRatio = data.monthly_average > 0 ? 
          data.total_spent / data.monthly_average : 0;
        
        // Calculate category diversity (number of categories / max categories)
        const categoryDiversity = activeCategories / 10; // Assuming max of 10 categories
        
        // Generate tips based on spending patterns
        const tips = [];
        if (spendingRatio < 0.7) {
          tips.push("Great job maintaining your spending habits!");
        } else {
          tips.push("Consider reviewing your discretionary spending");
        }
        
        if (categoryDiversity > 0.5) {
          tips.push("Good category distribution!");
        } else {
          tips.push("Try diversifying your spending categories");
        }

        setSpendingData({
          ...data,
          active_categories: activeCategories,
          spending_ratio: spendingRatio,
          category_diversity: categoryDiversity,
          tips
        });
        setError(null);
      })
      .catch(err => {
        console.error('Error fetching spending overview:', err);
        setError(err.message);
      })
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    if (!spendingData) return;

    if (valueRef.current) {
      animateValue(valueRef.current, 0, spendingData.total_spent, 2000);
    }
    if (averageRef.current) {
      animateValue(averageRef.current, 0, spendingData.monthly_average, 2000);
    }
    if (categoriesRef.current) {
      animateValue(categoriesRef.current, 0, spendingData.active_categories, 1500);
    }
  }, [spendingData]);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="w-8 h-8 animate-spin text-[#4ADE80]" />
      </div>
    );
  }

  if (error || !spendingData) {
    return (
      <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
        <div className="text-center text-red-400">
          <AlertTriangle className="w-8 h-8 mx-auto mb-2" />
          <p>{error || "Failed to load spending overview"}</p>
        </div>
      </Card>
    );
  }

  return (
    <div className="space-y-6 min-h-[500px]">
      {/* Total Spent */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <motion.div
          initial={{ opacity: 0, y: 20, scale: 0.95 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
          transition={{ duration: 0.5 }}
          whileHover={{ scale: 1.02 }}
        >
          <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#4ADE80]/20 relative group shadow-lg shadow-[#4ADE80]/5">
            <motion.div
              className="absolute inset-0 bg-gradient-radial from-[#4ADE80]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-all duration-500 blur-xl"
              initial={false}
              animate={{ rotate: [0, 360] }}
              transition={{ duration: 8, repeat: Infinity, ease: "linear" }}
            />
            <div className="relative">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium flex items-center text-white">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-br from-[#4ADE80]/20 to-[#4ADE80]/5 flex items-center justify-center backdrop-blur-xl mr-3">
                    <DollarSign className="w-5 h-5 text-[#4ADE80]" />
                  </div>
                  Total Spent
                </h3>
                <motion.span 
                  className="text-2xl font-bold text-[#4ADE80]"
                  initial={{ scale: 0.5, opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  transition={{ delay: 0.2, type: "spring" }}
                >
                  $<span ref={valueRef}>0.00</span>
                </motion.span>
              </div>
              <div className="text-sm text-zinc-400">
                {spendingData.spending_ratio < 0.7 ? "On Track" : "High Spending"}
              </div>
            </div>
          </Card>
        </motion.div>

        {/* Monthly Average */}
        <motion.div
          initial={{ opacity: 0, y: 20, scale: 0.95 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
          transition={{ duration: 0.5, delay: 0.1 }}
          whileHover={{ scale: 1.02 }}
        >
          <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#4ADE80]/20 relative group shadow-lg shadow-[#4ADE80]/5">
            <motion.div
              className="absolute inset-0 bg-gradient-radial from-[#4ADE80]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-all duration-500 blur-xl"
              initial={false}
              animate={{ rotate: [0, 360] }}
              transition={{ duration: 8, repeat: Infinity, ease: "linear" }}
            />
            <div className="relative">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium flex items-center text-white">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-br from-[#4ADE80]/20 to-[#4ADE80]/5 flex items-center justify-center backdrop-blur-xl mr-3">
                    <CreditCard className="w-5 h-5 text-[#4ADE80]" />
                  </div>
                  Monthly Average
                </h3>
                <motion.span 
                  className="text-2xl font-bold text-[#4ADE80]"
                  initial={{ scale: 0.5, opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  transition={{ delay: 0.3, type: "spring" }}
                >
                  $<span ref={averageRef}>0.00</span>
                </motion.span>
              </div>
              <div className="text-sm text-zinc-400">Based on last 30 days</div>
            </div>
          </Card>
        </motion.div>

        {/* Active Categories */}
        <motion.div
          initial={{ opacity: 0, y: 20, scale: 0.95 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
          transition={{ duration: 0.5, delay: 0.2 }}
          whileHover={{ scale: 1.02 }}
        >
          <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#4ADE80]/20 relative group shadow-lg shadow-[#4ADE80]/5">
            <motion.div
              className="absolute inset-0 bg-gradient-radial from-[#4ADE80]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-all duration-500 blur-xl"
              initial={false}
              animate={{ rotate: [0, 360] }}
              transition={{ duration: 8, repeat: Infinity, ease: "linear" }}
            />
            <div className="relative">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium flex items-center text-white">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-br from-[#4ADE80]/20 to-[#4ADE80]/5 flex items-center justify-center backdrop-blur-xl mr-3">
                    <LineChart className="w-5 h-5 text-[#4ADE80]" />
                  </div>
                  Active Categories
                </h3>
                <motion.span 
                  className="text-2xl font-bold text-[#4ADE80]"
                  initial={{ scale: 0.5, opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  transition={{ delay: 0.4, type: "spring" }}
                >
                  <span ref={categoriesRef}>0</span>
                </motion.span>
              </div>
              <div className="text-sm text-zinc-400">Spending categories this month</div>
            </div>
          </Card>
        </motion.div>
      </div>

      {/* Financial Health Overview */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.3 }}
        className="flex-1"
      >
        <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#4ADE80]/20 relative group shadow-lg shadow-[#4ADE80]/5 h-full">
          <motion.div
            className="absolute inset-0 bg-gradient-radial from-[#4ADE80]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-all duration-500 blur-xl"
            initial={false}
            animate={{ rotate: [0, 360] }}
            transition={{ duration: 8, repeat: Infinity, ease: "linear" }}
          />
          <div className="relative">
            <h3 className="text-xl font-medium flex items-center text-white mb-6">
              <LineChart className="w-5 h-5 text-[#4ADE80] mr-2" />
              Financial Health Overview
              <motion.span 
                className="text-[#4ADE80] ml-2"
                animate={{
                  scale: [1, 1.2, 1],
                  rotate: [0, 10, -10, 0],
                }}
                transition={{
                  duration: 2,
                  repeat: Infinity,
                  ease: "easeInOut"
                }}
              >
                ✨
              </motion.span>
            </h3>

            <div className="space-y-4">
              <div>
                <div className="flex justify-between items-center mb-2">
                  <span className="text-zinc-400">Spending Ratio</span>
                  <motion.span 
                    className="text-white"
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: 0.5 }}
                  >
                    {(spendingData.spending_ratio * 100).toFixed(1)}%
                  </motion.span>
                </div>
                <div className="h-2 bg-black/50 rounded-full overflow-hidden backdrop-blur-xl">
                  <motion.div
                    className="h-full bg-gradient-to-r from-[#4ADE80] to-emerald-600"
                    initial={{ width: 0 }}
                    animate={{ width: `${spendingData.spending_ratio * 100}%` }}
                    transition={{ duration: 1, delay: 0.5, ease: "easeOut" }}
                  />
                </div>
              </div>

              <div>
                <div className="flex justify-between items-center mb-2">
                  <span className="text-zinc-400">Category Diversity</span>
                  <motion.span 
                    className="text-white"
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: 0.7 }}
                  >
                    {(spendingData.category_diversity * 100).toFixed(1)}%
                  </motion.span>
                </div>
                <div className="h-2 bg-black/50 rounded-full overflow-hidden backdrop-blur-xl">
                  <motion.div
                    className="h-full bg-gradient-to-r from-[#4ADE80] to-emerald-600"
                    initial={{ width: 0 }}
                    animate={{ width: `${spendingData.category_diversity * 100}%` }}
                    transition={{ duration: 1, delay: 0.7, ease: "easeOut" }}
                  />
                </div>
              </div>
            </div>

            {/* Quick Tips */}
            <motion.div 
              className="mt-6"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.9 }}
            >
              <h4 className="text-white font-medium mb-3 flex items-center">
                Quick Tips 
                <motion.span
                  animate={{
                    rotate: [0, 15, -15, 0],
                    scale: [1, 1.2, 1]
                  }}
                  transition={{
                    duration: 2,
                    repeat: Infinity,
                    ease: "easeInOut",
                    delay: 1
                  }}
                >
                  💡
                </motion.span>
              </h4>
              <ul className="space-y-2">
                {spendingData.tips.map((tip, index) => (
                  <motion.li 
                    key={index} 
                    className="text-[#4ADE80]"
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: 1 + index * 0.1 }}
                  >
                    • {tip}
                  </motion.li>
                ))}
              </ul>
            </motion.div>
          </div>
        </Card>
      </motion.div>
    </div>
  );
};

export default SpendingOverview; 