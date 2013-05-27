var server_api = require('../server_api'),
	hydra = server_api.hydra;

module.exports = function(req, res){
	var appId = req.params.appId;

	try{
		var app = req.body.app;
		app.id = appId;

		hydra.app.getFromId(appId, function(item){
			if(item === null){
				hydra.app.create(app, function(app){
					if(app === null){
						res.send(400,'Bad request');
					} else {
						res.send(200,{});
					}
				});
			} else {
				hydra.app.update(app, function(app){
					if(app === null){
						res.send(400,'Bad request');
					} else {
						res.send(200,{});
					}
				});
			}
		});

	} catch (ex){
		console.log(ex);
		res.send(400,'Bad request');
	}
};
