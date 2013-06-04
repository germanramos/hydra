exports.paths = [

	{
		"method": "OPTIONS",
		"path": "/hydra",
		"handler": require("./path/hydra_options")
	},

		{
		"method": "OPTIONS",
		"path": "/app/:appId",
		"handler": require("./path/app_appid_options")
	},

	// HYDRA SERVERS
	{
		"method": "GET",
		"path": "/hydra",
		"handler": require("./path/hydra_get")
	},

	// APPS
	{
		"method": "GET",
		"path": "/app/:appId",
		"handler": require("./path/app_appid_get")
	}
];
