package api_test

import (
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"fmt"
	"testing"
)

func TestApi(t *testing.T) {
	fmt.Println("Entra en TestApi")
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}
