INTERVAL = 0; //5000;
//INIT_HYDRA_URL = "http://hydra.cloud1.com:7002/app/time";
INIT_HYDRA_URL = "http://2.hydra.innotechapp.com:443/app";
PROBE_PORT = "9099";

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
		var server_sysmon = splitted[0] + ":" + splitted[1] + ":" + PROBE_PORT + "/extended"
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
	if ($("#configShowConnections:checked").length > 0) {
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
	} else {
		paint_server(app, server, {state: 2}, true);
	}
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
			document.getElementById('main').appendChild(cloudElement);
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
	if (data.locked) {
		lockElement = document.createElement("span");
		lockElement.setAttribute("class","lockSymbol");
		serverElement.appendChild(lockElement);
	}
	
	serverElement.setAttribute('id', server.server);
	serverElement.setAttribute('checked', 'true');
	serverElement.appendChild(create_row("ID", app.appId));
	serverElement.appendChild(create_row("URL", server.server));
	serverElement.appendChild(create_row("PRZ", server.cost));
	serverElement.ondblclick = function(e) {
		$(this).remove();
		e.stopPropagation();
		return false;
	}
	$(serverElement).draggable();
	if (alive) {
		var cpuLoad = data.state == 0 ? data.cpuLoad : server.status.cpuLoad;
		var cpuElement = create_row("CPU", Math.round(cpuLoad).toString() + '%', cpuLoad);
		serverElement.appendChild(cpuElement);
		var memLoad = data.state == 0 ? data.memLoad : server.status.memLoad;
		var memElement = create_row("MEM", Math.round(memLoad).toString() + '%', memLoad);
		serverElement.appendChild(memElement);
		var serverClass = 'server app_' + app.appId;
		if (data.state == 0) {
			serverClass += ' active';
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
			//Count all connections of this application including all clouds
			var allConnections = 0;
			$('.server.active.app_' + app.appId + ' .CON .value').each(function() {
				allConnections += parseInt(this.innerHTML);
			});
			
			//Create connections row
			allConnections += filtered.length;
			serverElement.insertBefore(create_row("CON", filtered.length.toString() + ' of ' + allConnections), cpuElement);
			
			//Create connections percent row
			if (allConnections == 0)
				var percent = 0;
			else
				var percent = Math.round(filtered.length * 100 / allConnections);
			serverElement.appendChild(create_row("BAL", percent.toString() + '%', percent));
		} else if (data.state == 2) {
			serverClass += ' active';
		} else {
			serverClass += ' warning';
		}
		serverElement.setAttribute('class', serverClass);
		
	} else {
		serverElement.setAttribute('class', 'server');
	}
	create_menu(serverElement, server, app);
}


function create_menu_action(action, menuElement, server) {
	var actionElement = document.createElement("a");
	actionElement.appendChild(document.createTextNode(action.charAt(0).toUpperCase()));
	actionElement.setAttribute('href', '#');
	actionElement.setAttribute('title', action);
	menuElement.appendChild(actionElement);
	password = $('#password').val();
	
	actionElement.onclick = function() {
		var url_parts = server.server.split(':');
		log("Sending order '" + action + "' to "+ server.server);
		$.ajax({
			type: "GET",
			url : url_parts[0] + ":" + url_parts[1] + ":" + PROBE_PORT + "/" + action + "?password=" + password,
			timeout : 3000,
			success : function(data) {
				log("Succesfull response " + data + " from '" + server.server + "' to order '" + action + "'");
			},
			error : function(data) {
				log("Error response " + data + " from '" + server.server + "' to order '" + action + "'");
			}
		});
	}
}

function create_menu_action_delete(menuElement, server, app) {
	var deleteElement = document.createElement("a");
	deleteElement.appendChild(document.createTextNode("D"));
	deleteElement.setAttribute('title', "delete");
	deleteElement.setAttribute('href', '#');
	menuElement.appendChild(deleteElement);
	
	deleteElement.onclick = function() {
		var answer = window.prompt("Enter delay (ms):", "0");
		if (answer == null)
			return
		var delay = parseInt(answer);
		var now = (new Date).getTime();
		var when = now + delay;	
		var data = {
		    servers: [{
		            server: server.server,
		            status: {
		            	stateEvents: {
		            		//"0": 2
		            	}
		            }
		    }]
		}
		data.servers[0].status.stateEvents[when] = 2;	
		var url_parts = $("#infoServer").val().split('/');
		var url = url_parts[0] + '//' + url_parts[2];
		log("Sending order 'delete' to "+ server.server);
		$.ajax({
			type: "POST",
			url : url + "/app/" + app.appId,
			contentType: 'application/json',
			data : JSON.stringify(data),
			timeout : 3000,
			success : function(data) {
				log("Succesfull response from '" + server.server + "' to order 'delete'");
			},
			error : function(data) {
				log("Error response from '" + server.server + "' to order 'delete'");
			}
		});
	}
}

