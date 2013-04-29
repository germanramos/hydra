var hydr = {};
(function (){

	var _local		= 'localhost',
		_staging	= 'stagging-*.hydr.es',
		_pro		= '*.hydr.es';

	function _hostContains () {
        for(var i = 0 ; i< arguments.length; i++) {
            if ( document.location.host.indexOf(arguments[i]) > -1 ) {
                return true;
            }
        }
        return false;
	}
	hydr.hostContains = _hostContains;

	var _localMode	= _hostContains('local'),
		_debugMode	= false;

	var	_HTTP_STATE_DONE = 0,
		_HTTP_SUCCESS	= 200;

	var _ROOT_REFRESH_TIME = 3600000, // ms.
		_REFRESH_TIME = 60000; // ms.

	var _servers = [],
		_root_timer = null;

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
		httpRequest.overrideMimeType('application/json');
		httpRequest.withCredentials = true;
		return httpRequest;
	}

	function _async(p_method, p_url, f_success) {
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
					console.log('ajax error');
					f_success({ "status" : req.status, req : req },null);
				}
			}
		};
		req.send(null);
	}

	function _merge(p_to, p_from) {
		for ( var p in p_from ){
			if ( typeof p_from[p] === 'object' ){
				if ( !p_to[p] ) {
					p_to[p] = {};
				}
				p_to[p] = _merge(p_to[p], p_from[p]);
			}
			p_to[p] = p_from[p];
		}
		return p_to[p];
	}

	function _array_unique(array) {
	    var a = array.concat();
	    for(var i=0; i<a.length; ++i) {
	        for(var j=i+1; j<a.length; ++j) {
	            if(a[i] === a[j]) {
	                a.splice(j--, 1);
	            }
	        }
	    }
	    return a;
	};


	hydr.util = {
		merge : _merge
	};

	hydr.on		= _on;
	hydr.notify = _notify;

	hydr.httpGet = function (p_url, f_success) {
		_async('GET', p_url, f_success);
	};

	hydr.debug = function (p_debugMode) {
		if ( p_debugMode !== undefined ) {
			_debugMode = p_debugMode;
		}
		return _debugMode;
	};

	hydr.debugMode = function (p_debugMode){
		hydr.debug(p_debugMode);
		return this;
	};

	hydr.local = function (){
		return _localMode;
	};

	hydr._root_refresh = function (){
		var count = hydr._servers.length
		,	responses = [];
		for (var i=0, end=count; i<end; ++i) {
			hydr.httpGet(hydr._servers[i]+'/active/', function(err, response){
				if (!err){ // Error conditions are ignored (console logged in the httpGet call)
					reponses = hydr._arrayUnique(responses.concat(response.active));
				}
				if (!(count -= 1)){ // TODO: Proper control data handling 
					hydr._servers = responses;
				}
			});
		}
		hydr._root_timer = setTimeout(hydr._root_refresh, hydr._ROOT_REFRESH_TIME);
	}; // TODO: Modify ROOT_REFRESH_TIME in case of no responses at all
	hydr._root_refresh();

	function hydr.request = function (service_id, consumer_id){
		// Constructor initialization
		var server = '' // Storage of the found server
		,	status = '' // Storage of the last finding status
		,	_service_id = service_id // Requested service
		,	_consumer_id = consumer_id // Consumer if that service
		,	_cbk = null // Callback to notify further changes of the found server
		,	_timer = null; // Internal timer to verify if the server is still the right one or it has changed

		this._enough_responses_received = function(){
			_cbk(this);
		}

		this._get_server_from_service = function (){ // Gets the proper server to provide a service to a consumer
			var count = hydr._servers.length
			,	responses = [];
			for (var i=0, end=count; i<end; ++i) {
				hydr.httpGet(hydr._servers[i]+
					'/services/'+_service_id+
					'/server?consumer_id='+_consumer_id, function(err, response){ // TODO: Longpolling handling
					if (!err){ // Error conditions are ignored (console logged in the httpGet call)
						responses.push(response);
					}
					if (!(count -= 1)){ // TODO: complete behaviour. It'll return the first right one (to be optimized)
						for (var i=0, end=responses.length, not_found = true; i<end; ++i && not_found){
							if (responses[i].status = null) { // Ok Response
								if (server != reponses[i].server){
									server = reponses[i].server;
									status = reponses[i].status;
									_enough_responses_received();
								} // Non-changing server is not notified
								not_found = false;
							} // TODO: Other status to be processed
						}
					}
				});
			}
			_timer = setTimeout(_get_server_from_service, hydr._REFRESH_TIME);
		};

		this.start = function (cbk){
			_cbk = cbk;
			if (_timer){
				clearInterval(_timer);
			}
			_get_server_from_service();
		};

		this.stop = function (){
			if (_timer){
				clearInterval(_timer);
			}
		};

	};

})();

hydr.client = {};
(function (){

	var hydr.client._running_requests = {};

	function _start(service_id, consumer_id, cbk){
		var compound_id = service_id+consumer_id;
		var req=hydr.client._running_requests[compound_id];
		if (!req){ // If it's not in the already started services, create a new entry, 
			req = new hydr.request(service_id, consumer_id); // Constructor does not invoke anything yet. 
			hydr.client._running_requests[compound_id] = req;
		}
		req.start(cbk); // It will stop the timer if any and will do the actual invocation. Timers will be set when answers (or invocation timeout) will arrive.
		return true;
	}

	function _stop(service_id, consumer_id){
		var compound_id = service_id+consumer_id;
		var req=hydr.client._running_requests[compound_id];
		if (req){ // If it's already started, remove request and stop its timer
			hydr.client._running_requests.splice(hydr.client._running_requests.indexOf(req), 1);
			req.stop();
			return true;
		}
		else { // Otherwise do nothing
			return false;
		}
	}

	hydr.client.start = _start;
	hydr.client.stop = _stop;

})();
