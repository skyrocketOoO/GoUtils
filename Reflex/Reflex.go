package Reflex

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/oleiade/reflections"
)

const (
	CurrentCaller int = 1
	ParentCaller      = CurrentCaller + 1
)

func IsNil(it reflect.Value) (ok bool) {
	defer func() {
		throw := recover()
		if nil != throw {
			ok = false
		}
	}()

	if false == it.IsValid() {
		ok = true
		return
	}

	ok = it.IsNil()
	return
}

func ValueOf[Self any](it Self, deep bool) (out reflect.Value) {
	out = reflect.ValueOf(it)

	if deep {
		for reflect.Pointer == out.Kind() {
			out = out.Elem()
		}

		return
	}

	if reflect.Pointer == out.Kind() {
		out = out.Elem()
	}

	return
}

func GetField[Out any, Self any](it *Self, name string) (out *Out, throw error) {
	value, throw := reflections.GetField(it, name)
	if nil != throw {
		return
	}

	item, ok := value.(Out)
	if false == ok {
		throw = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"assert field type fail: %s",
			name,
		)
	}

	out = &item
	return
}

func CallField[Self any](
	it *Self,
	name string,
	arguments []reflect.Value,
) (out []reflect.Value, throw error) {
	value := ValueOf(it, false)

	_, ok := value.Type().MethodByName(name)
	if false == ok {
		throw = gerror.NewCodef(gcode.CodeInvalidParameter, "method is not exist: %s", name)
		return
	}

	method := value.MethodByName(name)

	count := len(arguments)
	if method.Type().NumIn() != count {
		throw = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"the method arguments is out of index: %s",
			method.String(),
		)
		return
	}

	out = method.Call(arguments)
	return
}

func GetMap[Out any, In any](it *In) (out map[string]Out, throw error) {
	anyOut, throw := reflections.ItemsDeep(it)
	if nil != throw {
		return
	}

	data := make(map[string]Out)
	for key, value := range anyOut {
		item, ok := value.(Out)
		if false == ok {
			throw = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				"assert field type fail: %s",
				key,
			)
			return
		}

		data[key] = item
	}

	out = data
	return
}

func GetFunctionName[Self any](data Self, shortName bool) (out string) {
	name := runtime.FuncForPC(reflect.ValueOf(data).Pointer()).Name()

	if shortName {
		items := strings.Split(name, ".")
		items = strings.Split(items[len(items)-1], "-")
		name = items[0]
	}

	out = name
	return
}

func GetCallerName(skip int, shortName bool) (out string) {
	out = ""
	pointer, _, _, ok := runtime.Caller(skip)
	if false == ok {
		return
	}

	name := runtime.FuncForPC(pointer).Name()

	if shortName {
		items := strings.Split(name, ".")
		name = items[len(items)-1]
	}

	out = name
	return
}

func GetCurrentCallerShortName() (out string) {
	out = GetCallerName(CurrentCaller+1, true)
	return
}

func GetCurrentCallerFullName() (out string) {
	out = GetCallerName(CurrentCaller+1, false)
	return
}
