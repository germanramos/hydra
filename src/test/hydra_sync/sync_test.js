var utils = require('../../lib/utils');

var clouds = {
	cloud1: 'http://server1.cloud1.com',
	cloud2: 'http://server1.cloud2.com',
	cloud3: 'http://server1.cloud3.com'
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
 }

function addApp(app, appData, toCloud, f_cbk) {
	console.log(toCloud + '/app/' + app, appData);
	utils.httpPost(toCloud + '/app/' + app, appData, f_cbk);
}

function configHydras() {
	app.servers[0].server = clouds.cloud2;
	app.servers[0].cloud = 'cloud2';
	addApp('hydra', app, clouds.cloud1, callback);

	app.servers[0].server = clouds.cloud3;
	app.servers[0].cloud = 'cloud3';
	addApp('hydra', app, clouds.cloud2, callback);

	app.servers[0].server = clouds.cloud1;
	app.servers[0].cloud = 'cloud1';
	addApp('hydra', app, clouds.cloud3, callback);
}

function configApp() {
	app.servers[0].server = clouds.cloud2;
	app.servers[0].cloud = 'cloud2';
	addApp('test', app, clouds.cloud1, callback);
}

function callback(){
	console.log('Response', arguments);
}

main();