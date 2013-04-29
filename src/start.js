var sys 	= require('sys')
,	exec 	= require('child_process').exec
;

var _servers = [
    {
        "port"      : 3000
    ,   "filePath"  : "./sokobank_auth/app/main.js"
    }
,
    {
        "port"      : 3001
    ,   "filePath"  : "./sokobank_api/app/main.js"
    }
,
    {
        "port"      : 3002
    ,   "filePath"  : "./sokobank_userconf/app/main.js"
    }
,
    {
        "port"      : 3003
    ,   "filePath"  : "./sokobank_message_manager/app/main.js"
    }
,
    {
        "port"      : 3004
    ,   "filePath"  : "./sokobank_monitor/app/main.js"
    }
,
    {
        "port"      : 3005
    ,   "filePath"  : "./sokobank_fake_feeder/app/main.js"
    }
,
    {
        "port"      : 3006
    ,   "filePath"  : "./sokobank_pusher/app/main.js"
    }
];

var _env = 'pro';
var _processes = [];
var _startingError = false;

for ( var f=0, F=_servers.length; f<F && !_startingError; f++ ) {
	console.log('starting '+_servers[f].filePath+' on port '+_servers[f].port);
	console.log(' -> node '+_servers[f].filePath+' --port='+_servers[f].port+' --env='+_env );
	_processes.push (
		exec(
			'node '+_servers[f].filePath+' --port='+_servers[f].port+' --env='+_env
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




