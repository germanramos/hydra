//Check that received ids exists
function checkIds(req, res, service_id, consumer_id) {
	//Check consumer id
	var consumer = consumers.filter(function (consumer) {
		return consumer.id == consumer_id
	});
	if (consumer.length <= 0) {
		res.send({err:err_enum.wrong_consumer, srv:null});
	}	

	//Check service id
	var service = services.filter(function (service) {
		return service.id == service_id
	});
	if (service.length <= 0) {
		res.send({err:err_enum.wrong_service, srv:null});
	}
}

//Bind urls
function bind(app, data) {
	//Look into its database for a server who satisfies the request.
	app.get('/post_start/:service_id/:consumer_id', function(req, res) {
		//Print log
		console.log("-------------------");
		console.log("command: post_start");
		console.log("service_id:" + req.params.service_id);
		console.log("consumer_id:" + req.params.consumer_id);
		
		//Check received ids
		checkIds(req, res, req.params.service_id, req.params.consumer_id);
	
		//Find link
		link = active_links.filter(function (link) {
			return link.service_id == req.params.service_id && link.consumer_id == req.params.consumer_id;
		});
			
		//Add link if not exists
		if (link.length <= 0) {
			link = {service_id: req.params.service_id, consumer_id: req.params.consumer_id};
			active_links.push(link);
		}

		//Find service in servers
		var server = data.servers.filter(function (server) {
			return server.services.indexOf(req.params.service_id) >= 0;
		});
		if (server.length > 0) {
			res.send({err:data.err_enum.ok, srv:server[0].name});
		} else {
			res.send({err:data.err_enum.not_here, srv:null});
		}	
	});
	
	//Remove link
	app.get('/post_stop/:service_id/:consumer_id', function(req, res) {
		//Print log
		console.log("------------------");
		console.log("command: post_stop");
		console.log("service_id: " + req.params.service_id);
		console.log("consumer_id: " + req.params.consumer_id);
		
		//Check received ids
		checkIds(req, res, req.params.service_id, req.params.consumer_id);
		
		//Remove from active links
		data.active_links = data.active_links.filter(function (link) {
			return link.service_id != req.params.service_id || link.consumer_id != req.params.consumer_id;
		});
		
		//Reply
		res.send({err:data.err_enum.ok});
	});
	
	//Return the known active siblings of that Hydra server
	app.get('/get_active', function(req, res) {
		console.log("-------------------");
		console.log("command: get_active");
		res.send({active:data.siblings, control:{} });
	});
}

module.exports.bind = bind;