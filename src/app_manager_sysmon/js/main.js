window.onload = function() {
    //var hydraServer = window.prompt("Enter an hydra server","http://");
    var hydraServer = "http://localhost:7002"
    var url = hydraServer + "/app"
    $.ajax({
        url: url,
        success: function ( data ) {
        	console.log("Getted apps from hydra succesfully")
        	paint_apps(data)
        },
        error: function ( data ) {
        	console.log("Error when gwtting apps from hydra: " + data);
        }
    })
    
    function paint_apps(apps) {
    	for (var i=0; i<apps.length; i++) {
    		var app = apps[i];
    		console.log("Detected app: " + app.appId);
    		paint_servers(app.servers)
    	}	
    }
    
    function paint_servers(servers) {
    	for (var i=0; i<servers.length; i++) {
    		var server = servers[i];
    		console.log("Detected server: " + server.server);
    		var splitted = server.server.split(":");
    		var server_sysmon = splitted[0] + ":" + splitted[1]+ ":7777/extended"
    		$.ajax({
    	        url: server_sysmon,
    	        success: function ( data ) {
    	        	console.log("Getted statics from server " + server_sysmon)
    	        	paint_server(server, data.connections);
    	        },
    	        error: function ( data ) {
    	        	console.log("Error when getting static from server " + server_sysmon);
    	        }
    	    })
    	}	
    }
    
    function paint_server(server, connections) {
    	//$("body").append("<div>" + server.server + "</div>");
    	var serverElement = document.createElement("div");
    	serverElement.setAttribute('class', 'server') 
    	serverElement.appendChild(create_row("Server", server.server));
    	serverElement.appendChild(create_row("CPU Load", server.status.cpuLoad));
    	serverElement.appendChild(create_row("Memory Load", server.status.memLoad));
    	serverElement.appendChild(create_row("Connections", connections.length));
    	document.body.appendChild(serverElement);
    }
    
    function create_row(key, value) {
    	var pElement = document.createElement("p");
    	
    	var keyElement = document.createElement("span");
    	keyElement.appendChild(document.createTextNode(key));	
    	
    	var valueElement = document.createElement("span");
    	valueElement.appendChild(document.createTextNode(value));
    	
    	pElement.appendChild(keyElement);
    	pElement.appendChild(valueElement);
    	return pElement
    }
}
