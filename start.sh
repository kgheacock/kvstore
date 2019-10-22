#!/bin/bash
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
docker run -p 13800:13800 -d --net=kv_subnet --ip=10.10.0.2 --name="main-instance" kv-store:2.0
r=13800
for i in {3..50} 
do
    ((r++))
    docker run -p "$r":13800 -d --net=kv_subnet --ip=10.10.0."$i" --name="follower-instance-$i" -e FORWARDING_ADDRESS="10.10.0.2:13800" kv-store:2.0
done