package jsonroute

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/antlinker/alog"
)

func init() {
	alog.RegisterAlog()
}

func callfun(funValue reflect.Value, jumpobj reflect.Value, param0 reflect.Value, ortherparam ...interface{}) error {
	funType := funValue.Type()
	num := funType.NumIn()
	//Debugf("参数总数：", num)
	if num > 0 {
		params := make([]reflect.Value, num)
		start := 1
		if jumpobj.IsValid() && !jumpobj.IsNil() {
			params[0] = jumpobj
			//Debugf("参数0:%v：", jumpobj)
			params[1] = param0
			//Debugf("参数1:%v：", param0)
			start = 2
		} else {
			params[0] = param0
			//Debugf("参数0:%v：", param0)
		}

		if len(ortherparam) > 0 {
			plen := len(ortherparam) + start
			for i := start; i < num && i < plen; i++ {
				pvalue := reflect.ValueOf(ortherparam[i-start])
				//	Warnf("参数%d：%v=====%v ", i, funType.In(i).Kind(), pvalue.Kind())
				ftype := funType.In(i)
				if ftype.Kind() == pvalue.Kind() {

					//	Debugf("参数%d:%v：", i, pvalue)
					params[i] = pvalue
					continue
				}
				if ftype.Kind() == reflect.Interface && pvalue.Type().Implements(ftype) {
					//	Debugf("参数%d:%v：", i, pvalue)
					params[i] = pvalue
					continue
				}

			}
		}
		//Warnf("执行函数：", funValue.Type())
		//Warnf("参数：", params)

		funValue.Call(params)
		return nil
	}
	return fmt.Errorf("注册函数%v参数数量错误", funValue)
}

// ExecFuner 执行者
type ExecFuner interface {
	ExecFun(data []byte, param ...interface{}) error
}
type execStructFun struct {
	key      string
	funValue reflect.Value
	param    reflect.Type
	jumpObj  reflect.Value
}

func (f *execStructFun) ExecFun(data []byte, param ...interface{}) error {
	obj := reflect.New(f.param).Interface()
	json.Unmarshal(data, &obj)
	param0 := reflect.Indirect(reflect.ValueOf(obj))
	return callfun(f.funValue, f.jumpObj, param0, param...)
}

type execStructPrtFun struct {
	key      string
	funValue reflect.Value
	param    reflect.Type
	jumpObj  reflect.Value
}

func (f *execStructPrtFun) ExecFun(data []byte, param ...interface{}) error {
	obj := reflect.New(f.param).Interface()
	json.Unmarshal(data, &obj)
	param0 := reflect.Indirect(reflect.ValueOf(obj))
	return callfun(f.funValue, f.jumpObj, param0, param...)
}

type execMapFun struct {
	key      string
	funValue reflect.Value
	param    reflect.Type
	jumpObj  reflect.Value
}

func (f *execMapFun) ExecFun(data []byte, param ...interface{}) error {
	var obj interface{}

	json.Unmarshal(data, &obj)
	pmap := obj.(map[string]interface{})

	nmap := reflect.MakeMap(f.param)
	switch f.param.Elem().Kind() {
	case reflect.Int:
		for k, v := range pmap {
			nmap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(int(v.(float64))))
		}
	case reflect.Int16:
		for k, v := range pmap {
			nmap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(int16(v.(float64))))
		}
	case reflect.Int32:
		for k, v := range pmap {
			nmap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(int32(v.(float64))))
		}
	case reflect.Int64:
		for k, v := range pmap {
			nmap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(int64(v.(float64))))
		}
	case reflect.Uint:
		for k, v := range pmap {
			nmap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(uint(v.(float64))))
		}
	case reflect.Uint16:
		for k, v := range pmap {
			nmap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(uint16(v.(float64))))
		}
	case reflect.Uint32:
		for k, v := range pmap {
			nmap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(uint32(v.(float64))))
		}
	case reflect.Uint64:
		for k, v := range pmap {
			nmap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(uint64(v.(float64))))
		}
	case reflect.Float32:
		for k, v := range pmap {
			nmap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(float32(v.(float64))))
		}
	default:
		for k, v := range pmap {
			nmap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
		}

	}
	return callfun(f.funValue, f.jumpObj, nmap, param...)
}

type execSliceFun struct {
	key      string
	funValue reflect.Value
	param    reflect.Type
	jumpObj  reflect.Value
}

func (f *execSliceFun) ExecFun(data []byte, param ...interface{}) error {
	slice := reflect.MakeSlice(f.param, 0, 0)
	obj := slice.Interface()
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return fmt.Errorf("解析json失败%v", err)
	}
	if reflect.TypeOf(obj) == f.param {

		return callfun(f.funValue, f.jumpObj, reflect.ValueOf(obj), param...)
	}
	return fmt.Errorf("解析json失败,不是预期的切片类型:%#v %#v", reflect.TypeOf(obj), f.param)
}

// CreateJSONAnalysisHandle 创建json解析器
func CreateJSONAnalysisHandle() JSONAnalysisHandler {
	return &jsonAnalysisHandle{funMap: make(map[string]ExecFuner, 10)}
}

// ErrNoRegKey 未注册的key
type ErrNoRegKey struct {
	Key string
}

// New 创建新的错误
func (e ErrNoRegKey) New(key string) ErrNoRegKey {
	e.Key = key
	return e
}
func (e ErrNoRegKey) Error() string {
	return "未注册的key:" + e.Key
}

// JSONAnalysisHandler json解析处理
type JSONAnalysisHandler interface {
	Exec(key string, data []byte) error
	AddHandle(key string, fun interface{}) bool
	IsValid(fun interface{}) (ExecFuner, bool)
}
type jsonAnalysisHandle struct {
	funMap map[string]ExecFuner
}

func (j *jsonAnalysisHandle) Exec(key string, data []byte) error {
	fun, ok := j.funMap[key]
	if ok {
		return fun.ExecFun(data)
	}
	//Debug("未注册的KEY:", key)
	return ErrNoRegKey{}.New(key)
}
func (j *jsonAnalysisHandle) AddHandle(key string, fun interface{}) bool {
	exefun, ok := j.IsValid(fun)
	if ok {
		j.funMap[key] = exefun
		alog.Debug("注册的KEY:", key, "成功")
		return true
	}
	alog.Debug("注册的KEY:", key, "失败")
	return false
}

func (j *jsonAnalysisHandle) IsValid(fun interface{}) (ExecFuner, bool) {
	if fun == nil {
		return nil, false
	}
	f := reflect.ValueOf(fun)

	if f.Kind() == reflect.Func {
		argnum := f.Type().NumIn()
		if argnum != 1 {
			alog.Warnf("函数参数错误:%#v", f)
			return nil, false
		}
		p := f.Type().In(0)
		kind := p.Kind()
		switch kind {
		case reflect.Struct:
			return &execStructFun{funValue: f, param: p}, true
		case reflect.Map:
			if p.Key().Kind() == reflect.String {
				return &execMapFun{funValue: f, param: p}, true
			}

		case reflect.Ptr:
			switch p.Elem().Kind() {
			case reflect.Struct:
				return &execStructPrtFun{funValue: f, param: p}, true
			case reflect.Slice:
				//if p.Elem().Kind() == reflect.Interface {
				return &execSliceFun{funValue: f, param: p.Elem()}, true
				//}

			}
		case reflect.Slice:
			return &execSliceFun{funValue: f, param: p}, true
		}
	}
	return nil, false
}
