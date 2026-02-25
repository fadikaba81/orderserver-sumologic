#!/bin/bash

for i in {1..100}
do
    START=$(date -u -d "$i minutes ago" +"%Y-%m-%dT%H:%M:%SZ")
    END=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

    curl -sk "https://localhost/order?start=$START&end=$END" > /dev/null

    echo "Request $i sent"
    sleep 0.5
done