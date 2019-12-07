export ADDRESS="localhost:8082"
export VIEW="localhost:8081,localhost:8082,localhost:8083,localhost:8084"
export REPL_FACTOR="2"
go run main.go &

export ADDRESS="localhost:8081"
export VIEW="localhost:8081,localhost:8082,localhost:8083,localhost:8084"
export REPL_FACTOR="2"
go run main.go &

export ADDRESS="localhost:8083"
export VIEW="localhost:8081,localhost:8082,localhost:8083,localhost:8084"
export REPL_FACTOR="2"
go run main.go &

export ADDRESS="localhost:8084"
export VIEW="localhost:8081,localhost:8082,localhost:8083,localhost:8084"
export REPL_FACTOR="2"
go run main.go 
