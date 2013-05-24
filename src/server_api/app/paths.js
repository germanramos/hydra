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
	{
		"method": "DELETE",
		"path": "/hydra/:url",
		"handler": require("./path/hydra_url_delete")
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
	},
	{
		"method": "DELETE",
		"path": "/app/:appId",
		"handler": require("./path/app_appid_delete")
	},
	{
		"method": "DELETE",
		"path": "/app/:appId/server/:serverUrl",
		"handler": require("./path/app_appid_server_serverurl_delete")
	}
];
