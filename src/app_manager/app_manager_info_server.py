#!/usr/bin/env python
# encoding: utf-8

import socket
import psutil
import sys
import BaseHTTPServer
import time
import json

if len(sys.argv) != 6:
    print "Usage: {0} CHECK_HOST CHECK_PORT LISTEN_HOST LISTEN_PORT PID".format(sys.argv[0])
    sys.exit()
else:
    CHECK_HOST = sys.argv[1]
    CHECK_PORT = int(sys.argv[2])
    LISTEN_HOST = sys.argv[3]
    LISTEN_PORT = int(sys.argv[4])
    PID = int(sys.argv[5])
         
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
        if isOpen(CHECK_HOST, CHECK_PORT):
            if (self.path == "/extended"):
                p = psutil.Process(PID)
                data = {
                        "state": 0,
                        "cpuLoad": psutil.cpu_percent(interval=0.1, percpu=False),
                        "memLoad": psutil.virtual_memory().percent,
                        "connections": p.get_connections(kind='inet')
                        #connection: fd, family, type, local_address=(ip, port), remote_address=(ip, port), status)
                        }
                
            else:
                data = {
                        "state": 0,
                        "cpuLoad": psutil.cpu_percent(interval=0.1, percpu=False),
                        "memLoad": psutil.virtual_memory().percent
                        }
        else:
            data = {
                    "state": 1
                    }
        self.wfile.write(json.dumps(data))

server_class = BaseHTTPServer.HTTPServer
httpd = server_class((LISTEN_HOST, LISTEN_PORT), MyHandler)
print time.asctime(), "Server Starts - %s:%s" % (LISTEN_HOST, LISTEN_PORT)
try:
    httpd.serve_forever()
except KeyboardInterrupt:
    pass
httpd.server_close()
print time.asctime(), "Server Stops - %s:%s" % (LISTEN_HOST, LISTEN_PORT)
