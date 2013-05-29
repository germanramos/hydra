var utils = require('./utils'),
	enums = require('./enums');

var defaultServer = {
	url: null,
	sibling: null
	//	status: {
	//		state: stateEnum.READY, //Current state of the server
	//		cpuLoad: 50, //Cpu load of the server 0-100
	//		memLoad: 50, //Memory load of the server 0-100
	//		timeStamp: 42374897239, //UTC time stamp of this info
	//		stateEvents: [{
	//			state: stateEnum.READY, //Future state of the serve
	//			applyTimeStamp: 42374897239 //UTC time stamp of this info
	//		}]
	//	}
};

module.exports = function(colServer){
	var self = {};
	self.create = function(p_server, p_cbk){
		var server = utils.merge({},defaultServer);
		server = utils.merge(server, p_server);

		//Si no tenemos url no creamos el server
		if(app.url === null){
			p_cbk(null);
			return;
		}

		colServer.insert(server, {w:1}, function(err, items){
			if(err || items.length === 0){
				p_cbk(null);
			} else {
				p_cbk(items[0]);
			}
		});
	};

	self.getAll = function(p_cbk){
		colServer.find({}).toArray(function(err, items){
			p_cbk(items);
		});
	};

	self.getFromUrl = function(p_url, p_cbk){
		var find = {
			url: p_url
		};

		colServer.findOne(find, {}, function(err, item){
			p_cbk(item);
		});
	};

	self.update = function(p_server, p_cbk){
		var find = {
			url: p_app.url
		};

		colServer.update(find, p_server, function(err){
			p_cbk();
		});
	};

	self.remove = function(p_url, p_cbk){
		var find = {
			url: p_url
		};

		colServer.remove(find, function(err, item){
			if(err) {
				p_cbk(err);
			}
			else {
				p_cbk(null);
			}
		});
	};

	return self;
};
