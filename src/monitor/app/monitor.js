var commons = require('../../lib/commons'),
	mongodb = commons.mongodb,
	async	= commons.async,
	ObjectID = commons.mongodb.ObjectID,
	hero	= commons.hero,
	app		= hero.app,
	express	= commons.express,
	hydra	= commons.hydra,
	os		= commons.os;

module.exports = hero.worker (
	function(self){
		var dbHydra = self.db('config', self.config.db);
		var interval = null;

		function _updateStatus(){
			var now = Date.now();
			var server = {
				url: self.config.url,
				status: {
					cpuLoad : os.loadavg()[0] * 100,
					memLoad : os.freemem() * 100 / os.totalmem(),
					timeStamp : now,
					stateEvents : {}
				}
			};
			server.status.stateEvents[now] = hydra.enums.server.stateEnum.READY;

			hydra.server.update(server, function(err){
				console.log('>>>>> update Status\n', server,'\n--------------------------\n');
			});
		}

		function _startUpdateStatus(){
			if(interval === null) {
				return setInterval(_updateStatus, self.config.server.update);
			} else {
				return interval;
			}
		}

		function _stopUpdateStatus(){
			if(interval !== null) {
				clearInterval(interval);
				interval = null;
			}
		}

		self.ready = function(p_cbk){
			async.series (
				[
					// mongoDb
					function(done){
						dbHydra.setup(
							function(err, client){
								hydra.init(client, self.config, done);
							}
						);
					},
					function (done) {
						self.start = _startUpdateStatus;
						self.stop = _stopUpdateStatus;
						_updateStatus();
						interval = _startUpdateStatus();
						done();
					}
				], function(err){
					p_cbk(err);
				}
			);
		};
	}
);

module.exports.hydra = hydra;
