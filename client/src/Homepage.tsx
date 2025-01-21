import { motion } from "framer-motion";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "./components/ui/tabs";
import { Card } from "./components/ui/card";
import SpendingOverview from "@/components/SpendingOverview";
import SpendingCategories from "@/components/SpendingCategories";
import SpendingPredictions from "@/components/SpendingPredictions";
import SpendingPatterns from "@/components/SpendingPatterns";
import SpendingInsights from "@/components/SpendingInsights";
import BillsOverview from "@/components/BillsOverview";
import { DollarSign, Wallet, BellDot } from "lucide-react";
import { dummyTransactions } from "./data";

const user = dummyTransactions[0];

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1
    }
  }
};



const Homepage = () => {
  return (
    <div className="min-h-screen bg-black text-white">
      {/* Header Section */}
      <div className="max-w-7xl mx-auto p-6">
        <div className="flex justify-between items-center mb-8">
          <motion.div
            initial={{ y: -20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.5 }}
          >
            <h1 className="text-3xl font-bold">Welcome back, {user.account.owner_name}</h1>
            <p className="text-gray-400 text-sm">Track your financial journey</p>
          </motion.div>
          
          <motion.div 
            initial={{ y: -20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.5 }}
            className="flex items-center space-x-6"
          >
            <div className="text-right">
              <p className="text-3xl font-bold">${user.account.balance.current}</p>
              <p className="text-[#00C805] text-sm font-medium">+ ${user.account.balance.current - user.account.balance.available} (8.4%)</p>
            </div>
            <button className="p-2 rounded-full hover:bg-zinc-900 transition-colors">
              <BellDot className="w-6 h-6 text-gray-400" />
            </button>
          </motion.div>
        </div>

        {/* Quick Stats */}
        <motion.div 
          className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8"
          variants={containerVariants}
          initial="hidden"
          animate="visible"
        >
          <motion.div 
            className="bg-zinc-900 rounded-2xl p-6 border border-zinc-800"
            variants={containerVariants}
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-400 text-sm">Available Balance</p>
                <p className="text-xl font-bold mt-1">${user.account.balance.available}</p>
              </div>
              <Wallet className="w-8 h-8 text-[#00C805]" />
            </div>
          </motion.div>

          <motion.div 
            className="bg-zinc-900 rounded-2xl p-6 border border-zinc-800"
            variants={containerVariants}
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-400 text-sm">Monthly Spending</p>
                <p className="text-xl font-bold mt-1">${user.account.balance.available}</p>
              </div>
              <DollarSign className="w-8 h-8 text-[#00C805]" />
            </div>
          </motion.div>

          <motion.div 
            className="bg-zinc-900 rounded-2xl p-6 border border-zinc-800"
            variants={containerVariants}
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-400 text-sm">Account Type</p>
                <p className="text-xl font-bold mt-1">{user.account.account_type}</p>
              </div>
              <div className="bg-[#00C805]/10 text-[#00C805] px-3 py-1 rounded-full text-sm font-medium">
                Premium
              </div>
            </div>
          </motion.div>
        </motion.div>

        {/* Main Dashboard Content */}
        <Tabs defaultValue="overview" className="space-y-4">
          <TabsList className="w-full bg-zinc-900 p-1 rounded-xl border border-zinc-800">
            <TabsTrigger 
              value="overview" 
              className="flex-1 data-[state=active]:bg-zinc-800 data-[state=active]:text-white text-gray-400"
            >
              Overview
            </TabsTrigger>
            <TabsTrigger 
              value="categories" 
              className="flex-1 data-[state=active]:bg-zinc-800 data-[state=active]:text-white text-gray-400"
            >
              Categories
            </TabsTrigger>
            <TabsTrigger 
              value="insights" 
              className="flex-1 data-[state=active]:bg-zinc-800 data-[state=active]:text-white text-gray-400"
            >
              Insights
            </TabsTrigger>
            <TabsTrigger 
              value="bills" 
              className="flex-1 data-[state=active]:bg-zinc-800 data-[state=active]:text-white text-gray-400"
            >
              Bills
            </TabsTrigger>
            <TabsTrigger 
              value="predictions" 
              className="flex-1 data-[state=active]:bg-zinc-800 data-[state=active]:text-white text-gray-400"
            >
              Predictions
            </TabsTrigger>
            <TabsTrigger 
              value="patterns" 
              className="flex-1 data-[state=active]:bg-zinc-800 data-[state=active]:text-white text-gray-400"
            >
              Patterns
            </TabsTrigger>
          </TabsList>

          <motion.div
            variants={containerVariants}
            initial="hidden"
            animate="visible"
            className="bg-zinc-900 rounded-2xl border border-zinc-800"
          >
            <TabsContent value="overview" className="p-6">
              <SpendingOverview />
            </TabsContent>

            <TabsContent value="categories" className="p-6">
              <SpendingCategories />
            </TabsContent>

            <TabsContent value="insights" className="p-6">
              <SpendingInsights />
            </TabsContent>

            <TabsContent value="bills" className="p-6">
              <BillsOverview />
            </TabsContent>

            <TabsContent value="predictions" className="p-6">
              <SpendingPredictions />
            </TabsContent>

            <TabsContent value="patterns" className="p-6">
              <SpendingPatterns />
            </TabsContent>
          </motion.div>
        </Tabs>
      </div>
    </div>
  );
};

export default Homepage;
