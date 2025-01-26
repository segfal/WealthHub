import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { useState } from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface Bill {
  name: string;
  category: string;
  amount: number;
  paycheckPercentage: number;
  lastPaidDate: string;
  merchant: string;
  location: string;
  monthlyIncome: number;
}

interface BillsData {
  bills: Bill[];
  totalBillsAmount: number;
  totalIncome: number;
  totalBillsPercentage: number;
}

interface BillPaymentGridProps {
  initialData: BillsData;
}

// Company logos mapping (you can expand this)
const companyLogos: Record<string, string> = {
  "Netflix": "/logos/netflix.png",
  "Spotify": "/logos/spotify.png",
  "ChatGPT": "/logos/chatgpt.png",
  "McDonald's": "/logos/mcdonalds.png",
  "Gas Station": "/logos/gas.png",
  "Car Insurance": "/logos/car-insurance.png",
  // Add more company logos as needed
};

const months = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December"
];

const BillPaymentGrid = ({ initialData }: BillPaymentGridProps) => {
  const [selectedMonth, setSelectedMonth] = useState(new Date().getMonth());
  
  // Group bills by month with monthly income
  const billsByMonth = initialData.bills.reduce((acc, bill) => {
    const date = new Date(bill.lastPaidDate);
    const month = date.getMonth();
    if (!acc[month]) {
      acc[month] = {
        bills: [],
        totalAmount: 0,
        income: bill.monthlyIncome || 0
      };
    }
    acc[month].bills.push(bill);
    acc[month].totalAmount += Math.abs(bill.amount);
    // Update monthly income if it's higher (in case of multiple transactions)
    if (bill.monthlyIncome > acc[month].income) {
      acc[month].income = bill.monthlyIncome;
    }
    return acc;
  }, {} as Record<number, { bills: Bill[], totalAmount: number, income: number }>);

  // Calculate monthly percentage
  const monthData = billsByMonth[selectedMonth] || { bills: [], totalAmount: 0, income: 0 };
  const monthlyPercentage = monthData.income > 0 ? (monthData.totalAmount / monthData.income) * 100 : 0;

  return (
    <div className="space-y-6">
      {/* Month Selection */}
      <div className="flex items-center justify-between">
        <button
          onClick={() => setSelectedMonth((prev) => (prev === 0 ? 11 : prev - 1))}
          className="p-2 hover:bg-gray-100 rounded-full"
        >
          <ChevronLeft className="w-6 h-6" />
        </button>
        <h2 className="text-2xl font-bold">{months[selectedMonth]} Bills</h2>
        <button
          onClick={() => setSelectedMonth((prev) => (prev === 11 ? 0 : prev + 1))}
          className="p-2 hover:bg-gray-100 rounded-full"
        >
          <ChevronRight className="w-6 h-6" />
        </button>
      </div>

      {/* Summary Card */}
      <motion.div
        key={selectedMonth}
        initial={{ scale: 0.95, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3 }}
      >
        <Card className="p-6 bg-gradient-to-br from-primary/10 to-primary/5">
          <h2 className="text-2xl font-bold mb-2">Monthly Summary</h2>
          <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
            <div>
              <p className="text-sm text-gray-500">Total Bills</p>
              <p className="text-xl font-bold">${monthData.totalAmount.toFixed(2)}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Monthly Income</p>
              <p className="text-xl font-bold">${monthData.income.toFixed(2)}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Bills Percentage</p>
              <p className="text-xl font-bold">{monthlyPercentage.toFixed(1)}%</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Number of Bills</p>
              <p className="text-xl font-bold">{monthData.bills.length}</p>
            </div>
          </div>
        </Card>
      </motion.div>

      {/* Bills Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {monthData.bills.map((bill, index) => (
          <motion.div
            key={`${bill.name}-${bill.lastPaidDate}`}
            initial={{ scale: 0, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            transition={{
              type: "spring",
              stiffness: 260,
              damping: 20,
              delay: index * 0.1,
            }}
          >
            <Card className="p-4 hover:shadow-lg transition-shadow duration-300">
              <div className="flex items-start space-x-4">
                <div className="w-12 h-12 flex-shrink-0">
                  {companyLogos[bill.merchant] ? (
                    <img
                      src={companyLogos[bill.merchant]}
                      alt={bill.merchant}
                      className="w-full h-full object-contain"
                    />
                  ) : (
                    <div className="w-full h-full bg-primary/10 rounded-full flex items-center justify-center">
                      <span className="text-lg font-bold text-primary">
                        {bill.name[0]}
                      </span>
                    </div>
                  )}
                </div>
                <div className="flex-1">
                  <h3 className="font-bold text-lg">{bill.name}</h3>
                  <p className="text-sm text-gray-500">{bill.category}</p>
                  <div className="mt-2 flex justify-between items-end">
                    <div>
                      <p className="text-2xl font-bold">${Math.abs(bill.amount).toFixed(2)}</p>
                      <p className="text-sm text-gray-500">
                        {((Math.abs(bill.amount) / monthData.income) * 100).toFixed(1)}% of income
                      </p>
                    </div>
                    <div className="text-right text-sm text-gray-500">
                      Paid on:<br />
                      {new Date(bill.lastPaidDate).toLocaleDateString()}
                    </div>
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

export default BillPaymentGrid; 
