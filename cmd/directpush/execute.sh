#! /bin/sh

nohup ./build.sh > /dev/null 2>&1 &
sleep 5 
./sendLogs.sh
