#!/usr/bin/env python
# encoding: utf-8
'''
app_manager -- Generic app manager for hydra

The app manager is in charge of check and monitor one or several servers and update the status information at one or several Hydra Servers using the restful server AP.

The basic functionality is to notify to one Hydra Server when an application is Started, Stopping, or Removed. In addition, it will provide information about the server health status like CPU and memory usage and any useful information like the size of the server or the prefered balance strategy.

All these information should be updated periodically. If not, the Hydra server will assume that the servers are shutted down.

@author:     German Ramos Garcia
            
@copyright:  2013 BBVA. All rights reserved.
            
@license:    license

@contact:    german.ramos.garcia@bbva.com
@deffield    updated: Updated
'''

import config

import time
import socket
import sys
import os
import json
import logging
from logging.config import fileConfig   
from optparse import OptionParser
import urllib2
import subprocess

__all__ = []
__version__ = 0.3
__date__ = '2013-05-29'
__updated__ = '2013-05-29'

DEBUG = 1
TESTRUN = 0
PROFILE = 0

class stateEnum:
    READY = 0
    UNAVAILABLE = 1
      
def main(argv=None):
    '''Command line options.'''
    
    program_name = os.path.basename(sys.argv[0])
    program_version = "v0.1"
    program_build_date = "%s" % __updated__
 
    program_version_string = '%%prog %s (%s)' % (program_version, program_build_date)
    program_longdesc = '''''' # optional - give further explanation about what the program does
    program_license = "Copyright 2013 - BBVA"
 
    if argv is None:
        argv = sys.argv[1:]
    try:
        # setup option parser
        parser = OptionParser(version=program_version_string, epilog=program_longdesc, description=program_license)
        parser.add_option("-v", "--verbose", dest="verbose", action="count", help="set verbosity level [default: %default]")

        # process options
        (opts, _args) = parser.parse_args(argv)
        
        if opts.verbose > 0:
            print("verbosity level = %d" % opts.verbose)
                    
        # MAIN BODY #      
        while True:
            try:
                servers = []
                logging.debug("*** BEGIN ITERATION ***")
                for server,user,command in config.SERVERS:
                    logging.debug("Getting info from " + server)
                    try:
                        if server != "127.0.0.1" and server != "localhost":
                            wrapper = config.SSH_CMD + " {0}@{1} \"{2}\"".format(user, server, command)
                        else:
                            wrapper = command
                        output = subprocess.check_output(wrapper, stdin=None, stderr=None, shell=False, universal_newlines=False)
                        lines = output.replace("\r","").split("\n")
                        state = lines[0]
                        cpuLoad = lines[1]
                        memLoad = lines[2]
                    except:
                        state = stateEnum.UNAVAILABLE
                        cpuLoad = 0
                        memLoad = 0
                    #Create server status object and append to the server list
                    server_status_item = { 
                            "state": state,
                            "cpuLoad": cpuLoad,
                            "memLoad": memLoad,
                            "timeStamp": int(round(time.time() * 1000))
                    }
                    logging.debug(server_status_item)
                    server_item = {
                              "server": server,
                              "status": server_status_item
                    }
                    servers.append(server_item)
                #End for
                localStrategiesEvents = [{
                    "localStrategy": config.LOCAL_STRATEGY,
                    "applyTimeStamp": int(round(time.time() * 1000))
                }]
                cloudStrategiesEvents = [{
                    "cloudStrategy": config.CLOUD_STRATEGY,
                    "applyTimeStamp": int(round(time.time() * 1000))
                }]
                data = {
                        "localStrategiesEvents": localStrategiesEvents,
                        "cloudStrategiesEvents": cloudStrategiesEvents,
                        "servers": servers
                }
                answer = json.dumps(data)
                #logging.debug(answer)
                #POST
                for hydra in config.HYDRAS:
                    logging.debug("Posting to " + hydra)                   
                    opener = urllib2.build_opener(urllib2.HTTPHandler)
                    request = urllib2.Request(hydra, answer)
                    request.get_method = lambda: 'POST'
                    url = opener.open(request)
                    if url.code != 200:
                        logging.error("Error connecting with hydra {0}: Code: {1}".format(hydra,url.code))
                    else:
                        logging.debug("Posted OK")
            except Exception, e:
                logging.error("Exception: " + str(e))
            time.sleep(config.SLEEP_TIME)
        
    except Exception, e:
        indent = len(program_name) * " "
        sys.stderr.write(program_name + ": " + repr(e) + "\n")
        sys.stderr.write(indent + "  for help use --help")
        return 2


if __name__ == "__main__":
    #Conf logging
    fileConfig('logging.conf')    
    if DEBUG:
        sys.argv.append("-v")
        logging.getLogger("root").setLevel(logging.DEBUG)
    else:
        sys.argv.append("-v")
    if TESTRUN:
        import doctest
        doctest.testmod()
    if PROFILE:
        import cProfile
        import pstats
        profile_filename = 'app_manager_profile.txt'
        cProfile.run('main()', profile_filename)
        statsfile = open("profile_stats.txt", "wb")
        p = pstats.Stats(profile_filename, stream=statsfile)
        stats = p.strip_dirs().sort_stats('cumulative')
        stats.print_stats()
        statsfile.close()
        sys.exit(0)
    sys.exit(main())