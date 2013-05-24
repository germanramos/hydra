var commons	= require("../../lib/commons"),
	express	= commons.express,
	hero	= commons.hero,
	head	= require('./head.js'),
	app		= hero.app;

hero.init(
	require("./paths.js").paths,

	function (){
		head.ready(function(err){
			if(err) {
				hero.error(err,'hydra,head,start');
			} else {
				app.listen( hero.port() );
				console.log('listening on port '+hero.port() );
			}
		});
	}
);