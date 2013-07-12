#!/usr/bin/env python
# encoding: utf-8

import socket
import psutil
import sys
from SocketServer import ThreadingMixIn
import BaseHTTPServer
import time
import json
import random

STRESS_TIME = 45;
HALT_TIME = 90
stress = 0;
halt = 0;

print len(sys.argv)
if len(sys.argv) < 6 or len(sys.argv) > 8:
    print "Usage: {0} CHECK_HOST CHECK_PORT LISTEN_HOST LISTEN_PORT PID [STRESS_TIME] [HALT_TIME]".format(sys.argv[0])
    sys.exit()
else:
    CHECK_HOST = sys.argv[1]
    CHECK_PORT = int(sys.argv[2])
    LISTEN_HOST = sys.argv[3]
    LISTEN_PORT = int(sys.argv[4])
    PID = int(sys.argv[5])
    if len(sys.argv) > 6:
        STRESS_TIME = int(sys.argv[6])
    if len(sys.argv) > 7:
        HALT_TIME = int(sys.argv[7])
        
print "STRESS_TIME = " + str(STRESS_TIME)
print "HALT_TIME = " + str(HALT_TIME)
         
def isOpen(ip,port):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        s.connect((ip, int(port)))
        s.shutdown(2)
        return True
    except:
        return False

class MyHandler(BaseHTTPServer.BaseHTTPRequestHandler):
    def do_OPTIONS(self):           
        self.send_response(200, "ok")       
        self.send_header('Access-Control-Allow-Origin', '*')                
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header("Access-Control-Allow-Headers", "X-Requested-With") 
    def do_HEAD(self):
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.end_headers()
    def do_GET(self):
        print "Path: " + self.path
        """Respond to a GET request."""
        self.send_response(200)
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header("Content-type", "application/json")
        self.end_headers()
        
        global stress
        global halt
        if self.path == "/ready":
            stress = 0;
            halt = 0;
            data = "OK"
        elif self.path == "/stress":
            stress = time.time()
            data = "OK"
        elif self.path == "/halt":
            halt = time.time()
            data = "OK"
        elif isOpen(CHECK_HOST, CHECK_PORT):
            data = {}
            if time.time() - halt < HALT_TIME:
                data["state"] = 1
            elif time.time() - stress < STRESS_TIME:
                data["state"] = 0
                data["cpuLoad"] = 90 + random.randint(0, 10)
                data["memLoad"] = 90 + random.randint(0, 10)
            else:
                data["state"] = 0
                data["cpuLoad"] = psutil.cpu_percent(interval=0.1, percpu=False)
                data["memLoad"] = psutil.virtual_memory().percent  
                
            if self.path == "/extended":
                p = psutil.Process(PID)
                data["connections"] = p.get_connections(kind='inet')
                #connection: fd, family, type, local_address=(ip, port), remote_address=(ip, port), status)
        else:
            data = {
                    "state": 1
            }
        self.wfile.write(json.dumps(data))

class ThreadedHTTPServer(ThreadingMixIn, BaseHTTPServer.HTTPServer):
    pass

#server_class = BaseHTTPServer.HTTPServer
server_class = ThreadedHTTPServer
httpd = server_class((LISTEN_HOST, LISTEN_PORT), MyHandler)
print time.asctime(), "Server Starts - %s:%s" % (LISTEN_HOST, LISTEN_PORT)
try:
    httpd.serve_forever()
except KeyboardInterrupt:
    pass
httpd.server_close()
print time.asctime(), "Server Stops - %s:%s" % (LISTEN_HOST, LISTEN_PORT)
