var assert	= require('assert');
var commons = require('../../lib/commons'),
	hero	= commons.hero,
	hydra	= commons.hydra,
	async	= commons.async
;

describe('Hydra Server', function(){
	var url = "http://test.local.com";
	var h = null;
	var now = Date.now();

	var app = {
		appId: 'test_app',
		localStrategyEvents : {
			//'42374897239' : localStrategyEnum.INDIFFERENT
		},
		cloudStrategyEvents : {
			//'42374897239': cloudStrategyEnum.INDIFFERENT
		},
		servers : [
			//{
			//	server: 'http://server3/app',
			//	cloud : 'nubeA',
			//  cost : 0,
			//	status: {
			//		cpuLoad: 50, //Cpu load of the server 0-100
			//		memLoad: 50, //Memory load of the server 0-100
			//		timeStamp: 42374897239, //UTC time stamp of this info
			//		stateEvents: {
			//			'42374897239' : state: stateEnum.READY, //Future state of the serve
			//		}
			//	}
			//}
		]
	};


	app.localStrategyEvents[now] = 0;
	app.cloudStrategyEvents[now] = 0;

	before(function(done){
		h = hero.worker(function(self){
			var dbHydra = self.db('config', self.config.db);

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
		});

		h.ready(function(err){
			if (!err){
				done();
			}
		});
	});

	before( function(done) {
		hydra.app.remove(url, done);
	});

	afterEach(function(done){
		hydra.app.remove(url, done);
	});

	describe('Add Apps', function(){
		it('Server available with future events', function(done){
			var server = {
				"url": url,
				"status": {
					"memLoad": 32.8,
					"cpuLoad": 9.8,
					"stateEvents": {
					}
				}
			};

			var now = Date.now();
			server.status.stateEvents[now - 20000] = 0;
			server.status.stateEvents[now - 10000] = 0;
			server.status.stateEvents[now - 100] = 0;
			server.status.stateEvents[now + 10000] = 0;
			server.status.stateEvents[now + 20000] = 0;

			app.servers = [server];

			hydra.app.update(server, function(item){
				assert.notEqual(item, null);
				done();
			});
		});

		it('Server with future events and all unavailable ', function(done){
			var server = {
				"url": url,
				"status": {
					"memLoad": 32.8,
					"cpuLoad": 9.8,
					"stateEvents": {
					}
				}
			};

			var now = Date.now();
			server.status.stateEvents[now - 20000] = 1;
			server.status.stateEvents[now - 10000] = 1;
			server.status.stateEvents[now - 100] = 1;
			server.status.stateEvents[now + 10000] = 1;
			server.status.stateEvents[now + 20000] = 1;

			app.servers = [server];

			hydra.app.update(server, function(item){
				assert.equal(item, null);
				done();
			});
		});

		it('Server without future state events but available', function(done){
			var server = {
				"url": url,
				"status": {
					"memLoad": 32.8,
					"cpuLoad": 9.8,
					"stateEvents": {
					}
				}
			};

			var now = Date.now();
			server.status.stateEvents[now - 20000] = 0;
			server.status.stateEvents[now - 10000] = 0;
			server.status.stateEvents[now - 100] = 0;

			app.servers = [server];

			hydra.app.update(server, function(item){
				assert.notEqual(item, null);
				done();
			});
		});

		it('Server without future events and unavailable ', function(done){
			var server = {
				"url": url,
				"status": {
					"memLoad": 32.8,
					"cpuLoad": 9.8,
					"stateEvents": {
					}
				}
			};

			var now = Date.now();
			server.status.stateEvents[now - 20000] = 0;
			server.status.stateEvents[now - 10000] = 0;
			server.status.stateEvents[now - 100] = 1;

			app.servers = [server];

			hydra.app.update(server, function(item){
				assert.equal(item, null);
				done();
			});

		});
	});
});