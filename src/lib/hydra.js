var commons = require('./commons'),
	mongodb = commons.mongodb,
	utils = require('./utils'),
	enums = require('./enums');

var hydra = module.exports;
hydra.enums = enums;

var colServer = null;
var colApp = null;
var config = {};

hydra.init = function(p_dbClient, p_config, p_cbk){
	config = utils.merge(config, p_config);

	colApp = new mongodb.Collection(p_dbClient, 'app');
	hydra.app = require('./dao/app')(colApp, config);

	p_cbk();
};
