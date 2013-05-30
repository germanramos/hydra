#!/usr/bin/env python
# encoding: utf-8

import socket
import psutil
import sys
         
def isOpen(ip,port):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        s.connect((ip, int(port)))
        s.shutdown(2)
        return True
    except:
        return False
      
 
if len(sys.argv) == 2:
    if isOpen("127.0.0.1", sys.argv[1]):
        print "0"
        print psutil.cpu_percent(interval=0.1, percpu=False)
        print psutil.virtual_memory().percent
    else:
        print "1"
