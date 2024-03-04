#!/bin/bash
if [ -e "proxy.log" ]; then
	echo "proxy.log Exists"
else
	echo "Touching proxy.log"
	touch proxy.log
fi
./init_docker.sh
docker-compose up
