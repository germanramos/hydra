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

import time
import sys
import os
import json
import logging
from logging.config import fileConfig   
from optparse import OptionParser
import urllib2
import ConfigParser

__all__ = []
__version__ = 1.0
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
        config = ConfigParser.ConfigParser()
        config.read(['app_manager.cfg', os.path.expanduser('~/app_manager.cfg'), '/etc/app_manager.cfg'])  
        while True:
            try:
                servers = []
                logging.debug("*** BEGIN ITERATION ***")
                for key,server in config.items("SERVERS"):
                    logging.debug("Getting info from " + server)
                    try:
                        response = urllib2.urlopen(server)
                        output = response.read()
                        lines = output.replace("\r","").split("\n")
                        state = lines[0]
                        cpuLoad = lines[1]
                        memLoad = lines[2]
                    except Exception, e:
                        logging.error("Exception: " + str(e))
                        state = stateEnum.UNAVAILABLE
                        cpuLoad = 0
                        memLoad = 0
                    #Create server status object and append to the server list
                    timestamp = int(round(time.time() * 1000))
                    server_item = {
                        "server": server,
                        "status": {
                                   "cost": config.get("MAIN", "cost"),
                                   "cpuLoad": cpuLoad,
                                   "memLoad": memLoad,
                                   "timeStamp": timestamp,
                                   "stateEvents": {
                                                   timestamp: state
                                                   }
                                   }
                    }
                    logging.debug(server_item)
                    servers.append(server_item)
                #End servers for
                timestamp = int(round(time.time() * 1000))
                localStrategyEvents = {
                                         timestamp: config.get("MAIN", "local_strategy")
                }
                cloudStrategyEvents = {
                                         timestamp: config.get("MAIN", "cloud_strategy")
                }
                data = {
                        "localStrategyEvents": localStrategyEvents,
                        "cloudStrategyEvents": cloudStrategyEvents,
                        "servers": servers
                }
                answer = json.dumps(data)
                #logging.debug(answer)
                #POST
                for key,hydra in config.items("HYDRAS"):
                    logging.debug("Posting to " + hydra)                   
                    opener = urllib2.build_opener(urllib2.HTTPHandler)
                    request = urllib2.Request(hydra + "/app/" + config.get("MAIN", "app_id"), answer)
                    request.add_header("content-type", "application/json")
                    #request.get_method = lambda: 'POST'
                    url = opener.open(request)
                    if url.code != 200:
                        logging.error("Error connecting with hydra {0}: Code: {1}".format(hydra,url.code))
                    else:
                        logging.debug("Posted OK")
            except Exception, e:
                logging.error("Exception: " + str(e))
            time.sleep(config.getint("MAIN", "sleep_time"));
        
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