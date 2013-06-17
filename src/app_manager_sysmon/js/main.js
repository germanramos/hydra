INTERVAL = 0; //5000;
INIT_HYDRA_URL = "http://hydra1.cloud1.com:7002/app";

var refresh = true;
var interval;
var watch = [];

$(document).ajaxStop(function() {
	console.log("Removing obsolete servers");
	// Delete all non refreshed servers
	$(".server").each(function() {
		if (this.getAttribute('checked') == 'false') {
			parent = this.parentNode;
			if ($(parent).find(".server").length < 2) {
				parent.parentNode.removeChild(parent) // Remove the cloud
			} else {
				parent.removeChild(this); // Remove the server
			}
		}
	});
	// Update all non refreshed watchers
	$(".watcher").each(function() {
		if (this.getAttribute('checked') == 'false') {
			$(this).find(".where").html("unknown");
		}
	});
	// Set up refresh interval
	if (refresh) {
		interval = setInterval(init_refresh, INTERVAL);
		
	}
});

function init_refresh() {
	clearInterval(interval);

	// Mark all server as not refreshed
	$(".server").each(function() {
		this.setAttribute('checked', 'false');
	});
	
	// Mark all watchers as not refreshed
	$(".watcher").each(function() {
		this.setAttribute('checked', 'false');
	});
	
	// Get app info from hydra
	$.ajax({
		url : $("#infoServer").val(),
		timeout : 3000,
		success : function(data) {
			console.log("Got app from hydra succesfully")
			data = [].concat(data);
			for (var i=0; i < data.length; i++) {
				process_app(data[i]);
			}
		},
		error : function(data) {
			console.log("Error when getting app from hydra: " + data);
		}
	});
}

function process_app(app) {
	servers = app.servers
	for ( var i = 0; i < servers.length; i++) {
		var server = servers[i];
		console.log("Detected server: " + server.server);
		var splitted = server.server.split(":");
		var server_sysmon = splitted[0] + ":" + splitted[1] + ":7777/extended"
		for ( var key in server.status.stateEvents) {
			if (server.status.stateEvents[key] == 0) {
				process_server(app, server, server_sysmon);
			} else if ($("#configShowUnavailable:checked").length > 0) {
				paint_server(app, server, [], false);
			}
			break;
		}
	}
}

function process_server(app, server, server_sysmon) {
	$.ajax({
		url : server_sysmon,
		timeout : 3000,
		success : function(data) {
			console.log("Getted statics from server " + server_sysmon);
			paint_server(app, server, data, true);
		},
		error : function(data) {
			console.log("Error when getting static from server "
					+ server_sysmon);
			paint_server(app, server, {state: 1}, true);
		}
	})
}

function paint_server(app, server, data, alive) {
	console.log("Painting server: " + server.server);
	var serverElement = document.getElementById(server.server)
	if (serverElement == null) {
		serverElement = document.createElement("div");
		create_server(app, serverElement, server, data, alive);
		// Create cloud if needed
		cloud = server.cloud;
		var cloudElement = document.getElementById(cloud)
		if (cloudElement == null) {
			cloudElement = document.createElement("div");
			cloudElement.setAttribute('id', cloud);
			cloudElement.setAttribute('class', 'cloud');
			cloudElement.ondblclick = function() {
				$(this).remove();
			}
			divElement = document.createElement("div");
			divElement.setAttribute('class', 'title');
			divElement.appendChild(document.createTextNode(cloud));
			cloudElement.appendChild(divElement);
			document.body.appendChild(cloudElement);
			//document.getElementById("main").appendChild(cloudElement);
			$(cloudElement).resizable();
			$(cloudElement).draggable();
		}
		// Append server
		cloudElement.appendChild(serverElement);
		// Set an appropiate width to cloud div if it is no manual resized and
		// have more than one server
		if ($(cloudElement).find(".server").length > 1
				&& cloudElement.style.width == "") {
			cloudElement.style.width = "605px";
		}
	} else {
		while (serverElement.hasChildNodes()) {
			serverElement.removeChild(serverElement.lastChild);
		}
		create_server(app, serverElement, server, data, alive);
	}
}

