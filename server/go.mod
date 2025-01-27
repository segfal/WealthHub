module server

go 1.21

require (
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
)

require github.com/felixge/httpsnoop v1.0.3 // indirect

replace (
	server/analytics => ./analytics
	server/bills => ./bills
	server/categories => ./categories
	server/handlers => ./handlers
	server/income => ./income
	server/types => ./types
)
