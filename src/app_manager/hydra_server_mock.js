var reqs = []
var http = require('http');

http.createServer(function (req, res) {
  res.writeHead(200, {'Content-Type': 'text/plain'});
  console.log(req.url);
  if (req.method == 'POST') {
        var body = '';
        req.on('data', function (data) {
            body += data;
        });
        req.on('end', function () {
			console.log(body)
        });
  }
  res.end(); 
}).listen(1337, '0.0.0.0');
console.log('Server running at http://0.0.0.0:1337/');