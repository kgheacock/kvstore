#!/bin/bash

rm puts.txt keycount.txt gets.txt

for i in {1..500} 
do
  curl -s --request   PUT \
          --header    "Content-Type: application/json" \
          --write-out "\n%{http_code}\n" \
          --data      "{\"value\": \"sampleValue$i\"}"\
          http://localhost:8080/kv-store/keys/sampleKey$i >> puts.txt
done

for i in {1..510} 
do
  curl -s --request   GET \
          --header    "Content-Type: application/json" \
          --write-out "\n%{http_code}\n" \
          http://localhost:8080/kv-store/keys/sampleKey$i >> gets.txt
done

curl -s --request   GET \
        --header    "Content-Type: application/json" \
        --write-out "\n%{http_code}\n" \
        http://localhost:8080/kv-store/key-count >> keycount.txt


curl -s --request   GET \
        --header    "Content-Type: application/json" \
        --write-out "\n%{http_code}\n" \
        http://localhost:8081/kv-store/key-count >> keycount.txt

curl -s --request   GET \
        --header    "Content-Type: application/json" \
        --write-out "\n%{http_code}\n" \
        http://localhost:8082/kv-store/key-count >> keycount.txt

curl -s --request   GET \
        --header    "Content-Type: application/json" \
        --write-out "\n%{http_code}\n" \
        http://localhost:8083/kv-store/key-count >> keycount.txt

curl -s --request   GET \
        --header    "Content-Type: application/json" \
        --write-out "\n%{http_code}\n" \
        http://localhost:8084/kv-store/key-count >> keycount.txt