require('http').globalAgent.maxSockets = 100000;
var request = require("request");

if (process.argv.length != 7) {
	console.log("Usage: node hydraStress.js app_url status_url number wait random");
	console.log("Example: node hydraStress.js http://hydra1.cloud1.com:7001/app/hydra http://hydra1.cloud1.com:7777 2000 10000 true");
	return;
}

var url = process.argv[2];
var statusUrl = process.argv[3];
var number = parseInt(process.argv[4]);
var randomWait = parseInt(process.argv[5]);
var random = process.argv[6] == "true";

console.log(url, statusUrl, number, randomWait, random);

var totalOk = 0;
var totalError = 0;
var totalUnknown = 0;
var totalTime = 0;
var index = 0;

var serverStatus = [];

function makeRequest() {
	//console.log(index);
	index++;
	
	var start = new Date().getTime()
	var options = {
		url : url,
		//agent: false
	}

	request(options, function(error, response, body) {
		//console.log(body);
		if (response && response.statusCode) {
			if (response.statusCode == 200) {
				totalOk++;
				var end = new Date().getTime();
				var time = end - start;
				totalTime += time;
				//console.log('OK, Execution time: ' + time + 'ms');
			}
			else {
				totalError++;
				//console.log("Error:" + response.statusCode);
			} 
		} else {
			totalUnknown++;
			//console.log("Unknow Error");
		}
		if (totalOk + totalError + totalUnknown >= number) {
			console.log('Total Error: ' + totalError);
			console.log('Total Unknown: ' + totalUnknown);
			console.log('Total Ok: ' + totalOk);
			console.log('Total Time: ' + totalTime + 'ms');
			console.log('Average Time: ' + totalTime/number + 'ms');
			
			var memLoad = 0;
			var cpuLoad = 0;
			for (var i=0; i<serverStatus.length; i++) {
				var item = JSON.parse(serverStatus[i]);
				memLoad += item.memLoad;
				cpuLoad += item.cpuLoad;
			}
			console.log('Average Mem Load: ' + Math.round(memLoad/serverStatus.length) + '%');
			console.log('Average Cpu Load: ' + Math.round(cpuLoad/serverStatus.length) + '%');
			process.exit(code=0);
		}	
	});
}

function cpuControl() {
	var options = {
			url : 'http://hydra1.cloud1.com:7777',
			//agent: false
		}
	
	request(options, function(error, response, body) {
		//console.log(body);
		if (response && response.statusCode) {
			if (response.statusCode == 200) {
				serverStatus.push(body);
			}
			else {
				console.log("Status Error:" + response.statusCode);
			} 
		} else {
			console.log("Status Unknow Error");
		}
		setTimeout(cpuControl, 1000);
	});
}

var step = randomWait/number;
for (var i = 0; i < number; i++) {
	if (random)
		setTimeout(makeRequest, Math.floor(Math.random()*randomWait));
	else
		setTimeout(makeRequest, Math.floor(i*step));
}

setTimeout(cpuControl, 1000);
