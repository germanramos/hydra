#!/usr/bin/env python
# encoding: utf-8
'''
app_manager -- Generin app manager for hydra

The app manager is in charge of check and monitor one or several servers and update the status information at one or several Hydra Servers using the restful server AP.

The basic functionality is to notify to one Hydra Server when an application is Started, Stopping, or Removed. In addition, it will provide information about the server health status like CPU and memory usage and any useful information like the size of the server or the prefered balance strategy.

All these information should be updated periodically. If not, the hydra server will assume that the servers are shutted down.

@author:     German Ramos Garcia
            
@copyright:  2013 Next Limit Technologies. All rights reserved.
            
@license:    license

@contact:    german.ramos.garcia@bbva.com
@deffield    updated: Updated
'''

import config

import time
import psutil
import socket
import sys
import os
import json
import logging
from logging.config import fileConfig 
import time        
from optparse import OptionParser
import urllib2

__all__ = []
__version__ = 0.3
__date__ = '2013-05-29'
__updated__ = '2013-05-29'

DEBUG = 1
TESTRUN = 0
PROFILE = 0
     
def isOpen(ip,port):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        s.connect((ip, int(port)))
        s.shutdown(2)
        return True
    except:
        return False

def remoteRunBackground(vm, command, logfilename, debug=None):
    return _remoteRun(vm, command, logfilename, debug, True)

def remoteRunForeground(vm, command, logfilename, debug=None):
    return _remoteRun(vm, command, logfilename, debug, False)

def _remoteRun(vm, command, logfilename, debug=None, background = False):
    """Execute a background command(command) in a remote virtual machine(vm) logging into logfilename"""
    if vm.ip == None:
        raise Exception("Virtual machine with no ip")
    if logfilename == None:
        extendedCommand = command + " >/dev/null 2>&1 </dev/null"
    else:
        extendedCommand = command + " >" + logfilename + " 2>&1 </dev/null"
    if background == True:
        extendedCommand = extendedCommand + " &"
    if debug:
        debug(vm.name + ": Launching background command: " + extendedCommand)
    
    sshwrapper = "ssh -o StrictHostKeyChecking=no {0}@{1} \"{2}\"".format(vm.userSSH, vm.ip, extendedCommand)
    returncode = os.system(sshwrapper)
    return returncode
      
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
                logging.log(logging.DEBUG, "Iteration")
                for server in config.SERVERS:
                    logging.log(logging.DEBUG, server)
                    if server == "127.0.0.1" or server == "localhost":
                        if isOpen(server, config.PORT):
                            state = 0
                            cpuLoad = psutil.cpu_percent(interval=0.1, percpu=False)
                            memLoad = psutil.virtual_memory().percent
                        else:
                            state = 1
                            cpuLoad = 0
                            memLoad = 0
                    else:
                        pass
                        #TODO
                    #Inform
                    server_status_item = {
                            
                            "state": state,
                            "cpuLoad": cpuLoad,
                            "memLoad": memLoad,
                            "timeStamp": int(round(time.time() * 1000))
                    }
                    server_item = {
                              "server": server,
                              "status": server_status_item
                    }
                    servers.append(server_item)
                data = {
                        "servers": servers
                }
                answer = json.dumps(data)
                print answer;                
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
        profile_filename = 'app_anager_profile.txt'
        cProfile.run('main()', profile_filename)
        statsfile = open("profile_stats.txt", "wb")
        p = pstats.Stats(profile_filename, stream=statsfile)
        stats = p.strip_dirs().sort_stats('cumulative')
        stats.print_stats()
        statsfile.close()
        sys.exit(0)
    sys.exit(main())