import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Card } from "./ui/card";
import { AlertTriangle, Receipt, Calendar, Loader2, DollarSign, PieChart } from "lucide-react";
import { getBillsIncomeAnalysis } from "../lib/api";

interface BillAnalysis {
  merchant: string;
  amount: number;
  percentage: number;
  date: string;
}

interface AnalysisResponse {
  monthlyIncome: number;
  totalBills: number;
  billsToIncomeRatio: number;
  remainingIncome: number;
  remainingIncomePercentage: number;
  bills: BillAnalysis[];
  year: number;
  month: number;
  monthName: string;
}

const BillsOverview = () => {
  const [analysis, setAnalysis] = useState<AnalysisResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const accountId = import.meta.env.VITE_ACCOUNT_ID;
    if (!accountId) {
      setError("No account ID provided");
      setLoading(false);
      return;
    }
    // TODO :
    

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error || !analysis) {
    return (
      <Card className="p-6">
        <div className="text-center text-destructive">
          <AlertTriangle className="w-8 h-8 mx-auto mb-2" />
          <p>{error || "Failed to load analysis"}</p>
        </div>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* Monthly Summary */}
      <motion.div
        initial={{ scale: 0.95, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3 }}
        className="grid grid-cols-1 md:grid-cols-2 gap-4"
      >
        <Card className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-medium flex items-center">
              <DollarSign className="w-5 h-5 mr-2" />
              Monthly Income
            </h3>
            <span className="text-2xl font-bold text-green-600">
              ${analysis.monthlyIncome.toFixed(2)}
            </span>
          </div>
          <div className="text-sm text-gray-500">
            {analysis.monthName} {analysis.year}
          </div>
        </Card>

        <Card className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-medium flex items-center">
              <Receipt className="w-5 h-5 mr-2" />
              Total Bills
            </h3>
            <span className="text-2xl font-bold text-red-600">
              ${analysis.totalBills.toFixed(2)}
            </span>
          </div>
          <div className="text-sm text-gray-500">
            {analysis.bills.length} recurring payments
          </div>
        </Card>
      </motion.div>

      {/* Income Distribution */}
      <motion.div
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.4, delay: 0.2 }}
      >
        <Card className="p-6">
          <div className="flex items-center mb-4">
            <PieChart className="w-5 h-5 mr-2" />
            <h3 className="text-lg font-medium">Income Distribution</h3>
          </div>
          <div className="relative h-4 bg-gray-200 rounded-full overflow-hidden mb-4">
            <motion.div
              className="absolute left-0 top-0 h-full bg-red-500"
              initial={{ width: 0 }}
              animate={{ width: `${analysis.billsToIncomeRatio}%` }}
              transition={{ duration: 0.8, delay: 0.4 }}
            />
          </div>
          <div className="flex justify-between text-sm">
            <div>
              <span className="font-medium">Bills: </span>
              <span className="text-red-600">{analysis.billsToIncomeRatio.toFixed(1)}%</span>
            </div>
            <div>
              <span className="font-medium">Remaining: </span>
              <span className="text-green-600">{analysis.remainingIncomePercentage.toFixed(1)}%</span>
            </div>
          </div>
        </Card>
      </motion.div>

      {/* Bills Breakdown */}
      <div className="space-y-4">
        <h3 className="text-lg font-medium ml-2">Bills Breakdown</h3>
        <AnimatePresence>
          {analysis.bills.map((bill, index) => (
            <motion.div
              key={bill.merchant}
              initial={{ x: -20, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              exit={{ x: 20, opacity: 0 }}
              transition={{ duration: 0.3, delay: index * 0.1 }}
            >
              <Card className="p-4">
                <div className="flex justify-between items-center">
                  <div>
                    <h4 className="font-medium">{bill.merchant}</h4>
                    <div className="flex items-center text-sm text-gray-500 mt-1">
                      <Calendar className="w-4 h-4 mr-1" />
                      <span>Last paid: {bill.date}</span>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-lg font-medium">${bill.amount.toFixed(2)}</div>
                    <div className="text-sm text-gray-500">
                      {bill.percentage.toFixed(1)}% of income
                    </div>
                  </div>
                </div>
                <div className="mt-2 h-1 bg-gray-200 rounded-full overflow-hidden">
                  <motion.div
                    className="h-full bg-blue-500"
                    initial={{ width: 0 }}
                    animate={{ width: `${bill.percentage}%` }}
                    transition={{ duration: 0.5, delay: 0.6 + index * 0.1 }}
                  />
                </div>
              </Card>
            </motion.div>
          ))}
        </AnimatePresence>
      </div>
    </div>
  );
};

export default BillsOverview; 