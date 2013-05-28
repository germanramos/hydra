var server_api = require('../server_api'),
	hydra = server_api.hydra;

module.exports = function(req, res){
	var appId = req.params.appId;

	try{
		hydra.app.getFromId(appId, function(item){
			if(item === null){
				res.send(400,'Bad request');
			} else {
				res.send(200,item);
			}
		});
	} catch (ex){
		console.log(ex);
		res.send(400,'Bad request');
	}
};
