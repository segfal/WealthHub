FROM golang:1.21

WORKDIR /app

# Copy the entire project
COPY . .

# Build the application
RUN go build -v -o main

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"] 
