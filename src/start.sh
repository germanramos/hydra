#!/bin/sh
killall -9 node
node ./client_api/app/main.js --port=7001 --env=local &
node ./server_api/app/main.js --port=7002 --env=local &
node ./client_server/app/main.js --port=7000 --env=local &
