#!/usr/bin/env python
# encoding: utf-8
# v1.0

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
import ConfigParser

# Config File Example:
# [REDIS]
# host = localhost
# port = 6379
# key = libertyStatusPing
# 
# [QLOG]
# client_id = grAUUne3Lv1eqYW5
# secret_key = F37P1Y9Xpr8TkltjyDHUsvIN2S48htfu
# url = 1.qlog.innotechapp.com
# port = 3001
# tags = probe

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
                data = {
                        "msg": {
                                "cpuLoad": psutil.cpu_percent(interval=0.1, percpu=False),
                                "memLoad": psutil.virtual_memory().percent,
                                "diskUsed": psutil.disk_usage('/').used,
                                "diskFree": psutil.disk_usage('/').free
                        },
                        "tags": QLOG_TAGS,
                        "time": int(item['data']),
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
    if len(sys.argv) != 2:
        print "Usage: {0} CONFIG_FILE".format(sys.argv[0])
        sys.exit()
    else:
        config = ConfigParser.ConfigParser()
        config.readfp(open(sys.argv[1]))
        REDIS_HOST = config.get('REDIS', 'host')
        REDIS_PORT = int(config.get('REDIS', 'port'))
        REDIS_KEY = config.get('REDIS', 'key')
        QLOG_CLIENT_ID = config.get('QLOG', 'client_id')
        QLOG_SECRET_KEY = config.get('QLOG', 'secret_key')
        QLOG_URL = config.get('QLOG', 'url')
        QLOG_PORT = int(config.get('QLOG', 'port'))
        QLOG_TAGS = config.get('QLOG', 'tags')
    print time.asctime(), "Server Starts"
    reconnect()
    sleepTillBreak()
    sys.exit()


