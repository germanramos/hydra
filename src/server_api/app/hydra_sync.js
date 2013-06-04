var commons = require('../../lib/commons'),
	hero	= commons.hero,
	hydra	= commons.hydra,
	utils	= require('../../lib/utils');

module.exports = new function (){
	var _Servers;
	var _Siblings;

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

	function _getServer(p_url, f_done, f_fail){
		_httpGet(p_url+'/hydra', f_done, f_fail);
	}

	function _getApps(p_url, f_done, f_fail){
		_httpGet(p_url+'/app', f_done, f_fail);
	}

	function _syncServer(p_url){
		_getServer(p_url, _serverDone, _serverFail);
		_getApps(p_url, _appsDone, _appsFail);
	}

	function _serverDone(p_servers){
		for(var f=0, F=p_servers.length; f<F; f++) {
			if ( _Servers.indexOf(p_servers[f].url) > -1 ){
				p_servers[f].sibling = true;		// Force sibling any server that is currently as sibling
			}
			else {
				p_servers[f].sibling = false;		// Force not sibling any server that is not currently as sibling
			}
			hydra.server.update( p_servers[f] );
		}
	}

	function _serverFail(err){
		console.log('hydra_sync:_serverFail', err);
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
		_Servers  = p_json;
		_Siblings = [];
		if (_Servers.length < 1) {
			hero.error("No hydra servers are configured. This hydra instance doesn't connect with other hydras!");
		}

		for ( var f=0, F=p_json.length; f<F;  f++ ) {
			if ( p_json[f].sibling === true ) {
				_Siblings.push( p_json[f].url );		// Save siblings
			}
			_syncServer(p_json[f].url + ':' + p_json[f].serverPort);
		}
	}

	this.sync = function () {
		hydra.server.getAll(
			_syncDone
		);
	};

};
