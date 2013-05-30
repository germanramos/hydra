#CONSTANTS
class localStrategyEnum:
    INDIFFERENT = 0
    ROUND_ROBIN = 1
    SERVER_LOAD = 2
    
class cloudStrategyEnum:
    INDIFFERENT = 0
    ROUND_ROBIN = 1
    CHEAPEST = 2
    CLOUD_LOAD = 3

#CONFIG
APP_ID = "time"     # Unique application id
SLEEP_TIME = 5      # Time to sleep after each iteration
LOCAL_STRATEGY = localStrategyEnum.INDIFFERENT  # 0 is indifferent
CLOUD_STRATEGY = cloudStrategyEnum.INDIFFERENT # 0 is indifferent
#List of app servers: Tuples of (server, user, command)
SERVERS = [ 
           ("127.0.0.1", None, "python D:\\german\\workspace\\hydra\\src\\app_manager\\app_manager_info.py 888"),
           ("vallecas.no-ip.org", "ranita", "python /home/ranita/app_manager_info.py 80")
           ]
#List of hydra servers
HYDRAS = ["http://127.0.0.1:1337"]
SSH_CMD = "ssh -i \"C:\Documents and Settings\e029359\.ssh\id_rsa\" -o StrictHostKeyChecking=no"