function create_menu(serverElement, server, app) {	
	var menuElement = document.createElement("div");
	menuElement.setAttribute('class', 'contextmenu');
	serverElement.appendChild(menuElement);
	
	create_menu_action("stress", menuElement, server);
	create_menu_action("halt", menuElement, server)
	create_menu_action("ready", menuElement, server);
	create_menu_action("lock", menuElement, server);
	create_menu_action("unlock", menuElement, server);
	create_menu_action_delete(menuElement, server, app);

}

function create_row(key, value, percent) {
	var pElement = document.createElement("p");
	pElement.setAttribute('class', key);

	var keyElement = document.createElement("span");
	keyElement.setAttribute('class', 'key');
	keyElement.appendChild(document.createTextNode(key));
	pElement.appendChild(keyElement);
	
	if (percent >= 0) {
		var progressBarElement = document.createElement("span");
		progressBarElement.setAttribute('class', 'progressBar');
		var width = percent * 220 / 100;
		progressBarElement.style.width = width.toString() + 'px';
		pElement.appendChild(progressBarElement);
		if (percent >= 15)
			progressBarElement.appendChild(document.createTextNode(value));	
	}
	if (percent == null || percent < 15) {
		var valueElement = document.createElement("span");
		valueElement.setAttribute('class', 'value');
		valueElement.appendChild(document.createTextNode(value));
		pElement.appendChild(valueElement);
		if (percent >= 0)
			valueElement.style.margin = '0px ' + (width+5).toString() + 'px';
	}
		
	return pElement
}

function log(data) {
	var date = new Date();
	var element = document.createElement("p");
	element.appendChild(document.createTextNode(date.toTimeString().split(" ")[0] + ": " + data));
	logElement = $('#log')[0];
	logElement.appendChild(element);
	logElement.scrollTop = logElement.scrollHeight;
	console.info(data);
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
		if (ip == null)
			return
		if (document.getElementById(ip) != null) {
			log("Watcher already exists");
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

		document.getElementById('main').appendChild(watcherElement);
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
				log("Succesfull response " + data + " from '" + server.server + "' to order 'add hydra'");
			},
			error : function(data) {
				log("Error response " + data + " from '" + server.server + "' to order 'add hydra'");
			}
		});
	});
	
	$("#lockButton, #unlockButton").click(function() {
		password = $('#password').val();
		action = this.value.toLowerCase();
		$('.server').each(function() {
			server_url = this.id
			var url_parts = server_url.split(':');	
			
			$.ajax({
				type: "GET",
				url : url_parts[0] + ":" + url_parts[1] + ":" + PROBE_PORT + "/" + action + "?password=" + password,
				timeout : 3000,
				success : function(data) {
					log("Succesfull response " + data + " from '" + server_url + "' to order '" + action + "'");
					console.info(data);
				},
				error : function(data) {
					log("Error response " + data + " from '" + server_url + "' to order '" + action + "'");
				}
			});
		});			
	});
	
	function togleLog() {
		if (this.checked) {
			$('#log').show();
			$('#main').css("bottom","229px");
		} else {
			$('#log').hide();
			$('#main').css("bottom","0");
		}
	}
	
	$("#configShowLog").change(togleLog);
	$("#configShowLog").change();
	
	function togleHelp() {
		if (this.checked) {
			$('.help').show();
			$("#leyendHelp").show();
		} else {
			$('.help').hide();
		}
	}
	$("#configShowHelp").change(togleHelp);
	$("#configShowHelp").change();
	$("#leyendHelp").draggable();
	$("#leyendHelp").dblclick(function() {
		$("#leyendHelp").hide();
	});
	$("#closeHelp").click(function() {
		$("#leyendHelp").hide();
	})

	init_refresh();
}
