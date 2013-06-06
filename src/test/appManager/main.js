var http = require('http'),
	utils = require('../../lib/utils');

var server_api = {
	host: 'localhost',
	port: 7002
};

// -----
// ENUMS
// -----

var localStrategyEnum = {
	INDIFFERENT: 0,
	ROUND_ROBIN: 1,
	SERVER_LOAD: 2
};

var cloudStrategyEnum = {
	INDIFFERENT: 0,
	ROUND_ROBIN: 1,
	CHEAPEST: 2,
	CLOUD_LOAD: 3
};

var stateEnum = {
	READY: 0,
	UNAVAILABLE : 1
};

// ----
// APPS
// ----

var now = new Date().getTime();

var apps = [
{
	appId: 1,
	localStrategyEvents : {},
	cloudStrategyEvents : {},
	servers : []
},
{
	appId: 2
}
];

apps[0].cloudStrategyEvents[now+10000] = cloudStrategyEnum.CHEAPEST;
apps[0].localStrategyEvents[now+10000] = localStrategyEnum.SERVER_LOAD;
apps[0].servers.push(generateServer('http://server1/app', 'cloudA', 1, now, 50));
apps[0].servers.push(generateServer('http://server2/app', 'cloudA', 2, now + 10000, 40));
apps[0].servers.push(generateServer('http://server3/app', 'cloudA', 3, now + 10000, 30));
apps[0].servers.push(generateServer('http://server4/app', 'cloudB', 4, now + 10000, 20));
apps[0].servers.push(generateServer('http://server5/app', 'cloudB', 5, now + 10000, 10));

function generateServer(p_url, p_cloud, p_cost, p_timeStamp, p_load){
	var server = {
		server: p_url,
		cloud: p_cloud,
		cost: p_cost,
		status: {
			cpuLoad: p_load, //Cpu load of the server 0-100
			memLoad: p_load, //Memory load of the server 0-100
			timeStamp: p_timeStamp, //UTC time stamp of this info
			stateEvents: {}
		}
	};

	server.status.stateEvents[p_timeStamp] = stateEnum.READY; //Future state of the serve

	return server;
}

// -----
// TESTS
// -----
function main(){
	for(var appIdx in apps){
		var app = apps[appIdx];
		registerApp(app, registerEnd);
	}
}

// Registering apps
var registered = 0;
function registerApp(app, cbk){
	registered++;
	var data = utils.merge({}, app);

	utils.httpPost('http://'+server_api.host+':'+server_api.port+'/app/'+app.appId, data ,function(status, data){
		if(status == 200){
			console.log('OK: App '+app.appId+' register');
		} else {
			console.log('FAIL: App '+app.appId+' register');
		}
		registered--;
		if(registered === 0){
			cbk();
		}
	});
}

function registerEnd(){
	getAllApps(getAllEnd);
}

// Getting all apps
function getAllApps(cbk){

	utils.httpGet('http://'+server_api.host+':'+server_api.port+'/app', function(status, data){
		if(status === 200){
			data = JSON.parse(data);
			var found = 0;
			for(var appIdx in apps){
				var app = apps[appIdx];

				for(var dataIdx in data){
					var dataApp = data[dataIdx];
					if(app.appId == dataApp.appId){
						found++;
						break;
					}
				}
			}

			if(found == apps.length) {
				console.log('OK: get all apps');
			} else {
				console.log('FAIL: get all apps');
			}

		} else {
			console.log('FAIL: get all apps');
		}

		cbk();
	});
}

function getAllEnd(){
	//updating server status

	var data = utils.merge({}, apps[0]);
	data.servers = [generateServer('http://server2/app', 'cloudA', 1, now + 10000, 30)];

	utils.httpPost('http://'+server_api.host+':'+server_api.port+'/app/'+apps[0].appId, data ,function(status, data){
		if(status == 200){
			console.log('OK: App '+apps[0].appId+' update');
		} else {
			console.log('FAIL: App '+apps[0].appId+' update');
		}
	});
}

main();
