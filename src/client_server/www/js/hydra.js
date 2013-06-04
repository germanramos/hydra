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

function hydra(appId, overrideCache, f_cbk) {

	_GetHydraServers(function(err){
		if(!err) {
			_GetApp(appId, f_cbk);
		} else {
			console.log('Hydra head have been chopped off!');
		}
	});

	//////////////////////////
	//     HYDRA UTILS      //
	//////////////////////////
	function _GetHydraServers(f_callback) {
		if((Date.now() - hydraServers.lastUpdate) > updateHydraDelta ){
			_async('GET', hydraServers.list[0] + '/hydra',
			function(err, data){
				console.log('_GetHydraServers response', err, data);
				if(!err) {
					hydraServers.list = data;
					hydraServers.lastUpdate = Date.now();
					console.log('Hydra Servers', hydraServers);
					f_callback(null);
				} else {
					f_callback(err);
				}
			});
		} else {
			console.log('_GetHydraServers cache response', null, hydraServers.list);
			f_callback(null);
		}
	}

	function _GetApp(appId, f_callback){
		if(overrideCache || !(appId in appServers) || (Date.now() - appServers[appId].lastUpdate > updateAppDelta)) {
			_async('GET', hydraServers.list[0] + '/app/'+ appId,
			function(err, data){
				console.log('_GetAppServers response', err, data);
				if(!err) {
					appServers[appId] = {
						list: data,
						lastUpdate: Date.now()
					};

					f_callback(err, data);
				} else {
					var srv = hydraServers.list.shift();
					hydraServers.list.push(srv);
					setTimeout(function() {
						_GetApp(appId, f_callback);
					}, 2000);
				}
			});
		} else {
			console.log('_GetAppServers cache response', null, appServers[appId].list);
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