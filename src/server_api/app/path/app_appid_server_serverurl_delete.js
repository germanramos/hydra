var server_api = require('../server_api'),
	hydra = server_api.hydra;

module.exports = function(req, res){
	try{
		var appId = req.params.appId;
		var serverUrl = req.params.serverUrl;
		hydra.app.getFromId(appId, function(app){
			if(app === null){
				res.send(400,'Bad request');
			} else {

				var found = false;
				for(var serverIdx in app.servers){
					var server = app.servers[serverIdx];
					if(server.server == serverUrl){
						app.servers.splice(serverIdx,1);
						found = true;
						break;
					}
				}
				if(!found){
					res.send(400,'Bad request');
				} else {
					hydra.app.update(app, function(){
						if(app === null){
							res.send(400,'Bad request');
						} else {
							res.send(200,{});
						}
					});
				}
			}
		});

	} catch (ex){
		console.log(ex);
		res.send(400,'Bad request');
	}
};
