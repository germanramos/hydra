var server_api = require('../server_api'),
	hydra = server_api.hydra;

function _statusStructOk(status) {
	if(	typeof status === 'object'&&
			'cpuLoad'		in status &&
			'memLoad'		in status &&
			'timeStamp'		in status &&
			'stateEvents'	in status &&
			typeof status.stateEvents === 'object' ) {
		return true;
	} else {
		return false;
	}
}

module.exports = function(req, res){
	try{
		var server = {
			url : req.body.url,
			sibling : req.body.sibling,
			status : req.body.status,
			clientPort : req.body.clientPort,
			serverPort : req.body.serverPort,
			cloud : req.body.cloud,
			cost : req.body.cost
		};

		if(_statusStructOk(server.status)) {
			hydra.server.update(server, function(err){
				if(err){
					res.send(400,'Bad request');
				} else {
					res.send(200,{});
				}
			});

		} else {
			res.send(400, 'Missing status parameters');
		}

	} catch (ex){
		console.log(ex);
		res.send(400,'Bad request');
	}
};
