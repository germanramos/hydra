var server_api = require('../server_api'),
	hydra = server_api.hydra;

module.exports = function(req, res){
	try{
		var app = {
			appId : req.params.appId,
			localStrategyEvents : req.body.localStrategyEvents || {},
			cloudStrategyEvents : req.body.cloudStrategyEvents || {},
			servers : req.body.servers || []
		};

		hydra.app.update(app, function(err){
			console.log(err);
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
};
