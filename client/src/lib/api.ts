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
        console.log('Fetching user from:', `${API_URL}/api/user/${accountId}`);
        const response = await apiService.get(`/api/user/${accountId}`);
        console.log('User API Response:', response.data);
        return response.data;
    } catch (error) {
        console.error('Error fetching user:', error);
        throw error;
    }
};

export const getSpendingOverview = async (accountId: string) => {
    const response = await apiService.get(`/api/analytics/${accountId}?timeRange=1%20month`);
    console.log('Spending Overview API Response:', response.data);
    console.log('URL:', `${API_URL}/api/analytics/${accountId}?timeRange=1%20month`);
    return response.data;
};

export const getBillsOverview = async (accountId: string) => {
    const response = await apiService.get(`/api/bills/${accountId}`);
    console.log('Bills Overview API Response:', response.data);
    console.log('URL:', `${API_URL}/api/bills/${accountId}`);
    return response.data;
};

export const getInsights = async (accountId: string) => {
    const response = await apiService.get(`/api/analytics/${accountId}?timeRange=1%20month`);
    return response.data;
};

export const getTransactions = async (accountId: string) => {
    const response = await apiService.get(`/api/transactions/${accountId}`);
    return response.data;
};

export const getAccountDetails = async () => {
    const response = await axios.get(`${import.meta.env.VITE_URL}/account?accountId=${import.meta.env.VITE_ACCOUNT_ID}`);
    return response.data;
};

export const getSpendingCategories = async (accountId: string) => {
    const response = await apiService.get(`/api/categories/${accountId}`);
    return response.data;
};

export const getSpendingPredictions = async (accountId: string) => {
    const response = await apiService.get(`/api/predictions/${accountId}`);
    return response.data;
};

export const getSpendingPatterns = async (accountId: string) => {
    const response = await apiService.get(`/api/patterns/${accountId}`);
    return response.data;
};

export const getSpendingInsights = async (accountId: string) => {
    const response = await apiService.get(`/api/insights/${accountId}`);
    console.log('Insights API Response:', response.data);
    console.log('URL:', `${API_URL}/api/insights/${accountId}`);
    return response.data;
};

export const getBillsIncomeAnalysis = async (accountId: string, year?: number, month?: number) => {
    const params = new URLSearchParams();
    if (year) params.append('year', year.toString());
    if (month) params.append('month', month.toString());
    
    const response = await apiService.get(`/api/analysis/bills-income/${accountId}?${params}`);
    console.log('Bills-Income Analysis Response:', response.data);
    return response.data;
};

