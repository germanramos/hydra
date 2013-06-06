var assert	= require('assert');
var hydra_server = require('../../lib/dao/server')();
var commons = require('../../lib/commons'),
	hero	= commons.hero,
	hydra	= commons.hydra,
	async	= commons.async
;

describe('Hydra Server', function(){
	var url = "http://test.local.com";
	var h = null;

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
		hydra.server.remove(url, done);
	});

	afterEach(function(done){
		hydra.server.remove(url, done);
	});

	describe('Add Servers', function(){
		it('Server available with future events', function(done){
			var server = {
				"url": url,
				"clientPort": 7001,
				"serverPort" : 7002,
				"sibling": true,
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

			hydra.server.update(server, function(item){
				assert.notEqual(item, null);
				done();
			});
		});

		it('Server with future events and all unavailable ', function(done){
			var server = {
				"url": url,
				"clientPort": 7001,
				"serverPort" : 7002,
				"sibling": true,
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

			hydra.server.update(server, function(item){
				assert.equal(item, null);
				done();
			});
		});

		it('Server without future state events but available', function(done){
			var server = {
				"url": url,
				"clientPort": 7001,
				"serverPort" : 7002,
				"sibling": true,
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

			hydra.server.update(server, function(item){
				assert.notEqual(item, null);
				done();
			});
		});

		it('Server without future events and unavailable ', function(done){
			var server = {
				"url": url,
				"clientPort": 7001,
				"serverPort" : 7002,
				"sibling": true,
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

			hydra.server.update(server, function(item){
				assert.equal(item, null);
				done();
			});

		});
	});
});