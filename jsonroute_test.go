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
		jump      = "abc"
		route     = CreateJSONRoute("MT", nil)
		routetest = CreateJSONRoute("MT", jump)
	)
	It("测试AddHandle", func() {
		Expect(routetest.GetRouteKey()).To(Equal("MT"))
		routetest.SetRouteKey("AA")
		Expect(routetest.GetRouteKey()).To(Equal("AA"))
	})
	It("测试JumpObj", func() {
		var modelbb = &ModelRouteTest{"jumpobj", "1111", "222", 33333}
		routetest.AddHandle("jumpobj", func(result *ModelRouteTest, jumpobj string) {
			Expect(result.Abc).To(Equal(modelbb.Abc))
			Expect(result.Bcd).To(Equal(modelbb.Bcd))
			Expect(jumpobj).To(Equal(jump))
			Expect(result.Cde).To(BeNumerically("==", modelbb.Cde))
		})
		jsonObjData, _ := json.Marshal(modelbb)
		routetest.Exec(jsonObjData)
	})
	It("测试AddHandle", func() {

		ok := route.AddHandle("dd", "a")
		Expect(ok).To(BeFalse())
		ok = route.AddHandle("dd", nil)
		Expect(ok).To(BeFalse())
		ok = route.AddHandle("dd", 11)
		Expect(ok).To(BeFalse())
		ok = route.AddHandle("dd", int64(11))
		Expect(ok).To(BeFalse())
		ok = route.AddHandle("dd", true)
		Expect(ok).To(BeFalse())
		ok = route.AddHandle("dd", false)
		Expect(ok).To(BeFalse())
		ok = route.AddHandle("dd", new(ModelRouteTest))
		Expect(ok).To(BeFalse())
		ok = route.AddHandle("dd", make(map[string]interface{}))
		Expect(ok).To(BeFalse())
		ok = route.AddHandle("dd", func() {})
		Expect(ok).To(BeFalse())
		ok = route.AddHandle("dd", func(result ModelRouteTest) {})
		Expect(ok).To(BeTrue())
		ok = route.AddHandle("dd", func(result *ModelRouteTest) {})
		Expect(ok).To(BeTrue())
		ok = route.AddHandle("dd", func(result map[string]interface{}) {})
		Expect(ok).To(BeTrue())
		ok = route.AddHandle("dd", func(result map[string]string) {})
		Expect(ok).To(BeTrue())
	})
	It("测试结构体", func() {

		var modelbb = &ModelRouteTest{"bb", "1111", "222", 33333}
		var modelcc = &ModelRouteTest{"cc", "1111", "222", 33333}
		var modeldd = &ModelRouteTest{"dd", "1111", "222", 33333}
		var modelEE = &ModelRouteTest{"EE", "1111", "222", 33333}
		var modelFF = &ModelRouteTest{"ff", "1111", "222", 33333}
		var modelGg = &ModelRouteTest{"Gg", "1111", "222", 33333}
		var modelHh = &ModelRouteTest{"Hh", "1111", "222", 33333}
		var modelNoreg = &ModelRouteTest{"noreg", "1111", "222", 33333}
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
		route.AddHandle("ee", func(result map[string]interface{}) {
			Expect(result["Abc"]).To(Equal(modelEE.Abc))
			Expect(result["Bcd"]).To(Equal(modelEE.Bcd))
			Expect(result["Cde"]).To(BeNumerically("==", modelEE.Cde))
		})
		route.AddHandle("FF", func(result map[string]interface{}) {
			Expect(result["Abc"]).To(Equal(modelFF.Abc))
			Expect(result["Bcd"]).To(Equal(modelFF.Bcd))
			Expect(result["Cde"]).To(BeNumerically("==", modelFF.Cde))
		})
		route.AddHandle("Gg", func(result map[string]interface{}) {
			Expect(result["Abc"]).To(Equal(modelGg.Abc))
			Expect(result["Bcd"]).To(Equal(modelGg.Bcd))
			Expect(result["Cde"]).To(BeNumerically("==", modelGg.Cde))
		})
		route.AddHandle("hh", func(result map[string]interface{}, data []byte) {
			Expect(result["Abc"]).To(Equal(modelHh.Abc))
			Expect(result["Bcd"]).To(Equal(modelHh.Bcd))
			Expect(result["Cde"]).To(BeNumerically("==", modelHh.Cde))
			// var r ModelRouteTest
			// err := json.Unmarshal(data, &r)
			// Expect(err).NotTo(HaveOccurred())
			// Expect(r).To(MatchJSON(modelHh))
		})

		jsonObjData, _ := json.Marshal(modelbb)

		err := route.Exec(jsonObjData)

		Expect(err).NotTo(HaveOccurred())
		jsonObjDatacc, _ := json.Marshal(modelcc)
		err = route.Exec(jsonObjDatacc)
		Expect(err).NotTo(HaveOccurred())
		jsonObjDatadd, _ := json.Marshal(modeldd)
		err = route.Exec(jsonObjDatadd)
		Expect(err).NotTo(HaveOccurred())
		jsonObjDataee, _ := json.Marshal(modelEE)
		err = route.Exec(jsonObjDataee)
		Expect(err).NotTo(HaveOccurred())
		jsonObjDataff, _ := json.Marshal(modelFF)
		err = route.Exec(jsonObjDataff)
		Expect(err).NotTo(HaveOccurred())
		jsonObjDatagg, _ := json.Marshal(modelGg)
		err = route.Exec(jsonObjDatagg)
		Expect(err).NotTo(HaveOccurred())
		jsonObjDatahh, _ := json.Marshal(modelHh)
		err = route.Exec(jsonObjDatahh)
		Expect(err).NotTo(HaveOccurred())

		jsonObjDataNoreg, _ := json.Marshal(modelNoreg)
		err = route.Exec(jsonObjDataNoreg)
		Expect(err).To(HaveOccurred())

		jsonObjDataErr, _ := json.Marshal(modelNoreg)
		jsonObjDataErr = jsonObjDataErr[1:]
		err = route.Exec(jsonObjDataErr)
		Expect(err).To(HaveOccurred())
	})

})
