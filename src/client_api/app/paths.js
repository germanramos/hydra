exports.paths = [
	// APPS
	{
		"method": "GET",
		"path": "/app/:appId",
		"handler": require("./path/app_appid_get")
	}
];
