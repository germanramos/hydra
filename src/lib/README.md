
------------
REDIS ON MAC
------------

Get source and install

$ wget http://redis.googlecode.com/files/redis-2.2.12.tar.gz 
$ tar xzf redis-2.2.12.tar.gz 
$ cd redis-2.2.12 
$ make
$ make install 
And it should be up and running on:

/usr/local/bin/redis-server
For Redis server to start on boot

$ sudo nano /Library/LaunchDaemons/org.redis.redis-server.plist
Copy the  following to the the file you just created. (go to link http://pastebin.com/C0r7L5U0)

Create the log dir if it doesn’t exists yet

sudo mkdir /var/log/redis
Load and launch the Daemon:

sudo launchctl load /Library/LaunchDaemons/org.redis.redis-server.plist
sudo launchctl start org.redis.redis-server

-------------
RABBIT ON MAC 
-------------
/usr/local/sbin/rabbitmq-server