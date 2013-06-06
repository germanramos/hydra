var commons = require('../../lib/commons'),
	mongodb = commons.mongodb,
	async	= commons.async,
	ObjectID = commons.mongodb.ObjectID,
	hero	= commons.hero,
	app		= hero.app,
	express	= commons.express,
	hydra	= commons.hydra
	utils	= commons.utils;

module.exports = hero.worker (
	function(self){
		var dbHydra = self.db('config', self.config.db);

		var colServers;

		// Configuration
		app.configure(function() {
			app.use(express.bodyParser());
			app.use(utils.addHeaders(self.config.client_api.allowedOrigins));
			app.use(app.router);
			app.use(express.errorHandler({
				dumpExceptions : true,
				showStack : true
			}));
		});

		self.ready = function(p_cbk){
			async.parallel (
				[
					// mongoDb
					function(done){
						dbHydra.setup(
							function(err, client){
								hydra.init(client, self.config, done);
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

module.exports.hydra = hydra;
