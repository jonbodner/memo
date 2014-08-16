package memo

import (
	"fmt"
	"reflect"
)

func Memoize(fptr interface{}) error {
	fValue := reflect.ValueOf(fptr)
	//fmt.Println(fptr)
	fActual := fValue.Elem()
	//fmt.Println(fActual)
	if !fActual.CanSet() {
		return fmt.Errorf("Cannot memoize; parameter must be a variable pointer to function")
	}
	fOrig := reflect.ValueOf(fActual.Interface())
	//fmt.Println(fValue.Kind(), fActual.Kind(), fOrig.Kind())
	fType := fActual.Type()
	count := fType.NumIn()
	var m memoHolder
	if count <= 3 {
		m = memoSmall{}
	} else {
		m = memoBig{}
	}
	v := reflect.MakeFunc(fType, func(in []reflect.Value) []reflect.Value {
		//check to see if we have seen this before
		key := m.buildKey(in)
		//fmt.Println("using key", key)
		if val, ok := m.hasVal(key); ok {
			//fmt.Printf("returning memoed value %v\n",val)
			return val
		}
		result := fOrig.Call(in)
		m.storeVal(key,result)
		//fmt.Printf("stored memoed value %v\n", result)
		return result
	})
	fActual.Set(v)
	return nil
}

type memoKey interface{}

type memoHolder interface {
	buildKey([]reflect.Value) memoKey
	hasVal(memoKey) ([]reflect.Value, bool)
	storeVal(memoKey, []reflect.Value)
}

type memoSmall map[[3]reflect.Value][]reflect.Value

func (ms memoSmall)buildKey(in []reflect.Value) memoKey {
	key := [3]reflect.Value{}
	copy(key[:],in)
	return memoKey(key)
}

func(ms memoSmall)hasVal(key memoKey) ([]reflect.Value, bool) {
	v, ok := ms[key.([3]reflect.Value)]
	return v, ok
}

func (ms memoSmall)storeVal(key memoKey, result []reflect.Value) {
	ms[key.([3]reflect.Value)] = result
}

type memoBig map[string][]reflect.Value

func (mb memoBig)buildKey(in []reflect.Value) memoKey {
        key := ""
        for _, val := range in {
            key += fmt.Sprintf("%v-", val.Interface())
        }
	return memoKey(key)
}

func (mb memoBig)hasVal(key memoKey) ([]reflect.Value, bool) {
	v, ok := mb[key.(string)]
	return v, ok
}

func (mb memoBig)storeVal(key memoKey, result []reflect.Value) {
	mb[key.(string)]=result
}

