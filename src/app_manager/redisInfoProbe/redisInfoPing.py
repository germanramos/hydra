#!/usr/bin/env python
# encoding: utf-8

import time
import redis
import ConfigParser
import sys

if len(sys.argv) != 2:
    print "Usage: {0} CONFIG_FILE".format(sys.argv[0])
    sys.exit()
else:
    config = ConfigParser.ConfigParser()
    config = ConfigParser.ConfigParser()
    config.readfp(open(sys.argv[1]))
    
    REDIS_HOST = config.get('REDIS', 'host')
    REDIS_PORT = int(config.get('REDIS', 'port'))
    REDIS_KEY = config.get('REDIS', 'key')
    
    t = int(time.time()*1000)
    print "Sending ping " + REDIS_KEY + " with value " + str(t)
        
    r = redis.Redis(REDIS_HOST, REDIS_PORT)
    r.publish(REDIS_KEY, t)
    
    


