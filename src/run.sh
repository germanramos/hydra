#!/bin/bash
killall -9 node
killall -9 python
node ./client_api/app/main.js --port=$2 --env=$1 &
python ./app_manager/app_manager_info_server.py localhost $2 0.0.0.0 7777 $! &
node ./server_api/app/main.js --port=$3 --env=$1 &
cd ./app_manager
python ./app_manager.py &
cd ..