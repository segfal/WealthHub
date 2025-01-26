import { motion, AnimatePresence } from "framer-motion";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "./components/ui/tabs";
import { useState, useEffect } from "react";
import { Card } from "./components/ui/card";
import SpendingOverview from "@/components/SpendingOverview";
import SpendingCategories from "@/components/SpendingCategories";
import SpendingPredictions from "@/components/SpendingPredictions";
import SpendingPatterns from "@/components/SpendingPatterns";
import SpendingInsights from "@/components/SpendingInsights";
import BillsOverview from "@/components/BillsOverview";
import { DollarSign, Wallet, BellDot, Loader2, AlertTriangle, ChevronUp, ChevronDown, Activity, PieChart, TrendingUp, Receipt, LineChart, Sparkles } from "lucide-react";
import { getUser } from "./lib/api";

const containerVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.8,
      ease: [0.6, -0.05, 0.01, 0.99],
      staggerChildren: 0.1
    }
  }
};

const itemVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.6,
      ease: [0.6, -0.05, 0.01, 0.99]
    }
  }
};

const cardVariants = {
  hidden: { opacity: 0, scale: 0.95 },
  visible: {
    opacity: 1,
    scale: 1,
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

interface UserData {
  account_id: string;
  account_name: string;
  account_type: string;
  account_number: string;
  balance_current: number;
  balance_available: number;
  balance_currency: string;
  owner_name: string;
}

const Homepage = () => {
  const [user, setUser] = useState<UserData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState("overview");

  useEffect(() => {
    const accountId = import.meta.env.VITE_ACCOUNT_ID;
    if (!accountId) {
      setError("No account ID provided");
      setLoading(false);
      return;
    }

    setLoading(true);
    getUser(accountId)
      .then(data => {
        setUser(data);
        setError(null);
      })
      .catch(err => {
        console.error('Error fetching user:', err);
        setError(err.message);
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen bg-black text-white flex items-center justify-center">
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
          <Loader2 className="w-16 h-16 text-[#00C805]" />
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

  if (error || !user) {
    return (
      <div className="min-h-screen bg-black text-white flex items-center justify-center">
        <motion.div
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ duration: 0.5 }}
        >
          <Card className="p-8 bg-[#2F3136] border-red-500/20 backdrop-blur-xl">
            <div className="text-center text-red-500">
              <AlertTriangle className="w-16 h-16 mx-auto mb-4" />
              <p className="text-lg font-medium">{error || "Failed to load user data"}</p>
            </div>
          </Card>
        </motion.div>
      </div>
    );
  }

  const getBalanceChangeIndicator = () => {
    if (user?.balance_current !== undefined && user?.balance_available !== undefined) {
      const difference = user.balance_current - user.balance_available;
      const percentageChange = (difference / user.balance_available * 100).toFixed(1);
      const isPositive = Number(percentageChange) >= 0;
      
      return (
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          className={`flex items-center space-x-1 ${isPositive ? 'text-[#00C805]' : 'text-red-500'}`}
        >
          {isPositive ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
          <span className="text-sm font-medium">
            ${Math.abs(difference).toFixed(2)} ({percentageChange}%)
          </span>
        </motion.div>
      );
    }
    return null;
  };

  const tabIcons = {
    overview: <Activity className="w-4 h-4" />,
    categories: <PieChart className="w-4 h-4" />,
    insights: <TrendingUp className="w-4 h-4" />,
    bills: <Receipt className="w-4 h-4" />,
    predictions: <LineChart className="w-4 h-4" />,
    patterns: <Activity className="w-4 h-4" />
  };

  return (
    <div className="min-h-screen bg-black text-white relative overflow-hidden">
      <motion.div 
        className="absolute inset-0 bg-[url('/grid.svg')] opacity-5"
        initial={{ opacity: 0 }}
        animate={{ opacity: 0.05 }}
        transition={{ duration: 1 }}
      />
      <div className="max-w-7xl mx-auto p-6 relative">
        <motion.div
          variants={containerVariants}
          initial="hidden"
          animate="visible"
          className="space-y-8"
        >
          {/* Header Section */}
          <div className="flex justify-between items-start">
            <motion.div variants={itemVariants} className="relative">
              <motion.div
                className="absolute -left-4 -top-4 w-20 h-20 bg-[#00C805]/10 rounded-full blur-2xl"
                animate={{
                  scale: [1, 1.2, 1],
                  opacity: [0.3, 0.5, 0.3]
                }}
                transition={{
                  duration: 3,
                  repeat: Infinity,
                  ease: "easeInOut"
                }}
              />
              <h1 className="text-5xl font-bold bg-gradient-to-r from-[#00C805] to-emerald-500 bg-clip-text text-transparent relative">
                Welcome back, {user.owner_name}
                <motion.span
                  className="absolute -top-6 -right-6"
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
                  <Sparkles className="w-5 h-5 text-[#00C805]" />
                </motion.span>
              </h1>
              <p className="text-zinc-400 mt-2 text-lg">Your financial insights await</p>
            </motion.div>
            
            <motion.div variants={itemVariants} className="flex items-center space-x-8">
              <div className="text-right relative">
                <motion.div
                  className="absolute -right-4 -top-4 w-32 h-32 bg-[#00C805]/10 rounded-full blur-3xl"
                  animate={{
                    scale: [1, 1.2, 1],
                    opacity: [0.2, 0.4, 0.2]
                  }}
                  transition={{
                    duration: 4,
                    repeat: Infinity,
                    ease: "easeInOut"
                  }}
                />
                <motion.p 
                  className="text-5xl font-bold"
                  initial={{ opacity: 0, y: -20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.5, delay: 0.2 }}
                >
                  ${user?.balance_current?.toFixed(2) || '0.00'}
                </motion.p>
                {getBalanceChangeIndicator()}
              </div>
              <motion.button 
                className="p-3 rounded-full hover:bg-[#2F3136] transition-colors relative"
                whileHover={{ scale: 1.1 }}
                whileTap={{ scale: 0.95 }}
              >
                <BellDot className="w-6 h-6 text-[#00C805]" />
                <motion.span 
                  className="absolute top-0 right-0 w-2 h-2 bg-red-500 rounded-full"
                  animate={{
                    scale: [1, 1.2, 1],
                    opacity: [1, 0.8, 1]
                  }}
                  transition={{
                    duration: 2,
                    repeat: Infinity,
                    ease: "easeInOut"
                  }}
                />
              </motion.button>
            </motion.div>
          </div>

          {/* Quick Stats */}
          <motion.div 
            className="grid grid-cols-1 md:grid-cols-3 gap-6"
            variants={containerVariants}
          >
            <motion.div 
              className="bg-[#1a1d21] backdrop-blur-xl rounded-2xl p-8 border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5"
              variants={cardVariants}
              whileHover="hover"
            >
              <motion.div
                className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-2xl transition-opacity duration-500 blur-xl"
                initial={false}
              />
              <div className="flex items-center justify-between relative">
                <div>
                  <p className="text-zinc-400 text-sm font-medium">Available Balance</p>
                  <p className="text-3xl font-bold mt-2 text-white">
                    ${user?.balance_available?.toFixed(2) || '0.00'}
                  </p>
                </div>
                <div className="w-14 h-14 rounded-full bg-gradient-to-br from-[#00C805]/20 to-[#00C805]/5 flex items-center justify-center backdrop-blur-xl">
                  <Wallet className="w-7 h-7 text-[#00C805]" />
                </div>
              </div>
            </motion.div>

            <motion.div 
              className="bg-[#1a1d21] backdrop-blur-xl rounded-2xl p-8 border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5"
              variants={cardVariants}
              whileHover="hover"
            >
              <motion.div
                className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-2xl transition-opacity duration-500 blur-xl"
                initial={false}
              />
              <div className="flex items-center justify-between relative">
                <div>
                  <p className="text-zinc-400 text-sm font-medium">Current Balance</p>
                  <p className="text-3xl font-bold mt-2 text-white">
                    ${user?.balance_current?.toFixed(2) || '0.00'}
                  </p>
                </div>
                <div className="w-14 h-14 rounded-full bg-gradient-to-br from-[#00C805]/20 to-[#00C805]/5 flex items-center justify-center backdrop-blur-xl">
                  <DollarSign className="w-7 h-7 text-[#00C805]" />
                </div>
              </div>
            </motion.div>

            <motion.div 
              className="bg-[#1a1d21] backdrop-blur-xl rounded-2xl p-8 border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5"
              variants={cardVariants}
              whileHover="hover"
            >
              <motion.div
                className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-2xl transition-opacity duration-500 blur-xl"
                initial={false}
              />
              <div className="flex items-center justify-between relative">
                <div>
                  <p className="text-zinc-400 text-sm font-medium">Account Type</p>
                  <p className="text-3xl font-bold mt-2 text-white">{user?.account_type || 'N/A'}</p>
                </div>
                <div className="bg-gradient-to-br from-[#00C805]/20 to-[#00C805]/5 text-[#00C805] px-6 py-3 rounded-full text-sm font-medium backdrop-blur-xl">
                  {user?.account_type || 'N/A'}
                </div>
              </div>
            </motion.div>
          </motion.div>

          {/* Main Dashboard Content */}
          <motion.div variants={itemVariants}>
            <Tabs 
              defaultValue="overview" 
              className="space-y-6"
              value={activeTab}
              onValueChange={setActiveTab}
            >
              <TabsList className="w-full bg-[#2F3136] backdrop-blur-xl p-1.5 rounded-2xl border border-zinc-800/50">
                {Object.entries(tabIcons).map(([value, icon]) => (
                  <TabsTrigger
                    key={value}
                    value={value}
                    className="flex-1 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#00C805] data-[state=active]:to-emerald-600 data-[state=active]:text-black text-zinc-400 capitalize py-3"
                  >
                    <motion.div
                      className="flex items-center space-x-2"
                      initial={false}
                      animate={{ 
                        scale: activeTab === value ? 1.1 : 1,
                        y: activeTab === value ? -2 : 0
                      }}
                      transition={{ duration: 0.2 }}
                    >
                      {icon}
                      <span>{value}</span>
                    </motion.div>
                  </TabsTrigger>
                ))}
              </TabsList>

              <AnimatePresence mode="wait">
                <motion.div
                  key={activeTab}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -20 }}
                  transition={{ duration: 0.3 }}
                  className="bg-[#2F3136] backdrop-blur-xl rounded-2xl border border-zinc-800/50 relative overflow-hidden"
                >
                  <motion.div
                    className="absolute inset-0 bg-gradient-to-b from-[#00C805]/5 to-transparent opacity-0 group-hover:opacity-100"
                    initial={false}
                  />
                  <div className="p-8 relative">
                    <TabsContent value="overview">
                      <SpendingOverview />
                    </TabsContent>

                    <TabsContent value="categories">
                      <SpendingCategories />
                    </TabsContent>

                    <TabsContent value="insights">
                      <SpendingInsights />
                    </TabsContent>

                    <TabsContent value="bills">
                      <BillsOverview />
                    </TabsContent>

                    <TabsContent value="predictions">
                      <SpendingPredictions />
                    </TabsContent>

                    <TabsContent value="patterns">
                      <SpendingPatterns />
                    </TabsContent>
                  </div>
                </motion.div>
              </AnimatePresence>
            </Tabs>
          </motion.div>
        </motion.div>
      </div>
    </div>
  );
};

export default Homepage;
