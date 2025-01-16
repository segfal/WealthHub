
export interface CategoryData {
    [key: string]: number;
}

export interface Analytics {
    totalSpent: number;
    averageTransaction: number;
    transactionCount: number;
    spendingTrend: { date: string; amount: number }[];
  }


export interface Prediction {
    category: string;
    predictedAmount: number;
    confidence: number;
    trend: { date: string; amount: number }[];
    warning?: string;
  }



export interface Pattern {
    timeOfDay: {
      morning: number;
      afternoon: number;
      evening: number;
      night: number;
    };
    dayOfWeek: {
      [key: string]: number;
    };
    recurringTransactions: {
      merchant: string;
      amount: number;
      frequency: string;
    }[];
  }