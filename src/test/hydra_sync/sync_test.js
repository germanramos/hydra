var utils = require('../../lib/utils'),
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
	configHydras();
	setTimeout(function(){
		getApps(clouds.cloud1.server, function(status, data1){
			data1 = JSON.parse(data1);
			getApps(clouds.cloud2.server, function(status, data2){
				data2 = JSON.parse(data2);
				getApps(clouds.cloud3.server, function(status, data3){
					data3 = JSON.parse(data3);
					servers1 = getServers(data1[0].servers);
					servers2 = getServers(data2[0].servers);
					servers3 = getServers(data3[0].servers);

					assert.deepEqual(servers1, servers2);
					assert.deepEqual(servers2, servers3);
				});
			});
		});
	},10000);
 }

function addApp(app, appData, toCloud, f_cbk) {
	console.log(toCloud + '/app/' + app, appData);
	utils.httpPost(toCloud + '/app/' + app, appData, f_cbk);
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

function configHydras() {
	app.servers[0].server = clouds.cloud1.client;
	app.servers[0].cloud = 'cloud2';
	addApp('hydra', app, clouds.cloud1.server, callback);

	app.servers[0].server = clouds.cloud3.client;
	app.servers[0].cloud = 'cloud3';
	addApp('hydra', app, clouds.cloud2.server, callback);

	app.servers[0].server = clouds.cloud1.client;
	app.servers[0].cloud = 'cloud1';
	addApp('hydra', app, clouds.cloud3.server, callback);
}

function callback(status, data){
	assert.equal(status, 200);
}

main();