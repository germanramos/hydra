var hydra = hydra || function () {
	var hydraServers = {
		list : [_GetJSUrl('hydra.js')],
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
		_Initialize();
		_GetApp(appId, override, f_cbk);
	}

	function _Config(p_servers, p_options) {
		p_options = p_options || {};

		hydraServers.list = p_servers;
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
		_Async('GET', hydraServers.list[0] + '/app/hydra',
		function(err, data){
			if(!err) {
				if (data.length > 0) {
					hydraServers.list = data;
					hydraServers.lastUpdate = Date.now();
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
			_Async('GET', hydraServers.list[0] + '/app/'+ appId,
			function(err, data){
				if(!err) {
					// Store the app in the local cache
					appServers[appId] = {
						list: data,
						lastUpdate: Date.now()
					};

					retryTimeout = null;
					f_callback(err, data);
				} else {
					// If the app doesn't exist return the error
					if(err.status === _HTTP_BAD_REQUEST) {
						f_callback(err, null);
					} else {
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
				}
			});
		} else {
			f_callback(null, appServers[appId].list);
		}
	}

	function _CycleHydraServer() {
		var srv = hydraServers.list.shift();
		hydraServers.list.push(srv);
	}

	//////////////////////////
	//    GENERIC UTILS     //
	//////////////////////////
	function _InstanceHttpReq(){
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

		return httpRequest;
	}

	function _Async(p_method, p_url, f_success, data) {
		var req = _InstanceHttpReq();
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

	function _GetJSUrl(file){
		var scripts = document.getElementsByTagName('script');
		for (var i = 0, L = scripts.length; i<L; i++){
			var url = scripts[i].src || '';
			//if(url.indexOf(file) > -1 && (url.indexOf(file) + file.length ===  url.length)){
			if(url.indexOf(file) > -1){
				var fields = url.match( /(.*)[:/]{3}([^:/]+)[:]?([^/]*)([^?]*)[?]?(.*)/ );
				return fields[1] + '://' + fields[2] + (fields[3].length > 0 ? ':' + fields[3] : '');
			}
		}
		return null;
	}

	//////////////////////////////
	//     EXTERNAL METHODS     //
	//////////////////////////////
	return {
		get: _Get,
		config: _Config
	};
}();