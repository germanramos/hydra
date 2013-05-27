var commons		= require("../../lib/commons"),
	express		= commons.express,
	hero		= commons.hero,
	server_api	= require('./server_api.js'),
	app			= hero.app;

hero.init(
	require("./paths.js").paths,

	function (){
		server_api.ready(function(err){
			if(err) {
				hero.error(err,'hydra,server_api,start');
			} else {
				app.listen( hero.port() );
				console.log('server_api listening on port '+hero.port() );
			}
		});
	}
);