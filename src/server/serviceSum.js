var express = require('express');

if (process.argv.length != 3) {
	console.log("Usage: serviceSum.js PORT");
	process.exit();
}

PORT = process.argv[2];

var app = express();
app.get('/service_sum/:x/:y', function(req, res) {
		//Print log
		console.log("-------------------");
		console.log("service_sum");
		
		var sum = parseInt(req.params.x) + parseInt(req.params.y);
		res.send({result: sum});
	});

app.listen(PORT);
console.log('Listening on port ' + PORT);
