var utils = require('../utils'),
enums = require('../enums');

var defaultApp = {
	appId: null,
	localStrategyEvents : {
		//'42374897239' : localStrategyEnum.INDIFFERENT
	},
	cloudStrategyEvents : {
		//'42374897239': cloudStrategyEnum.INDIFFERENT
	},
	servers : [
		//{
		//	server: 'http://server3/app',
		//	cloud : 'nubeA',
		//  cost : 0,
		//	status: {
		//		cpuLoad: 50, //Cpu load of the server 0-100
		//		memLoad: 50, //Memory load of the server 0-100
		//		timeStamp: 42374897239, //UTC time stamp of this info
		//		stateEvents: {
		//			'42374897239' : state: stateEnum.READY, //Future state of the serve
		//		}
		//	}
		//}
	]
};

module.exports = function(colApp, config){
	var self = {};

	self.create = function(p_app, p_cbk){
		var app = utils.merge({},defaultApp);

		app = utils.merge(app, p_app);

		//Si no tenemos id no creamos la app
		if(app.appId === null){
			p_cbk(null);
			return;
		}

		colApp.insert(app, {w:1}, function(err, items){
			if(p_cbk) p_cbk(err);
		});
	};

	self.getAll = function(p_cbk){
		colApp.find({}).toArray(function(err, items){
			for(var i in items){
				var modified = clean(items[i]);
				if(modified) self.update(items[i]);
			}
			p_cbk(items);
		});
	};

	self.getFromId = function(p_appId, p_cbk){
		var find = {
			appId: p_appId
		};

		colApp.findOne(find, {}, function(err, item){
			if(item !== null) {
				var modified = clean(item);
				if(modified){
					self.update(item);
				}
			}
			p_cbk(item);
		});
	};

	function clean(p_app){
		if(p_app === null) return false;
		var now = new Date().getTime();

		var modified = false;

		//clean localStrategy
		var previousLocal;
		for(var local in p_app.localStrategyEvents){
			if(local < now){
				if(previousLocal > 0){
					delete p_app.localStrategyEvents[previousLocal];
					modified = true;
				}
				previousLocal = local;
			}
		}

		//clean cloudStrategy
		var previousCloud;
		for(var cloud in p_app.cloudStrategyEvents){
			if(cloud < now){
				if(previousCloud > 0){
					delete p_app.cloudStrategyEvents[previousCloud];
					modified = true;
				}
				previousCloud = cloud;
			}
		}

		//clean servers
		var server, previousState;
		var s, S = p_app.servers.length;
		for(s=0;s<S;s++){
			server = p_app.servers[s];
			previousState = -1;
			for(var serverState in server.status.stateEvents){
				if(serverState < now){
					if(serverState < (now - config.app.timeout) && server.status.stateEvents[serverState] != enums.app.stateEnum.UNAVAILABLE){
						server.status.stateEvents[now] = enums.app.stateEnum.UNAVAILABLE;
						modified = true;
					}
					if(previousState > 0){
						delete server.status.stateEvents[previousState];
						modified = true;
					}
					previousState = serverState;
				}
			}
		}

		return modified;
	}

	self.update = function(p_app, p_cbk){
		var find = {
			appId: p_app.appId
		};

		colApp.findOne(find, {}, function(err, oldApp){
			p_app = utils.merge(utils.merge({},defaultApp), p_app);
			if(err || oldApp === null){
				self.create(p_app, p_cbk);
			} else {

				//merging & sort localStrategies schedule
				// merge
				for(var localStrategyEventsIdx in p_app.localStrategyEvents){
					oldApp.localStrategyEvents[localStrategyEventsIdx] = p_app.localStrategyEvents[localStrategyEventsIdx];
				}
				//sort
				oldApp.localStrategyEvents = utils.sortObj(oldApp.localStrategyEvents);

				//merging & sort cloudStrategies schedule
				//merge
				for(var cloudStrategyEventsIdx in p_app.cloudStrategyEvents){
					oldApp.cloudStrategyEvents[cloudStrategyEventsIdx] = p_app.cloudStrategyEvents[cloudStrategyEventsIdx];
				}
				//sort
				oldApp.cloudStrategyEvents = utils.sortObj(oldApp.cloudStrategyEvents);

				//merging servers
				var newServer, serverFound, oldServer;
				var ns, NS = p_app.servers.length;
				var os, OS = oldApp.servers.length;
				for(ns=0;ns<NS;ns++){
					newServer = p_app.servers[ns];

					serverFound = false;
					if(newServer.status !== undefined){
						for(os=0;os<OS;os++){
							oldServer = oldApp.servers[os];
							if(newServer.server == oldServer.server){
								for(var stateEventsIdx in newServer.status.stateEvents){
									oldServer.status.stateEvents[stateEventsIdx] = parseInt(newServer.status.stateEvents[stateEventsIdx]);
								}
								oldServer.status.stateEvents = utils.sortObj(oldServer.status.stateEvents);

								// Copies info
								for(var info in newServer){
									if(info == 'status') continue;
									oldServer[info] = newServer[info];
								}

								// Checks timestamp for cpu/mem updates
								if(newServer.status.timeStamp > oldServer.status.timeStamp || oldServer.status.timeStamp === undefined){
									for(var serverStatusFieldIdx in newServer.status){
										if(serverStatusFieldIdx == 'stateEvents') continue;
										oldServer.status[serverStatusFieldIdx] = newServer.status[serverStatusFieldIdx];
									}
								}

								serverFound = true;
								break;
							}
						}
					}
					if(!serverFound) {
						oldApp.servers.push(newServer);
					}

				}

				clean(oldApp);

				colApp.update(find, oldApp, function(err){
					if(p_cbk) p_cbk(err);
				});
			}
		});
	};

	self.remove = function(p_id, p_cbk){
		var find = {
			appId: p_id
		};

		colApp.remove(find, function(err, item){
			if(err) {
				p_cbk(err);
			}
			else {
				p_cbk(null);
			}
		});
	};

	self.availableServers = function (p_appId, p_cbk){
		self.getFromId(p_appId, function(app){
			if(app === null){
				p_cbk(null);
				return;
			}

			self.balanceServers(app, p_cbk);
		});
	};

	function onlineClouds(p_app){
		var clouds = [];

		for(var serverIdx in p_app.servers){
			var server = p_app.servers[serverIdx];
			for(var serverStateIdx in server.status.stateEvents){
				if(server.status.stateEvents[serverStateIdx] == enums.app.stateEnum.READY){
					if(clouds.indexOf(server.cloud) == -1) clouds.push(server.cloud);
				}

				break;
			}
		}
		return clouds;
	}

	function onlineCloudsLoad(p_app){
		var clouds = onlineClouds(p_app);
		var loads = [];

		var servers, load;
		var c,C= clouds.length;
		for(c=0;c<C;c++){
			servers = onlineServersLoad(p_app, clouds[c]);
			load = 0;
			var s,S=servers.length;
			for(s=0;s<S;s++){
				load += servers[s];
			}
			load=load/S;
			loads.push(load);
		}
		return loads;
	}

	function onlineCloudsCost(p_app){
		var clouds = onlineClouds(p_app);
		var costs = [];

		var servers, cost;
		var c,C= clouds.length;
		for(c=0;c<C;c++){
			servers = onlineServersCost(p_app, clouds[c]);
			cost = 0;
			var s,S=servers.length;
			for(s=0;s<S;s++){
				cost += servers[s];
			}
			cost = cost / S;
			costs.push(cost);
		}
		return costs;
	}


	function onlineServers(p_app, p_cloud){
		var servers = [];

		for(var serverIdx in p_app.servers){
			var server = p_app.servers[serverIdx];
			if(p_cloud && server.cloud != p_cloud) continue; // not in current cloud
			for(var serverStateIdx in server.status.stateEvents){
				if(server.status.stateEvents[serverStateIdx] == enums.app.stateEnum.READY){
					servers.push(server.server);
				}

				break;
			}
		}

		return servers;
	}

	function onlineServersLoad(p_app, p_cloud){
		var servers = [];

		for(var serverIdx in p_app.servers){
			var server = p_app.servers[serverIdx];
			if(p_cloud && server.cloud != p_cloud) continue; // not in current cloud
			for(var serverStateIdx in server.status.stateEvents){
				if(server.status.stateEvents[serverStateIdx] == enums.app.stateEnum.READY){
					servers.push(parseFloat(server.status.cpuLoad) + parseFloat(server.status.memLoad));
				}

				break;
			}
		}
		return servers;
	}

	function onlineServersCost(p_app, p_cloud){
		var servers = [];

		for(var serverIdx in p_app.servers){
			var server = p_app.servers[serverIdx];
			if(p_cloud && server.cloud != p_cloud) continue; // not in current cloud
			for(var serverStateIdx in server.status.stateEvents){
				if(server.status.stateEvents[serverStateIdx] == enums.app.stateEnum.READY){
					servers.push(parseInt(server.cost));
				}

				break;
			}
		}
		return servers;
	}

	function localStrategy(p_app){
		//current strategy
		var currentStrategy = enums.app.localStrategyEnum.INDIFFERENT;
		for(var localStrategyIdx in p_app.localStrategyEvents){
			currentStrategy = p_app.localStrategyEvents[localStrategyIdx];
			break;
		}
		return currentStrategy;
	}

	function cloudStrategy(p_app){
		//current strategy
		var currentStrategy = enums.app.cloudStrategyEnum.INDIFFERENT;
		for(var cloudStrategyIdx in p_app.cloudStrategyEvents){
			currentStrategy = p_app.cloudStrategyEvents[cloudStrategyIdx];
			break;
		}
		return currentStrategy;
	}

	var cloudCurrentRoundRobin = {};
	var localCurrentRoundRobin = {};
	self.balanceServers = function(p_app, p_cbk){
		var appId = p_app.appId;
		var servers = [];
		var clouds = onlineClouds(p_app);
		var currentLocalStrategy = parseInt(localStrategy(p_app));
		var currentCloudStrategy = parseInt(cloudStrategy(p_app));

		var cutPoint,pre,post,loads, costs;

		var c,C, currentCloud;
		// -------------
		// CLOUD BALANCE
		// -------------

		switch(currentCloudStrategy){
			case enums.app.cloudStrategyEnum.INDIFFERENT:
				cutPoint = Math.floor(Math.random()*clouds.length);
				//cortamos la baraja :)
				pre = clouds.slice(0,cutPoint);
				post = clouds.slice(cutPoint);
				clouds = post.concat(pre);
				break;

			case enums.app.cloudStrategyEnum.ROUND_ROBIN:
				if(cloudCurrentRoundRobin[appId] === undefined) {
					cloudCurrentRoundRobin[appId] = 0;
				}

				if(cloudCurrentRoundRobin[appId] >= clouds.length) {
					cloudCurrentRoundRobin[appId] = 0;
				}

				//cortamos la baraja :)
				pre = clouds.slice(0,cloudCurrentRoundRobin[appId]);
				post = clouds.slice(cloudCurrentRoundRobin[appId]);
				clouds = post.concat(pre);
				cloudCurrentRoundRobin[appId]++;
				break;

			case enums.app.cloudStrategyEnum.CHEAPEST:
				costs = onlineCloudsCost(p_app);
				loads = onlineCloudsLoad(p_app);

				var cloudCosts = [];
				C = clouds.length;
				for(c=0;c<C;c++){
					cloudCosts.push({k:clouds[c],c:costs[c],l:loads[c]});
				}
				cloudCosts.sort(function(a,b){
					// if(a.l > 160 && b.l > 160) return a.c - b.c;
					if(a.l > 160 && b.l < 160) return 1;
					if(a.l < 160 && b.l > 160) return -1;
					return a.c - b.c;
				});
				clouds=[];
				for(c=0;c<C;c++){
					clouds.push(cloudCosts[c].k);
				}

				break;

			case enums.app.cloudStrategyEnum.CLOUD_LOAD:
				loads = onlineCloudsLoad(p_app);

				var cloudLoads = [];
				C = clouds.length;
				for(c=0;c<C;c++){
					cloudLoads.push({k:clouds[c],v:loads[c]});
				}
				cloudLoads.sort(function(a,b){
					return a.v - b.v;
				});
				clouds=[];
				for(c=0;c<C;c++){
					clouds.push(cloudLoads[c].k);
				}

				break;

			default:
				break;
		}

		servers = onlineServers(p_app, clouds[0]);

		// --------------
		// SERVER BALANCE
		// --------------

		switch(currentLocalStrategy){

			//LOCAL INDIFFERENT
			case enums.app.localStrategyEnum.INDIFFERENT:
				cutPoint = Math.floor(Math.random()*servers.length);
				//cortamos la baraja :)
				pre = servers.slice(0,cutPoint);
				post = servers.slice(cutPoint);
				servers = post.concat(pre);
				break;

			//LOCAL ROUND ROBIN
			case enums.app.localStrategyEnum.ROUND_ROBIN:
				if(localCurrentRoundRobin[appId] === undefined){
					localCurrentRoundRobin[appId] = {};
				}
				if(localCurrentRoundRobin[appId][clouds[0]] === undefined){
					localCurrentRoundRobin[appId][clouds[0]] = 0;
				}

				if(localCurrentRoundRobin[appId][clouds[0]] >= servers.length){
					localCurrentRoundRobin[appId][clouds[0]] = 0;
				}

				//cortamos la baraja :)
				pre = servers.slice(0,localCurrentRoundRobin[appId][clouds[0]]);
				post = servers.slice(localCurrentRoundRobin[appId][clouds[0]]);
				servers = post.concat(pre);
				localCurrentRoundRobin[appId][clouds[0]]++;
				break;

			//LOCAL SERVER LOAD
			case enums.app.localStrategyEnum.SERVER_LOAD:
				loads = onlineServersLoad(p_app, clouds[0]);

				var serverLoads = [];
				var s,S = servers.length;
				for(s=0;s<S;s++){
					serverLoads.push({k:servers[s],v:loads[s]});
				}

				serverLoads.sort(function(a,b){
					return a.v - b.v;
				});

				servers = [];
				for(s=0;s<S;s++){
					servers.push(serverLoads[s].k);
				}

				break;

			case enums.app.localStrategyEnum.CHEAPEST:
				costs = onlineServersCost(p_app, clouds[0]);

				var serverCosts = [];
				var s,S = servers.length;
				for(s=0;s<S;s++){
					serverCosts.push({k:servers[s],v:costs[s]});
				}

				serverCosts.sort(function(a,b){
					return a.v - b.v;
				});

				servers = [];
				for(s=0;s<S;s++){
					servers.push(serverCosts[s].k);
				}

				break;

			default:
				break;
		}

		//adding servers of other clouds
		C = clouds.length;
		var cloudServers = [];
		for(c = 1; c<C;c++){
			cloudServers = onlineServers(p_app,clouds[c]);
			servers = servers.concat(cloudServers);
		}

		p_cbk(servers);
	};

	return self;
};
