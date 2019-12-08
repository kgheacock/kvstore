docker network rm kv_subnet
docker network create --subnet=10.10.0.0/16 kv_subnet
docker build -t kv-store:4.0 /home/kheacock2/go/src/github.com/colbyleiske/cse138_assignment2/
