#!/usr/bin/env python
# encoding: utf-8

import socket
import psutil
import sys
from SocketServer import ThreadingMixIn
import BaseHTTPServer
import time
import json

#print len(sys.argv)
if len(sys.argv) != 3:
    print "Usage: {0} LISTEN_HOST LISTEN_PORT".format(sys.argv[0])
    sys.exit()
else:

    LISTEN_HOST = sys.argv[1]
    LISTEN_PORT = int(sys.argv[2])

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
        data = {
                "cpuLoad": psutil.cpu_percent(interval=0.1, percpu=False),
                "memLoad": psutil.virtual_memory().percent  
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
