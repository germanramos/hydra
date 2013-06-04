INTERVAL = 5000
HYDRA_URL = "http://localhost:7002/app/time"

window.onload = function() {   
    init_refresh();
    interval = setInterval(init_refresh, INTERVAL);
    
    $("#refreshButton").click(function () {
    	if (this.innerHTML == "Start Refresh") {
    		init_refresh();
    		interval = setInterval(init_refresh, INTERVAL);
    		this.innerHTML = "Stop Refresh"
    		this.style.backgroundColor = "Red";
    	} else {
    		clearInterval(interval);
    		this.innerHTML = "Start Refresh";
    		this.style.backgroundColor = "Green"
    	}
    });
    
    function init_refresh() {
	    $.ajax({
	        url: HYDRA_URL,
	        success: function ( data ) {
	        	console.log("Getted apps from hydra succesfully")
	        	process_app(data)
	        },
	        error: function ( data ) {
	        	console.log("Error when gwtting apps from hydra: " + data);
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
    		process_server(app, server, server_sysmon);
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
    		document.body.appendChild(serverElement);
    	} else {
    		while (serverElement.hasChildNodes()) {
    			serverElement.removeChild(serverElement.lastChild);
    		}
    		create_server(app, serverElement, server, connections);
    	}  	
    }
    
    function create_server(app, serverElement, server, connections) {
    	serverElement.setAttribute('id', server.server);
    	serverElement.setAttribute('class', 'server');
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
}
