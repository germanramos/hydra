var commons		= require("../../lib/commons"),
	express		= commons.express,
	hero		= commons.hero,
	monitor		= require('./monitor.js'),
	app			= hero.app;

monitor.ready(function(err){
	if(err) {
		hero.error(err,'hydra,monitor,start');
	} else {
		hero.log('Hydra Monitor Launched', 'hydra,monitor,start');
	}
});