package jsonroute

import (
	"encoding/json"
	"reflect"

	. "github.com/antlinker/alog"
)

func init() {
	RegisterAlog()
}

var unJsonFun = reflect.ValueOf(json.Unmarshal)

func callfun(funValue reflect.Value, jumpobj reflect.Value, param0 reflect.Value, ortherparam ...interface{}) {
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
		return
	}
	Warnf("注册函数%v参数数量错误", funValue)
}

type ExecFuner interface {
	ExecFun(data []byte, param ...interface{})
}
type ExecStructFun struct {
	key      string
	funValue reflect.Value
	param    reflect.Type
	jumpObj  reflect.Value
}

func (f *ExecStructFun) ExecFun(data []byte, param ...interface{}) {
	obj := reflect.New(f.param).Interface()
	json.Unmarshal(data, &obj)
	param0 := reflect.Indirect(reflect.ValueOf(obj))
	callfun(f.funValue, f.jumpObj, param0, param...)
}

type ExecStructPrtFun struct {
	key      string
	funValue reflect.Value
	param    reflect.Type
	jumpObj  reflect.Value
}

func (f *ExecStructPrtFun) ExecFun(data []byte, param ...interface{}) {
	obj := reflect.New(f.param).Interface()
	json.Unmarshal(data, &obj)
	param0 := reflect.Indirect(reflect.ValueOf(obj))
	callfun(f.funValue, f.jumpObj, param0, param...)
}

type ExecMapFun struct {
	key      string
	funValue reflect.Value
	param    reflect.Type
	jumpObj  reflect.Value
}

func (f *ExecMapFun) ExecFun(data []byte, param ...interface{}) {
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
	callfun(f.funValue, f.jumpObj, nmap, param...)
}

type ExecSliceFun struct {
	key      string
	funValue reflect.Value
	param    reflect.Type
	jumpObj  reflect.Value
}

func (f *ExecSliceFun) ExecFun(data []byte, param ...interface{}) {
	slice := reflect.MakeSlice(f.param, 0, 0)
	obj := slice.Interface()
	json.Unmarshal(data, &obj)
	callfun(f.funValue, f.jumpObj, reflect.ValueOf(obj), param...)
}
func CreateJsonAnalysisHandle() *JsonAnalysisHandle {
	return &JsonAnalysisHandle{funMap: make(map[string]ExecFuner, 10)}
}

type JsonAnalysisHandle struct {
	funMap map[string]ExecFuner
}

func (j *JsonAnalysisHandle) Exec(key string, data []byte) {
	fun, ok := j.funMap[key]
	if ok {
		fun.ExecFun(data)
		Debug(key, "执行完成")
		return
	}
	Debug("未注册的KEY:", key)
}
func (j *JsonAnalysisHandle) AddHandle(key string, fun interface{}) bool {
	exefun, ok := j.IsValid(fun)
	if ok {
		j.funMap[key] = exefun
		Debug("注册的KEY:", key, "成功")
		return true
	}
	Debug("注册的KEY:", key, "失败")
	return false
}

func (j *JsonAnalysisHandle) IsValid(fun interface{}) (ExecFuner, bool) {
	if fun == nil {
		return nil, false
	}
	f := reflect.ValueOf(fun)

	if f.Kind() == reflect.Func {
		argnum := f.Type().NumIn()
		if argnum != 1 {
			Warn(f, "函数参数错误")
			return nil, false
		}
		p := f.Type().In(0)
		kind := p.Kind()
		switch kind {
		case reflect.Struct:
			return &ExecStructFun{funValue: f, param: p}, true
		case reflect.Map:
			if p.Key().Kind() == reflect.String {
				return &ExecMapFun{funValue: f, param: p}, true
			}

		case reflect.Ptr:
			switch p.Elem().Kind() {
			case reflect.Struct:
				return &ExecStructPrtFun{funValue: f, param: p}, true
			case reflect.Slice:
				//if p.Elem().Kind() == reflect.Interface {
				return &ExecSliceFun{funValue: f, param: p.Elem()}, true
				//}

			}
		case reflect.Slice:
			return &ExecSliceFun{funValue: f, param: p}, true
		}
	}
	return nil, false
}
