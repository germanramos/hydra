var commons = require('../../lib/commons'),
	mongodb = commons.mongodb,
	async	= commons.async,
	ObjectID = commons.mongodb.ObjectID,
	hero	= commons.hero,
	app		= hero.app,
	express	= commons.express,
	hydra	= commons.hydra,
	hydra_sync = require('./hydra_sync.js'),
	utils	= commons.utils;

module.exports = hero.worker (
	function(self){
		var dbHydra = self.db('config', self.config.db);

		var colServers;

		// Configuration
		app.configure(function() {
			app.use(express.bodyParser());
			app.use(utils.addHeaders(self.config.server_api.allowedOrigins));
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

					if ( err === null ) {
						// Start to sync hydra
						hydraSync();
						setInterval( hydraSync, self.config.server.sync );
					}

					p_cbk(err);
				}
			);
		};

		function hydraSync(){
			hydra_sync.sync(self.config);
		}

	}
);

module.exports.hydra = hydra;
