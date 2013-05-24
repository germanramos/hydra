var commons		= require("../../lib/commons"),
	express		= commons.express,
	hero		= commons.hero,
	client_api	= require('./client_api.js'),
	app			= hero.app;

hero.init(
	require("./paths.js").paths,

	function (){
		client_api.ready(function(err){
			if(err) {
				hero.error(err,'hydra,client_api,start');
			} else {
				app.listen( hero.port() );
				console.log('listening on port '+hero.port() );
			}
		});
	}
);