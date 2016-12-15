package jsonroute

import (
	"encoding/json"
	"reflect"
	"strings"

	. "github.com/antlinker/alog"
)

func CreateJsonRoute(routeKey string, jump interface{}) *JsonRoute {
	return &JsonRoute{routeKey: routeKey, funMap: make(map[interface{}]ExecFuner, 16), jumpObj: reflect.ValueOf(jump)}
}

type JsonRoute struct {
	routeKey string
	funMap   map[interface{}]ExecFuner
	jumpObj  reflect.Value
}

func (j *JsonRoute) Exec(data []byte, param ...interface{}) {
	var obj map[string]interface{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		Warn("解码错误:", err)
	}
	value, ok := obj[j.routeKey]
	if !ok {
		ukey := strings.ToUpper(j.routeKey)
		value, ok = obj[ukey]
		if !ok {
			lkey := strings.ToLower(j.routeKey)
			value, ok = obj[lkey]
		}
	}
	if ok {
		fun, ok := j.funMap[value]
		if ok {
			fun.ExecFun(data, param...)
			Debugf("[%s:%v]执行完成", j.routeKey, value)
			return
		}
		Debugf("未注册的%s:%v", j.routeKey, value)
		Debug("所有注册信息:", j.funMap)
		return
	}

	Debug("未注册的"+j.routeKey+":", value)

}
func (j *JsonRoute) GetRouteKey() string {
	return j.routeKey
}
func (j *JsonRoute) SetRouteKey(key string) {
	j.routeKey = key
}
func (j *JsonRoute) AddHandle(value interface{}, fun interface{}) bool {
	exefun, ok := j.IsValid(fun)
	if ok {
		j.funMap[value] = exefun
		Debugf("注册[%s:%v]成功", j.routeKey, value)
		return true
	}
	Debugf("注册[%s:%v]失败", j.routeKey, value)
	return false
}

func (j *JsonRoute) IsValid(fun interface{}) (ExecFuner, bool) {
	if fun == nil {
		return nil, false
	}
	f := reflect.ValueOf(fun)

	if f.Kind() == reflect.Func {
		ftype := f.Type()
		argnum := ftype.NumIn()

		if argnum == 0 {
			Warn(f, "函数参数错误")
			return nil, false
		}
		Debugf("注册函数%v", ftype)
		for n := 0; n < argnum; n++ {

			Debugf("\t参数%d:%v", n, ftype.In(n))
		}

		p := ftype.In(0)
		jumpObj := reflect.ValueOf(nil)

		if j.jumpObj.IsValid() && !j.jumpObj.IsNil() {
			Debugf("跳过对象%v", j.jumpObj.Type())

			if reflect.DeepEqual(j.jumpObj.Type(), p) {
				if argnum == 1 {
					Warn(f, "函数参数错误")
					return nil, false
				}
				p = ftype.In(1)
				jumpObj = j.jumpObj
			}

		}
		kind := p.Kind()
		Debug("p.Kind:", kind)
		switch kind {
		case reflect.Struct:
			return &ExecStructFun{jumpObj: jumpObj, funValue: f, param: p}, true
		case reflect.Map:
			if p.Key().Kind() == reflect.String {
				return &ExecMapFun{jumpObj: jumpObj, funValue: f, param: p}, true
			}

		case reflect.Ptr:
			switch p.Elem().Kind() {
			case reflect.Struct:
				return &ExecStructPrtFun{jumpObj: jumpObj, funValue: f, param: p}, true
			}
		case reflect.Slice:
			return &ExecSliceFun{jumpObj: jumpObj, funValue: f, param: p}, true
		}
	}
	return nil, false
}
