import { useState, useMemo, useEffect } from 'react';

interface Transaction {
    transactionId: string;
    category: string;
    location: string;
    company: string;
    amount: number;
    date: string;
    description: string;
    status: string;
}

const HomePage = () => {
    const currentDate = new Date(); 
    const formattedDate = currentDate.toLocaleDateString();
    const [date, setDate] = useState(formattedDate);
    const [dummyTransactions, setDummyTransactions] = useState<Transaction[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchTransactions = async () => {
            try {
                setIsLoading(true);
                // Simulating an API call with local data
                const response = await import('./data.js');
                setDummyTransactions(response.dummyTransactions);
                setError(null);
            } catch (err) {
                setError('Failed to load transactions');
                console.error('Error loading transactions:', err);
            } finally {
                setIsLoading(false);
            }
        };

        fetchTransactions();
    }, []); // Empty dependency array means this runs once on mount

    // Handler functions with correct TypeScript syntax
    const handleDate = (transaction: Transaction): string => {
        return transaction.date;
    }

    const handleMoney = (transaction: Transaction): number => {
        return transaction.amount;
    }

    const handleLocation = (transaction: Transaction): string => {
        return transaction.location;
    }

    const handleCompany = (transaction: Transaction): string => {
        return transaction.company;
    }

    const handleCategory = (transaction: Transaction): string => {
        return transaction.category;
    }

    const handleDescription = (transaction: Transaction): string => {
        return transaction.description;
    }

    const handleStatus = (transaction: Transaction): string => {
        return transaction.status;
    }

    //useMemo uses previous data in localStorage. useMemo stores/retrieves data from the browser
    // Analytics calculations with null checks
    const analytics = useMemo(() => {
        if (!dummyTransactions.length) {
            return {
                totalSpent: 0,
                spendingByCategory: {},
                mostFrequentLocation: 'No data',
                highestExpense: null
            };
        }

        return {
            totalSpent: dummyTransactions.reduce((sum, transaction) => 
                sum + transaction.amount, 0),
            
            spendingByCategory: dummyTransactions.reduce((acc, transaction) => {
                acc[transaction.category] = (acc[transaction.category] || 0) + transaction.amount;
                return acc;
            }, {} as { [key: string]: number }),
            
            mostFrequentLocation: Object.entries(
                dummyTransactions.reduce((acc, transaction) => {
                    acc[transaction.location] = (acc[transaction.location] || 0) + 1;
                    return acc;
                }, {} as { [key: string]: number })
            ).sort((a, b) => b[1] - a[1])[0]?.[0] || 'No data',
            
            highestExpense: dummyTransactions.reduce((max, transaction) => 
                transaction.amount > max.amount ? transaction : max, 
                dummyTransactions[0]
            )
        };
    }, [dummyTransactions]);

    if (isLoading) {
        return (
            <div className="min-h-screen bg-gray-900 text-white p-4 flex items-center justify-center">
                <div className="text-xl">Loading transactions...</div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="min-h-screen bg-gray-900 text-white p-4 flex items-center justify-center">
                <div className="text-xl text-red-400">{error}</div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-900 text-white p-4">
            <h1 className="text-3xl font-bold mb-6 text-white">Finances Bros Dashboard</h1>
            
            {/* Analytics Overview */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
                <div className="bg-gray-800 p-4 rounded-lg shadow-lg border border-gray-700">
                    <h3 className="font-bold text-gray-300">Total Spent</h3>
                    <p className="text-2xl text-green-400">${analytics.totalSpent.toFixed(2)}</p>
                </div>
                
                <div className="bg-gray-800 p-4 rounded-lg shadow-lg border border-gray-700">
                    <h3 className="font-bold text-gray-300">Most Frequent Location</h3>
                    <p className="text-2xl text-blue-400">{analytics.mostFrequentLocation}</p>
                </div>
                
                <div className="bg-gray-800 p-4 rounded-lg shadow-lg border border-gray-700">
                    <h3 className="font-bold text-gray-300">Highest Single Expense</h3>
                    <p className="text-2xl text-red-400">
                        {analytics.highestExpense 
                            ? `$${analytics.highestExpense.amount.toFixed(2)}`
                            : 'No data'
                        }
                    </p>
                    <p className="text-sm text-gray-400">
                        {analytics.highestExpense?.description || 'No description'}
                    </p>
                </div>
            </div>

            {/* Transactions List */}
            <div>
                <h2 className="text-xl font-bold mb-4 text-white">Recent Transactions</h2>
                {dummyTransactions.length === 0 ? (
                    <div className="text-gray-400">No transactions available</div>
                ) : (
                    <div className="grid gap-4">
                        {dummyTransactions.map((transaction) => (
                            <div key={transaction.transactionId} 
                                 className="bg-gray-800 p-4 rounded-lg shadow-lg border border-gray-700 
                                          flex justify-between items-center hover:bg-gray-700 transition-colors">
                                <div>
                                    <h3 className="font-bold text-white">{transaction.company}</h3>
                                    <p className="text-sm text-gray-400">{transaction.description}</p>
                                    <p className="text-sm text-gray-500">{transaction.date}</p>
                                </div>
                                <div className="text-right">
                                    <p className="font-bold text-green-400">${transaction.amount.toFixed(2)}</p>
                                    <p className="text-sm text-gray-400">{transaction.category}</p>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}

export default HomePage;
