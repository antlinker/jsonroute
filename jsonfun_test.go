package jsonroute_test

import (
	"encoding/json"
	"fmt"

	. "github.com/antlinker/jsonroute"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ModelTest struct {
	Abc int32 `json:"Abc"`
	Bcd int32 `json:"Bcd"`
	Cde int32 `json:"Cde"`
}

func testMapPointfull(result map[string]interface{}) {
	fmt.Printf("testMapPointfull result:%s\n", result)
}

func testObjectPointfull(result *ModelTest) {
	fmt.Printf("testObjectPointfull result:%s\n", result)
}
func testMapfull(result map[string]interface{}) {
	fmt.Printf("testMapfull result:%s\n", result)
}

func testObjectfull(result ModelTest) {
	fmt.Printf("testObjectfull result:%s\n", result)
}

var _ = Describe("测试json解析路由", func() {
	var (
		model         = &ModelTest{1111, 1111, 33333}
		models        = []interface{}{"aaa", "bbb", "cccc", "bbb", "cccc", "cccc"}
		jsonObjData   []byte
		jsonArrayData []byte
		analys        = CreateJsonAnalysisHandle()
	)
	BeforeSuite(func() {
		jsonObjData, _ = json.Marshal(model)
		jsonArrayData, _ = json.Marshal(models)
		analys.AddHandle("aa", func(result []interface{}) {
			Ω(result).Should(Equal(models))
		})
		analys.AddHandle("bb", func(result *ModelTest) {
			Expect(result.Abc).To(Equal(model.Abc))
			Expect(result.Bcd).To(Equal(model.Bcd))
			Expect(result.Cde).To(BeNumerically("==", model.Cde))
		})
		analys.AddHandle("cc", func(result ModelTest) {
			Expect(result.Abc).To(Equal(model.Abc))
			Expect(result.Bcd).To(Equal(model.Bcd))
			Expect(result.Cde).To(BeNumerically("==", model.Cde))
		})
		analys.AddHandle("dd", func(result map[string]int32) {
			Expect(result["Abc"]).To(Equal(model.Abc))
			Expect(result["Bcd"]).To(Equal(model.Bcd))
			Expect(result["Cde"]).To(BeNumerically("==", model.Cde))
		})
	})

	It("测试结构体", func() {
		analys.Exec("aa", jsonArrayData)
		analys.Exec("bb", jsonObjData)
		analys.Exec("cc", jsonObjData)
		analys.Exec("dd", jsonObjData)
	})

})
