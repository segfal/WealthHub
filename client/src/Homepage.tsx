import { useState } from 'react';

const HomePage = () => {
    const currentDate = new Date(); 
    const formattedDate = currentDate.toLocaleDateString();
    const [date, setDate] = useState(formattedDate);
    const [money, setMoney] = useState(null); 
    const [location, setLocation] = useState(null); 
    const [item, setItem] = useState(null);
    const [category, setCategory] = useState(null);
    const dummyTransactions = [
        {
            transactionId: "TXN001",
            category: "Groceries",
            location: "Walmart, New York",
            company: "Walmart",
            amount: 120.50,
            date: "2024-12-01",
            description: "Weekly groceries",
            status: "Completed"
        },
        {
            transactionId: "TXN002",
            category: "Dining",
            location: "Chipotle, San Francisco",
            company: "Chipotle",
            amount: 25.75,
            date: "2024-12-02",
            description: "Lunch with friends",
            status: "Completed"
        },
        {
            transactionId: "TXN003",
            category: "Transportation",
            location: "Uber",
            company: "Uber",
            amount: 18.00,
            date: "2024-12-03",
            description: "Ride to office",
            status: "Completed"
        },
        {
            transactionId: "TXN004",
            category: "Entertainment",
            location: "AMC Theatres, Los Angeles",
            company: "AMC Theatres",
            amount: 15.00,
            date: "2024-12-03",
            description: "Movie night",
            status: "Completed"
        },
        {
            transactionId: "TXN005",
            category: "Shopping",
            location: "Amazon",
            company: "Amazon",
            amount: 75.30,
            date: "2024-12-04",
            description: "Bought books and gadgets",
            status: "Completed"
        },
        {
            transactionId: "TXN006",
            category: "Health",
            location: "Walgreens, Chicago",
            company: "Walgreens",
            amount: 42.00,
            date: "2024-12-05",
            description: "Medication purchase",
            status: "Completed"
        },
        {
            transactionId: "TXN007",
            category: "Utilities",
            location: "Online",
            company: "Con Edison",
            amount: 95.60,
            date: "2024-12-05",
            description: "Electricity bill",
            status: "Completed"
        },
        {
            transactionId: "TXN008",
            category: "Travel",
            location: "Delta Airlines",
            company: "Delta Airlines",
            amount: 350.00,
            date: "2024-12-06",
            description: "Flight to Miami",
            status: "Completed"
        },
        {
            transactionId: "TXN009",
            category: "Groceries",
            location: "Trader Joe's, Boston",
            company: "Trader Joe's",
            amount: 87.25,
            date: "2024-12-07",
            description: "Weekly groceries",
            status: "Completed"
        },
        {
            transactionId: "TXN010",
            category: "Dining",
            location: "Starbucks, Seattle",
            company: "Starbucks",
            amount: 12.45,
            date: "2024-12-07",
            description: "Morning coffee",
            status: "Completed"
        },
        {
            transactionId: "TXN011",
            category: "Shopping",
            location: "Target, Houston",
            company: "Target",
            amount: 58.90,
            date: "2024-12-08",
            description: "Bought home essentials",
            status: "Completed"
        },
        {
            transactionId: "TXN012",
            category: "Transportation",
            location: "Lyft",
            company: "Lyft",
            amount: 22.50,
            date: "2024-12-08",
            description: "Ride to airport",
            status: "Completed"
        },
        {
            transactionId: "TXN013",
            category: "Health",
            location: "CVS, San Diego",
            company: "CVS",
            amount: 30.00,
            date: "2024-12-09",
            description: "Vitamins and supplements",
            status: "Completed"
        },
        {
            transactionId: "TXN014",
            category: "Utilities",
            location: "Online",
            company: "AT&T",
            amount: 65.00,
            date: "2024-12-09",
            description: "Internet bill",
            status: "Completed"
        },
        {
            transactionId: "TXN015",
            category: "Entertainment",
            location: "Spotify",
            company: "Spotify",
            amount: 9.99,
            date: "2024-12-10",
            description: "Monthly subscription",
            status: "Completed"
        },
        {
            transactionId: "TXN016",
            category: "Travel",
            location: "Airbnb",
            company: "Airbnb",
            amount: 250.00,
            date: "2024-12-10",
            description: "Stay in Miami",
            status: "Completed"
        },
        {
            transactionId: "TXN017",
            category: "Groceries",
            location: "Costco, Atlanta",
            company: "Costco",
            amount: 105.80,
            date: "2024-12-11",
            description: "Bulk groceries",
            status: "Completed"
        },
        {
            transactionId: "TXN018",
            category: "Dining",
            location: "Pizza Hut, Chicago",
            company: "Pizza Hut",
            amount: 22.35,
            date: "2024-12-11",
            description: "Family dinner",
            status: "Completed"
        },
        {
            transactionId: "TXN019",
            category: "Shopping",
            location: "Best Buy, San Jose",
            company: "Best Buy",
            amount: 120.00,
            date: "2024-12-12",
            description: "Bought a new headset",
            status: "Completed"
        },
        {
            transactionId: "TXN020",
            category: "Transportation",
            location: "Metro, Boston",
            company: "MBTA",
            amount: 2.75,
            date: "2024-12-12",
            description: "Subway ticket",
            status: "Completed"
        },
        {
            transactionId: "TXN021",
            category: "Health",
            location: "Rite Aid, Dallas",
            company: "Rite Aid",
            amount: 20.00,
            date: "2024-12-13",
            description: "First aid kit",
            status: "Completed"
        },
        {
            transactionId: "TXN022",
            category: "Utilities",
            location: "Online",
            company: "Verizon",
            amount: 80.00,
            date: "2024-12-13",
            description: "Phone bill",
            status: "Completed"
        },
        {
            transactionId: "TXN023",
            category: "Entertainment",
            location: "Netflix",
            company: "Netflix",
            amount: 15.49,
            date: "2024-12-14",
            description: "Monthly subscription",
            status: "Completed"
        },
        {
            transactionId: "TXN024",
            category: "Travel",
            location: "Hertz, Orlando",
            company: "Hertz",
            amount: 150.00,
            date: "2024-12-14",
            description: "Car rental",
            status: "Completed"
        },
        {
            transactionId: "TXN025",
            category: "Groceries",
            location: "Whole Foods, Austin",
            company: "Whole Foods",
            amount: 75.60,
            date: "2024-12-15",
            description: "Organic groceries",
            status: "Completed"
        }
    ];

    const handleDate = () => {
        return null
    }

    const handleMoney = () => {
        return null
    }  

    const handleLocation = () => {
        return null
    }  

    const handleItem = () => {
        return null
    } 

    const handleCategory = () => {
        return null
    } 

    return (
      //data being rendered NOT printed -> console.log() 
      //map is a for loop in a function
            <div>
                <h1>Home Page</h1>
                <h2>Today's Date: {date}</h2>
                {dummyTransactions.map((transaction) => (
                    <div key={transaction.transactionId}>
                        <h3>Transaction ID: {transaction.transactionId}</h3> 
                        <p>Category: {transaction.category}</p>
                        <p>Location: {transaction.location}</p> 
                        <p>Company: {transaction.company}</p>
                        <p>Amount: {transaction.amount}</p>
                        <p>Date: {transaction.date}</p>
                        <p>Description: {transaction.description}</p>
                        <p>Status: {transaction.status}</p> 
                    </div>
                ))}
            </div>
        
    );
}

export default HomePage;
