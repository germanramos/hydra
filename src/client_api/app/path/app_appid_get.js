var client_api = require('../client_api'),
	hydra = client_api.hydra;

module.exports = function(req, res){
	try{
		var appId = req.params.appId;
		hydra.app.availableServers(appId, function(servers){
			if(servers === null){
				res.send(400,'Bad request');
			} else {
				res.send(200,servers);
			}
		});
	} catch (ex){
		console.log(ex);
		res.send(400,'Bad request');
	}
};
