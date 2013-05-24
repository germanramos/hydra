var commons = require('./commons'),
	mongodb = commons.mongodb;

var hydra = module.exports;

var colHydra = null;
var colApp = null;

hydra.init = function(p_dbClient, p_cbk){
	colHydra = new mongodb.Collection(p_dbClient, 'hydra');
	colApp = new mongodb.Collection(p_dbClient, 'app');
	p_cbk();
};

