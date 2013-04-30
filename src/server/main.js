var express = require('express');
var data = require('./data');
var front = require('./front');
var back = require('./back');

var allowCrossDomain = function(req, res, next) {
	var baseurl = req.get('origin');
	var referer = req.get('Referer');
	res.header('Access-Control-Allow-Origin',baseurl);
	res.header('Access-Control-Allow-Credentials', true);
	next();
};

if (process.argv.length != 3) {
	console.log("Usage: main.js PORT");
	process.exit();
}

PORT = process.argv[2];

var app = express();

app.configure(function() {
    app.use(allowCrossDomain);
    //app.use(express.static(__dirname + '/public'));
});

front.bind(app, data);
back.bind(app, data);

app.listen(PORT);
console.log('Listening on port ' + PORT);

/* Examples:
http://127.0.0.1:3000/post_start/Service1/Pepe
http://127.0.0.1:3000/post_stop/Service1/Pepe
http://127.0.0.1:3000/get_active
http://127.0.0.1:3000/post_start_service/Service5
http://127.0.0.1:3000/post_stop_service/Service5
http://127.0.0.1:3000/get_services
http://127.0.0.1:3000/post_start_server/Server4/Service7
http://127.0.0.1:3000/post_stop_server/Server4/Service7
http://127.0.0.1:3000/get_servers
http://127.0.0.1:3000/get_servers/Service1
*/