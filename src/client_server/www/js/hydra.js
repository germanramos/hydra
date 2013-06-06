var hydra = hydra || function () {
	var hydraServers = {
		list : ['http://localhost:7001'],
		lastUpdate : 0
	},
		appServers = {
		/* 
		app : {
			list : [],
			lastUpdate : Date.now();
		}
		*/
	},
		updateHydraDelta = 60000, //timeout de cache de hydra servers
		updateAppDelta = 10000, //timeout de cache de app servers
		retryOnFail = 2000,
		overrideCache = false;

	var	_HTTP_STATE_DONE = 0,
		_HTTP_SUCCESS	= 200;

	//////////////////////////
	//     HYDRA  ENTRY     //
	//////////////////////////
	function _get(appId, override, f_cbk) {
		overrideCache = override;
		_GetHydraServers(function(){
			_GetApp(appId, f_cbk);
		});
	}

	//////////////////////////
	//     HYDRA UTILS      //
	//////////////////////////
	function _GetHydraServers(f_callback) {
		if((Date.now() - hydraServers.lastUpdate) > updateHydraDelta ){
			_async('GET', hydraServers.list[0] + '/hydra',
			function(err, data){
				if(!err) {
					hydraServers.list = data;
					hydraServers.lastUpdate = Date.now();

					f_callback();
				} else {
					// In case hydra server doesn't reply, push it to the back 
					// of the list and try another
					var srv = hydraServers.list.shift();
					hydraServers.list.push(srv);

					setTimeout(function() {
						_GetHydraServers(f_callback);
					}, retryOnFail);
				}
			});
		} else {
			f_callback(null);
		}
	}

	function _GetApp(appId, f_callback){
		// Get Apps from server if we specify to override the cache, it's not on the list or the cache is outdated
		var getFromServer = overrideCache || !(appId in appServers) || (Date.now() - appServers[appId].lastUpdate > updateAppDelta);

		if(getFromServer) {
			_async('GET', hydraServers.list[0] + '/app/'+ appId,
			function(err, data){
				if(!err) {
					// Store the app in the local cache
					appServers[appId] = {
						list: data,
						lastUpdate: Date.now()
					};

					f_callback(err, data);
				} else {
					// In case hydra server doesn't reply, push it to the back 
					// of the list and try another
					var srv = hydraServers.list.shift();
					hydraServers.list.push(srv);

					setTimeout(function() {
						_get(appId, overrideCache, f_callback);
					}, retryOnFail);
				}
			});
		} else {
			f_callback(null, appServers[appId].list);
		}
	}

	//////////////////////////
	//    GENERIC UTILS     //
	//////////////////////////
	function _instanceHttpReq(){
		var httpRequest;
		if ( window.XMLHttpRequest ) {
			httpRequest = new XMLHttpRequest();
		}
		else if ( window.ActiveXObject ) {
			try {
				httpRequest = new ActiveXObject('MSXML2.XMLHTTP');
			}
			catch (err1) {
				try {
					httpRequest = new ActiveXObject('Microsoft.XMLHTTP');
				}
				catch (err2) {
					if ( window.console && window.console.error ) {
						console.error('Fatal error', err2);
					}
				}
			}
		}
		if ( !httpRequest ) {
			if ( window.console && window.console.error ) {
				console.error('Fatal error, object httpRequest is not available');
			}
		}

		//httpRequest.overrideMimeType('application/json');
		//httpRequest.withCredentials = true;
		return httpRequest;
	}

	function _async(p_method, p_url, f_success, data) {
		var req = _instanceHttpReq();
		req.open(p_method, p_url+'?_='+(new Date().getTime()), true);
		req.onreadystatechange  = function() {
			if ( req.readyState === 0 || req.readyState === 4 ){
				if (req.status === _HTTP_SUCCESS) {
					if ( req.responseText !== null ) {
						f_success( null, JSON.parse(req.responseText) );
					}
					else {
						f_success(null, null);
					}
				}
				else {
					f_success({ "status" : req.status, req : req },null);
				}
			}
		};

		if(data !== null)
		{
			req.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
			data = JSON.stringify(data);
		}

		req.send(data);
	}

	//////////////////////////////
	//     EXTERNAL METHODS     //
	//////////////////////////////
	return {
		get: _get
	};
}();