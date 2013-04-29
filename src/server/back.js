function bind(app, data) {
	//Add a new service
	app.get('/post_start_service/:service_id', function(req, res) {
		//Print log
		console.log("-------------------");
		console.log("command: post_start_service");
		console.log("service_id:" + req.params.service_id);
		
		//Find service
		var service = data.services.filter(function (service) {
			return service.id == req.params.service_id;
		});
		
		//Reply and add if necessary
		if (service.length > 0) {
			res.send({err:data.err_back.already_started});
		} else {
			services.push({id: req.params.service_id});
			res.send({err:data.err_back.ok});
		}	
	});
	
	//Remove a service
	app.get('/post_stop_service/:service_id', function(req, res) {
		//Print log
		console.log("-------------------");
		console.log("command: post_stop_service");
		console.log("service_id:" + req.params.service_id);
		
		//Find service
		var found = false;
		data.services = data.services.filter(function (service) {
			if (service.id == req.params.service_id) {
				found = true;
				return false;
			} else {
				return true;
			}
		});
		
		//Reply
		if (found) {
			res.send({err:data.err_back.ok});
		} else {
			res.send({err:data.err_back.already_stopped});
		}	
	});
	
	//Return the list of services
	app.get('/get_services', function(req, res) {
		console.log("-------------------");
		console.log("command: get_services");
		res.send({local_services:data.services, remote_known_services:[] });
	});
	
	//Add a new service to previous or new server
	app.get('/post_start_server/:server_name/:service_id', function(req, res) {
		//Print log
		console.log("-------------------");
		console.log("command: post_start_server");
		console.log("server_name:" + req.params.server_name);
		console.log("service_id:" + req.params.service_id);
		
		//Find server
		var server = data.servers.filter(function (server) {
			return server.name == req.params.server_name;
		});
		
		//Reply and add if necessary
		if (server.length <= 0) {
			data.servers.push({name: req.params.server_name, services: [req.params.service_id]});
			res.send({err:data.err_back.ok});
		} else {
			if (server[0].services.indexOf(req.params.service_id) >= 0) {
				res.send({err:data.err_back.already_started});
			} else {
				server[0].services.push(req.params.service_id);
				res.send({err:data.err_back.ok});
			}		
		}	
	});
	
	//Remove a service from server
	app.get('/post_stop_server/:server_name/:service_id', function(req, res) {
		//Print log
		console.log("-------------------");
		console.log("command: post_stop_server");
		console.log("server_name:" + req.params.server_name);
		console.log("service_id:" + req.params.service_id);
		
		//Find server
		var server = data.servers.filter(function (server) {
			return server.name == req.params.server_name;
		});
		
		//Reply and add if necessary
		if (server.length <= 0) {
			res.send({err:data.err_back.already_stopped});
		} else {
			var index = server[0].services.indexOf(req.params.service_id)
			if ( index >= 0) {
				if (server[0].services.length > 1) {
					server[0].services.splice(index,1);
				} else {
					var index = data.servers.indexOf(server[0]);
					data.servers.splice(index, 1)
				}
				res.send({err:data.err_back.ok});
			} else {
				res.send({err:data.err_back.already_stopped});
			}		
		}	
	});
	
	//Return the list of servers
	app.get('/get_servers', function(req, res) {
		console.log("-------------------");
		console.log("command: get_servers");
		res.send({servers:data.servers});
	});
	
	//Return the list of servers for one specific service
	app.get('/get_servers/:service_id', function(req, res) {
		console.log("-------------------");
		console.log("command: get_servers");
		var servers = data.servers.filter(function (server) {
			if (server.services.indexOf(req.params.service_id) >= 0) {
				return true
			}
		});
		res.send({servers:servers});
	});
}

module.exports.bind = bind;