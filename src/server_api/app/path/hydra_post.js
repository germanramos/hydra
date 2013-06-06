var server_api = require('../server_api'),
	hydra = server_api.hydra;

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

		hydra.server.update(server, function(err){
			console.log('updated', item);
			if(err){
				res.send(400,'Bad request');
			} else {
				res.send(200,{});
			}
		});
	} catch (ex){
		console.log(ex);
		res.send(400,'Bad request');
	}

	res.send(400,'Bad request');
};
