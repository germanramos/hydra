var express = require('express');

if (process.argv.length != 3) {
	console.log("Usage: serviceTime.js PORT");
	process.exit();
}

PORT = process.argv[2];

var app = express();
app.get('/service_time', function(req, res) {
		//Print log
		console.log("-------------------");
		console.log("service_time");
		
		var date = new Date();
		res.send({hours:date.getHours(), mins:date.getMinutes(), secs:date.getSeconds()});
	});

app.listen(PORT);
console.log('Listening on port ' + PORT);
