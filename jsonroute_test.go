package jsonroute_test

import (
	"encoding/json"

	. "github.com/antlinker/jsonroute"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ModelRouteTest struct {
	Mt  string `json:"MT"`
	Abc string `json:"Abc"`
	Bcd string `json:"Bcd"`
	Cde int32  `json:"Cde"`
}

var _ = Describe("测试jsonroute解析路由测试", func() {
	var (
		route = CreateJsonRoute("MT", nil)
	)

	It("测试结构体", func() {

		var modelbb = &ModelRouteTest{"bb", "1111", "222", 33333}
		var modelcc = &ModelRouteTest{"cc", "1111", "222", 33333}
		var modeldd = &ModelRouteTest{"dd", "1111", "222", 33333}

		route.AddHandle("bb", func(result *ModelRouteTest) {
			Expect(result.Abc).To(Equal(modelbb.Abc))
			Expect(result.Bcd).To(Equal(modelbb.Bcd))
			Expect(result.Cde).To(BeNumerically("==", modelbb.Cde))
		})
		route.AddHandle("cc", func(result ModelRouteTest) {
			Expect(result.Abc).To(Equal(modelcc.Abc))
			Expect(result.Bcd).To(Equal(modelcc.Bcd))
			Expect(result.Cde).To(BeNumerically("==", modelcc.Cde))
		})
		route.AddHandle("dd", func(result map[string]interface{}) {
			Expect(result["Abc"]).To(Equal(modeldd.Abc))
			Expect(result["Bcd"]).To(Equal(modeldd.Bcd))
			Expect(result["Cde"]).To(BeNumerically("==", modeldd.Cde))
		})

		jsonObjData, _ := json.Marshal(modelbb)

		route.Exec(jsonObjData)
		jsonObjDatacc, _ := json.Marshal(modelcc)
		route.Exec(jsonObjDatacc)
		jsonObjDatadd, _ := json.Marshal(modeldd)
		route.Exec(jsonObjDatadd)
	})

})
