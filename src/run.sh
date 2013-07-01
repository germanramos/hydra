#!/bin/bash
killall -9 node
killall -9 python
/usr/local/bin/node ./client_api/app/main.js --port=$2 --env=$1 &
python ./app_manager/app_manager_info_server.py localhost $2 0.0.0.0 $4 $! &
/usr/local/bin/node ./server_api/app/main.js --port=$3 --env=$1 &
cd ./app_manager
python ./app_manager.py -c $4 &
cd ..