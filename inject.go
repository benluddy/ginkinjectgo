package ginkinjectgo

import (
	"reflect"
)

type env struct {
	parent    *env
	providers map[reflect.Type]provider
}

var global env

type provider reflect.Value

func (p provider) Provide() reflect.Value {
	return reflect.Value(p).Call(nil)[0]
}

func (e *env) Env() *env {
	return &env{parent: e}
}

func (e *env) GetProvider(t reflect.Type) provider {
	if p, ok := e.providers[t]; ok {
		return p
	}
	if e.parent != nil {
		return e.parent.GetProvider(t)
	}
	pt := reflect.FuncOf(nil, []reflect.Type{t}, false)
	return provider(reflect.MakeFunc(pt, func(args []reflect.Value) []reflect.Value {
		return []reflect.Value{reflect.Zero(t)}
	}))
}

func RegisterProvider(p interface{}) {
	curEnv.RegisterProvider(p)
}

func (e *env) RegisterProvider(p interface{}) {
	tp := reflect.TypeOf(p)
	if tp.Kind() != reflect.Func {
		if e.providers == nil {
			e.providers = make(map[reflect.Type]provider)
		}
		tfn := reflect.FuncOf(nil, []reflect.Type{tp}, false)
		e.providers[tp] = provider(reflect.MakeFunc(tfn, func(args []reflect.Value) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(p)}
		}))
		return
	}
	if tp.NumIn() != 0 {
		panic("provider takes params")
	}
	if tp.NumOut() != 1 {
		panic("provider does not return a single value")
	}
	to := tp.Out(0)
	if e.providers == nil {
		e.providers = make(map[reflect.Type]provider)
	}
	e.providers[to] = provider(reflect.ValueOf(p))
}

func (e *env) Inject(body interface{}) func() {
	return func() {
		tb := reflect.TypeOf(body)
		if tb.Kind() != reflect.Func {
			panic("body not func")
		}
		args := make([]reflect.Value, tb.NumIn())
		for i := range args {
			ti := tb.In(i)
			args[i] = e.GetProvider(ti).Provide()
		}
		reflect.ValueOf(body).Call(args)
	}
}
