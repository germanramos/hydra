var commons = require('../../lib/commons'),
	mongodb = commons.mongodb,
	async	= commons.async,
	ObjectID = commons.mongodb.ObjectID,
	hero	= commons.hero,
	app		= hero.app,
	express	= commons.express;

module.exports = hero.worker (
	function(self){
		var dbHydra = self.db('config', self.config.db);

		var colServers;

		// Configuration
		app.configure(function() {
			app.use(express.bodyParser());
			app.use(addHeaders);
			app.use(app.router);
			app.use(express.errorHandler({
				dumpExceptions : true,
				showStack : true
			}));
		});

		function addHeaders(req, res, next){
			var allowedOrigins = self.config.api.auth.allowedOrigins;

			var baseurl = req.get('origin');
			var referer = req.get('Referer');
			var i;I=allowedOrigins.length;
			for(i = 0; i < I; i++){
				if((baseurl && baseurl.indexOf(allowedOrigins[i]) !== -1) || (referer && referer.indexOf(allowedOrigins[i]) !== -1)){
					res.header('Access-Control-Allow-Origin',baseurl);
					res.header('Access-Control-Allow-Credentials', true);
				}
			}
			next();
		}

		self.ready = function(p_cbk){
			async.parallel (
				[
					// mongoDb
					function(done){
						dbHydra.setup(
							function(err, client){
								colServers = new mongodb.Collection(client, 'servers');
								done(null);
							}
						);
					}
				], function(err){
					p_cbk(err);
				}
			);
		};

	}
);