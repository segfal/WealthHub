import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { AlertTriangle, Calendar, Loader2, TrendingUp } from "lucide-react";
import { getSpendingPredictions } from "../lib/api";
// Categories that are considered bills/rent and should be excluded
const EXCLUDED_CATEGORIES = new Set([
  'Rent', 
  'Income',
  'Utilities',
  'Insurance',
  'Phone Bill',
  'Internet',
  'Mortgage',
  'Water Bill',
  'Electric Bill',
  'Gas Bill'
]);

interface PredictedSpend {
  category: string;
  likelihood: number;
  predictedDate: string;
  warning: string;
  amount: number;
}

const SpendingPredictions = () => {
  const [predictions, setPredictions] = useState<PredictedSpend[]>([]);
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
    getSpendingPredictions(accountId)
      .then(data => {
        if (!Array.isArray(data)) {
          throw new Error('Invalid data format received from server');
        }
        // Filter out bills and rent from predictions
        const filteredPredictions = data.filter(pred => !EXCLUDED_CATEGORIES.has(pred.category));
        setPredictions(filteredPredictions);
        setError(null);
      })
      .catch(err => {
        console.error('Error fetching predictions:', err);
        setError(err.message);
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
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
    <div className="space-y-4">
      <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
        <motion.div
          className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-opacity duration-500 blur-xl"
          initial={false}
        />
        <div className="relative">
          <div className="flex items-center mb-4">
            <TrendingUp className="w-5 h-5 mr-2 text-[#00C805]" />
            <h3 className="text-lg font-medium text-white">Discretionary Spending Predictions</h3>
          </div>
          <p className="text-sm text-zinc-400 mb-4">
            Based on your spending patterns, here are likely upcoming expenses
            (excluding bills and recurring payments)
          </p>
        </div>
      </Card>

      {predictions.map((prediction, index) => (
        <motion.div
          key={prediction.category}
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.3, delay: index * 0.1 }}
        >
          <Card className="p-6 bg-[#1a1d21] backdrop-blur-xl border border-[#00C805]/20 relative group shadow-lg shadow-[#00C805]/5">
            <motion.div
              className="absolute inset-0 bg-gradient-radial from-[#00C805]/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 rounded-lg transition-opacity duration-500 blur-xl"
              initial={false}
            />
            <div className="relative">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h3 className="text-lg font-medium text-white">{prediction.category}</h3>
                  <div className="flex items-center text-sm text-zinc-400 mt-1">
                    <Calendar className="w-4 h-4 mr-1" />
                    <span>Expected around {new Date(prediction.predictedDate).toLocaleDateString()}</span>
                  </div>
                  {prediction.amount && (
                    <div className="text-sm font-medium mt-1 text-white">
                      Estimated amount: ${prediction.amount.toFixed(2)}
                    </div>
                  )}
                </div>
                {prediction.warning && (
                  <div className="flex items-center text-yellow-400">
                    <AlertTriangle className="w-5 h-5 mr-2" />
                    <span className="text-sm">{prediction.warning}</span>
                  </div>
                )}
              </div>

              <div className="space-y-4">
                <div>
                  <div className="text-sm text-zinc-400 mb-1">
                    Likelihood: {(prediction.likelihood * 100).toFixed(0)}%
                  </div>
                  <div className="h-2 bg-black/50 rounded-full overflow-hidden backdrop-blur-xl">
                    <motion.div 
                      className="h-full bg-gradient-to-r from-[#00C805] to-emerald-600"
                      initial={{ width: 0 }}
                      animate={{ width: `${prediction.likelihood * 100}%` }}
                      transition={{ duration: 1, ease: "easeOut" }}
                    />
                  </div>
                </div>
              </div>
            </div>
          </Card>
        </motion.div>
      ))}
    </div>
  );
};

export default SpendingPredictions; 