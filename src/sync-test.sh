#!/bin/bash
killall -9 node
clear

node ./server_api/app/main.js --port=7002 --env=local1 &
node ./server_api/app/main.js --port=7012 --env=local2 &
node ./server_api/app/main.js --port=7022 --env=local3 &


sleep 5
node ./test/hydra_sync/sync_test.js &