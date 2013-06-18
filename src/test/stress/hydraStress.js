require('http').globalAgent.maxSockets = 100000;
var request = require("request");

if (process.argv.length != 5) {
	console.log("Usage: ...");
	return;
}

var url = process.argv[2];
var number = parseInt(process.argv[3]);
var randomWait = parseInt(process.argv[4]);

console.log(url, number, randomWait)

var totalOk = 0;
var totalError = 0;
var totalUnknown = 0;
var totalTime = 0;
var index = 0;

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
			console.log('Average Time: ' + totalTime/number + 'ms')
		}
		
	});

}

for (var i = 0; i < number; i++) {
	setTimeout(makeRequest, Math.floor(Math.random()*randomWait));
}
