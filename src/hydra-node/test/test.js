var hydra = require('../lib/hydra-node');

hydra.config(['http://localhost:7001']);

setInterval(function(){
	hydra.get('hydra', false, function(err, data){
		console.log('>>>> Error', err);
		console.log('>>>> Data', data);
	});
}, 2000);