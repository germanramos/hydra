package api_test

import (
	. "github.com/innotech/hydra/tests/helpers"
)

var s *ServiceTester = NewServiceTester("127.0.0.1:8082", "application")

var _ = s.DefineServiceTests()
