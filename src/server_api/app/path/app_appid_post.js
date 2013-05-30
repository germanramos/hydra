var server_api = require('../server_api'),
	hydra = server_api.hydra;

module.exports = function(req, res){
	try{
		var appId = req.params.appId;
		var app = req.body;
		app.appId = appId;

		hydra.app.update(app, function(newApp){
			if(newApp === null){
				res.send(400,'Bad request');
			} else {
				res.send(200,{});
			}
		});
	} catch (ex){
		console.log(ex);
		res.send(400,'Bad request');
	}
};
