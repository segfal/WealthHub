import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Card } from "./ui/card";
import { AlertTriangle, Receipt, Calendar, Loader2, DollarSign, PieChart } from "lucide-react";
import { getBillsOverview, getMonthlyIncome } from "../lib/api";

interface BillTransaction {
  transaction_id: string;
  account_id: string;
  date: string;
  amount: number;
  category: string;
  merchant: string;
  location: string;
}

interface ProcessedBill {
  merchant: string;
  amount: number;
  percentage: number;
  date: string;
  nextDueDate?: string;
  isRecurring: boolean;
}

interface AnalysisData {
  monthlyIncome: number;
  totalBills: number;
  billsToIncomeRatio: number;
  remainingIncome: number;
  remainingIncomePercentage: number;
  bills: ProcessedBill[];
  year: number;
  month: number;
  monthName: string;
}

const BillsOverview = () => {
  const [analysis, setAnalysis] = useState<AnalysisData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const accountId = import.meta.env.VITE_ACCOUNT_ID;
    if (!accountId) {
      setError("No account ID provided");
      setLoading(false);
      return;
    }

    const currentDate = new Date();
    const year = currentDate.getFullYear();
    const month = currentDate.getMonth() + 1;

    setLoading(true);
    
    // Fetch both bills and income data
    Promise.all([
      getBillsOverview(accountId),
      getMonthlyIncome(accountId, year, month)
    ])
      .then(([billsData, incomeData]) => {
        // Process bills data
        const monthlyIncome = Array.isArray(incomeData) ? 
          incomeData.reduce((sum, t) => sum + Math.abs(t.amount), 0) : 0;

        const totalBills = billsData.reduce((sum, t) => sum + Math.abs(t.amount), 0);
        const billsToIncomeRatio = monthlyIncome > 0 ? (totalBills / monthlyIncome) * 100 : 0;
        
        // Group bills by merchant and calculate totals
        const billsByMerchant = new Map<string, ProcessedBill>();
        billsData.forEach(bill => {
          const existing = billsByMerchant.get(bill.merchant);
          if (existing) {
            existing.amount += Math.abs(bill.amount);
            if (new Date(bill.date) > new Date(existing.date)) {
              existing.date = bill.date;
            }
          } else {
            billsByMerchant.set(bill.merchant, {
              merchant: bill.merchant,
              amount: Math.abs(bill.amount),
              percentage: 0, // Will calculate after
              date: bill.date,
              isRecurring: true, // Assuming all bills are recurring for now
              nextDueDate: new Date(bill.date).toISOString() // You might want to calculate this differently
            });
          }
        });

        // Calculate percentages
        const processedBills = Array.from(billsByMerchant.values()).map(bill => ({
          ...bill,
          percentage: monthlyIncome > 0 ? (bill.amount / monthlyIncome) * 100 : 0
        }));

        const analysisData: AnalysisData = {
          monthlyIncome,
          totalBills,
          billsToIncomeRatio,
          remainingIncome: monthlyIncome - totalBills,
          remainingIncomePercentage: 100 - billsToIncomeRatio,
          bills: processedBills,
          year,
          month,
          monthName: new Date(year, month - 1).toLocaleString('default', { month: 'long' })
        };

        setAnalysis(analysisData);
        setError(null);
      })
      .catch(err => {
        console.error('Error fetching bills analysis:', err);
        setError(err.message);
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="w-8 h-8 animate-spin text-[#00C805]" />
      </div>
    );
  }

  if (error || !analysis) {
    return (
      <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
        <div className="text-center text-red-400">
          <AlertTriangle className="w-8 h-8 mx-auto mb-2" />
          <p>{error || "Failed to load analysis"}</p>
        </div>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* Monthly Summary */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <motion.div
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.5 }}
        >
          <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
            <motion.div
              className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-opacity duration-500 blur-xl"
              initial={false}
            />
            <div className="relative">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium flex items-center text-white">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-br from-[#00C805]/20 to-[#00C805]/5 flex items-center justify-center backdrop-blur-xl mr-3">
                    <DollarSign className="w-5 h-5 text-[#00C805]" />
                  </div>
                  Monthly Income
                </h3>
                <span className="text-2xl font-bold text-[#00C805]">
                  ${analysis.monthlyIncome.toFixed(2)}
                </span>
              </div>
              <div className="text-sm text-zinc-400">
                {analysis.monthName} {analysis.year}
              </div>
            </div>
          </Card>
        </motion.div>

        <motion.div
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.5, delay: 0.1 }}
        >
          <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
            <motion.div
              className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-opacity duration-500 blur-xl"
              initial={false}
            />
            <div className="relative">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium flex items-center text-white">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-br from-[#00C805]/20 to-[#00C805]/5 flex items-center justify-center backdrop-blur-xl mr-3">
                    <Receipt className="w-5 h-5 text-[#00C805]" />
                  </div>
                  Total Bills
                </h3>
                <span className="text-2xl font-bold text-red-400">
                  ${analysis.totalBills.toFixed(2)}
                </span>
              </div>
              <div className="text-sm text-zinc-400">
                {analysis.bills.length} recurring payments
              </div>
            </div>
          </Card>
        </motion.div>
      </div>

      {/* Income Distribution */}
      <motion.div
        initial={{ y: -20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.5, delay: 0.2 }}
      >
        <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
          <motion.div
            className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-opacity duration-500 blur-xl"
            initial={false}
          />
          <div className="relative">
            <div className="flex items-center mb-4">
              <div className="w-10 h-10 rounded-full bg-gradient-to-br from-[#00C805]/20 to-[#00C805]/5 flex items-center justify-center backdrop-blur-xl mr-3">
                <PieChart className="w-5 h-5 text-[#00C805]" />
              </div>
              <h3 className="text-lg font-medium text-white">Income Distribution</h3>
            </div>
            <div className="relative h-4 bg-black/50 rounded-full overflow-hidden backdrop-blur-xl mb-4">
              <motion.div
                className="absolute left-0 top-0 h-full bg-gradient-to-r from-red-500 to-red-600"
                initial={{ width: 0 }}
                animate={{ width: `${analysis.billsToIncomeRatio}%` }}
                transition={{ duration: 0.8, delay: 0.4 }}
              />
            </div>
            <div className="flex justify-between text-sm">
              <div>
                <span className="text-zinc-400">Bills: </span>
                <span className="text-red-400 font-medium">{analysis.billsToIncomeRatio.toFixed(1)}%</span>
              </div>
              <div>
                <span className="text-zinc-400">Remaining: </span>
                <span className="text-[#00C805] font-medium">{analysis.remainingIncomePercentage.toFixed(1)}%</span>
              </div>
            </div>
          </div>
        </Card>
      </motion.div>

      {/* Bills Breakdown */}
      <div className="space-y-4">
        <h3 className="text-lg font-medium text-white ml-2">Bills Breakdown</h3>
        <AnimatePresence>
          {analysis.bills.map((bill, index) => (
            <motion.div
              key={bill.merchant}
              initial={{ y: -20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              exit={{ y: 20, opacity: 0 }}
              transition={{ duration: 0.3, delay: index * 0.1 }}
            >
              <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
                <motion.div
                  className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-opacity duration-500 blur-xl"
                  initial={false}
                />
                <div className="relative">
                  <div className="flex justify-between items-center">
                    <div>
                      <h4 className="font-medium text-white">{bill.merchant}</h4>
                      <div className="flex items-center text-sm text-zinc-400 mt-1">
                        <Calendar className="w-4 h-4 mr-1" />
                        <span>Last paid: {new Date(bill.date).toLocaleDateString()}</span>
                        {bill.nextDueDate && (
                          <span className="ml-2">Next due: {new Date(bill.nextDueDate).toLocaleDateString()}</span>
                        )}
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="text-lg font-medium text-white">${bill.amount.toFixed(2)}</div>
                      <div className="text-sm text-zinc-400">
                        {bill.percentage.toFixed(1)}% of income
                      </div>
                    </div>
                  </div>
                  <div className="mt-2 h-2 bg-black/50 rounded-full overflow-hidden backdrop-blur-xl">
                    <motion.div
                      className="h-full bg-gradient-to-r from-[#00C805] to-emerald-600"
                      initial={{ width: 0 }}
                      animate={{ width: `${bill.percentage}%` }}
                      transition={{ duration: 0.5, delay: 0.6 + index * 0.1 }}
                    />
                  </div>
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