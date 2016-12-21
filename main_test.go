package jsonroute_test

import (
	"testing"

	"github.com/antlinker/alog/log"

	"github.com/antlinker/alog"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMain(t *testing.T) {
	alog.SetEnabled(false)
	alog.GALog.SetLogLevel(log.INFO)
	RegisterFailHandler(Fail)

	RunSpecs(t, "测试")
}
