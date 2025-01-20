import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { AlertTriangle, Loader2 } from "lucide-react";
import { Prediction } from "./types";

const SpendingPredictions = () => {
  const [predictions, setPredictions] = useState<Prediction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    fetch('http://localhost:8080/api/analytics/predictions?accountId=1234567891')
      .then(res => {
        if (!res.ok) {
          throw new Error(`Failed to fetch predictions: ${res.status} ${res.statusText}`);
        }
        return res.json();
      })
      .then(data => {
        if (!data || !data.categoryPredictions) {
          throw new Error('Invalid data format received from server');
        }
        // Transform the data to match the Prediction interface
        const transformedPredictions: Prediction[] = data.categoryPredictions.map((p: any) => ({
          category: p.category,
          predictedAmount: p.predictedAmount || 0,
          confidence: (p.confidence || 0) * 100, // Convert to percentage
          trend: p.trend || [{amount: 0, date: 'Previous'}, {amount: 0, date: 'Predicted'}],
          warning: p.warning || (p.predictedAmount > p.trend?.[0]?.amount ? 
            `${((p.predictedAmount - p.trend[0].amount) / p.trend[0].amount * 100).toFixed(1)}% increase expected` : 
            undefined)
        }));
        setPredictions(transformedPredictions);
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
      {predictions.map((prediction, index) => (
        <motion.div
          key={prediction.category}
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.3, delay: index * 0.1 }}
        >
          <Card className="p-6">
            <div className="flex justify-between items-start mb-4">
              <div>
                <h3 className="text-lg font-medium">{prediction.category}</h3>
                <p className="text-sm text-gray-500">
                  Predicted spending: ${prediction.predictedAmount.toFixed(2)}
                </p>
              </div>
              {prediction.warning && (
                <div className="flex items-center text-yellow-500">
                  <AlertTriangle className="w-5 h-5 mr-2" />
                  <span className="text-sm">{prediction.warning}</span>
                </div>
              )}
            </div>

            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-gray-500">Last Month</p>
                  <p className="font-medium">${prediction.trend[0]?.amount.toFixed(2)}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-500">Next Month (Predicted)</p>
                  <p className="font-medium">${prediction.trend[prediction.trend.length - 1]?.amount.toFixed(2)}</p>
                </div>
              </div>

              <div>
                <div className="text-sm text-gray-500 mb-1">
                  Confidence: {prediction.confidence}%
                </div>
                <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                  <div 
                    className="h-full bg-blue-500 transition-all duration-500"
                    style={{ width: `${prediction.confidence}%` }}
                  />
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