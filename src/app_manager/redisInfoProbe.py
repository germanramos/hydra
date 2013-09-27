#!/usr/bin/env python
# encoding: utf-8

import socket
import psutil
import sys
from SocketServer import ThreadingMixIn
import BaseHTTPServer
import time
import json
import redis
import threading

REDIS_HOST = "localhost"
REDIS_PORT = 6379
REDIS_KEY = "libertyStatusPing"

class Listener(threading.Thread):
    def __init__(self, r, channels):
        try:
            threading.Thread.__init__(self)
            self.daemon = True
            self.redis = r
            self.pubsub = self.redis.pubsub()
            self.pubsub.subscribe(channels) 
        except Exception, e:
            print "init", str(e)
    def work(self, item):
        try:
            print item['channel'], ":", item['data']
            data = {
                    "cpuLoad": psutil.cpu_percent(interval=0.1, percpu=False),
                    "memLoad": psutil.virtual_memory().percent
            }
            msg = json.dumps(data);
            #TODO: Post to QLOG insted of print
            print msg
        except Exception, e:
            print "work", str(e)
    
    def run(self):
        try:
            for item in self.pubsub.listen():
                if item['data'] == "KILL":
                    self.pubsub.unsubscribe()
                    print self, "unsubscribed and finished"
                    break
                else:
                    self.work(item)
        except Exception, e:
            #print "run", str(e)
            time.sleep(1);
            reconnect()

print time.asctime(), "Server Starts"

def reconnect():
    print "Connecting to Redis..."
    r = redis.Redis(REDIS_HOST, REDIS_PORT)
    client = Listener(r, [REDIS_KEY])
    client.start()

def sleepTillBreak():
    try:
        while True:
            time.sleep(1000)
    except KeyboardInterrupt:
        print ""

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print "Usage: {0} REDIS_HOST REDIS_PORT REDIS_KEY".format(sys.argv[0])
        sys.exit()
    else:
        LISTEN_HOST = sys.argv[1]
        LISTEN_PORT = int(sys.argv[2])
    reconnect()
    sleepTillBreak()
    sys.exit()


