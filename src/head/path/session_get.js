var head = require('../head');

function handler(req, res){
	head.setSession(req, res, req.param('s'));
}

module.exports = handler;
