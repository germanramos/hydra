var commons = require('./commons'),
	mongodb = commons.mongodb,
	utils = require('./utils'),
	enums = require('./enums');

var hydra = module.exports;

var colServer = null;
var colApp = null;

hydra.init = function(p_dbClient, p_cbk){
	colApp = new mongodb.Collection(p_dbClient, 'app');
	hydra.app = require('./dao/app')(colApp);

	colServer = new mongodb.Collection(p_dbClient, 'server');
	hydra.server = require('./dao/server')(colServer);

	p_cbk();
};
