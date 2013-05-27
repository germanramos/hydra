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

var apps = [
{
	id: 1,
	localStrategy: localStrategyEnum.INDIFFERENT,
	cloudStrategy: cloudStrategyEnum.ROUND_ROBIN,
	servers : [
	{
		server: 'http://server1/app',
		status: stateEnum.READY
	},
	{
		server: 'http://server2/app',
		status: stateEnum.READY
	},
	{
		server: 'http://server3/app',
		status: stateEnum.UNAVAILABLE
	}
	]
},
{
	id: 2
}
];

// -----
// TESTS
// -----
function main(){
	for(var appIdx in apps){
		var app = apps[appIdx];
		removeApp(app, removeEnd);
	}
}

// Removing previous apps
var removed = 0;
function removeApp(app, cbk){
	removed++;

	utils.httpDelete('http://'+server_api.host+':'+server_api.port+'/app/'+app.id, function(status, data){
		if(status == 200){
			console.log('OK: App '+app.id+' remove');
		} else {
			console.log('FAIL: App '+app.id+' remove');
		}
		removed--;
		if(removed === 0 ) cbk();
	});
}

function removeEnd(){
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
		app : {
			localStrategy: app.localStrategy,
			cloudStrategy: app.cloudStrategy,
			servers: app.servers
		}
	};

	utils.httpPost('http://'+server_api.host+':'+server_api.port+'/app/'+app.id, data ,function(status, data){
		if(status == 200){
			console.log('OK: App '+app.id+' register');
		} else {
			console.log('FAIL: App '+app.id+' register');
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
					if(app.id == dataApp.id){
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
	removeServers();
}

// Removing a server from app
function removeServers(){
	var app = apps[0];
	var serverUrl = encodeURIComponent(app.servers[0].server);
	utils.httpDelete('http://'+server_api.host+':'+server_api.port+'/app/'+app.id+'/server/'+serverUrl, function(status, data){
		if(status == 200) {
			console.log('OK: remove server from app');
		} else {
			console.log('FAIL: remove server from app');
		}
	});
}

main();