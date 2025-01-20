#!/bin/bash

# Function to check if a command exists
check_command() {
    if ! command -v $1 &> /dev/null; then
        echo "❌ Error: $1 is not installed"
        case $1 in
            "npm")
                echo "Please install Node.js and npm from: https://nodejs.org/"
                ;;
            "node")
                echo "Please install Node.js from: https://nodejs.org/"
                ;;
            "go")
                echo "Please install Go from: https://golang.org/doc/install"
                ;;
        esac
        exit 1
    else
        echo "✅ Found $1"
    fi
}

echo "🔍 Checking prerequisites..."

# Check for required tools
check_command "npm"
check_command "node"
check_command "go"

echo "📦 Installing dependencies..."

# Install client dependencies
echo "📱 Setting up client..."
cd client
if npm install; then
    echo "✅ Client dependencies installed successfully"
else
    echo "❌ Failed to install client dependencies"
    exit 1
fi

# Return to root and install server dependencies
cd ..
echo "🖥️  Setting up server..."
cd server
if go mod download; then
    echo "✅ Server dependencies installed successfully"
else
    echo "❌ Failed to install server dependencies"
    exit 1
fi

# Return to root
cd ..

echo "🎉 Setup completed successfully!"
echo "👉 Run 'npm run dev' to start both client and server" 