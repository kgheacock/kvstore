#!/bin/sh
docker stop main-instance
docker stop follower-instance
docker rm main-instance
docker rm follower-instance
docker run -p 13800:13800 --net=kv_subnet --ip=10.10.0.2 --name="main-instance" kv-store:2.0
docker run -p 13801:13800 --net=kv_subnet --ip=10.10.0.2 --name="follower-instance" -e FORWARDING_ADDRESS="10.10.0.2:13800" kv-store:2.0