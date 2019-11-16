export ADDRESS="localhost:8080"
export VIEW="localhost:8080,localhost:8081,localhost:8082"
go run main.go &
export ADDRESS="localhost:8081"
export VIEW="localhost:8080,localhost:8081,localhost:8082"
go run main.go &
export ADDRESS="localhost:8082"
export VIEW="localhost:8080,localhost:8081,localhost:8082"
go run main.go 