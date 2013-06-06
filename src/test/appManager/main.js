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

apps[0].localStrategyEvents[now+10000] = localStrategyEnum.ROUND_ROBIN;
apps[0].cloudStrategyEvents[now+10000] = cloudStrategyEnum.ROUND_ROBIN;
apps[0].servers.push(generateServer('http://server1/app', now + 10000));
apps[0].servers.push(generateServer('http://server2/app', now + 10000));
apps[0].servers.push(generateServer('http://server3/app', now + 10000));
apps[0].servers.push(generateServer('http://server4/app', now + 10000));

function generateServer(url, timeStamp){
	var server = {
		server: url,
		status: {
			cpuLoad: Math.floor(Math.random()*100), //Cpu load of the server 0-100
			memLoad: Math.floor(Math.random()*100), //Memory load of the server 0-100
			timeStamp: timeStamp, //UTC time stamp of this info
			stateEvents: {}
		}
	};

	server.status.stateEvents[timeStamp] = stateEnum.READY; //Future state of the serve

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
	var data = {
		localStrategyEvents : app.localStrategyEvents,
		cloudStrategyEvents : app.cloudStrategyEvents,
		servers: app.servers
	};

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
	var data = {
		localStrategyEvents : apps[0].localStrategyEvents,
		cloudStrategyEvents : apps[0].cloudStrategyEvents,
		servers: [
			generateServer('http://server2/app', now+15)
		]
	};

	utils.httpPost('http://'+server_api.host+':'+server_api.port+'/app/'+apps[0].appId, data ,function(status, data){
		if(status == 200){
			console.log('OK: App '+apps[0].appId+' update');
		} else {
			console.log('FAIL: App '+apps[0].appId+' update');
		}
	});
}

main();
