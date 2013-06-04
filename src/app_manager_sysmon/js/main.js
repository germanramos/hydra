INTERVAL = 5000;
INIT_HYDRA_URL = "http://localhost:7002/app/time";

var refresh = true;
var interval;

$(document).ajaxStop(function() {
	console.log("Removing obsolete servers");
	//Delete all non refreshed servers
	$(".server").each(function() {
		if (this.getAttribute('checked') == 'false') {
			parent = this.parentNode;
			if (parent.childNodes.length < 6) {
				parent.parentNode.removeChild(parent) //Remove the cloud
			} else {
				parent.removeChild(this); //Remove the server
			}
		}
	});
	//Set up refresh interval
	if (refresh) {
		interval = setInterval(init_refresh, INTERVAL);
	}
});

function init_refresh() {
	clearInterval(interval);
	
	//Mark all server as not refreshed
	$(".server").each(function() {
		this.setAttribute('checked', 'false');
	});
	//Get app info from hydra
    $.ajax({
        url: $("#infoServer").val(),
        success: function ( data ) {
        	console.log("Getted app from hydra succesfully")
        	process_app(data)
        },
        error: function ( data ) {
        	console.log("Error when getting app from hydra: " + data);
        }
    });
}

function process_app(app) {
	servers = app.servers
	for (var i=0; i<servers.length; i++) {
		var server = servers[i];
		console.log("Detected server: " + server.server);
		var splitted = server.server.split(":");
		var server_sysmon = splitted[0] + ":" + splitted[1]+ ":7777/extended"
		for (var key in server.status.stateEvents) {
			if (server.status.stateEvents[key] == 0) {
				process_server(app, server, server_sysmon);
			}
			break;
		}
	}	
}

function process_server(app, server, server_sysmon) {   	
	$.ajax({
        url: server_sysmon,
        success: function ( data ) {
        	console.log("Getted statics from server " + server_sysmon)
        	paint_server(app, server, data.connections);
        },
        error: function ( data ) {
        	console.log("Error when getting static from server " + server_sysmon);
        }
    })
}

function paint_server(app, server, connections) { 
	console.log("Painting server: " + server.server);
	var serverElement = document.getElementById(server.server)
	if (serverElement == null) {
		serverElement = document.createElement("div");
		create_server(app, serverElement, server, connections);
		//Create cloud if needed
		cloud = parseCloud(server.server);
		var cloudElement = document.getElementById(cloud)
		if (cloudElement == null) {
			cloudElement = document.createElement("div");
			cloudElement.setAttribute('id', cloud);
			cloudElement.setAttribute('class', 'cloud');
			divElement = document.createElement("div");
			divElement.setAttribute('class', 'center');
			divElement.appendChild(document.createTextNode(cloud));
			cloudElement.appendChild(divElement);
			document.body.appendChild(cloudElement);
			$(cloudElement).resizable();
			$(cloudElement).draggable();
		}
		//Append server
		cloudElement.appendChild(serverElement);
		//Set an appropiate width to cloud div if it is no manual resized and have more than one server
		if (cloudElement.childNodes.length > 6 && cloudElement.style.width == "" ) {
			cloudElement.style.width = "610px";
		}
	} else {
		while (serverElement.hasChildNodes()) {
			serverElement.removeChild(serverElement.lastChild);
		}
		create_server(app, serverElement, server, connections);
	}  	
}

function parseCloud(url) {
	var found = 0;
	for (var i=url.length-1; i>=0; i--) {
		if (url[i] == '.') {
			found++
			if (found == 2) {
				return url.substring(i+1, url.lastIndexOf(':'));
			}
		}
	}
}

function create_server(app, serverElement, server, connections) {
	serverElement.setAttribute('id', server.server);
	serverElement.setAttribute('class', 'server');
	serverElement.setAttribute('checked', 'true');
	serverElement.appendChild(create_row("ID", app.appId));
	serverElement.appendChild(create_row("URL", server.server));
	serverElement.appendChild(create_row("CPU", server.status.cpuLoad));
	serverElement.appendChild(create_row("MEM", server.status.memLoad));
	filtered = connections.filter(function(element, index, array) {
		return element[5]=="ESTABLISHED"
		});
	serverElement.appendChild(create_row("CON", filtered.length));
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
	
	$("#refreshButton").click(function () {
		if (this.innerHTML == "Start Refresh") {
			init_refresh();
			refresh = true;
			this.innerHTML = "Stop Refresh"
			this.style.backgroundColor = "Red";
		} else {
			clearInterval(interval);
			refresh = false;
			this.innerHTML = "Start Refresh";
			this.style.backgroundColor = "Green"
		}
	});
	
    init_refresh();
    
}
