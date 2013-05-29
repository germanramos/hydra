var client_api = require('../client_api'),
	hydra = client_api.hydra;

module.exports = function(req, res){

	try{
		hydra.server.getAll(appId, function(items){
			if(items === null){
				res.send(400,'Bad request');
			} else {
				res.send(200,items);
			}
		});
	} catch(ex) {
		console.log(ex);
		res.send(400,'Bad request');
	}

};