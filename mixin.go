package tango

import (
	"reflect"
)

var mixins = []MixinInterface{}

type MixinInterface interface {
	InitMixin()
}

type BaseMixin struct{}

func (m *BaseMixin) InitMixin() {
	// pass
}

func Mixin(m MixinInterface) {
	//LogInfo.Printf("Mixin INIT: %q", reflect.TypeOf(*m).Name())
	mixins = append(mixins, m)

	// As soon as the Mixin is added, initialize it.
	m.InitMixin()
}

func runMixinPrepare(handler HandlerInterface) {
	// Top to bottom.
	for i := 0; i < len(mixins); i++ {
		m := mixins[i]
		dispatchMixinCall(handler, "Prepare"+reflect.TypeOf(m).Elem().Name())
	}
}

func runMixinFinish(handler HandlerInterface) {
	// Bottom to top.
	for i := len(mixins) - 1; i >= 0; i-- {
		m := mixins[i]
		dispatchMixinCall(handler, "Finish"+reflect.TypeOf(m).Elem().Name())
	}
}

func dispatchMixinCall(handler HandlerInterface, name string) {
	value := reflect.ValueOf(handler)
	valueType := value.Type()
	numMethods := valueType.NumMethod()

	for i := 0; i < numMethods; i++ {
		method := valueType.Method(i)
		if (method.PkgPath == "") && (method.Type.NumIn() == 1) {
			if method.Name == name {
				callArgs := make([]reflect.Value, 1)
				callArgs[0] = value
				method.Func.Call(callArgs)
				return
			}
		}
	}
}
