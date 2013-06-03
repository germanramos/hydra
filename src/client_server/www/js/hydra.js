function hydra(appId, cache, f_cbk) {
	var self = this;
	var hydraServers = {
		list : ['http://localhost:7001'],
		lastUpdate : 0
	};

	var appServers = {
		/* 
		app : {
			list : [],
			lastUpdate : Date.now();
		}

		*/
	};

	var updateHydraDelta = 60000;
	var updateAppDelta = 6000;

	var	_HTTP_STATE_DONE = 0,
	_HTTP_SUCCESS	= 200;

	// Update hydra servers if needed
	//if ((Date.now() - hydraServers.lasUpdate) > updateHydraDelta);


	if((Date.now() - hydraServers.lastUpdate) > updateHydraDelta ){
		// ask for apps again and get the app
		_GetHydraServers(function(err){
			_GetApp(appId, f_cbk);
		});
	} else {
		// get the app servers
		_GetApp(appId, f_cbk);
	}



	//////////////////////////
	//     HYDRA UTILS      //
	//////////////////////////

	function _GetHydraServers(f_callback) {
		_async('GET', 'http://localhost:7001/hydra',
		function(err, data){
			console.log('_GetHydraServers response', err, data);
			f_callback(err);
		});

		/*for(var server in hydraServers.list){
			console.log('hydraServer', hydraServers.list[server]);
		}*/
	}

	function _GetApp(appId, f_callback){
		_async('GET', 'http://localhost:7001/app/'+appId,
		function(err, data){
			console.log('_GetAppServers response', err, data);
			f_callback(err, data);
		});

		/*if(appId in appServers) {
			for (var server in appServers[appId].list){
				console.log('appServer', server);
			}
			f_callback(null, null);
		}
		f_callback(null, null);*/
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
		console.log('_async', arguments);
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
}