package jsonroute

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/antlinker/alog"
)

// CreateJSONRoute 创建json路由
func CreateJSONRoute(routeKey string, jump interface{}) JSONRouter {
	return &jsonRoute{routeKey: routeKey, funMap: make(map[interface{}]ExecFuner, 16), jumpObj: reflect.ValueOf(jump)}
}

// JSONRouter json路由
type JSONRouter interface {
	Exec(data []byte, param ...interface{}) error
	GetRouteKey() string
	SetRouteKey(key string)
	AddHandle(value interface{}, fun interface{}) bool
}

type jsonRoute struct {
	routeKey string
	funMap   map[interface{}]ExecFuner
	jumpObj  reflect.Value
}

func (j *jsonRoute) Exec(data []byte, param ...interface{}) error {
	var obj map[string]interface{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return fmt.Errorf("json解码错误:%v", err)
	}
	value, ok := obj[j.routeKey]

	if ok {
		v := fmt.Sprintf("%v", value)
		fun, ok := j.funMap[v]
		if !ok {
			ukey := strings.ToUpper(v)
			fun, ok = j.funMap[ukey]
			if !ok {
				lkey := strings.ToLower(v)
				fun, ok = j.funMap[lkey]
			}
		}
		if ok {
			//	alog.Debugf("[%s:%v]执行完成", j.routeKey, param)
			return fun.ExecFun(data, param...)
			//	Debugf("[%s:%v]执行完成", j.routeKey, value)

		}
		return fmt.Errorf("未注册的%s:%v", j.routeKey, value)
	}
	return fmt.Errorf("报文格式错误：缺少%s字段", j.routeKey)
}
func (j *jsonRoute) GetRouteKey() string {
	return j.routeKey
}
func (j *jsonRoute) SetRouteKey(key string) {
	j.routeKey = key
}
func (j *jsonRoute) AddHandle(value interface{}, fun interface{}) bool {
	exefun, ok := j.IsValid(fun)
	if ok {
		j.funMap[value] = exefun
		alog.Debugf("注册[%s:%v]成功", j.routeKey, value)
		return true
	}
	alog.Debugf("注册[%s:%v]失败", j.routeKey, value)
	return false
}

func (j *jsonRoute) IsValid(fun interface{}) (ExecFuner, bool) {
	if fun == nil {
		return nil, false
	}
	f := reflect.ValueOf(fun)

	if f.Kind() == reflect.Func {
		ftype := f.Type()
		argnum := ftype.NumIn()

		if argnum == 0 {
			alog.Warn(f, "函数参数错误")
			return nil, false
		}
		alog.Debugf("注册函数%v", ftype)
		for n := 0; n < argnum; n++ {

			alog.Debugf("\t参数%d:%v", n, ftype.In(n))
		}

		p := ftype.In(0)
		jumpObj := reflect.ValueOf(nil)
		if j.jumpObj.IsValid() {

			if reflect.DeepEqual(j.jumpObj.Type(), p) {
				if argnum == 1 {
					alog.Warn(f, "函数参数错误")
					return nil, false
				}
				p = ftype.In(1)
				jumpObj = j.jumpObj
			}

		}
		kind := p.Kind()
		switch kind {
		case reflect.Struct:
			return &execStructFun{jumpObj: jumpObj, funValue: f, param: p}, true
		case reflect.Map:
			if p.Key().Kind() == reflect.String {
				return &execMapFun{jumpObj: jumpObj, funValue: f, param: p}, true
			}

		case reflect.Ptr:
			switch p.Elem().Kind() {
			case reflect.Struct:
				return &execStructPrtFun{jumpObj: jumpObj, funValue: f, param: p}, true
			}
		case reflect.Slice:
			return &execSliceFun{jumpObj: jumpObj, funValue: f, param: p}, true
		}
	}
	return nil, false
}
