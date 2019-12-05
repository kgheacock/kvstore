export ADDRESS="localhost:8082"
export VIEW="localhost:8081,localhost:8082"
export REPL_FACTOR="1"
go run main.go &

export ADDRESS="localhost:8081"
export VIEW="localhost:8081,localhost:8082"
export REPL_FACTOR="1"
go run main.go
