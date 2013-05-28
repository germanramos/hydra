var utils = require('./utils'),
enums = require('./enums');

var defaultApp = {
	id: null,
	localStrategy: enums.app.localStrategyEnum.INDIFFERENT,
	cloudStrategy: enums.app.cloudStrategyEnum.INDIFFERENT,
	servers : [
		//{
		//	server: 'http://server3/app',
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
		//}
	]
};

module.exports = function(colApp){
	var self = {};

	self.create = function(p_app, p_cbk){
		var app = utils.merge({},defaultApp);

		app = utils.merge(app, p_app);

		//Si no tenemos id no creamos la app
		if(app.id === null){
			p_cbk(null);
			return;
		}

		colApp.insert(app, {w:1}, function(err, items){
			if(err || items.length === 0){
				p_cbk(null);
			} else {
				p_cbk(items[0]);
			}
		});
	};

	self.getAll = function(p_cbk){
		colApp.find({}).toArray(function(err, items){
			p_cbk(items);
		});
	};

	self.getFromId = function(p_id, p_cbk){
		var find = {
			id: p_id
		};

		colApp.findOne(find, {}, function(err, item){
			p_cbk(item);
		});
	};

	self.update = function(p_app, p_cbk){
		var find = {
			id: p_app.id
		};

		colApp.update(find, p_app, function(err){
			p_cbk();
		});
	};

	self.remove = function(p_id, p_cbk){
		var find = {
			id: p_id
		};

		colApp.remove(find, function(err, item){
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
