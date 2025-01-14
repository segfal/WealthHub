import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Card } from "./ui/card";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { AlertTriangle } from "lucide-react";
import axios from 'axios';
import { dummyTransactions } from "../data";
import { Prediction } from "./types";


const user = dummyTransactions[0];



const SpendingPredictions = () => {
  const [predictions, setPredictions] = useState<Prediction[]>([]);

  useEffect(() => {
    axios.get('http://localhost:8080/api/predictions/1234567890')
      .then(res => setPredictions(res.data))
      .catch(err => console.error('Error fetching predictions:', err));
  }, []);

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
                  Predicted spending: ${prediction.likelihood}
                </p>
              </div>
              {prediction.warning && (
                <div className="flex items-center text-yellow-500">
                  <AlertTriangle className="w-5 h-5 mr-2" />
                  <span className="text-sm">{prediction.warning}</span>
                </div>
              )}
            </div>

            <div className="h-[200px]">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={prediction.trend}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" />
                  <YAxis />
                  <Tooltip />
                  <Line 
                    type="monotone" 
                    dataKey="amount" 
                    stroke="#8884d8" 
                    strokeWidth={2}
                    dot={false}
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>

            <div className="mt-4">
              <div className="text-sm text-gray-500">
                Confidence: {prediction.confidence}%
              </div>
            </div>
          </Card>
        </motion.div>
      ))}
    </div>
  );
};

export default SpendingPredictions; 