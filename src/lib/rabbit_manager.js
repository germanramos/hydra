var amqp = require('amqp'),
q		= require('q')
;

var rabbit = function(p_config){
	var self = this;

	this.config = p_config;
	this.connection = null;
	this.exchange = null;
	this.queue = null;

	function createConnection() {
		var host = 'amqp://' + self.config.host + (!!self.config.port && self.config.port.length > 0 ? ':' + self.config.port : '');
		var promise = q.defer();
		if(!self.connection){
			self.connection = amqp.createConnection({url : host});
			self.connection.on('ready', function(){
				promise.resolve();
			});
		} else {
			promise.resolve();
		}
		return promise.promise;
	}

	this.on = function(callback){
		createConnection().then(function(){
			if(!self.exchange){
				self.exchange = self.connection.exchange(self.config.exchange);
			}
			if(!this.queue){
				self.connection.queue(self.config.queue, function(q){
					self.queue = q;
					q.bind(self.exchange, '*');
					q.subscribe(callback);
				});
			}
		});
	};

	this.notify = function(data){
		self.route(self.config.queue, data);
	};

	this.route = function(key, data){
		createConnection().then(function(){
			if(!self.exchange){
				self.exchange = self.connection.exchange(self.config.exchange);
			}
			self.exchange.publish(key, data);
		});
	};

	this.reset = function(){
		/*if (self.queue){
			self.queue.unbind(self.config.queue);
		}
		if (self.exchange){
			self.exchange.unbind(self.config.exchange);
		}*/
		if (self.connection){
			self.connection.end();
			self.queue = null;
			self.exchange = null;
			self.connection = null;
		}
	};
};

module.exports = rabbit;
