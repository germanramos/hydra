var crypto = require('crypto'), //http://nodejs.org/api/crypto.html
	http = require('http'); //http://nodejs.org/api/http.html

// Pauses streams
exports.pause = require('pause');

// Merge object b into object a.
exports.merge = function(a, b){
	if (a && b) {
		for (var key in b) {
			a[key] = b[key];
		}
	}
	return a;
};

// Random token generator
exports.createToken = function(len){
	var tkn = crypto
	.randomBytes(Math.ceil(len * 3 / 4))
	.toString('base64')
	.slice(0, len)
	.replace(/\//g, '-')
	.replace(/\+/g, '_')
	;
	return tkn;
};

// Cipher a string
exports.cipherString = function(data, key){
	var cipher = crypto.createCipher('aes-256-cbc', key);
	var crypted = cipher.update(data,'utf8','hex');
	crypted += cipher.final('hex');
	return crypted;
};

// Decipher a string
exports.decipherString = function(crypted, key){
	var decipher = crypto.createDecipher('aes-256-cbc', key);
	var data = decipher.update(crypted,'hex','utf8');
	data += decipher.final('utf8');
	return data;
};

function request(method, url, data, cbk){
	var fields = url.match( /(.*)[:/]{3}([^:/]+)[:]?([^/]*)([^?]*)[?]?(.*)/ );
	if(fields === null){
		throw new Error('bar url param');
	}
	var protocol = fields[1];
	var host = fields[2];
	var port = fields[3];
	var path = fields[4];
	var query = fields[5];

	var opt = {
		hostname : host,
		port: port,
		path: path,
		method: method
	};

	var sData = null;
	if(method == 'POST'){
		sData = JSON.stringify(data);
		opt.headers = {
			'Content-Type': 'application/json',
			'Content-Length': sData.length
		};
	}

	var req = http.request(opt, function(res){
		var chunks = [];
		res.setEncoding('utf8');
		res.on('data', function(chunk){
			chunks.push(chunk);
		});
		res.on('end', function(){
			if(cbk !== undefined && cbk !== null){
				cbk(res.statusCode, chunks.join(''));
			}
		});
	});

	if(method == 'POST'){
		req.write(sData);
	}

	req.end();
}

exports.httpGet = function(url, cbk){
	request('GET', url, null, cbk);
};

exports.httpPost = function(url, data, cbk){
	request('POST', url, data, cbk);
};

exports.httpPut = function(url, data, cbk){
	request('PUT', url, data, cbk);
};

exports.httpDelete = function(url, cbk){
	request('DELETE', url, null, cbk);
};
