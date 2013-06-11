#!/bin/bash
killall -9 node
killall -9 python
node ./client_api/app/main.js --port=7001 --env=local &
python ./app_manager/app_manager_info_server.py localhost 7001 0.0.0.0 7777 $! &
node ./server_api/app/main.js --port=7002 --env=local &
#node ./monitor/app/main.js --port=0 --env=local &
node ./client_server/app/main.js --port=7000 --env=local &
cd ./app_manager
python ./app_manager.py &
cd ..
