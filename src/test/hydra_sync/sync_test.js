var utils = require('../../lib/utils'),
	async = require('../../lib/node_modules/async'),
	assert = require('assert');

var clouds = {
	cloud1: {
		client : 'http://hydra1.cloud1.com:7001',
		server : 'http://hydra1.cloud1.com:7002'
	},
	cloud2: {
		client : 'http://hydra1.cloud2.com:7001',
		server : 'http://hydra1.cloud2.com:7002'
	},
	cloud3: {
		client : 'http://hydra1.cloud3.com:7001',
		server : 'http://hydra1.cloud3.com:7002'
	}
};

var app = {
  "localStrategyEvents": {"1370601215953": 0},
  "cloudStrategyEvents": {"1370601215953": 0},
  "servers":[
    {
      "server":"",
      "status": {
        "memLoad": 32.8,
        "cpuLoad": 9.8,
        "stateEvents": {
          "1370350570694": 0
        }
      },
      "cost": 3,
      "cloud": ""
     }
   ]
 };


 function main(){
	async.parallel(
		configHydras,
		function(error, result){
			if(!error) {
				console.log('Waiting for Hydra Server Sync');
				setTimeout(function(){
					console.log('Checking Hydra server sync');
					async.parallel(
						hydraInfo,
						function(error, result){
							if(!error) {
								servers1 = getServers(result.cloud1[0].servers);
								servers2 = getServers(result.cloud2[0].servers);
								servers3 = getServers(result.cloud3[0].servers);

								assert.deepEqual(servers1, servers2);
								assert.deepEqual(servers2, servers3);
								console.log('Servers synced correctly');
							}
						});
				},10000);
			}
		}
	);
 }

function addApp(app, appData, toCloud, f_cbk) {
	console.log('Adding app', app, 'to', toCloud);
	utils.httpPost(toCloud + '/app/' + app, appData,
	function(status, data){
		console.log('Response for ', toCloud,':', status);
		f_cbk(status, data);
	});
}

function getApps(from, f_cbk){
	utils.httpGet(from + '/app/', f_cbk);
}

function getServers(servers) {
	var s = [];
	for(var index in servers){
		s.push(servers[index].server);
	}

	s.sort();
	return s;
}

var configHydras =  {
	cloud1: function(done){
		app.servers[0].server = clouds.cloud2.client;
		app.servers[0].cloud = 'cloud2';
		addApp('hydra', app, clouds.cloud1.server, function(status, data){
			console.log('Response for ', clouds.cloud1.server,':', status);
			done(status === 200 ? null : new Error(data));
		});

	},
	cloud2: function(done){
		app.servers[0].server = clouds.cloud3.client;
		app.servers[0].cloud = 'cloud3';
		addApp('hydra', app, clouds.cloud2.server, function(status, data){
			console.log('Response for ', clouds.cloud2.server,':', status);
			done(status === 200 ? null : new Error(data));
		});
	},
	cloud3: function(done){
		app.servers[0].server = clouds.cloud1.client;
		app.servers[0].cloud = 'cloud1';
		addApp('hydra', app, clouds.cloud3.server, function(status, data){
			console.log('Response for ', clouds.cloud3.server,':', status);
			done(status === 200 ? null : new Error(data));
		});

	}
};

var hydraInfo = {
	cloud1: function(done){
		getApps(clouds.cloud1.server, function(status, data){
			var error = status === 200 ? null : new Error(data);
			var result = status === 200 ? JSON.parse(data) : null;
			done(error, result);
		});
	},
	cloud2: function(done){
		getApps(clouds.cloud2.server, function(status, data){
			var error = status === 200 ? null : new Error(data);
			var result = status === 200 ? JSON.parse(data) : null;
			done(error, result);
		});
	},
	cloud3: function(done){
		getApps(clouds.cloud3.server, function(status, data){
			var error = status === 200 ? null : new Error(data);
			var result = status === 200 ? JSON.parse(data) : null;
			done(error, result);
		});
	}
};

main();