var commons = require('../../lib/commons'),
	hero	= commons.hero,
	hydra	= commons.hydra,
	utils	= require('../../lib/utils');

module.exports = new function (){
	var _Servers;
	var _Siblings;
	var _config;

	function _httpGet(p_url, f_done, f_fail){
		console.log("synchronizing with", p_url);
		utils.httpGet( p_url, function(status, data){
			if(status === 200){
				data = JSON.parse(data);
				f_done(data);
			} else {
				console.log('FAIL: get all apps');
				f_fail(status);
			}
		});
	}

	function _getApps(p_url, f_done, f_fail){
		_httpGet(p_url+'/app', f_done, f_fail);
	}

	function _syncServer(p_url){
		_getApps(p_url, _appsDone, _appsFail);
	}

	function _appsDone(p_apps){
		for(var f=0, F=p_apps.length; f<F; f++) {
			hydra.app.update( p_apps[f] );
		}
	}

	function _appsFail(err){
		console.log('hydra_sync:_appsFail', err);
	}

	function _syncDone(p_json){
		p_json = p_json || [];
		_Servers  = p_json.servers || [];

		if (_Servers.length < 1) {
			hero.error("No hydra servers are configured. This hydra instance doesn't connect with other hydras!");
		}

		for ( var f=0, F=_Servers.length; f<F;  f++ ) {
			//ignoring self on sync servers
			if(_Servers[f].server.indexOf(_config.publicUrl) > -1) continue;
			var urlParts = utils.splitUrl(_Servers[f].server);
			var url = urlParts.protocol + '://' + urlParts.host + ':' + _config.serverPort;
			_syncServer(url);
		}
	}

	this.sync = function (config) {
		_config = config;
		hydra.app.getFromId('hydra',
			_syncDone
		);
	};
};
