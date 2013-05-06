var express = require('express');

if (process.argv.length != 3) {
	console.log("Usage: serviceUuid.js PORT");
	process.exit();
}

PORT = process.argv[2];

function s4() {
	return Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1);
}

function guid() {
  return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
         s4() + '-' + s4() + s4() + s4();
}

function allowCrossDomain(req, res, next) {
	var baseurl = req.get('origin');
	var referer = req.get('Referer');
	res.header('Access-Control-Allow-Origin',baseurl);
	res.header('Access-Control-Allow-Credentials', true);
	next();
};

var app = express();
app.configure(function() {
    app.use(allowCrossDomain);
});
app.get('/service_uuid', function(req, res) {
		//Print log
		console.log("-------------------");
		console.log("service_sum");
		
		var uuid = guid();
		res.send({uuid: uuid});
	});

app.listen(PORT);
console.log('Listening on port ' + PORT);
