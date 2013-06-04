var server_api = require('../server_api'),
	hydra = server_api.hydra;

module.exports = function(req, res){
	try{
		hydra.server.getAll(function(items){
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
