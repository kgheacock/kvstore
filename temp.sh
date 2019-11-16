export ADDRESS="localhost:8080"
export VIEW="localhost:8080,localhost:8081"
go run main.go &

export ADDRESS="localhost:8081"
export VIEW="localhost:8080,localhost:8081"
go run main.go 
