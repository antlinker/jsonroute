package jsonroute_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMain(t *testing.T) {

	RegisterFailHandler(Fail)

	RunSpecs(t, "测试")
}