function create_server(app, serverElement, server, data, alive) {
	serverElement.setAttribute('id', server.server);
	serverElement.setAttribute('checked', 'true');
	serverElement.appendChild(create_row("ID", app.appId));
	serverElement.appendChild(create_row("PRZ", server.cost));
	serverElement.appendChild(create_row("URL", server.server));
	serverElement.ondblclick = function(e) {
		$(this).remove();
		e.stopPropagation();
		return false;
	}
	$(serverElement).draggable();
	if (alive) {		
		serverElement.appendChild(create_row("CPU", server.status.cpuLoad));
		serverElement.appendChild(create_row("MEM", server.status.memLoad));
		if (data.state == 0) {
			serverElement.setAttribute('class', 'server active');
			filtered = data.connections.filter(function(element, index, array) {
				if (element[5] == "ESTABLISHED") {
					var ip = element[4][0];
					if ($.inArray(ip, watch) >= 0) {
						$('#' + ip.replace(/\./g, "\\.") + " .where").each(
								function() {
									this.innerHTML = server.server;
									this.parentNode.setAttribute('checked', 'true');
								});
					}
					return true;
				} else {
					return false;
				}
			});
			serverElement.appendChild(create_row("CON", filtered.length));
		} else {
			serverElement.setAttribute('class', 'server warning');
		}
		
		
	} else {
		serverElement.setAttribute('class', 'server');
	}
}

function create_row(key, value) {
	var pElement = document.createElement("p");

	var keyElement = document.createElement("span");
	keyElement.setAttribute('class', 'key');
	keyElement.appendChild(document.createTextNode(key));

	var valueElement = document.createElement("span");
	valueElement.setAttribute('class', 'value');
	valueElement.appendChild(document.createTextNode(value));

	pElement.appendChild(keyElement);
	pElement.appendChild(valueElement);
	return pElement
}

window.onload = function() {
	$("#infoServer").val(INIT_HYDRA_URL);
	$("#title").html("Hydra System Monitor");

	$("#refreshButton").click(function() {
		if (this.value == "Start Refresh") {
			init_refresh();
			refresh = true;
			this.value = "Stop Refresh"
		} else {
			clearInterval(interval);
			refresh = false;
			this.value = "Start Refresh";
		}
	});

	$("#addWatcherButton").click(function() {
		var ip = window.prompt("Enter client ip:", "127.0.0.1");
		if (document.getElementById(ip) != null) {
			alert("Watcher already exists");
			return;
		}

		watch.push(ip);
		watcherElement = document.createElement("div");
		watcherElement.setAttribute('id', ip);
		watcherElement.setAttribute('class', 'watcher');
		watcherElement.setAttribute('checked', 'true');

		var ipElement = document.createElement("span");
		ipElement.setAttribute('class', 'ip');
		ipElement.appendChild(document.createTextNode(ip));
		watcherElement.appendChild(ipElement);

		var whereElement = document.createElement("span");
		whereElement.setAttribute('class', 'where');
		whereElement.appendChild(document.createTextNode("unknown"));
		watcherElement.appendChild(whereElement);

		watcherElement.ondblclick = function() {
			var index = watch.indexOf(ip)
			watch.splice(index, 1);
			$(this).remove();
		}

		document.body.appendChild(watcherElement);
		//document.getElementById("main").appendChild(watcherElement);
		$(watcherElement).draggable();
	});
	
	$("#addHydraButton").click(function() {
		var ip = window.prompt("Enter new hydra url:", "http://");
		if (ip == null)
			return;
		
		var data = {
		    servers: [
		        {
		            server: ip
		        }
		    ]
		}
		
		var url_parts = $("#infoServer").val().split('/');
		var url = url_parts[0] + '//' + url_parts[2];
		
		$.ajax({
			type: "POST",
			url : url + "/app/hydra",
			contentType: 'application/json',
			data : JSON.stringify(data),
			timeout : 3000,
			success : function(data) {
				//alert("OK");
			},
			error : function(data) {
				alert("Error:" + data);
			}
		});
	});

	init_refresh();
}
