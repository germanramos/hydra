var utils = require('../utils'),
	enums = require('../enums');

var defaultServer = {
	url: null,
	sibling: false
	//	status: {
	//		cpuLoad: 50, //Cpu load of the server 0-100
	//		memLoad: 50, //Memory load of the server 0-100
	//		timeStamp: 42374897239, //UTC time stamp of this info
	//		stateEvents: {
	//			'42374897239' : stateEnum.READY //Future state of the serve
	//		}
	//	}
};

module.exports = function(colServer, config){
	var self = {};
	self.create = function(p_server, p_cbk){
		var server = utils.merge({},defaultServer);
		server = utils.merge(server, p_server);

		//Si no tenemos url no creamos el server
		if(server.url === null){
			p_cbk(null);
			return;
		}

		colServer.insert(server, {w:1}, function(err, items){
			if(err || items.length === 0){
				p_cbk(null);
			} else {
				p_cbk(items[0]);
			}
		});
	};

	self.getAll = function(p_cbk){
		colServer.find({}).toArray(function(err, items){
			for(var i in items){
				var modified = clean(items[i]);
				if(modified) self.update(items[i]);
			}
			p_cbk(items);
		});
	};

	self.getFromUrl = function(p_url, p_cbk){
		var find = {
			url: p_url
		};

		colServer.findOne(find, {}, function(err, item){
			var modified = clean(item);
			if(modified){
				self.update(item);
			}
			p_cbk(item);
		});
	};

	function clean(p_server){
		var now = new Date().getTime();

		var modified = false;

		//clean states
		var previousState;
		for(var serverState in p_server.status.stateEvents){
			if(serverState < now){
				if(serverState < (now - config.server.timeout) && p_server.status.stateEvents[serverState] != enums.server.stateEnum.UNAVAILABLE){
					p_server.status.stateEvents[now] = enums.server.stateEnum.UNAVAILABLE;
					modified = true;
				}
				if(previousState > 0){
					delete p_server.status.stateEvents[previousState];
					modified = true;
				}
				previousState = serverState;
			}
		}

		return modified;
	}

	self.update = function(p_server, p_cbk){
		var find = {
			url: p_server.url
		};

		colServer.findOne(find, {}, function(err, oldServer){
			p_server = utils.merge(utils.merge({},defaultServer), p_server);
			if(err || oldServer === null){
				self.create(p_server, p_cbk);
			} else {
				//merge state schedule
				for(var stateEventsIdx in p_server.status.stateEvents){
					oldServer.status.stateEvents[stateEventsIdx] =  p_server.status.stateEvents[stateEventsIdx];
				}
				oldServer.status.stateEvents = utils.sortObj(oldServer.status.stateEvents);

				// Checks timestamp for cpu/mem updates
				if(p_server.status.timeStamp > oldServer.status.timeStamp){
					for(var serverStatusFieldIdx in p_server.status){
						if(serverStatusFieldIdx == 'stateEvents') continue;
						oldServer.status[serverStatusFieldIdx] = p_server.status[serverStatusFieldIdx];
					}
				}

				clean(oldServer);

				colServer.update(find, oldServer, function(err){
					if(p_cbk) p_cbk();
				});
			}
		});

	};

	self.remove = function(p_url, p_cbk){
		var find = {
			url: p_url
		};

		colServer.remove(find, function(err, item){
			if(err) {
				p_cbk(err);
			}
			else {
				p_cbk(null);
			}
		});
	};

	self.availableServers = function (p_cbk){
		self.getAll(function(knownServers){
			var servers = [];
			for(var serverIdx in knownServers){
				var server = knownServers[serverIdx];
				for(var serverStateIdx in server.status.stateEvents){
					if(server.status.stateEvents[serverStateIdx] == enums.app.stateEnum.READY){
						servers.push(server.url);
					}

					break;
				}
			}
			p_cbk(servers);
		});
	};

	return self;
};
