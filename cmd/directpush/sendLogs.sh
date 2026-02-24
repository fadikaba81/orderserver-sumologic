#!/bin/bash


i=0
while [ $i -lt 300 ]
do  
	curl -k https://localhost -H "X-Correlation-ID: OENI7V9p1kqe"
       	i=$((i+1))
done
