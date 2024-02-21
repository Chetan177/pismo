#!/bin/bash

docker build . -t pismo

if [ $? -ne 0 ]; then
    echo "Docker build failed"
    exit 1
fi

docker-compose up -d
