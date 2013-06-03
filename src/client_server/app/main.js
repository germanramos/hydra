var commons		= require("../../lib/commons"),
	express		= commons.express,
	hero		= commons.hero,
	app			= hero.app;

hero.init(require("./paths.js").paths, function(){
	// Configuration
	app.configure(function() {
		app.use( express.static(__dirname + '/../www') );
		app.use(express.bodyParser());
		app.use(app.router);
		app.use(express.errorHandler({
			dumpExceptions : true,
			showStack : true
		}));
	});
	app.listen(hero.port());
	console.log('Client server listening on port', hero.port());
});