package api_test

import (
	. "github.com/innotech/hydra/tests/helpers"
)

var instanceBaseId string = "Instance"

var instance1 = map[string]interface{}{
	instanceBaseId + "1": map[string]interface{}{
		"cpu": "75.89",
		"mem": "44.21",
	},
}

var instance2 = map[string]interface{}{
	instanceBaseId + "2": map[string]interface{}{
		"cpu": "23.77",
		"mem": "65.33",
	},
}

var instanceServiceTester *ServiceTester = NewServiceTester("127.0.0.1:8082", "apps/App1/instances", "instance", instanceBaseId)

var _ = instanceServiceTester.DefineServiceTests(instance1, instance2)
