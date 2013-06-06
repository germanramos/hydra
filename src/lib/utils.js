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

exports.sortObj = function(obj){
	var kvs = [];
	for(var key in obj){
		kvs.push({k:key,v:obj[key]});
	}

	kvs.sort(function(a,b){
		return a.k.localeCompare(b.k);
	});

	var ret = {};
	for(var idx in kvs){
		ret[kvs[idx].k] = kvs[idx].v;
	}
	return ret;
};

exports.addHeaders = function addHeaders(allowedOrigins){
	return function (req, res, next){
		var baseurl = req.get('origin');
		var referer = req.get('Referer');
		var i;I=allowedOrigins.length;
		if(allowedOrigins[0] == "*"){
			res.header('Access-Control-Allow-Origin',baseurl);
			res.header('Access-Control-Allow-Headers','content-type');
			//res.header('Access-Control-Allow-Credentials', true);
		} else {
			for(i = 0; i < I; i++){
				if((baseurl && baseurl.indexOf(allowedOrigins[i]) !== -1) || (referer && referer.indexOf(allowedOrigins[i]) !== -1)){
					res.header('Access-Control-Allow-Origin',baseurl);
					res.header('Access-Control-Allow-Credentials', true);
				}
			}
		}

		if(req.method == 'OPTIONS'){
			res.send(200,{});
		} else {
			next();
		}
	};
};

// Count the number of items on an array
// or keys in an object
exports.count = function (item) {
	if('length' in item){
		return item.length;
	}

	var c = 0;
	for(var key in item){
		c++;
	}
	return c;
};

// Return the keys in an object
exports.keys = function (item) {
	var k = [];
	for(var key in item){
		k.push(key);
	}
	return k;
};