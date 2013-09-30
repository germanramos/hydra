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
import urllib2

REDIS_HOST = "localhost"
REDIS_PORT = 6379
REDIS_KEY = "libertyStatusPing"

QLOG_CLIENT_ID = "grAUUne3Lv1eqYW5"
QLOG_SECRET_KEY = "F37P1Y9Xpr8TkltjyDHUsvIN2S48htfu"
QLOG_URL = "1.qlog.innotechapp.com"
QLOG_PORT = 3001

#TODO: QLOG parameters

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
            if int(item['data']) > 1:
                print item['channel'], ":", item['data']
                systemStatus = {
                                "cpuLoad": psutil.cpu_percent(interval=0.1, percpu=False),
                                "memLoad": psutil.virtual_memory().percent
                }
                data = {
                        "msg": json.dumps(systemStatus),
                        "tags": "Probe",
                        "time": item['data'],
                        "secretKey": QLOG_SECRET_KEY
                }
                msg = json.dumps(data);
                url = "http://" + QLOG_URL + ":" + str(QLOG_PORT) + "/app/" + QLOG_CLIENT_ID + "/log"
                print "Posting data to " + url
                print data
                req = urllib2.Request(url, msg, {'Content-Type': 'application/json'})
                req.get_method = lambda: 'PUT'
                f = urllib2.urlopen(req)
                response = f.read()
                f.close()
                print "Done OK"
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
    print "Connected"

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


