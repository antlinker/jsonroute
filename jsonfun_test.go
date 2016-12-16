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
		analys        = CreateJSONAnalysisHandle()
	)
	BeforeSuite(func() {
		jsonObjData, _ = json.Marshal(model)
		jsonArrayData, _ = json.Marshal(models)
		ok := analys.AddHandle("aa", func(result []interface{}) {
			Ω(result).Should(Equal(models))
		})

		Expect(ok).To(BeTrue())
		ok = analys.AddHandle("bb", func(result *ModelTest) {
			Expect(result.Abc).To(Equal(model.Abc))
			Expect(result.Bcd).To(Equal(model.Bcd))
			Expect(result.Cde).To(BeNumerically("==", model.Cde))
		})
		Expect(ok).To(BeTrue())
		ok = analys.AddHandle("cc", func(result ModelTest) {
			Expect(result.Abc).To(Equal(model.Abc))
			Expect(result.Bcd).To(Equal(model.Bcd))
			Expect(result.Cde).To(BeNumerically("==", model.Cde))
		})
		Expect(ok).To(BeTrue())
		ok = analys.AddHandle("dd", func(result map[string]int32) {
			Expect(result["Abc"]).To(Equal(model.Abc))
			Expect(result["Bcd"]).To(Equal(model.Bcd))
			Expect(result["Cde"]).To(BeNumerically("==", model.Cde))
		})
		Expect(ok).To(BeTrue())
		ok = analys.AddHandle("ff", func(result map[string]int32, data []byte) {

		})
		Expect(ok).To(BeFalse())
		ok = analys.AddHandle("dd", nil)
		Expect(ok).To(BeFalse())
	})

	It("测试结构体", func() {
		err := analys.Exec("aa", jsonArrayData)
		Expect(err).NotTo(HaveOccurred())
		err = analys.Exec("aa", jsonObjData)
		Expect(err).To(HaveOccurred())
		err = analys.Exec("bb", jsonObjData)
		Expect(err).NotTo(HaveOccurred())
		err = analys.Exec("cc", jsonObjData)
		Expect(err).NotTo(HaveOccurred())
		err = analys.Exec("dd", jsonObjData)
		Expect(err).NotTo(HaveOccurred())
		err = analys.Exec("ff", jsonObjData)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("未注册的key:ff"))
		err = analys.Exec("ff", jsonObjData[1:])
		Expect(err).To(HaveOccurred())
		//err = analys.Exec("aaa", jsonArrayData)
		//	Expect(err).NotTo(HaveOccurred())
	})

})
