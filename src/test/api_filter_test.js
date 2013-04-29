var commons	= require('../lib/commons')
,	assert	= commons.assert
,	async	= commons.async
,	hero	= commons.hero
;

describe('Sokobank_filter', function(){
	var queues;

	before(function(done){
		queues = hero.worker(function(self){
			self.feed = self.mq('feed', self.config.mq.feed);
			self.push = self.mq('push', self.config.mq.push);
			self.noPush = self.mq('noPush', self.config.mq.message_nopush);
			self.monitor = self.mq('monitor', self.config.mq.monitor);
		});
		done();
	});	
	
	beforeEach( function(done) {
		done();
	});
	
	describe('Send message to monitor queue', function(){
		it('Use queue to receive message', function(done){
			async.series([
				function(callback){
					queues.push.on(function(data){
						console.log('Push', data);
						done();
					});
					callback(null, 1);
				}
			,
				function(callback){
					queues.noPush.on(function(data){
						console.log('noPush', data);
						done();
			
					});
					callback(null, 1);
				}
			,
				function(callback){
					queues.monitor.on(function(data){
						console.log('Monitor', data);
						done();
					});
					callback(null, 1);
				}
			,
				function(callback){
					queues.feed.on(function(data){
						console.log('feed', data);
						done();
					});
					callback(null, 1);
				}
			], function(err, result){
			});

		});
	});
});