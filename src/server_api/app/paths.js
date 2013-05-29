exports.paths = [

	// HYDRA SERVERS

	{
		"method": "GET",
		"path": "/hydra",
		"handler": require("./path/hydra_get")
	},
	{
		"method": "POST",
		"path": "/hydra",
		"handler": require("./path/hydra_post")
	},

	// APPS
	{
		"method": "GET",
		"path": "/app",
		"handler": require("./path/app_get")
	},
	{
		"method": "GET",
		"path": "/app/:appId",
		"handler": require("./path/app_appid_get")
	},
	{
		"method": "POST",
		"path": "/app/:appId",
		"handler": require("./path/app_appid_post")
	}
];
