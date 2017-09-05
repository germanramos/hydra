package api_test

import (
	. "github.com/innotech/hydra/tests/helpers"
)

var appBaseId string = "App"

var app1 = map[string]interface{}{
	appBaseId + "1": map[string]interface{}{
		"Cala": map[string]interface{}{
			"Status": "5",
		},
	},
}

var app2 = map[string]interface{}{
	appBaseId + "2": map[string]interface{}{
		"Cloud": "amazon",
		"WWW":   "48.50",
	},
}

var appServiceTester *ServiceTester = NewServiceTester("127.0.0.1:8082", "apps", "app", appBaseId)

var _ = appServiceTester.DefineServiceTests(app1, app2)
