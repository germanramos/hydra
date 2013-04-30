var sys 	= require('sys')
,	exec 	= require('child_process').exec
;

var _servers = [
    {
        "port"      : 3000
    ,   "filePath"  : "./server/main.js"
    }
,
    {
        "port"      : 4001
    ,   "filePath"  : "./server/serviceTime.js"
    }
,
    {
        "port"      : 4002
    ,   "filePath"  : "./server/serviceSum.js"
    }
,
    {
        "port"      : 4003
    ,   "filePath"  : "./server/serviceUuid.js"
    }
];

var _processes = [];
var _startingError = false;

for ( var f=0, F=_servers.length; f<F && !_startingError; f++ ) {
	console.log('starting '+_servers[f].filePath+' on port '+_servers[f].port);
	_processes.push (
		exec(
			'node '+_servers[f].filePath+' '+_servers[f].port
		,
			function (err, stdout, stderr) {
				if (err !== null) {
					_startingError = true;
					console.log('exec error: ' + err);
				}
				else {
					console.log('started');
				}
			}
		)
	);
}

process.on('exit', function () {
  console.log('About to exit.');
});


if ( _startingError ) {
	console.log('an error occurs when starting servers')
	console.log('start killing all processes');
	for ( var f=0, F=_processes.length; f<F; f++ ) {
		console.log('start killing process ['+_processes[f].pid+']' );
		exec(
			'kill '+_processes[f].pid 
		,
			function (err, stdout, stderr) {
				if (err !== null) {
					console.log('WARNING!!! you have to kill manually the process ['+_processes[f].pid+']' );
				}
			}
		);
	}
	console.log('start killing main process');
	process.exit(1);
}




