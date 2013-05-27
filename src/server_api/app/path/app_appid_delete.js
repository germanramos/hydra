var server_api = require('../server_api'),
	hydra = server_api.hydra;

module.exports = function(req, res){
	var appId = req.params.appId;

	try{
		hydra.app.remove(appId, function(err){
			if(err === null){
				res.send(200,{});
			} else {
				res.send(400, 'Bad request');
			}
		});
	} catch (ex) {
		console.log(ex);
		res.send(400,'Bad request');
	}
};
