import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { Progress } from "./ui/progress";
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts';
import { CategoryData } from "./types";

// Colors for different categories
const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8'];

const SpendingCategories = () => {
  const [categories, setCategories] = useState<CategoryData>({});

  useEffect(() => {
    fetch('http://localhost:8080/api/categories/1234567890')
      .then(res => res.json())
      .then(data => setCategories(data))
      .catch(err => console.error('Error fetching categories:', err));
  }, []);

  // Convert categories data for charts
  const pieData = Object.entries(categories).map(([name, value]) => ({
    name,
    value: Math.abs(value) // Use absolute value for positive numbers
  }));

  // Calculate total spending for percentages
  const totalSpending = Object.values(categories).reduce((sum, value) => sum + Math.abs(value), 0);

  return (
    <div className="grid gap-4 grid-cols-1 lg:grid-cols-2">
      {/* Pie Chart */}
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3 }}
      >
        <Card className="p-6">
          <h3 className="text-lg font-medium mb-4">Spending by Category</h3>
          <div className="h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={pieData}
                  dataKey="value"
                  nameKey="name"
                  cx="50%"
                  cy="50%"
                  outerRadius={100}
                  label={(entry) => entry.name}
                >
                  {pieData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </Card>
      </motion.div>

      {/* Category Breakdown */}
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.3, delay: 0.1 }}
      >
        <Card className="p-6">
          <h3 className="text-lg font-medium mb-4">Category Breakdown</h3>
          <div className="space-y-4">
            {Object.entries(categories).map(([category, amount], index) => {
              const percentage = (Math.abs(amount) / totalSpending) * 100;
              return (
                <div key={category} className="space-y-2">
                  <div className="flex justify-between text-sm">
                    <span>{category}</span>
                    <span>${Math.abs(amount).toFixed(2)} ({percentage.toFixed(1)}%)</span>
                  </div>
                  <Progress value={percentage} className="h-2" 
                    style={{ backgroundColor: COLORS[index % COLORS.length] + '40' }}>
                    <div className="h-full" style={{ backgroundColor: COLORS[index % COLORS.length] }} />
                  </Progress>
                </div>
              );
            })}
          </div>
        </Card>
      </motion.div>
    </div>
  );
};

export default SpendingCategories; 