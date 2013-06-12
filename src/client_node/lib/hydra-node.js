var request = require('request');

module.exports =  function () {
	var appServers = {
		hydra : {
			list: [],
			lastUpdate : 0
		}
	},
		hydraTimeOut		= 60000,  //timeout de cache de hydra servers
		appTimeOut			= 20000,  //timeout de cache de app servers
		retryOnFail			= 500,
		retryTimeout		= null,
		initialized			= false;

	var	_HTTP_STATE_DONE	= 0,
		_HTTP_SUCCESS		= 200,
		_HTTP_BAD_REQUEST	= 400;

	//////////////////////////
	//     HYDRA  ENTRY     //
	//////////////////////////
	function _Get(appId, override, f_cbk){
		if(!initialized) {
			throw Error('Hydra client not initialized. Use hydra.config([<server list>], {<options>});');
		}

		_GetApp(appId, override, f_cbk);
	}

	function _Config(p_servers, p_options) {
		p_options = p_options || {};

		appServers['hydra'].list = p_servers;

		hydraTimeOut	= (p_options.hydraTimeOut && p_options.hydraTimeOut	> hydraTimeOut ? p_options.hydraTimeOut : hydraTimeOut);
		appTimeOut		= (p_options.appTimeOut   && p_options.appTimeOut   > appTimeOut ? p_options.appTimeOut   : appTimeOut);
		retryOnFail		= (p_options.retryOnFail  && p_options.retryOnFail	> retryOnFail  ? p_options.retryOnFail  : retryOnFail);

		_Initialize();
	}

	//////////////////////////
	//     HYDRA UTILS      //
	//////////////////////////
	function _Initialize(){
		if(initialized) return;

		initialized = true;
		_GetHydraServers();
		setInterval(_GetHydraServers, hydraTimeOut);
	}


	function _GetHydraServers() {
		request.get(appServers['hydra'].list[0] + '/app/hydra',
		function(err, res, data){
			if(!err && res.statusCode === _HTTP_SUCCESS) {
				data = JSON.parse(data);
				if (data.length > 0) {
					appServers['hydra'].list = data;
					appServers['hydra'].lastUpdate = Date.now();
				}

				retryTimeout = null;
			} else {
				// In case hydra server doesn't reply, push it to the back 
				// of the list and try another
				if(!retryTimeout) {
					_CycleHydraServer();
				}

				retryTimeout = setTimeout(function() {
					retryTimeout = null;
					_GetHydraServers();
				}, retryOnFail);
			}
		});
	}

	function _GetApp(appId, overrideCache, f_callback){
		// Get Apps from server if we specify to override the cache, it's not on the list or the list is empty or the cache is outdated
		var getFromServer = overrideCache ||
							!(appId in appServers) ||
							appServers[appId].list.length === 0 ||
							(Date.now() - appServers[appId].lastUpdate > appTimeOut);

		if(getFromServer) {
			request.get(appServers['hydra'].list[0] + '/app/'+ appId,
			function(err, res, data){
				if(!err && res.statusCode === _HTTP_SUCCESS) {
					// Store the app in the local cache
					data = JSON.parse(data);
					appServers[appId] = {
						list: data,
						lastUpdate: Date.now()
					};

					retryTimeout = null;
					f_callback(err, data);
				} else if(!err && res.statusCode === _HTTP_BAD_REQUEST){
					// If the app doesn't exist return the error
					f_callback(new Error(data), null);
				} else if(err) {
					// In case hydra server doesn't reply, push it to the back 
					// of the list and try another
					if(!retryTimeout) {
						_CycleHydraServer();
					}

					retryTimeout = setTimeout(function() {
						retryTimeout = null;
						_Get(appId, overrideCache, f_callback);
					}, retryOnFail);
				}
			});
		} else {
			f_callback(null, appServers[appId].list);
		}
	}

	function _CycleHydraServer() {
		var srv = appServers['hydra'].list.shift();
		appServers['hydra'].list.push(srv);
	}

	//////////////////////////////
	//     EXTERNAL METHODS     //
	//////////////////////////////
	return {
		get: _Get,
		config: _Config
	};
}();