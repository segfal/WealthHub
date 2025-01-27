import axios from 'axios';

const API_URL = import.meta.env.VITE_URL || 'http://localhost:8080';

export const apiService = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add response interceptor for better error handling
apiService.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', {
      status: error.response?.status,
      data: error.response?.data,
      config: error.config,
    });
    return Promise.reject(error);
  }
);

export const getUser = async (accountId: string) => {
    try {
        const response = await apiService.get(`/api/user/${accountId}`);
        console.log('User API Response:', response.data);
        return response.data;
    } catch (error) {
        console.error('Error fetching user:', error);
        throw error;
    }
};

// Analytics Endpoints
export const getSpendingOverview = async (accountId: string) => {
    const response = await apiService.get(`/api/analytics/${accountId}`);
    console.log('Spending Overview API Response:', response.data);
    return response.data;
};

export const getSpendingPredictions = async (accountId: string) => {
    const response = await apiService.get(`/api/predictions/${accountId}`);
    console.log('Predictions API Response:', response.data);
    return response.data;
};

export const getSpendingPatterns = async (accountId: string) => {
    const response = await apiService.get(`/api/patterns/${accountId}`);
    console.log('Patterns API Response:', response.data);
    return response.data;
};

export const getSpendingInsights = async (accountId: string) => {
    const response = await apiService.get(`/api/insights/${accountId}`);
    console.log('Insights API Response:', response.data);
    return response.data;
};

// Bills Endpoints
export const getBillsOverview = async (accountId: string) => {
    const response = await apiService.get(`/api/bills/${accountId}`);
    console.log('Bills Overview API Response:', response.data);
    return response.data;
};

export const getBillsIncomeAnalysis = async (accountId: string, year?: number, month?: number) => {
    const params = new URLSearchParams();
    if (year) params.append('year', year.toString());
    if (month) params.append('month', month.toString());
    
    const response = await apiService.get(`/api/bills/${accountId}/analysis?${params}`);
    console.log('Bills-Income Analysis Response:', response.data);
    return response.data;
};

// Categories Endpoints
export const getSpendingCategories = async (accountId: string) => {
    const response = await apiService.get(`/api/categories/${accountId}`);
    console.log('Categories API Response:', response.data);
    return response.data;
};

export const getCategoryTotals = async (accountId: string) => {
    const response = await apiService.get(`/api/categories/${accountId}/totals`);
    console.log('Category Totals API Response:', response.data);
    return response.data;
};

// Income Endpoints
export const getIncome = async (accountId: string) => {
    const response = await apiService.get(`/api/income/${accountId}`);
    console.log('Income API Response:', response.data);
    return response.data;
};

export const getMonthlyIncome = async (accountId: string, year: number, month: number) => {
    const response = await apiService.get(`/api/income/${accountId}/monthly?year=${year}&month=${month}`);
    console.log('Monthly Income API Response:', response.data);
    return response.data;
};

